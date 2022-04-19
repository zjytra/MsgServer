/*
创建时间: 2020/2/10
作者: zjy
功能介绍:

*/

package csvparse

import (
	"encoding/csv"
	"fmt"
	"github.com/zjytra/MsgServer/devlop/xutil"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const Invalid_Row = 3

// 获取csv 数据
// return map对象  gamemap[行数]gamemap[csv第一行每列的名称]对应列具体的数据
func GetCsvMapData(csvfile string) map[interface{}]map[string]interface{} {
	// 配置文件最好一次读完,一次写完
	alldata := GetCsvSliceData(csvfile)
	if alldata == nil {
		return nil
	}
	if len(alldata) < Invalid_Row {
		xlog.Error("解析csv %s 数据小于%d 行",csvfile,Invalid_Row)
		return nil
	}
	// 第一行的数据是字段 第二行类型
	nameRow  := alldata[0]
	typeRow  := alldata[1]
	csvdata := make(map[interface{}]map[string]interface{})
	for i, rows := range alldata {
		if i < Invalid_Row { // 排除前三行
			continue
		}
		filedInfo := make(map[string]interface{})
		for j, colval := range rows {
			SetFiledMap(filedInfo, nameRow[j], typeRow[j], colval)
		}
		// 转换id类型
		id := CsvStrToInterfaceStrT(typeRow[0], rows[0])
		// 第一个id
		csvdata[id] = filedInfo
	}
	return csvdata
}

// 获取csv 数据
// return 二维切片
func GetCsvSliceData(csvfile string) [][]string {
	// 创建csv文件
	fs, err := os.Open(csvfile)
	if xutil.IsError(err) {
		return nil
	}
	defer fs.Close() //这里关闭文件
	csvReader := csv.NewReader(fs)
	if csvReader == nil {
		return nil
	}
	// 配置文件最好一次读完,一次写完
	alldata, erro := csvReader.ReadAll()
	if xutil.IsError(erro) {
		return nil
	}
	return alldata
}

func SetFiledMap(filedInfo map[string]interface{}, filedName, filedType, filedval string) {
	if strutil.StringIsNil(filedName) || strutil.StringIsNil(filedType) {
		fmt.Println("csv字段数据解析 filedName = ", filedName, "filedType = ", filedType)
		return
	}
	// 首字母大写与字段对齐
	filedInfo[xutil.Capitalize(filedName)] = CsvStrToInterfaceStrT(filedType, filedval)
}

// 通过csv定义的类型,转换为go的内置类型
func CsvStrToInterfaceStrT(fliedtype string, strval string) interface{} {
	switch fliedtype {
	case "int":
		return strutil.StrToInt(strval)
	case "int8":
		return strutil.StrToInt8(strval)
	case "uint8":
		return strutil.StrToUint8(strval)
	case "int16":
		return strutil.StrToInt16(strval)
	case "uint16":
		return strutil.StrToUint16(strval)
	case "int32":
		return strutil.StrToInt32(strval)
	case "uint32":
		return strutil.StrToUint32(strval)
	case "float64":
		flt64, erro := strconv.ParseFloat(strval, 64)
		if !xutil.IsError(erro) {
			return flt64
		}
	case "string":
		return strval
	case "bool":
		b,erro := strconv.ParseBool(strval)
		if xutil.IsError(erro) {
		
		}
		return b
	case "[]int":
		intArr := StringsToIntArr(RepleaceBrackets(strval))
		return intArr
	case "[]string":
		return strings.Split(RepleaceBrackets(strval), ",")
	default:
		fmt.Println(fliedtype, "is an unknown type.")
		return nil
	}
	return nil
}

//替换掉中括号
func RepleaceBrackets(str string) string  {
	str = strings.ReplaceAll(str,"[","")
	str = strings.ReplaceAll(str,"]","")
	return str
}

//string ,号分割字符串转换为[]int
func StringsToIntArr(str string) []int {
	strarr := strings.Split(str, ",")
	if strarr == nil {
		return nil
	}
	var intarr []int
	for _, str := range strarr {
		inval, erro := strconv.Atoi(str)
		if xutil.IsError(erro) {
			continue
		}
		intarr = append(intarr, inval)
	}
	return intarr
}

// 通过反射的方式设置字段
// param obj 需要设置的结构体 这里必须是引用类型,可以取地址的
// param name 结构体字段名称
// param value 给结构体设置的值
func ReflectSetField(obj interface{}, name string, value interface{}) error {
	if obj == nil {
		return fmt.Errorf("obj is nil", name)
	}
	paramval := reflect.ValueOf(obj) // 参数值
	if !ValueCanElem(paramval) {     // Elem前验证下
		return fmt.Errorf("obj 无效", obj)
	}
	// won't work if I remove .Elem()
	structValue := paramval.Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	// won't work either if I add .Elem() to the end
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return fmt.Errorf("Provided value %v type %v didn't match obj field type %v", val, val.Type(), structFieldType)
	}
	structFieldValue.Set(val)
	return nil
}

// 验证val对象是否能取值
func ValueCanElem(value reflect.Value) bool {
	return value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr
}



//检查类型
func CheckType(filed interface{}, csvType string,csvfieldName string,csvName string)  bool{
	t := reflect.TypeOf(filed)
	if t == nil {
		return  false
	}
	fieldTypeName := t.Name()
	if fieldTypeName != csvType {
		fmt.Printf("表%v 字段%v 类型%v 不匹配结构体字段%v", csvName,csvfieldName,csvType,fieldTypeName)
		return false
	}
	t = nil
	return  true
}