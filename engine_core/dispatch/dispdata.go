/*
创建时间: 2020/12/20 0:07
作者: zjy
功能介绍:
定义变量
*/

package dispatch

import (
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

var (
	MainQueue *SyncQueue
)




func InitMainQueue()  {
	MainQueue = NewSyncQueue(5000) // 主队列只有一个线程调用
	if MainQueue == nil {
		panic("InitMainQueue is nil")
	}
}


func AddEventToQueue(event Evener)  {
	err := MainQueue.AddEvent(event)
	if err != nil {
		xlog.Warning("AddEventToQueue %v" ,err)
	}

}



func CheckTime(modeStr string,startT int64,checkTime int64)  {
	since := mathutil.MaxInt64(0, timeutil.GetCurrentTimeMs()-startT)
	if since >= checkTime {
		xlog.Warning("%s 执行时间%v ms" ,modeStr, since)
	}
}
