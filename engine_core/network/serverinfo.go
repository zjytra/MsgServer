/*
创建时间: 2020/5/17
作者: zjy
功能介绍:

*/

package network

//服务器信息主要用于转发
type SeverInfo struct {
	AppId   int32     //服务器id 每个服务器id 为一个,同组服务器id 不能重复
	Status  int32    //状态
	Num     int32   //连接数
	GroupId int32   //组id 表示一个区的服务器
	IP      string  //ip地址
	Port    int32   //服务器端口
	Kind    int32  //服务器类型
	Name    string //appName
}

func (this *SeverInfo) SetServerInfo(appid int32, _status int32, num int32)  {
	this.AppId = appid
	this.Status = _status
	this.Num = num
}

const (
	ServerStauts_None = 0 //未知
	ServerStauts_Online  = 1 //在线
	ServerStauts_Offline = 2 //离线
	ServerStauts_Pause   = 3 //服务器维护
	ServerStauts_BeBusy   = 4 //忙碌
)

