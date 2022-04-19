/*
创建时间:
作者: zjy
功能介绍:
消息注册自动生成
*/
package apploginsv

import (

	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/app/apploginsv/session"
	"github.com/zjytra/MsgServer/app/apploginsv/Login"
	"github.com/zjytra/MsgServer/app/apploginsv/MsgPro"
)

func RegisterMsg(){

	//---------- 模块 Login 开始 ----------
	//连接登录服成功后发送确认
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_LoginConnectAck,Login.C2L_LoginConnectAck)
	//登陆
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_LoginMsg,Login.C2L_LoginMsg)
	//创建角色
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_CreateRoleReq,Login.C2L_CreateRoleReq)
	//---------- 模块 Login 结束 ----------
	//---------- 模块 MsgPro 开始 ----------
	//发送消息
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_SendMsgReq,MsgPro.C2L_SendMsgReq)
	//获取消息列表
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_GetRoomMsgListReq,MsgPro.C2L_GetRoomMsgListReq)
	//获取玩家信息
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_GetRoleInfoReq,MsgPro.C2L_GetRoleInfoReq)
	//获取最近10分钟发送最频繁的消息
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_PopularMsgReq,MsgPro.C2L_PopularMsgReq)
	//切换房间
	session.ClientSessionMgr.RegisterHandle(Cmd.C2L_SwitchRoomReq,MsgPro.C2L_SwitchRoomReq)
	//---------- 模块 MsgPro 结束 ----------
}