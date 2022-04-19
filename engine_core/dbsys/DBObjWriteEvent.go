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

type OnWriteEventCb func(result DBErrorCode, param *DBObjWriteParam)

//数据库对象查询
type DBObjWriteEvent struct {
	store   *MySqlDBStore
	logicQe dispatch.WaitQueue //逻辑队列
	Param   *DBObjWriteParam
	writeCb OnWriteEventCb //删除数据回调
	//避免sql超长
	writeIndex     int //已经写到哪个下标未包含
	hasNum         int //需要写多少次

	valToDBValIndex int //已经写到哪个下标未包含
	writeSuccessNum int //需要写多少次
	allSuccessNum   int //需要写多少次
}

//队列调度
func (this *DBObjWriteEvent) Execute() {
	//先计算出需要写多少次
	this.hasNum = len(this.Param.DbObjs)
	this.writeIndex = 0
	this.valToDBValIndex = 0
	resCode := this.doExecuteSql()
	//如果这些参数是nil的就直接返回了
	if this.logicQe == nil || this.writeCb == nil {
		//xlog.Debug("DBObjWriteEvent 队列或回调为nil")
		return
	}
	data := new(DBObjWriteEventCb)
	data.dbEvent = this
	data.code = resCode
	this.logicQe.AddEvent(data)
}

func (this *DBObjWriteEvent) doExecuteSql() DBErrorCode {
	resCode := NotDo
	//没有全部写入成功再继续写
	for this.allSuccessNum != this.hasNum   {
		startT := timeutil.GetCurrentTimeMs() //计算当前时间
		sql := ""
		switch this.Param.writeTpe {
		case DBInsertDBObj:
			sql = this.getInsertObjsSql()
		case DBUpdateDBObj:
			sql = this.getUpdateObjsSql()
		case DBDeleteDBObj:
			sql = this.getDelObjsSql()
		}
		if len(sql) == 0 {
			xlog.Debug("DBObjWriteEvent Execute sql is nil")
			return resCode
		}
		rest, err := this.store.Execute(sql)
		dispatch.CheckTime("数据库 :"+this.EvenName()+sql, startT, 200)
		//错误不拦截投递到主线程去处理,万一要重试呢？
		if err == nil {
			id, resErro := rest.RowsAffected()
			if resErro == nil && id > 0 {
				resCode = DBSUCCESS
			} else {
				xlog.Debug("DBObjWriteEvent Execute %s \n resErro %v", sql, resErro)
				resCode = NotDo
			}
		} else {
			resCode = DBSQLERRO
		}

		if resCode == DBSUCCESS {
			if this.Param.writeTpe != DBDeleteDBObj {
				this.onWriteSuccess()
			}
		}
		//在这里设置因为中途有可能不成功
		//设置下一次开始设置位置
		this.valToDBValIndex = this.writeIndex
	}
	return resCode
}


func (this *DBObjWriteEvent) onWriteSuccess() {
	num :=  this.valToDBValIndex + this.writeSuccessNum
	xlog.Debug("设置dbVal位置 %d, 本次写到多少 %d",this.valToDBValIndex,num -1)
	for  i := this.valToDBValIndex; i < num; i++  {
		data := this.Param.DbObjs[i]
		if this.Param.writeTpe == DBInsertDBObj {
			data.setCreateToDB()
		}
		if this.Param.writeTpe != DBDeleteDBObj {
			//更新了需要把逻辑值设置为数据的值
			data.valToDBVal()
		}
	}

}

func (this *DBObjWriteEvent) getInsertObjsSql() string {

	filedSql := ""
	valsSql := ""
	this.writeSuccessNum = 0
	for  i := this.writeIndex; i < this.hasNum; i++  {
		data := this.Param.DbObjs[i]
		this.store.SetDBObjDefault(data)
		//第一个才添加
		if filedSql == "" {
			this.writeSuccessNum ++
			this.allSuccessNum ++
			filedSql = data.InsertSql(false)
			//第一次把字段加上
			valsSql = filedSql
			continue
		}

		add := "," + data.MoreInsertVal()
		//超过了长度就放
		if sqlLenIsMax(&valsSql, &add) {
			this.writeIndex = i //记录已经写到什么位置
			valsSql += ";"
			return valsSql
		}
		this.writeSuccessNum ++
		this.allSuccessNum ++
		valsSql += add
	}
	if valsSql != "" {
		valsSql += ";"
	}
	return valsSql
}

func (this *DBObjWriteEvent) getDelObjsSql() string {
	sql := ""
	for  i := this.writeIndex; i < this.hasNum; i++  {
		data := this.Param.DbObjs[i]
		this.store.SetDBObjDefault(data)
		add := data.DeleteSql()
		//超过了长度就放
		if sqlLenIsMax(&sql, &add) {
			this.writeIndex = i //记录已经写到什么位置
			return sql
		}
		this.allSuccessNum ++
		sql += add
	}
	return sql
}

func (this *DBObjWriteEvent) getUpdateObjsSql() string {
	this.writeSuccessNum = 0
	sql := ""
	for  i := this.writeIndex; i < this.hasNum; i++  {
		data := this.Param.DbObjs[i]
		this.store.SetDBObjDefault(data)
		add := data.UpdateSql()
		//超过了长度就放
		if sqlLenIsMax(&sql, &add) {
			this.writeIndex = i //记录已经写到什么位置
			xlog.Debug("失败 开始位置 %d , 本次成功数量 %d",this.writeIndex,this.writeSuccessNum)
			return sql
		}
		this.writeSuccessNum ++
		this.allSuccessNum ++
		sql += add
	}
	return sql
}

func (this *DBObjWriteEvent) EvenName() string {
	return "DBObjWriteEvent"
}

type DBObjWriteEventCb struct {
	dbEvent *DBObjWriteEvent
	code    DBErrorCode
}

//逻辑线程执行
func (this *DBObjWriteEventCb) Execute() {
	//线程队列回调
	if this.dbEvent.writeCb != nil {
		this.dbEvent.writeCb(this.code, this.dbEvent.Param)
	}
}

func (this *DBObjWriteEventCb) EvenName() string {
	return "DBObjWriteEventCb"
}
