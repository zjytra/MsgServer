/*
创建时间: 2020/7/14
作者: zjy
功能介绍:
客户端连接对象
//这里需要验证
1.自己主动关闭连接流程 向写的通道里写空,在写的线程关闭连接
2.远端关闭，被动关闭
3.远端直接关机
*/

package network

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
	"github.com/zjytra/MsgServer/msgcode"
	"github.com/zjytra/MsgServer/protomsg"
	"net"
)



type TcpConn struct {
	Connect               //继承 连接对象
	clnEv ClientNetEvent
	recMsgPs    int          // 每秒收包个数
}


func newClientConn(conn net.Conn, connId uint32, netOb ClientNetEvent, appid int32, msgParser *MsgParser) *TcpConn {
	if netOb == nil {
		xlog.Error("svnet 网络处理接口为空")
		return nil
	}
	if msgParser == nil {
		xlog.Error("msgParser 数据解析对象为null")
		return nil
	}
	clientConn := new(TcpConn)
	// tcp的方法
	erro := clientConn.initConnData(conn, connId, appid, msgParser, clientConn.OnClose)
	if erro != nil {
		xlog.Error("newServiceConn %v", erro)
		return nil
	}
	clientConn.clnEv = netOb
	return clientConn
}

func (this *TcpConn) setTcpConn(conn net.Conn, connId uint32,netOb ClientNetEvent,appID int32,msgParser  *MsgParser) *TcpConn {
	if netOb == nil {
		xlog.Error("netOb 网络处理接口为空")
		return nil
	}
	if msgParser == nil {
		xlog.Error("msgParser 数据解析对象为null")
		return nil
	}
	// tcp的方法
	erro := this.initConnData(conn,connId,appID,msgParser,this.OnClose)
	if erro != nil {
		xlog.Error("newClientConn %v", erro)
		return nil
	}
	this.clnEv = netOb
	//erro = netOb.OnNetWorkConnect(clntConn) // 投递到事件队列中去
	//if erro != nil {
	//	xlog.Error("newClientConn %v", erro)
	//	return nil
	//}
	return this
}

//通知逻辑层连接成功单独写一个方法是因为 参数是TcpConn子类,逻辑层类型断言时好使用
func (this *TcpConn) notifyConnect(conn Conner) {
	connEvent := newClientConnect()
	connEvent.conn = conn
	connEvent.clientEvent = this.clnEv
	dispatch.AddEventToQueue(connEvent)
}

// 用解析对象读,单协程调用
func (this *TcpConn) ReadMsg() error {

	data, err := this._msgParser.DoRead(this)
	if err != nil {
		return err
	}
	// 查看连接每秒发多少
	if  !this.checkCanRead() {
		return errors.New("每秒超过最大包数")
	}
	maincmd, msgdata, erro := this._msgParser.UnpackOne(data)
	if erro != nil {
		return erro
	}
	//設置當前收包时间
	this.lastRecTime = timeutil.GetCurrentTimeS()
	//查看命令是否错误
	err = this.clnEv.CMDIsError(maincmd)
	if err != nil {
		return err
	}
	// 投递到主逻辑线程去处理
	msgEvent := newClientMsg(this,this.clnEv,maincmd,msgdata)
	dispatch.AddEventToQueue(msgEvent)
	//this.clnEv.OnClientMsg(this,maincmd,msgdat)
	return nil
}

////发送到主逻协程
//func (this *TcpConn) SendMainLogic(maincmd uint32, msgdat []byte) error {
//	//查看命令是否错误
//	erro := this.Connect.CMDIsErro(maincmd)
//	if erro != nil {
//		return erro
//	}
//	// TODO 优化对象创建
//	msgData := NewMsgData(this, maincmd, msgdat)
//	if msgData == nil {
//		return errors.New("NewMsgData is nil")
//	}
//	// 这里应该进入队列
//	return dispatch.MainQueue.AddEvent(msgData)
//}
//
//
////网关服务器客户端消息
//func (this *TcpConn) doGateWayMsg(maincmd uint32, msg []byte) error  {
//	//如果是网关 需要转发数据 转发消息还需要确定转发到哪个服务器
//	_,ok := ClientToDcMsg[maincmd]
//	if ok {
//		//给数据中心转发数据
//		//必须保证不为nil
//		//需要处理
//		transMsge := this.GetTransmitReq(msg)
//		//客户端的主命令
//		return 	DCClient.WritePBMsg(maincmd,transMsge)
//	}
//
//	//游戏服就需要区分哪个游戏服了
//	_,okGame := ClientToGameMsg[maincmd]
//	if okGame {
//		//transMsge := this.GetTransmitReq(msg)
//
//		return nil
//	}
//
//	return this.SendMainLogic(maincmd, msg)
//}

