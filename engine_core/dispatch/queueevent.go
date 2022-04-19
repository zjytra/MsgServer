/*
创建时间: 2020/4/26
作者: zjy
功能介绍:
系统事件调度队列,跨gocroutinue
队列事件，接收各种事件进队，可以多个协程处理，也可以单个协程处理
*/

package dispatch

import (
	"errors"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/xlog"
	"sync"
)


//事件往调度系统塞的
type Evener interface {
	Execute()
	EvenName() string
}

//等待队列接口
type WaitQueue interface {
	AddEvent(event Evener) error
	Start()
	Release()
}


//异步队列多个协程消费
//当队列中没有数据时会阻塞协程
type AsyncQueue struct {
	eventQueue   chan Evener// 存放事件的的队列
	isClose      bool
	wg           sync.WaitGroup
	goroutineNum int  //执行协程的数量
}


func NewAsyncQueue(queueSize int,_goroutineNum int) *AsyncQueue {
	qe := new(AsyncQueue)
	qe.eventQueue = make(chan Evener,queueSize)
	qe.goroutineNum = _goroutineNum
	if qe.goroutineNum == 0 {
		xlog.Warning("创建队列服务器 没开协程")
		qe.goroutineNum = 1
	}
	return qe
}

func (this *AsyncQueue) Start(){
	for i := 0; i< 	this.goroutineNum ;i ++ {
		this.wg.Add(1)
		go this.doEvent()
	}
}

func (this *AsyncQueue) AddEvent(event Evener) error {
	if this.isClose { //关闭了就不能写了
		return errors.New( "添加事件:"+ event.EvenName()+ " 时QueueEvent已经关闭")
	}
	this.eventQueue <- event
	return nil
 }

func (this *AsyncQueue) doEvent() {
	// 拉起错误避免宕机
	defer xlog.RecoverToLog(func() {
		this.wg.Add(1)
		go	this.doEvent()
	})
	defer this.wg.Done()
	for event :=  range this.eventQueue {
		// 当事件处理完并且没有数据
		qlen := len(this.eventQueue)
		if qlen > 100 {
			xlog.Warning( "AsyncQueue.eventQueue剩余需要处理的数据 = %d", qlen)
		}
		if event == nil {
			if this.canRelease() {
				break
			}
			continue
		}
		startT := timeutil.GetCurrentTimeMs()		//计算当前时间
		event.Execute()
		CheckTime("调度",startT,50)
	}

	xlog.Debug(" AsyncQueue 执行完毕 ")
}

//判断是否能释放队列
func (this *AsyncQueue) canRelease() bool  {
	if	len(this.eventQueue) == 0 && this.isClose {
		xlog.Debug("关闭")
		close(this.eventQueue)
		return true
	}
	return false
}

//关闭队列服务 发个空数据避免队列中没有数据唤醒队列
func (this *AsyncQueue) Release() {
	if this.isClose {
		return
	}
	this.isClose = true
	//提醒一个协程关闭chan,当关闭chan 所有 for range chan 会退出
	this.eventQueue <- nil
	this.wg.Wait()
}
