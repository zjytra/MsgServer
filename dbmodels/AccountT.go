/*
创建时间: 2022/4/18 15:57
作者: zjy
功能介绍:
//支持中间删除添加字段,不支持修改表字段名,对象被删除的字段,数据库会删除哦
QFlag 字段作为查询条件的标志2进制位  1 select 2 update与del共用   3表示select与update insert
//账号表
*/

package dbmodels

import (
	"github.com/zjytra/MsgServer/engine_core/dbsys"
)

type AccountT struct {
	dbsys.DBObj
	AccID  dbsys.DBInt64 `indexTP:"PRIMARY" node:"账号id" QFlag:"2"`
	LoginName  dbsys.DBStr `tp:"varchar(32)" indexTP:"UNIQUE" node:"登录账号"  QFlag:"1"`
	RegisterTime dbsys.DBInt64  `node:"注册时间"  `
	RegisterIp dbsys.DBStr `tp:"varchar(32)" node:"注册ip"`
}



func NewAccountT(loginName string) *AccountT {
	pAcc := new(AccountT)
	pAcc.LoginName.SetVal(loginName)
	return pAcc
}

func (this *AccountT)CreateObj() dbsys.DBObJer {
	return new(AccountT)
}

func (this *AccountT)GetUID() int64 {
	return this.AccID.GetVal()
}


func (this *AccountT)GetLoginName() string {
	return this.LoginName.GetVal()
}



