/*
创建时间: 2020/08/2020/8/23
作者: Administrator
功能介绍:
消息相关回复
*/
package msgcode

import (
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
)


func IsVerify(Username string) uint32 {
	//长度验证
	lencode := VerifyStrLen(Username)
	if lencode != Succeed{
		return lencode
	}
	//账号包含空格或者非单词字符
	isMatch := strutil.StringHasSpaceOrSpecialChar(Username)
	if isMatch {
		return AccountCode_UserNameFormatErro
	}
	//sql注入验证
	isMatch = strutil.StringHasSqlKey(Username)
	if isMatch {
		return AccountCode_SqlZhuRu
	}
	return Succeed
}

//验证长度
func VerifyStrLen(username string) uint32 {
	strLen := len(username)
	if strLen <= 4 {
		return AccountCode_UserNameShort
	}
	if strLen > 11 {
		return AccountCode_UserNameLong
	}

	return Succeed
}