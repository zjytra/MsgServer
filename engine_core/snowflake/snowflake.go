/*
创建时间: 2020/3/28
作者: zjy
功能介绍:
Snowflake 算法
毫秒级 分布式生成全局唯一id
*/

package snowflake

import (
	"errors"
	"fmt"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"sync"
)

const (
	Twepoch       = int64(1584623409000) // 默认起始的时间 2020/03/19 21:30:00时间戳 1584624600 。计算时，减去这个值用41位存
	maxNextIdsNum = 100                  // 单次获取ID的最大数量
	MaxRightBit   = uint(22)             // 右边最大位数
)

// 1bit |41bit 	              |10bit |12bit
// 符号位|时间戳(当前时间-纪元时间)|机器id|自增序列
type NodeGID struct {
	sequence      int64 // 序号
	lastTimestamp int64 // 最后时间戳
	nodeId        int64 // 节点ID
	
	nodeIdBits         uint  // 节点 所占位置对应机器的位置 1023
	sequenceBits       uint  // 自增ID 所占用位置
	maxNodeId          int64 // 节点 GID 最大范围
	sequenceMask       int64 // 序列的最大范围
	nodeIdShift        uint  // 左移次数
	timestampLeftShift uint  // 时间戳最终要移动的位数
	mtx                sync.Mutex
}

var GUID *NodeGID

func init()  {
	var err error
	GUID,err = NewNodeGID(int16(appdata.AppID),10)
	if err != nil {
		xlog.Error("%v",err)
	}
}


// NewNodeGID new a snowflake id generator object.
// @param NodeId 节点 根据位数设置节点的值 0 到  -1 ^ (-1 << nodeidbits)1023,区分不同服务器的
// @param nodeidbits 节点的位数 最好在 10 - 12之间
func NewNodeGID(nodeId int16, nodeidbits uint) (*NodeGID, error) {
	fid := new(NodeGID)
	fid.setNodeIdBits(nodeidbits)
	if int64(nodeId) > fid.maxNodeId || nodeId < 0 {
		return nil, errors.New(fmt.Sprintf("nodeId Id can't be greater than %d or less than 0", fid.maxNodeId))
	}
	fid.nodeId = int64(nodeId)
	fid.lastTimestamp = timeutil.GetCurrentTimeMs()// 当前时间作为最后一次更新
	fid.sequence = 0
	return fid, nil
}

// 节点位数 + 序列位数  = 10 + 12 = 22 位
func (fid *NodeGID) setNodeIdBits(nodeidbits uint) {
	fid.nodeIdBits = nodeidbits
	fid.sequenceBits = MaxRightBit - nodeidbits
	if fid.sequenceBits < 5 {
		fmt.Println("  sequenceBits 位数太少只有 ", nodeidbits)
		// 进行重置
		fid.nodeIdBits = 10
		fid.sequenceBits = 12
	}
	fid.maxNodeId = -1 ^ (-1 << fid.nodeIdBits)                // 节点 GID 最大范围
	fid.nodeIdShift = fid.sequenceBits                         // 节点左移次数
	fid.sequenceMask = -1 ^ (-1 << fid.sequenceBits)           // 序列的最大范围
	fid.timestampLeftShift = fid.nodeIdBits + fid.sequenceBits // 时间搓要左移的次数
}

// timeGen generate a unix millisecond.


// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeutil.GetCurrentTimeMs()
	for timestamp <= lastTimestamp { //循环等待到下一毫秒
		timestamp = timeutil.GetCurrentTimeMs()
	}
	return timestamp
}

// NextIds get snowflake ids.
func (id *NodeGID) NextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		return nil, errors.New(fmt.Sprintf("NextIds num can't be greater than %d or less than 0", maxNextIdsNum))
	}
	ids := make([]int64, num)
	id.mtx.Lock()
	for i := 0; i < num; i++ {
		ids[i] = id.NextId()
	}
	id.mtx.Unlock()
	return ids, nil
}

func (this *NodeGID) nextId() (int64)  {
	timestamp := timeutil.GetCurrentTimeMs()
	// 时间靠前了
	if timestamp < this.lastTimestamp {
		fmt.Println(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", this.lastTimestamp-timestamp))
		return 0
	}
	if this.lastTimestamp == timestamp {
		// 在同一毫秒内叠加 序列，并不能超过最大值
		this.sequence = (this.sequence + 1) & this.sequenceMask
		if this.sequence == 0 { // 如果同一秒内序列用完，就到下一秒
			timestamp = tilNextMillis(this.lastTimestamp)
		}
	} else {
		this.sequence = 0
	}
	this.lastTimestamp = timestamp
	return ((timestamp - Twepoch) << this.timestampLeftShift) | (this.nodeId << this.nodeIdShift) | this.sequence
}

func (this *NodeGID) NextId() (int64) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	return  this.nextId()
}

func (this *NodeGID) Time(id int64) int64 {
	return (id >> this.timestampLeftShift) + Twepoch
}

func (this *NodeGID) Node() int64 {
	return this.nodeId
}

func (this *NodeGID) Step(id int64) int64 {
	return id & this.sequenceMask
}

// 生成的id
type GID int64

// Time returns an int64 unix timestamp in milliseconds of the snowflake GID time
// DEPRECATED: the below function will be removed in a future release.
func (id GID) Time(fast *NodeGID) int64 {
	return (int64(id) >> fast.timestampLeftShift) + Twepoch
}

// Node returns an int64 of the snowflake GID node number
// DEPRECATED: the below function will be removed in a future release.
func (id GID) Node(fast *NodeGID) int64 {
	return fast.nodeId
}

// Step returns an int64 of the snowflake step (or sequence) number
// DEPRECATED: the below function will be removed in a future release.
func (id GID) Step(fast *NodeGID) int64 {
	return int64(id) & fast.sequenceMask
}
