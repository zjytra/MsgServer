/*
创建时间: 2021/7/11 15:57
作者: zjy
功能介绍:

*/

package gamemap

import "math"

type Node struct {
	PosX         PosType
	PosY         PosType
	fVal         int32 //方块的总移动代价
	gVal         int32 //开始点到当前方块的移动代价
	hVal         int32 //当前方块到结束点的预估移动代价
	addOpenCount int16 //添加进开启列表中计数
	Father       *Node //父节点
}

func NewNode() *Node {
	node := new(Node)
	return node
}

func (this *Node) SetFather(father *Node) {
	this.Father = father

}

func (this *Node) SetPos(x PosType, y PosType) {
	this.PosX = x
	this.PosY = y
}

//计算开始点到当前点的移动代价
//14 10 14
//10 0  10
//14 10 14
func (this *Node) CalcGVal(startx PosType, starty PosType) {
	xVal := int32(math.Abs(float64(startx - this.PosX)))
	yVal := int32(math.Abs(float64(starty - this.PosY)))
	var temG int32
	if xVal == 0 || yVal == 0 {
		temG = 10
	}else if xVal == 1 && yVal == 1 { //走的斜角
		temG = 14
	}
	if this.Father != nil {
		temG += this.Father.gVal
	}
	this.gVal = temG
}

//计算到终点的移动代价 曼哈顿街区算法
func (this *Node) CalcHVal(endx PosType, endy PosType) {
	xVal := endx - this.PosX
	yVal := endy - this.PosY
	this.hVal = int32(math.Abs(float64(xVal)) + math.Abs(float64(yVal))) * 10
}

//计算移动代价
func (this *Node) CalcFVal() {
	this.fVal = this.gVal + this.hVal
}

func (this *Node) GetXYVal() int32 {
	return XYToVal(this.PosX,this.PosY)
}


func (this *Node) Reset()  {
	this.PosX = 0
	this.PosY = 0
	this.fVal = 0
	this.gVal = 0
	this.hVal = 0
	this.addOpenCount = 0
	this.Father = nil
}
