/*
创建时间: 2020/12/20 1:20
作者: zjy
功能介绍:

*/

package invoker

type CallID int

const (
	AppClose  = 1
	OnServerDisconnect = 2 //服务器关闭
)
