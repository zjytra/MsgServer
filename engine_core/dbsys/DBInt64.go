/*
创建时间: 2021/4/18 2:14
作者: zjy
功能介绍:

*/

package dbsys


type DBInt64 struct {
	col   *ColumnInfo //指向的字段信息
	dBVal int64       //数据库中的数据
	val   int64       //逻辑数据
}

//初始化字段
func (this *DBInt64) initDBField() {

}
//查询出来的结果同步
func (this *DBInt64) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBInt64) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBInt64) DBGetVal() interface{}{
	return  this.val
}

func (this *DBInt64) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBInt64) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBInt64) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBInt64) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBInt64) GetVal() int64{
	return  this.val
}
func (this *DBInt64) SetVal(val int64){
	this.val = val
}