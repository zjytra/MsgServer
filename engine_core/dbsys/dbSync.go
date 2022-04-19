/*
创建时间: 2021/4/5 18:04
作者: zjy
功能介绍:
数据库同步
*/

package dbsys

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
	"strings"
)


const (
	//主键
	PRIKey  = "PRIMARY"
	//唯一索引
	UNIQUE  = "UNIQUE"
	//普通索引数据库查询出来的是小写
	MUl = "MUl"

	//数据类型
	DBUnSigned = "unsigned"
	DBInt8Type = "tinyint"
	DBUint8Type = DBInt8Type + " " + DBUnSigned
	DBInt16Type = "smallint"
	DBUint16Type = DBInt16Type + " " + DBUnSigned
	DBInt32Type = "int"
	DBUint32Type = DBInt32Type + " " + DBUnSigned
	DBInt64Type = "bigint"
	DBUint64Type = DBInt64Type + " " + DBUnSigned
	DBFloat = "float(11, 5)"
	DBVarchar = "varchar"
	DBDATETIME = "datetime"
	WhereFlagSelect = 1
	WhereFlagUpdate = 2
	auto_increment = "auto_increment"
)



func RunTypeToDBType(runTypeName string) string  {
	switch runTypeName {
	case "DBInt8":
		return DBInt8Type
	case "DBUInt8":
		return DBUint8Type
	case "DBInt16":
		return DBInt16Type
	case "DBUInt16":
		return DBUint16Type
	case "DBInt32":
		return DBInt32Type
	case "DBUInt32":
		return DBUint32Type
	case "DBInt64":
		return DBInt64Type
	case "DBUInt64":
		return DBUint64Type
	case "DBStr":
		return "varchar(30)"
	case "DBBool":
		return DBInt8Type
	case "DBFloat32":
		return DBFloat
	case "DBIntMap": //当字符串处理
		return "varchar(2048)"
	}
	return runTypeName
}

func IsDBStr(dbtypeName string) bool  {
	return strings.Contains(dbtypeName, "char")||
		strings.Contains(dbtypeName, "varchar")||
		strings.Contains(dbtypeName, "longtext")||
		strings.Contains(dbtypeName, "text")
	//strings.Contains(dbtypeName, "binary")||
	//strings.Contains(dbtypeName, "blob")
}

func GetDefaultVal(dbtypeName string) string  {
	switch dbtypeName {
	case DBInt8Type,DBUint8Type,DBInt16Type,DBUint16Type,DBInt32Type,DBUint32Type,DBInt64Type,DBUint64Type:
		return "0"
	case DBFloat:
		return "0.00000"
	default:
		if IsDBStr(dbtypeName) {
			return " '' "
		}
	}
	return ""
}


func (this *MySqlDBStore) GetDBTable(tabName string) *DBTable {
	tab, ok := this.dbTables[tabName]
	if ok {
		return tab
	}
	xlog.Warning("not find table %s  ", tabName)
	return nil
}
func (this *MySqlDBStore) addTableIndex(colInfo *ColumnInfo) {
	idexs, ok := this.tableIdxs[colInfo.TABLE_NAME]
	if ok {
		idexs = append(idexs, colInfo)
		this.tableIdxs[colInfo.TABLE_NAME] = idexs
		return
	}
	var idex []*ColumnInfo
	idex = append(idex, colInfo)
	this.tableIdxs[colInfo.TABLE_NAME] = idex
}


//支持中间删除添加字段,不支持修改表字段名,对象被删除的字段,数据库会删除哦
//创建一个表
func (this *MySqlDBStore) RegisterTable(dbmodel interface{}) *DBTable {
	tb := this.RegisterTableByName(dbmodel,"")
	return tb
}