// 查看是否可以读取
func (this *TcpConn) checkCanRead() bool {
	
	if !this.isCheck() { // 查看是否检查每秒包量
		this.lastRecTime = timeutil.GetCurrentTimeS() //这里要设置下时间
		return  true
	}
	cfg := this.GetNetCfg()
	if cfg != nil {
		return true
	}
	currentTime := timeutil.GetCurrentTimeS()

	// 在同一秒内
	if this.lastRecTime == currentTime {
		// 收包的数量超过每秒最大限制数量
		if this.recMsgPs >= cfg.Max_rec_msg_ps {
			xlog.Error("每秒收包 %v 个超过 %v个",this.recMsgPs,cfg.Max_rec_msg_ps)
			return false
		}
		this.recMsgPs ++
		return true
	}
	// 过了一秒重置变量
	this.recMsgPs = 0
	this.lastRecTime = currentTime
	return true
}

//查看是否检查每秒包量
//服务器内部通信不需要检测
func (this *TcpConn) isCheck() bool {
	cfg := this.GetNetCfg()
	if cfg != nil {
		return false
	}
	return cfg.App_kind != model.APP_DataCenter && cfg.App_kind != model.APP_GameServer
}


// 写单个消息
func (this *TcpConn) WritePBMsg(maincmd uint32, pb proto.Message){
	data, erro := this.GetPBByteArr(maincmd, pb)
	if erro != nil {
		xlog.Debug( "WritePBMsg erro : %v", erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func  (this *TcpConn) GetPBByteArr(maincmd uint32,pb proto.Message) (sendMsg []byte, erro error) {
	if !xutil.InterFaceIsNil(pb) {
		sendMsg, erro = proto.Marshal(pb)
	}
	if erro != nil {
		xlog.ErrorLog(appdata.GetSceneName(), "GetPBByteArr %v", erro)
		return nil,erro
	}
	sendMsg, erro = this._msgParser.PackOne(maincmd, sendMsg)
	return sendMsg, erro
}

// 写单个消息
func (this *TcpConn) WriteOneMsg(maincmd uint32, msg []byte) {
	data, erro := this.GetOneMsgByteArr(maincmd, msg)
	if erro != nil {
		xlog.Debug( "WriteOneMsg erro : %v", erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}


// 将消息体构建为[]byte数组，最终要发出去的单包
func  (this *TcpConn) GetOneMsgByteArr(maincmd uint32, msg []byte) ([]byte, error) {
	return this._msgParser.PackOne(maincmd, msg)
}





// 服务器 连接实现的不一样的
func (this *TcpConn) OnClose() {
	xlog.Debug( "TcpConn OnClose")
	// 通知其他模块
	//this.clnEv.OnNetWorkClose(this)
	event := newClientClose()
	event.conn = this
	event.clientEvent = this.clnEv
	dispatch.AddEventToQueue(event)
}



//发送最终构建 MsgRes
func (this *TcpConn) WritePBToMsgRes(maincmd uint32, pb proto.Message) {
	data, erro := this.createMsgResByteArr(maincmd, msgcode.Succeed, pb)
	if erro != nil {
		xlog.Debug("BaseSession WritePBByMsgRes maincmd %d erro : %v", maincmd, erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

//发送protobuf消息 直接发送的MsgRes
func (this *TcpConn) WriteMsgRes(maincmd uint32, pb *protomsg.MsgRes)  {
	data, err := this.GetPBByteArr(maincmd, pb)
	if err != nil {
		xlog.Debug("%v", err)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}

//发送返回结果码
func (this *TcpConn) WritePBMsgAndCode(maincmd uint32, resCode uint32, pb proto.Message) {
	data, erro := this.createMsgResByteArr(maincmd, resCode, pb)
	if erro != nil {
		xlog.Debug("TcpConn WritePBMsgAndCode maincmd %d erro : %v", maincmd, erro)
		return
	}
	// 向写通道投递数据
	this.Write(data)
}



func (this *TcpConn) GetMsgResCodeArr(maincmd uint32, resCode uint32, pb proto.Message) ([]byte, error) {
	return this.createMsgResByteArr(maincmd, resCode, pb)
}


func (this *TcpConn) createMsgResByteArr(maincmd uint32, resCode uint32, pb proto.Message) ([]byte, error) {
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