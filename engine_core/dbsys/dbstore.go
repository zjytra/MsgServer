/*
创建时间: 2020/3/3
作者: zjy
功能介绍:
数据库逻辑封装读写分离
*/

package dbsys

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/devlop/xutil/osutil"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/timingwheel"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
	"time"
)

const (
	DBQueryRowsCB_Event              = 1 // 查询返回原始的Rows
	DBQueryRowToMapCb_Event          = 2 // 查询单行数据返回map
	DBQueryRowsToMapArrCb_Event      = 3 // 查询多行数据返回map切片
	DBQueryRowToStructCb_Event       = 4 // 查询返回单行结构体事件
	DBQueryRowsToStructSliceCb_Event = 5 // 查询多行结构体事件
	DBQueryMoreReusltMapCb_Event = 6 // 多个结果集返回map
	DBQueryMoreReusltStructArrCb_Event = 7 // 多个结果集返回多个结构体切片
	DBQueryDBObj      = 8 // 查询单个结构体
	DBQueryMoreDBObj      = 9 // 查询多个结构体
	DBDeleteDBObj         = 10 // 删除结构体
	DBUpdateDBObj         = 11 // 更新
	DBInsertDBObj         = 12 // 插入
	DBEvent_max                      = 13
)



// 封装数据库处理
type MySqlDBStore struct {
	db          *sql.DB
	dbName      string               //动态获取配置避免程序运行过程中修改配置
	quereyEvent *dispatch.AsyncQueue // 数据库查询队列
	writeEvent  *dispatch.AsyncQueue // 数据库写队列
	//表操作任务,等到数据批量更新
	tableTasks map[string]*DBTableTask
	dbTables   map[string]*DBTable //当前数据库注册的所有的表
	tableIdxs  map[string][]*ColumnInfo
	updateTimer *timingwheel.Timer
}

func NewMySqlDBStore(dbname string) *MySqlDBStore {
	dbstore := new(MySqlDBStore)
	dbconf := csvdata.GetDbCfgPtr(dbname)
	if dbconf == nil {
		panic(dbname + "数据库配置为null")
	}
	dbstore.dbName = dbconf.Dbname
	xlog.Debug(" NewMySqlDBStore dbname %v", dbconf.Dbname)
	dbstore.tableTasks = make(map[string]*DBTableTask)
	dbstore.dbTables = make(map[string]*DBTable)
	dbstore.tableIdxs = make(map[string][]*ColumnInfo)
	return dbstore
}


func (this *MySqlDBStore)StartTimer() {
	//开启延迟任务
	scheduler := &timingwheel.EveryScheduler{time.Second * 10}
	scheduler2 := &timingwheel.EveryScheduler{time.Second * 10}
	scheduler3 := &timingwheel.EveryScheduler{time.Second * 10}
	this.updateTimer = timingwheel.ScheduleFuncInQueue(scheduler, this.executeInsert, dispatch.MainQueue)
	timingwheel.ScheduleFuncInQueue(scheduler2, this.executeDelete, dispatch.MainQueue)
	timingwheel.ScheduleFuncInQueue(scheduler3, this.executeUpdate, dispatch.MainQueue)
}

func (this *MySqlDBStore) OpenDB() error {
	// 异步连接
	dbconf := csvdata.GetDbCfgPtr(this.dbName)
	if dbconf == nil {
		return errors.New(fmt.Sprintf("OpenDB %s Cfg is nil", this.dbName))
	}
	if err := this.openDB(); err != nil {
		xlog.Error("open db error: %v ", err)
		return err
	}

	if dbconf.Readnum == 0 || dbconf.Writenum == 0 {
		xlog.Debug("数据库 %v 读线程数: %v ,数据库写线程数 %v  ",dbconf.Dbname, dbconf.Readnum,dbconf.Writenum)
	}

	xlog.Debug("数据库 %v  连接succeed", dbconf.Dbname)
	//开启读写队列
	if dbconf.Readnum > 0 && this.quereyEvent == nil {
		this.quereyEvent = dispatch.NewAsyncQueue(500,int(dbconf.Readnum))
		//开启读写队列
		this.quereyEvent.Start()
	}

	if dbconf.Writenum > 0 && this.writeEvent == nil {
		//多个线程写
		this.writeEvent = dispatch.NewAsyncQueue(500,int(dbconf.Writenum))
		this.writeEvent.Start()
	}

	return nil
}

