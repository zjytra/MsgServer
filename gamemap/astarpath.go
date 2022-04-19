/*
创建时间: 2021/7/11 14:01
作者: zjy
功能介绍:

*/

package gamemap

import (
	"fmt"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"math"
	"sync"
)

//a星寻路
//14 10 14
//10 0  10
//14 10 14
type AStarPath struct {
	openList  map[int32]*Node //等待查找的列表
	visitList map[int32]*Node //已经被检查过的节点
	startX PosType            //开始节点
	startY PosType
	endX PosType              //结束节点
	endY PosType
	nodePool  sync.Pool
	addOpenCount int16 //添加进开启列表中计数
}

var (
	AStar *AStarPath
	//14 10 14
	//10    10
	//14 10 14
	AroundArr = [8][2]int{
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
	}
)

func NewAStarPath() {
	AStar = new(AStarPath)
	AStar.visitList = make(map[int32]*Node)
	AStar.openList = make(map[int32]*Node)
	AStar.nodePool.New = func() interface{} {
		return new(Node)
	}
}

func (receiver *AStarPath) CreateNode(x PosType, y PosType) *Node {
	pdata := receiver.nodePool.Get()
	node, ok := pdata.(*Node)
	if ok {
		node.SetPos(x, y)
		return node
	}
	fmt.Printf("node x=%d y=%d create is  nil  \n", x, y)
	return nil
}
func (this *AStarPath) FindPathByInt16(startX int16, startY int16, endX int16, endY int16)[]*Pos2D {
	return this.FindPath(PosType(startX),PosType(startY),PosType(endX),PosType(endY))
}
func (this *AStarPath) FindPath(startX PosType, startY PosType, endX PosType, endY PosType)[]*Pos2D {
	//是否是同一个位置
	if 	IsOnePos(startX, startY, endX, endY) {
		return nil
	}

	if !GMap.CanMoveTo(startX, startY) || !GMap.CanMoveTo(endX, endY) {
		return nil
	}

	this.startX = startX
	this.startY = startY
	this.endX = endX
	this.endY = endY
	pStart := this.CreateNode(startX, startY)
	pStart.CalcHVal(endX, endX)
	//直接加入关闭列表中
	this.addToVisitList(pStart)
	//最开始获取start openlist
	this.setOpenList(pStart)
	for {
		minNode := this.getMinF()
		//xlog.Debug("minNode %v",minNode)
		//没有找到最小节点
		if minNode == nil {
			return nil
		}
		//把目标格添加进了关闭列表(注解)，这时候路径被找到，或者
		if len(this.openList) == 0 || 	IsOnePos(minNode.PosX, minNode.PosY, endX, endY) {
			//寻路结束
			this.findPathEnd()
			return this.getFindPath(minNode)
		}
		//开启列表数据到检查列表
		this.openToVisit(minNode)
		//以当前最小F的节点 查找周围的节点
		this.setOpenList(minNode)

	}

	//寻路结束
	this.findPathEnd()
	return nil
}


func (this *AStarPath)getFindPath(endNode *Node) []*Pos2D{
	var paths []*Pos2D
	paths = append(paths,NewPos2D(endNode.PosX,endNode.PosY))
	father := endNode.Father
	if father == nil {
		return paths
	}
	for {
		paths = append(paths,NewPos2D(father.PosX,father.PosY))
		father = father.Father
		if father == nil {
			break
		}
	}
	return paths
}
//路径寻找完毕
func (this *AStarPath) findPathEnd(){
	//已经访问过的列表
	for key, node := range this.visitList {
		node.Reset()
		this.nodePool.Put(node)
		delete(this.visitList,key)
	}
	//开启的列表
	for key, node := range this.openList {
		node.Reset()
		this.nodePool.Put(node)
		delete(this.openList,key)
	}
}


func (this *AStarPath) getMinF() *Node {
	var minFNode *Node 	 = nil
	for _, openNode := range this.openList {
		//time.Sleep(time.Second * 3)
		//xlog.Debug("当前节点 x= %d , y=%d f =%d , g = %d , h = %d",openNode.PosX,openNode.PosY,openNode.fVal,openNode.gVal,openNode.hVal)
		//设置最小的节点
		if minFNode == nil || openNode.fVal < minFNode.fVal {
			minFNode = openNode
			continue
		}
		//当f值相同,就看到达终点耗费最少
		if openNode.fVal == minFNode.fVal && openNode.addOpenCount >  minFNode.addOpenCount {
			minFNode = openNode
			continue
		}
	}
	if minFNode != nil {
		//time.Sleep(time.Second * 3)
		//xlog.Debug("选中 x= %d , y=%d f =%d , g = %d , h = %d",minFNode.PosX,minFNode.PosY,minFNode.fVal,minFNode.gVal,minFNode.hVal)
	}

	return minFNode

}

