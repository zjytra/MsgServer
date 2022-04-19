/*
创建时间: 2020/2/6
作者: zjy
功能介绍:
配置相关功能
*/

package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerJson struct {
	LogConf    VolatileLogModel `json:"LogConf"`
	DCServerID int32            `json:DCServerID`
	MonitorID  int32            `json:MonitorID`
	ResUrl     string           `json:"ResUrl"`
	InnerKey   string           `json:"InnerKey"`
}


var (
	 SvJson *ServerJson
	 PathModelPtr   *PathInfo //最先有路径对象
)



func InitConf()  {
	SvJson = new(ServerJson)
	SvJson.DCServerID = 15 // 默认值
	RedisCfg = new(RootConf)
	setAppPath()
}

//设置app程序路径
func setAppPath() {
	//创建对象在前
	PathModelPtr = newPathModel()
	if PathModelPtr == nil {
		panic("创建 PathModelPtr 失敗")
	}
	pwd, _ := os.Getwd()
	//测试写死路径
	fmt.Println("服务器当前路径: ",pwd)
	PathModelPtr.SetRootPath(pwd)
	PathModelPtr.InitPathModel()
}

//设置app程序路径
func SetAppPath(pwd string) {
	//创建对象在前
	PathModelPtr = newPathModel()
	if PathModelPtr == nil {
		panic("创建 PathModelPtr 失敗")
	}
	//测试写死路径
	fmt.Println("服务器当前路径: ",pwd)
	PathModelPtr.SetRootPath(pwd)
	PathModelPtr.InitPathModel()
}

func StartConf()  {
	//获取配置
	erro := readJson(PathModelPtr.ServerConfPath)
	if erro != nil {
		panic(erro)
	}
	erro = readRedisConf(PathModelPtr.RedisConfPath)
	if erro != nil {
		panic(erro)
	}
}

func readJson(iniPath string) error  {
	filePtr, err := os.Open(iniPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(SvJson)
	if err != nil {
		return err
	}
	fmt.Printf("readJson %v \n",SvJson)
	return nil
}

func OnAppClose()  {
	SvJson = nil
	RedisCfg = nil
	PathModelPtr = nil
}