package xlist

import "github.com/jrmarcco/jit/internal/errs"

var _ List[any] = (*LinkedList[any])(nil)

// LinkedList 为基于双向链表 ( 带头尾哨兵节点 ) 实现的列表。
type LinkedList[T any] struct {
	head *linkedListNode[T]
	tail *linkedListNode[T]

	size int
}

// NewLinkedList 创建空链表。
func NewLinkedList[T any]() *LinkedList[T] {
	head := &linkedListNode[T]{}
	tail := &linkedListNode[T]{
		prev: head,
		next: head,
	}
	head.next = tail
	head.prev = tail

	return &LinkedList[T]{
		head: head,
		tail: tail,
	}
}

// LinkedListOf 由 vals 创建链表。
func LinkedListOf[T any](vals []T) *LinkedList[T] {
	ll := NewLinkedList[T]()
	_ = ll.Append(vals...)
	return ll
}

// Insert 在 index 位置插入元素，越界时返回错误。
func (ll *LinkedList[T]) Insert(index int, val T) error {
	if index < 0 || index > ll.size {
		return errs.NewErrIndexOutOfBounds(ll.size, index)
	}

	if index == ll.size {
		// 插入尾部。
		return ll.Append(val)
	}

	currN := ll.findNode(index)
	newN := &linkedListNode[T]{
		val:  val,
		prev: currN.prev,
		next: currN,
	}

	newN.prev.next, newN.next.prev = newN, newN
	ll.size++
	return nil
}

// findNode 返回 index 位置的节点 ( 按就近方向遍历 )。
func (ll *LinkedList[T]) findNode(index int) *linkedListNode[T] {
	if index <= ll.size/2 {
		// 从头部向后遍历。
		currN := ll.head
		for i := 0; i <= index; i++ {
			currN = currN.next
		}
		return currN
	}

	// 从尾部向前遍历。
	currN := ll.tail
	for i := ll.size - 1; i >= index; i-- {
		currN = currN.prev
	}
	return currN
}

// Append 在尾部追加元素。
func (ll *LinkedList[T]) Append(vals ...T) error {
	for _, val := range vals {
		newN := &linkedListNode[T]{
			val:  val,
			prev: ll.tail.prev,
			next: ll.tail,
		}
		newN.prev.next, newN.next.prev = newN, newN
		ll.size++
	}
	return nil
}

// Del 删除 index 位置的节点，越界时返回错误。
func (ll *LinkedList[T]) Del(index int) error {
	if index < 0 || index >= ll.size {
		return errs.NewErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)

	currN.prev.next = currN.next
	currN.next.prev = currN.prev
	currN.prev, currN.next = nil, nil

	ll.size--
	return nil
}

// Set 覆盖 index 位置的元素，越界时返回错误。
func (ll *LinkedList[T]) Set(index int, val T) error {
	if index < 0 || index >= ll.size {
		return errs.NewErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)
	currN.val = val
	return nil
}

// Get 返回 index 位置的元素，越界时返回错误。
func (ll *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, errs.NewErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)
	return currN.val, nil
}

// Iter 遍历元素，visitFunc 返回非 nil error 时提前终止并返回该 error。
func (ll *LinkedList[T]) Iter(visitFunc func(idx int, val T) error) error {
	for currN, index := ll.head.next, 0; currN != ll.tail; currN, index = currN.next, index+1 {
		if err := visitFunc(index, currN.val); err != nil {
			return err
		}
	}
	return nil
}

// ToSlice 返回元素的副本切片。
func (ll *LinkedList[T]) ToSlice() []T {
	res := make([]T, ll.size)
	for currN, index := ll.head.next, 0; currN != ll.tail; currN, index = currN.next, index+1 {
		res[index] = currN.val
	}
	return res
}

// Cap 返回链表节点数，与 Len 相同 ( 链表无冗余容量概念 )。
func (ll *LinkedList[T]) Cap() int {
	return ll.size
}

// Len 返回元素个数。
func (ll *LinkedList[T]) Len() int {
	return ll.size
}

// linkedListNode 为双向链表节点。
type linkedListNode[T any] struct {
	val  T
	prev *linkedListNode[T]
	next *linkedListNode[T]
}
