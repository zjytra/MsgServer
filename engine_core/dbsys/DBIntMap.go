/*
创建时间: 2021/5/18 21:26
作者: zjy
功能介绍:

*/

package dbsys

import (
	"fmt"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strings"
)

type DBIntMap struct {
	col   *ColumnInfo     //指向的字段信息
	dBVal map[int32]int64 //数据库中的数据
	val   map[int32]int64 //逻辑数据
	str   string          //转换成str
}

func (this *DBIntMap) initDBField() {
	this.dBVal = make(map[int32]int64)
	this.val = make(map[int32]int64)
}

//查询出来的结果同步
//str 转成map
func (this *DBIntMap) dbValSetVal() {
	//分割外层
	arr := strings.Split(this.str, ";")
	if arr != nil && len(arr) > 0 {
		for _, s := range arr {
			if s == "" {
				continue
			} //分内部
			valArr := strings.Split(s, ",")
			if valArr == nil || len(valArr) < 2 {
				continue
			}
			this.dBVal[strutil.StrToInt32(valArr[0])] = strutil.StrToInt64(valArr[1])
		}
		//将数据库的值设置到运行中
		for u, i := range this.dBVal {
			this.val[u] = i
		}
	}

}

////str 转成map
//func (this *DBIntMap) DbValSetVal() {
//	this.dbValSetVal()
//}
////执行了sql后要将val值与数据库值同步
//func (this *DBIntMap) ValSetDBVal() {
//	this.valSetDBVal()
//}

//执行了sql后要将val值与数据库值同步
func (this *DBIntMap) valSetDBVal() {
	for u, i := range this.val {
		this.dBVal[u] = i
	}
}

func (this *DBIntMap) DBGetVal() interface{} {
	//拼接字符串
	this.str = ""
	for u, i := range this.val {
		this.str += fmt.Sprintf("%v,%v;", u, i)
	}
	return this.str
}

func (this *DBIntMap) GetDBValAddr() interface{} {
	return &this.str
}
func (this *DBIntMap) GetCol() *ColumnInfo {
	return this.col
}

func (this *DBIntMap) SetColInfo(info *ColumnInfo) {
	this.col = info
}

func (this *DBIntMap) IsChange() bool {
	//长度不相同
	if len(this.dBVal) != len(this.val) {
		return true
	}
	for id, val := range this.val {
		runVal, ok := this.dBVal[id]
		if !ok {
			return true
		}
		if val != runVal {
			return true
		}
	}
	return false
}

func (this *DBIntMap) GetVal() map[int32]int64 {
	//长度不相同
	return this.val
}

func (this *DBIntMap) SetValToMap(dis map[int32]int64) {
	if dis == nil {
		return
	}
	for id, val := range this.val {
		dis[id] = val
	}
}

func (this *DBIntMap) SetValFromMap(res map[int32]int64) {
	if res == nil {
		return
	}
	for id, val := range res {
		this.val[id] = val
	}
}

func (this *DBIntMap) AddVal(id int32, val int64) {
	this.val[id] = val
}

func (this *DBIntMap) GetValByKey(id int32) int64 {
	val, ok := this.val[id]
	if ok {
		return val
	}
	return 0
}
