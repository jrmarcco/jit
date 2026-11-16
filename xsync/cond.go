package xsync

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Cond 为支持 context 的条件变量。
//
// 与标准库 sync.Cond 的差异:
// Wait 接收 context，可在超时或取消时提前返回。
type Cond struct {
	// L 为 Wait 使用的外部锁 ( 与 sync.Cond.L 相同 )。
	L sync.Locker

	notifyList *notifyList
	checker    unsafe.Pointer
	once       sync.Once

	noCopy noCopy
}

// NewCond 创建 Cond，locker 为配套的外部锁。
func NewCond(locker sync.Locker) *Cond {
	return &Cond{
		L: locker,
	}
}

// Wait 类似标准库 sync.Cond.Wait:
// 调用前必须已持有锁 L。
//
// Wait 会先把当前 goroutine 加入等待队列，随后释放锁 L 并进入等待
// ( 等待其他 goroutine 调用 Signal 或 Broadcast )。
//
// 被唤醒后 Wait 会重新加锁再返回，
// 因此返回后调用方已重新持有锁，可以安全地检查条件与处理数据。
//
// context 超时或取消时 Wait 返回 ctx.Err() ( 同样会重新加锁 )。
func (c *Cond) Wait(ctx context.Context) error {
	c.checkCopy()
	c.initNotifyListOnce()

	// 把当前 goroutine 加入等待队列。
	n := c.notifyList.add()

	// 等待前必须先释放锁。
	c.L.Unlock()
	defer c.L.Lock()

	return c.notifyList.wait(ctx, n)
}

// checkCopy 检测 Cond 是否被复制，被复制时 panic。
func (c *Cond) checkCopy() {
	// 校验 checker 指针是否指向当前地址。
	// 所有读取必须原子化，避免与下方的 CAS 构成数据竞争。
	//nolint:gosec // G103: 复制检测需要故意使用 unsafe。
	if atomic.LoadPointer(&c.checker) != unsafe.Pointer(c) &&
		// 确保 checker 指针未被改写 ( checker 指针首次赋值 )。
		// 新建时 checker 指针为 nil。
		//nolint:gosec // G103: 复制检测需要故意使用 unsafe。
		!atomic.CompareAndSwapPointer(&c.checker, nil, unsafe.Pointer(c)) &&
		// 再次校验 checker 指针，
		// 必须等于当前地址，否则说明发生了复制。
		//nolint:gosec // G103: 复制检测需要故意使用 unsafe。
		atomic.LoadPointer(&c.checker) != unsafe.Pointer(c) {
		panic("xsync: Cond is copied")
	}
}

// initNotifyListOnce 惰性初始化 notifyList。
func (c *Cond) initNotifyListOnce() {
	c.once.Do(func() {
		if c.notifyList == nil {
			c.notifyList = newNotifyList()
		}
	})
}

// Signal 唤醒一个等待中的 goroutine ( 无等待者时不产生任何效果 )。
func (c *Cond) Signal() {
	c.checkCopy()
	c.initNotifyListOnce()

	c.notifyList.notifyOne()
}

// Broadcast 唤醒全部等待中的 goroutine。
func (c *Cond) Broadcast() {
	c.checkCopy()
	c.initNotifyListOnce()

	c.notifyList.notifyAll()
}

// notifyList 为等待队列: 以互斥锁保护的双向链表管理各等待者的通知 channel。
type notifyList struct {
	lock sync.Mutex
	list *chanList
}

// add 新建一个节点并加入等待队列尾部，返回该节点。
func (l *notifyList) add() *node {
	l.lock.Lock()
	defer l.lock.Unlock()

	n := l.list.allocate()
	l.list.pushBack(n)
	return n
}

