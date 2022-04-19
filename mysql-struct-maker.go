/*
创建时间: 2020/5/1
作者: zjy
功能介绍:
 根据数据库表生成gofile
*/

package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"reflect"
	"strings"
	"sync"
)
//TABLE_NAME
//COLUMN_NAME,COLUMN_TYPE,COLUMN_DEFAULT, IS_NULLABLE,CHARACTER_MAXIMUM_LENGTH,COLUMN_KEY,COLUMN_COMMENT
// 列信息

var (
	sqlDB      *sql.DB
	gameDBConf *csvdata.DbCfg
	dbwg       sync.WaitGroup
)



//func Report() {
//	queryResult := QueryTables()
//	if queryResult == nil {
//		return
//	}
//	for tableName, columnInfos := range queryResult {
//		dbwg.Add(1)
//		go ParseColumn(tableName,columnInfos)
//	}
//}

//func ParseColumn(tableName string, columnInfos []*ColumnInfo)  {
//	defer  dbwg.Done()
//	// 创建csv文件
//	fs, err := os.Create(filepath.Join("./dbmodels",strings.ToLower(tableName) +"_dbfeild.go"))
//	if xutil.IsError(err) {
//		return
//	}
//	defer fs.Close()
//	fs.WriteString(fmt.Sprintf("//生成的文件建议不要改动,详见mysql-struct-maker.go ParseColumn方法源码生成格式 \n"))
//	fs.WriteString(fmt.Sprintf("package dbmodels \n"))
//	fs.WriteString(fmt.Sprintf("\ntype %s struct {\n", xutil.Capitalize(tableName)))
//	for _, info := range columnInfos {
//		// translate to go struct foramt
//		vname := strings.Title(info.columName) // 字段名称
//		retype := DBTypeToGoT(info.columType)  //
//		if strings.Compare(info.isNullable,"YES") == 0{
//			//fix by zjy 20200826
//			//如果字段可以为null就需要给字段设置为指针,解决数据库查询Scan给结构体地址赋值错误问题
//			//sql: Scan error on column index x, name “x”: converting NULL to int64 is unsupported
//			//retype = fmt.Sprintf("*%s",retype) // 字段名称
//			retype = fmt.Sprintf("%s",retype) // 字段名称 还是常规使用,数据库不要有null值
//		}
//
//		fs.WriteString(fmt.Sprintf("\t%s %s `sql:\"%s\"` // 数据库注释:%s \n ", vname, retype,info.columName, info.columComment))
//	}
//	fs.WriteString("}\n")
//
//}

func DBTypeToGoT(dbtype string) string {
	if dbtype == "" || strings.Compare(dbtype, "") == 0 {
		return ""
	}
	resstr := dbtype
	if strings.Contains(dbtype, "varchar") ||
		strings.Contains(dbtype, "longtext")  {
		resstr = "string"
	} else if strings.Contains(dbtype, "date") ||
		      strings.Contains(dbtype, "datetime") ||
		       strings.Contains(dbtype,"timestamp") {
		// resstr = "time.Time" //这里要string类型
		resstr = "string"
	} else if strings.Contains(dbtype, "tinyint") {
		resstr = "int8"
	} else if strings.Contains(dbtype, "smallint") {
		resstr = "int16"
	} else if strings.Contains(dbtype, "integer") {
		resstr = "int32"
	} else if strings.Contains(dbtype, "bigint") {
		resstr = "int64"
	} else if strings.Contains(dbtype, "int") {
		resstr = "int32"
	} else if strings.Contains(dbtype, "double") {
		resstr = "float64"
	} else if strings.Contains(dbtype, "float") {
		resstr = "float32"
	}
	
	// 查看是否是无符号类型
	if strings.Contains(dbtype, "unsigned") {
		resstr = "u" + resstr
	}
	return resstr
}

//需要注册
type Account struct {
	ID  *dbsys.DBInt8
}

type Test interface {
	Name()
}

type Base struct {
	name string
}

func (this *Base)Name()   {
	fmt.Println("Base",this.name)
}

type Sub struct {
	Base
}

//func (this *Sub)Name()   {
//	fmt.Println("Sub",this.name)
//}

