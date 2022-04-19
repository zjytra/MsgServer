//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvparse"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strings"
)

var DbCfgMap map[string]*DbCfg

type  DbCfg struct {
	Dbname string //#数据库名称 字段名称  dbname
	Ip string //ip地址 字段名称  ip
	Dbport string //端口号 字段名称  dbport
	Dbusername string //用户名 字段名称  dbusername
	Dbpwd string //密码 字段名称  dbpwd
	Maxopenconns int //最大链接数 字段名称  maxopenconns
	Maxidleconns int //闲置连接数 字段名称  maxidleconns
	Readnum int8 //读协程数量 字段名称  readnum
	Writenum int8 //写协程数量 字段名称  writenum
	Char_set string //字符集 字段名称  char_set
}

func SetDbCfgMapData(csvpath  string ) {
   DbCfgMap = loadDbCfgCsv(csvpath)
}

func loadDbCfgCsv(csvpath  string ) map[string]*DbCfg{

	csvName := "/DbCfg.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}
	tem := make(map[string]*DbCfg)
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
		one := new(DbCfg)
		
		done = csvparse.CheckType(one.Dbname, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Dbname = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Ip, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Ip = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Dbport, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Dbport = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Dbusername, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Dbusername = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Dbpwd, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Dbpwd = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Maxopenconns, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Maxopenconns = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Maxidleconns, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Maxidleconns = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Readnum, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Readnum = strutil.StrToInt8(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Writenum, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Writenum = strutil.StrToInt8(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Char_set, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Char_set = oneRow[col]
		col++
		
		if _,ok := tem[one.Dbname]; ok {
			fmt.Println(one.Dbname,"重复")
		}
		tem[one.Dbname] = one
	}
	return tem
}

func GetDbCfgPtr(dbname string) *DbCfg{
   if data, ok := DbCfgMap[dbname]; ok {
		return data
	}
	return nil
}

func GetAllDbCfg() map[string]*DbCfg{
	return DbCfgMap
}
