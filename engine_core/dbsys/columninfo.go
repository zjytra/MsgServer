/*
创建时间: 2021/4/20 21:33
作者: zjy
功能介绍:

*/

package dbsys

import (
	"database/sql"
	"fmt"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"strings"
)

type ColumnInfo struct {
	TABLE_NAME     string //表名
	COLUMN_NAME    string //列名
	COLUMN_TYPE    string //列类型
	COLUMN_DEFAULT sql.NullString //默认值
	IS_NULLABLE string    //是否为null
	COLUMN_COMMENT sql.NullString //
	INDEX_NAME string  //索引名称
	INDEX_Type_NAME string  //索引类型名称 primkey
	More_INDEX_NAME string  //联合索引名称
	More_Type_NAME string  //联合索引类型 是UNIQUE 还是 NORMAL
	MoreINDEX_Cols string  //联合索引列
	FiledOrder uint8  //字段顺序
	WhereFlag uint8  //字段作为查询条件的标志2进制位  1 select 2 update与del共用   3表示select与update insert
	EXTRA   string   //主键是否自增
}

//创建字段
func (this *ColumnInfo) GetColumnSql() string  {
	str := " " + this.COLUMN_NAME + " " + this.COLUMN_TYPE
	if IsDBStr(this.COLUMN_TYPE) { //字符串要加上字符集
		str += " CHARACTER SET utf8mb4 COLLATE utf8mb4_bin "
	}
	str += this.NotNull()

	return str
}

//是否改变 1 标识字段名称不一样
//2是其他属性不一样
//0表示无需修改
func (this *ColumnInfo) GetFieldChangeSql(db *ColumnInfo) string  {
	changeSql := ""
	if this.COLUMN_NAME != db.COLUMN_NAME || this.HasAUTO_INCR() != db.HasAUTO_INCR()  {
		changeSql += this.GetChangeColumnSql(db)
		xlog.Debug("数据库比较 %s 表  注册字段 %s  数据库字段 %s 名称不一", this.TABLE_NAME,this.COLUMN_NAME,db.COLUMN_NAME)
	} else if (this.COLUMN_TYPE != db.COLUMN_TYPE ||
				this.IS_NULLABLE != db.IS_NULLABLE ||
				this.COLUMN_DEFAULT.String != db.COLUMN_DEFAULT.String||
				this.COLUMN_COMMENT.String != db.COLUMN_COMMENT.String ||
				this.HasAUTO_INCR() != db.HasAUTO_INCR()) {
		changeSql += this.GetModifyColumnSql()
		xlog.Debug("数据库比较 %s 表  注册字段 %s  数据库字段 %s  不能匹配", this.TABLE_NAME,this.COLUMN_NAME,db.COLUMN_NAME)
	}

	if (this.INDEX_NAME != db.INDEX_NAME || this.INDEX_Type_NAME != db.INDEX_Type_NAME) ||
		(this.More_INDEX_NAME != db.More_INDEX_NAME || this.More_Type_NAME != db.More_Type_NAME){
		xlog.Debug("数据库比较 %s 表  注册索引 %s  数据库 %s 索引", this.TABLE_NAME,this.INDEX_NAME,db.INDEX_NAME)
		changeSql += this.GetColumnAlterIndexSql(db)
	}
	return changeSql
}

func (this *ColumnInfo) NotNull() string {
	str := "  NOT NULL "
	//自增
	str += this.ADD_AUTO_INCR()

	if !this.HasAUTO_INCR() &&  this.COLUMN_DEFAULT.String != "" {
		str +=  " DEFAULT " + GetDefaultVal(this.COLUMN_TYPE)
	}
	if this.COLUMN_COMMENT.String != "" {
		str += " COMMENT '" + this.COLUMN_COMMENT.String + "' "
	}
	return str
}


