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
	"github.com/zjytra/MsgServer/app/apploginsv/RoleMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoomMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/session"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/snowflake"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/msgcode"
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
	writerParam := &dbsys.DBObjWriteParam{
		ClientConnId: conn.GetConnID(),
		ParamObj:     reqMsg,
	}
	writerParam.AddObjs(msg)
	dbsys.GameAccountDB.AsyncInsertObj(writerParam,nil)
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
	j := 0
	i := len(msgs) - 1 //到着遍历最新50条
	for  ;i >= 0 ; i -- {
		msgPro := new(protomsg.MsgInfo)
		msgs[i].BuildPro(msgPro)
		sendMsg.Msgs = append(sendMsg.Msgs,msgPro)
		j ++
		if j >= 50 {
			break
		}
	}

	conn.WritePBToMsgRes(Cmd.L2C_GetRoomMsgListRes,sendMsg)
}

//获取玩家信息
func C2L_GetRoleInfoReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_GetRoleInfoReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_GetRoleInfoReq 错误:" + erro.Error())
		conn.Close()
		return
	}
	sendMsg := &protomsg.L2C_GetUserInfoRes{}
	//账号包含空格或者非单词字符
	isMatch := strutil.StringHasSpaceOrSpecialChar(reqMsg.RoleName)
	if isMatch {
		conn.WritePBMsgAndCode(Cmd.L2C_GetUserInfoRes, msgcode.RoleNameError, sendMsg)
		return
	}
	//sql注入验证
	isMatch = strutil.StringHasSqlKey(reqMsg.RoleName)
	if isMatch {
		conn.WritePBMsgAndCode(Cmd.L2C_GetUserInfoRes, msgcode.RoleNameError, sendMsg)
		return
	}
	pRole := RoleMgr.GetNameRole(reqMsg.RoleName)
	if pRole != nil {
		sendMsg.RoleName = pRole.RoleName.GetVal()
		sendMsg.LoginTime = pRole.LonginTime.GetVal()
		sendMsg.OnlineTime = pRole.GetOnlineTime()
		sendMsg.RoomID = pRole.RoomID.GetVal()
		//内存中有数据直接返回
		conn.WritePBToMsgRes(Cmd.L2C_GetUserInfoRes, sendMsg)
		return
	}
	dbParam := new(dbsys.DBParam)
	dbParam.CltConn = conn.GetConnID()
	//框架还待完善,由于角色表的查询字段被账号id占了,只有用原来封装的接口了
	sql := "SELECT RoleName,LonginTime,OnlineTime,RoomID,OffLineTime FROM RoleT WHERE RoleName = ?"
	dbsys.GameAccountDB.AsyncQueryStruct(dbParam,dbmodels.QueryRole{},onQueryRoleCb,sql,reqMsg.RoleName)

}

func onQueryRoleCb(param  *dbsys.DBQueryParam,result *dbsys.DBQueryResult)  {
	if param == nil {
		xlog.Debug("onQueryRoleCb 未传递参数 ")
		return
	}
	sendMsg := &protomsg.L2C_GetUserInfoRes{}
	queryRole, pbOk := result.QueryObj.(*dbmodels.QueryRole)
	clSession := session.ClientSessionMgr.GetConn(param.CltConn)
	if !pbOk {
		xlog.Debug("onQueryRoleCb  dbmodels.QueryRole erro")
		clSession.WritePBMsgAndCode(Cmd.L2C_GetUserInfoRes, msgcode.UserNotFind, sendMsg)
		return
	}
	if clSession == nil  {
		xlog.Debug("onQueryRoleCb 玩家已离线")
		return
	}

	sendMsg.RoleName = queryRole.RoleName
	sendMsg.LoginTime = queryRole.LonginTime
	sendMsg.OnlineTime = 0
	sendMsg.RoomID = queryRole.RoomID
	//内存中有数据直接返回
	clSession.WritePBToMsgRes(Cmd.L2C_GetUserInfoRes, sendMsg)
}

//获取最近10分钟发送最频繁的消息
func C2L_PopularMsgReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_PopularMsgReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("C2L_PopularMsgReq 错误:" + erro.Error())
		conn.Close()
		return
	}
	sendMsg := &protomsg.L2C_PopularMsgRes{}
	if RoomMgr.GetRoom(reqMsg.RoomID) == nil {
		conn.WritePBMsgAndCode(Cmd.L2C_PopularMsgRes,msgcode.RoomNotFind, sendMsg)
		return
	}
	nowMs := timeutil.GetCurrentTimeMs()
	msgs :=	 RoomMgr.GetRoomMsgPro(reqMsg.RoomID)
	if msgs != nil {
		for _,msg := range msgs {
			//最近10分钟发送最频繁的消息
			if msg.CreateTime.GetVal() >= nowMs - 10 * 60 * 1000 {
				sendMsg.PopMsg = append(sendMsg.PopMsg,msg.MsgContent.GetVal())
			}
		}
	}
	conn.WritePBToMsgRes(Cmd.L2C_PopularMsgRes,sendMsg)
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
		j := 0
		i := len(msgs) - 1 //到着遍历最新50条
		for  ;i >= 0 ; i -- {
			msgPro := new(protomsg.MsgInfo)
			msgs[i].BuildPro(msgPro)
			sendMsg.Msgs = append(sendMsg.Msgs,msgPro)
			j ++
			if j >= 50 {
				break
			}
		}
	}
	conn.WritePBToMsgRes(Cmd.L2C_SwitchRoomRes,sendMsg)
}