//支持中间删除添加字段,不支持修改表字段名,对象被删除的字段,数据库会删除哦
//争对动态创建表
func (this *MySqlDBStore) RegisterTableByName(dbmodel interface{}, tbName string) *DBTable {
	runType := reflect.TypeOf(dbmodel)
	if runType.Kind() != reflect.Struct {
		xlog.Debug("%s 数据类型 不是结构体", runType.Name())
		return nil
	}
	runObjVal := reflect.New(runType)
	dbObj,ok := runObjVal.Interface().(DBObJer)
	if !ok {
		xlog.Error("%s 不是数据库类型 DBObJer ", runType.Name())
		return nil
	}
	tb := NewDBTable(dbObj)
	if tbName != "" {
		tb.TbName = tbName
	}else {
		tb.TbName = runType.Name()
	}
	_, isRegister := this.dbTables[tb.TbName]
	if isRegister {
		xlog.Debug("%s 表已经注册", tb.TbName)
		return nil
	}
	tb.dbStore = this
	tb.ParesFields(runType)

	//保存表的信息
	this.dbTables[tb.TbName] = tb
	//放到内存中在注册对象
	return tb
}



//sql 是否超过长度
func sqlLenIsMax(oldSql *string,addSql *string) bool {
	sqlby := []byte(*oldSql + *addSql)

	//test
	//if len(sqlby) > 1000 {
	//	return true
	//}

	////最长1Mb
	if len(sqlby) > 1024000 {
		return true
	}
	return false
}

//同步数据字段
func (this *MySqlDBStore) DoSyncDiffSql(dbtbls map[string]map[string]*ColumnInfo) {
	//比较表sql
	Sql := ""
	//检测sql长度
	for s, regTb := range this.dbTables {
		dbfields, ok := dbtbls[s]
		if !ok { //数据库没有需要创建表
			addSql := regTb.GetCreateTableSql()
			//最长1Mb
			if sqlLenIsMax(&Sql,&addSql) {
				if Sql != "" {
					_, err := this.db.Exec(Sql)
					xlog.Debug("sqlLenIsMax  GetCreateTableSql %s", Sql)
					Sql = ""
					if err != nil {
						Sql += addSql
						xlog.Error("%v", err)
						continue
					}
				}
			}
			if addSql != "" {
				xlog.Debug("CreateTableSql %s", addSql)
			}
			Sql += addSql
			continue
		}
		//存在表 匹配字段
		addSql := regTb.GetModifyTableSql(dbfields)
		if sqlLenIsMax(&Sql,&addSql) {
			if Sql != "" {
				_, err := this.db.Exec(Sql)
				xlog.Debug("sqlLenIsMax GetModifyTableSql sql %s", Sql)
				Sql = ""
				if err != nil {
					Sql += addSql
					xlog.Debug("%v", err)
					continue
				}
			}
		}
		if addSql != "" {
			xlog.Debug("ModifyTableSql %s", addSql)
		}
		Sql += addSql
	}

	if Sql != "" {
		_, err := this.db.Exec(Sql)
		modiArr :=  strings.Split(Sql,";")
		if modiArr != nil && len(modiArr) > 0 {
			for _, s := range modiArr {
				xlog.Debug("db sync sql %s", s)
			}
		}else {
			xlog.Debug("db sync sql %s", Sql)
		}
		Sql = ""
		if err != nil {
			xlog.Error("%v", err)
		}
	}
}
//同步数据库
//在注册完表后执行
func (this *MySqlDBStore)SyncBD()  {
	//先连接mysql的库
	mysqlDBCof := csvdata.GetDbCfgPtr(MysqlDBName)//"mysql"
	if mysqlDBCof == nil {
		panic("mysqlDB == nil")
	}
	sqlDB, err := sql.Open("mysql", GetMysqlDataSourceName(mysqlDBCof))
	if err != nil {
		xlog.Debug("open  DB  %v", err)
		return
	}
	dBConf := csvdata.GetDbCfgPtr(this.dbName)
	if dBConf == nil {
		panic(this.dbName + " conf == nil")
	}
	if !CheckHasDB(sqlDB,this.dbName) {
		createErro := DoCreateDB(sqlDB,dBConf)
		if createErro != nil {
			xlog.Debug( "create db %s err: %v",this.dbName,createErro)
		}else {
			xlog.Debug("DoCreateDB %v  success",this.dbName)
		}
	}
	sqlDB.Close()
	//切换到对应的库
	err = this.OpenDB()
	if err != nil {
		xlog.Debug("open  DB  %v", err)
		return
	}
	//查询数据表的信息
	dbtbls := this.QueryTables(this.dbName)
	for s, _ := range dbtbls {
		xlog.Debug("db has tab %s",s)
	}
	//数据库表信息与内存注册的表信息对比需要执行的sql
	this.DoSyncDiffSql(dbtbls)
	xlog.Debug(" syncDB dbname %v", this.dbName)
}


