/*
创建时间: 2020/08/2020/8/29
作者: Administrator
功能介绍:

*/
package apploginsv

import (
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

// 每秒定时器
func (this *LoginServer) PerOneSTimer() {
	xlog.Debug("PerOneSTimer")
}



// 分钟定时器
func (this *LoginServer) PerOneMinuteTimer()  {

}

// 小时定时器
func (this *LoginServer) PerOneHourTimer() {

}