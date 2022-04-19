/*
创建时间: 2020/2/3
作者: zjy
功能介绍:

*/

package global

import (
	"fmt"
	"runtime"
)





// 拉起宕机标准输出
func GrecoverToStd() {
	if rec := recover(); rec != nil {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		fmt.Printf("%v\n%s \n", rec, buf[:l])
	}
}
