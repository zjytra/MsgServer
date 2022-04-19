/*
创建时间: 2021/4/6 19:44
作者: zjy
功能介绍:

*/

package dbsys

import (
	"fmt"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"sync"
)

//处理 创建对象时,注册字段信息 动态拼接 update insert delete select
//1.创建对象时反射注册设置字段类型
//2.反射转换字段成对应接口,方便数据库比较值,获取字段sql,设置值
//3.获得数据库值得地址方便数据库执行字段扫描设置数据库值

//待完成
//4.写语句应该定时多久批量写入,每个表批量,关服时要同步数据
//5.
//逻辑层投递对象,及对应的操作类型

type DBErrorCode int8

const (
	DBSUCCESS DBErrorCode = 0
	NODATA    DBErrorCode = 1
	DBSQLERRO DBErrorCode = 2 //执行sql错误
	NotDo     DBErrorCode = 3 //执行sql无效
)

type CreateObjFun func() DBObJer

//数据库线程会使用到的接口
type DBObJer interface {
	Init()
	GetUID() int64
	GetDBRunFields() []DBRunColumnar
	GetDBValesAddr() []interface{}
	InsertSql(isOneEnd bool) string
	UpdateSql() string
	DeleteSql() string
	SelectSql() string
	//加载该表所有数据
	LoadAll() string
	MoreInsertVal() string
	//查询回调
	OnLoadForm(result DBErrorCode, param *DBObjQueryParam)
	//数据库删除
	OnDBDel(result DBErrorCode)
	//是否创建至数据库
	//是否创建至数据库
	GetTabName() string
	//设置表名
	SetTableName(tbName string)
	//逻辑值设置到数据库的值
	CreateObj() DBObJer
	isCreateToDB() bool
	//执行了insert操作就要写入数据库
	//与select
	setCreateToDB()
	dBValToVal()
	valToDBVal()
	//设置数据库
	setDB(dbStore *MySqlDBStore)
	//添加运行时的字段
	addRunCol(runCol DBRunColumnar)
	//设置select 的查询条件
	setSelectWhere(filedName string, runCol DBRunColumnar)
	//设置update 与delete条件
	setUpdateWhere(filedName string, runCol DBRunColumnar)
	//是否设置数据库操作接口
	isRegisterDBInterface() bool
	//设置数据库操作接口
	setIsRegister()

	//保证安全读map
	lock()
	unLock()
}

//数据库线程使用到的结构体
type DBObj struct {
	tbName         string
	dBFields       []DBRunColumnar
	delUpdateWhere map[string]DBRunColumnar
	selectWhere    map[string]DBRunColumnar
	//字段地址方便数据库查询结果scan
	dBAddr       []interface{}
	isCreateDB   bool
	dbStore      *MySqlDBStore
	isRegister   bool //是否注册数据库操作接口
	rwLock       sync.RWMutex //避免map竞争
	//子表任务 这里记录主要是为了方便玩家离线后找到子表数据库任务对象快速执行并删除批量未执行的任务
	subTableTasks map[string]*DBTableTask
}

//初始化
func (this *DBObj) Init() {
	this.delUpdateWhere = make(map[string]DBRunColumnar)
	this.selectWhere = make(map[string]DBRunColumnar)
	this.subTableTasks = make(map[string]*DBTableTask)
}

func (this *DBObj) GetUID() int64 {
	return 0
}
func (this *DBObj) GetTabName() string {
	return this.tbName
}

func (this *DBObj) OnLoadForm(result DBErrorCode, param *DBObjQueryParam) {

}

func (this *DBObj) OnDBDel(result DBErrorCode) {

}

func (this *DBObj) isCreateToDB() bool {
	return this.isCreateDB
}

func (this *DBObj) setCreateToDB() {
	this.isCreateDB = true
}

func (this *DBObj) setDB(store *MySqlDBStore) {
	this.dbStore = store
}
func (this *DBObj) SetTableName(tbName string) {
	this.tbName = tbName
}
func (this *DBObj) addRunCol(runCol DBRunColumnar) {
	//保存字段
	this.dBFields = append(this.dBFields, runCol)
	//保存设置数据库值地址
	this.dBAddr = append(this.dBAddr, runCol.GetDBValAddr())
}

func (this *DBObj) setUpdateWhere(filedName string, runCol DBRunColumnar) {
	if this.delUpdateWhere == nil {
		this.delUpdateWhere = make(map[string]DBRunColumnar)
	}
	this.delUpdateWhere[filedName] = runCol
}

