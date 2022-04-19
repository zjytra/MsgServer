/*
创建时间: 2021/9/2 23:11
作者: zjy
功能介绍:
服务器session管理
*/

package network

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/MsgServer/Cmd"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
	"github.com/zjytra/MsgServer/msgcode"
	"github.com/zjytra/MsgServer/protomsg"
	"strconv"
)

//转发的消息
type HandleTransMsg func(psSession *ServiceConn, msgdata []byte, msg_form_connid uint32)
type HandleServerMsg func(psSession *ServiceConn, msgdata []byte)

type ServerConnMgr struct {
	sessionS           map[uint32]*ServiceConn  //根据connid 关联
	serverAppIDSession map[int32]*ServiceConn   //根据appid 关联
	groupSession       map[int32][]*ServiceConn //根据group 关联
	//连接其他服务器对象
	serverClient map[int32]*ServiceClient
	//转发消息处理
	msgTransHandler map[uint32]HandleTransMsg
	//处理服务器消息
	msgServerHandler map[uint32]HandleServerMsg
}

func (this *ServerConnMgr) InitSessionMgr() {
	this.sessionS =  make(map[uint32]*ServiceConn)
	this.serverAppIDSession = make(map[int32]*ServiceConn)
	this.groupSession = make(map[int32][]*ServiceConn)
	this.serverClient = make(map[int32]*ServiceClient)
	this.msgTransHandler = make(map[uint32]HandleTransMsg)
	this.msgServerHandler = make(map[uint32]HandleServerMsg)
}


func (this *ServerConnMgr) AddServerConn(conn *ServiceConn) {
	_, ok := this.sessionS[conn.GetConnID()]
	if ok {
		xlog.Debug("AddClientSession is exist %d", conn.GetConnID())
		return
	}
	this.sessionS[conn.GetConnID()] = conn
}

//移除连接
func (this *ServerConnMgr) RemoveConnByConnID(connid uint32) {
	_, ok := this.sessionS[connid]
	if !ok {
		return
	}
	delete(this.sessionS, connid)
}

//添加连接的服务器信息
func (this *ServerConnMgr) addServerConn(pServer *ServiceConn) bool {
	if pServer == nil {
		xlog.Warning("addServerConn is nil ")
		return false
	}
	_, ok := this.serverAppIDSession[pServer.GetAppID()]
	if ok {
		xlog.Debug("addServerConn appid %d 已经添加", pServer.GetAppID())
		return false
	}
	this.serverAppIDSession[pServer.GetAppID()] = pServer
	this.AddServerSessionToGroup(pServer)
	xlog.Debug("addServerConn appid %d", pServer.GetAppID())
	return true
}

func (this *ServerConnMgr) AddServerSessionToGroup(info *ServiceConn) {
	if info == nil {
		return
	}
	infos, ok := this.groupSession[info.GetGroupID()]
	if ok {
		for _, temInfo := range infos {
			if temInfo.GetAppID() == info.GetAppID() {
				xlog.Debug("添加 AddServerSessionToGroup appid %d 服务器已经存在", info.GetAppID())
				return
			}
		}
		infos = append(infos, info)
		this.groupSession[info.GetGroupID()] = infos
		return
	}

	var tem []*ServiceConn
	tem = append(tem, info)
	this.groupSession[info.GetGroupID()] = tem
}

func (this *ServerConnMgr) RemoveSessionByServerSession(info *ServiceConn) {
	if info == nil {
		return
	}
	RemoveServerInfo(info.GetAppID())
	this.RemoveConnByConnID(info.GetConnID())
	infos, ok := this.groupSession[info.GetGroupID()]
	if !ok {
		return
	}
	for i, temInfo := range infos {
		if temInfo.GetAppID() == info.GetAppID() {
			infos[i] = nil
			newInfo := append(infos[:i], infos[i+1:]...)
			this.groupSession[info.GetGroupID()] = newInfo
			return
		}
	}
}

// 移除某个服务器连接
func (this *ServerConnMgr) RemoveServerConnByAppID(appId int32) {
	pInfo, ok := this.serverAppIDSession[appId]
	if !ok {
		return
	}
	this.RemoveSessionByServerSession(pInfo)
	delete(this.serverAppIDSession, appId)
}

//获取服务器信息
func (this *ServerConnMgr) GetServerSession(appId int32) *ServiceConn {
	pInfo, ok := this.serverAppIDSession[appId]
	if !ok {
		return nil
	}
	return pInfo
}

