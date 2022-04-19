/*
创建时间: 2021/4/18 2:13
作者: zjy
功能介绍:

*/

package dbsys

type DBUInt16 struct {
	col   *ColumnInfo //指向的字段信息
	dBVal uint16      //数据库中的数据
	val   uint16      //逻辑数据
}

//初始化字段
func (this *DBUInt16) initDBField() {

}

//查询出来的结果同步
func (this *DBUInt16) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBUInt16) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBUInt16) DBGetVal() interface{} {
	return this.val
}



func (this *DBUInt16) GetDBValAddr() interface{} {
	return &this.dBVal
}

func (this *DBUInt16) GetCol() *ColumnInfo {
	return this.col
}

func (this *DBUInt16) SetColInfo(info *ColumnInfo) {
	this.col = info
}

func (this *DBUInt16) IsChange() bool {
	return this.dBVal != this.val
}


func (this *DBUInt16) GetVal() uint16{
	return  this.val
}


func (this *DBUInt16) SetVal(val uint16){
	this.val = val
}