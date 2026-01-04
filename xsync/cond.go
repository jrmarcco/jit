package xsync

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Cond struct {
	L sync.Locker

	notifyList *notifyList
	checker    unsafe.Pointer
	once       sync.Once

	noCopy noCopy
}

func NewCond(locker sync.Locker) *Cond {
	return &Cond{
		L: locker,
	}
}

// Wait like Go standard library's sync.Cond.Wait.
// have to lock the lock before calling Wait.
//
// Wait will lock the lock and add the current goroutine to the wait list.
//
// When Wait is invoked, the current goroutine will release the lock first and go to the wait list.
//
//	( waiting for other goroutines to invoke Signal or Broadcast ).
//
// After being notified, Wait will lock the lock and return.
// The point of this is that can hold the lock again after the Wait method returns.
// It is safe to check conditions and process data.
func (c *Cond) Wait(ctx context.Context) error {
	c.checkCopy()
	c.initNotifyListOnce()

	// put the current goroutine to the wait list.
	n := c.notifyList.add()

	// before waiting, must lock the lock.
	c.L.Unlock()
	defer c.L.Lock()

	return c.notifyList.wait(ctx, n)
}

// checkCopy checks whether the Cond is copied.
func (c *Cond) checkCopy() {
	// check the checker pointer is to the current address.
	if c.checker != unsafe.Pointer(c) &&
		// ensure the checker pointer is not changed ( checker pointer first assigned ).
		// the checker pointer is nil when first created.
		!atomic.CompareAndSwapPointer(&c.checker, nil, unsafe.Pointer(c)) &&
		// check the checker pointer again.
		// the checker pointer has to the current address.
		c.checker != unsafe.Pointer(c) {
		panic("xsync: Cond is copied")
	}
}

// initNotifyListOnce initializes the notifyList.
func (c *Cond) initNotifyListOnce() {
	c.once.Do(func() {
		if c.notifyList == nil {
			c.notifyList = newNotifyList()
		}
	})
}

// Signal notify one goroutine that wait in Cond.
func (c *Cond) Signal() {
	c.checkCopy()
	c.initNotifyListOnce()

	c.notifyList.notifyOne()
}

// Broadcast notifies all goroutines that wait in Cond.
func (c *Cond) Broadcast() {
	c.checkCopy()
	c.initNotifyListOnce()

	c.notifyList.notifyAll()
}

type notifyList struct {
	lock sync.Mutex
	list *chanList
}

// add adds a new channel to the list.
func (l *notifyList) add() *node {
	l.lock.Lock()
	defer l.lock.Unlock()

	n := l.list.allocate()
	l.list.pushBack(n)
	return n
}

// wait waits for a channel to be notified.
func (l *notifyList) wait(ctx context.Context, n *node) error {
	ch := n.Value
	// put the node back in the pool.
	// it's useful to reduce GC and memory allocation.
	defer l.list.free(n)

	select {
	case <-ctx.Done():
		// timeout or canceled.
		l.lock.Lock()
		defer l.lock.Unlock()

		select {
		// double-check: notified before locked.
		// if notified, notify next.
		// why notify the next one?
		//		if the notified channel is not consumed by the current goroutine,
		//		then the other waiters will never receive the signal.
		//		then the queue will be blocked.
		//		you've been received a signal here means that you should be notified.
		//		but you missed it ( context canceled or timeout ).
		//		so need to notify the next one.
		case <-ch:
			if l.list.len() != 0 {
				l.notifyNext()
			}
		default:
			// did not receive a signal until context canceled or timeout.
			// means the invoker given up waiting for the signal.
			// this node will never be notified anymore.
			l.list.remove(n)
		}
		return ctx.Err()
	case <-ch:
		// notified.
		return nil
	}
}

// notifyOne notifies one channel in the list.
// if the list is empty, notifyOne does nothing.
func (l *notifyList) notifyOne() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.list.len() == 0 {
		return
	}

	l.notifyNext()
}

// notifyNext notifies the next channel in the list.
func (l *notifyList) notifyNext() {
	front := l.list.front()
	ch := front.Value
	l.list.remove(front)

	// send a signal to the channel.
	ch <- struct{}{}
}

// notifyAll notifies all channels in the list.
// if the list is empty, notifyAll does nothing.
func (l *notifyList) notifyAll() {
	l.lock.Lock()
	defer l.lock.Unlock()

	for l.list.len() != 0 {
		l.notifyNext()
	}
}

func newNotifyList() *notifyList {
	return &notifyList{
		lock: sync.Mutex{},
		list: newChanList(),
	}
}

// node is used to store the channels.
// is a linked list element.
type node struct {
	prev  *node
	next  *node
	Value chan struct{}
}

// chanList is a double-linked list for saving the channels.
// use pool to reuse list elements.
type chanList struct {
	// in the double-linked list,
	// sentinel is in the middle of the list.
	// the list front is sentinel.next,
	// and the list back is sentinel.prev.
	// [sentinel] <-> node1 <-> node2 <-> [sentinel]
	//            		↑         ↑
	//          	front        tail
	// when sentinel.prev and sentinel.next are sentinel itself means the list is empty.
	sentinel *node
	size     int64
	pool     *sync.Pool
}

func newChanList() *chanList {
	sentinel := &node{}
	// list is empty, sentinel.prev and sentinel.next are sentinel itself.
	sentinel.next = sentinel
	sentinel.prev = sentinel

	return &chanList{
		sentinel: sentinel,
		size:     0,
		pool: &sync.Pool{
			New: func() any {
				return &node{
					Value: make(chan struct{}, 1),
				}
			},
		},
	}
}

// allocate returns a node from the pool.
// if the pool is empty, a new node is created.
func (l *chanList) allocate() *node {
	n, ok := l.pool.Get().(*node)
	if !ok {
		panic("[jit] pool contained non-node value")
	}
	return n
}

// pushBack adds a new node to the back of the list.
func (l *chanList) pushBack(n *node) {
	n.next = l.sentinel
	n.prev = l.sentinel.prev
	n.prev.next = n
	n.next.prev = n
	l.size++
}

// remove removes a node from the list.
// n must not be nil.
func (l *chanList) remove(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	l.size--
}

// free put the node back in the pool.
func (l *chanList) free(n *node) {
	l.pool.Put(n)
}

// front returns the first node of the list.
// if the list is empty, nil is returned.
func (l *chanList) front() *node {
	return l.sentinel.next
}

func (l *chanList) len() int64 {
	return l.size
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
