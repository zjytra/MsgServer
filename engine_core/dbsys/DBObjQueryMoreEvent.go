/*
创建时间: 2021/8/15 17:14
作者: zjy
功能介绍:
数据库对象查询,查询多条数据
*/

package dbsys

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

type OnLoadMoreFormCb func(result DBErrorCode,param *DBObjQueryParam,MoreObj []DBObJer)
//数据库对象查询
type DBObjQueryMoreEvent struct {
	store *MySqlDBStore
	logicQe dispatch.WaitQueue //逻辑队列
	Param *DBObjQueryParam
	loadCb OnLoadMoreFormCb   //查询回调,当设置这个回调就不会回调DBObJer的回调
	MoreObj []DBObJer  //返回的结果集
	isAllLoad  bool        //加载该表所有数据
}


//队列调度
func (this *DBObjQueryMoreEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	//这里注册
	this.store.SetDBObjDefault(this.Param.DbObj)
	pTable := this.store.GetDBTable(this.Param.DbObj.GetTabName())
	if pTable == nil {
		return
	}
	var sql string
	if this.isAllLoad { //加载所有
		sql = this.Param.DbObj.LoadAll()
	}else {
		sql = this.Param.DbObj.SelectSql()
	}
	//向监听消息者投递
	rows,err :=	this.store.Query(sql)
	defer rows.Close()
	dispatch.CheckTime("数据库 :"+ this.EvenName() + sql,startT,200)
	resCode := NODATA
	//错误不拦截投递到主线程去处理,万一要重试呢？
	if err == nil {
		this.MoreObj = this.store.RowsToDBObjArr(rows, pTable, &resCode)
	}else {
		resCode = DBSQLERRO
	}
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil {
	 	xlog.Error("DBObjQueryMoreEvent 队列和回调为nil")
		return
	}

	data := new(DBObjQueryMoreEventCb)
	data.dbEvent = this
	data.code = resCode
	this.logicQe.AddEvent(data)
}

func (this *DBObjQueryMoreEvent)EvenName() string {
	return "DBObjQueryMoreEvent"
}


type DBObjQueryMoreEventCb struct {
	dbEvent *DBObjQueryMoreEvent
	code DBErrorCode
}

//逻辑线程执行
func (this *DBObjQueryMoreEventCb) Execute(){
	//线程队列回调
	if this.dbEvent.loadCb != nil {
		this.dbEvent.loadCb(this.code,this.dbEvent.Param,this.dbEvent.MoreObj)
	}
	this.dbEvent.Param = nil
	this.dbEvent.store = nil
	this.dbEvent.MoreObj = nil //清空数据
}


func (this *DBObjQueryMoreEventCb)EvenName() string {
	return "DBObjQueryMoreEventCb"
}