//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvparse"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strings"
)

var MapCfgMap map[int]*MapCfg

type  MapCfg struct {
	Id int //#地图id 字段名称  id
	Mapname string //地图名字 字段名称  mapname
	Navmeshres string //寻路资源 字段名称  navmeshres
}

func SetMapCfgMapData(csvpath  string ) {
   MapCfgMap = loadMapCfgCsv(csvpath)
}

func loadMapCfgCsv(csvpath  string ) map[int]*MapCfg{

	csvName := "/MapCfg.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}
	tem := make(map[int]*MapCfg)
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
		one := new(MapCfg)
		
		done = csvparse.CheckType(one.Id, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Id = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Mapname, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Mapname = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Navmeshres, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Navmeshres = oneRow[col]
		col++
		
		if _,ok := tem[one.Id]; ok {
			fmt.Println(one.Id,"重复")
		}
		tem[one.Id] = one
	}
	return tem
}

func GetMapCfgPtr(id int) *MapCfg{
   if data, ok := MapCfgMap[id]; ok {
		return data
	}
	return nil
}

func GetAllMapCfg() map[int]*MapCfg{
	return MapCfgMap
}