// wait 等待节点 n 的通知，返回 nil 表示被唤醒，
// context 超时或取消时返回 ctx.Err()。
func (l *notifyList) wait(ctx context.Context, n *node) error {
	ch := n.Value
	// 把节点归还对象池，减少 GC 与分配开销。
	defer l.list.free(n)

	select {
	case <-ctx.Done():
		// 超时或取消。
		l.lock.Lock()
		defer l.lock.Unlock()

		select {
		// 二次确认: 加锁前已收到通知。
		// 若已收到通知，则向下传递唤醒下一个等待者。
		// 为什么要唤醒下一个 ?
		//   如果发给当前 channel 的信号没有被当前 goroutine 消费，
		//   其他等待者将永远等不到信号，队列就此卡死。
		//   此处既然收到了信号，说明本应唤醒的是当前 goroutine，
		//   但它因 context 取消/超时放弃了等待，错过了信号，
		//   因此需要把信号传递给下一个等待者。
		case <-ch:
			if l.list.len() != 0 {
				l.notifyNext()
			}
		default:
			// 直到 context 取消/超时仍未收到信号，
			// 说明调用方已放弃等待，该节点不会再被通知。
			l.list.remove(n)
		}
		return ctx.Err()
	case <-ch:
		// 收到通知。
		return nil
	}
}

// notifyOne 唤醒队列中最先入队的等待者，队列为空时不产生任何效果。
func (l *notifyList) notifyOne() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.list.len() == 0 {
		return
	}

	l.notifyNext()
}

// notifyNext 唤醒队首等待者并将其移出队列。
func (l *notifyList) notifyNext() {
	front := l.list.front()
	ch := front.Value
	l.list.remove(front)

	// 发送唤醒信号 ( channel 容量为 1，且每个节点至多发一次，不会阻塞 )。
	ch <- struct{}{}
}

// notifyAll 唤醒队列中的全部等待者，队列为空时不产生任何效果。
func (l *notifyList) notifyAll() {
	l.lock.Lock()
	defer l.lock.Unlock()

	for l.list.len() != 0 {
		l.notifyNext()
	}
}

// newNotifyList 创建 notifyList。
func newNotifyList() *notifyList {
	return &notifyList{
		lock: sync.Mutex{},
		list: newChanList(),
	}
}

// node 为存储通知 channel 的双向链表节点。
type node struct {
	prev  *node
	next  *node
	Value chan struct{}
}

// chanList 为双向链表 + 对象池复用节点。
type chanList struct {
	// 双向循环链表以哨兵为中心:
	// 链表头为 sentinel.next，链表尾为 sentinel.prev。
	// [sentinel] <-> node1 <-> node2 <-> [sentinel]
	//              ↑          ↑
	//            front       tail
	// sentinel.prev 与 sentinel.next 均指向 sentinel 自身时链表为空。
	sentinel *node
	size     int64
	pool     *sync.Pool
}

// newChanList 创建 chanList。
func newChanList() *chanList {
	sentinel := &node{}
	// 空链表: sentinel 的 prev 与 next 均指向自身。
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

// allocate 从对象池取出一个节点。
func (l *chanList) allocate() *node {
	n, ok := l.pool.Get().(*node)
	if !ok {
		panic("[jit] pool contained non-node value")
	}
	return n
}

// pushBack 把节点加入链表尾部。
func (l *chanList) pushBack(n *node) {
	n.next = l.sentinel
	n.prev = l.sentinel.prev
	n.prev.next = n
	n.next.prev = n
	l.size++
}

// remove 把节点移出链表，n 不能为 nil。
func (l *chanList) remove(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	l.size--
}

// free 把节点归还对象池。
func (l *chanList) free(n *node) {
	l.pool.Put(n)
}

// front 返回链表头节点，链表为空时返回哨兵自身。
func (l *chanList) front() *node {
	return l.sentinel.next
}

// len 返回链表长度。
func (l *chanList) len() int64 {
	return l.size
}

// noCopy 用于配合 go vet 的 copylocks 检测 Cond 被复制。
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
