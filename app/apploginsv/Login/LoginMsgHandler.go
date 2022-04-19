/*
创建时间:
作者: zjy
功能介绍:
消息处理自动生成
下面生成的方法可以随便改,自定义的方法不要包含2不然会被过滤掉
//登录服与客户端交互走http请求
*/

package Login

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/app/apploginsv/AccountMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoleMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/RoomMgr"
	"github.com/zjytra/MsgServer/app/apploginsv/session"
	"github.com/zjytra/MsgServer/conf"
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


//连接登录服成功后发送确认
func C2L_LoginConnectAck(pSession network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_LoginConnectAck{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Error("C2L_LoginConnectAck 错误:" + erro.Error())
		pSession.Close()
		return
	}
	if reqMsg.Code != conf.SvJson.InnerKey {
		xlog.Debug("验证码错误: %v", reqMsg.Code)
		pSession.Close()
		return
	}
	xlog.Debug("收到客户端消息")
	pSession.SetContAck()
	send := protomsg.L2C_LoginConnectAck{}
	send.UnixNano = timeutil.GetCurrentTimeMs()
	pSession.WritePBToMsgRes(Cmd.L2C_LoginConnectAck, &send)
}

//登陆
func C2L_LoginMsg(pSession network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_LoginMsg{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Error("C2L_LoginMsg 错误:" + erro.Error())
		return
	}
	sendMsg := &protomsg.L2C_LoginMsg{}
	//已经登陆过的连接
	pOldSession := session.ClientSessionMgr.GetClientSession(pSession.GetConnID())
	if pOldSession != nil {
		sendMsg.RoleID = 0
		if pOldSession.PRole != nil {
			pOldSession.PRole.BuildRolePro(sendMsg)
		}
		pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.AccountCode_IsLogined, sendMsg)
		return
	}
	//已经登录过
	//acc := AccountMgr.GetAccountInfoByUserName(reqMsg.Acc)
	//if acc != nil {
	//	pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.AccountCode_IsLogined, sendMsg)
	//	return
	//}
	//防止异步调用的时候重复投递
	if AccountMgr.IsAsyncLogin(reqMsg.Acc) {
		pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.AccountCode_IsLogining, sendMsg)
		return
	}
	//账号合法验证
	code := msgcode.IsVerify(reqMsg.Acc)
	if code != msgcode.Succeed {
		pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, code, sendMsg)
		return
	}

	//在redis中查询是否有注册
	intCmd := dbsys.RedisCli.SAdd(dbsys.CtxBg, "registerAccount", reqMsg.Acc)
	if intCmd == nil { //redis 查询不成功再查询数据库
		queryAccount(pSession.GetConnID(), reqMsg)
		return
	}
	if intCmd.Err() != nil {
		queryAccount(pSession.GetConnID(), reqMsg)
		return
	}
	//插入不成功已经注册
	if intCmd.Val() == 0 {
		queryAccount(pSession.GetConnID(), reqMsg)
		return
	}
	//插入成功就创建账号
	createAccount(pSession, reqMsg)
}

//异步查询账号
func queryAccount(ClientConnId uint32, reqMsg *protomsg.C2L_LoginMsg) {
	AccountMgr.SetAsyncLogin(reqMsg.Acc)
	pAccount := dbmodels.NewAccountT(reqMsg.Acc)
	queryParam := &dbsys.DBObjQueryParam{
		ClientConnId: ClientConnId,
		ParamObj:     reqMsg,
		DbObj:        pAccount,
	}
	dbsys.GameAccountDB.AsyncLoadObJerFromDBAndCb(queryParam, OnQueryAccCb)
}

func createAccount(pSession network.Conner,reqMsg *protomsg.C2L_LoginMsg) {
	sendMsg := &protomsg.L2C_LoginMsg{}
	//自增id
	accID := dbsys.RedisCli.Incr(dbsys.CtxBg, "AccountID")
	if accID == nil {
		pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RegisterRedisError, sendMsg)
		return
	}
	if accID.Err() != nil {
		xlog.Error("C2L_RegisterAccountMsg 错误: %v", accID.Err())
		pSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RegisterRedisError, sendMsg)
		return
	}
	//id
	//异步创建账号中
	AccountMgr.SetAsyncLogin(reqMsg.Acc)
	pAccount := dbmodels.NewAccountT(reqMsg.Acc)
	pAccount.AccID.SetVal(accID.Val())
	pAccount.LoginName.SetVal(reqMsg.Acc)
	nowTime := timeutil.GetCurrentTimeMs()
	pAccount.RegisterTime.SetVal(nowTime)
	pAccount.RegisterIp.SetVal(pSession.RemoteAddrIp())
	writerParam := &dbsys.DBObjWriteParam{
		ClientConnId: pSession.GetConnID(),
		ParamObj:     reqMsg,
	}
	writerParam.AddObjs(pAccount)
	dbsys.GameAccountDB.AsyncInsertObj(writerParam, OnDBRegister)
}

