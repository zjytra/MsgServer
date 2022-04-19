/*
创建时间: 2021/3/22 21:19
作者: zjy
功能介绍:
1.同步表的功能
2.存储变更的数据
数据库表映射
需要把所有的表信息存入
包含表的字段

*/

package dbsys

import (
	"fmt"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
	"strings"
	"sync"
)



//注册表的
type DBTable struct {
	TbName string //表的名称
	//与结构体映射类型
	SynFields []*ColumnInfo          //与数据库同步的字段
	FieldsMap map[string]*ColumnInfo //与数据库同步的字段
	DelUpdateWhere map[string]*ColumnInfo
	SelectWhere    map[string]*ColumnInfo
	mulIndex  map[string][]string    //联合索引
	dbStore *MySqlDBStore   //数据库对象
	subTables []*DBTable //子表
	subTabLock sync.Mutex
	mainTable *DBTable //当作为子表时指向主表
	createFun CreateObjFun   //映射的对象主要方便数据库线程创建对象
}

func NewDBTable(Obj DBObJer)*DBTable {
	tb := new(DBTable)
	tb.createFun = Obj.CreateObj
	tb.DelUpdateWhere = make(map[string]*ColumnInfo)
	tb.SelectWhere = make(map[string]*ColumnInfo)
	return tb
}


func (this *DBTable) GetCreateTableSql() string {
	if this.SynFields == nil {
		return ""
	}
	createSql := "create table " + this.TbName + " ("
	//遍了字段
	flen := len(this.SynFields)
	for i, info := range this.SynFields {
		createSql += info.GetColumnSql()
		if i < flen-1 {
			createSql += ","
		}
	}
	regIndex, regOk := this.dbStore.tableIdxs[this.TbName]
	if regOk {
		//遍历索引
		for _, info := range regIndex {
			if info.INDEX_Type_NAME != "" {
				createSql += "," + info.GetColumnIndexSql()
			}
			if info.More_INDEX_NAME != "" {
				createSql += "," + info.GetColumnMoreIndexSql()
			}
		}
	}
	createSql += ")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;"
	return createSql
}


//支持中间删除添加字段,不支持修改表字段
func (this *DBTable) GetModifyTableSql(dbfields map[string]*ColumnInfo) string {
	if this.SynFields == nil {
		return ""
	}
	mdSql := ""
	afterCol := ""
	for i, info := range this.SynFields {
		DBCol, findok := dbfields[info.COLUMN_NAME]
		if !findok { //没有找到
			if  i > 0 { //添加到对应字段后面
				afterCol = this.SynFields[i-1].COLUMN_NAME
			}
			mdSql += info.GetADDColumnSql(afterCol) + info.ADD_AUTO_INCR()
		} else {
			//修改
			mdSql += info.GetFieldChangeSql(DBCol)
		}
	}
	//删除多余的字段
	for name,dbCol := range dbfields{
		//注册的是否有对应的数据库字段，没有就删除
		if !this.registerHasDBField(name) {
			mdSql += dbCol.dropColSql()
		}
	}
	//ALTER TABLE `gamedb`.`RoleT`
	//DROP COLUMN `Job`

	return mdSql
}


func (this *DBTable) registerHasDBField(name string) bool {
	for _, regInfo := range this.SynFields {
		if name == regInfo.COLUMN_NAME {
			return true
		}
	}
	return false
}

