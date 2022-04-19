/*
创建时间:
作者: zjy
功能介绍:
消息处理自动生成
下面生成的方法可以随便改,自定义的方法不要包含2不然会被过滤掉
*/

package MsgPro

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/protomsg"
)


//发送消息
func L2C_SendMsgRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_SendMsgRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_SendMsgRes 错误:" + erro.Error())
		return
	}
}

//获取消息列表回复
func L2C_GetRoomMsgListRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_GetRoomMsgListRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_GetRoomMsgListRes 错误:" + erro.Error())
		return
	}
}

func L2C_PopularMsgRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_PopularMsgRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_PopularMsgRes 错误:" + erro.Error())
		return
	}
}

//切换房间回复
func L2C_SwitchRoomRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_SwitchRoomRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_SwitchRoomRes 错误:" + erro.Error())
		return
	}
}

//获取玩家信息回复
func L2C_GetUserInfoRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_GetUserInfoRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_GetUserInfoRes 错误:" + erro.Error())
		return
	}
}