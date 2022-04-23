// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 1.处理mysql相关逻辑
// 2.简单orm操作
// 3.多协程操作
// 论数据库操作
// 两种拼sql 方案,主线程拼sql,写任务线程拼,

// 先讲写任务线程拼sql有可能导致逻辑线程创建的数据,同时写线程正在写数据,然后后面又跟来一个update操作进写协程队,
// 执行更新任务时查看是否还未创建进数据库,未创建什么也不做,指针会拿到最新的内存数据（这个可以验证下）
// 解决以上问题方案dbcoroutine协程处理dbtabletask的任务 分对应的协程处理 按玩家uid 取余处理,保证单协程写入是某一部分的玩家同步的,可以不用加锁 ,
// 然后就是配置定时器多久写一次数据,写入后就从需要写的任务队列中移除

// 读的任务
// 1.玩家登录查询一次,中途查询较少
// 2.中途写入与更新较多
// 3.查单个对象,
// 4.代条件查多个对象
package dbsys

import (
	"database/sql"
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
	"strings"
	"sync"
)

const(
	GameStatisticsDBName =  "game_statisticsdb"
	GameAccountLogDBName = "accountlog"
	GameAccountdbName =  "accountdb"
	MysqlDBName = "mysql"
)






var (
	GameDB            *MySqlDBStore //游戏库
	GameAccountDB     *MySqlDBStore //账号库
	GameAccountLogDB  *MySqlDBStore //账号日志库
	LogDB             *MySqlDBStore //日志库
	Game_statisticsDB *MySqlDBStore //分析库
	MysqlDB           *MySqlDBStore //mysql库 用来同步表格
	AccountMutex      sync.Mutex   //账号锁保证创建账号唯一性
)


func InitGameDB() {
	if GameDB != nil {
		GameDB.CloseDB()
	}
	// 创建数据库相关操作
	GameDB = NewMySqlDBStore(csvdata.OutNetConf.Dbname)
}


func InitLogDB() {
	if LogDB != nil {
		LogDB.CloseDB()
	}
	// 创建数据库相关操作
	LogDB = NewMySqlDBStore(csvdata.OutNetConf.Logdbname)
}

func InitStatisticsDB() {
	if Game_statisticsDB != nil {
		Game_statisticsDB.CloseDB()
	}
	// 创建数据库相关操作
	Game_statisticsDB = NewMySqlDBStore(GameStatisticsDBName)
}

func InitAccountDB() {
	if GameAccountDB != nil {
		GameAccountDB.CloseDB()
	}
	// 创建数据库相关操作
	GameAccountDB = NewMySqlDBStore(csvdata.OutNetConf.Dbname)
}

func InitAccountLogDB() {
	if GameAccountLogDB != nil {
		GameAccountLogDB.CloseDB()
	}
	// 创建数据库相关操作
	GameAccountLogDB = NewMySqlDBStore(csvdata.OutNetConf.Logdbname)
}

//获取连接字符串
func GetMysqlDataSourceName(dbinfo *csvdata.DbCfg) string {
	if dbinfo == nil {
		fmt.Println("dbinfo is nil")
		return ""
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=1s&charset=%s&multiStatements=true",
		dbinfo.Dbusername,
		dbinfo.Dbpwd,
		dbinfo.Ip,
		dbinfo.Dbport,
		dbinfo.Dbname,
		dbinfo.Char_set)
}

//将多行查询的结果转换为string 数组
func RowsToStrArrSlice(rows *sql.Rows) [][]string{
	defer rows.Close()
	if rows == nil{
		xlog.Error("RowsToStrArrSlice rows is nil ")
		return nil
	}
	columns,erro := rows.Columns()  //获取查询出的字段
	if erro != nil {
		xlog.Error("RowsToStrArrSlice Columns %v",erro)
		return nil
	}
	columnsCount := len(columns) //字段个数
	if columnsCount == 0 {
		xlog.Error("RowsToStrArrSlice columnsCount == 0")
		return nil
	}

	var strrows [][]string
	for rows.Next() {
		//拼接[][]string 这里作为桶装东西
		values := make([]string, columnsCount)  //每行的值
		scans := make([]interface{}, columnsCount)
		for i := range values {
			scans[i] = &values[i]  //这里存储装数据的地址 给扫描的时候使用
		}
		erro = rows.Scan(scans...) //传入values的地址进行赋值
		if erro != nil {
			xlog.Error("RowsToStrArrSlice Scan %v",erro)
		    continue
		}
		strrows = append(strrows, values)
	}
	return strrows
}

//将单行查询的结果转换为string 数组
func RowToStringSlice(rows *sql.Rows) []string{
	defer rows.Close()
	_, erro, columnsCount := CheckOneRow("RowToStringSlice",rows)
	if erro != nil || columnsCount == 0 {
		return nil
	}
	//装数据的
	values := make([]string, columnsCount)
	//把每个地址赋值进去
	scans := make([]interface{}, columnsCount)
	for i := range scans {
		scans[i] = &values[i]
	}
	erro = rows.Scan(scans...)
	if erro != nil {
		xlog.Error("RowToStringSlice Scan %v",erro)
		return nil
	}
	return values
}


