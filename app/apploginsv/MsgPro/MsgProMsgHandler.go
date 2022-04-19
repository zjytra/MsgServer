/*
创建时间:
作者: zjy
功能介绍:
消息处理自动生成
下面生成的方法可以随便改,自定义的方法不要包含2不然会被过滤掉
*/

package MsgPro

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/app/apploginsv/RoomMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/session"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/snowflake"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/protomsg"
)


//发送消息
func C2L_SendMsgReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_SendMsgReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		conn.Close()
		xlog.Debug("C2L_SendMsgReq 错误:" + erro.Error())
		return
	}
	clSession := session.ClientSessionMgr.GetClientSession(conn.GetConnID())
	//绑定对象已经离线
	if clSession == nil || clSession.PRole == nil {
		conn.Close()
		return
	}

	msg := new(dbmodels.MsgT)
	msg.MsgID.SetVal(snowflake.GUID.NextId())
	msg.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	msg.RoomID.SetVal(clSession.PRole.RoomID.GetVal())
	msg.MsgContent.SetVal(reqMsg.Conent)
	msg.SenderName.SetVal(clSession.PRole.RoleName.GetVal())
	RoomMgr.AddMsg(msg)
	//延迟插入
	dbsys.GameAccountDB.DelayInsert(msg)
	sendMsg := &protomsg.L2C_SendMsgRes{}
	sendMsg.MsgInfo = new(protomsg.MsgInfo)
	msg.BuildPro(sendMsg.MsgInfo)
	//通知当前房间的其他玩家
 	users :=  RoomMgr.GetRoomUsers(clSession.PRole.RoomID.GetVal())
	for conID,_ := range users{
		otherUser := session.ClientSessionMgr.GetClientSession(conID)
		//绑定对象已经离线
		if otherUser == nil || otherUser.PRole == nil || otherUser.PConn == nil {
			RoomMgr.OnRoleLeave(conID,clSession.PRole.RoomID.GetVal())
			continue
		}
		otherUser.PConn.WritePBToMsgRes(Cmd.L2C_SendMsgRes,sendMsg)
	}

}

//获取消息列表
func C2L_GetRoomMsgListReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_GetRoomMsgListReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_GetRoomMsgListReq 错误:" + erro.Error())
		return
	}
	msgs :=	 RoomMgr.GetRoomMsgPro(reqMsg.RoomID)
	if msgs == nil {
		return
	}
	sendMsg := &protomsg.L2C_GetRoomMsgListRes{}
	sendMsg.RoomID = reqMsg.RoomID
	for _,msg := range msgs {
		msgPro := new(protomsg.MsgInfo)
		msg.BuildPro(msgPro)
		sendMsg.Msgs = append(sendMsg.Msgs,msgPro)
	}

	conn.WritePBToMsgRes(Cmd.L2C_GetRoomMsgListRes,sendMsg)
}

//获取玩家信息
func C2L_GetRoleInfoReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_GetRoleInfoReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_GetRoleInfoReq 错误:" + erro.Error())
		return
	}
}

//获取最近10分钟发送最频繁的消息
func C2L_PopularMsgReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_PopularMsgReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_PopularMsgReq 错误:" + erro.Error())
		return
	}
}

//切换房间
func C2L_SwitchRoomReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_SwitchRoomReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_SwitchRoomReq 错误:" + erro.Error())
		return
	}
	clSession := session.ClientSessionMgr.GetClientSession(conn.GetConnID())
	//绑定对象已经离线
	if clSession == nil || clSession.PRole == nil {
		conn.Close()
		return
	}
	//房间相同
	if clSession.PRole.RoomID.GetVal() == reqMsg.RoomID {
		return
	}
	//玩家离开老房间
	RoomMgr.OnRoleLeave(conn.GetConnID(),clSession.PRole.RoomID.GetVal())
	clSession.PRole.RoomID.SetVal(reqMsg.RoomID)
	writer := new(dbsys.DBObjWriteParam)
	writer.AddObjs(clSession.PRole)
	//异步更新无回调
	dbsys.GameAccountDB.AsyncUpdateObj(writer,nil)
	//进入新房间
	RoomMgr.OnRoleEnter(conn.GetConnID(),clSession.PRole)
	sendMsg := &protomsg.L2C_SwitchRoomRes{}
	sendMsg.RoomID = reqMsg.RoomID
	msgs :=	 RoomMgr.GetRoomMsgPro(reqMsg.RoomID)
	if msgs != nil {
		for _,msg := range msgs {
			msgPro := new(protomsg.MsgInfo)
			msg.BuildPro(msgPro)
			sendMsg.Msgs = append(sendMsg.Msgs,msgPro)
		}
	}
	conn.WritePBToMsgRes(Cmd.L2C_SwitchRoomRes,sendMsg)
}