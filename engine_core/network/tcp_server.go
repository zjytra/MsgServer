/*
创建时间: 2020/7/14
作者: zjy
功能介绍:
处理客户端的服务器
*/
package network

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/timingwheel"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	ln         net.Listener                // 服务器监听对象
	connetsize *model.AtomicInt32FlagModel // 已经连接的数量
	connMaps   map[uint32]Conner           // 客户端连接对象map
	mutexConns sync.RWMutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup
	clnEv      ClientNetEvent // 网络事件观察者
	linkTicker *timingwheel.Timer
	writeChan  chan *GroupMessage // 写的通道，我服务器写的消息先写入通道再用连接传出去
	msgPool    sync.Pool
	isClose    *model.AtomicBool //是否关闭
	msgParser  *MsgParser        //消息解析
	appID      int32             // 保存appID
}

// 创建tcp Sever服务器
func NewTcpServer(netobs ClientNetEvent, appID int32) *TCPServer {
	tcpsv := new(TCPServer)
	tcpsv.clnEv = netobs
	err := tcpsv.initServerData(appID)
	if err != nil {
		xlog.Error("newClientConn %v", err)
		return nil
	}
	return tcpsv
}

func (server *TCPServer) initServerData(appID int32) error {
	server.appID = appID
	cfg := server.GetNetCfg()
	if cfg == nil {
		return errors.New("网络配置数据为nil appID: " + string(appID))
	}
	server.connetsize = model.NewAtomicInt32Flag()
	server.isClose = model.NewAtomicBool()
	server.isClose.SetFalse()

	server.writeChan = make(chan *GroupMessage, cfg.Write_cap_num)
	server.msgParser = NewMsgParser(cfg.Max_msglen, cfg.Msg_isencrypt)
	server.msgPool.New = func() interface{} {
		return new(GroupMessage)
	}
	return nil
}

func (server *TCPServer) Start() {
	xlog.Debug("TCPServer start")
	server.init()
	go server.run()
	go server.serverEvent()
	////检测客户端是否是活的
	//if !server.isInnerServer() {
	//	server.linkTicker = timersys.NewTimeTicker(time.Second, server.checkLink,server.qe)
	//}
}

//服务器内部通信不需要检测
func (server *TCPServer) isInnerServer() bool {
	cfg := server.GetNetCfg()
	if cfg == nil {
		return true
	}
	return cfg.App_kind == model.APP_DataCenter || cfg.App_kind == model.APP_GameServer
}

func (server *TCPServer) init() {
	cfg := server.GetNetCfg()
	if cfg == nil {
		return
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Out_prot))
	if err != nil {
		xlog.Debug("%v", err)
		panic(err)
	}
	xlog.Debug("TCPServer listen Addr:%v", ln.Addr())

	if cfg.Max_connect <= 0 {
		xlog.Warning("invalid MaxConnNum, reset to %v", cfg.Max_connect)
	}
	if cfg.Write_cap_num <= 0 {
		xlog.Warning("invalid ConnNum, reset to %v", cfg.Write_cap_num)
	}
	// if server.NewAgent == nil {
	// 	xlog.WarningLog(appdata.GetSceneName(),"NewAgent must not be nil")
	// }
	server.ln = ln
	server.connMaps = make(map[uint32]Conner)
}

func (server *TCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()
	defer xlog.RecoverToLog(func() {
		// 出錯要关闭远端连接
		go server.run()
	})
	xlog.Debug("TCPServer Accept Addr:%v", server.ln.Addr())
	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			// 临时错误才继续，其他错误就关闭监听
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				xlog.Warning("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			xlog.Warning("TCPServer Accept erro:%v", err)
			return
		}
		tempDelay = 0
		// 添加连接
		if !server.addConn(conn) {
			continue
		}
	}
}

// 添加链接信息
func (server *TCPServer) addConn(conn net.Conn) bool {

	server.mutexConns.Lock() //互斥锁
	maxNum := int32(10000)
	cfg := server.GetNetCfg()
	if cfg != nil {
		maxNum = cfg.Max_connect
	}
	if server.GetConnectSize() >= maxNum {
		server.mutexConns.Unlock()
		erro := conn.Close()
		if erro != nil {
			xlog.Warning("超过连接关闭链接错误 %v ", erro)
		}
		xlog.Warning("超过最大链接数,当前连接数%d", server.connetsize)
		return false
	}
	// 创建封装的连接
	ConnID := nextID()
	clConn := newClientConn(conn, ConnID, server.clnEv, server.appID, server.msgParser)
	if clConn != nil {
		clConn.notifyConnect(clConn)
	}
	//tcpConn := newTcpConn(conn, ConnID, server.clnEv, server.appID, server.msgParser)
	server.connMaps[ConnID] = clConn // 存储连接
	server.mutexConns.Unlock()        //解锁

	server.connetsize.AddInt32()
	server.wgConns.Add(1)

	go server.ReceiveData(clConn)

	xlog.Debug("当前连接数addConn %d,连接标识%d", server.GetConnectSize(), ConnID)
	return true
}

