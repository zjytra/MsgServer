/*
创建时间: 2020/2/17
作者: zjy
功能介绍:

*/

package csvdata

//初始化登陆服数据
func LoadLoginCsvData( )  {
	LoadCommonCsvData()
}

func ReLoadLoginCsvData()  {
	go LoadLoginCsvData()
}


