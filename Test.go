/*
创建时间: 2020/3/29
作者: zjy
功能介绍:

*/

package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"math/rand"
	_ "net/http/pprof"
	"reflect"
	"time"
)

//import (
//	"container/list"
//	"database/sql"
//	"encoding/binary"
//	"fmt"
//	"github.com/zjytra/devlop/xcontainer/queue"
//	"github.com/zjytra/MsgServer/devlop/xutil"
//	"github.com/zjytra/MsgServer/csvsys/csvdata"
//	"github.com/zjytra/MsgServer/dbmodels"
//	"github.com/zjytra/MsgServer/engine_core/dbsys"
//	"github.com/zjytra/MsgServer/engine_core/network"
//	"github.com/zjytra/MsgServer/engine_core/timingwheel"
//	"math"
//	"reflect"
//	"strings"
//	"sync"
//	"unsafe"
//
//	// "sync/atomic"
//	"time"
//)
//
//var locker = new(sync.Mutex)
//var cond = sync.NewCond(locker)
//var queue1 *queue.Queue
//var wg sync.WaitGroup
//
//var pool sync.Pool
//
//type Person struct {
//	name string
//	age  int
//}
//
//var i int
//
//func Test(t *timingwheel.Timer) {
//	fmt.Println("%p  ", t, time.Now())
//	if i == 5 {
//		t.Stop()
//		fmt.Println("Test stop")
//	}
//	i++
//}
//
//// param |加密key| 主命令 | 字命令 | datalen | Account |
//func UnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
//	msglen := len(readb)
//	if readb == nil || msglen == 0 {
//		return
//	}
//	var start, end int
//	start, end = GetNextIndex(end, maincmd)
//	maincmd = binary.LittleEndian.Uint16(readb[start:end])
//	start, end = GetNextIndex(end, subcmd)
//	subcmd = binary.LittleEndian.Uint16(readb[start:end])
//	var datalen uint32
//	start, end = GetNextIndex(end, datalen)
//	datalen = binary.LittleEndian.Uint32(readb[start:end])
//	msg = make([]byte, datalen)
//	copy(msg, readb[end:])
//	return
//}
//
//func GetNextIndex(end int, Account interface{}) (head, tail int) {
//	head = end                            //尾部变成头部
//	tail = head + xutil.IntDataSize(Account) //新的尾部=头加上数据的长度
//	return
//}
//
//// 打单包
//func PackOne(maincmd uint32, msg []byte) ([]byte, error) {
//	msglen := uint32(len(msg))
//	// 加密key
//	var alllen int
//	alllen += xutil.IntDataSize(maincmd)
//	alllen += xutil.IntDataSize(subcmd)
//	alllen += xutil.IntDataSize(msglen)
//	alllen += int(msglen)
//	writeBuf := make([]byte, alllen)
//	var start, end int
//	start, end = GetNextIndex(end, maincmd)
//	binary.LittleEndian.PutUint16(writeBuf[start:end], maincmd)
//	start, end = GetNextIndex(end, subcmd)
//	binary.LittleEndian.PutUint16(writeBuf[start:end], subcmd)
//	start, end = GetNextIndex(end, msglen)
//	binary.LittleEndian.PutUint32(writeBuf[start:end], msglen)
//	copy(writeBuf[end:], msg)
//	return writeBuf, nil
//}
//
//type RowStringMap gamemap[string]interface{}
//
//func testQuery(rows *sql.Rows, to interface{}) {
//	acc := dbsys.RowsToStructSlice(rows, to)
//	for _, Account := range acc {
//		ptr := Account.(*dbmodels.Accounts)
//		fmt.Println(ptr)
//	}
//}
//

