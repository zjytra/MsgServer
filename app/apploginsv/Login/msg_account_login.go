package Login

import (
	"github.com/zjytra/MsgServer/engine_core/network"
)

//////////////////////////////////账号登录//////////////////////////////////

//登录向中心服务器请求账号登录
func LSDCLoginAccountMsgHandler(pSession network.Conner, msgdata []byte)  {
	//loginAccout := &protomsg.LS_DC_ReqLoginMsg{}
	//erro := proto.Unmarshal(msgdata, loginAccout)
	//if erro != nil {
	//	xlog.Error("ServerHeartBeat %v", erro)
	//	return erro
	//}
	//pAccounts := PaccountMgr.GetAccountByUserName(loginAccout.Username)
	//if pAccounts != nil { //账号已经存在
	//	//校验密码
	//	if !CheckPassWord(conn, loginAccout) {
	//		return nil
	//	}
	//	state := pAccounts.AccountState()
	//	xlog.Debug("账号登录状态 %v", dcmodel.AccountStateToStr(state))
	//	switch state {
	//	case dcmodel.AccountState_None, dcmodel.AccountState_Leave: //注册了没有登录，或者主动退出的玩家
	//		//1.玩家没有登陆过登录
	//		DoAccountNotExistLogin(conn, loginAccout)
	//	case dcmodel.AccountState_Offline:
	//		//2.玩家离线 需要查看玩家是在那个服务器
	//		OnAccountOfflineLogin(conn, loginAccout)
	//	case dcmodel.AccountState_Online:
	//		//3.玩家在线过查看是否在登录服,在登录服是否是同一个连接
	//		OnAccountOnlineLogin(conn, loginAccout)
	//	default:
	//		xlog.Warning("账号%v状态%v错误", pAccounts.PDBAccountData.LoginName, pAccounts.AccountState())
	//	}
	//	//找到账号逻辑处理完
	//	return nil
	//}
	////在内存中未找到数据
	//LoginAccountNotFindInMemory(conn,loginAccout)
}

//校验密码
//func CheckPassWord(pSession network.Conner,regAccout *protomsg.LS_DC_ReqLoginMsg) bool {
//	//pAccounts := PaccountMgr.GetAccountByUserName(regAccout.Username)
//	//if strings.Compare(pAccounts.PDBAccountData.LoginPwd, regAccout.Password) != 0 { //登录密码错误
//	//	respn := &protomsg.DC_LS_RespnLoginAccoutMsg{
//	//		ClientConnID: regAccout.GetClientConnID(),
//	//		Username:     regAccout.GetUsername(),
//	//		AccountID:    pAccounts.PDBAccountData.AccountID,
//	//		RestCode :    msgcode.AccountCode_IsLoginPassWordIsErro,
//	//	}
//	//	conn.WritePBMsg(Cmd.CMD_Account, Cmd.DC_LS_LoginAccount, respn)
//	//	return false
//	//}
//	return true
//}


