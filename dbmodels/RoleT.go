/*
创建时间: 2022/4/18 22:48
作者: zjy
功能介绍:
角色表
QFlag 字段作为查询条件的标志2进制位  1 select 2 update与del共用   3表示select与update insert
*/

package dbmodels

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/protomsg"
)

type RoleT struct {
	dbsys.DBObj
	RoleID  dbsys.DBInt64 `indexTP:"PRIMARY" node:"角色id" QFlag:"2"`
	AccID  dbsys.DBInt64 `indexTP:"MUl" node:"账号id" QFlag:"1"`
	RoleName  dbsys.DBStr `tp:"varchar(32)" indexTP:"UNIQUE" node:"角色名"`
	CreateTime dbsys.DBInt64  `node:"创建时间" `
	LonginTime dbsys.DBInt64  `node:"登录时间"`
	OnlineTime dbsys.DBInt64  `node:"在线时长"`
	RoomID dbsys.DBInt32  `node:"房间id" `
	OffLineTime dbsys.DBInt64  `node:"离线时间"`
}



func NewRoleT() *RoleT {
	pAcc := new(RoleT)
	return pAcc
}

func (this *RoleT)CreateObj() dbsys.DBObJer {
	return new(RoleT)
}

func (this *RoleT)GetUID() int64 {
	return this.RoleID.GetVal()
}

func (this RoleT)BuildRolePro(pro *protomsg.L2C_LoginMsg)  {
	pro.RoleID = this.RoleID.GetVal()
	pro.RoleName = this.RoleName.GetVal()
	pro.RoomID = this.RoomID.GetVal()
}

func (this *RoleT)GetOnlineTime() int64  {
	this.OnlineTime.SetVal(timeutil.GetCurrentTimeMs() - this.LonginTime.GetVal())
	return this.OnlineTime.GetVal()
}