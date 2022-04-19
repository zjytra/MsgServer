/*
创建时间: 2020/7/14
作者: zjy
功能介绍:
处理服务器内部连接的功能
*/

package network

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"net"
	"time"
)

//服务器内部的服务器
type TCPServiceServer struct {
	TCPServer                         //继承客户端服务器的功能部分方法重写
	svOb      ServiceNetEvent         // 网络事件观察者
}



// 创建tcp Sever服务器
func NewTCPServiceServer(netobs ServiceNetEvent,appid int32) *TCPServiceServer {
	tcpsv := new(TCPServiceServer)
	tcpsv.svOb = netobs
	err := tcpsv.initServerData(appid)
	if err != nil {
		xlog.Error("newClientConn %v", err)
		return nil
	}
	return tcpsv
}


func (this *TCPServiceServer) Start() {
	xlog.Debug("TCPServiceServer start")
	this.init()
	go this.serverRun()
	go this.serverEvent()
}

func (server *TCPServiceServer) serverRun() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()
	var tempDelay time.Duration
	xlog.Debug("TCPServiceServer Accept Addr:%v", server.ln.Addr())
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

// 服务器添加连接
func (this *TCPServiceServer) addConn(conn net.Conn) bool {

	this.mutexConns.Lock()//互斥锁
	maxNum := int32(10000)
	cfg := this.GetNetCfg()
	if cfg != nil {
		maxNum = cfg.Max_connect
	}
	if this.GetConnectSize() >= maxNum {
		this.mutexConns.Unlock()
		erro := conn.Close()
		if erro != nil {
			xlog.Warning("超过连接关闭链接错误 %v ", erro)
		}
		xlog.Warning("超过最大链接数,当前连接数%d", this.connetsize)
		return false
	}
	// 创建封装的连接
	ConnID := nextID()
	tcpConn := newServiceConn(conn, ConnID, this.svOb, this.appID, this.msgParser)
	this.connMaps[ConnID] = tcpConn // 存储连接
	this.mutexConns.Unlock() //解锁

	this.connetsize.AddInt32()
	this.wgConns.Add(1)
	go	this.ReceiveData(tcpConn)
	xlog.Debug("当前连接数addConn %d,连接标识%d", this.GetConnectSize(), ConnID)
	return true
}

func (this *TCPServiceServer) SendAllPb(maincmd uint32, pb proto.Message) {
	if this.isClose.IsTrue() {
		return
	}
	sendBuff, erro := this.CreatePBMsg(maincmd , pb)
	if erro != nil {
		xlog.Error("SendSomeConnPb 失败,%v",pb.String())
		return
	}
	this.SendAllConn(sendBuff)
}

func (this *TCPServiceServer) SendSomeConnPb(ConnIDs []uint32,maincmd uint32, pb proto.Message) {
	if this.isClose.IsTrue() {
		return
	}
	sendBuff, erro := this.CreatePBMsg(maincmd , pb)
	if erro != nil {
		xlog.Error("SendSomeConnPb 失败,%v",pb.String())
		return
	}
	this.SendSomeConn(ConnIDs,sendBuff)
}
// 根据命令及protobuf创建包
func (server *TCPServiceServer) CreatePBMsg(maincmd uint32, pb proto.Message) (sendMsg []byte, erro error) {
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
func (this *TCPServiceServer) CreatePackage(maincmd uint32, msg []byte) ([]byte, error) {
	return this.msgParser.PackOne(maincmd, msg)
}
