/*
创建时间:
作者: zjy
功能介绍:
消息处理自动生成
下面生成的方法可以随便改,自定义的方法不要包含2不然会被过滤掉
*/

package Role

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/protomsg"
)


//返回创建结果
func D2C_CreateRoleRes(pSession network.Conner, msgdata []byte) {
	reqMsg := &protomsg.D2C_CreateRoleRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("D2C_CreateRoleRes 错误:" + erro.Error())
		return
	}
}

//返回创建结果
func D2C_SelectRoleEnterGameRes(pSession network.Conner, msgdata []byte) {
	reqMsg := &protomsg.D2C_SelectRoleEnterGameRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("D2C_SelectRoleEnterGameRes 错误:" + erro.Error())
		return
	}
}