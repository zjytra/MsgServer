/*
创建时间: 2020/09/2020/9/11
作者: Administrator
功能介绍:

*/
package network

import (
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"github.com/zjytra/MsgServer/model"
)

var (
	appsById map[int32]*SeverInfo //<appid,SeverInfo>根据appId关联服务器
	appsByGroups map[int32][]*SeverInfo //<group,[]SeverInfo>把相同或组的服务器放在一起
)

func init() {
	appsById = make(map[int32]*SeverInfo)
	appsByGroups = make(map[int32][]*SeverInfo)
}

//添加连接的服务器信息
func AddServerInfo(appid int32, _status int32, num int32) *SeverInfo {
	Cfg := csvdata.GetNetWorkCfgPtr(appid)
	if Cfg == nil {
		xlog.Warning(xutil.GetTraceStackStr("AddServerInfo 配置错误"))
		return nil
	}
	info, ok := appsById[appid]
	if ok {
		info.SetServerInfo(appid,_status,num)
		info.GroupId = Cfg.Group
		info.IP = Cfg.Out_addr
		info.Port = Cfg.Out_prot
		info.Kind = Cfg.App_kind
		info.Name = Cfg.App_name
	} else {
		info = new(SeverInfo)
		info.SetServerInfo(appid,_status,num)
		info.GroupId = Cfg.Group
		info.IP = Cfg.Out_addr
		info.Port = Cfg.Out_prot
		info.Kind = Cfg.App_kind
		info.Name = Cfg.App_name
		appsById[appid] = info
	}
	AddServerInfoToGroup(info)
	xlog.Debug("添加 appid %d  名称 %s ", appid, Cfg.App_name)
	return info
}

func AddServerInfoToGroup(info *SeverInfo)  {
	if info == nil {
		return
	}
	infos, ok := appsByGroups[info.GroupId]
	if ok {
		for _,temInfo := range infos {
			if temInfo.AppId == info.AppId {
				xlog.Debug("添加 AddServerInfoToGroup appid %d  名称 %s 服务器已经存在", temInfo.AppId, temInfo.Name)
				return
			}
		}
		infos = append(infos,info)
		appsByGroups[info.GroupId] = infos
		return
	}

	var tem []*SeverInfo
	tem = append(tem,info)
	appsByGroups[info.GroupId] = tem
}



func RemoveServerInfoToGroup(info *SeverInfo)  {
	if info == nil {
		return
	}
	infos, ok := appsByGroups[info.GroupId]
	if !ok {
		return
	}
	for i,temInfo := range infos {
		if temInfo.AppId == info.AppId {
			infos[i] = nil
			newInfo := append(infos[:i], infos[i+1:]...)
			appsByGroups[info.GroupId] = newInfo
			return
		}
	}
}

// 移除某个服务器连接
func RemoveServerInfo(appId int32) {
	info, ok := appsById[appId]
	if !ok {
		return
	}
	xlog.Debug("移除服务器 appid %d  名称 %s ", info.AppId, info.Name)
	RemoveServerInfoToGroup(info)
	delete(appsById, appId)
}

//获取服务器信息
func GetServerInfo(appId int32) *SeverInfo {
	pInfo, ok := appsById[appId]
	if !ok {
		return nil
	}
	return pInfo
}

//获取服务器信息
func GetServerInfosByGroup(group int32) []*SeverInfo {
	infos, ok := appsByGroups[group]
	if !ok {
		return nil
	}
	return infos
}

//获取每个区负载最小的
func GetMinNumGateWays() []*SeverInfo{
	//获取负载最小的网关
	var tem []*SeverInfo
	//获取负载最小的
	for _, apps := range appsByGroups { //每组服务器
		var minnum *SeverInfo // 最少人数的服务器
		for _, info := range apps {
			if info.Kind != model.APP_GATEWAY {
				continue
			}
			//// 不是在线的不算
			//if info.Status != ServerStauts_Online {
			//	continue
			//}
			//先设置第一个为最小负载
			if minnum == nil {
				minnum = info
				continue
			}
			if info.Num < minnum.Num {
				minnum = info
			}
		}
		//没满足上面的条件
		if minnum == nil {
			continue
		}
		tem = append(tem, minnum)
	}

	return tem
}


//获取某个区的某类型服务器
func GetGroupKindServer(group int32,kind int32) []*SeverInfo{
	var servers []*SeverInfo
	apps := GetServerInfosByGroup(group)
	//获取负载最小的
	for _, info := range apps {
		if info.Kind != kind {
			continue
		}
		//// 不是在线的不算
		//if info.Status != ServerStauts_Online {
		//	continue
		//}
		servers = append(servers, info)
	}

	return servers
}

//获取某个区的某类型单个服务器
func GetGroupKindOneServer(group int32,kind int32) *SeverInfo{
	apps := GetServerInfosByGroup(group)
	//获取负载最小的
	for _, info := range apps {
		if info.Kind != kind {
			continue
		}
		//// 不是在线的不算
		//if info.Status != ServerStauts_Online {
		//	continue
		//}
		return info
	}
	return nil
}