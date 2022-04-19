// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 1.
package model

// app基本信息


// app 的网络信息
type AppNetWorkModel struct {
	OutAddr     string // 外部访问地址
	OutPort     int    // 对外端口号
	MaxConnet   int    // 限制最大连接数
	SendMaxSize int    // 最大发送多少字节
	RecMaxSize  int    // 最大接受多少字节
}