//走登录成功的逻辑
//func DoAccountNotExistLogin(pSession network.Conner, loginAccount *protomsg.LS_DC_ReqLoginMsg) {
//	pAccounts := PaccountMgr.GetAccountByUserName(loginAccount.Username)
//	if  pAccounts == nil {
//		return
//	}
//	xlog.Debug("%v oldpaccount=== %v",osutil.GetRuntimeFileAndLineStr(0), pAccounts)
//	pAccounts.SetAccountState(dcmodel.AccountState_Online)
//	pAccounts.SetClientConnID(loginAccount.GetClientConnID())
//	pAccounts.SetExprationTime(0)
//	//获取登录服的信息
//	pAccounts.SetAccountServerInfo(PServerInfoMgr.GetServerInfoByConnID(conn.GetConnID()))
//
//	var accountID uint64
//	accountData := pAccounts.PDBAccountData
//	if accountData != nil {
//		accountData.LoginPwd = loginAccount.Password
//		accountData.LoginIp = &loginAccount.ClientIp
//		nowTimeStr := timeutil.GetTimeALLStr(timeutil.GetTimeNow())
//		accountData.LoginTime = &nowTimeStr
//		accountData.LoginMacAddr = &loginAccount.MacAddr
//		accountID = accountData.AccountID
//	}
//
//	pAccounts = PaccountMgr.GetAccountByUserName(loginAccount.Username)
//	xlog.Debug("%v newpaccount=== %v",osutil.GetRuntimeFileAndLineStr(0), pAccounts)
//	//更新账号信息并且记录日志
//	UpdateAccountsLoginToDBAndRecored(accountData,loginAccount.ClientType)
//	//返回登录成功消息
//	respn := &protomsg.DC_LS_RespnLoginAccoutMsg{
//		ClientConnID: loginAccount.GetClientConnID(),
//		Username:     loginAccount.GetUsername(),
//		AccountID:    accountID,
//		RestCode:     msgcode.AccountCode_Login_Succeed,
//	}
//	//TODO 这里是否要发送网关的连接信息 先发送下去不连接
//    Account,_ := conn.GetPBByteArr(Cmd.CMD_Account, Cmd.DC_LS_LoginAccount, respn)
//	conn.WriteMsg(Account,GetGateWayServerInfoByte(accountID))
//	//conn.WritePBMsg(Cmd.CMD_Account, Cmd_Account.DC_LS_LoginAccount, respn)
//	//发送多个消息
//}
//
//func UpdateAccountsLoginToDBAndRecored(accountData *dbmodels.Accounts,logintype uint32)  {
//	if accountData != nil && accountData.AccountID > 0 {
//		strSql := fmt.Sprintf("UPDATE accounts set LoginIp ='%s',LoginTime='%s',LoginMacAddr='%s'  WHERE AccountID =%d;",
//			*accountData.LoginIp,*accountData.LoginTime,*accountData.LoginMacAddr,accountData.AccountID)
//		dbsys.GameDB.AsyncExecute(strSql,nil,nil)
//		//写日志
//		ymf := timeutil.GetYearMonthFromatStrByTimeString(*accountData.LoginTime)
//		strSql = fmt.Sprintf("INSERT INTO accounts_login_record_%s(AccountID, LoginTime, LoginIp, Phone,LoginType,LoginMacAddr) VALUES (%d,'%s','%s','%s',%d,'%s')",
//			ymf,accountData.AccountID,*accountData.LoginTime,*accountData.LoginIp,accountData.Phone,logintype,*accountData.LoginMacAddr)
//		dbsys.LogDB.AsyncExecute(strSql,nil,nil)
//	}
//
//}
//
////玩家离线登录
//func OnAccountOfflineLogin(pSession network.Conner, loginAccount *protomsg.LS_DC_ReqLoginMsg) {
//	//2.玩家离线 需要查看玩家是在那个服务器
//	pAccounts := PaccountMgr.GetAccountByUserName(loginAccount.Username)
//	if  pAccounts == nil {
//		return
//	}
//	appKind := PServerInfoMgr.GetAppKindByAppID(pAccounts.GetServerAppID())
//	switch appKind { //查看玩家原来在哪一个服务器
//	case model.APP_LoginServer : //玩家在登录服
//		DoAccountNotExistLogin(conn, loginAccount) //直接让其登录成功
//	case model.APP_GATEWAY:     //如果在网关
//	//TODO 告诉网关重连
//	}
//}
//
////玩家在线登录
//func OnAccountOnlineLogin(pSession network.Conner, loginAccount *protomsg.LS_DC_ReqLoginMsg) {
//	//3.玩家在线过查看是否在登录服,在登录服是否是同一个连接
//	pAccounts := PaccountMgr.GetAccountByUserName(loginAccount.Username)
//	if  pAccounts == nil {
//		return
//	}
//	appInfo := PServerInfoMgr.GetServerInfoByAppID(pAccounts.GetServerAppID())
//	switch appInfo.AppKind { //查看玩家在哪一个服务器
//	case model.APP_LoginServer :
//		//通知老连接
//		respn := &protomsg.DC_LS_RespnLoginAccoutMsg{
//			ClientConnID: pAccounts.GetClientConnID(),
//			Username:     loginAccount.GetUsername(),
//			AccountID:    pAccounts.PDBAccountData.AccountID,
//			RestCode:     msgcode.AccountCode_IsLogined,
//		}
//		PdataCenter.tcpserver.WritePBMsgByConnID(appInfo.ConnID, Cmd.CMD_Account, Cmd.DC_LS_LoginAccount,respn)
//		//断开远端连接
//		CloseFarEndClientConn(appInfo.ConnID, pAccounts.PDBAccountData.AccountID)
//
//		//通知新连接
//		respn.ClientConnID = loginAccount.GetClientConnID()
//		PdataCenter.tcpserver.WritePBMsgByConnID(conn.GetConnID(), Cmd.CMD_Account, Cmd.DC_LS_LoginAccount,respn)
//		CloseFarEndClientConn(conn.GetConnID(), pAccounts.PDBAccountData.AccountID)
//	case model.APP_GATEWAY:
//
//	}
//
//}
//
//
////在内存中未找到账号数据
//func LoginAccountNotFindInMemory(pSession network.Conner, loginAccount *protomsg.LS_DC_ReqLoginMsg){
//	//accountLogin := PaccountLoginPool.Pop()
//	param := dbsys.PDBParamPool.Pop()
//	param.ClientConnID = loginAccount.GetClientConnID()
//	param.CbDispSys = PdataCenter.dispSys
//	param.ServerConnID = conn.GetConnID()
//	sql := fmt.Sprintf("SELECT * FROM accounts where LoginName = '%s'; ", loginAccount.Username)
//	param.ReflectObj = &dcdbmodel.DB_Req_LoginAccount{
//		Username:   loginAccount.GetUsername(),
//		Password:   loginAccount.GetPassword(),
//		ClientType: loginAccount.GetClientType(),
//		ClientIp :  loginAccount.GetClientIp(),
//		MacAddr:    loginAccount.GetMacAddr(),
//	}
//	param.QueryObj = new(dbmodels.Accounts)
//	dbsys.GameDB.AsyncRowToStructQuery(param,OnDBLoginAccount,sql)
//}
//
//
//func OnDBLoginAccount (dbParam *dbsys.DBEventParam) {
//	if dbParam == nil ||  dbParam.QueryObj == nil {
//		return errors.New("OnDBLoginAccount 投递的参数错误")
//	}
//	accountData,ok := dbParam.QueryObj.(*dbmodels.Accounts)
//	if !ok || accountData == nil {
//		dbsys.PDBParamPool.Recycle(dbParam)
//		return xutil.SprintfAssertObjErro("*dbmodels.Accounts")
//	}
//	loginAccount,isok := dbParam.ReflectObj.(*dcdbmodel.DB_Req_LoginAccount)
//	if !isok || loginAccount == nil {
//		dbsys.PDBParamPool.Recycle(dbParam)
//		return xutil.SprintfAssertObjErro("*dcdbmodel.DB_Req_LoginAccount")
//	}
//
//	respn := &protomsg.DC_LS_RespnLoginAccoutMsg{
//		ClientConnID: dbParam.ClientConnID,
//		Username:     loginAccount.Username,
//		AccountID:    accountData.AccountID,
//		ClientType:   loginAccount.ClientType,
//	}
//	if accountData.AccountID == 0 { //没有数据
//		respn.RestCode = msgcode.AccountCode_NotExsist
//		erro := PdataCenter.tcpserver.WritePBMsgByConnID(dbParam.ServerConnID, Cmd.CMD_Account, Cmd.DC_LS_LoginAccount,respn)
//		dbsys.PDBParamPool.Recycle(dbParam)
//		return erro
//	}
//
//	//有数据比较账号密码
//	if strings.Compare(accountData.LoginPwd,loginAccount.Password) != 0 {  //密码错误
//		respn.RestCode = msgcode.AccountCode_PassWordError
//		erro := PdataCenter.tcpserver.WritePBMsgByConnID(dbParam.ServerConnID, Cmd.CMD_Account, Cmd.DC_LS_LoginAccount,respn)
//		dbsys.PDBParamPool.Recycle(dbParam)
//		return erro
//	}
//	//更新内存信息
//	accountData.LoginPwd = loginAccount.Password
//	accountData.LoginIp = &loginAccount.ClientIp
//	nowTimeStr := timeutil.GetTimeALLStr(timeutil.GetTimeNow())
//	accountData.LoginTime = &nowTimeStr
//	accountData.LoginMacAddr = &loginAccount.MacAddr
//	//更新账号信息并且记录日志
//	UpdateAccountsLoginToDBAndRecored(accountData,loginAccount.ClientType)
//
//	//向缓存中添加数据
//	paccount := PaccountMgr.AddAccunts(accountData)
//	if paccount != nil {
//		paccount.SetAccountState(dcmodel.AccountState_Online)
//		paccount.SetClientConnID(dbParam.ClientConnID)
//		paccount.SetAccountServerInfo(PServerInfoMgr.GetServerInfoByConnID(dbParam.ServerConnID)) //获取登录服的信息
//		xlog.Debug("%v  paccount=== %v",osutil.GetRuntimeFileAndLineStr(0), paccount)
//	}
//
//	//登录成功
//	respn.RestCode = msgcode.AccountCode_Login_Succeed
//	msg1,_ := PdataCenter.tcpserver.CreatePBMsg(Cmd.CMD_Account, Cmd.DC_LS_LoginAccount,respn)
//	msg2 := GetGateWayServerInfoByte(accountData.AccountID)
//	erro := PdataCenter.tcpserver.WriteMoreMsgByConnID(dbParam.ServerConnID,msg1,msg2)
//	dbsys.PDBParamPool.Recycle(dbParam)
//	return erro
//}

//////////////////////////////////账号登录结束//////////////////////////////////




//func GetGateWayServerInfoByte(accountID uint64) []byte {
//	//serverInfo := PServerInfoMgr.GetGateWayServerInfo()
//	//if serverInfo == nil {
//	//	return nil
//	//}
//	//sendMsg := &protomsg.ServerInfoToUserMsg{
//	//	Serverinfo : &protomsg.ServerInfoMsg{
//	//		AppId: serverInfo.AppId,
//	//		AppKind: serverInfo.AppKind,
//	//		OutAddr:serverInfo.OutAddr,
//	//		OutPort: serverInfo.OutPort,
//	//	},
//	//	Account_ID: accountID,
//	//}
//	//Account,erro := PdataCenter.tcpserver.CreatePBMsg(Cmd.CMD_DataCenter, Cmd.Sub_DC_LS_GateWayInfo,sendMsg)
//	//if erro != nil {
//	//	xlog.Error("%v,错误 =%v",osutil.GetRuntimeFileAndLineStr(0),erro)
//	//	return nil
//	//}
//	//return Account
//}