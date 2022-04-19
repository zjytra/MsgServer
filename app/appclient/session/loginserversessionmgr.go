/*
创建时间: 2021/7/3 23:29
作者: zjy
功能介绍:
管理服务器会话列表
*/

package session

import (
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)


var(
	ServerSessionMgr *CServerSessionMgr
)

//登录服会话管理
type CServerSessionMgr struct {
	//客户端服务器
	MonitorClient *network.ServiceClient
	network.ServerConnMgr
}

func (this *CServerSessionMgr)Start() {
	this.InitSessionMgr()
	//this.MonitorClient = network.NewServiceClient(this, conf.SvJson.MonitorID)
	//this.MonitorClient.Start()
}



func (this *CServerSessionMgr)OnServiceLink(conn *network.ServiceConn) {
	pSession := this.GetServerConnByConnId(conn.GetConnID())
	if pSession != nil {
	    xlog.Debug("当前连接id 已经存在")
		return
	}
	//创建session对象
	//this.CreateConn(conn)


	//连接监控服后发送消息
	//if pSession.GetCfgKind() == model.APP_Monitor {
	//	sendMsg := &protomsg.L2M_ConnectAck{
	//		AppId: appdata.AppID,
	//		Code:  conf.SvJson.InnerKey,
	//		Num:   network.ClientServer.GetConnectSize(),
	//	}
	//	conn.WritePBMsg(Cmd.L2M_ConnectAck,sendMsg)
	//	return nil
	//}
}

func (this *CServerSessionMgr)OnServiceClose(conn *network.ServiceConn) {
	this.RemoveConnByConnID(conn.GetConnID())
}

//来至服务器的消息
func (this *CServerSessionMgr)OnServerMsg(conn *network.ServiceConn,maincmd uint32, msg []byte) {
	this.OnServerHandlerMsg(conn,maincmd,msg)
}