func (this *DBObj) setSelectWhere(filedName string, runCol DBRunColumnar) {
	if this.selectWhere == nil {
		this.selectWhere = make(map[string]DBRunColumnar)
	}
	this.selectWhere[filedName] = runCol
}

func (this *DBObj) GetDBRunFields() []DBRunColumnar {
	return this.dBFields
}

func (this *DBObj) GetDBValesAddr() []interface{} {
	return this.dBAddr
}

func (this *DBObj) InsertSql(isOneEnd bool) string {
	if this.dBFields == nil {
		return ""
	}
	flen := len(this.dBFields)
	if flen == 0 {
		return ""
	}
	filedSql := ""
	valsSql := ""
	for i, field := range this.dBFields {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		//自增不写入
		if pCol.HasAUTO_INCR() {
			continue
		}
		filedSql += pCol.COLUMN_NAME
		valsSql += pCol.GetValStr(field.DBGetVal())
		if i < flen-1 {
			filedSql += ","
			valsSql += ","
		}
	}
	if isOneEnd {
		return fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", this.tbName, filedSql, valsSql)
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", this.tbName, filedSql, valsSql)
}

//func (this *DBObj) InsertFieldSql() string {
//	if this.dBFields == nil {
//		return ""
//	}
//	flen := len(this.dBFields)
//	if flen == 0 {
//		return ""
//	}
//	filedSql := ""
//	for i, field := range this.dBFields {
//		pCol := field.GetCol()
//		if pCol == nil {
//			xlog.Error("%s field index %s is nil", this.tbName, i)
//			continue
//		}
//		filedSql += pCol.COLUMN_NAME
//		if i < flen-1 {
//			filedSql += ","
//		}
//	}
//	return fmt.Sprintf("INSERT INTO %s(%s) ", this.tbName, filedSql)
//}

func (this *DBObj) MoreInsertVal() string {
	if this.dBFields == nil {
		return ""
	}
	flen := len(this.dBFields)
	if flen == 0 {
		return ""
	}
	valsSql := ""
	for i, field := range this.dBFields {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		//自增不写入
		if pCol.HasAUTO_INCR() {
			continue
		}
		valsSql += pCol.GetValStr(field.DBGetVal())
		if i < flen-1 {
			valsSql += ","
		}

	}
	return fmt.Sprintf("(%s)", valsSql)
}

func (this *DBObj) UpdateSql() string {
	if this.dBFields == nil {
		return ""
	}
	flen := len(this.dBFields)
	if flen == 0 {
		return ""
	}
	if this.delUpdateWhere == nil {
		xlog.Error("%s not delUpdateWhere", this.tbName)
		return ""
	}
	whereLen := len(this.delUpdateWhere)
	if whereLen == 0 {
		xlog.Error("%s not delUpdateWhere", this.tbName)
		return ""
	}
	//没有多余的字段更新
	if flen <= whereLen {
		xlog.Error("%s not more field update", this.tbName)
		return ""
	}
	//查看改变字段的数量
	changeLen := this.getChangeLen()
	if changeLen == 0 {
		return ""
	}
	filedSql := ""
	for i, field := range this.dBFields {
		if !field.IsChange() {
			continue
		}
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		//作为条件的字段不用更新
		_, ok := this.delUpdateWhere[pCol.COLUMN_NAME]
		if ok {
			continue
		}
		filedSql += pCol.COLUMN_NAME + " = " + pCol.GetValStr(field.DBGetVal())
		changeLen--
		if changeLen > 0 {
			filedSql += ","
		}
	}
	//没有改变的字段
	if filedSql == "" {
		return ""
	}
	//条件语句
	whereSql := ""
	var i int
	for _, field := range this.delUpdateWhere {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			i++
			continue
		}
		whereSql += pCol.COLUMN_NAME + " = " + pCol.GetValStr(field.DBGetVal())
		if i < whereLen-1 {
			whereSql += " AND "
		}
		i++
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", this.tbName, filedSql, whereSql)
}

func (this *DBObj) getChangeLen() int {
	changeLen := 0
	for i, field := range this.dBFields {
		if !field.IsChange() {
			continue
		}
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		//作为条件的字段不用更新
		_, ok := this.delUpdateWhere[pCol.COLUMN_NAME]
		if ok {
			continue
		}
		changeLen++
	}
	return changeLen
}

func (this *DBObj) DeleteSql() string {
	if this.delUpdateWhere == nil {
		xlog.Error("%s not delUpdateWhere", this.tbName)
		return ""
	}
	//没有多余的字段更新
	whereLen := len(this.delUpdateWhere)
	if whereLen <= 0 {
		xlog.Error("%s not del where", this.tbName)
		return ""
	}
	//条件语句
	whereSql := ""
	var i int
	for _, field := range this.delUpdateWhere {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			i++
			continue
		}
		whereSql += pCol.COLUMN_NAME + " = " + pCol.GetValStr(field.DBGetVal())
		if i < whereLen-1 {
			whereSql += " AND "
		}
		i++
	}
	return fmt.Sprintf("DELETE FROM %s WHERE %s;", this.tbName, whereSql)
}

