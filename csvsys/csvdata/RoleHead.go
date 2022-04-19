//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvparse"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strings"
)

var RoleHeadMap map[int]*RoleHead

type  RoleHead struct {
	Id int //#头像id 字段名称  id
	RoleHead string //头像资源 字段名称  roleHead
}

func SetRoleHeadMapData(csvpath  string ) {
   RoleHeadMap = loadRoleHeadCsv(csvpath)
}

func loadRoleHeadCsv(csvpath  string ) map[int]*RoleHead{

	csvName := "/RoleHead.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}
	tem := make(map[int]*RoleHead)
	nameRow  := csvmapdata[0]
	typeRow  := csvmapdata[1]
	var col int
    var done bool
	for rowNum, oneRow := range csvmapdata {
		if rowNum < csvparse.Invalid_Row { // 排除前三行
			continue
		}
		col = 0 //重置变量
		//第一个是#的字符行忽略掉
		if strings.Index(oneRow[col],"#") == 0 {
			continue
		}
		one := new(RoleHead)
		
		done = csvparse.CheckType(one.Id, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Id = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.RoleHead, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.RoleHead = oneRow[col]
		col++
		
		if _,ok := tem[one.Id]; ok {
			fmt.Println(one.Id,"重复")
		}
		tem[one.Id] = one
	}
	return tem
}

func GetRoleHeadPtr(id int) *RoleHead{
   if data, ok := RoleHeadMap[id]; ok {
		return data
	}
	return nil
}

func GetAllRoleHead() map[int]*RoleHead{
	return RoleHeadMap
}