//操作数据库后回调
func OnDBRegister(result dbsys.DBErrorCode, param *dbsys.DBObjWriteParam) {
	if param == nil || param.ParamObj == nil {
		xlog.Debug("OnDBRegister 未传递参数 ")
		return
	}
	reqMsg, pbOk := param.ParamObj.(*protomsg.C2L_LoginMsg)
	if !pbOk {
		xlog.Debug("转换C2L_RegisterAccountMsg失败")
		return
	}
	AccountMgr.DelAsyncLogin(reqMsg.Acc)

	if param.DbObjs == nil || len(param.DbObjs) < 1 {
		xlog.Debug("创建账号数据未返回")
		return
	}

	clSession := session.ClientSessionMgr.GetConn(param.ClientConnId)
	if clSession == nil {
		xlog.Debug("OnDBLogin 玩家已离线")
		return
	}
	//才创建账号通知客户端创建角色
	sendMsg := &protomsg.L2C_LoginMsg{
		RoleID: 0,
	}
	//创建账号失败
	if result != dbsys.DBSUCCESS {
		//删除redis
		intCmd := dbsys.RedisCli.SRem(dbsys.CtxBg, "registerAccount",reqMsg.Acc)
		if intCmd == nil {
			clSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RegisterRedisError, sendMsg)
			return
		}
		if intCmd.Err() != nil {
			xlog.Error("OnDBRegister SRem 错误: %v", intCmd.Err())
			clSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RegisterRedisError, sendMsg)
			return
		}
		return
	}
	pAcc := param.DbObjs[0].(*dbmodels.AccountT)
	//session绑定对象
	session.ClientSessionMgr.CreateConn(pAcc,clSession)
	//创建账号成功
	clSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.Succeed, sendMsg)
}



//异步查询账号回调
func OnQueryAccCb(result dbsys.DBErrorCode, param *dbsys.DBObjQueryParam) {
	if param == nil || param.DbObj == nil || param.ParamObj == nil {
		xlog.Debug("OnDBLogin 未传递参数 ")
		return
	}
	reqMsg, pbOk := param.ParamObj.(*protomsg.C2L_LoginMsg)
	if !pbOk {
		xlog.Debug("转换C2L_LoginMsg失败")
		return
	}
	pAccunt, accOk := param.DbObj.(*dbmodels.AccountT)
	if !accOk {
		xlog.Debug("OnDBLogin dbmodels.AccountT 转换失败 ")
		return
	}

	clSession := session.ClientSessionMgr.GetConn(param.ClientConnId)
	if clSession == nil {
		xlog.Debug("OnDBLogin 玩家已离线")
		return
	}
	if result == dbsys.NODATA {
		//没有数据创建账号
		createAccount(clSession,reqMsg)
	} else if result == dbsys.DBSQLERRO {
		sendMsg := &protomsg.L2C_LoginMsg{
			RoleID: 0,
		}
		clSession.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.Login_DBERRO, sendMsg)
	} else if result == dbsys.DBSUCCESS {
		//将账号保存在内存中
		AccountMgr.AddAccount(pAccunt)
		//session绑定对象
		session.ClientSessionMgr.CreateConn(pAccunt,clSession)
		//还需要去查询角色
		pRole := dbmodels.NewRoleT()
		pRole.AccID.SetVal(pAccunt.GetUID()) //用id查询角色
		queryParam := &dbsys.DBObjQueryParam{
			ClientConnId: param.ClientConnId,
			ParamObj:     reqMsg,
			DbObj:        pRole,
		}
		dbsys.GameAccountDB.AsyncLoadObJerFromDBAndCb(queryParam, OnQueryRoleCb)
	}
}

//异步查询角色回调
func OnQueryRoleCb(result dbsys.DBErrorCode, param *dbsys.DBObjQueryParam) {
	if param == nil || param.DbObj == nil || param.ParamObj == nil {
		xlog.Debug("OnDBLogin 未传递参数 ")
		return
	}
	reqMsg, pbOk := param.ParamObj.(*protomsg.C2L_LoginMsg)
	if !pbOk {
		xlog.Debug("转换C2L_LoginMsg失败")
		return
	}
	AccountMgr.DelAsyncLogin(reqMsg.Acc)
	pRole, accOk := param.DbObj.(*dbmodels.RoleT)
	if !accOk {
		xlog.Debug("OnDBLogin dbmodels.AccountT 转换失败 ")
		return
	}

	clSession := session.ClientSessionMgr.GetClientSession(param.ClientConnId)
	if clSession == nil || clSession.PConn == nil {
		xlog.Debug("OnDBLogin 玩家已离线")
		return
	}
	sendMsg := &protomsg.L2C_LoginMsg{}
	if result == dbsys.NODATA {
		sendMsg.RoleID = 0
		clSession.PConn.WritePBToMsgRes(Cmd.L2C_LoginMsg,sendMsg)
	} else if result == dbsys.DBSQLERRO {
		clSession.PConn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.Login_DBERRO, sendMsg)
	} else if result == dbsys.DBSUCCESS {
		//连接绑定角色
		clSession.PRole = pRole
		RoleMgr.AddRole(pRole)
		RoomMgr.OnRoleEnter(param.ClientConnId,clSession.PRole)
		pRole.BuildRolePro(sendMsg)
		clSession.PConn.WritePBToMsgRes(Cmd.L2C_LoginMsg,sendMsg)
	}
}

