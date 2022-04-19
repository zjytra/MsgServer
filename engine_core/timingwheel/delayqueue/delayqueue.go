package delayqueue

import (
	"container/heap"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"sync"
	"sync/atomic"
	"time"
)

// The start of PriorityQueue implementation.
// Borrowed from https://github.com/nsqio/nsq/blob/master/internal/pqueue/pqueue.go

type Item struct {
	Value    interface{}
	Priority int64
	Index    int
}

// this is a priority queue as implemented by a min heap
// ie. the 0th element is the *lowest* value
type PriorityQueue []*Item

func NewPriorityQueue(capacity int) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	if pq.CheckIndex(i) || pq.CheckIndex(j) {
		return true
	}
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) CheckIndex(i int) bool {
	return  i < 0 || i >= pq.Len()
}

func (pq PriorityQueue) Swap(i, j int) {
	if pq.CheckIndex(i) || pq.CheckIndex(j) {
		return
	}
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq,item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	if n <= 0 {
		return nil
	}
	//取最后一个
	item := old[n-1]
	item.Index = -1
	//移除最后一个
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) PeekAndShift(max int64) (*Item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}

// The end of PriorityQueue implementation.

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C chan interface{}

	mu sync.RWMutex
	pq PriorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}
}

// New creates an instance of delayQueue with the specified size.
func New(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		pq:      NewPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
}

// Offer inserts the element into the current queue.
func (dq *DelayQueue) Offer(elem interface{}, expiration int64) {
	item := &Item{Value: elem, Priority: expiration}

	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mu.Unlock()

	if index == 0 {
		// A new Item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

// Poll starts an infinite loop, in which it continually waits for an element
// to expire and then send the expired element to the channel C.
func (dq *DelayQueue) Poll(exitC chan struct{}) {
	for {

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(timeutil.GetCurrentTimeMs())
		if item == nil {
			// No items left or at least one Item is pending.

			// We must ensure the atomicity of the whole operation, which is
			// composed of the above PeekAndShift and the following StoreInt32,
			// to avoid possible race conditions between Offer and Poll.
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 { //没有任务可以执行
				// No items left.
				select {
				case <-dq.wakeupC:
					// Wait until a new Item is added.
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 { //需要等待到延迟才能执行任务
				// At least one Item is pending.
				select {
				case <-dq.wakeupC:
					// A new Item with an "earlier" expiration than the current "earliest" one is added.
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// The current "earliest" Item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		select {
		case dq.C <- item.Value:
			// The expired element has been sent out successfully.
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}


func (dq *DelayQueue) Len() int{
	 dq.mu.RLock()
	 qlen := dq.pq.Len()
	 dq.mu.RUnlock()
	 return qlen
}

func (dq *DelayQueue) Clear()  []*Item{
	dq.mu.Lock()
	var items []*Item
	for dq.pq.Len() > 0 {
		item := heap.Pop(&dq.pq).(*Item)
		if item == nil {
			continue
		}
		items = append(items,item)
	}
	dq.mu.Unlock()
	return items

}