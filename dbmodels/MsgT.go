/*
创建时间: 2022/4/18 22:48
作者: zjy
功能介绍:
消息数据表
*/

package dbmodels

import (
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/protomsg"
)

type MsgT struct {
	dbsys.DBObj
	MsgID  dbsys.DBInt64 `indexTP:"PRIMARY" node:"角色id" QFlag:"2"`
	CreateTime dbsys.DBInt64  `node:"创建时间" `
	RoomID dbsys.DBInt32  `node:"房间id"  indexTP:"MUl"  QFlag:"1"`
	MsgContent  dbsys.DBStr `tp:"varchar(1024)" node:"消息内容" `
	SenderName  dbsys.DBStr `tp:"varchar(32)" node:"发送者名称" `
}




func (this *MsgT)CreateObj() dbsys.DBObJer {
	return new(MsgT)
}

func (this *MsgT)GetUID() int64 {
	return this.MsgID.GetVal()
}

func (this *MsgT)BuildPro(msgPro *protomsg.MsgInfo)  {
	msgPro.RoomID = this.RoomID.GetVal()
	msgPro.Conent = this.MsgContent.GetVal()
	msgPro.RoleName = this.SenderName.GetVal()
	msgPro.SendTime = this.CreateTime.GetVal()
}