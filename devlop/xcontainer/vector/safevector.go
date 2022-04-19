package vector

import (
	"sync"
)

type SafeVector struct {
	vec *Vector
	lock sync.RWMutex
}

// 创建工厂函数
func NewSafeVector(cap int) *SafeVector {
	this := new(SafeVector)
	this.vec = NewVector(cap)
	return this
}

func (this *SafeVector) IsEmpty() bool {
	this.lock.RLock()
	isEmpty := this.vec.IsEmpty()
	this.lock.RUnlock()
	return isEmpty
}

// 元素数量
func (this *SafeVector) Size() int {
	this.lock.RLock()
	size := this.vec.Size()
	this.lock.RUnlock()
	return size
}

//// 向队头添加元素
//// 慎用，效率没得操作队尾好
//func (this *Vector) PushFront(Account interface{}) bool {
//	this.values = append([]interface{}{Account},this.values...)
//	return true
//}

// 追加单个元素
func (this *SafeVector) PushBack(data interface{}) bool {
	this.lock.Lock()
	this.vec.PushBack(data)
	this.lock.Unlock()
	return true
}

// 获取队头的元素
func (this *SafeVector) Front() interface{} {
	return this.GetValue(0)
}

// 获取队尾元素
func (this *SafeVector) Back() interface{} {
	return this.GetValue(this.Size() - 1)
}

// 队头弹出元素
func (this *SafeVector) PopFront() interface{} {
	this.lock.Lock()
	if this.IsEmpty() {
		this.lock.Unlock()
		return nil
	}
	data := this.vec.PopFront()
	this.lock.Unlock()
	return data
}

// 队尾弹出元素
func (this *SafeVector) PopBack() interface{} {
	this.lock.Lock()
	if this.IsEmpty() {
		this.lock.Unlock()
		return nil
	}
	data := this.vec.PopBack()
	this.lock.Unlock()
	return data
}

// 追加元素切片
func (this *SafeVector) AppendAll(values []interface{}) bool {
	this.lock.Lock()
	if values == nil || len(values) == 0 {
		this.lock.Unlock()
		return false
	}
	this.vec.AppendAll(values)
	this.lock.Unlock()
	return true
}

// 插入单个元素
func (this *SafeVector) Insert(index int, data interface{}) bool {
	this.lock.Lock()
	if index < 0 || index >= this.vec.Size() {
		this.lock.Unlock()
		return false
	}
	this.vec.Insert(index,data)
	this.lock.Unlock()
	return true
}

// 插入元素切片
func (this *SafeVector) InsertAll(index int, values []interface{}) bool {
	this.lock.Lock()
	if index < 0 || index >= this.vec.Size() || values == nil || len(values) < 1 {
		this.lock.Unlock()
		return false
	}
	this.vec.InsertAll(index,values)
	this.lock.Unlock()
	return true
}

// 移除
func (this *SafeVector) RemoveAt(index int) bool {
	this.lock.Lock()
	if index < 0 || index >= this.vec.Size() {
		this.lock.Unlock()
		return false
	}
	this.vec.RemoveAt(index)
	this.lock.Unlock()
	return true
}

// 范围移除 从 fromIndex(包含) 到 toIndex(不包含) 之间的元素 [fromIndex,toIndex)
func (this *SafeVector) RemoveRange(fromIndex, toIndex int) bool {
	this.lock.Lock()
	if fromIndex < 0 || fromIndex >= this.vec.Size() || toIndex > this.vec.Size() || fromIndex > toIndex {
		this.lock.Unlock()
		return false
	}
	this.vec.RemoveRange(fromIndex,toIndex)
	this.lock.Unlock()
	return true
}

// 全部移除
func (this *SafeVector) RemoveAll() {
	this.lock.Lock()
	this.vec.RemoveAll()
	this.lock.Unlock()
}

func (this *SafeVector) getIndex(data interface{}) int {
	this.lock.RLock()
	i := this.vec.getIndex(data)
	this.lock.RUnlock()
	return i
}

// 是否存在该元素值 还是不要为好效率不高
func (this *SafeVector) Contains(data interface{}) bool {
	return this.getIndex(data) >= 0
}

// 获取元素值第一次出现的索引
func (this *SafeVector) IndexOf(data interface{}) int {
	return this.getIndex(data)
}

// 获取元素值最后一次出现的索引
func (this *SafeVector) LastIndexOf(data interface{}) int {
	this.lock.RLock()
	i := this.vec.LastIndexOf(data)
	this.lock.RUnlock()
	return i
}

// 得到索引对应的元素值
func (this *SafeVector) GetValue(index int) interface{} {
	this.lock.RLock()
	if index < 0 || index >= this.Size() {
		this.lock.RUnlock()
		return nil
	}
	data := this.vec.GetValue(index)
	this.lock.RUnlock()
	return data
}

// 设置值
func (this *SafeVector) SetValue(index int, data interface{}) bool {
	this.lock.Lock()
	if index < 0 || index >= this.Size() {
		this.lock.Unlock()
		return false
	}
	this.vec.SetValue(index,data)
	this.lock.Unlock()
	return true
}

//返回的是拷贝
func (this *SafeVector) ToArray() []interface{} {
	this.lock.RLock()
	dst := this.vec.ToArray()
	this.lock.RUnlock()
	return dst
}

//返回全部
func (this *SafeVector) GetData() []interface{} {
	this.lock.RLock()
	data := this.vec.GetData()
	this.lock.RUnlock()
	return data
}