func (this *DBObj) SelectSql() string {
	if this.dBFields == nil {
		return ""
	}
	flen := len(this.dBFields)
	if flen == 0 {
		return ""
	}
	if this.selectWhere == nil {
		xlog.Error("%s not delUpdateWhere", this.tbName)
		return ""
	}
	whereLen := len(this.selectWhere)
	if whereLen == 0 {
		xlog.Error("%s not selectWhere", this.tbName)
		return ""
	}

	filedSql := ""
	for i, field := range this.dBFields {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		filedSql += pCol.COLUMN_NAME
		if i < flen-1 {
			filedSql += ","
		}
	}
	//没有改变的字段
	if filedSql == "" {
		return ""
	}
	//条件语句
	whereSql := ""
	var i int
	for _, field := range this.selectWhere {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			i++
			continue
		}
		whereSql += pCol.COLUMN_NAME + " = " + pCol.GetValStr(field.DBGetVal())
		if i < whereLen-1 {
			whereSql += "AND"
		}
		i++
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s;", filedSql, this.tbName, whereSql)
}

func (this *DBObj) LoadAll() string {
	if this.dBFields == nil {
		return ""
	}
	flen := len(this.dBFields)
	if flen == 0 {
		return ""
	}
	if this.selectWhere == nil {
		xlog.Error("%s not delUpdateWhere", this.tbName)
		return ""
	}
	whereLen := len(this.selectWhere)
	if whereLen == 0 {
		xlog.Error("%s not selectWhere", this.tbName)
		return ""
	}

	filedSql := ""
	for i, field := range this.dBFields {
		pCol := field.GetCol()
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.tbName, i)
			continue
		}
		filedSql += pCol.COLUMN_NAME
		if i < flen-1 {
			filedSql += ","
		}
	}
	//没有改变的字段
	if filedSql == "" {
		return ""
	}
	return fmt.Sprintf("SELECT %s FROM %s;", filedSql, this.tbName)
}

//当数据从数据库加载出来后
func (this *DBObj) dBValToVal() {
	if this.dBFields == nil {
		return
	}
	for _, field := range this.dBFields {
		field.dbValSetVal()
	}
}

//执行了更新操作后更新数据库值
func (this *DBObj) valToDBVal() {
	if this.dBFields == nil {
		return
	}
	for _, field := range this.dBFields {
		field.valSetDBVal()
	}
}

//
//func (this *DBObj)Cb()  {
//	fmt.Println("DBObj Cb")
//}

//是否设置数据库操作接口
func (this *DBObj) isRegisterDBInterface() bool {
	return this.isRegister
}

//设置数据库操作接口
func (this *DBObj) setIsRegister() {
	this.isRegister = true
}

//是否设置数据库操作接口
func (this *DBObj) lock()  {
	this.rwLock.Lock()
}

//设置数据库操作接口
func (this *DBObj) unLock() {
	this.rwLock.Unlock()
}



//子类必须实现这个方法好让数据库去创建
func (this *DBObj) CreateObj() DBObJer {
	return nil
}

//添加子表任务
func (this *DBObj) AddSubInsertData(in DBObJer) {
	tasks, ok := this.subTableTasks[in.GetTabName()]
	if !ok {
		tasks = this.AddTbTask(in)
	}
	tasks.AddInsertData(in)
}

func (this *DBObj) AddTbTask(in DBObJer) *DBTableTask {
	tasks := NewDBTableTask(in.GetTabName())
	this.subTableTasks[in.GetTabName()] = tasks
	return tasks
}

func (this *DBObj) AddSubUpdateData(in DBObJer) {
	//还未创建至数据库就不用更新
	tasks, ok := this.subTableTasks[in.GetTabName()]
	if !ok {
		tasks = this.AddTbTask(in)
	}
	tasks.AddUpdateData(in)
}

func (this *DBObj) AddSubDelData(in DBObJer) {
	//还未创建至数据库就不用再创建了
	tasks, ok := this.subTableTasks[in.GetTabName()]
	if !ok {
		tasks = this.AddTbTask(in)
	}
	tasks.AddDelData(in)
}
