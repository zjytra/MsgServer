//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvparse"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strings"
)

var AttrCfgMap map[int]*AttrCfg

type  AttrCfg struct {
	Id int //#属性id 字段名称  id
	AttrName string //属性 字段名称  attrName
}

func SetAttrCfgMapData(csvpath  string ) {
   AttrCfgMap = loadAttrCfgCsv(csvpath)
}

func loadAttrCfgCsv(csvpath  string ) map[int]*AttrCfg{

	csvName := "/AttrCfg.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}
	tem := make(map[int]*AttrCfg)
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
		one := new(AttrCfg)
		
		done = csvparse.CheckType(one.Id, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Id = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.AttrName, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.AttrName = oneRow[col]
		col++
		
		if _,ok := tem[one.Id]; ok {
			fmt.Println(one.Id,"重复")
		}
		tem[one.Id] = one
	}
	return tem
}

func GetAttrCfgPtr(id int) *AttrCfg{
   if data, ok := AttrCfgMap[id]; ok {
		return data
	}
	return nil
}

func GetAllAttrCfg() map[int]*AttrCfg{
	return AttrCfgMap
}