//将单行查询的结果转换为map key=字段名称
func RowToMap(rows *sql.Rows) map[string]string{
	columns, erro, columnsCount := CheckOneRow("RowToMap",rows)
	if erro != nil || columnsCount == 0 {
		return nil
	}
	//装数据的
	values := make([]string, columnsCount)
	//把每个地址赋值进去
	scans := make([]interface{}, columnsCount)
	for i := range scans {
		scans[i] = &values[i]
	}
	erro = rows.Scan(scans...)
	if erro != nil {
		xlog.Error("RowToStringSlice Scan %v", erro)
		return nil
	}
	fieldAndVal := make(map[string]string)
	for i, column := range columns {
		fieldAndVal[column] = values[i]
	}
	return fieldAndVal
}

//将多行查询的结果转换为map数组,其中一行去slice下表map[数据库字段名称]字段值
//这个api记得最后调用 rows.Close()
func RowsToStrMapArr(rows *sql.Rows) []map[string]string{
	if rows == nil{
		xlog.Error("RowsToStrArrSlice rows is nil ")
		return nil
	}
	columns,erro := rows.Columns()  //获取查询出的字段
	if erro != nil {
		xlog.Error("RowsToStrArrSlice Columns %v",erro)
		return nil
	}
	columnsCount := len(columns) //字段个数
	if columnsCount == 0 {
		xlog.Error("RowsToStrArrSlice columnsCount == 0")
		return nil
	}

	var strrows []map[string]string
	for rows.Next() {
		values := make([]string, columnsCount)  //每行的值
		//拼接[][]string 这里作为桶装东西
		scans := make([]interface{}, columnsCount)
		for i := range scans {
			scans[i] = &values[i]  //这里存储装数据的地址 给扫描的时候使用
		}
		erro = rows.Scan(scans...) //传入values的地址进行赋值
		if erro != nil {
			xlog.Error("RowsToStrArrSlice Scan %v",erro)
			continue
		}
		//将查询结果赋值给map
		oneRowMap := make(map[string]string)
		for i, column := range columns {
			oneRowMap[column] = values[i]
		}
		strrows = append(strrows, oneRowMap)
	}
	return strrows
}

//多个结果集,多行查询转换为map数组,其中一行去slice下表map[表明第几个结果集]gamemap[数据库字段名称]字段值
func MoreResultRowsToStrMapArr(rows *sql.Rows) map[int][]map[string]string{
	allRes :=make(map[int][]map[string]string)
	oneResult := RowsToStrMapArr(rows)
	i := 1
	allRes[i] = oneResult
	//如果还有下一个结果集继续遍历
	for rows.NextResultSet() {
		oneResult := RowsToStrMapArr(rows)
		i ++
		allRes[i] = oneResult
	}
	return allRes
}


func CheckOneRow(funName string,rows *sql.Rows) ([]string, error, int) {
	if rows == nil {
		xlog.Debug(funName + " rows is nil ")
		return nil, nil, 0
	}
	if !rows.Next() {
		return nil, nil, 0
	}
	columns, erro := rows.Columns() //获取查询出的字段
	if erro != nil {
		xlog.Error(funName + " Columns %v", erro)
		return nil, nil, 0
	}
	columnsCount := len(columns) //字段个数
	if columnsCount == 0 {
		xlog.Error(funName + " columnsCount == 0")
		return nil, nil, 0
	}
	return columns, erro, columnsCount
}


//查询结果转为结构体
//@param to 结构体模型
//@param rows 数据库查询结果
func RowToStruct(rows *sql.Rows, to interface{}) interface{} {
	if rows == nil{
		xlog.Error("RowToStruct rows is nil ")
		return nil
	}
	if to == nil {
		xlog.Error("RowToStruct to= nil")
		return nil
	}
	tp := reflect.TypeOf(to)
	if tp.Kind()  != reflect.Struct {
		xlog.Warning("%s 不是结构体",tp.Name())
		return nil
	}
	val := reflect.New(tp)//创建个对象
	if val.Kind() != reflect.Ptr || val.CanAddr() {
		xlog.Error("赋值对象不是指针 %v",val.Kind())
		return nil
	}
	valElem := val.Elem()
	column_names,column_types := GetRowsColsAndTypes(rows)
	if column_names == nil || column_types == nil {
		return nil
	}
	colLen := len(column_names)
	valElemLen := valElem.NumField()
	if colLen != valElemLen {
		xlog.Error("数据库列数 %v 与结构体列数 %v 不匹配 ",colLen,valElemLen)
		return nil
	}
	scan_dest := GetStructScans(&valElem,column_names,column_types) //获取结构体的扫描地址
	if  scan_dest == nil{
		return nil
	}
	for rows.Next() {
		erro := rows.Scan(scan_dest...) //传入values的地址进行赋值
		if erro != nil {
			xlog.Error("RowToStruct Scan %v",erro)
			continue
		}
	}

	return val.Interface()
}


