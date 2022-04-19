/*
创建时间: 2020/7/5
作者: zjy
功能介绍:
服务器连接对象 主要与客户端端之间的连接进行区分
复用 连接部分方法
重写连接建立 通知不同的接口
重写读取数据 服务器内部通讯不加密
重写连接关闭
*/

package network

import (
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/msgcode"
	"github.com/zjytra/MsgServer/protomsg"
	"net"
)

//复用连接部分方法
type ServiceConn struct {
	Connect
	svNet ServiceNetEvent //服务器相关的网络事件
	connKind int32
	appID int32
	appGroup int32
}

func newServiceConn(conn net.Conn, connId uint32, svnet ServiceNetEvent, appid int32, msgParser *MsgParser) *ServiceConn {
	if svnet == nil {
		xlog.Error("svnet 网络处理接口为空")
		return nil
	}
	if msgParser == nil {
		xlog.Error("msgParser 数据解析对象为null")
		return nil
	}
	serviceConn := new(ServiceConn)
	// tcp的方法
	erro := serviceConn.initConnData(conn, connId, appid, msgParser, serviceConn.OnClose)
	if erro != nil {
		xlog.Error("newServiceConn %v", erro)
		return nil
	}
	serviceConn.svNet = svnet
	event := newServerConnectEvent()
	event.conn = serviceConn
	event.observe = serviceConn.svNet
	dispatch.AddEventToQueue(event)
	//erro = svnet.OnServiceLink(serviceConn)// 通知其他模块已经连接 这里要用服务器的通知模块
	//if erro != nil {
	//	xlog.Error( "newServiceConn %v", erro)
	//	return nil
	//}
	return serviceConn
}

//为了不同客户端连接定义
func (this *ServiceConn) setTcpConn(conn net.Conn, connId uint32,netOb ClientNetEvent,appID int32,msgParser  *MsgParser) *TcpConn {

	return nil
}

// 服务器 连接实现的不一样的
func (this *ServiceConn) OnClose() {
	xlog.Debug("ServiceConn 关闭服务器连接")
	// 通知其他模块
	//erro := this.svNet.OnServiceClose(this)
	//if erro != nil {
	//	xlog.Debug( "doClose OnServiceClose %v", erro)
	//}
	//弄到主逻辑线程去
	event := newServerCloseEvent()
	event.conn = this
	event.observe = this.svNet
	dispatch.AddEventToQueue(event)
}

// 服务器读取 不用加密
func (this *ServiceConn) ReadMsg() error {
	// 查看连接每秒发多少
	data, err := this._msgParser.DoRead(this)
	if err != nil {
		return err
	}
	maincmd, msgdata, erro := this._msgParser.UnpackOne(data)
	if erro != nil {
		return erro
	}
	//設置當前收包时间
	this.lastRecTime = timeutil.GetCurrentTimeS()
	//查看命令是否错误
	erro = this.svNet.CMDIsError(maincmd)
	if erro != nil {
		return erro
	}
	// 投递到主线程去处理
	msgEvent := newServerMsg(this, this.svNet, maincmd, msgdata)
	dispatch.AddEventToQueue(msgEvent)
	//return this.svNet.OnServerMsg(this,maincmd,msgdat)
	return nil
}

