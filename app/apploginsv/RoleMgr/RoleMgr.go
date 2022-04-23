/*
创建时间: 2022/4/19 16:54
作者: zjy
功能介绍:
玩家管理相关
*/

package RoleMgr

import (
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
)

var(
	idRoles  map[int64]*dbmodels.RoleT   //根据id保存
	nameRoles map[string]*dbmodels.RoleT   //根据username保存
	asyncCreate  map[string]int64          //检测异步注册
)




func Init() {
	idRoles = make(map[int64]*dbmodels.RoleT)
	nameRoles = make(map[string]*dbmodels.RoleT)
	asyncCreate = make(map[string]int64)
}

func AddRole(pRole *dbmodels.RoleT)  {
	if pRole == nil {
		return
	}
	idRoles[pRole.RoleID.GetVal()] = pRole
	nameRoles[pRole.RoleName.GetVal()] = pRole
}

func Remove(pRole *dbmodels.RoleT)  {
	if pRole == nil {
		return
	}
	delete(idRoles,pRole.RoleID.GetVal())
	delete(nameRoles,pRole.RoleName.GetVal())
}


//查看是否正在登录中防止客户端重复发请求
func IsAsyncCreate(roleName string) bool {
	_, isOk := 	asyncCreate[roleName]
	if !isOk {
		return false
	}
	return true
}
//设置正在登录的标志
func SetAsyncCreate(roleName string) {
	//保证
	asyncCreate[roleName] = timeutil.GetCurrentTimeS()
}
//删除数据
func DelAsyncCreate(roleName string) {
	delete(asyncCreate,roleName)
}


func GetNameRole(roleName string)*dbmodels.RoleT{
 	role,ok := nameRoles[roleName]
	if ok {
		return role
	}
	return nil
}