//
//type Users struct {
//	adas []byte
//	Age  int32
//}
//
//type Call func()
//
//func TestBuf(pointer unsafe.Pointer) {
//	if pointer == nil {
//		fmt.Println(pointer)
//		return
//	}
//	buf := make([]byte, 1024)
//	netbuf := network.NewNetBuffer(buf, 0)
//	netbuf.WriteInt64(int64(9223372036854775707))
//	netbuf.WriteUint64(uint64(922337203685477000))
//	netbuf.WriteInt32(int32(-2099555222))
//	netbuf.WriteUint32(math.MaxUint32 - 1)
//	netbuf.WriteInt16(-32768)
//	netbuf.WriteUint16(65535)
//	netbuf.WriteInt8(int8(-120))
//	netbuf.WriteUInt8(255)
//	netbuf.WriteString("中国")
//	netbuf.WriteFloat(999999.1234567)
//	netbuf.WriteDouble(1.123123123123123 + 1.123123)
//	fmt.Println(netbuf.ReadInt64(), netbuf.ReadUint64(), netbuf.ReadInt32(), netbuf.ReadUint32())
//	fmt.Println(netbuf.ReadInt16())
//	fmt.Println(netbuf.ReadUint16())
//	fmt.Println(netbuf.ReadInt8())
//	fmt.Println(netbuf.ReadUint8())
//	fmt.Println(netbuf.ReadBytes())
//	fmt.Println(netbuf.ReadFloat())
//	fmt.Println(netbuf.ReadDouble())
//}
//
////将str
//func StrToByteArrByPtr(str *string) []byte {
//	pointer := unsafe.Pointer(str)
//	return *(*[]byte)(pointer)
//}
//
////将byte数组转换为string
//func ByteArrToStr(Account []byte) string {
//	pointer := unsafe.Pointer(&Account[0])
//	return *(*string)(pointer)
//}
//
//type Temp struct {
//	name string
//}
//
//func (this *Temp)PrintName()  {
//	fmt.Println(this.name)
//}
//
//func TestPT(){
//	fmt.Println("Test")
//}
//
//type Caller struct {
//	ptr unsafe.Pointer
//	cb func()
//}
//
//
//
//
//func TestPtr(cb func())  {
//	temmmm := unsafe.Pointer(&cb)
//	cal := *(*func())(temmmm)
//	cal()
//}
//
//
//
//func TestList() {
//	//初始化一个list
//	l := list.New()
//	l.PushBack(1)
//	l.PushBack(2)
//	l.PushBack(3)
//	l.PushBack(4)
//
//	fmt.Println("Before Removing...")
//	//遍历list，删除元素
//	var n *list.Element
//	for e := l.Front(); e != nil; e = n {
//		fmt.Println("removing", e.Value)
//		n = e.Next()
//		l.Remove(e)
//	}
//	fmt.Println("After Removing...")
//	//遍历删除完元素后的list
//	for e := l.Front(); e != nil; e = e.Next() {
//		fmt.Println(e.Value)
//	}
//}
//func Test2(p interface{}) {
//	v := reflect.ValueOf(p)
//	typ := reflect.TypeOf(p)
//	elmv := reflect.Indirect(v)
//	fmt.Println("name", elmv.Type().Name())
//
//	allname := typ.String()
//	splitname := strings.Split(allname, ".")
//	if len(splitname) > 1 {
//		fmt.Println("name2", splitname[1])
//	}
//
//	tagv := v.Elem()
//	fmt.Println(v.Kind(), v.Elem().Kind())
//	for i := 0; i < tagv.NumField(); i++ {
//		feildval := tagv.Field(i)
//		fmt.Println(feildval.Type().Name())
//	}
//}
//
//func TestTaskPool() {
//	for {
//		fmt.Println("TestTaskPool")
//		time.Sleep(time.Second)
//	}
//
//}
//func TestTaskPool2() {
//	fmt.Println("TestTaskPool2")
//}
//
//func TestPool() {
//	pool.New = func() interface{} {
//		return new(Person)
//	}
//	Account := pool.Get()
//	fmt.Println(Account)
//}
//
//func TestCond() {
//	queue1 = queue.NewQueue()
//	_, err := queue1.PopFront()
//	if err != nil {
//		fmt.Println(err)
//	}
//	wg.Add(20)
//	// 10个消费
//	for i := 0; i < 10; i++ {
//		go func(x int) {
//			defer wg.Done()
//			cond.L.Lock() // 获取锁
//			for queue1.Len() == 0 {
//				cond.Wait() // 等待通知，阻塞当前 goroutine
//			}
//			cond.L.Unlock() // 释放锁
//			val, erro := queue1.PopFront()
//			if erro != nil {
//				fmt.Println(erro)
//				return
//			}
//			// do something. 这里仅打印
//			fmt.Println("队列的值 ", val, "type=", reflect.TypeOf(val))
//		}(i)
//	}
//	for i := 0; i < 10; i++ {
//		go func(x int) {
//			defer wg.Done()
//			queue1.PushBack(x)
//			cond.Signal() // 通知其他线程
//		}(i)
//	}
//
//	wg.Wait()
//	fmt.Printf("end")
//}

