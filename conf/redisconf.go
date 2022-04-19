/*
创建时间: 2021/7/3 22:23
作者: zjy
功能介绍:

*/

package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var(
	RedisCfg      *RootConf
)

type RootConf struct {
	RedisConfigs []*RedisConf `json:"RedisConfigs"`
}

//redis 配置
type RedisConf struct {
	ClientName string `json:"clientName"`  //客户端名称
	IP  string      `json:"IP"`     	//ip地址
	Dbport  string  `json:"dbport"`  	//端口号
	Dbpwd  string	`json:"dbpwd"`		//密码
	Ismaster  bool	`json:"ismaster"`	//客户端名称
}

func readRedisConf(iniPath string) error  {
	filePtr, err := os.Open(iniPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	buf,erro := ioutil.ReadAll(filePtr)
	if erro != nil {
		return  erro
	}
	err = json.Unmarshal(buf, RedisCfg)
	if err != nil {
		return err
	}
	return nil
}

func GetRedisConfByName(clientName string) *RedisConf {
	if RedisCfg == nil || RedisCfg.RedisConfigs == nil || len(RedisCfg.RedisConfigs) == 0 {
		return nil
	}
	for _, conf := range RedisCfg.RedisConfigs {
		if conf.ClientName == clientName {
			return  conf
		}
	}
	return  nil
}