package timingwheel

import (
	"errors"
	"fmt"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"github.com/zjytra/MsgServer/engine_core/timingwheel/delayqueue"
	"sync"
	"time"
)

// TimingWheel is an implementation of Hierarchical Timing Wheels.
type TimingWheel struct {
	tick      int64 // in milliseconds
	wheelSize int64

	interval    int64 // in milliseconds
	currentTime int64 // in milliseconds
	buckets     []*bucket
	queue       *delayqueue.DelayQueue

	// The higher-level overflow wheel.
	//
	// NOTE: This field may be updated and read concurrently, through Add().
	overflowWheel *TimingWheel // type: *TimingWheel


	isClose bool

	exitC     chan struct{}
	waitGroup 	sync.WaitGroup
	//保证只用低层的队列
	addQueue chan *Timer
	delQueue chan *Timer
}
//默认时间轮
var defWheel *TimingWheel

func InitTimeWheel() {
	defWheel = NewTimingWheel(time.Millisecond,20)
	defWheel.Start()
}

// NewTimingWheel creates an instance of TimingWheel with the given tick and wheelSize.
func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be greater than or equal to 1ms"))
	}

	startMs := timeutil.GetCurrentTimeMs()

	tw := newTimingWheel(
		tickMs,
		wheelSize,
		startMs,
		delayqueue.New(int(wheelSize)),
		make(chan *Timer,wheelSize),
	)
	tw.delQueue = 	make(chan *Timer,wheelSize)
	return tw
}

// newTimingWheel is an internal helper function that really creates an instance of TimingWheel.
func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, queue *delayqueue.DelayQueue,addQueue chan *Timer) *TimingWheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	return &TimingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * wheelSize,
		buckets:     buckets,
		queue:       queue,
		exitC:       make(chan struct{}),
		addQueue:  addQueue,
	}
}

// add inserts the timer t into the current timing wheel.
func (tw *TimingWheel) add(t *Timer) bool {
	currentTime := tw.currentTime
	if t.expiration < currentTime+tw.tick {
		// Already expired
		return false
	} else if t.expiration < currentTime+tw.interval {
		// Put it into its own bucket
		virtualID := t.expiration / tw.tick
		b := tw.buckets[virtualID%tw.wheelSize]
		b.Add(t)
		// Set the bucket expiration time
		if b.SetExpiration(virtualID * tw.tick) {
			// The bucket needs to be enqueued since it was an expired bucket.
			// We only need to enqueue the bucket when its expiration time has changed,
			// i.e. the wheel has advanced and this bucket get reused with a new expiration.
			// Any further calls to set the expiration within the same wheel cycle will
			// pass in the same value and hence return false, thus the bucket with the
			// same expiration will not be enqueued multiple times.
			tw.queue.Offer(b, b.Expiration())
		}
		return true
	} else {
		// Out of the interval. Put it into the overflow wheel
		if tw.overflowWheel == nil {
			//tw.queue 这里是传的最低层的队列进去,所有的桶进的是一个队列
			tw.overflowWheel = newTimingWheel(
				tw.interval,
				tw.wheelSize,
				currentTime,
				tw.queue,
				tw.addQueue,
			)
		}
		return tw.overflowWheel.add(t)
	}
}

// addOrRun inserts the timer t into the current timing wheel, or run the
// timer's task if it has already expired.
func (tw *TimingWheel) addOrRun(t *Timer) {
	if tw.isClose {
		return
	}
	if !tw.add(t) {
		// Already expired
		//t.doTask()
		//如果有队列对象让队列对象执行
		if t.executeQueue != nil {
			 t.executeQueue.AddEvent(t)
			return
		}
		tw.waitGroup.Add(1)
		go func() {
			t.task()
			tw.waitGroup.Done()
		}()
	}
}

func (tw *TimingWheel) advanceClock(expiration int64) {
	currentTime := tw.currentTime
	if expiration >= currentTime+tw.tick {
		currentTime = truncate(expiration, tw.tick)
		tw.currentTime = currentTime
		// Try to advance the clock of the overflow wheel if present
		if tw.overflowWheel != nil {
			tw.overflowWheel.advanceClock(currentTime)
		}
	}
}

//只有最低层的时间轮在调
// Start starts the current timing wheel.
func (tw *TimingWheel) Start() {

	tw.waitGroup.Add(1)
	go func() {
		tw.queue.Poll(tw.exitC)
		tw.waitGroup.Done()
		fmt.Println("TimingWheel Poll end")
	}()


	tw.waitGroup.Add(1)
	go func() {
		for {
			select {
			case elem := <-tw.queue.C:
				b := elem.(*bucket)
				tw.advanceClock(b.Expiration())
				//每次都先从低层的加
				b.Flush(tw.addOrRun)
			case <-tw.exitC:
				tw.OnStop()
				close(tw.addQueue)
				close(tw.delQueue)
				tw.waitGroup.Done()
				fmt.Println("TimingWheel tw.exitC")
				return
			case  timer := <-tw.addQueue:
				tw.addOrRun(timer)
			case  timer := <-tw.delQueue:
				timer.Stop()
			}
		}
	}()
}