//自增sql
func (this *ColumnInfo) HasAUTO_INCR() bool {
	//主键才有
	if this.INDEX_Type_NAME != PRIKey {
		return false
	}
	if this.EXTRA == auto_increment {
		return true
	}
	return false
}
//自增sql
func (this *ColumnInfo) ADD_AUTO_INCR() string {
	if this.HasAUTO_INCR() {
		return " AUTO_INCREMENT "
	}
	return ""
}

//是否存在索引
func (this *ColumnInfo)HasIndex()bool  {
	return  this.INDEX_Type_NAME != "" || this.More_INDEX_NAME != ""
}

func (this *ColumnInfo)GetColumnIndexSql()  string {
	if this.INDEX_Type_NAME == PRIKey {
		return " PRIMARY KEY(" + this.COLUMN_NAME + ")"
	}
	if this.INDEX_Type_NAME == UNIQUE {
		return " UNIQUE INDEX "+this.INDEX_NAME+"(" + this.COLUMN_NAME + ")"
	}
	if this.INDEX_Type_NAME != "" {
		return " INDEX "+this.INDEX_NAME +"(" + this.COLUMN_NAME + ")"
	}
	return ""
}

func (this *ColumnInfo) SetIndexName()   {
	if this.INDEX_Type_NAME == PRIKey {
		this.INDEX_NAME = "PRIMARY"
	}else if this.INDEX_Type_NAME == UNIQUE {
		this.INDEX_NAME =  "uq_" + this.COLUMN_NAME
	}else if this.INDEX_Type_NAME != "" {
		this.INDEX_NAME =  "idx_"+this.COLUMN_NAME
	}
}

func (this *ColumnInfo)GetColumnMoreIndexSql()  string {
	if this.More_Type_NAME == UNIQUE {
		return " UNIQUE INDEX "+this.More_INDEX_NAME+"(" + this.MoreINDEX_Cols + ")"
	}
	if this.More_Type_NAME != "" {
		return " INDEX " + this.More_INDEX_NAME +"(" + this.MoreINDEX_Cols + ")"
	}
	return ""
}

//修改字段
func (this *ColumnInfo)GetADDColumnSql(after string) string  {
	sql := " ALTER TABLE " + this.TABLE_NAME
	sql += " ADD COLUMN "
	sql += this.GetColumnSql()
	if after != "" {
		sql +=  " AFTER " + after
	}
	sql += ";"

	if this.More_Type_NAME != "" || this.INDEX_NAME != "" {
		sql += "ALTER TABLE " + this.TABLE_NAME + " ADD " + this.GetColumnIndexSql() + ";"
	}
	return sql
}


//删除字段
func (this *ColumnInfo)dropColSql() string  {
	sql := " ALTER TABLE " + this.TABLE_NAME
	sql += " DROP COLUMN "
	sql += this.COLUMN_NAME + ";"
	return sql
}

//修改字段
func (this *ColumnInfo)GetChangeColumnSql(db *ColumnInfo) string  {
	sql := " ALTER TABLE " + this.TABLE_NAME
	sql += " CHANGE COLUMN  " + db.COLUMN_NAME + " "
	sql += this.GetColumnSql() + ";"
	return sql
}
//修改字段
func (this *ColumnInfo)GetModifyColumnSql() string  {
	sql := " ALTER TABLE " + this.TABLE_NAME
	sql += " MODIFY COLUMN  "
	sql += this.GetColumnSql() + ";"
	return sql
}


