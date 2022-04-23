/*
创建时间: 2021/8/26 22:04
作者: zjy
功能介绍:

*/

package dbsys_test

import (
	"fmt"
	"github.com/zjytra/MsgServer/app/appdata"
	"github.com/zjytra/MsgServer/conf"
	"github.com/zjytra/MsgServer/csvsys/csvdata"
	"github.com/zjytra/MsgServer/dbmodels"
	"github.com/zjytra/MsgServer/engine_core/dbsys"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/timingwheel"
	"testing"
)

func TestMySqlDBStore_Execute(t *testing.T) {
	conf.InitConf()

	conf.SetAppPath("D:/work/WorldGame/Server/wengo/")
	csvdata.InitCsv()
	conf.StartConf()
	csvdata.StartCsv()
	appdata.InitAppDataByAppArgs(15)
	dispatch.InitMainQueue()
	dbsys.InitAccountDB()
	dbsys.InitGameDB()
	timingwheel.InitTimeWheel()

	//dbsys.GameAccountDB.RegisterTable(dbmodels.AccountT{})
	//dbsys.GameAccountDB.SyncBD()
	//dbsys.GameDB.RegisterTable(dbmodels.AccountLog{})
	//pTable := dbsys.GameDB.RegisterTable(dbmodels.RoleT{})
	//pTable.RegisterSubTable(dbmodels.SkillT{})
	//pTable.RegisterSubTable(dbmodels.ItemT{})
	//pTable.RegisterSubTable(dbmodels.MoneyT{})
	dbsys.GameDB.SyncBD()
	//dbsys.GameDB.StartTimer()
	//delayUpdate()
	dispatch.MainQueue.Start()
	//dbsys.GameDB.RegisterTable(dbmodels.RoleT{})
	//dbsys.GameDB.SyncBD()

}

func loadMore() {
	//Role := new(dbmodels.RoleT)
	//Role.OID.SetVal(2)
	//parm := new(dbsys.DBObjQueryParam)
	//parm.DbObj = Role
	//dbsys.GameDB.AsyncLoadSubTables(parm, OnLoadMoreSubFormCb)
}

func asyncUpdate() {
	//Role := new(dbmodels.RoleT)
	//Role.AccID.SetVal(5)
	//Role.OID.SetVal(197900535565647872)
	//Role.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	//Role.Name.SetVal("我爱你")
	//Role.Job.SetVal(1)
	//Role.Gender.SetVal(1)
	//Role.Lvl.SetVal(30)
	//Role.VipLv.SetVal(1)
	//
	//Role2 := new(dbmodels.RoleT)
	//Role2.AccID.SetVal(5)
	//Role2.OID.SetVal(197900535565647873)
	//Role2.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	//Role2.Name.SetVal("爱着你")
	//Role2.Job.SetVal(1)
	//Role2.Gender.SetVal(1)
	//Role2.Lvl.SetVal(20)
	//Role2.VipLv.SetVal(10)
	//Role.AccID.SetVal(2)
	//param := &dbsys.DBObjQueryParam{
	//	DbObj: Role,
	//}
	//dbsys.GameDB.AsyncLoadMoreObjs(param,OnLoadMoreFormCb)
	//writeParam := new(dbsys.DBObjWriteParam)
	//writeParam.DbObjs = append(writeParam.DbObjs,Role,Role2)
	//dbsys.GameDB.AsyncInsertObj(writeParam,nil)
}

func delayUpdate() {
	//Role := new(dbmodels.RoleT)
	//Role.AccID.SetVal(5)
	//Role.OID.SetVal(197900535565647872)
	//Role.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	//Role.Name.SetVal("我爱你")
	//Role.Job.SetVal(1)
	//Role.Gender.SetVal(1)
	//Role.Lvl.SetVal(30)
	//Role.VipLv.SetVal(1)
	//
	//startT := timeutil.GetCurrentTimeMs() //计算当前时间
	//writeParam := new(dbsys.DBObjWriteParam)
	//for i := 1; i <= 10000; i++ {
	//	money := new(dbmodels.MoneyT)
	//	money.RoleID.SetVal(10)
	//	money.MoneyID.SetVal(int16(i))
	//	money.MoneyVal.SetVal(int64(i))
	//	writeParam.AddObjs(money)
	//	//dbsys.GameDB.DelayInsert(Role2)
	//}
	//dispatch.CheckTime("构建数据 :", startT, 200)
	//dbsys.GameDB.AsyncInsertObj(writeParam,nil)

	//update := new(dbsys.DBObjWriteParam)
	//for i := 1; i <= 10000; i++ {
	//	money := new(dbmodels.MoneyT)
	//	money.RoleID.SetVal(10)
	//	money.MoneyID.SetVal(int16(i))
	//	money.MoneyVal.SetVal(int64(90))
	//	update.AddObjs(money)
	//}
	//dbsys.GameDB.AsyncDelObj(update,nil)


	//newWrite := new(dbsys.DBObjWriteParam)
	//for i := 1; i < 10000; i++ {
	//	Role2 := new(dbmodels.RoleT)
	//	Role2.AccID.SetVal(5)
	//	Role2.OID.SetVal(snowflake.GUID.NextId())
	//	Role2.CreateTime.SetVal(timeutil.GetCurrentTimeMs())
	//	Role2.Name.SetVal(fmt.Sprintf("qqqq%d",i))
	//	Role2.Job.SetVal(1)
	//	Role2.Gender.SetVal(1)
	//	Role2.Lvl.SetVal(20)
	//	Role2.VipLv.SetVal(10)
	//	newWrite.AddObjs(Role2)
	//	//dbsys.GameDB.DelayInsert(Role2)
	//}
	//dbsys.GameDB.AsyncUpdateObj(newWrite,nil)


	//newWrite := new(dbsys.DBObjWriteParam)
	//for i := 1; i < 100000; i++ {
	//	Role2 := new(dbmodels.AccountLog)
	//	Role2.AccID.SetVal(5)
	//	Role2.LoginName.SetVal("小明")
	//	Role2.LoginTime.SetVal(timeutil.GetCurrentTimeMs())
	//	Role2.LoginIp.SetVal("192.168.1.14")
	//	Role2.LoginMacAddr.SetVal("192.168.1.14")
	//	newWrite.AddObjs(Role2)
	//	//dbsys.GameDB.DelayInsert(Role2)
	//}
	//dbsys.GameDB.AsyncInsertObj(newWrite,nil)


}