//解析字段
func (this *DBTable) ParesFields(RunType reflect.Type) {
	this.FieldsMap = make(map[string]*ColumnInfo)
	fieldLen := RunType.NumField()
	//从 1 开始
	var j uint8 = 1
	for i := 0; i < fieldLen; i++ {
		filed := RunType.Field(i)
		temType := filed.Type //字段类型
		if filed.Name == "DBObj" {
			continue
		}
		//有忽略字段
		_, hasIgnore := filed.Tag.Lookup("ign")
		if hasIgnore {
			continue
		}
		column := new(ColumnInfo)
		column.FiledOrder = j
		j++
		//字段描述
		dbcol := strings.TrimSpace(filed.Tag.Get("col"))
		tp := strings.TrimSpace(filed.Tag.Get("tp"))           //字段类型
		idx := strings.TrimSpace(filed.Tag.Get("indexTP"))       //索引类型
		mIdx := strings.TrimSpace(filed.Tag.Get("moreIdxCol"))    //联合索引字段,在联合索引第一个字段设置,不能设置到后面的字段去
		mIdxType := strings.TrimSpace(filed.Tag.Get("moreIdxTp")) //联合索引类型
		node := strings.TrimSpace(filed.Tag.Get("node"))
		whereFlagTag := strings.TrimSpace(filed.Tag.Get("QFlag"))
		column.WhereFlag = strutil.StrToUint8(whereFlagTag)
		//fmt.Println("tag ", "tp = ", tp, "idx =", idx, "node", node)
		column.TABLE_NAME = this.TbName
		if dbcol == "" {
			column.COLUMN_NAME = filed.Name
		} else {
			column.COLUMN_NAME = dbcol
		}
		if temType.Kind() == reflect.Ptr {
			temType = temType.Elem()
		}
		if tp == "" { //取字段类型
			column.COLUMN_TYPE = RunTypeToDBType(temType.Name())
		} else {
			column.COLUMN_TYPE = tp
		}
		column.COLUMN_DEFAULT.String = GetDefaultVal(column.COLUMN_TYPE)
		//不为nil
		column.IS_NULLABLE = "NO"
		column.INDEX_Type_NAME = idx
		//主键是否自增 其他 字段自增无效
		if column.INDEX_Type_NAME == PRIKey  {
			//是否由自增标识
			_, hasauto := filed.Tag.Lookup("auto_increment")
			if hasauto {
				column.EXTRA = auto_increment
				//设置了自增的列没有默认值
				column.COLUMN_DEFAULT.String = ""
			}
		}
		column.SetIndexName()

		//多列组合成一个索引
		if mIdx != "" {
			cols := strings.Split(mIdx, ",")
			//换成下划线
			column.More_INDEX_NAME = strutil.MergeStrSlice(cols, "_")
			column.MoreINDEX_Cols = mIdx
			column.More_Type_NAME = mIdxType
		}

		if idx != "" || mIdx != "" {
			this.dbStore.addTableIndex(column)
		}
		column.COLUMN_COMMENT.String = node
		this.SynFields = append(this.SynFields, column)

		//查询字段
		if column.IsQueryWhere() { //select 字段
			this.SelectWhere[column.COLUMN_NAME] = column
		}

		if column.IsUpdateOrDelWhere() { //del 字段
			this.DelUpdateWhere[column.COLUMN_NAME] = column
		}

		//字段不会重名吧,重名玩坏
		_, ok := this.FieldsMap[column.COLUMN_NAME]
		if ok {
			xlog.Error("表 %s 字段重名 %s", this.TbName, column.COLUMN_NAME)
		}
		this.FieldsMap[column.COLUMN_NAME] = column
	}
}

//支持中间删除添加字段,不支持修改表字段名,对象被删除的字段,数据库会删除哦
//注册一个子表
func (this *DBTable) RegisterSubTable(dbmodel interface{}) *DBTable {
	tb := this.RegisterSubTableByName(dbmodel,"")
	return tb
}

//支持中间删除添加字段,不支持修改表字段名,对象被删除的字段,数据库会删除哦
//注册一个子表
func (this *DBTable) RegisterSubTableByName(dbmodel interface{},tbName string) *DBTable {
	if this.mainTable != nil && this.mainTable.mainTable != nil && this.mainTable.mainTable.mainTable != nil {
		xlog.Error("主次表关系最多关联三层")
		return nil
	}

	tb := this.dbStore.RegisterTableByName(dbmodel,tbName)
	this.AddSubTable(tb)
	return tb
}

func (this *DBTable) AddSubTable(subTable *DBTable){
	if subTable == nil {
		return
	}
	for _, table := range this.subTables {
		if table.TbName == subTable.TbName {
			xlog.Error("主表 %s 已经注册过子表 %s", this.TbName,subTable.TbName)
			return
		}
	}
	this.subTables = append(this.subTables,subTable)
	//设置主表的子表
	subTable.mainTable = this
}

func (this *DBTable) HasSubTables() bool{
	if this.subTables == nil {
		return false
	}
	return  len(this.subTables) > 0
}

func (this *DBTable) GetSubTables() []*DBTable{
	return  this.subTables
}


func (this *DBTable) SelectSql(uid int64) string {
	if this.SynFields == nil {
		return ""
	}
	flen := len(this.SynFields)
	if flen == 0 {
		return ""
	}
	if this.SelectWhere == nil  {
		xlog.Error("%s not DelUpdateWhere", this.TbName)
		return ""
	}
	whereLen := len(this.SelectWhere)
	if whereLen == 0 {
		xlog.Error("%s not SelectWhere", this.TbName)
		return ""
	}

	filedSql := ""
	for i, pCol := range this.SynFields {
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.TbName, i)
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
	for _, pCol := range this.SelectWhere {
		if pCol == nil {
			xlog.Error("%s field index %d is nil", this.TbName, i)
			i++
			continue
		}
		whereSql += pCol.COLUMN_NAME + " = " + pCol.GetValStr(uid)
		if i < whereLen - 1 {
			whereSql += ","
		}
		i++
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s;",filedSql, this.TbName,whereSql)
}
