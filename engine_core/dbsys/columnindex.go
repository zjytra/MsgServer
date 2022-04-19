/*
创建时间: 2021/4/20 21:34
作者: zjy
功能介绍:
数据库表索引描述
*/

package dbsys


type ColumnIndex struct {
	TABLE_NAME  string //表名
	NON_UNIQUE  int //是否是唯一索引
	INDEX_NAME  string //索引名称
	INDEX_Type_NAME  string //索引类型
	SEQ_IN_INDEX int   //多索引排序
	COLUMN_NAME string //列名
	More_INDEX_NAME string  //联合索引名称
	More_Type_NAME string  //联合索引类型 是UNIQUE 还是 NORMAL
	MoreINDEX_Cols string  //联合索引列
}