func (this *ServerConnMgr) GetServerConnByConnId(connid uint32) *ServiceConn {
	session, ok := this.sessionS[connid]
	if !ok {
		return nil
	}
	return session
}

//获取服务器信息
func (this *ServerConnMgr) SendGroupServerMsg(group int32, kind int32, maincmd uint32, pb proto.Message) {
	apps := this.GetServerSessionsByGroup(group)
	for _, pSesson := range apps {
		if pSesson == nil {
			continue
		}
		if pSesson.GetConnKind() != kind {
			continue
		}
		pSesson.WritePBMsg(maincmd, pb)
	}
}

//获取某个区的某类型服务器
func (this *ServerConnMgr) GetGroupKindServersSession(group int32, kind int32) []*ServiceConn {
	//获取负载最小的网关
	var servers []*ServiceConn
	apps := this.GetServerSessionsByGroup(group)
	//获取负载最小的
	var minnum *ServiceConn
	for _, info := range apps {
		if info.GetConnKind() != kind {
			continue
		}
		servers = append(servers, minnum)
	}
	return servers
}

//获取某个区的某类型服务器
func (this *ServerConnMgr) GetGroupOneServerSession(group int32, kind int32) *ServiceConn {
	apps := this.GetServerSessionsByGroup(group)
	for _, info := range apps {
		if info.GetConnKind() != kind {
			continue
		}
		return info
	}
	return nil
}

//获取服务器信息
func (this *ServerConnMgr) GetServerSessionsByGroup(group int32) []*ServiceConn {
	infos, ok := this.groupSession[group]
	if !ok {
		return nil
	}
	return infos
}

func (this *ServerConnMgr) OnServerAck(pSession *ServiceConn, appid int32) *SeverInfo {
	return this.OnServerAckAndNum(pSession, appid, 0)
}

//num 玩家数量
func (this *ServerConnMgr) OnServerAckAndNum(pSession *ServiceConn, appid int32, num int32) *SeverInfo {
	pSession.SetContAck()
	// 缓存游戏服
	netCfg := csvdata.GetNetWorkCfgPtr(appid)
	if netCfg != nil {
		pSession.SetConnKind(netCfg.App_kind)
		pSession.SetAppID(appid, netCfg.Group)
		this.addServerConn(pSession)

		return AddServerInfo(appid, ServerStauts_Online, num)
	}
	return nil
}

//处理服务器的消息
func (this *ServerConnMgr) OnServerHandlerMsg(conn *ServiceConn, maincmd uint32, msg []byte) {
	pSession := this.GetServerConnByConnId(conn.GetConnID())
	if pSession == nil {
		xlog.Debug("MonitorSessionMgr OnServerMsg 命令 %d 连接已断开 ", maincmd)
		return
	}
	//分消息处理
	handle := this.GetServerMsgHandle(maincmd)
	if handle == nil {
		return
	}
	//服务器消息就不用再加结果码了
	//pMsg := this.PassMsgRes(msg)
	//if pMsg == nil {
	//	xlog.Debug("解析命令 = %d 消息错误", maincmd)
	//	return
	//}
	//if pMsg.ResCode != msgcode.Succeed {
	//	xlog.Debug("命令 = %d 错误码 %d ", maincmd, pMsg.ResCode)
	//	return
	//}
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	handle(pSession, msg)          //查看是否注册对应的处理函数
	dispatch.CheckTime("命令:"+strconv.FormatInt(int64(maincmd), 10), startT, timeutil.FrameTimeMs)
}

//处理服务器转发至客户端的消息
func (this *ServerConnMgr) OnMsgTransClientReq(conn *ServiceConn, maincmd uint32, msg []byte) {
	pSession := this.GetServerConnByConnId(conn.GetConnID())
	if pSession == nil {
		xlog.Debug("OnMsgTransClientReq 命令 %d 连接已断开 ", maincmd)
		return
	}
	pMsg := this.PassMsgTransReq(msg)
	if pMsg == nil {
		xlog.Debug("OnMsgTransClientReq 解析命令 = %d 消息错误", maincmd)
		return
	}
	handle := this.GetMsgTransHandle(pMsg.MainCmd)
	if handle == nil { //转发消息无处理
		pSession.WriteToClientTransResAndCode(pMsg.MainCmd,msgcode.MsgNotHandler,nil,pMsg.FromConId)
		return
	}
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	//处理客户端转发消息
	handle(pSession, pMsg.PbMsg, pMsg.FromConId)
	dispatch.CheckTime("命令:"+strconv.FormatInt(int64(maincmd), 10), startT, timeutil.FrameTimeMs)
}

