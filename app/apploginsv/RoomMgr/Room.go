/*
创建时间: 2022/4/19 21:53
作者: zjy
功能介绍:

*/

package RoomMgr

import "github.com/zjytra/MsgServer/dbmodels"

type Room struct {
	ID  int32
	Msgs []*dbmodels.MsgT
	users map[uint32]int64  //玩家连接id与玩家id
}

func NewRoom(ID int32) *Room {
	 pRoom := new(Room)
	 pRoom.users = make(map[uint32]int64)
	 pRoom.ID = ID
	 return pRoom
}



