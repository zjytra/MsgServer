/*
创建时间: 2021/7/30 0:12
作者: zjy
功能介绍:
数据库的任务,便于分用户存插入删除,更新任务
一个表的的任务
*/

package dbsys

import (
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

type DBTableTask struct {
	tableName string
	//可以用一般的map 单线程处理
	insertData map[DBObJer] DBObJer
	updateData map[DBObJer] DBObJer
	delData    map[DBObJer] DBObJer
}

func NewDBTableTask(name string) *DBTableTask {
	tb := new(DBTableTask)
	tb.tableName = name
	tb.insertData = make(map[DBObJer] DBObJer)
	tb.updateData = make(map[DBObJer] DBObJer)
	tb.delData = make(map[DBObJer] DBObJer)
	return tb
}

func (this *DBTableTask) AddInsertData(in DBObJer) {
	//在删除列表中就不要加了
	if this.hasDel(in) {
		return
	}
	_, ok := this.insertData[in]
	if ok {
		xlog.Warning("%s is exist %d", in.GetTabName(), in.GetUID())
	}
	this.insertData[in] = in
}

//如果再插入队列里面就更新
func (this *DBTableTask) hasInsert(in DBObJer) bool {
	_, ok := this.insertData[in]
	return ok
}

//如果再插入队列里面就更新
func (this *DBTableTask) hasDel(in DBObJer) bool {
	_, ok := this.delData[in]
	return ok
}

func (this *DBTableTask) AddUpdateData(in DBObJer) {
	//还未创建至数据库就不用更新
	if this.hasInsert(in) {
		return
	}
	if this.hasDel(in) {
		return
	}
	this.updateData[in] =  in
}

func (this *DBTableTask) AddDelData(in DBObJer) {
	//还未创建至数据库就不用反正要删除
	if this.hasInsert(in) {
		//删除还未写进数据库的
		this.RemoveInsertData(in)
		return
	}
	this.delData[in] = in
}

func (this *DBTableTask) RemoveInsertData(in DBObJer) {
	delete(this.insertData,in)
}

func (this *DBTableTask) RemoveUpdateData(in DBObJer) {
	delete(this.updateData,in)
}

func (this *DBTableTask) RemoveDelData(in DBObJer) {
	delete(this.delData,in)
}

func (this *DBTableTask) InsertMoreSql() string {
	filedSql := ""
	valsSql := ""
	for _, data := range this.insertData {
		//第一个才添加
		if filedSql == "" {
			filedSql = data.InsertSql(false)
			data.setCreateToDB()
			data.valToDBVal()
			this.RemoveInsertData(data)
			//第一次把字段加上
			valsSql = filedSql
			continue
		}
		add := "," + data.MoreInsertVal()
		//超过了长度就放
		if sqlLenIsMax(&valsSql, &add) {
			valsSql += ";"
			return valsSql
		}
		data.setCreateToDB()
		data.valToDBVal()
		this.RemoveInsertData(data)
		valsSql += add
	}
	if valsSql != "" {
		valsSql += ";"
	}
	return valsSql
}

func (this *DBTableTask) MoreUpdateSql() string {
	updateSql := ""
	for _, data := range this.updateData {
		//if !Account.isCreateToDB() {
		//	this.RemoveUpdateData(Account)
		//	continue
		//}
		add := data.UpdateSql()
		//超过了长度就放
		if sqlLenIsMax(&updateSql, &add) {
			return updateSql
		}
		data.valToDBVal()
		this.RemoveUpdateData(data)
		updateSql += add
	}

	return updateSql
}

func (this *DBTableTask) MoreDelSql() string {
	delSql := ""
	for _, data := range this.delData {
		//if !Account.isCreateToDB() {
		//	this.RemoveDelData(Account)
		//	continue
		//}
		add := data.DeleteSql()
		//超过了长度就放
		if sqlLenIsMax(&delSql, &add) {
			return delSql
		}
		this.RemoveDelData(data)
		delSql += add
	}
	return delSql
}