func (this *MySqlDBStore)QueryTables(dbName string) map[string]map[string]*ColumnInfo {
	//组装好索引，再填充字段
	indexs := this.queryIndex(dbName)
	queryResult := make(map[string]map[string]*ColumnInfo)
	//这里还是增加排序,虽然反射是通过字段名称读取的,为了方便查看字段位置，还是匹配下顺序为好
	sql := `SELECT TABLE_NAME,COLUMN_NAME,COLUMN_TYPE,COLUMN_DEFAULT,IS_NULLABLE,COLUMN_COMMENT,ORDINAL_POSITION,EXTRA FROM information_schema.COLUMNS   
WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME,ORDINAL_POSITION;`
	rows, err := this.db.Query(sql,dbName)
	defer rows.Close()
	if err != nil {
		panic(err)
		return nil
	}

	for rows.Next() {
		column := new(ColumnInfo)
		erro := rows.Scan(&column.TABLE_NAME,
			&column.COLUMN_NAME,
			&column.COLUMN_TYPE,
			&column.COLUMN_DEFAULT,
			&column.IS_NULLABLE,
			&column.COLUMN_COMMENT,
			&column.FiledOrder,
			&column.EXTRA)
		if erro != nil {
			xlog.Debug("查询数据库  表%s 字段名%s 错误 %v",column.TABLE_NAME,column.COLUMN_NAME,erro)
		}
		if  IsDBStr(column.COLUMN_TYPE) && column.COLUMN_DEFAULT.String == "" {
			column.COLUMN_DEFAULT.String = " '' "
		}
		column.SetDBColIndex(indexs)
		//找表
		tabInfo, ok := queryResult[column.TABLE_NAME]
		if !ok { //没找到表创建表的字段映射
			tabInfo = make(map[string]*ColumnInfo)
			tabInfo[column.COLUMN_NAME] = column
			queryResult[column.TABLE_NAME] = tabInfo
			continue
		}//找到了直接赋值
		tabInfo[column.COLUMN_NAME] = column
	}
	return queryResult
}

