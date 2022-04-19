/*
创建时间: 2020/12/20 0:10
作者: zjy
功能介绍:

*/

package network

import (
	"github.com/zjytra/MsgServer/model"
)

var (
	// 连接自增ID给其他模块使用
	connID      *model.AtomicUInt32FlagModel
)

func init() {
	connID = model.NewAtomicUInt32Flag()
}

//生成下一个连接ID
func  nextID() uint32 {
	connID.AddUint32()
	if connID.GetUInt32() == 0 {
		connID.AddUint32()
	}
	return connID.GetUInt32()
}
