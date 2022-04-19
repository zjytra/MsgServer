package vector

import (
	"reflect"
)

type Vector struct {
	values []interface{}
}

// 创建工厂函数
func NewVector(cap int) *Vector {
	this := new(Vector)
	this.values = make([]interface{}, 0, cap)
	return this
}

func (this *Vector) IsEmpty() bool {
	return len(this.values) == 0
}

// 元素数量
func (this *Vector) Size() int {
	return len(this.values)
}

//// 向队头添加元素
//// 慎用，效率没得操作队尾好
//func (this *Vector) PushFront(Account interface{}) bool {
//	this.values = append([]interface{}{Account},this.values...)
//	return true
//}

// 追加单个元素
func (this *Vector) PushBack(data interface{}) bool {
	this.values = append(this.values, data)
	return true
}

// 获取队头的元素
func (this *Vector) Front() interface{} {
	return this.GetValue(0)
}

// 获取队尾元素
func (this *Vector) Back() interface{} {
	return this.GetValue(this.Size() - 1)
}

// 队头弹出元素
func (this *Vector) PopFront() interface{} {
	if this.IsEmpty() {
		return nil
	}
	data := this.values[0]
	this.values[0] = nil
	this.values = this.values[1:]
	return data
}

// 队尾弹出元素
func (this *Vector) PopBack() interface{} {
	if this.IsEmpty() {
		return nil
	}
	popindex := this.Size() - 1
	data := this.values[popindex]
	this.values[popindex] = nil
	this.values = this.values[:popindex]
	return data
}

// 追加元素切片
func (this *Vector) AppendAll(values []interface{}) bool {
	if values == nil || len(values) == 0 {
		return false
	}
	this.values = append(this.values, values...)
	return true
}

// 插入单个元素
func (this *Vector) Insert(index int, data interface{}) bool {
	if index < 0 || index >= this.Size() {
		return false
	}
	this.values = append(this.values[:index], append([]interface{}{data}, this.values[index:]...)...)
	return true
}

// 插入元素切片
func (this *Vector) InsertAll(index int, values []interface{}) bool {
	if index < 0 || index >= this.Size() || values == nil || len(values) < 1 {
		return false
	}
	this.values = append(this.values[:index], append(values, this.values[index:]...)...)
	return true
}

// 移除
func (this *Vector) RemoveAt(index int) bool {
	if index < 0 || index >= len(this.values) {
		return false
	}
	// 重置为 nil 防止内存泄漏
	this.values[index] = nil
	this.values = append(this.values[:index], this.values[index+1:]...)
	return true
}

// 范围移除 从 fromIndex(包含) 到 toIndex(不包含) 之间的元素 [fromIndex,toIndex)
func (this *Vector) RemoveRange(fromIndex, toIndex int) bool {
	if fromIndex < 0 || fromIndex >= len(this.values) || toIndex > len(this.values) || fromIndex > toIndex {
		return false
	}
	// 重置为 nil 防止内存泄漏
	for i := fromIndex; i < toIndex; i++ {
		this.values[i] = nil
	}
	this.values = append(this.values[:fromIndex], this.values[toIndex:]...)
	return true
}

// 全部移除
func (this *Vector) RemoveAll() {
	// 重置为 nil 防止内存泄漏
	for i := 0; i < this.Size(); i++ {
		this.values[i] = nil
	}
	this.values = this.values[0:0]
}

func (this *Vector) getIndex(data interface{}) int {
	lens := this.Size()
	for i := 0; i < lens; i++ {
		if reflect.DeepEqual(this.values[i], data) {
			return i
		}
	}
	return -1
}

// 是否存在该元素值 还是不要为好效率不高
func (this *Vector) Contains(data interface{}) bool {
	return this.getIndex(data) >= 0
}

// 获取元素值第一次出现的索引
func (this *Vector) IndexOf(data interface{}) int {
	return this.getIndex(data)
}

// 获取元素值最后一次出现的索引
func (this *Vector) LastIndexOf(data interface{}) int {
	lens := this.Size() - 1
	for i := lens; i >= 0; i-- {
		if reflect.DeepEqual(this.values[i], data) {
			return i
		}
	}
	return -1
}

// 得到索引对应的元素值
func (this *Vector) GetValue(index int) interface{} {
	if index < 0 || index >= this.Size() {
		return nil
	}
	return this.values[index]
}

// 设置值
func (this *Vector) SetValue(index int, data interface{}) bool {
	if index < 0 || index >= this.Size() {
		return false
	}
	this.values[index] = data
	return true
}

//返回的是拷贝
func (this *Vector) ToArray() []interface{} {
	dst := make([]interface{}, this.Size())
	copy(dst, this.values)
	return dst
}

//返回全部
func (this *Vector) GetData() []interface{} {
	return this.values
}
