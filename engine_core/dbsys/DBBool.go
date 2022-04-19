/*
创建时间: 2021/4/5 21:36
作者: zjy
功能介绍:

*/

package dbsys

type DBBool struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    bool //数据库中的数据
	val      bool //逻辑数据
}

//初始化字段
func (this *DBBool) initDBField() {

}
//查询出来的结果同步
func (this *DBBool) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBBool) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBBool) DBGetVal() interface{}{
	return  this.val
}



func (this *DBBool) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBBool) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBBool) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBBool) IsChange() bool{
	return this.dBVal != this.val
}