func (this *MySqlDBStore) openDB() error {
	dbconf := csvdata.GetDbCfgPtr(this.dbName)
	if dbconf == nil {
		return errors.New("dbconf is nil")
	}
	DataSoureName := GetMysqlDataSourceName(dbconf)
	if strutil.StringIsNil(DataSoureName) {
		return errors.New(fmt.Sprintf("%v 数据库连接信息为nil",dbconf.Dbname))
	}
	var Erro error
	this.db, Erro = sql.Open("mysql", DataSoureName)
	if xutil.IsError(Erro) {
		return Erro
	}
	this.db.SetMaxOpenConns(dbconf.Maxopenconns)
	this.db.SetMaxIdleConns(dbconf.Maxidleconns)
	if erro := this.db.Ping(); xutil.IsError(erro) {
		this.CloseDB()
		return erro
	}
	return Erro
}

// 关闭数据库
func (this *MySqlDBStore) CloseDB() {
	 erro := this.db.Close()
	if erro != nil {
		xlog.Error("CloseDB %v",erro)
	}
	 //退出
	if this.quereyEvent != nil {
		this.quereyEvent.Release()
	}
	if this.writeEvent != nil {
		this.writeEvent.Release()
	}
	xlog.Debug("CloseDB %v",this.dbName)
}

func (this *MySqlDBStore) Query(query string, args ...interface{}) (row *sql.Rows, erro error) {
	row, erro = this.db.Query(query, args ...)
	if erro != nil {
		xlog.Debug( "%s", osutil.GetRuntimeFileAndLineStr(1))
		xlog.Debug( "db.Query sql =%s", query)
		xlog.Debug( "erro %v", erro)
		return
	}
	// erro = row.Close() 解析完数据才能关闭
	return
}

func (this *MySqlDBStore) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	row = this.db.QueryRow(query, args ...)
	if row.Err() != nil {
		xlog.Debug( "%s", osutil.GetRuntimeFileAndLineStr(1))
		xlog.Debug( "db.Query sql =%s", query)
		xlog.Debug( "erro %v", row.Err())
		return
	}
	// erro = row.Close() 解析完数据才能关闭
	return
}

func (this *MySqlDBStore) Execute(query string, args ...interface{}) (result sql.Result, erro error) {
	result, erro = this.db.Exec(query, args ...)
	if erro != nil {
		xlog.Debug( "%s", osutil.GetRuntimeFileAndLineStr(1))
		xlog.Debug( "db.Execute sql =%s", query)
		xlog.Debug( "erro %v", erro)
	}
	return
}

