/*
创建时间: 2021/8/2 22:33
作者: zjy
功能介绍:

*/

package dbsys

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

type DBWriteDelayEvent struct {
	store *MySqlDBStore
	query string
}

//数据库线程使用了还不能回收，需要逻辑线程搞完了才能回收
func NewDBWriteDelayEvent() *DBWriteDelayEvent  {
	poolObj := writeDelayEventPool.Get()
	data,ok := poolObj.(*DBWriteDelayEvent)
	if !ok || data == nil {
		xlog.Error("NewDBWriteEvent 失败")
		return nil
	}
	return data
	//return new(DBWriteEvent)
}

//队列调度
func (this *DBWriteDelayEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	_,err := this.store.Execute(this.query)
	dispatch.CheckTime("数据库 :" + this.EvenName(),startT,200)
	if err != nil {
		xlog.Error("数据库Execute%s, 出错 %v",this.query,err)
	}
	writeDelayEventPool.Put(this)
}

func (this *DBWriteDelayEvent)EvenName() string {
	return "DBWriteDelayEvent"
}



