package xlist

import "github.com/jrmarcco/jit/internal/errs"

var _ List[any] = (*LinkedList[any])(nil)

type LinkedList[T any] struct {
	head *linkedListNode[T]
	tail *linkedListNode[T]
	size int
}

func (ll *LinkedList[T]) Insert(index int, val T) error {
	if index < 0 || index > ll.size {
		return errs.ErrIndexOutOfBounds(ll.size, index)
	}

	if index == ll.size {
		// insert to tail
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

func (ll *LinkedList[T]) findNode(index int) *linkedListNode[T] {
	if index <= ll.size/2 {
		currN := ll.head
		for i := 0; i <= index; i++ {
			currN = currN.next
		}
		return currN
	}

	currN := ll.tail
	for i := ll.size - 1; i >= index; i-- {
		currN = currN.prev
	}
	return currN
}

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

func (ll *LinkedList[T]) Del(index int) error {
	if index < 0 || index >= ll.size {
		return errs.ErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)

	currN.prev.next = currN.next
	currN.next.prev = currN.prev
	currN.prev, currN.next = nil, nil

	ll.size--
	return nil
}

func (ll *LinkedList[T]) Set(index int, val T) error {
	if index < 0 || index >= ll.size {
		return errs.ErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)
	currN.val = val
	return nil
}

func (ll *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, errs.ErrIndexOutOfBounds(ll.size, index)
	}

	currN := ll.findNode(index)
	return currN.val, nil
}

func (ll *LinkedList[T]) Iter(visitFunc func(idx int, val T) error) error {
	for currN, index := ll.head.next, 0; currN != ll.tail; currN, index = currN.next, index+1 {
		if err := visitFunc(index, currN.val); err != nil {
			return err
		}
	}
	return nil
}

func (ll *LinkedList[T]) ToSlice() []T {
	res := make([]T, ll.size)
	for currN, index := ll.head.next, 0; currN != ll.tail; currN, index = currN.next, index+1 {
		res[index] = currN.val
	}
	return res
}

func (ll *LinkedList[T]) Cap() int {
	return ll.size
}

func (ll *LinkedList[T]) Len() int {
	return ll.size
}

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

func LinkedListOf[T any](vals []T) *LinkedList[T] {
	ll := NewLinkedList[T]()
	_ = ll.Append(vals...)
	return ll
}

type linkedListNode[T any] struct {
	val  T
	prev *linkedListNode[T]
	next *linkedListNode[T]
}
