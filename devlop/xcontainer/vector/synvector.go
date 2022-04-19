package vector

import "sync"

// Synchronous FIFO queue
type SyncVector struct {
	lock    sync.Mutex  //锁
	popable *sync.Cond  //条件变量
	buffer  *Vector      // 数据buff
	closed  bool
}

// Create a new SyncVector
func NewSyncVector(cap int) *SyncVector {
	ch := &SyncVector{
		buffer: NewVector(cap),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

// 提取队列头变量
// 注意:当队列中没有数据的时候会阻塞直到队列中有数据或者等待关闭
func (q *SyncVector) WaitPop() (v interface{}) {
	q.lock.Lock()
	for q.buffer.Size() == 0 && !q.closed {
		q.popable.Wait()
	}
	if q.buffer.Size() > 0 {
		v = q.buffer.PopFront()
	}
	q.lock.Unlock()
	return
}

// 立即返回,不阻塞,如果队列中数据为空,返回的nil,false
func (q *SyncVector) TryPop() (v interface{}, ok bool) {

	q.lock.Lock()
	ok = false
	if q.buffer.Size() > 0 {
		v = q.buffer.PopFront()
		ok = true
	} else if q.closed {
		ok = true
	}

	q.lock.Unlock()
	return
}

// Push an item to SyncVector. Always returns immediately without blocking
func (q *SyncVector) PushBack(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.PushBack(v)
		q.popable.Signal()
	}
	q.lock.Unlock()
}

// Get the length of SyncVector
func (q *SyncVector) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Size()
	q.lock.Unlock()
	return
}

//关闭对象要让等待放开
func (q *SyncVector) Close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.popable.Signal()
	}
	q.lock.Unlock()
}

//查看队列是否关闭
func (q *SyncVector) IsClose()(b bool) {
	q.lock.Lock()
	b = q.closed
	q.lock.Unlock()
	return
}

//清除队列的数据
func (safeq *SyncVector)Clear(){
	safeq.lock.Lock()
	safeq.buffer.RemoveAll()
	safeq.lock.Unlock()
}
