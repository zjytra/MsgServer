/*
创建时间: 2020/4/28
作者: zjy
功能介绍:
TODO 重新加载需要单独设计
*/

package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/conf"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"time"
)


var (
	csvPath string
	OutNetConf     *NetWorkCfg // 本服务器网络配置在服务器解析参数的时候就获得
	DcNetConf      *NetWorkCfg // 数据服连接配置
	MonitorNetConf *NetWorkCfg // 数据服连接配置
)

func SetCsvPath(csvpath string) {
	if strutil.StringIsNil(csvpath) {
		fmt.Println("csvpath is nil")
	}
	csvPath = csvpath
}

func InitCsv()  {
	SetCsvPath(conf.PathModelPtr.CsvPath)
}

func StartCsv() {
	LoadCommonCsvData() // 读取公共的csv
}


func OnAppClose() {

}

func SetAppCfg(appid int32)  {
	for {
		OutNetConf = GetNetWorkCfgPtr(appid)
		if OutNetConf == nil {
			xlog.Error("serverID 未找到")
		} else {
			break
		}
		time.Sleep(time.Second * 5)
	}

	MonitorNetConf = GetNetWorkCfgPtr(conf.SvJson.MonitorID)
	if MonitorNetConf == nil {
		panic("MonitorNetConf  == nil")
	}
}

//初始化登陆服数据
func LoadCommonCsvData()  {
	SetNetWorkCfgMapData(csvPath)
	SetDbCfgMapData(csvPath)
}


func ReLoadCommonCsvData()  {
	 go LoadCommonCsvData()
}