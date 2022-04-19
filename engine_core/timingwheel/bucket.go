package timingwheel

import (
	"container/list"
	"github.com/zjytra/MsgServer/engine_core/dispatch"
	"sync"
	"sync/atomic"
	"unsafe"
)


type TimeEvent func()
// Timer represents a single event. When the Timer expires, the given
// task will be executed.
type Timer struct {
	expiration int64 // in milliseconds
	task       TimeEvent

	// The bucket that holds the list to which this timer's pItem belongs.
	//
	// NOTE: This field may be updated and read concurrently,
	// through Timer.Stop() and Bucket.Flush().
	b unsafe.Pointer // type: *bucket

	// The timer's pItem.
	element *list.Element
    //执行队列,多半是主线程执行
	executeQueue dispatch.WaitQueue
}

func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}

// Stop prevents the Timer from firing. It returns true if the call
// stops the timer, false if the timer has already expired or been stopped.
//
// If the timer t has already expired and the t.task has been started in its own
// goroutine; Stop does not wait for t.task to complete before returning. If the caller
// needs to know whether t.task is completed, it must coordinate with t.task explicitly.
func (t *Timer) Stop() bool {
	stopped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		// If b.Remove is called just after the timing wheel's goroutine has:
		//     1. removed t from b (through b.Flush -> b.remove)
		//     2. moved t from b to another bucket ab (through b.Flush -> b.remove and ab.Add)
		// this may fail to remove t due to the change of t's bucket.
		stopped = b.Remove(t)

		// Thus, here we re-get t's possibly new bucket (nil for case 1, or ab (non-nil) for case 2),
		// and retry until the bucket becomes nil, which indicates that t has finally been removed.
	}
	return stopped
}

func (t *Timer) doTask() error {
	//如果有队列对象让队列对象执行
	if t.executeQueue != nil {
		return t.executeQueue.AddEvent(t)
	}

	go t.task()

	return nil
}

//队列执行任务需要
func (t *Timer) Execute() {
	t.task()
	return
}

func (t *Timer) EvenName() string {
	return "Timer"
}

type bucket struct {
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers do not ensure it. So we must keep the 64-bit field
	// as the first field of the struct.
	//
	// For more explanations, see https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	// and https://go101.org/article/memory-layout.html.
	expiration int64
	
	mu     sync.Mutex
	timers *list.List
}

func newBucket() *bucket {
	return &bucket{
		timers:     list.New(),
		expiration: -1,
	}
}

func (b *bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

func (b *bucket) SetExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}

func (b *bucket) Add(t *Timer) {
	b.mu.Lock()

	e := b.timers.PushBack(t)
	t.setBucket(b)
	t.element = e

	b.mu.Unlock()
}

func (b *bucket) remove(t *Timer) bool {
	if t.getBucket() != b {
		// If remove is called from t.Stop, and this happens just after the timing wheel's goroutine has:
		//     1. removed t from b (through b.Flush -> b.remove)
		//     2. moved t from b to another bucket ab (through b.Flush -> b.remove and ab.Add)
		// then t.getBucket will return nil for case 1, or ab (non-nil) for case 2.
		// In either case, the returned value does not equal to b.
		return false
	}
	b.timers.Remove(t.element)
	t.setBucket(nil)
	t.element = nil
	return true
}

func (b *bucket) Remove(t *Timer) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.remove(t)
}

func (b *bucket) Flush(reinsert func(*Timer)) {
	var ts []*Timer
	b.mu.Lock()
	for e := b.timers.Front(); e != nil; {
		next := e.Next()
		t := e.Value.(*Timer)
		b.remove(t)
		ts = append(ts, t)
		e = next
	}
	b.SetExpiration(-1) // TODO: Improve the coordination with b.Add()
	b.mu.Unlock()
	//从最底层的时间轮重新插入
	for _, t := range ts {
		reinsert(t)
	}
}


func (b *bucket) Clear() {
	var ts []*Timer
	b.mu.Lock()
	for e := b.timers.Front(); e != nil; {
		next := e.Next()
		t := e.Value.(*Timer)
		ts = append(ts, t)
		e = next
	}
	b.mu.Unlock()

	for _, t := range ts {
		t.Stop()
	}

}