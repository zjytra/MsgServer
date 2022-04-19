/*
创建时间: 2019/12/25
作者: zjy
功能介绍:
路径相关管理
*/

package conf

import (
	"path"
)

type PathInfo struct {
	AppRootPath    string //  程序(main)根路径
	CsvPath        string
	LogsPath       string
	ConfPath       string
	ServerConfPath string
	RedisConfPath string
}

// 创建PathModel
func newPathModel() *PathInfo {
	return new(PathInfo)
}

// 路径管理相关函数
func (pthpro *PathInfo) SetRootPath(pwd string ) {
	pthpro.AppRootPath = pwd
}

func (pthpro *PathInfo) InitPathModel() {
	pthpro.ConfPath = path.Join(pthpro.AppRootPath, "cfgs")
	pthpro.CsvPath = path.Join(pthpro.AppRootPath, "csv")
	pthpro.LogsPath = path.Join(pthpro.AppRootPath, "logs")
	pthpro.ServerConfPath = path.Join(pthpro.ConfPath, "serverconf.json")
	pthpro.RedisConfPath = path.Join(pthpro.ConfPath, "redis.json")
}

