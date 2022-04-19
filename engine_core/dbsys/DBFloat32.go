/*
创建时间: 2021/4/18 2:12
作者: zjy
功能介绍:

*/

package dbsys

type DBFloat32 struct {
	col      *ColumnInfo //指向的字段信息
	dBVal    float32 //数据库中的数据
	val      float32 //逻辑数据
}

//初始化字段
func (this *DBFloat32) initDBField() {

}
//查询出来的结果同步
func (this *DBFloat32) dbValSetVal() {
	this.val = this.dBVal
}

//执行了sql后要将val值与数据库值同步
func (this *DBFloat32) valSetDBVal() {
	this.dBVal = this.val
}

func (this *DBFloat32) DBGetVal() interface{}{
	return  this.val
}



func (this *DBFloat32) GetDBValAddr() interface{}{
	return &this.dBVal
}

func (this *DBFloat32) GetCol() *ColumnInfo{
	return this.col
}

func (this *DBFloat32) SetColInfo(info *ColumnInfo){
	this.col = info
}

func (this *DBFloat32) IsChange() bool{
	return this.dBVal != this.val
}

func (this *DBFloat32) GetVal() float32{
	return  this.val
}


func (this *DBFloat32) SetVal(val float32){
	this.val = val
}