/*
创建时间: 2021/7/3 23:29
作者: zjy
功能介绍:

*/

package session

import (
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/app/apploginsv/AccountMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoleMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoomMgr"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

var(
	ClientSessionMgr *LoginClientSessionMgr
)

func InitSessionMgr()  {
	ClientSessionMgr = new(LoginClientSessionMgr)
	ClientSessionMgr.Start()
}


//登录服会话管理
type LoginClientSessionMgr struct {
	network.ConnMgr
	//客户端服务器
	ClientServer *network.TCPServer
	//客户端连接绑定
	ClientSessions  map[uint32]*LsClientSession
}

func (this *LoginClientSessionMgr)Start() {
	this.InitSessionMgr()
	this.ClientSessions = make( map[uint32]*LsClientSession)
	this.ClientServer = network.NewTcpServer(this,appdata.AppID)
	this.ClientServer.Start()
}

func (this *LoginClientSessionMgr)CreateConn(acc *dbmodels.AccountT,conn network.Conner)*LsClientSession {
	psession := this.GetClientSession(conn.GetConnID())
	if psession == nil { //绑定账号连接
		psession = new(LsClientSession)
		psession.PConn = conn
		psession.PAcc = acc
		this.ClientSessions[conn.GetConnID()] = psession
	}
	return psession
}

func (this *LoginClientSessionMgr) OnNetWorkConnect(conn network.Conner)  {
    psession := this.GetConn(conn.GetConnID())
	if psession != nil {
		xlog.Debug("当前连接id 已经存在")
		return
	}
	this.AddConn(conn)
	//创建session对象
	//psession = this.CreateConn(conn)
	xlog.Debug("客户端登录服连接成功 %d addr %s ",conn.GetConnID(),conn.RemoteAddrStr())
}

func (this *LoginClientSessionMgr) OnNetWorkClose(conn network.Conner)  {
	this.RemoveConn(conn.GetConnID())
	psession := this.GetClientSession(conn.GetConnID())
	if psession == nil {
		return
	}
	accName := ""
	if psession.PAcc != nil {
		accName = psession.PAcc.LoginName.GetVal()
	}
	//客户端连接断开要通知
	xlog.Debug("客户端断开登录服  %d addr %s 账号:%s",conn.GetConnID(),conn.RemoteAddrStr(),accName)
	this.RemoveClientSession(conn.GetConnID())
}



//来至客户端的消息
func (this *LoginClientSessionMgr)OnClientMsg(conn network.Conner,maincmd uint32, msg []byte) {
	pSession := this.GetConn(conn.GetConnID())
	if pSession == nil {
		xlog.Debug("OnClientMsg 命令 %d 连接已断开 ",maincmd)
		return
	}
	handle := this.GetMsgHandle(maincmd)
	if handle == nil {
		//踢掉
		conn.Close()
		xlog.Error("命令 = %d 未注册",maincmd)
		return
	}
	startT := timeutil.GetCurrentTimeMs()		//计算当前时间
	handle(pSession,msg) //查看是否注册对应的处理函数
	since := mathutil.MaxInt64(0,timeutil.GetCurrentTimeMs() - startT)
	if since > timeutil.FrameTimeMs { //大于50毫秒
		xlog.Debug("命令 = %d  耗时 %v",maincmd, since)
	}
}


func (this *LoginClientSessionMgr) GetClientSession(connid uint32) *LsClientSession {
	session := this.ClientSessions[connid]
	if session == nil {
		return nil
	}
	return session
}


func (this *LoginClientSessionMgr) GetConnSize() int32 {
	if this.ClientServer == nil {
		return 0
	}
	return this.ClientServer.GetConnectSize()
}



func (this *LoginClientSessionMgr) RemoveClientSession(connid uint32){
	session := this.ClientSessions[connid]
	if session == nil {
		return
	}
	if session.PAcc != nil  {
		AccountMgr.DelAccount(session.PAcc.LoginName.GetVal())
	}
	if session.PRole != nil  {
		session.PRole.OffLineTime.SetVal(timeutil.GetCurrentTimeMs())
		updateRole := new(dbsys.DBObjWriteParam)
		updateRole.AddObjs(session.PRole)
		dbsys.GameAccountDB.AsyncUpdateObj(updateRole, nil)
		RoleMgr.Remove(session.PRole)
		RoomMgr.OnRoleLeave(connid,session.PRole.RoomID.GetVal())
	}


}