func (this *ColumnInfo)GetColumnAlterIndexSql(db *ColumnInfo) string {
	sql := ""
	if this.INDEX_NAME == "" && db.INDEX_NAME != "" { //自己没得表有要删除索引
		if db.INDEX_Type_NAME == PRIKey { //主键删除
			sql += " ALTER TABLE "+ db.TABLE_NAME + " DROP PRIMARY KEY;"
		}else {
			sql += " DROP INDEX "+ db.INDEX_NAME + " ON "+ this.TABLE_NAME +";"
		}
	}else if db.INDEX_NAME != "" && (this.INDEX_NAME != db.INDEX_NAME || this.INDEX_Type_NAME != db.INDEX_Type_NAME) {
		if db.INDEX_Type_NAME == PRIKey { //主键删除
			sql += " ALTER TABLE "+ db.TABLE_NAME + " DROP PRIMARY KEY;"
		}else {
			sql += " DROP INDEX "+ db.INDEX_NAME + " ON "+ this.TABLE_NAME +";"
		}
		//后面要加索引
		sql += " ALTER TABLE " + this.TABLE_NAME + " add " + this.GetColumnIndexSql() + ";"
	}else if db.INDEX_NAME == "" && this.INDEX_NAME != "" {
		//后面要加索引
		sql += " ALTER TABLE " + this.TABLE_NAME + " add " + this.GetColumnIndexSql() + ";"
	}

	// 联合索引
	if this.More_INDEX_NAME == "" && db.More_INDEX_NAME != "" { //自己没得表有要删除索引
		sql += " DROP INDEX "+ db.More_INDEX_NAME + " ON "+ this.TABLE_NAME +";"
	}else if db.More_INDEX_NAME != "" && (this.More_INDEX_NAME != db.INDEX_NAME || this.More_Type_NAME != db.More_Type_NAME)  {
		sql += " DROP INDEX "+ db.More_INDEX_NAME + " ON "+ this.TABLE_NAME +";"
		//后面要加索引
		sql += " ALTER TABLE " + this.TABLE_NAME + " add " + this.GetColumnMoreIndexSql() + ";"
	}else if db.More_INDEX_NAME == "" && this.More_INDEX_NAME != "" {
		//后面要加索引
		sql += " ALTER TABLE " + this.TABLE_NAME + " add " + this.GetColumnMoreIndexSql() + ";"
	}
	return sql
}

//设置数据库查询出的字段索引
func (this *ColumnInfo)SetDBColIndex(indexs map[string]map[string][] *ColumnIndex){
	tbIdxs,idxOk := indexs[this.TABLE_NAME]
	//有索引
	if idxOk {
		for _, idx := range tbIdxs {
			if idx == nil || len(idx) == 0  {
				continue
			}
			colIdx := idx[0]
			if colIdx.COLUMN_NAME != this.COLUMN_NAME {
				continue
			}
			if colIdx.INDEX_NAME != "" {
				this.INDEX_Type_NAME = colIdx.INDEX_Type_NAME
				this.INDEX_NAME = colIdx.INDEX_NAME
			}
			if colIdx.More_INDEX_NAME != "" {
				this.More_INDEX_NAME  = colIdx.More_INDEX_NAME  //联合索引名称
				this.More_Type_NAME  = colIdx.More_Type_NAME   //联合索引类型 是UNIQUE 还是 NORMAL
				this.MoreINDEX_Cols  = colIdx.MoreINDEX_Cols   //联合索引列
			}
		}
	}
}

func (this *ColumnInfo)GetValStr(val interface{}) string {
	if strings.Contains(this.COLUMN_TYPE,"int") {
		return fmt.Sprintf("%d",val)
	}
	if strings.Contains(this.COLUMN_TYPE,"varchar") {
		return fmt.Sprintf("'%s'",val)
	}
	if strings.Contains(this.COLUMN_TYPE,"text") {
		return fmt.Sprintf("'%s'",val)
	}
	if strings.Contains(this.COLUMN_TYPE,"float") {
		return fmt.Sprintf("%f",val)
	}
	return ""
}

//是否是删除的where 条件
func (this *ColumnInfo)IsUpdateOrDelWhere() bool {
	return (this.WhereFlag & WhereFlagUpdate)  != 0
}

//是否是查询的where 条件
func (this *ColumnInfo)IsQueryWhere() bool {
	return (this.WhereFlag & WhereFlagSelect)  != 0
}

