/*
创建时间:
作者: zjy
功能介绍:
消息注册自动生成
*/
package appclient

import (

	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/app/appclient/session"
	"github.com/zjytra/MsgServer/app/appclient/Login"
	"github.com/zjytra/MsgServer/app/appclient/MsgPro"
)

func RegisterMsg(){

	//---------- 模块 Login 开始 ----------
	//连接登录服成功后登录服回复
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_LoginConnectAck,Login.L2C_LoginConnectAck)
	//登录返回角色
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_LoginMsg,Login.L2C_LoginMsg)
	//返回创建的角色id与房间号
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_CreateRoleRes,Login.L2C_CreateRoleRes)
	//---------- 模块 Login 结束 ----------
	//---------- 模块 MsgPro 开始 ----------
	//发送消息
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_SendMsgRes,MsgPro.L2C_SendMsgRes)
	//获取消息列表回复
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_GetRoomMsgListRes,MsgPro.L2C_GetRoomMsgListRes)
	//获取玩家信息回复
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_GetUserInfoRes,MsgPro.L2C_GetUserInfoRes)
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_PopularMsgRes,MsgPro.L2C_PopularMsgRes)
	//切换房间回复
	session.ClientSessionMgr.RegisterHandle(Cmd.L2C_SwitchRoomRes,MsgPro.L2C_SwitchRoomRes)
	//---------- 模块 MsgPro 结束 ----------
}