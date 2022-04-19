/*
创建时间: 2021/7/11 14:02
作者: zjy
功能介绍:

*/

package gamemap

import (
	"fmt"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

type GameMap struct {
	//地图的格子
	mapGrids [][]*Grid
}

var (
	//地图管理对象
	GMap *GameMap

)

const(
	//最大x值
	MAXPosX PosType = 2000
	//最大Y值
	MAXPosY PosType = 2000
)

//坐标点是否有效
func PosIsValid(x PosType,y PosType) bool {
	if x < 0 || x >= MAXPosX || y < 0 ||  y >= MAXPosY  {
		return false
	}
	return true
}

func PosErrLog(x PosType,y PosType,fix string)  {
	xlog.Debug("%s x = %d  y = %d is err",fix, x,y)
}

//初始化
func GameMapInit()  {
	GMap = new(GameMap)
	GMap.CreateMap()
	NewAStarPath()
}

//根据地图文件创建地图
func (this *GameMap) CreateMap()  {
	var y PosType
	var x PosType
	this.mapGrids = make([][]*Grid,MAXPosY)

	for y = 0; y < MAXPosY; y++ {
		this.mapGrids[y] = make([]*Grid,MAXPosX)
		for x = 0; x < MAXPosX; x++ {
			grid := NewGrid(x,y)
			if grid == nil {
				continue
			}
			////rand := rand.Int31n(int32(MAXPosY) *int32(MAXPosX))
			//if 200000  <= rand {
			//	grid.SetCanMove(1)
			//}
			this.AddGrid(grid)
		}
	}
	this.ShowMap()
}

func (this *GameMap) ShowMap() {
	var y PosType
	var x PosType
	for y = MAXPosY - 1; y >= 0 ; y-- {
		fmt.Printf("y = %d  ",y)
		for x = 0; x < MAXPosX; x++ {
			gd := this.GetGrid(x,y)
			fmt.Printf("[%2d,%2d,%2d]",gd.PosX,gd.PosY,gd.canMove)
		}
		fmt.Printf("\n")
	}
}

func (this *GameMap) AddGrid(grid *Grid)  {
	if	grid == nil {
		return
	}
	this.mapGrids[grid.PosY][grid.PosX] = grid
}

func (this *GameMap) GetGrid(x PosType,y PosType)*Grid  {
	if	!PosIsValid(x,y) {
		//PosErrLog(x,y,"GetGrid")
		return nil
	}
	return this.mapGrids[y][x]
}

//根据x y 合并值获取格子
func (this *GameMap) GetGridByVal(Val int32)*Grid  {
	x,y := ValToXY(Val)
	return this.GetGrid(x,y)
}

//根据x y 判断该格子是否可以移入
func (this *GameMap) CanMoveTo(x PosType,y PosType)bool  {
 	gd  := this.GetGrid(x,y)
	if gd == nil {
		return false
	}
	return gd.CanMoveInto()
}

//根据x y 合并值判断该格子是否可以移入
func (this *GameMap) CanMoveToByVal(Val int32)bool  {
	gd  := this.GetGridByVal(Val)
	if gd == nil {
		return false
	}
	return gd.CanMoveInto()
}