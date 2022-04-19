/*
创建时间: 2021/7/3 23:29
作者: zjy
功能介绍:

*/

package network

import (
	"errors"
	"fmt"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

// 处理消息的函数对象
type HandleMsg func(session Conner, msgdata []byte)



//登录服会话管理
type ConnMgr struct {
	sessionS map[uint32]Conner
	//消息处理者
	msgHandler map[uint32]HandleMsg
}

func (this *ConnMgr) InitSessionMgr() {
	this.sessionS = make(map[uint32]Conner)
	this.msgHandler = make(map[uint32]HandleMsg)
}

func (this *ConnMgr) GetConn(connid uint32) Conner {
	session, ok := this.sessionS[connid]
	if !ok {
		return nil
	}
	return session
}


func (this *ConnMgr) AddConn(conn Conner) {
	_, ok := this.sessionS[conn.GetConnID()]
	if ok {
		xlog.Debug("AddClientSession is exist %d", conn.GetConnID())
		return
	}
	this.sessionS[conn.GetConnID()] = conn
}

//移除连接
func (this *ConnMgr) RemoveConn(connid uint32) {
	_, ok := this.sessionS[connid]
	if !ok {
		return
	}
	delete(this.sessionS, connid)
}

func (this *ConnMgr) CMDIsError(maincmd uint32) error {
	if !this.IsRegister(maincmd) {
		return errors.New(fmt.Sprintf("CMDIsErro 主命令%d 未注册", maincmd))
	}
	return nil
}

func (this *ConnMgr) RegisterHandle(_cmd uint32, handler HandleMsg) {
	if this.IsRegister(_cmd) {
		xlog.Debug("主命令 = %d 已经注册过", _cmd)
		return
	}
	this.msgHandler[_cmd] = handler
}

func (this *ConnMgr) IsRegister(_cmd uint32) bool {
	_, ok := this.msgHandler[_cmd]
	return ok
}

func (this *ConnMgr) GetRegisterCmdS() map[uint32]HandleMsg {
	return this.msgHandler
}

func (this *ConnMgr) Release() {
	this.msgHandler = nil
}

func (this *ConnMgr) GetMsgHandle(cmdID uint32) HandleMsg {
	handle, ok := this.msgHandler[cmdID]
	if !ok {
		xlog.Error("命令 = %d 未注册", cmdID)
		return nil
	}
	return handle
}


