/*
创建时间: 2021/2/25 21:19
作者: zjy
功能介绍:

*/

package network

import (
	_ "github.com/zjytra/MsgServer/engine_core/dispatch"
	"sync"
)

//都改成传地址高效点吧
//一般是客户端的网络事件向外传递
type ClientNetEvent interface {
	//可能各个服务器的客户端连接对象不一样所以要有一个创建
	OnNetWorkConnect(conn Conner)
	OnClientMsg(conn Conner,maincmd uint32, msg []byte)
	OnNetWorkClose(conn Conner)
	CMDIsError(maincmd uint32) error
}



var(
	//对象池
	clConnectPool sync.Pool
	clClosePool sync.Pool
	clMsgPool sync.Pool
)

func init()  {
	clConnectPool.New = func() interface{} {
		return new(ClientConnect)
	}
	clClosePool.New = func() interface{} {
		return new(ClientClose)
	}
	clMsgPool.New = func() interface{} {
		return new(ClientMsg)
	}
}

func newClientConnect()*ClientConnect {
 	data := clConnectPool.Get()
	cl,ok := data.(*ClientConnect)
	if !ok {
		return nil
	}
	return cl
}

func newClientClose()*ClientClose {
	data := clClosePool.Get()
	cl,ok := data.(*ClientClose)
	if !ok {
		return nil
	}
	return cl
}

//客户端连接
type ClientConnect struct {
	clientEvent ClientNetEvent
	conn Conner
}

//队列调度
func (this *ClientConnect) Execute(){
	//向监听消息者投递
	this.clientEvent.OnNetWorkConnect(this.conn)
	clConnectPool.Put(this)
}

func (this *ClientConnect)EvenName() string {
	return "ClientConnect"
}

//客户端关闭
type ClientClose struct {
	clientEvent ClientNetEvent
	conn Conner
}

//队列调度
func (this *ClientClose) Execute(){
	//向监听消息者投递
	this.clientEvent.OnNetWorkClose(this.conn)
	//回收对象
	clClosePool.Put(this)
}

func (this *ClientClose)EvenName() string {
	return "ClientClose"
}


func newClientMsg(conn Conner,clientEvent ClientNetEvent,maincmd uint32, msg []byte)*ClientMsg {
	data := clMsgPool.Get()
	cl,ok := data.(*ClientMsg)
	if !ok {
		return nil
	}
	cl.conn = conn
	cl.clientEvent = clientEvent
	cl.cmdID = maincmd
	cl.msg = msg
	return cl
}

//客户端消息
type ClientMsg struct {
	clientEvent ClientNetEvent
	conn Conner
	cmdID uint32
	msg   []byte
}

//队列调度
func (this *ClientMsg) Execute(){
	//向监听消息者投递
	this.clientEvent.OnClientMsg(this.conn,this.cmdID,this.msg)
	//回收对象
	clMsgPool.Put(this)
}

func (this *ClientMsg)EvenName() string {
	return "ClientMsg"
}





