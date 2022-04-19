/*
创建时间: 2021/4/4 22:18
作者: zjy
功能介绍:

*/

package dispatch

import (
	"errors"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
)

//同步队列在会Start 会阻塞协程
type SyncQueue struct {
	eventQueue   chan Evener// 存放事件的的队列
	isClose      bool
}

func NewSyncQueue(queueSize int) *SyncQueue {
	qe := new(SyncQueue)
	qe.eventQueue = make(chan Evener,queueSize)
	return qe
}

func (this *SyncQueue) AddEvent(event Evener) error {
	if this.isClose { //关闭了就不能写了
		return errors.New( "添加事件:"+ event.EvenName()+ " 时QueueEvent已经关闭")
	}
	this.eventQueue <- event
	return nil
}
func (this *SyncQueue) Start(){
	this.doEvent()
}

func (this *SyncQueue) doEvent() {
	// 拉起错误避免宕机
	defer xlog.RecoverToLog(func() {
		this.doEvent()
	})
	for event :=  range this.eventQueue {
		// 当事件处理完并且没有数据
		qlen := len(this.eventQueue)
		if qlen > 100 {
			xlog.Debug( "SyncQueue.eventQueue剩余需要处理的数据 = %d", qlen)
		}
		if event == nil { //执行了关闭命令要检测下
			if this.canRelease() {
				break
			}
			continue
		}
		startT := timeutil.GetCurrentTimeMs()		//计算当前时间
		event.Execute()
		CheckTime("调度",startT,50)
		//任务执行完了需要检测下
	}

	xlog.Debug("SyncQueue 执行完毕 ")
}

//如果没有数据并且执行了关闭命令可以释放队列
func (this *SyncQueue) canRelease() bool  {
	if	len(this.eventQueue) == 0 && this.isClose {
		xlog.Debug("关闭")
		close(this.eventQueue)
		return true
	}
	return false
}

// 先关闭队列等待任务处理完
func (this *SyncQueue) Release() {
	if this.isClose {
		return
	}
	this.isClose = true
	this.eventQueue <- nil
	xlog.Debug( "还有未处理的事件%d", len(this.eventQueue))
}
