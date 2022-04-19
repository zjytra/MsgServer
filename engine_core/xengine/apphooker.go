/*
创建时间: 2020/12/20 1:09
作者: zjy
功能介绍:
要让其他层调用App模块的方法
*/

package xengine

type AppHooker interface {
	CloseApp()
}
