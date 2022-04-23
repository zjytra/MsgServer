/*
创建时间: 2021/9/19 17:24
作者: zjy
功能介绍:

*/

package dbsys

import "sync"

var (
	dbParamPool       sync.Pool //投递参数对象池
	writeEventPool    sync.Pool //写事件对象池
	writeDelayEventPool    sync.Pool //延迟写事件
	writeCbPool       sync.Pool //查询结果回调对象池
	queryEventPool    sync.Pool //查询事件对象池
	queryCbPool       sync.Pool
	queryCusTomPool   sync.Pool
	queryCusTomCbPool sync.Pool
	queryParamPool    sync.Pool
	dbObjQueryPool  sync.Pool   //数据库对象查询对象池
	dbObjQueryCbPool  sync.Pool  //数据库对象查询回调对象池

	dbObjQueryMorePool  sync.Pool   //数据库对象查询对象池
	dbObjQueryMoreCbPool  sync.Pool  //数据库对象查询回调对象池
)

////数据库执行返回
type DBExecuteCallback func(dbResult  *DBWriteCb)

func init()  {
	dbParamPool.New = func() interface{} {
		return new(DBParam)
	}
	writeEventPool.New = func() interface{} {
		return new(DBWriteEvent)
	}
	writeCbPool.New = func() interface{} {
		return new(DBWriteCb)
	}
	queryEventPool.New = func() interface{} {
		qevnt :=  new(DBQueryEvent)
		qevnt.result = new(DBQueryResult)
		return qevnt
	}
	queryCbPool.New = func() interface{} {
		return new(DBQueryCb)
	}
	queryCusTomPool.New = func() interface{} {
		return new(CustomDBEvent)
	}
	queryCusTomCbPool.New = func() interface{} {
		return new(CustomDBEventCb)
	}
	queryParamPool.New = func() interface{} {
		return new(DBQueryParam)
	}
	writeDelayEventPool.New = func() interface{} {
		return new(DBWriteDelayEvent)
	}
	dbObjQueryPool.New = func() interface{} {
		return new(DBObjQueryEvent)
	}
	dbObjQueryCbPool.New = func() interface{} {
		return new(DBObjQueryEventCb)
	}

	dbObjQueryMorePool.New = func() interface{} {
		return new(DBObjQueryMoreSubEvent)
	}
	dbObjQueryMoreCbPool.New = func() interface{} {
		return new(DBObjQueryMoreSubEventCb)
	}
}
