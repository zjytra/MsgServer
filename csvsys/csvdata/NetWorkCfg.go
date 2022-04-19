//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/MsgServer/csvsys/csvparse"
	"github.com/zjytra/MsgServer/devlop/xutil/strutil"
	"strconv"
	"strings"
)

var NetWorkCfgMap map[int32]*NetWorkCfg

type  NetWorkCfg struct {
	App_id int32 //#服务器id 字段名称  app_id
	App_kind int32 //服务器类型1.客户端2.聊天服3.登录服4.网关5.数据中6.游戏逻辑服7.监控服,监控所有的服务器，管理服务器之间的连接 字段名称  app_kind
	Group int32 //服务器分组，好做分区 字段名称  group
	App_name string //服务器名称 字段名称  app_name
	Out_addr string //外部连接的地址 字段名称  out_addr
	Inner_addr string //内部连接地址 字段名称  inner_addr
	Out_prot int32 //外部连接端口 字段名称  out_prot
	Max_connect int32 //最大连接数 字段名称  max_connect
	Goroutines_size int //go协程池数量连接数的两倍多一点暂时没用 字段名称  goroutines_size
	Max_msglen int32 //消息最大长度 字段名称  max_msglen
	Write_cap_num int //连接写的包队列大小 字段名称  write_cap_num
	Checklink_s int //检查连接存活时间间隔秒 字段名称  checklink_s
	Max_rec_msg_ps int //每秒最大收包数量 字段名称  max_rec_msg_ps
	Http_port int //http端口 字段名称  http_port
	Msg_isencrypt bool //消息是否加密 字段名称  msg_isencrypt
	Dbname string //关联的数据库 字段名称  dbname
	Logdbname string //关联的日志库 字段名称  logdbname
}

func SetNetWorkCfgMapData(csvpath  string ) {
   NetWorkCfgMap = loadNetWorkCfgCsv(csvpath)
}

func loadNetWorkCfgCsv(csvpath  string ) map[int32]*NetWorkCfg{

	csvName := "/NetWorkCfg.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}
	tem := make(map[int32]*NetWorkCfg)
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
		one := new(NetWorkCfg)
		
		done = csvparse.CheckType(one.App_id, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.App_id = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.App_kind, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.App_kind = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Group, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Group = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.App_name, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.App_name = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Out_addr, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Out_addr = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Inner_addr, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Inner_addr = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Out_prot, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Out_prot = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Max_connect, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Max_connect = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Goroutines_size, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Goroutines_size = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Max_msglen, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Max_msglen = strutil.StrToInt32(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Write_cap_num, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Write_cap_num = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Checklink_s, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Checklink_s = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Max_rec_msg_ps, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Max_rec_msg_ps = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Http_port, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Http_port = strutil.StrToInt(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Msg_isencrypt, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Msg_isencrypt ,_ = strconv.ParseBool(oneRow[col])
		col++
		
		done = csvparse.CheckType(one.Dbname, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Dbname = oneRow[col]
		col++
		
		done = csvparse.CheckType(one.Logdbname, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.Logdbname = oneRow[col]
		col++
		
		if _,ok := tem[one.App_id]; ok {
			fmt.Println(one.App_id,"重复")
		}
		tem[one.App_id] = one
	}
	return tem
}

func GetNetWorkCfgPtr(app_id int32) *NetWorkCfg{
   if data, ok := NetWorkCfgMap[app_id]; ok {
		return data
	}
	return nil
}

func GetAllNetWorkCfg() map[int32]*NetWorkCfg{
	return NetWorkCfgMap
}
