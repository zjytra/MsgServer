/*
创建时间: 2019/11/23
作者: zjy
功能介绍:
给其他库提供
*/

package appdata

import (
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/xengine"
	"github.com/zjytra/MsgServer/model"
)

var (
	AppID          int32             //serverId
	appKind        model.AppKind     // app类型 通过外部传递参数确定
	AppHook        xengine.AppHooker //让app下层的模块能调用App方法
)

func InitData(hooker xengine.AppHooker)  {
	AppHook = hooker
}

//根据参数初始化服务器
func InitAppDataByAppArgs(appid int32) {
	AppID = appid
	csvdata.SetAppCfg(appid)
	appKind = model.ItoAppKind(csvdata.OutNetConf.App_kind)

}

func GetAppKind() model.AppKind {
	return appKind
}

func GetSceneName() string {
	switch appKind {
	//gameserver需要区分场景
	case model.APP_GameServer:
		return csvdata.OutNetConf.App_name
	//这些服务器器都没有场景名称
	// case model.APP_NONE,model.APP_Client,model.APP_MsgServer,model.APP_DataCenter:
	// 	return ""
	default:
		return ""
	}
	return ""
}
