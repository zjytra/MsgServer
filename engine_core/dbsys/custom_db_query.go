/*
创建时间: 2021/2/26 19:06
作者: zjy
功能介绍:
自定以查询封装
*/

package dbsys

import (
	"github.com/zjytra/MsgServer/devlop/xutil/mathutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
)

//逻辑线程处理自定义查询
type CustomDBOperate interface {
	//执行数据库查询
	ExecuteQueryFun() error
	//查询回调
	OnQueryCB() error
}

type CustomDBEvent struct {
	op CustomDBOperate
	logicQe dispatch.WaitQueue //逻辑队列
}

//数据库线程执行查询
func (this *CustomDBEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	err := this.op.ExecuteQueryFun()
	since := mathutil.MaxInt64(0, timeutil.GetCurrentTimeMs()-startT)
	if since >= 200 {
		reFlecttype := reflect.TypeOf(this.op)
		xlog.Warning("自定义接口%v,执行时间%v ms", reFlecttype.String(), since)
	}
	if err != nil {
		queryCusTomPool.Put(this)
		return
	}
	poolObj := queryCbPool.Get()
	data,ok := poolObj.(*CustomDBEventCb)
	if !ok || data == nil {
		xlog.Debug("创建CustomDBEventCb 失败")
		return
	}
	if this.logicQe == nil {
	 	xlog.Debug("CustomDBEvent 逻辑队列为nil")
		return
	}
	data.dbEvent = this
	this.logicQe.AddEvent(data)
}

func (this *CustomDBEvent)EvenName() string {
	return "CustomDBEvent"
}

type CustomDBEventCb struct {
	dbEvent *CustomDBEvent
}

//逻辑线程执行
func (this *CustomDBEventCb) Execute(){
	this.dbEvent.op.OnQueryCB()
	queryCusTomPool.Put(this.dbEvent)
	queryCbPool.Put(this)
}


func (this *CustomDBEventCb)EvenName() string {
	return "CustomDBEventCb"
}
