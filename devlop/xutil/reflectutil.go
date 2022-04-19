/*
创建时间: 2020/4/28
作者: zjy
功能介绍:
反射相关处理
*/

package xutil

import (
	"fmt"
	"reflect"
)

// 验证内置类型数组
func ValidArrIndex(arr interface{}, index int) bool {
	if arr == nil {
		return false
	}
	// 下标为负
	if index < 0 {
		return false
	}
	switch val := arr.(type) {
	case []int:
		return index < len(val)
	case []string:
		return index < len(val)
	case []float32:
		return index < len(val)
	default:
		fmt.Println(arr, "is an unknown type. ")
		return false
	}
	return true
}

//这里判断接口是否为nil
func InterFaceIsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	switch vi.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:

		return vi.IsNil()
	}
	return false
}