func TestPb(){
	//Warp := new(protomsg.MsgWarp)
	//info := new(protomsg.ServerInfoPb)
	//info.AppId = 100
	//info.AppName ="小明"
	//Account,erro := proto.Marshal(info)
	//if erro != nil {
	//	fmt.Println(erro)
	//}
	//Warp.PbMsg = Account
	//
	//newWarpData,_ :=  proto.Marshal(Warp)
	//
	//Parse := new(protomsg.MsgWarp)
	//proto.Unmarshal(newWarpData,Parse)
	//
	//fanWarp := new(protomsg.ServerInfoPb)
	//erro = proto.Unmarshal(Parse.PbMsg,fanWarp)
	//if erro != nil {
	//	fmt.Println(erro)
	//}
	//fmt.Println(fanWarp)
}

// 遍历时删除所有的偶数,结果:确实删除了所有的偶数
func fun1() {
	x := sync.Map{}
	// 构建
	for i := 0; i < 10; i++ {
		x.Store(i, i)
	}
	// 遍历时删除偶数
	x.Range(func(k, v interface{}) bool {
		if k.(int)%2 == 0 {
			x.Delete(k)
		}
		return true
	})
	// 遍历打印剩下的
	cout := 0
	x.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		cout++
		return true
	})
	// 会发现是50个,说明删除了所有的偶数
	fmt.Println("fun 1 删除偶数后,剩余元素数,cout:", cout)
}

// 遍历时删除所有元素,结果:确实删除了所有的元素
func fun2() {
	x := sync.Map{}
	// 构建
	for i := 0; i < 100; i++ {
		x.Store(i, i)
	}
	// 遍历时删除偶数
	x.Range(func(k, v interface{}) bool {
		x.Delete(k)
		return true
	})
	// 遍历打印剩下的
	cout := 0
	x.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		cout++
		return true
	})
	// 会发现是0个,说明删除了所有的元素
	fmt.Println("fun 2全部删除后,剩余元素数,cout:", cout)
}

func main() {

	//// set the file path that result save in
	//
	////fmt.Println("Prase Scuess!")
	//dbsys.AsyncOpenRedis("zjy1")

	//dbsys.RegisterTable(dbmodels.AccountT{})
	//dbsys.RegisterTable(dbmodels.RoleT{})
	//dbsys.RegisterTable(dbmodels.ItemT{})
	//dbsys.RegisterTable(dbmodels.MoneyT{})
	//dbsys.SyncBD("accountdb")
	//dBConf := csvdata.GetDbCfgPtr("accountdb")
	////切换新建的库以便创建表
	//pDB, err := sql.Open("mysql", dbsys.GetMysqlDataSourceName(dBConf))
	//if err != nil {
	//	fmt.Println("open  DB  %v", err)
	//	return
	//}
	//Account := dbmodels.NewAccountT()
	//Account.LoginName.val = "aaa"
	//TestSql(pDB,Account)
	//
	//sizelen := 	binary.Size(Account)
	//pDB.Close()
	//fmt.Println(sizelen,Account)

	//t := reflect.TypeOf(dbmodels.RoleT{})
	//newT := reflect.New(t)
	//obj,ok := newT.Interface().(dbsys.DBObJer)
	//if ok {
	//	fmt.Println(obj)
	//}
	dbmodels.NewItemT()
	dbmodels.NewMoneyT()
	var a = make(map[int32]*dbmodels.MoneyT)
	temp := new(dbmodels.MoneyT)
	var temcccc interface{}
	temcccc = a
	mapVal := reflect.ValueOf(temcccc)
	maps := reflect.MakeMap(mapVal.Type())
	maps.SetMapIndex(reflect.ValueOf(int32(100)),reflect.ValueOf(temp))
	keys := mapVal.MapKeys()
	for i, key := range keys {
		fmt.Println(i,key)
	}
	fmt.Println(a)
	//tx,erro := sqlDB.Begin()
	//if erro != nil {
	//	fmt.Println(erro)
	//	return
	//}
	//
	//_,erro = sqlDB.Exec(createSql)
	//if erro != nil {
	//	fmt.Println(erro)
	//	tx.Rollback()
	//	return
	//}
	//tx.Commit()
	//sqlDB.Close()
}

func TestSql(pDB *sql.DB,objer dbsys.DBObJer)  {
	sql := objer.SelectSql()
	rows,erro := pDB.Query(sql)
	defer rows.Close()
	if erro != nil {
		return
	}
	for rows.Next() {
		val := objer.GetDBValesAddr()
		rows.Scan(val...)
	}
}