// 写单个消息
func (this *ServiceConn) WritePBMsg(maincmd uint32, pb proto.Message)  {
	data, erro := this.GetPBByteArr(maincmd, pb)
	if erro != nil {
		xlog.Debug("WritePBMsg erro : %v", erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}


// 将消息体构建为[]byte数组，最终要发出去的单包
func (this *ServiceConn) GetPBByteArr(maincmd uint32, pb proto.Message) (sendMsg []byte, erro error) {
	if !xutil.InterFaceIsNil(pb) {
		sendMsg, erro = proto.Marshal(pb)
	}
	if erro != nil {
		xlog.ErrorLog(appdata.GetSceneName(), "GetPBByteArr %v", erro)
		return nil, erro
	}
	sendMsg, erro = this._msgParser.PackOne(maincmd, sendMsg)
	return
}

// 写单个消息
func (this *ServiceConn) WriteOneMsg(maincmd uint32, msg []byte) {
	data, erro := this.GetOneMsgByteArr(maincmd, msg)
	if erro != nil {
		xlog.Debug("WriteOneMsgByConnID erro : %v", erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func (this *ServiceConn) GetOneMsgByteArr(maincmd uint32, msg []byte) ([]byte, error) {
	return this._msgParser.PackOne(maincmd, msg)
}

func (this *ServiceConn) GetConnKind() int32 {
	return this.connKind
}

func (this *ServiceConn) SetConnKind(kind int32) {
	this.connKind = kind
}


func (this *ServiceConn) SetAppID(appid int32,group int32) {
	this.appID = appid
	this.appGroup = group
}

func (this *ServiceConn) GetAppID() int32 {
	return 	this.appID
}

func (this *ServiceConn) GetGroupID() int32 {
	return 	this.appGroup
}


//将客户端的消息转发至
func (this *ServiceConn) WriteClientTransReq(maincmd uint32, clnmsg []byte, pClConnId uint32)  {
	sendCMsg := &protomsg.MsgTransReq{}
	sendCMsg.FromConId = pClConnId
	sendCMsg.PbMsg = clnmsg
	sendCMsg.MainCmd = maincmd
	data, err := this.GetPBByteArr(Cmd.CmdMsgTransReq, sendCMsg)
	if err != nil {
		xlog.Debug("WriteClientTransReq GetPBByteArr maincmd = %d %v",maincmd, err)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

//将客户端的消息转发至
func (this *ServiceConn) WriteToClientTransRes(maincmd uint32, pb proto.Message, pClConnId ...uint32) {
	this.WriteToClientTransResAndCode(maincmd, msgcode.Succeed, pb, pClConnId...)
}

//将客户端的消息转发至
func (this *ServiceConn) WriteToClientTransResAndCode(maincmd uint32, resCode uint32, pb proto.Message, pClConnId ...uint32) {
	erro := this.CheckConn(maincmd)
	if erro != nil {
		xlog.Debug("%v", erro)
		return
	}
	var msg []byte
	if pb != nil {
		msg, erro = proto.Marshal(pb)
		if erro != nil {
			xlog.Debug("%v", erro)
			return
		}
	}
	sendCMsg := &protomsg.MsgTransRes{}
	sendCMsg.MainCmd = maincmd
	sendCMsg.SendTo = pClConnId
	sendCMsg.ToClientMsg = new(protomsg.MsgRes)
	sendCMsg.ToClientMsg.ResCode = resCode
	sendCMsg.ToClientMsg.PbMsg = msg
	data, err := this.GetPBByteArr(Cmd.CmdMsgTransRes, sendCMsg)
	if err != nil {
		xlog.Debug("%v", erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}


//发送最终构建 MsgRes
func (this *ServiceConn) WritePBToMsgRes(maincmd uint32, pb proto.Message) {
	data, erro := this.createMsgResByteArr(maincmd, msgcode.Succeed, pb)
	if erro != nil {
		xlog.Debug("BaseSession WritePBByMsgRes maincmd %d erro : %v", maincmd, erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

//发送protobuf消息 直接发送的MsgRes
func (this *ServiceConn) WriteMsgRes(maincmd uint32, pb *protomsg.MsgRes)  {
	data, err := this.GetPBByteArr(maincmd, pb)
	if err != nil {
		xlog.Debug("%v", err)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

//发送返回结果码
func (this *ServiceConn) WritePBMsgAndCode(maincmd uint32, resCode uint32, pb proto.Message) {
	data, erro := this.createMsgResByteArr(maincmd, resCode, pb)
	if erro != nil {
		xlog.Debug("TcpConn WritePBMsgAndCode maincmd %d erro : %v", maincmd, erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}



func (this *ServiceConn) createMsgResByteArr(maincmd uint32, resCode uint32, pb proto.Message) ([]byte, error) {
	erro := this.CheckConn(maincmd)
	if erro != nil {
		return nil, erro
	}
	var msg []byte
	if pb != nil {
		msg, erro = proto.Marshal(pb)
		if erro != nil {
			return nil, erro
		}
	}
	sendCMsg := &protomsg.MsgRes{}
	sendCMsg.ResCode = resCode
	sendCMsg.PbMsg = msg
	return this.GetPBByteArr(maincmd, sendCMsg)
}