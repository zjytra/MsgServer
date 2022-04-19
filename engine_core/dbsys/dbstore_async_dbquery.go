/*
创建时间: 2020/08/2020/8/26
作者: Administrator
功能介绍:
数据库异步查询事件
*/
package dbsys

import (
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

// 异步写
//dbcb 回调方法
func (this *MySqlDBStore) AsyncExecute(dbParam *DBParam,
	//logicQe dispatch.WaitQueue,
	dbcb DBExecuteCallback,
	query string,
	args ...interface{}) {
	if strutil.StringIsNil(query) {
		return
	}
	event := NewDBWriteEvent()
	if event == nil {
		return
	}
	event.store = this
	event.query = query
	event.args = args
	event.logicQe = dispatch.MainQueue  //目前默认返回主逻辑队列
	event.cb = dbcb
	event.param = dbParam
	this.writeEvent.AddEvent(event)
}

// 异步查询返回原始结果 取 DBQueryResult.DBRows
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
//@param args 拼装字符串参数
func (this *MySqlDBStore) AsyncQueryRows(dbP *DBParam,
	//logicQe dispatch.WaitQueue,
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, nil, nil)
	dbParam.Data = dbP.Data
	this.asyncQuery(DBQueryRowsCB_Event, dbParam, dbcb, query, args...)
}

// 异步查询返回单行结果 取 DBQueryResult.MapRow
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
//@param args 拼装字符串参数
func (this *MySqlDBStore) AsyncQueryRowMap(dbP *DBParam,
	//logicQe dispatch.WaitQueue,
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, nil, nil)
	this.asyncQuery(DBQueryRowToMapCb_Event, dbParam,  dbcb, query, args...)
}



// 异步查询多个结果集 取 DBQueryResult.MapRowS
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
//@param args 拼装字符串参数
func (this *MySqlDBStore) AsyncQueryRowsMap(dbP *DBParam,
	//logicQe dispatch.WaitQueue,
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, nil, nil)
	this.asyncQuery(DBQueryRowsToMapArrCb_Event, dbParam, dbcb, query, args...)
}

// 异步查询返回多行结果集 取 DBQueryResult.MoreMapRowS
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
//@param args 拼装字符串参数
func (this *MySqlDBStore) AsyncQueryMoreRowsMap(dbP *DBParam,
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, nil, nil)
	this.asyncQuery(DBQueryMoreReusltMapCb_Event, dbParam, dbcb, query, args...)
}

// 异步查询返回多行结果集 取 DBQueryResult.MoreObjs
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
//@param args 拼装字符串参数
func (this *MySqlDBStore) AsyncQueryMoreStructArr(dbP *DBParam,
	reflArr []interface{},
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, nil, reflArr)
	this.asyncQuery(DBQueryMoreReusltStructArrCb_Event, dbParam,  dbcb, query, args...)
}

// 异步查询返回单个结构体 DBQueryResult.QueryObj
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncQueryStruct(dbP *DBParam,
	refl interface{},
	//logicQe dispatch.WaitQueue,
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, refl, nil)
	this.asyncQuery(DBQueryRowToStructCb_Event, dbParam, dbcb, query, args...)
}

// 异步查询返回单个结构体 DBQueryResult.QueryObj
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncQueryDBObj(dbP *DBParam,
	refl interface{},
	dbcb OnDBQueryCB) {
	//dbParam := NewDBQueryParam(dbP, refl, nil)
	////this.asyncQuery(DBQueryDBObj, dbParam, dbcb, query, args...)
}

// 异步查询返回结构体数组 DBQueryResult.Objs
//@param dbP 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由logicQe 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncQueryStructArr(dbP *DBParam,
	refl interface{},
	dbcb OnDBQueryCB,
	query string,
	args ...interface{}) {
	dbParam := NewDBQueryParam(dbP, refl, nil)
	this.asyncQuery(DBQueryRowsToStructSliceCb_Event, dbParam, dbcb, query, args...)
}

////投递查询数据
func (this *MySqlDBStore) asyncQuery(dtype int, dbParam *DBQueryParam,  dbcb OnDBQueryCB, query string, args ...interface{}) {
	if strutil.StringIsNil(query) {
		xlog.Error("AsyncQuery  query == nil")
		return
	}
	//由于逻辑线程使用了对象池，这里投递的时候就不用池子了避免提前被回收
	data := queryEventPool.Get()
	queryEvent, ok := data.(*DBQueryEvent)
	if !ok || queryEvent == nil {
		xlog.Error("创建DBQueryEvent 失败")
		return
	}
	queryEvent.queryType = dtype
	queryEvent.param = dbParam
	queryEvent.store = this
	queryEvent.query = query
	queryEvent.args = args
	queryEvent.logicQe = dispatch.MainQueue
	queryEvent.cb = dbcb
	this.quereyEvent.AddEvent(queryEvent)
}

// 异步查询
//@param oneRow 自定义查询 具体的在 CustomDBOperate 类型中实现
func (this *MySqlDBStore) AsyncCustomOneRowQuery(oneRow CustomDBOperate) {
	if oneRow == nil {
		return
	}
	poolObj := queryCusTomPool.Get()
	data, ok := poolObj.(*CustomDBEvent)
	if !ok || data == nil {
		xlog.Error("创建CustomDBEvent 失败")
		return
	}
	data.op = oneRow
	data.logicQe = dispatch.MainQueue
	this.quereyEvent.AddEvent(data)
}




