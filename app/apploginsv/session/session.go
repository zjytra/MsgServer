/*
创建时间: 2020/12/20 0:10
作者: zjy
功能介绍:

*/

package session

import (
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/engine_core/network"
)



//逻辑层操作连接对象,加逻辑数据
type LsClientSession struct {
	PConn network.Conner
	PAcc *dbmodels.AccountT
	PRole *dbmodels.RoleT
}

