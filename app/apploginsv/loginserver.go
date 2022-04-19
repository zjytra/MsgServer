/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/ 

package apploginsv

import (
	"github.com/zjytra/MsgServer/app/apploginsv/AccountMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoleMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoomMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/session"
	"github.com/zjytra/MsgServer/app/comm"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
)



type LoginServer struct {
	oneMinuteTimeID uint32//定时器id
	sendheart  uint32    //给世界服发送心跳
}

// 程序启动
func (this *LoginServer)OnStart() {
	initOK := this.OnInit()
	if !initOK {
		panic("LoginServer 初始化失败")
	}
}

//初始化
func (this *LoginServer)OnInit() bool{
	// csvdata.LoadLoginCsvData()
	AccountMgr.Init()
	RoleMgr.Init()
	RoomMgr.Init()
	this.InitDB()
	dbsys.AsyncOpenRedis("zjy1")
	comm.InitMsgVerify()
	//接收客户端的消息
	session.InitSessionMgr()
	RegisterMsg()
	this.AddTimer()
	return true
}

//先注册数据库表
func (this *LoginServer)InitDB(){
	//初始化数据库
	dbsys.InitAccountDB()
	//注册账号表
	dbsys.GameAccountDB.RegisterTable(dbmodels.AccountT{})
	dbsys.GameAccountDB.RegisterTable(dbmodels.RoleT{})
	dbsys.GameAccountDB.RegisterTable(dbmodels.MsgT{})
	dbsys.GameAccountDB.SyncBD()
	dbsys.GameAccountDB.StartTimer()
	//加载所有的消息
	loadAll := new(dbsys.DBObjQueryParam)
	loadAll.DbObj = new(dbmodels.MsgT)
	dbsys.GameAccountDB.AsyncLoadAllObjs(loadAll,LoadAllMsgCb)
	timeutil.GetTimeNow()
}



//读取客戶端发来的消息


// 关闭
func (this *LoginServer)OnRelease(){
	comm.ReleaseData()
}


func (this *LoginServer)AddTimer(){
	//timersys.Timers.AfterFunc(time.Minute.Milliseconds(),-1, this.PerOneMinuteTimer) //每分钟调用
	//this.sendheart = timersys.NewWheelTimer(time.Second * 30,this.SendHeartToWS,dispatch.ServerDisp)
}



func (this *LoginServer)ReleaseTimer(){
}

//加载所有消息到内存中
func LoadAllMsgCb(result dbsys.DBErrorCode,param *dbsys.DBObjQueryParam,MoreObj []dbsys.DBObJer)  {

	if result == dbsys.DBSUCCESS {
		if MoreObj != nil && len(MoreObj) >0 {
			for _,obj :=range MoreObj{
				pMsg := obj.(*dbmodels.MsgT)
				RoomMgr.AddMsg(pMsg)
			}
		}
	}
}