//查询结果转为结构体
//@param to 最好是结构体类型
//@return []interface{}  遍历时必须用指针断言
//这个api记得最后调用 rows.Close()
func RowsToStructSlice(rows *sql.Rows,to interface{})[]interface{}{
	if rows == nil{
		xlog.Error("RowsToStructSlice rows is nil ")
		return nil
	}
	if to == nil {
		xlog.Error("RowsToStructSlice to= nil")
		return nil
	}
	tp := reflect.TypeOf(to)
	if tp.Kind()  != reflect.Struct {
		xlog.Warning("%s 不是结构体",tp.Name())
		return nil
	}
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	column_names,column_types := GetRowsColsAndTypes(rows)
	if column_names == nil || column_types == nil {
		return  nil
	}
	colLen := len(column_names)
	valElemLen := tp.NumField()
	if colLen != valElemLen {
		xlog.Error("数据库列数 %v 与结构体列数 %v 不匹配 ",colLen,valElemLen)
		return nil
	}
	
	var allData  []interface{}
	for rows.Next() {
		val := reflect.New(tp)//创建个对象
		if val.Kind() != reflect.Ptr || val.CanAddr() {
			xlog.Error("赋值对象不是指针 %v",val.Kind())
			return nil
		}
		// xlog.Debug("val.Elem() = %v",val.Elem())
		valElem := val.Elem()
		//if valElem.Kind() != reflect.Struct {
		//	xlog.Error("valElem.Kind() %v",valElem.Kind())
		//	continue
		//}
		scan_dest := GetStructScans(&valElem,column_names,column_types)
		if scan_dest == nil {
			continue
		}
		erro := rows.Scan(scan_dest...) //传入values的地址进行赋值
		if erro != nil {
			xlog.Error("RowsToStrArrSlice Scan %v",erro)
			continue
		}
		// xlog.Error("Account = %v ",val)
		allData = append(allData,val.Interface())
	}
	return allData
}

//多个结果集赋值
//@param to 可以是结构体类型 也可以是指针类型,主要创建实例使用
//@return []interface{}  遍历时必须用指针断言
func MoreResultRowsToStructArr(rows *sql.Rows,toArr []interface{})map[int][]interface{}  {
	if rows == nil{
		xlog.Error("RowsToStructSlice rows is nil ")
		return nil
	}
	if toArr == nil {
		xlog.Error("MoreResultRowsToStructArr to= nil")
		return nil
	}
	allDataMap := make(map[int][]interface{})
	//先查第零个
	for i, i2 := range toArr {
		objs := RowsToStructSlice(rows,i2)
		allDataMap[i+1] = objs
		if !rows.NextResultSet() {
			break
		}
	}
	return allDataMap
}

//获取db结果集 的列与类型
func GetRowsColsAndTypes(rows *sql.Rows)(column_names []string,column_types []*sql.ColumnType){
	var erro error
 	column_names, erro = rows.Columns()
	if erro != nil {
		xlog.Error("GetRowsColsAndTypes获取数据库列错误 %v",erro)
		return nil,nil
	}
	column_types,erro = rows.ColumnTypes()
	if erro != nil {
		xlog.Error("GetRowsColsAndTypes获取数据库类型错误 %v",erro)
		return  nil,nil
	}
	// 打印数据库类型
	// for _,v := range column_types{
	// 	xlog.Debug("db =%v,scan =%v",v.DatabaseTypeName() ,v.ScanType())
	// }
	return
}

//获取结构体地址为rows.Scan提供容器
func GetStructScans(valElem *reflect.Value,column_names []string,column_types []*sql.ColumnType)[]interface{}{
	if valElem == nil {
		xlog.Debug("valElem = %v",valElem)
		return nil
	}
	scan_dest := []interface{}{} //扫描赋值的
	colLen := len(column_names)
	//结构体字段的地址
	for i := 0; i < colLen; i++ {
		dbColName := strings.Title(column_names[i]) //数据库字段生成结构体首字母大写
		one_value := valElem.FieldByName(dbColName) //用名字去找结构体字段可以更好的匹配
		if !DBTypeMatchFieldType(column_types[i].ScanType().String(),one_value.Kind().String()) {
			xlog.Error("字段 = %v 数据库字段类型 %v 与结构体字段类型 %v不匹配 ",dbColName,column_types[i].ScanType().String(),one_value.Kind().String())
			// continue
			//return nil //这里没有匹配上字段就直接返回了,因为这里出错scan也会出错的
		}
		//将结构体的地址赋值给
		scan_dest = append(scan_dest,  one_value.Addr().Interface())
	}
	return  scan_dest
}

//数据库结果字段扫描类型与结构类型进行匹配
func DBTypeMatchFieldType(dbscantp,fieldTp string) bool{
	switch dbscantp {
	case "sql.RawBytes","mysql.NullTime": //fix by zjy 20200826 这里增加ptr 判断
		if fieldTp != "string" && fieldTp != "ptr" {
			return false
		}
	case "uint64","int64","uint32","int32","uint16","int16","uint8","int8","float64","float32":
		if dbscantp != fieldTp && fieldTp != "ptr" {
			return false
		}
	default:
		xlog.Debug("DBTypeMatchFieldType 未处理类型%v", dbscantp)
		return false
	}
	return true
}