// 连接中读取数据
func (server *TCPServer) ReceiveData(tcpConn Conner) {
	defer server.wgConns.Done() // 关闭连接waitgorup减一

	defer xlog.RecoverToLog(func() {
		// 出錯要关闭远端连接
		server.closeConn(tcpConn)
	})
	for {
		err := tcpConn.ReadMsg()
		if err != nil { // 这里读到错误消息,关闭
			xlog.Debug("read message err: %v ", err)
			break // 关闭连接
		}
	}
	// 处理远端关闭
	server.closeConn(tcpConn)
}

// 根据连接id断开连接
func (server *TCPServer) CloseConnID(connID uint32) {
	conn := server.GetTcpConnect(connID)
	if conn != nil { //让写协程关闭,这样流程才正确 保证closeConn只被调用一次
		conn.Close()
	}
}

// 被动断开连接 read EOF
func (server *TCPServer) closeConn(tcpConn Conner) {
	if tcpConn == nil {
		return
	}

	tcpConn.Close() // 关闭写协程
	id := tcpConn.GetConnID()
	server.mutexConns.Lock() //互斥锁
	delete(server.connMaps, id)
	server.mutexConns.Unlock() //互斥锁

	server.connetsize.SubInt32()
	xlog.Debug("被动断开连接 当前连接数closeConn %d,移除连接标识%d", server.GetConnectSize(), id)
}

// 获取连接数
func (server *TCPServer) GetConnectSize() int32 {
	return server.connetsize.GetInt32()
}

//监测连接数
func (server *TCPServer) checkLink() error {
	//xlog.Warning("TCPServercheckLink")
	// 关闭所有连接
	server.mutexConns.Lock()
	for _, conn := range server.connMaps {
		if !conn.IsAlive() {
			sec := 50
			cfg := server.GetNetCfg()
			if cfg != nil {
				sec = cfg.Checklink_s
			}
			xlog.Error("connid %d 已经超过 %d 秒没发包", conn.GetConnID(), sec)
			conn.Close()
		}
	}
	server.mutexConns.Unlock()
	return nil
}

// 获取连接数
func (server *TCPServer) serverEvent() {

	for groupmsg := range server.writeChan {
		if groupmsg == nil {
			close(server.writeChan)
			break //退出协程
		}
		server.doSend(groupmsg)
	}
	// xlog.Debug("tcp serverEvent")
}

func (server *TCPServer) Close() {
	server.isClose.SetTrue()
	if server.linkTicker != nil { //只有开启了才关闭
		timingwheel.StopTimer(server.linkTicker)
	}
	server.writeChan <- nil // 关闭服务器事件

	erro := server.ln.Close() // 关闭监听
	if erro != nil {
		xlog.Warning("TCPServer关闭监听错误 %v", erro)
	}
	server.wgLn.Wait()

	// 关闭所有连接
	for _, conn := range server.connMaps {
		conn.Close()
	}
	server.connMaps = nil
	server.wgConns.Wait()
	fmt.Println("TCPServer doClose")
}

// 发送多个消息
func (server *TCPServer) WriteMoreMsgByConnID(ConnID uint32, msg ...[]byte) {
	tcpconn := server.GetTcpConnect(ConnID)
	if tcpconn == nil {
		xlog.Debug("SendMsg 未找到连接 %v", ConnID)
		return
	}
	tcpconn.WriteMsg(msg...)

}

// 写单个消息
func (server *TCPServer) WriteOneMsgByConn(conn Conner, maincmd uint32, msg []byte) {
	server.WriteOneMsgByConnID(conn.GetConnID(), maincmd, msg)
}

// 写单个消息
func (server *TCPServer) WriteOneMsgByConnID(ConnID uint32, maincmd uint32, msg []byte) {
	tcpconn := server.GetTcpConnect(ConnID)
	if tcpconn == nil {
		xlog.Debug("SendMsg 未找到连接 %v", ConnID)
		return
	}
	tcpconn.WriteOneMsg(maincmd, msg)
}

// 用protubuf的方式写单个消息
func (server *TCPServer) WritePBMsgByConn(conn Conner, maincmd uint32, pb proto.Message) {
	ConnID := conn.GetConnID()
	tcpconn := server.GetTcpConnect(ConnID)
	if tcpconn == nil {
		xlog.Debug("TCPServer WritePBMsgByConn cmd %d 未找到连接 %v", maincmd, ConnID)
		return
	}
	tcpconn.WritePBMsg(maincmd, pb)
}

