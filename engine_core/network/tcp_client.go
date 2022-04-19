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

// 每個对象维持ConnNum个连接,这个当测试客户端连接
type TCPClient struct {
	sync.RWMutex
	Addr            string        // 服务器连接地址
	ConnectInterval time.Duration // 重连时间
	ConnNum         int           // 连接数量
	AutoReconnect   bool
	clnEv           ClientNetEvent // 客户端事件观察者
	tcpConn         Conner       // 连接对象
	useConnIdex     int            // 使用的连接下标
	wg              sync.WaitGroup
	closeFlag       *model.AtomicBool
	appID           int32      // 保存appID
	msgParser       *MsgParser // 数据包解析对象 共用一个对象
}

// 创建tcp 客戶端
func NewTCPClient(clnEv ClientNetEvent, appID int32) *TCPClient {
	cfg := csvdata.GetNetWorkCfgPtr(appID)
	if cfg == nil {
		xlog.Error("appid %v 配置为nil", appID)
		return nil
	}
	if clnEv == nil {
		xlog.Warning("服务器 消息处理 svNet is nil")
		return nil
	}
	client := new(TCPClient)
	client.appID = appID
	client.clnEv = clnEv
	client.closeFlag = model.NewAtomicBool()
	client.msgParser = NewMsgParser(cfg.Max_msglen, cfg.Msg_isencrypt)
	return client
}

func (client *TCPClient) Start() {
	client.init()
	client.wg.Add(1)
	go client.connect() // 开启一个连接

}

func (client *TCPClient) init() {
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
		xlog.Error("appid %v 配置为nil", client.appID)
		return
	}
	client.Addr = fmt.Sprintf("%s:%d", cfg.Out_addr, cfg.Out_prot)
	xlog.Debug("client.Addr connet %v ", client.Addr)
}

func (client *TCPClient) dial() net.Conn {
	for {
		if client.closeFlag.IsTrue() {
			return nil
		}
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil {
			return conn
		}

		xlog.Debug("TCPClient dial to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {
	defer client.wg.Done()
redial:
	conn := client.dial()
	if conn == nil {
		return
	}
	if !client.setConn(conn) {
		return
	} // 执行连接及读
	// 没有关闭才进行重连
	if client.closeFlag.IsFalse() && client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto redial
	}
}

// 添加链接信息
func (client *TCPClient) setConn(conn net.Conn) bool {
	client.Lock()
	if client.closeFlag.IsTrue() {
		client.Unlock()
		conn.Close()
		return false
	}
	cclient := newClientConn(conn, nextID(), client.clnEv, client.appID, client.msgParser)
	if cclient != nil {
		client.tcpConn = cclient
		cclient.notifyConnect(client.tcpConn)
	}
	client.Unlock()

	xlog.Debug("连接远程 %v 地址成功", conn.RemoteAddr())
	// 连接成功,将阻塞读取数据
	client.ReceiveData(client.tcpConn)
	xlog.Debug("TCPClient结束读取")
	return true
}

// 连接中读取数据
func (client *TCPClient) ReceiveData(conn Conner) {
	for {
		err := conn.ReadMsg()
		if err != nil { // 这里读到错误消息,关闭
			xlog.Warning("read message: %v ", err)
			break // 关闭连接
		}
	}
	// cleanup
	client.closeConn(conn)
}

func (client *TCPClient) closeConn(conn Conner) {
	client.Lock()
	conn.Close()
	client.Unlock()
	xlog.Debug("关闭远程连接")
}

// 写单个消息
func (client *TCPClient) WriteOneMsg(maincmd uint32, msg []byte) {
	if client.closeFlag.IsTrue() {
		xlog.Debug("TCPClient WriteOneMsg %v", colseErro)
		return
	}
	if client.tcpConn == nil {
		xlog.Debug("TCPClient WriteOneMsg 未建立连接 %v", client.Addr)
		return
	}
	client.tcpConn.WriteOneMsg(maincmd, msg)
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func (client *TCPClient) GetOneMsgByteArr(maincmd uint32, msg []byte) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.tcpConn == nil {
		return nil, errors.New(fmt.Sprintf("GetOneMsgByteArr未建立连接 %v", client.Addr))
	}
	return client.tcpConn.GetOneMsgByteArr(maincmd, msg)
}

// 写单个消息pb实现
func (client *TCPClient) WritePBMsg(maincmd uint32, pb proto.Message)  {
	if client.closeFlag.IsTrue() {
		xlog.Debug("TCPClient WritePBMsg cmd = %d  错误 %v",maincmd, colseErro)
		return
	}
	if client.tcpConn == nil || client.tcpConn.IsClose() {
		xlog.Debug("TCPClient WritePBMsg cmd = %d 未建立连接 地址 %v",maincmd, client.Addr)
		return
	}
	client.tcpConn.WritePBMsg(maincmd, pb)
}

// 将消息体构建为[]byte数组，最终要发出去的单包 pb实现
func (client *TCPClient) GetPBByteArr(maincmd uint32, pb proto.Message) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.tcpConn == nil {
		return nil, errors.New(fmt.Sprintf("GetPBByteArr未建立连接 %v", client.Addr))
	}
	return client.tcpConn.GetPBByteArr(maincmd, pb)
}

// 一起写多个数据包
// 每个包的数据 由GetOneMsgByteArr构建
func (client *TCPClient) WriteMsg(args ...[]byte) {
	if client.closeFlag.IsTrue() {
		xlog.Debug("TCPClient WriteOneMsg %v", colseErro)
		return
	}
	if client.tcpConn == nil {
		xlog.Debug("TCPClient WriteOneMsg 未建立连接 %v", client.Addr)
		return
	}
	client.tcpConn.WriteMsg(args...)
}

// 是否存活 没有存活会被提下线
func (client *TCPClient) IsAlive() bool {
	return client.tcpConn.IsAlive()
}

// 是否关闭
func (client *TCPClient) IsClose() (isAllClose bool) {
	return client.tcpConn.IsClose()
}

//关闭服务
func (client *TCPClient) Close() {
	client.closeFlag.SetTrue()
	client.tcpConn.Close()
	client.wg.Wait()
}
