/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/

package appclient

import (
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"time"
)

type AppClient struct {
	//netmsgsys   *network.NetMsgHandler
	sndHeartTimer uint32 //定时起
	span       time.Duration
}

// 程序启动
func (this *AppClient) OnStart() {
	this.OnInit()

	// dispSys.SetServiceNet(this) // 模拟服务器连接
	//this.tcpclient = network.NewTCPClient(this.dispSys, appdata.OutNetConf)
	//this.tcpclient.Start()
	//this.span = time.Second * 5
	//for i:= 0; i < 3; i++ {
	//	time.Sleep(time.Second * 10)
	//	go func() {
	//		clent := network.NewTCPClient(this,appdata.AppID)
	//		clent.Start()
	//		// 连接成功发送登陆命令
	//		//time.Sleep(time.Second)
	//		//this.TestRegisterAccount(clent.NextCon())
	//		//账号登录
	//		////con := clent.NextCon()
	//		//if con != nil {
	//		//	this.TestLoginAccount(con)
	//		//}
	//
	//	}()
	//}
	
	//this.sndHeartTimer = timersys.NewWheelTimer(time.Second * 30,this.TestTimer,this.dispSys)
}

func (this *AppClient) TestRegisterAccount(pSession network.Conner) {
	//reqCreateAccount := &protomsg.C2L_RegisterAccountMsg{
	//	Username:   "zjy082",
	//	Password:   "jp3411952",
	//	ClientType: model.ClientType_Test,
	//	MacAddr:    osutil.GetUpMacAddr(),
	//	Version:    "",
	//}
	//if conn != nil {
	//	//erro = conn.WritePBMsg(Cmd., Cmd.Sub_C_LS_RegisterAccount, reqCreateAccount)
	//}

}
func (this *AppClient) TestLoginAccount(pSession network.Conner) {
	//reqMsg := &protomsg.C2L_LoginMsg{
	//	Username:   "zty111uuy",
	//	Password:   "jp3411952",
	//	ClientType: model.ClientType_Test,
	//	MacAddr:    osutil.GetUpMacAddr(),
	//	Version:    "",
	//}
	//var erro
	////if conn != nil {
	////	conn.WritePBMsg(Cmd.Cmd_Sockect_Event, Cmd.C2L_ReqLoginAck, nil)
	////	erro = conn.WritePBMsg(Cmd.CMD_Account, Cmd.C_LS_LoginAccount, reqMsg)
	////}
	//if erro != nil {
	//	xlog.Error("TestLoginAccount write erro %v ", erro)
	//}
}


// 发送心跳给世界服
func (this *AppClient) TestTimerTestTimer() {
	defer xlog.RecoverToLog(func() {
		//timersys.StopTimer(this.sndHeartTimer)
	})
	//账号登录
	this.TestLoginAccount(nil)
	//this.TestRegisterAccount()

}
// 初始化
func (this *AppClient) OnInit() bool {
	csvdata.LoadCommonCsvData()
	//this.netmsgsys = network.NewNetMsgHandler()
	return true
}

// 程序运行
func (this *AppClient) OnUpdate() bool {
	// xlog.DebugLog("","run LoginApp")
	
	return true
}

// 关闭
func (this *AppClient) OnRelease() {
}


func (this *AppClient) SendCreateAccount()  {
	// 连接成功发送登陆命令
	//reqCreateAccount := &protomsg.C2L_RegisterAccountMsg{
	//	Username:   "zty111uuy",
	//	Password:   "jp3411952",
	//	ClientType: model.ClientType_Test,
	//	MacAddr:    osutil.GetUpMacAddr(),
	//	Version: "",
	//}
	//erro := this.tcpclient.WritePBMsg(Cmd.CMD_Account, Cmd.Sub_C_LS_RegisterAccount, reqCreateAccount)
	//if erro != nil {
	//	xlog.Error("OnNetWorkConnect write erro %v ", erro.Error())
	//}
}

