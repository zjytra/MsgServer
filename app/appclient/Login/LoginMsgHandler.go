/*
创建时间:
作者: zjy
功能介绍:
消息处理自动生成
下面生成的方法可以随便改,自定义的方法不要包含2不然会被过滤掉
*/

package Login

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/engine_core/network"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/protomsg"
)


//连接登录服成功后登录服回复
func L2C_LoginConnectAck(pSession network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_LoginConnectAck{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_LoginConnectAck 错误:" + erro.Error())
		return
	}
}

func L2C_LoginMsg(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_LoginMsg{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_LoginMsg 错误:" + erro.Error())
		return
	}
}

//返回创建的角色id与房间号
func L2C_CreateRoleRes(conn network.Conner, msgdata []byte) {
	reqMsg := &protomsg.L2C_CreateRoleRes{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Debug("L2C_CreateRoleRes 错误:" + erro.Error())
		return
	}
}