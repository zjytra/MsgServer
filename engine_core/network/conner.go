package network

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/protomsg"
	"net"
)

//连接接口
type Conner interface {
	Read(b []byte) (int, error)
	ReadMsg() (error)
	//一次性发送多个消息
	WriteMsg(args ...[]byte)
	WriteOneMsg(maincmd uint32, msg []byte)
	//将命令和内容打成一个包
	GetOneMsgByteArr(maincmd uint32, msg []byte) ([]byte, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteAddrStr() string
	RemoteAddrIp() string
	RemoteAddrPort() string
	Close()
	Destroy()
	//获取连接id
	GetConnID() uint32
	//多个消息转换为一个
	ConnMorePackageToOne(args ...[]byte) ([]byte, error)
	//发送protobuf消息
	WritePBMsg(maincmd uint32, pb proto.Message)
	//发送返回消息
	WritePBToMsgRes(maincmd uint32, pb proto.Message)
	//发送protobuf消息 直接发送的MsgRes
	WriteMsgRes(maincmd uint32, pb *protomsg.MsgRes)
	//发送返回结果码服务器内部发送如果有错误码就会不处理该消息
	WritePBMsgAndCode(maincmd uint32,resCode uint32,pb proto.Message)
	//发送protobuf消息
	GetPBByteArr(maincmd uint32,pb proto.Message)([]byte,  error)
	//客户端确认连接
	SetContAck()
	//获取配置类型
	GetCfgKind() int32
	//下面的方法是统一接口
	Write(b []byte)
	IsAlive() bool
	IsClose() bool
	//为了不同客户端连接定义
	setTcpConn(conn net.Conn, connId uint32,netOb ClientNetEvent,appID int32,msgParser  *MsgParser) *TcpConn
}