func (this *MySqlDBStore)queryIndex(dbName string) map[string]map[string][]*ColumnIndex {
	//这里还是增加排序,虽然反射是通过字段名称读取的,为了方便查看字段位置，还是匹配下顺序为好
	sql := `SELECT TABLE_NAME,NON_UNIQUE,INDEX_NAME,SEQ_IN_INDEX,COLUMN_NAME FROM information_schema.statistics WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME,SEQ_IN_INDEX;`
	rows, err := this.db.Query(sql,dbName)
	defer rows.Close()
	if err != nil {
		panic(err)
		return nil
	}
	//记录表的联合索引
	queryResult := make(map[string]map[string][]*ColumnIndex)
	for rows.Next() {
		column := new(ColumnIndex)
		erro := rows.Scan(
			&column.TABLE_NAME,
			&column.NON_UNIQUE,
			&column.INDEX_NAME,
			&column.SEQ_IN_INDEX,
			&column.COLUMN_NAME)
		if erro != nil {
			xlog.Debug("查询数据库索引  表%s 字段名%s 错误 %v",column.TABLE_NAME,column.COLUMN_NAME,erro)
		}
		//构建成一样的索引类型名称
		if  column.INDEX_NAME == PRIKey && column.NON_UNIQUE == 0 { //主键
			column.INDEX_Type_NAME = PRIKey
		}else if column.INDEX_NAME != PRIKey && column.NON_UNIQUE == 0 { //唯一索引
			column.INDEX_Type_NAME = UNIQUE
		}else {
			column.INDEX_Type_NAME = MUl
		}
		SetMoreIdx(queryResult,column)
		//queryResult[column.TABLE_NAME] = colIndex
	}
	//看一下修改指针是否起效果
	for _, indexs := range queryResult {
		for indexName, col := range indexs {
			var collen  = len(col)
			if collen == 1 {
				continue
			}
			//拼接字段名称
			cols := col[0].COLUMN_NAME
			//联合索引处理
			for i := 1; i < collen; i++ {
				column := col[i]
				//构建成一样的索引类型名称
				cols += "," + column.COLUMN_NAME
			}
			//联合索引只保留一个字段 保留第一个字段,方便匹配
			column := col[0]
			column.MoreINDEX_Cols = cols
			column.More_INDEX_NAME = column.INDEX_NAME
			if column.INDEX_NAME != PRIKey && column.NON_UNIQUE == 0 { //唯一索引
				column.More_Type_NAME = UNIQUE
			}else {
				column.More_Type_NAME = MUl
			}
			//如果这个索引是联合索引,那就不使用单索引了
			column.INDEX_NAME = ""
			column.INDEX_Type_NAME = ""
			//保留第一个字段 保证联合索引只有一个字段
			//设置回去
			indexs[indexName] = col[0:1]
		}
	}
	return  queryResult
}

func SetMoreIdx(moreIdx map[string]map[string][]*ColumnIndex,column *ColumnIndex){
	//没有找到表
	tabInfo, ok := moreIdx[column.TABLE_NAME]
	if !ok {
		tabInfo = make(map[string][]*ColumnIndex)
		var temslice []*ColumnIndex
		temslice = append(temslice, column)
		tabInfo[column.INDEX_NAME] = temslice
		moreIdx[column.TABLE_NAME] = tabInfo
		return
	}
	colIndex,idxok := tabInfo[column.INDEX_NAME]
	if !idxok {
		var temslice []*ColumnIndex
		temslice = append(temslice, column)
		tabInfo[column.INDEX_NAME] = temslice
		return
	}
	//找到索引的map
	colIndex = append(colIndex,column)
	tabInfo[column.INDEX_NAME] = colIndex
}

func CheckHasDB(pdb *sql.DB,dbName string) bool {
	rows, erro := pdb.Query("select SCHEMA_NAME from information_schema.SCHEMATA where SCHEMA_NAME = ?; ",dbName)
	if erro != nil {
		xlog.Error("%v",erro)
		return false
	}
	dbs := RowToMap(rows)
	if dbs == nil {
		return  false
	}
	hasName, ok :=  dbs["SCHEMA_NAME"]
	if ok && hasName == dbName {
		return true
	}
	return false
}

func DoCreateDB(pdb *sql.DB,cfg *csvdata.DbCfg) error  {
	if cfg == nil {
		return  errors.New( "  数据库配置为null")
	}
	charOrder := cfg.Char_set + "_bin"
	_,erro := pdb.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET %s  COLLATE %s;",cfg.Dbname,cfg.Char_set,charOrder))
	if erro != nil {
		xlog.Error("%v",erro)
		return erro
	}
	return nil
}


func FindDBColByOrderId(fields map[string]*ColumnInfo,fieldOrder uint8)*ColumnInfo  {
	for _, field := range fields {
		if field.FiledOrder == fieldOrder {
			return field
		}
	}
	return nil
}


