/*
创建时间: 2021/8/17 22:17
作者: zjy
功能介绍:

*/

package AccountMgr

import (
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
)



var(
	asyncLogin  map[string]int64                  //检测是否异步登录中
	userNameAccount map[string]*dbmodels.AccountT //根据username判断
)



func Init() {
	userNameAccount = make(map[string]*dbmodels.AccountT)
	asyncLogin = make(map[string]int64)
}


//获取账号信息
func  GetAccountInfoByUserName(userName string) *dbmodels.AccountT {
	return userNameAccount[userName]
}

//获取账号信息
func  AddAccount(pacc *dbmodels.AccountT){
	accName := pacc.LoginName.GetVal()
	_,ok := userNameAccount[accName]
	if ok {
		return
	}
	userNameAccount[accName] = pacc
}


//查看是否正在登录中防止客户端重复发请求
func IsAsyncLogin(account string) bool {
	_, isOk := 	asyncLogin[account]
	if !isOk {
		return false
	}
	return true
}
//设置正在登录的标志
func SetAsyncLogin(account string) {
	//保证
	asyncLogin[account] = timeutil.GetCurrentTimeS()
}
//删除数据
func DelAsyncLogin(account string) {
	delete(asyncLogin,account)
}


//获取账号信息
func  DelAccount(account string){
	_,ok := userNameAccount[account]
	if ok {
		delete(userNameAccount,account)
		return
	}
}


func  ReleaseData(){
	userNameAccount = nil
	asyncLogin = nil
}