func TestParam(arr ...interface{}) {
	fmt.Println(arr)
	t := reflect.TypeOf(arr)
	fmt.Println(t.Kind())
}

func DBTest() {
	//csvdata.SetDbCfgMapData("./csv")
	//conf := csvdata.GetDbCfgPtr("accountdb")
	//datasource := dbsys.GetMysqlDataSourceName(conf)
	//db, Erro := sql.Open("mysql", datasource)
	//if Erro != nil {
	//	fmt.Printf("%v \n", Erro)
	//	return
	//}
	//res, Erro := db.Query("select SCHEMA_NAME from information_schema.SCHEMATA where SCHEMA_NAME = ?; ", "accountdb")
	//if Erro != nil {
	//	fmt.Println(Erro)
	//	return
	//}
	//dbs := dbsys.RowToMap(res)
	//if dbs == nil {
	//	return
	//}
	//hasName, ok := dbs["SCHEMA_NAME"]
	//if ok && hasName == "accountdb" {
	//	return
	//}
	//if Erro != nil {
	//	fmt.Println(Erro)
	//	return
	//}
	//fmt.Println(res)
	//pa := GetParam(1)
	//rows, erro := db.Query("SELECT * FROM RoleT WHERE accid = ?; ",pa...)
	//if erro != nil {
	//	fmt.Printf("%v \n", erro)
	//	return
	//}
	//SELECT * FROM SkillT WHERE roleid = 1;
	// pacc := dbsys.RowsToStructSlice(rows,reflect.TypeOf(&dbmodels.Accounts{}))
	//var arr []interface{}
	//arr := dbsys.RowToStruct(rows,dbmodels.RoleT{})
	//Account := maps[4]
	//if Account == nil {
	//	fmt.Println("第四个为null")
	//}
	//fmt.Println(arr)
}

func GetParam(pram ...interface{}) []interface{} {
	return pram
}
func TestLua() {
	LuaState := lua.NewState()
	defer LuaState.Close()
	err := LuaState.DoFile("LuaScripts/main.lua")
	Lfunc, ok := LuaState.GetGlobal("AAA").(*lua.LFunction)
	if ok {
		erro := LuaState.PCall(0, lua.MultRet, Lfunc)
		if erro != nil {
			fmt.Println(erro)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	if err := LuaState.DoString(`print("hello")`); err != nil {
		panic(err)
	}
}

type TimerTask struct {
	timeMS int
	id     int
}

func (this *TimerTask) DoTask() {
	fmt.Println(this.id)
}

func DoTask() error {
	fmt.Println("DoTask", rand.Int())
	return nil
}
func DoTask2() error {
	fmt.Println("DoTask2")
	return nil
}


type EveryScheduler struct {
	Interval time.Duration
}

func (s *EveryScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

func main() {

	fmt.Println(timeutil.GetTimeNow())
	//gamemap.GameMapInit()
	//for i := 0; i < 50; i++ {
	//	startTime := timeutil.GetCurrentTimeMs()
	//	arr := gamemap.AStar.FindPathByInt16(1, 2, 1900, 1900)
	//	dispatch.CheckTime("寻路", startTime, 10)
	//	if arr != nil {
	//		xlog.Debug("find ok")
	//		//for _, d := range arr {
	//		//	fmt.Printf("[%d,%d] \n",d.PosX,d.PosY)
	//		//}
	//	}
	//	arr = nil
	//}

	//TestReflect()
}

type TestEvent struct {
}

func (this *TestEvent) Execute() {
	fmt.Println(this.EvenName())
}

func (this *TestEvent) EvenName() string {
	return "TestEvent"
}
