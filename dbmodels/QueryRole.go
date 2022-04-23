/*
创建时间: 2022/4/23 1:21
作者: zjy
功能介绍:

*/

package dbmodels

type QueryRole struct {
	RoleName  string
	LonginTime int64
	OnlineTime int64
	RoomID   int32
	OffLineTime int64
}