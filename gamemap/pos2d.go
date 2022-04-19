/*
创建时间: 2021/7/11 15:50
作者: zjy
功能介绍:

*/


package gamemap

type PosType int16

func (receiver PosType) ToInt() int16  {
	return int16(receiver)
}


type Pos2D struct {
	PosX PosType
	PosY PosType
}


func NewPos2D(x PosType,y PosType)*Pos2D {
	if !PosIsValid(x,y) {
		return nil
	}
	p2d := new(Pos2D)
	p2d.PosX = x
	p2d.PosY = y
	return p2d
}