func OnLoadFormCb (result dbsys.DBErrorCode,dbObj dbsys.DBObJer,) {

	//pObj,objOk := dbObj.(*dbmodels.AccountT)
	//if !objOk {
	//	return
	//}
	//
	//reqMsg,pbOk := param.(*protomsg.C2L_LoginMsg)
	//if !pbOk {
	//	xlog.Error("转换C2L_LoginMsg失败")
	//	return
	//}
	//
	//if result == dbsys.NODATA {
	//	//没有数据
	//}else if result == dbsys.DBSQLERRO {
	//}else if  result == dbsys.DBSUCCESS {
	//
	//	if strings.Compare(pObj.GetPwd(),reqMsg.Password) != 0  {
	//		//密码错误
	//		return
	//	}
	//	//存到redis 方便网关使用
	//	erro := dbsys.RedisCli.HSet(dbsys.CtxBg,"LoginAccount",strconv.FormatInt(pObj.GetUID(),10),reqMsg.Username).Err()
	//	if erro != nil {
	//		return
	//	}
	//	//登录成功
	//
	//}
}


func OnLoadMoreFormCb (result dbsys.DBErrorCode,param *dbsys.DBObjQueryParam,MoreObj []dbsys.DBObJer) {

	pObj,objOk := param.DbObj.(*dbmodels.RoleT)
	if !objOk {
		return
	}
	fmt.Printf("%v",pObj)
	if result == dbsys.NODATA {
		//没有数据
	}else if result == dbsys.DBSQLERRO {
	}else if  result == dbsys.DBSUCCESS {
		//for _, jer := range MoreObj {
		//	OneObj,oneOk := jer.(*dbmodels.RoleT)
		//	if oneOk {
		//		xlog.Debug("uid =%d,accid %d,名字 %s 等级 %d",OneObj.GetUID(),OneObj.AccID.GetVal(),OneObj.Name.GetVal(),OneObj.Lvl.GetVal())
		//	}
		//
		//}

	}
}


func OnLoadMoreSubFormCb (result dbsys.DBErrorCode,param *dbsys.DBObjQueryParam,objArr map[int][]dbsys.DBObJer) {

	//pObj,objOk := param.DbObj.(*dbmodels.RoleT)
	//if !objOk {
	//	return
	//}
	//xlog.Debug("%v",pObj)
	//if result == dbsys.NODATA {
	//	//没有数据
	//}else if result == dbsys.DBSQLERRO {
	//}else if  result == dbsys.DBSUCCESS {
	//    skillObj,ok := objArr[1]
	//	if ok {
	//		//技能
	//		for _, jer := range skillObj {
	//			pSkill,ok := jer.(*dbmodels.SkillT)
	//			if ok {
	//				xlog.Debug("pSkill uid = %d, cid = %d  lv = %d",pSkill.GetUID(),pSkill.SkillCid.GetVal(),pSkill.SkillLv.GetVal())
	//			}
	//		}
	//	}
	//
	//	items,ok := objArr[2]
	//	if ok {
	//		//技能
	//		for _, jer := range items {
	//			pitem,itemok := jer.(*dbmodels.ItemT)
	//			if itemok {
	//				xlog.Debug("item roleid %d uid = %d, cid = %d num =%d",pitem.RoleID.GetVal(),pitem.GetUID(),pitem.CID.GetVal(),pitem.Num.GetVal())
	//			}
	//		}
	//	}
	//
	//	moneys,ok := objArr[3]
	//	if ok {
	//		//技能
	//		for _, jer := range moneys {
	//			pmoney,mok := jer.(*dbmodels.MoneyT)
	//			if mok {
	//				xlog.Debug("pmoney uid = %d, MoneyID = %d ,num = %d",pmoney.RoleID.GetVal(),pmoney.MoneyID.GetVal(),pmoney.MoneyVal.GetVal())
	//			}
	//		}
	//	}
	//
	//}
}