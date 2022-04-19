/*
创建时间: 2021/2/25 22:52
作者: zjy
功能介绍:
数据库查询事件执行，执行完再抛回主逻辑线程
*/

package dbsys

import (
	"database/sql"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)


//投递参数
type DBParam struct {
	UID     int64  //账号id 或者 userID
	CltConn uint32 //网络连接id
	SvConn  uint32 //服务器连接id
	Data    interface{}
}

func NewDBParam(mark int64,cltConn uint32,svCon uint32) *DBParam {
	data := dbParamPool.Get()
	param,ok := data.(*DBParam)
	if !ok || param == nil {
		return nil
	}
	param.UID = mark
	param.CltConn = cltConn
	param.SvConn = svCon
	return param
}

func NewParamByCln(cltConn uint32) *DBParam {
	param := new(DBParam)
	param.CltConn = cltConn
	return param
}


func NewParamByMark(mark int64) *DBParam {
	param := new(DBParam)
	param.UID = mark
	return param
}
func NewDBParamNotInPool(mark int64,cltConn uint32,svCon uint32) *DBParam {
	param := new(DBParam)
	param.UID = mark
	param.CltConn = cltConn
	param.SvConn = svCon
	return param
}

type DBWriteEvent struct {
	store *MySqlDBStore
	query string
	args  []interface{}
	logicQe dispatch.WaitQueue //逻辑队列
	cb DBExecuteCallback
	param  *DBParam
}

//数据库线程使用了还不能回收，需要逻辑线程搞完了才能回收
func NewDBWriteEvent() *DBWriteEvent  {
	poolObj := writeEventPool.Get()
	data,ok := poolObj.(*DBWriteEvent)
	if !ok || data == nil {
		xlog.Error("NewDBWriteEvent 失败")
		return nil
	}
	data.args = nil
	data.logicQe = nil //目前默认返回主逻辑队列
	data.cb = nil
	data.param = nil
	return data
	//return new(DBWriteEvent)
}

//队列调度
func (this *DBWriteEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	var rest sql.Result
	var err error
	if this.args != nil && len(this.args) > 0 {
		rest,err =	this.store.Execute(this.query,this.args ...)
	}else {
		rest,err =	this.store.Execute(this.query)
	}
	dispatch.CheckTime("数据库 :" + this.EvenName(),startT,200)
	//rest,err :=	this.store.Execute(this.query)
	//err应该不拦截投递到主线程处理 ,万一要重试呢？
	if err != nil {
		//if this.param != nil {
		//	dbParamPool.Put(this.param)
		//}
		//writeEventPool.Put(this)
	}
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil || this.cb == nil {
		if this.param != nil {
			dbParamPool.Put(this.param)
		}
		writeEventPool.Put(this)
		return
	}
	data := NewDBWriteCb()
	if data == nil {
		if this.param != nil {
			dbParamPool.Put(this.param)
		}
		writeEventPool.Put(this)
		xlog.Debug("创建 DBWriteCb 失败")
		return
	}
	data.Result = rest
	data.Event = this
	data.Err = err
	this.logicQe.AddEvent(data)
}

func (this *DBWriteEvent)EvenName() string {
	return "DBWriteEvent"
}


//逻辑线程回调
type DBWriteCb struct {
	Result sql.Result
	Event *DBWriteEvent //传过来回收用
	Err error //错误也传回去
}

//逻辑线程搞完了才能回收
func NewDBWriteCb() *DBWriteCb  {
	poolObj := writeCbPool.Get()
	data,ok := poolObj.(*DBWriteCb)
	if !ok || data == nil {
		return nil
	}
	return data
	//return new(DBWriteEvent)
}

//逻辑线程执行
func (this *DBWriteCb) Execute(){
	//err :=	this.cb(this.param,this.result)
	this.Event.cb(this)
	if this.Event.param != nil {
		dbParamPool.Put(this.Event.param)
	}
	writeEventPool.Put(this.Event)
	writeCbPool.Put(this)
}

func (this *DBWriteCb)EvenName() string {
	return "DBWriteCb"
}

