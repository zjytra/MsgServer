/*
创建时间: 2022/4/19 21:52
作者: zjy
功能介绍:

*/

package RoomMgr

import (
	"github.com/zjytra/MsgServer/dbmodels"
)

var rooms map[int32]*Room

//目前先创建3个
func Init()  {
	rooms = make(map[int32]*Room)
    var i int32
	for i  = 1; i <= 3; i ++  {
		rooms[i] = NewRoom(i)
	}
}

func AddMsg(t * dbmodels.MsgT)  {
    pRoom := GetRoom(t.RoomID.GetVal())
	if pRoom == nil {
		return
	}
	pRoom.Msgs = append(pRoom.Msgs,t)
}

func GetRoom(id int32) *Room {
	pRoom := rooms[id]
	if pRoom == nil {
		return nil
	}
	return pRoom
}


func OnRoleEnter(connID uint32,pRole * dbmodels.RoleT) {
	pRoom := GetRoom(pRole.RoomID.GetVal())
	if pRoom == nil {
		return
	}
	pRoom.users[connID] = pRole.GetUID()
}


func OnRoleLeave(connID uint32,roomID int32) {
	pRoom := GetRoom(roomID)
	if pRoom == nil {
		return
	}
	delete(pRoom.users,connID)
}

func GetRoomUsers(roomID int32) map[uint32]int64 {
	pRoom := GetRoom(roomID)
	if pRoom == nil {
		return nil
	}
	return pRoom.users
}

func GetRoomMsgPro(idRoomId int32)[]*dbmodels.MsgT  {
	pRoom := GetRoom(idRoomId)
	if pRoom == nil {
		return nil
	}
	if pRoom.Msgs == nil || len(pRoom.Msgs)  == 0{
		return nil
	}
	return  pRoom.Msgs
}