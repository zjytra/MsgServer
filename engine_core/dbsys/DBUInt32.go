/*
创建时间: 2021/4/18 2:13
作者: zjy
功能介绍:

*/

package dbsys

type DBUInt32 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    uint32 //数据库中的数据
	val      uint32 //逻辑数据
}

//初始化字段
func (this *DBUInt32) initDBField() {

}
//查询出来的结果同步
func (this *DBUInt32) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBUInt32) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBUInt32) DBGetVal() interface{}{
	return  this.val
}



func (this *DBUInt32) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBUInt32) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBUInt32) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBUInt32) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBUInt32) GetVal() uint32{
	return  this.val
}


func (this *DBUInt32) SetVal(val uint32){
	this.val = val
}