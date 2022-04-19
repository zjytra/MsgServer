/*
创建时间: 2021/4/18 2:12
作者: zjy
功能介绍:

*/

package dbsys

type DBInt16 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    int16 //数据库中的数据
	val      int16 //逻辑数据
}

//初始化字段
func (this *DBInt16) initDBField() {

}
//查询出来的结果同步
func (this *DBInt16) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBInt16) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBInt16) DBGetVal() interface{}{
	return  this.val
}



func (this *DBInt16) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBInt16) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBInt16) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBInt16) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBInt16) GetVal() int16{
	return  this.val
}


func (this *DBInt16) SetVal(val int16){
	this.val = val
}