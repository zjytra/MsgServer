/*
创建时间: 2021/4/5 21:36
作者: zjy
功能介绍:

*/

package dbsys

type DBInt8 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    int8 //数据库中的数据
	val      int8 //逻辑数据
}

//初始化字段
func (this *DBInt8) initDBField() {

}

//查询出来的结果同步
func (this *DBInt8) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBInt8) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBInt8) DBGetVal() interface{}{
	return  this.val
}


func (this *DBInt8) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBInt8) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBInt8) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBInt8) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBInt8) GetVal() int8{
	return  this.val
}


func (this *DBInt8) SetVal(val int8){
	this.val = val
}