//  用protubuf的方式写单个消息
func (server *TCPServer) WritePBMsgByConnID(ConnID uint32, maincmd uint32, pb proto.Message) {
	tcpconn := server.GetTcpConnect(ConnID)
	if tcpconn == nil {
		xlog.Debug("TCPServer WritePBMsgByConnID cmd %d 未找到连接 %v", maincmd, ConnID)
		return
	}
	tcpconn.WritePBMsg(maincmd, pb)
}

// 根据命令及protobuf创建包
func (server *TCPServer) CreatePBMsg(maincmd uint32, pb proto.Message) (sendMsg []byte, erro error) {
	if pb != nil {
		sendMsg, erro = proto.Marshal(pb)
	}
	if erro != nil {
		xlog.ErrorLog(appdata.GetSceneName(), "CreatePBMsg %v", erro)
		return nil, erro
	}
	sendMsg, erro = server.CreatePackage(maincmd, sendMsg)
	return
}

// 根据命令创建包
func (server *TCPServer) CreatePackage(maincmd uint32, msg []byte) ([]byte, error) {
	return server.msgParser.PackOne(maincmd, msg)
}

// 将多个包合并成一个
func (server *TCPServer) MorePackageToOne(args ...[]byte) ([]byte, error) {
	return server.msgParser.MorePackageToOne(args...)
}

// 给所有连接发送消息
func (server *TCPServer) SendAllConn(msg []byte) {
	if server.isClose.IsTrue() {
		return
	}
	groupmsg := server.createGroupMessage(nil, msg)
	if groupmsg == nil {
		return
	}
	server.writeChan <- groupmsg
}

func (server *TCPServer) SendAllPb(maincmd uint32, pb proto.Message) {
	if server.isClose.IsTrue() {
		return
	}
	sendBuff, erro := server.CreatePBMsg(maincmd, pb)
	if erro != nil {
		xlog.Error("SendSomeConnPb 失败,%v", pb.String())
		return
	}
	server.SendAllConn(sendBuff)
}

func (server *TCPServer) SendSomeConnPb(ConnIDs []uint32, maincmd uint32, pb proto.Message) {
	if server.isClose.IsTrue() {
		return
	}
	sendBuff, erro := server.CreatePBMsg(maincmd, pb)
	if erro != nil {
		xlog.Error("SendSomeConnPb 失败,%v", pb.String())
		return
	}
	server.SendSomeConn(ConnIDs, sendBuff)
}

// 给所有连接发送消息
func (server *TCPServer) SendSomeConn(ConnIDs []uint32, msg []byte) {
	if server.isClose.IsTrue() {
		return
	}
	groupmsg := server.createGroupMessage(ConnIDs, msg)
	if groupmsg == nil {
		return
	}
	server.writeChan <- groupmsg
}

// 从池子里面创建消息
func (server *TCPServer) createGroupMessage(ConnIDs []uint32, msg []byte) *GroupMessage {
	if msg == nil {
		xlog.Error("群发消息为空")
		return nil
	}
	groupval := server.msgPool.Get()
	if groupval == nil {
		xlog.Error("获取消息体错误")
		return nil
	}
	groupmsg, ok := groupval.(*GroupMessage)
	if !ok {
		xlog.Error("获取消息体错误")
		return nil
	}
	groupmsg.Msgdata = msg
	groupmsg.ConnIDs = ConnIDs
	return groupmsg
}

func (server *TCPServer) doSend(message *GroupMessage) {
	if message == nil || message.Msgdata == nil {
		xlog.Error("群发消息为空")
		return
	}
	// 发送给一部分
	if message.ConnIDs != nil && len(message.ConnIDs) > 0 {
		for _, connId := range message.ConnIDs {
			erro := server.sendOneMsg(connId, message.Msgdata)
			if erro != nil {
				xlog.Debug("ConnId 发送消息错误 %v", connId)
			}
		}
		return
	}

	// 发送给全部连接
	server.mutexConns.Lock()
	for _, conn := range server.connMaps {
		conn.Write(message.Msgdata)
	}
	server.mutexConns.Unlock()
}

//获取网络配置,动态获取,方便更新
func (this *TCPServer) GetNetCfg() *csvdata.NetWorkCfg {
	cfg := csvdata.GetNetWorkCfgPtr(this.appID)
	if cfg == nil {
		return nil
	}
	return cfg
}

// 写单个消息
func (server *TCPServer) sendOneMsg(ConnID uint32, msg []byte) error {
	tcpconn := server.GetTcpConnect(ConnID)
	if tcpconn == nil {
		return errors.New(fmt.Sprintf("SendMsg 未找到连接 %v", ConnID))
	}
	// 向写通道投递数据
	tcpconn.Write(msg)
	return nil
}

//获取tcp连接
func (server *TCPServer) GetTcpConnect(ConnID uint32) Conner {
	server.mutexConns.RLock()
	tcpconn := server.connMaps[ConnID]
	server.mutexConns.RUnlock()
	return tcpconn
}
