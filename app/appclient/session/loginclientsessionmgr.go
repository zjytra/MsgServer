/*
创建时间: 2021/7/3 23:29
作者: zjy
功能介绍:

*/

package session

import (
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

var(
	ClientSessionMgr *LoginClientSessionMgr
)




//登录服会话管理
type LoginClientSessionMgr struct {
	//客户端服务器
	Client *network.TCPClient
	network.ConnMgr
}

func (this *LoginClientSessionMgr)Start() {
	this.Client = network.NewTCPClient(this,appdata.AppID)
	this.Client.Start()

}

func (this *LoginClientSessionMgr)CreateConn() network.Conner {
	return nil
}

func (this *LoginClientSessionMgr) OnNetWorkConnect(conn network.Conner)  {
	pSession := this.GetConn(conn.GetConnID())
	if pSession != nil {
		xlog.Debug("当前连接id 已经存在")
		return
	}
	//创建session对象
	this.AddConn(conn)
	xlog.Debug("LoginClientSessionMgr 连接成功 %d addr %s ",conn.GetConnID(),conn.RemoteAddrStr())


}

func (this *LoginClientSessionMgr) OnNetWorkClose(conn network.Conner)  {
	psession := this.GetConn(conn.GetConnID())
	if psession == nil {
		xlog.Debug("当前连接id 不存在")
		return
	}
	this.RemoveConn(conn.GetConnID())
}


//来至客户端的消息
func (this *LoginClientSessionMgr)OnClientMsg(conn network.Conner,maincmd uint32, msg []byte) {
	pSession := this.GetConn(conn.GetConnID())
	if pSession == nil {
		xlog.Debug("OnClientMsg 命令 %d 连接已断开 ",maincmd)
		return
	}

	//分消息处理
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
		xlog.Error("命令 = %d  耗时 %v",maincmd, since)
	}
}