//从开启列表移动到关闭列表中
func (this *AStarPath) openToVisit(minFNode *Node){
	if minFNode == nil {
		return
	}
	delete(this.openList, minFNode.GetXYVal())
	//把当前的节点放到关闭列表中
	this.addToVisitList(minFNode)
}


//获取周边的格子
//并把他们添加open列表中
//14 10 14
//10 0  10
//14 10 14
func (this *AStarPath) setOpenList(node *Node) {
	var minFNode *Node 	 = nil
	var posX, posY PosType
	for _, ints := range AroundArr {
		posX = node.PosX + PosType(ints[0])
		posY = node.PosY + PosType(ints[1])
		gd := GMap.GetGrid(posX, posY)
		if gd == nil { //排除没有的格子
			continue
		}
		//在关闭列表中跳过
		if this.IsInVisitList(posX,posY) {
			continue
		}
		newNone := this.CreateNode(posX, posY)
		if !gd.CanMoveInto() {
			//不可以移动放到关闭列表中
			this.addToVisitList(newNone)
			continue
		}
		pOpenNode := this.GetOpen(posX, posY)
		if pOpenNode != nil {
			//如果它已经在开启列表中，用G值为参考检查新的路径是否更好。
			//更低的G值意味着更好的路径。
			//如果是这样，就把这一格的父节点改成当前格，并且重新计算这一格的G和F值。如果你保持你的开启列表按F值排序，改变之后你可能需要重新对开启列表排序。

			xVal := int32(math.Abs(float64(node.PosX - pOpenNode.PosX)))
			yVal := int32(math.Abs(float64(node.PosY - pOpenNode.PosY)))
			var temG int32
			if xVal == 0 || yVal == 0 {
				temG = 10
			}else if xVal == 1 && yVal == 1 { //走的斜角
				temG = 14
			}
			temG += node.gVal
			if temG >= pOpenNode.gVal {
				//G值不会更好就跳过
				continue
			}
			//重新计算G值与F值
			pOpenNode.SetFather(node)
			pOpenNode.CalcGVal(node.PosX, node.PosY)
			pOpenNode.CalcFVal()
			continue
		}
		//不再待检查的列表中设置为待检查列表
		newNone.SetFather(node)
		//
		newNone.CalcHVal(this.endX, this.endY)
		//计算g节点
		newNone.CalcGVal(node.PosX, node.PosY)
		newNone.CalcFVal()
		this.addOpenCount++
		newNone.addOpenCount = this.addOpenCount
		this.AddToOpenList(newNone)
		//设置最小的节点
		if minFNode == nil || newNone.fVal < minFNode.fVal {
			minFNode = newNone
			continue
		}
		//当f值相同,就看到达终点耗费最少
		if newNone.fVal == newNone.fVal && newNone.hVal <=  minFNode.hVal {
			minFNode = newNone
			continue
		}
		//已经到了结束的位置
		if IsOnePos(newNone.PosX,newNone.PosY,this.endX, this.endY) {
			break
		}
	}
}

func (receiver *AStarPath) AddToOpenList(node *Node) {
	pnode := receiver.GetOpen(node.PosX, node.PosY)
	if pnode != nil {
		xlog.Error("x = %d y = %d is in openList", node.PosX, node.PosY)
		return
	}
	receiver.openList[node.GetXYVal()] = node
}

//获取已经在open列表中的数据
func (receiver *AStarPath) GetOpen(x PosType, y PosType)*Node {
	xyVal := XYToVal(x,y)
	pnode, ok := receiver.openList[xyVal]
	if ok {
		//xlog.Debug("GetOpen x = %d y = %d is in openList", x, y)
		return pnode
	}
	return nil
}

func (receiver *AStarPath) addToVisitList(node *Node) {
	xyVal := XYToVal(node.PosX, node.PosY)
	_, ok := receiver.visitList[xyVal]
	if ok {
		xlog.Error("x = %d y = %d is in visitList", node.PosX, node.PosY)
		return
	}
	receiver.visitList[xyVal] = node
}

//是否在关闭列表中
func (receiver *AStarPath)IsInVisitList(posX, posY PosType)bool {
	xyVal := XYToVal(posX, posY)
	_, ok := receiver.visitList[xyVal]
	if	ok {
		return true
	}
	return false
}


func XYToVal(x PosType, y PosType) int32 {
	val := int32(y) << 16
	val = val | int32(x)
	return val
}

func ValToXY(val int32) (x PosType, y PosType) {
	y = PosType(val >> 16)
	x = PosType(val)
	return
}


func IsOnePos(startX PosType, startY PosType, endX PosType, endY PosType) bool {
	if startX == endX && startY == endY {
		return true
	}
	return false
}