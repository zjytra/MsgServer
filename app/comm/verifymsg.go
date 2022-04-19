/*
创建时间: 2021/3/4 22:09
作者: zjy
功能介绍:

*/

package comm

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
)

var(
	AccountLogin map[int64]int64//账号登录验证
)

func InitMsgVerify()  {
	AccountLogin = make(map[int64]int64)
}

//查看是否正在登录中防止客户端重复发请求
func IsLoginQuery(accId int64) bool {
	_, isOk := 	AccountLogin[accId]
	if !isOk {
		return false
	}
	return true
}
//设置正在登录的标志
func SetLogin(accId int64) {
	//保证
	AccountLogin[accId] = timeutil.GetCurrentTimeS()
}
//删除数据
func DelLogin(accId int64) {
	delete(AccountLogin,accId)
}

func ReleaseData(){
	for s, _ := range AccountLogin {
		delete(AccountLogin,s)
	}
	AccountLogin = nil
}