func (this *ServerConnMgr) GetServerClient(appid int32) *ServiceClient {
	pServerClient, ok := this.serverClient[appid]
	if ok {
		return pServerClient
	}
	return nil
}

func (this *ServerConnMgr) AddServerClient(pServerClient *ServiceClient) {
	if pServerClient == nil {
		return
	}
	this.serverClient[pServerClient.GetAppId()] = pServerClient
}

func (this *ServerConnMgr) DoCreateServiceClient(svNet ServiceNetEvent, appid int32) {
	dcCfg := csvdata.GetNetWorkCfgPtr(appid)
	if dcCfg == nil {
		xlog.Debug("没找到数据世界服配置%d", appid)
		return
	}
	pServerClient := this.GetServerClient(appid)
	if pServerClient != nil {
		return
	}
	xlog.Debug("创建 %d 服务器连接 ", appid)
	pServerClient = NewServiceClient(svNet, appid)
	if pServerClient == nil {
		return
	}
	pServerClient.Start()
	this.AddServerClient(pServerClient)
}

func (this *ServerConnMgr) GetMyGroupDataCenterSession() *ServiceConn {
	return this.GetGroupOneServerSession(csvdata.OutNetConf.Group, model.APP_DataCenter)
}

func (this *ServerConnMgr) GetMyGroupWorldSession() *ServiceConn {
	return this.GetGroupOneServerSession(csvdata.OutNetConf.Group, model.APP_DataCenter)
}


func (this *ServerConnMgr) IsRegisterTrans(_cmd uint32) bool {
	_, ok := this.msgTransHandler[_cmd]
	return ok
}

func (this *ServerConnMgr) RegisterTransHandle(_cmd uint32, handler HandleTransMsg) {
	if this.IsRegisterTrans(_cmd) {
		xlog.Debug("转发命令 = %d 已经注册过", _cmd)
		return
	}
	this.msgTransHandler[_cmd] = handler
}

func (this *ServerConnMgr) GetMsgTransHandle(cmdID uint32) HandleTransMsg {
	handle, ok := this.msgTransHandler[cmdID]
	if !ok {
		xlog.Error("转发命令 = %d 未注册", cmdID)
		return nil
	}
	return handle
}



func (this *ServerConnMgr) CMDIsError(maincmd uint32) error {
	//转发的命令不需要注册
	if maincmd == Cmd.CmdMsgTransReq {
		return nil
	}
	if !this.IsServerRegister(maincmd) {
		return errors.New(fmt.Sprintf("ServerConnMgr CMDIsErro 主命令%d 未注册", maincmd))
	}
	return nil
}

func (this *ServerConnMgr) RegisterServerHandle(_cmd uint32, handler HandleServerMsg) {
	if this.IsServerRegister(_cmd) {
		xlog.Debug("主命令 = %d 已经注册过", _cmd)
		return
	}
	this.msgServerHandler[_cmd] = handler
}

func (this *ServerConnMgr) IsServerRegister(_cmd uint32) bool {
	_, ok := this.msgServerHandler[_cmd]
	return ok
}


func (this *ServerConnMgr) Release() {
	//this.ConnMgr.Release()
	this.msgTransHandler = nil
	this.msgServerHandler = nil
}

func (this *ServerConnMgr) GetServerMsgHandle(cmdID uint32) HandleServerMsg {
	handle, ok := this.msgServerHandler[cmdID]
	if !ok {
		xlog.Error("命令 = %d 未注册", cmdID)
		return nil
	}
	return handle
}

func (this *ServerConnMgr) GetTransHandleCmdS() []uint32 {
	var cmds []uint32
	for cmd, _ := range this.msgTransHandler {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (this *ServerConnMgr) PassMsgRes(msgdata []byte) *protomsg.MsgRes {
	reqMsg := new(protomsg.MsgRes)
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Error("PassMsgRes erro %v ", erro)
		return nil
	}
	return reqMsg
}

func (this *ServerConnMgr) PassMsgTransReq(msgdata []byte) *protomsg.MsgTransReq {
	reqMsg := new(protomsg.MsgTransReq)
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Error("PassMsgRes erro %v ", erro)
		return nil
	}
	return reqMsg
}

func (this *ServerConnMgr) PassMsgMsgTransRes(msgdata []byte) *protomsg.MsgTransRes {
	reqMsg := new(protomsg.MsgTransRes)
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.Error("PassMsgRes erro %v ", erro)
		return nil
	}
	return reqMsg
}
