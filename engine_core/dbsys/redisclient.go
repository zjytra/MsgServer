/*
创建时间: 2020/12/19 23:36
作者: zjy
功能介绍:

*/

package dbsys

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zjytra/MsgServer/conf"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

var (
	RedisCli *redis.Client
	CtxBg = context.Background()
)


func AsyncOpenRedis(clientName string) {
	reConf := conf.GetRedisConfByName(clientName)
	if reConf == nil {
		xlog.Error("没有找到 %s 的redis 配置",clientName)
		return
	}
	//异步连接避免阻塞协程
	go func() {
		addr := reConf.IP + ":" + reConf.Dbport
		RedisCli = redis.NewClient(&redis.Options{
			Addr:  addr   ,
			Password: reConf.Dbpwd, // no password set
			DB:       0,  // use default DB
			DialTimeout : 1,
			ReadTimeout : 1,
			WriteTimeout : 1,
		})
		_, erro := RedisCli.Ping(CtxBg).Result()
		if erro != nil {
			fmt.Print(erro)
			return
		}
		fmt.Println("连接redis clientName",clientName," 地址 ",addr,"成功")
	}()
}


