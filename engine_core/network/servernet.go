/*
创建时间: 2021/2/25 21:21
作者: zjy
功能介绍:

*/

package network

import "sync"

//其他服务器触发的相关事件
//提供此接口的主要目的是区分服务器与服务器之间交互 与 服务器与客户端的交互
type ServiceNetEvent interface {
	OnServiceLink(conn *ServiceConn)
	OnServerMsg(conn *ServiceConn,maincmd uint32, msg []byte)
	OnServiceClose(conn *ServiceConn)
	CMDIsError(maincmd uint32) error
}


var(
	//对象池
	svConnectPool sync.Pool
	svClosePool sync.Pool
	svMsgPool  sync.Pool
)

func init()  {
	svConnectPool.New = func() interface{} {
		return new(ServerConnectEvent)
	}
	svClosePool.New = func() interface{} {
		return new(ServerClose)
	}
	svMsgPool.New = func() interface{} {
		return new(ServerMsg)
	}
}

func newServerConnectEvent()*ServerConnectEvent {
	data := svConnectPool.Get()
	cl,ok := data.(*ServerConnectEvent)
	if !ok {
		return nil
	}
	return cl
}

func newServerCloseEvent()*ServerClose {
	data := svClosePool.Get()
	cl,ok := data.(*ServerClose)
	if !ok {
		return nil
	}
	return cl
}


type ServerConnectEvent struct {
	observe ServiceNetEvent
	conn *ServiceConn
}

//队列调度
func (this *ServerConnectEvent) Execute(){
	//向监听消息者投递
	this.observe.OnServiceLink(this.conn)
	svConnectPool.Put(this)
}

func (this *ServerConnectEvent)EvenName() string {
	return "ServerConnect"
}

type ServerClose struct {
	observe ServiceNetEvent
	conn *ServiceConn
}


//队列调度
func (this *ServerClose) Execute(){
	this.observe.OnServiceClose(this.conn)
	svClosePool.Put(this)
}

func (this *ServerClose)EvenName() string {
	return "ServerClose"
}



func newServerMsg(conn *ServiceConn,observe ServiceNetEvent,maincmd uint32, msg []byte)*ServerMsg {
	data := svMsgPool.Get()
	sv,ok := data.(*ServerMsg)
	if !ok {
		return nil
	}
	sv.conn = conn
	sv.observe = observe
	sv.cmdID = maincmd
	sv.msg = msg
	return sv
}



type ServerMsg struct {
	observe ServiceNetEvent
	conn *ServiceConn
	cmdID uint32
	msg   []byte
}

//队列调度
func (this *ServerMsg) Execute(){
	//向监听消息者投递
	this.observe.OnServerMsg(this.conn,this.cmdID,this.msg)
	//回收对象
	svMsgPool.Put(this)
}

func (this *ServerMsg)EvenName() string {
	return "ServerMsg"
}



