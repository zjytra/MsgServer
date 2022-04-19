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
	"strings"
)
//查询子表回调
type OnLoadSubTableFormCb func(result DBErrorCode,Param *DBObjQueryParam,objArr map[int][]DBObJer)
//数据库对象查询
type DBObjQueryMoreSubEvent struct {
	store *MySqlDBStore
	logicQe dispatch.WaitQueue //逻辑队列
	Param *DBObjQueryParam
	objArr map[int][]DBObJer
	loadSubCb OnLoadSubTableFormCb
}


//队列调度
func (this *DBObjQueryMoreSubEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	//这里注册
	this.store.SetDBObjDefault(this.Param.DbObj)
	uid := this.Param.DbObj.GetUID()
	pTable := this.store.GetDBTable(this.Param.DbObj.GetTabName())
	if pTable == nil {
		dbObjQueryMorePool.Put(this)
		return
	}
	//获得子表
	subTables := pTable.GetSubTables()
	if subTables == nil || len(subTables) == 0 {
		dbObjQueryMorePool.Put(this)
		return
	}

	//结果集设置为nil
	this.objArr = make(map[int][]DBObJer)
	var sb strings.Builder
	for _, table := range subTables {
		sb.WriteString(table.SelectSql(uid))
	}
	sql := sb.String()
	dispatch.CheckTime("数据库 :"+ this.EvenName() + sql,startT,200)
	//向监听消息者投递
	rows,err :=	this.store.Query(sql)
	defer rows.Close()
	resCode := NODATA
	if err == nil {
		for i, table := range subTables {
			objs := this.store.RowsToDBObjArr(rows, table, &resCode)
			this.objArr[i+1] = objs
			if !rows.NextResultSet() {
				break
			}
		}
	} else {
		resCode = DBSQLERRO
	}
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil {
		dbObjQueryMorePool.Put(this)
		xlog.Error("DBObjQueryMoreSubEvent 队列和回调为nil")
		return
	}

	poolObj := dbObjQueryMoreCbPool.Get()
	data, ok := poolObj.(*DBObjQueryMoreSubEventCb)
	if !ok || data == nil {
		dbObjQueryMorePool.Put(this)
		xlog.Error("创建 DBObjQueryMoreSubEventCb 失败")
		return
	}
	data.dbEvent = this
	data.code = resCode
	this.logicQe.AddEvent(data)
}



func (this *DBObjQueryMoreSubEvent)EvenName() string {
	return "DBQueryEvent"
}


type DBObjQueryMoreSubEventCb struct {
	dbEvent *DBObjQueryMoreSubEvent
	code DBErrorCode
}

//逻辑线程执行
func (this *DBObjQueryMoreSubEventCb) Execute(){
	this.dbEvent.loadSubCb(this.code,this.dbEvent.Param,this.dbEvent.objArr)
	this.dbEvent.Param = nil
	this.dbEvent.objArr = nil //清空数据
	this.dbEvent.store = nil
	dbObjQueryMorePool.Put(this.dbEvent)
	dbObjQueryMoreCbPool.Put(this)
}


func (this *DBObjQueryMoreSubEventCb)EvenName() string {
	return "CustomDBEventCb"
}