//创建角色
func C2L_CreateRoleReq(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.C2L_CreateRoleReq{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		conn.Close()
		xlog.Debug("C2L_CreateRoleReq 错误:" + erro.Error())
		return
	}
	sendMsg := &protomsg.C2L_CreateRoleReq{}

	if RoleMgr.IsAsyncCreate(reqMsg.Username) {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RoleCreateIng, sendMsg)
		return
	}

	clSession := session.ClientSessionMgr.GetClientSession(conn.GetConnID())
	if clSession == nil  { //玩家已经离线
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.Login_DBERRO, sendMsg)
		return
	}
	if clSession.PRole != nil {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RoleIsExists, sendMsg)
		return
	}
	//账号包含空格或者非单词字符
	isMatch := strutil.StringHasSpaceOrSpecialChar(reqMsg.Username)
	if isMatch {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RoleNameError, sendMsg)
		return
	}
	//sql注入验证
	isMatch = strutil.StringHasSqlKey(reqMsg.Username)
	if isMatch {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RoleNameError, sendMsg)
		return
	}
	//在redis中查询是否有注册
	intCmd := dbsys.RedisCli.SAdd(dbsys.CtxBg, "roleName", reqMsg.Username)
	if intCmd == nil { //redis
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.SysError, sendMsg)
		return
	}
	if intCmd.Err() != nil {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.SysError, sendMsg)
		return
	}
	//插入不成功已经创建过
	if intCmd.Val() == 0 {
		conn.WritePBMsgAndCode(Cmd.L2C_LoginMsg, msgcode.RoleNameRepeat, sendMsg)
		return
	}
	RoleMgr.SetAsyncCreate(reqMsg.Username)
	pRole := dbmodels.NewRoleT()
	pRole.RoleID.SetVal(snowflake.GUID.NextId())
	pRole.AccID.SetVal(clSession.PAcc.AccID.GetVal())
	pRole.RoleName.SetVal(reqMsg.Username)
	pRole.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	pRole.LonginTime.SetVal(timeutil.GetCurrentTimeMs())
	pRole.RoomID.SetVal(1)
	writerParam := &dbsys.DBObjWriteParam{
		ClientConnId: conn.GetConnID(),
		ParamObj:     reqMsg,
	}
	writerParam.AddObjs(pRole)
	dbsys.GameAccountDB.AsyncInsertObj(writerParam, OnDBCreateRole)
}



//数据库创建角色回调
func OnDBCreateRole(result dbsys.DBErrorCode, param *dbsys.DBObjWriteParam) {
	if param == nil || param.ParamObj == nil {
		xlog.Debug("OnDBRegister 未传递参数 ")
		return
	}
	reqMsg, pbOk := param.ParamObj.(*protomsg.C2L_CreateRoleReq)
	if !pbOk {
		xlog.Debug("转换C2L_RegisterAccountMsg失败")
		return
	}
	RoleMgr.DelAsyncCreate(reqMsg.Username)
	if param.DbObjs == nil || len(param.DbObjs) < 1 {
		xlog.Debug("创建账号数据未返回")
		return
	}

	clSession := session.ClientSessionMgr.GetClientSession(param.ClientConnId)
	if clSession == nil  || clSession.PConn == nil {
		xlog.Debug("OnDBLogin 玩家已离线")
		return
	}
	//才创建账号通知客户端创建角色
	sendMsg := &protomsg.L2C_CreateRoleRes{
		RoleID: 0,
	}
	//创建账号失败
	if result != dbsys.DBSUCCESS {
		//删除redis
		intCmd := dbsys.RedisCli.SRem(dbsys.CtxBg, "roleName",reqMsg.Username)
		if intCmd == nil {
			clSession.PConn.WritePBMsgAndCode(Cmd.L2C_CreateRoleRes, msgcode.CreateRoleRedisError, sendMsg)
			return
		}
		if intCmd.Err() != nil {
			xlog.Error("OnDBRegister SRem 错误: %v", intCmd.Err())
			clSession.PConn.WritePBMsgAndCode(Cmd.L2C_CreateRoleRes, msgcode.CreateRoleRedisError, sendMsg)
			return
		}
		clSession.PConn.WritePBMsgAndCode(Cmd.L2C_CreateRoleRes, msgcode.CreateRoleDBError, sendMsg)
		return
	}
	clSession.PRole = param.DbObjs[0].(*dbmodels.RoleT)
	sendMsg.RoleID = clSession.PRole.RoleID.GetVal()
	sendMsg.RoomID = clSession.PRole.RoomID.GetVal()
	sendMsg.RoleName = clSession.PRole.RoleName.GetVal()
	RoleMgr.AddRole(clSession.PRole)
	RoomMgr.OnRoleEnter(param.ClientConnId,clSession.PRole)
	//session绑定对象
	//创建账号成功
	clSession.PConn.WritePBMsgAndCode(Cmd.L2C_CreateRoleRes, msgcode.Succeed, sendMsg)
}