func TestReflect()  map[string][]*dbsys.ColumnInfo {
	runD := make(map[string][]*dbsys.ColumnInfo)
	atype := reflect.TypeOf(dbmodels.AccountT{})
	tabName := atype.Name()
	fieldLen := atype.NumField()
	for i := 0; i < fieldLen; i++ {
		column := new(dbsys.ColumnInfo)
		filed := atype.Field(i)
		temType := filed.Type //字段类型
		//字段描述
		dbcol := strings.TrimSpace(filed.Tag.Get("col"))
		tp := strings.TrimSpace(filed.Tag.Get("tp"))
		idx := strings.TrimSpace(filed.Tag.Get("idx"))
		node := strings.TrimSpace(filed.Tag.Get("node"))
		fmt.Println("tag ","tp = ",tp,"idx =",idx,"node",node)
		column.TABLE_NAME = tabName
		if dbcol == "" {
			column.COLUMN_NAME = filed.Name
		}else {
			column.COLUMN_NAME = dbcol
		}
		if temType.Kind() == reflect.Ptr {
			temType = temType.Elem()
		}

		if tp == "" { //取字段类型
			column.COLUMN_TYPE = dbsys.RunTypeToDBType(temType.Name())
		}else {
			column.COLUMN_TYPE = tp
		}
		//不为nil
		column.IS_NULLABLE = "NO"
		////字符串才处理长度
		//if strings.Contains(tp,"varchar") {
		//	tp = strings.Replace(tp,"varchar(","",-1)
		//	tp = strings.Replace(tp,")","",-1)
		//	column.CHARACTER_MAXIMUM_LENGTH = strutil.StrToUint64(tp)
		//}
		//column.COLUMN_KEY = idx
		column.COLUMN_COMMENT.String = node
		tabInfo, ok := runD[tabName]
		if !ok {
			var temslice []*dbsys.ColumnInfo
			temslice = append(temslice, column)
			runD[column.TABLE_NAME] = temslice
			continue
		}

		runD[column.TABLE_NAME] = append(tabInfo, column)
	}


	return runD
	//pacc := new(dbmodels.Account)
	//pacc.AccID = new(dbsys.DBInt8)
	//pacc.AccID.SetVal(100)
	//fmt.Println(pacc.AccID)
	//pacc.AccID.AddVal(100)
	//fmt.Println(pacc.AccID)
	//val := reflect.ValueOf(pacc)
	//valElm := val.Elem()
	//tpElm := valElm.Type()
	//fieldLen = valElm.NumField()
	//for i := 0; i < fieldLen; i++ {
	//	valfiled := valElm.Field(i)
	//	tp := tpElm.Field(i)
	//	fmt.Println("--------",tp.Name,tp.Tag)
	//	//不是指针
	//	if valfiled.Kind() == reflect.Ptr  {
	//		valfiled = valfiled.Elem()
	//	}
	//	if valfiled.Kind() != reflect.Struct {
	//		return
	//	}
	//	inType := valfiled.Type()
	//	fieldTypeLen := valfiled.NumField()
	//	for j := 0; j < fieldTypeLen; j++ {
	//		tFiled := valfiled.Field(j)
	//		inFieldType := inType.Field(j)
	//		fmt.Println(inFieldType.Name,inFieldType.Type)
	//		fmt.Println(tFiled.Int())
	//	}
	//}
}


func CheckHasDB(dbName string) bool {
	rows, erro := sqlDB.Query("SHOW DATABASES;")
	if erro != nil {
		fmt.Println(erro)
		return false
	}
	dbs := dbsys.RowsToStrArrSlice(rows)
	if dbs == nil {
		return  false
	}
	for _, db := range dbs {
		for _, s := range db {
			if dbName == s {
				return  true
			}
		}
	}
	return false
}

func CreateDB(cfg *csvdata.DbCfg) error  {
	if cfg == nil {
		return  errors.New( "  数据库配置为null")
	}
	charOrder := cfg.Char_set + "_bin"
	_,erro := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET %s  COLLATE %s;",cfg.Dbname,gameDBConf.Char_set,charOrder))
	if erro != nil {
		return erro
	}
	return nil
}



