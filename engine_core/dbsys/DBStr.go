/*
创建时间: 2021/4/18 0:14
作者: zjy
功能介绍:

*/

package dbsys

type DBStr struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    string //数据库中的数据
	val      string //逻辑数据
}
//初始化字段
func (this *DBStr) initDBField() {

}
//查询出来的结果同步
func (this *DBStr) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBStr) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBStr) DBGetVal() interface{}{
	return  this.val
}


func (this *DBStr) GetDBValAddr() interface{}{
	return &this.dBVal
}
func (this *DBStr) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBStr) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBStr) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBStr) GetVal() string{
	return  this.val
}


func (this *DBStr) SetVal(val string){
	this.val = val
}