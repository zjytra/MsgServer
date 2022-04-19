/*
创建时间: 2021/7/11 14:02
作者: zjy
功能介绍:
//地图格子
*/

package gamemap

type Grid struct {
	PosX PosType   //x位置
	PosY PosType   //y位子
	canMove int8  //是否能移动
}

func NewGrid(x PosType,y PosType)*Grid  {
	if !PosIsValid(x,y)  {
		PosErrLog(x,y,"NewGrid")
		return nil
	}
	grid := new(Grid)
	grid.PosX = x
	grid.PosY = y
	return grid
}

//设置是否能移动
func (receiver *Grid) SetCanMove(canmove int8)  {
	receiver.canMove = canmove
}

//是否能移入该格子
func (receiver *Grid) CanMoveInto() bool {
	return receiver.canMove == 0
}

