/*
创建时间: 2021/8/22 22:14
作者: zjy
功能介绍:
数学工具包
*/

package mathutil

import "math"

func AddInt8(oldVal, val int8) int8  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxInt8 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxInt8
	}
	return oldVal + val
}

func AddUInt8(oldVal,val uint8) uint8  {
	//看以前还差多少到最大值
	diff := math.MaxUint8 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxUint8
	}
	return oldVal + val
}

func SubInt8(oldVal,val int8) int8  {
	if val < 0 { //支持负数
		return AddInt8(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func SubUInt8(oldVal,val uint8) uint8  {
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func AddInt64(oldVal, val int64) int64  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxInt64 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxInt64
	}
	return oldVal + val
}

func AddUInt64(oldVal,val uint64) uint64  {
	//看以前还差多少到最大值
	diff := math.MaxUint64 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxUint64
	}
	return oldVal + val
}

func SubInt64(oldVal,val int64) int64  {
	if val < 0 { //支持负数
		return AddInt64(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func SubUInt64(oldVal,val uint64) uint64  {
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func AddInt32(oldVal, val int32) int32  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxInt32 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxInt32
	}
	return oldVal + val
}

func AddUInt32(oldVal,val uint32) uint32  {
	//看以前还差多少到最大值
	diff := math.MaxUint32 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxUint32
	}
	return oldVal + val
}

func SubInt32(oldVal,val int32) int32  {
	if val < 0 { //支持负数
		return AddInt32(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func SubUInt32(oldVal,val uint32) uint32  {
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func AddInt16(oldVal, val int16) int16  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxInt16 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxInt16
	}
	return oldVal + val
}

func AddUInt16(oldVal,val uint16) uint16  {
	//看以前还差多少到最大值
	diff := math.MaxUint16 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxUint16
	}
	return oldVal + val
}

func SubInt16(oldVal,val int16) int16  {
	if val < 0 { //支持负数
		return AddInt16(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func SubUInt16(oldVal,val uint16) uint16  {
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}


func AddInt(oldVal, val int) int  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxInt32 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxInt32
	}
	return oldVal + val
}


func AddFloat32(oldVal, val float32) float32  {
	//为负数
	if val < 0 {
		//扣完了
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxFloat32 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxFloat32
	}
	return oldVal + val
}

func SubFloat32(oldVal,val float32) float32  {
	if val < 0 { //支持负数
		return AddFloat32(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}

func AddFloat64(oldVal, val float64) float64  {
	//为负数
	if val < 0 {
		tem := oldVal + val
		if tem < 0 {
			return 0
		}
		//能够扣的情况
		return oldVal + val
	}
	//看以前还差多少到最大值
	diff := math.MaxFloat64 - oldVal
	if diff  < val { //加上新的值绝对超过
		return math.MaxFloat64
	}
	return oldVal + val
}

func SubFloat64(oldVal,val float64) float64  {
	if val < 0 { //支持负数
		return AddFloat64(oldVal,val)
	}
	if oldVal >  val {
		return oldVal - val
	}
	return 0
}


func MaxInt64(one, two int64) int64 {
	if one >= two {
		return one
	}
	return two
}

func MinUInt64(one, two uint64) uint64 {
	if one <= two {
		return one
	}
	return two
}

func MaxUInt64(one, two uint64) uint64 {
	if one >= two {
		return one
	}
	return two
}

func MinInt64(one, two int64) int64 {
	if one <= two {
		return one
	}
	return two
}

func MaxUint32(one, two uint32) uint32 {
	return uint32(MaxUInt64(uint64(one),uint64(two)))
}

func MinUint32(one, two uint32) uint32 {
	return uint32(MinUInt64(uint64(one),uint64(two)))
}

func MaxInt(one, two int) int {
	if one >= two {
		return one
	}
	return two
}

func MinInt(one, two int) int {
	if one <= two {
		return one
	}
	return two
}

func MaxInt32(one, two int32) int32 {
	if one >= two {
		return one
	}
	return two
}

func MinInt32(one, two int32) int32 {
	if one <= two {
		return one
	}
	return two
}

func MaxInt16(one, two int16) int16 {
	if one >= two {
		return one
	}
	return two
}

func MinInt16(one, two int16) int16 {
	if one <= two {
		return one
	}
	return two
}

func MaxUInt16(one, two uint16) uint16 {
	if one >= two {
		return one
	}
	return two
}

func MinUInt16(one, two uint16) uint16 {
	if one <= two {
		return one
	}
	return two
}

func MaxInt8(one, two int8) int8 {
	if one >= two {
		return one
	}
	return two
}

func MinInt8(one, two int8) int8 {
	if one <= two {
		return one
	}
	return two
}

func MaxUInt8(one, two uint8) uint8 {
	if one >= two {
		return one
	}
	return two
}

func MinUInt8(one, two uint8) uint8 {
	if one <= two {
		return one
	}
	return two
}


// 将两个16位命令组合在一起
func MakeUint32(main, sub uint16) uint32 {
	return uint32(main)<<16 | uint32(sub)
}

// 将32位拆为两个 16位
func UnUint32(cmd uint32) (mn, sub uint16) {
	mn = uint16(cmd >> 16)
	// 去掉高位
	sub = uint16(cmd & math.MaxUint16)
	return
}