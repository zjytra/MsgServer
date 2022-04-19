package network

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
	"net"
	"sync"
	"time"
)

var (
	colseErro = errors.New("已经关闭连接")
)

// 每個对象维持一个连接,连接其他服务器使用
type ServiceClient struct {
	sync.Mutex
	Addr            string        // 服务器连接地址
	ConnectInterval time.Duration // 重连时间
	ConnNum         int           // 连接数
	AutoReconnect   bool
	svNet           ServiceNetEvent // 服务器事件观察者
	conn            *ServiceConn    // 连接对象
	appID           int32           // 保存appID
	wg              sync.WaitGroup
	closeFlag       *model.AtomicBool
	msgParser       *MsgParser // 数据包解析对象 共用一个对象

}

// 创建tcp 客戶端
func NewServiceClient(svNet ServiceNetEvent, appID int32) *ServiceClient {
	cfg := csvdata.GetNetWorkCfgPtr(appID)
	if cfg == nil {
		return nil
	}
	if svNet == nil {
		xlog.Warning("服务器 消息处理 svNet is nil")
		return nil
	}
	client := new(ServiceClient)
	client.svNet = svNet
	client.appID = appID
	client.closeFlag = model.NewAtomicBool()
	client.msgParser = NewMsgParser(cfg.Max_msglen, cfg.Msg_isencrypt)
	return client
}

func (client *ServiceClient) Start() {
	client.init()
	client.wg.Add(1)
	go client.connect() // 开启一个连接
}

func (client *ServiceClient) init() {
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 5 * time.Second
		xlog.Debug("invalid ConnectInterval, reset to %v", client.ConnectInterval)
	}
	if client.ConnNum <= 0 {
		client.ConnNum = 1
		xlog.Debug("invalid ConnNum, reset to %v", client.ConnNum)
	}
	client.closeFlag.SetFalse()
	client.AutoReconnect = true
	cfg := csvdata.GetNetWorkCfgPtr(client.appID)
	if cfg == nil {
		xlog.Warning("server %d conf is nil", client.appID)
		return
	}
	//服务器使用内网地址
	client.Addr = fmt.Sprintf("%s:%d", cfg.Inner_addr, cfg.Out_prot)
	xlog.Debug("connet to %v ", client.Addr)
}

func (client *ServiceClient) dial() net.Conn {
	for {
		if client.closeFlag.IsTrue() {
			return nil
		}
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil {
			return conn
		}

		xlog.Debug("ServiceClient dial to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *ServiceClient) connect() {
	defer client.wg.Done()
redial:
	client.doconnect() // 执行连接及读
	// 没有关闭才进行重连
	if client.closeFlag.IsFalse() && client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto redial
	}
}

func (client *ServiceClient) doconnect() bool {
	conn := client.dial()
	if conn == nil {
		return false
	}
	if !client.setConn(conn) {
		return false
	}
	return true
}

// 添加链接信息
func (client *ServiceClient) setConn(conn net.Conn) bool {
	if client.closeFlag.IsTrue() {
		conn.Close()
		return false
	}
	svconn := newServiceConn(conn, nextID(), client.svNet, client.appID, client.msgParser)
	client.conn = svconn
	xlog.Debug("连接远程 %v 地址成功", conn.RemoteAddr())
	// 连接成功,将阻塞读取数据
	client.ReceiveData(client.conn)
	xlog.Debug("ServeiceClient结束读取")
	return true
}

// 连接中读取数据
func (client *ServiceClient) ReceiveData(svcon *ServiceConn) {
	for {
		err := svcon.ReadMsg()
		if err != nil { // 这里读到错误消息,关闭
			xlog.Debug("read message: ", err)
			break // 关闭连接
		}
	}
	// cleanup
	client.closeConn(svcon)
}

func (client *ServiceClient) closeConn(conn *ServiceConn) {
	conn.Close()
	xlog.Debug("关闭远程服务器连接")
}

// 写单个消息
func (client *ServiceClient) WriteOneMsg(maincmd uint32, msg []byte) {
	if client.closeFlag.IsTrue() {
		xlog.Debug("WriteOneMsg  命令 %d %v", maincmd, colseErro)
		return
	}
	if client.conn == nil || client.conn.IsClose() {
		xlog.Debug("ServiceClient WriteOneMsg未建立连接 %v", client.Addr)
		return
	}
	client.conn.WriteOneMsg(maincmd, msg)
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func (client *ServiceClient) GetOneMsgByteArr(maincmd uint32, msg []byte) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.conn == nil {
		return nil, errors.New(fmt.Sprintf("GetOneMsgByteArr未建立连接 %v", client.Addr))
	}
	return client.conn.GetOneMsgByteArr(maincmd, msg)
}

// 写单个消息pb实现
func (client *ServiceClient) WritePBMsg(maincmd uint32, pb proto.Message)  {
	if client.closeFlag.IsTrue() {
		xlog.Debug("ServiceClient WriteMsg  cmd %d err  %v",maincmd, colseErro)
		return
	}
	if client.conn == nil || client.IsClose() {
		xlog.Debug("ServiceClient WriteMsg  未建立连接  cmd %d 地址  %v",maincmd, client.Addr)
		return
	}
	client.conn.WritePBMsg(maincmd, pb)
}

// 将消息体构建为[]byte数组，最终要发出去的单包 pb实现
func (client *ServiceClient) GetPBByteArr(maincmd uint32, pb proto.Message) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.conn == nil {
		return nil, errors.New(fmt.Sprintf("GetOneMsgByteArr未建立连接 %v", client.Addr))
	}
	return client.conn.GetPBByteArr(maincmd, pb)
}

// 一起写多个数据包
// 每个包的数据 由GetOneMsgByteArr构建
func (client *ServiceClient) WriteMsg(args ...[]byte) {
	if client.closeFlag.IsTrue() {
		xlog.Debug("ServiceClient WriteMsg   %v", colseErro)
		return
	}
	client.conn.WriteMsg(args...)
}

// 是否存活 没有存活会被提下线
func (client *ServiceClient) IsAlive() bool {
	return client.conn.IsAlive()
}

// 是否关闭
func (client *ServiceClient) IsClose() bool {
	return client.conn.IsClose()
}

// 获取连接对象id
func (client *ServiceClient) GetConnID() uint32 {
	return client.conn.GetConnID()
}

// 关闭服务
func (client *ServiceClient) Close() {
	client.closeFlag.SetTrue()
	client.conn.Close()
	client.wg.Wait()
}

// 关闭连接
func (client *ServiceClient) DoCloseConn() {
	client.conn.Close()
}

func (client *ServiceClient) GetServiceConn() *ServiceConn {
	return client.conn
}

func (client *ServiceClient) GetAppId() int32 {
	return client.appID
}