func (tw *TimingWheel) OnStop() {
	items := tw.queue.Clear()
	if items != nil {
		for _, item := range items {
			b,ok := item.Value.(*bucket)
			if ok {
				b.Clear()
			}
		}
	}
}

// Stop stops the current timing wheel.
//
// If there is any timer's task being running in its own goroutine, Stop does
// not wait for the task to complete before returning. If the caller needs to
// know whether the task is completed, it must coordinate with the task explicitly.
func (tw *TimingWheel) Stop() {
	tw.isClose = true
	close(tw.exitC)
	tw.waitGroup.Wait()
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
func (tw *TimingWheel) AfterFuncInQueue(d time.Duration, f TimeEvent,exeQueue dispatch.WaitQueue) *Timer {
	if tw.isClose {
		return nil
	}
	t := &Timer{
		expiration: timeutil.NowTimeAdd(d),
		task:       f,
		executeQueue:exeQueue,
	}

	//tw.addOrRun(t)
	tw.addQueue <- t
	return t
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
func (tw *TimingWheel) AfterFunc(d time.Duration, f TimeEvent) *Timer {
	return tw.AfterFuncInQueue(d,f,nil)
}

// Scheduler determines the execution plan of a task.
type Scheduler interface {
	// Next returns the next execution time after the given (previous) time.
	// It will return a zero time if no next time is scheduled.
	//
	// All times must be UTC.
	Next(time.Time) time.Time
}

// ScheduleFunc calls f (in its own goroutine) according to the execution
// plan scheduled by s. It returns a Timer that can be used to cancel the
// call using its Stop method.
//
// If the caller want to terminate the execution plan halfway, it must
// stop the timer and ensure that the timer is stopped actually, since in
// the current implementation, there is a gap between the expiring and the
// restarting of the timer. The wait time for ensuring is short since the
// gap is very small.
//
// Internally, ScheduleFunc will ask the first execution time (by calling
// s.Next()) initially, and create a timer if the execution time is non-zero.
// Afterwards, it will ask the next execution time each time f is about to
// be executed, and f will be called at the next execution time if the time
// is non-zero.
func (tw *TimingWheel) ScheduleFuncInQueue(s Scheduler, f TimeEvent,exeQueue dispatch.WaitQueue) (t *Timer) {
	if tw.isClose {
		return nil
	}
	expiration := s.Next(timeutil.GetTimeNow())
	if expiration.IsZero() {
		// No time is scheduled, return nil.
		return
	}

	t = &Timer{
		expiration: timeutil.TimeToMs(expiration),
		task: func()  {
			if tw.isClose {
				return
			}
			// Schedule the task to execute at the next time if possible.
			expi := s.Next(timeutil.MsToTime(t.expiration))
			if !expi.IsZero() {
				t.expiration = timeutil.TimeToMs(expi)
				//tw.addOrRun(t)
				tw.addQueue <- t
			}
			// Actually execute the task.
		    f()
		},
		executeQueue: exeQueue,
	}
	//tw.addOrRun(t)
	tw.addQueue <- t
	return t
}

func (tw *TimingWheel) ScheduleFunc(s Scheduler, f TimeEvent) (t *Timer){
	return tw.ScheduleFuncInQueue(s,f,nil)
}

//解决并发添加与删除问题
func (tw *TimingWheel) StopTimer(t *Timer) {
	if tw.isClose {
		return
	}
	tw.delQueue <- t
}


func AfterFunc(d time.Duration, f TimeEvent) *Timer {
	return 	AfterFuncInQueue(d,f,nil)
}

func ScheduleFunc(s Scheduler, f TimeEvent)  *Timer {
	return 	ScheduleFuncInQueue(s,f,nil)
}

//向事件队列中添加定时任务
func AfterFuncInQueue(d time.Duration, f TimeEvent,exeQueue dispatch.WaitQueue) *Timer {
	if defWheel == nil {
		InitTimeWheel()
		return nil
	}
	return defWheel.AfterFuncInQueue(d,f,exeQueue)
}

//向事件队列中添加定时任务
func ScheduleFuncInQueue(s Scheduler, f TimeEvent,exeQueue dispatch.WaitQueue)  *Timer {
	if defWheel == nil {
		InitTimeWheel()
		return nil
	}
	return 	defWheel.ScheduleFuncInQueue(s,f,exeQueue)
}


func Stop() {
	if defWheel == nil {
		return
	}
	defWheel.Stop()
}


func StopTimer(t *Timer) {
	if defWheel == nil {
		return
	}
	defWheel.StopTimer(t)
}