func (this *MySqlDBStore) CheckTableExists(tableName string) bool {
	if this.db == nil {
		xlog.Debug("CheckTableExists this.db is nil")
		return false
	}
	rows, erro := this.db.Query("SELECT t.TABLE_NAME FROM information_schema.TABLES AS t WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ", this.dbName, tableName)
	if xutil.IsError(erro) {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

//设置数据库可以操作的对象
func (this *MySqlDBStore) SetDBObjDefault(runObj DBObJer) {
	//已经注册过相关接口就不用再处理
	this.SetDBObj(runObj,"")
}

//设置对象数据库相关接口
func (this *MySqlDBStore) SetDBObj(runObj DBObJer,tbName string) {
	//已经注册过相关接口就不用再处理
	runObj.lock()
	if runObj.isRegisterDBInterface() {
		runObj.unLock()
		return
	}
	runObj.setDB(this)
	runType := reflect.TypeOf(runObj)
	elT := runType.Elem()
	if runObj.GetTabName() == "" {
		if tbName == "" {
			runObj.SetTableName(elT.Name())
		}else {
			runObj.SetTableName(tbName)
		}
	}
	//注册表的字段方便所有对象字段共用col信息
	dbTab := this.GetDBTable(runObj.GetTabName())
	if dbTab == nil {
		return
	}
	valOf := reflect.ValueOf(runObj)
	//先把指针转成结构体
	elVal := valOf.Elem()
	//反射解析字段
	fieldLen := elT.NumField()
	for i := 0; i < fieldLen; i++ {
		filed := elT.Field(i)
		filedName := filed.Name
		if filedName == "DBObj" {
			continue
		}
		//有忽略字段
		_, hasIgnore := filed.Tag.Lookup("ign")
		if hasIgnore {
			continue
		}
		colInfo, ok := dbTab.FieldsMap[filedName]
		if !ok {
			xlog.Warning("%s 未找到字段%s", dbTab.TbName, filedName)
			continue
		}
		//字段是结构体
		valFild := elVal.Field(i)
		if valFild.Kind() != reflect.Struct {
			xlog.Warning("%s Filed %s not Struct", dbTab.TbName, filedName)
			continue
		}
		//获取结构体的地址
		if !valFild.CanAddr() {
			xlog.Warning("%s Filed %s not addr", dbTab.TbName, filedName)
			continue
		}
		fildAddr := valFild.Addr()
		//地址转指针
		//断言字段接口
		//方便后续处理
		runCol, colIOk := fildAddr.Interface().(DBRunColumnar)
		if !colIOk {
			xlog.Warning("%s Filed %s not DBRunColumnar", dbTab.TbName, filedName)
			continue
		}
		//设置字段信息
		runCol.SetColInfo(colInfo)
		//查询字段
		if colInfo.IsQueryWhere() { //select 字段
			runObj.setSelectWhere(filedName, runCol)
		}
		if colInfo.IsUpdateOrDelWhere() { //del 字段
			runObj.setUpdateWhere(filedName, runCol)
		}
		runCol.initDBField()
		runObj.addRunCol(runCol)
	}
	runObj.setIsRegister()
	runObj.unLock()
}



func (this *MySqlDBStore) RowsToDBObjArr(rows *sql.Rows, table *DBTable, resCode *DBErrorCode) []DBObJer {
	var objs []DBObJer
	for rows.Next() {
		//创建一个新对象
		obj := table.createFun()
		if obj == nil {
			xlog.Debug("%s表对象没有实现创建的对象方法",table.TbName)
			break
		}
		this.SetDBObj(obj, table.TbName)
		val := obj.GetDBValesAddr()
		sErro := rows.Scan(val...)
		if sErro != nil {
			xlog.Debug("查询表 %s 赋值 错误 %v",obj.GetTabName(), sErro)
		}
		obj.setCreateToDB()
		obj.dBValToVal()
		objs = append(objs, obj)
		*resCode = DBSUCCESS
	}
	return objs
}



func (this *MySqlDBStore)DelayInsert(in DBObJer) {
	if in == nil {
		return
	}
	this.SetDBObjDefault(in)
	task := this.getTableTask(in)
	if task == nil {
		xlog.Error("DelayInsert %s 表任务未找到",in.GetTabName())
		return
	}
	task.AddInsertData(in)
}

func (this *MySqlDBStore) getTableTask(in DBObJer) *DBTableTask {
	tabName := in.GetTabName()
	tableTask, ok := this.tableTasks[tabName]
	if !ok {
		tableTask = NewDBTableTask(tabName)
		this.tableTasks[tabName] = tableTask
	}
	return tableTask
}

func (this *MySqlDBStore)DelayUpdate(in DBObJer) {
	if in == nil {
		return
	}
	this.SetDBObjDefault(in)
	task := this.getTableTask(in)
	if task == nil {
		xlog.Error("DelayUpdate %s 表任务未找到",in.GetTabName())
		return
	}
	task.AddUpdateData(in)
}

func (this *MySqlDBStore)DelayDelData(in DBObJer) {
	if in == nil {
		return
	}
	this.SetDBObjDefault(in)
	task := this.getTableTask(in)
	if task == nil {
		xlog.Error("DelayDelData %s 表任务未找到",in.GetTabName())
		return
	}
	task.AddDelData(in)
}


//执行任务延迟任务
func (this *MySqlDBStore)executeInsert()  {
	for _, table := range this.tableTasks {
		sql := table.InsertMoreSql()
		this.asyncExecuteDelay(sql)
	}
}

//执行任务延迟任务
func (this *MySqlDBStore)executeUpdate()  {
	for _, table := range this.tableTasks {
		sql := table.MoreUpdateSql()
		this.asyncExecuteDelay(sql)
	}

}

//执行任务延迟任务
func (this *MySqlDBStore)executeDelete()  {
	for _, table := range this.tableTasks {
		sql := table.MoreDelSql()
		this.asyncExecuteDelay(sql)
	}
}

func (this *MySqlDBStore) asyncExecuteDelay(query string)  {
	if strutil.StringIsNil(query) {
		return
	}
	event := NewDBWriteEvent()
	if event == nil {
		return
	}
	event.store = this
	event.query = query
	this.writeEvent.AddEvent(event)
}
