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
type DBQueryParam struct {
	DBParam
	ReflectObj    interface{} //传入的需要反射的参数,
	ReflectObjArr []interface{} //当为多个结果集的时候就要传递多个对象映射不同的表
}


func NewDBQueryParam(param *DBParam,refl  interface{},reflArr  []interface{}) *DBQueryParam {
	data := queryParamPool.Get()
	query,ok := data.(*DBQueryParam)
	if !ok || query == nil {
		return nil
	}
	query.UID = param.UID
	query.CltConn = param.CltConn
	query.SvConn = param.SvConn
	query.Data = param.Data
	query.ReflectObj = refl
	query.ReflectObjArr = reflArr
	return query
}

//查询结果
type DBQueryResult struct {
	DBRestCode int32           //返回的结果码
	//根据不同的数据使用不同的对象
	QueryObj interface{}
	//返回结构体数组对象
	Objs     []interface{}
	DBRows   *sql.Rows     //返回的原始结果集
	MapRow map[string]string //单行数据 gamemap[数据库字段]值
	MapRowS []map[string]string //多行数据 gamemap[数据库字段]值
	MoreMapRowS map[int][]map[string]string // 多个结果集map列表map[结果集序号从1开始]多行数据
	MoreObjs map[int][]interface{} // 多个结果集map列表map[第几个结果集]多行数据
	DBErr error
}

//数据库线程直接查询返回结构体不过结构体的要与数据库表字段匹配
type OnDBQueryCB func(param  *DBQueryParam,result *DBQueryResult)

type DBQueryEvent struct {
	queryType int
	store *MySqlDBStore
	query string
	args  []interface{}
	logicQe dispatch.WaitQueue //逻辑队列
	cb OnDBQueryCB
	param  *DBQueryParam
	result  *DBQueryResult
}


//队列调度
func (this *DBQueryEvent) Execute(){
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	defer dispatch.CheckTime("数据库 :"+ this.EvenName(),startT,200)
	//向监听消息者投递
	var rows *sql.Rows
	var err error

	if this.args != nil && len(this.args)  > 0 {
		rows,err =	this.store.Query(this.query,this.args ...)
	}else {
		rows,err =	this.store.Query(this.query)
	}
	//错误不拦截投递到主线程去处理,万一要重试呢？
	if err == nil {
		//queryParamPool.Put(this.param)
		//queryEventPool.Put(this)
		//return err
		defer rows.Close()
		//处理的方式
		switch this.queryType {
		case DBQueryRowsCB_Event:
			this.result.DBRows = rows
		case DBQueryRowToStructCb_Event:
			this.result.QueryObj = RowToStruct(rows, this.param.ReflectObj)
		case DBQueryRowsToStructSliceCb_Event:
			this.result.Objs = RowsToStructSlice(rows, this.param.ReflectObj)
		case DBQueryRowToMapCb_Event:
			this.result.MapRow = RowToMap(rows)
		case DBQueryRowsToMapArrCb_Event:
			this.result.MapRowS = RowsToStrMapArr(rows)
		case DBQueryMoreReusltMapCb_Event:
			this.result.MoreMapRowS = MoreResultRowsToStrMapArr(rows)
		case DBQueryMoreReusltStructArrCb_Event:
			this.result.MoreObjs = MoreResultRowsToStructArr(rows,this.param.ReflectObjArr)
		default:
			xlog.Error("DBQueryEvent查询类型%v未处理", this.queryType)
			return
		}
	}else {
		xlog.Error("数据库Query %s, 出错 %v",this.query,err)
	}
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil || this.cb == nil {
		if this.param != nil {
			queryParamPool.Put(this.param)
		}
		queryEventPool.Put(this)
		xlog.Debug("DBQueryEvent 队列和回调为nil")
		return
	}

	poolObj := queryCbPool.Get()
	data,ok := poolObj.(*DBQueryCb)
	if !ok || data == nil {
		if this.param != nil {
			queryParamPool.Put(this.param)
		}
		queryEventPool.Put(this)
		xlog.Debug("创建 DBQueryCb 失败")
		return
	}
	this.result.DBErr = err
	data.DBEvent = this
	this.logicQe.AddEvent(data)
}

func (this *DBQueryEvent)EvenName() string {
	return "DBQueryEvent"
}


//查询结果回调
type DBQueryCb struct {
	DBEvent *DBQueryEvent
}

//逻辑线程执行
func (this *DBQueryCb) Execute(){
	this.DBEvent.cb(this.DBEvent.param,this.DBEvent.result)
	if this.DBEvent.param != nil {
		queryParamPool.Put(this.DBEvent.param)
	}
	queryEventPool.Put(this.DBEvent)
	queryCbPool.Put(this)
}

func (this *DBQueryCb)EvenName() string {
	return "DBQueryCb"
}



