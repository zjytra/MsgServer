/*
创建时间: 2021/4/18 2:14
作者: zjy
功能介绍:

*/

package dbsys
//运行时字段,方便数据库线程操作
type DBRunColumnar interface {
	initDBField() //初始化数据库字段
	GetCol() *ColumnInfo
	SetColInfo(*ColumnInfo)
	GetDBValAddr() interface{}
	dbValSetVal()  //数据库值到值的设置
	valSetDBVal()  //值到数据库
	DBGetVal() interface{} //数据库线程使用获得值
	IsChange() bool
}


type DBUInt64 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    uint64 //数据库中的数据
	val      uint64 //逻辑数据
}


//初始化字段
func (this *DBUInt64) initDBField() {

}

//查询出来的结果同步
func (this *DBUInt64) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBUInt64) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBUInt64) DBGetVal() interface{}{
	return  this.val
}

func (this *DBUInt64) GetDBValAddr() interface{}{
	return &this.dBVal
}
func (this *DBUInt64) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBUInt64) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBUInt64) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBUInt64) GetVal() uint64{
	return  this.val
}

func (this *DBUInt64) SetVal(val uint64){
	this.val = val
}