//////////////////////////////////////////////////////////////////////
//按规定创建的数据库对象

//加载单个对象自定义回调
func (this *MySqlDBStore) AsyncLoadObJerFromDBAndCb(param *DBObjQueryParam,cb OnLoadFormCb) {
	if param == nil || param.DbObj == nil {
		return
	}
	poolObj := dbObjQueryPool.Get()
	data, ok := poolObj.(*DBObjQueryEvent)
	if !ok || data == nil {
		xlog.Error("创建CustomDBEvent 失败")
		return
	}

	if cb == nil {
		xlog.Error("AsyncLoadObJerFromDBAndCb 查询无回调")
		return
	}
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue //主线程回调
	data.loadCb = cb
	this.quereyEvent.AddEvent(data)
}

//加载单个对象 对象回调
func (this *MySqlDBStore) AsyncLoadFromDB(param *DBObjQueryParam) {
	if param == nil || param.DbObj == nil {
		return
	}
	poolObj := dbObjQueryPool.Get()
	data, ok := poolObj.(*DBObjQueryEvent)
	if !ok || data == nil {
		xlog.Error("创建CustomDBEvent 失败")
		return
	}
	if param == nil {
		return
	}
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	this.quereyEvent.AddEvent(data)
}



//加载关联的所有子表数据
func (this *MySqlDBStore) AsyncLoadSubTables(param *DBObjQueryParam,cb OnLoadSubTableFormCb) {
	if param == nil || param.DbObj == nil {
		return
	}
	if param.DbObj.GetUID() == 0 {
		xlog.Warning("AsyncLoadSubTables obj uid is 0")
		return
	}
	if !param.DbObj.isRegisterDBInterface() {
		this.SetDBObjDefault(param.DbObj)
	}
	pTable := this.GetDBTable(param.DbObj.GetTabName())
	if pTable == nil {
		xlog.Debug("AsyncLoadSubTables gettable %s is nil",param.DbObj.GetTabName())
		return
	}
	if !pTable.HasSubTables() {
		xlog.Debug("AsyncLoadSubTables table %s not subTable",param.DbObj.GetTabName())
		return
	}
	poolObj := dbObjQueryMorePool.Get()
	data, ok := poolObj.(*DBObjQueryMoreSubEvent)
	if !ok || data == nil {
		xlog.Error("创建CustomDBEvent 失败")
		return
	}
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	data.loadSubCb = cb
	this.quereyEvent.AddEvent(data)
}


//查询单个表多个数据
func (this *MySqlDBStore) AsyncLoadMoreObjs(param *DBObjQueryParam,cb OnLoadMoreFormCb) {
	this.asyncLoadMoreObjs(param, cb,false)
}

//查询单个表多个数据
func (this *MySqlDBStore) AsyncLoadAllObjs(param *DBObjQueryParam,cb OnLoadMoreFormCb) {
	this.asyncLoadMoreObjs(param, cb,true)
}

func (this *MySqlDBStore) asyncLoadMoreObjs(param *DBObjQueryParam, cb OnLoadMoreFormCb,isAll bool) {
	if param == nil || param.DbObj == nil {
		return
	}
	//这里注册
	if !param.DbObj.isRegisterDBInterface() {
		this.SetDBObjDefault(param.DbObj)
	}
	pTable := this.GetDBTable(param.DbObj.GetTabName())
	if pTable == nil {
		xlog.Debug("AsyncLoadMoreObjs gettable %s is nil", param.DbObj.GetTabName())
		return
	}
	data := new(DBObjQueryMoreEvent)
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	data.loadCb = cb
	data.isAllLoad = isAll
	this.quereyEvent.AddEvent(data)
}

//异步写数据
func (this *MySqlDBStore) AsyncInsertObj(param *DBObjWriteParam,cb OnWriteEventCb) {
	if param == nil || !param.CheckParam() {
		return
	}
	param.writeTpe = DBInsertDBObj
	data := new(DBObjWriteEvent)
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	data.writeCb = cb
	this.quereyEvent.AddEvent(data)
}

//异步更新数据
func (this *MySqlDBStore) AsyncUpdateObj(param *DBObjWriteParam,cb OnWriteEventCb) {
	if param == nil || !param.CheckParam() {
		return
	}
	//这里注册
	param.writeTpe = DBUpdateDBObj
	data := new(DBObjWriteEvent)
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	data.writeCb = cb
	this.quereyEvent.AddEvent(data)
}

//异步删除数据
func (this *MySqlDBStore) AsyncDelObj(param *DBObjWriteParam,cb OnWriteEventCb) {
	if param == nil || !param.CheckParam(){
		return
	}
	//这里注册
	param.writeTpe = DBDeleteDBObj
	data := new(DBObjWriteEvent)
	data.Param = param
	data.store = this
	data.logicQe = dispatch.MainQueue
	data.writeCb = cb
	this.quereyEvent.AddEvent(data)
}