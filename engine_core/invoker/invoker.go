/*
创建时间: 2020/12/20 1:18
作者: zjy
功能介绍:
中间调度，避免循环导包问题,不是线程安全的
*/

package invoker

import (
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"reflect"
	"strconv"
)

type CallBack func(interface{})

var (
	callMapping map[CallID][]CallBack
)

func init() {
	callMapping = make(map[CallID][]CallBack)
}

func RegisterInvoker(callid CallID, cb CallBack) {
	if cb == nil {
		return
	}
	handles, ok := callMapping[callid]
	if !ok { //没有找到
		var newHandles []CallBack
		newHandles = append(newHandles, cb)
		callMapping[callid] = newHandles
		return
	}
	for _, handle := range handles {
		val1 :=	reflect.ValueOf(handle)
		val2 := reflect.ValueOf(cb)
		if val1.Pointer() == val2.Pointer() {
			xlog.Debug("CallID %d  函数%v已经注册",callid,cb)
			return
		}
	}
	handles = append(handles, cb)
	callMapping[callid] = handles
}

func DelInvokerMapping(callid CallID, cb CallBack) {
	handles, ok := callMapping[callid]
	if !ok {
		return
	}
	for i, handle := range handles {
		val1 :=	reflect.ValueOf(handle)
		val2 := reflect.ValueOf(cb)
		if val1.Pointer() == val2.Pointer() {
			handles = append(handles[0:i], handles[i+1:]...)
			if len(handles) == 0 { //如果清空完了需要删除mapkey
				delete(callMapping,callid)
				return
			}
			callMapping[callid] = handles
		}
	}

}

func ReleaseallMapping() {
	callMapping = nil
}

func DoCall(callid CallID, dataPtr interface{})  {
	handles := getCallBack(callid)
	if handles == nil || len(handles) == 0 {
		xlog.Debug("CallID %d  为注册对应的处理",callid)
		return
	}
	for _, handle := range handles {
		startT := timeutil.GetCurrentTimeMs() //计算当前时间
		handle(dataPtr)                    //查看是否注册对应的处理函数
		dispatch.CheckTime("CallID = "+strconv.FormatInt(int64(callid), 10), startT, timeutil.FrameTimeMs)
	}
	return
}

func getCallBack(callid CallID) []CallBack {
	handles, ok := callMapping[callid]
	if !ok {
		xlog.Debug("callid = %d 未找到对应的方法",callid)
		return nil
	}
	return handles
}
