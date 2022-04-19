/*
创建时间: 2021/4/18 2:10
作者: zjy
功能介绍:

*/

package dbsys

type DBUInt8 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    uint8 //数据库中的数据
	val      uint8 //逻辑数据
}

//初始化字段
func (this *DBUInt8) initDBField() {

}

//查询出来的结果同步
func (this *DBUInt8) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBUInt8) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBUInt8) DBGetVal() interface{}{
	return  this.val
}


func (this *DBUInt8) GetDBValAddr() interface{}{
	return &this.dBVal
}
func (this *DBUInt8) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBUInt8) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBUInt8) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBUInt8) GetVal() uint8{
	return  this.val
}


func (this *DBUInt8) SetVal(val uint8){
	this.val = val
}