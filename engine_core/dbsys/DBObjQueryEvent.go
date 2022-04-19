/*
创建时间: 2021/8/15 17:14
作者: zjy
功能介绍:
数据库对象查询,只需要传数据库对象
*/

package dbsys

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

type OnLoadFormCb func(result DBErrorCode,param *DBObjQueryParam)
//数据库对象查询
type DBObjQueryEvent struct {
	store *MySqlDBStore
	logicQe dispatch.WaitQueue //逻辑队列
	Param *DBObjQueryParam
	loadCb OnLoadFormCb   //查询回调,当设置这个回调就不会回调DBObJer的回调
}


//队列调度
func (this *DBObjQueryEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	//这里注册
	this.store.SetDBObjDefault(this.Param.DbObj)
	sql := this.Param.DbObj.SelectSql()

	rows,err :=	this.store.Query(sql)
	defer rows.Close()
	dispatch.CheckTime("数据库 :"+ this.EvenName() + sql,startT,200)
	resCode := NODATA
	//错误不拦截投递到主线程去处理,万一要重试呢？
	if err == nil {
		for rows.Next() {
			val := this.Param.DbObj.GetDBValesAddr()
			sErro := rows.Scan(val...)
			if sErro != nil {
				xlog.Debug("DBObjQueryEvent 查询表 %s 赋值 错误 %v",this.Param.DbObj.GetTabName(), sErro)
			}
			resCode = DBSUCCESS
		}
		if  resCode == DBSUCCESS {
			this.Param.DbObj.setCreateToDB()
			this.Param.DbObj.dBValToVal()
		}
	}else {
		resCode = DBSQLERRO
	}
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil {
		dbObjQueryPool.Put(this)
	 	xlog.Error("DBObjQueryEvent 队列和回调为nil")
		return
	}

	poolObj := dbObjQueryCbPool.Get()
	data,ok := poolObj.(*DBObjQueryEventCb)
	if !ok || data == nil {
		dbObjQueryPool.Put(this)
		xlog.Error("创建 DBQueryCb 失败")
		return
	}
	data.dbEvent = this
	data.code = resCode
	//向监听消息者投递
	this.logicQe.AddEvent(data)
}

func (this *DBObjQueryEvent)EvenName() string {
	return "DBObjQueryEvent"
}


type DBObjQueryEventCb struct {
	dbEvent *DBObjQueryEvent
	code DBErrorCode
}

//逻辑线程执行
func (this *DBObjQueryEventCb) Execute(){
	//由回调使用回调
	if this.dbEvent.loadCb != nil {
		this.dbEvent.loadCb(this.code,this.dbEvent.Param)
	}else {
		this.dbEvent.Param.DbObj.OnLoadForm(this.code,this.dbEvent.Param)
	}
	this.dbEvent.Param = nil
	this.dbEvent.store = nil
	dbObjQueryPool.Put(this.dbEvent)
	dbObjQueryCbPool.Put(this)
}


func (this *DBObjQueryEventCb)EvenName() string {
	return "CustomDBEventCb"
}