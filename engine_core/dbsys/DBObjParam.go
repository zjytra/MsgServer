/*
创建时间: 2021/9/13 21:27
作者: zjy
功能介绍:

*/

package dbsys

//投递参数
type DBObjQueryParam struct {
	ClientConnId uint32      //网络连接id
	ServerConnId uint32      //服务器连接id
	ParamObj     interface{} //需要带回的参数,
	DbObj        DBObJer     //传入的对象构建Sql对象
}

func (this *DBObjQueryParam) CheckParam() bool {
	return this.DbObj != nil
}

//投递参数
type DBObjWriteParam struct {
	ClientConnId uint32    //网络连接id
	ServerConnId uint32    //服务器连接id
	DbObjs       []DBObJer //传入的对象构建Sql对象 传入的数量尽量一次写完
	writeTpe     int8      //写操作类型
	ParamObj     interface{} //需要带回的参数,
}

func (this *DBObjWriteParam) CheckParam() bool {
	return this.DbObjs != nil && len(this.DbObjs) > 0
}

func (this *DBObjWriteParam) AddObjs(obj ...DBObJer) {
	this.DbObjs = append(this.DbObjs, obj...)
}
