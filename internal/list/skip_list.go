package list

import (
	"math/rand/v2"

	"github.com/jrmarcco/jit"
)

const MaxLevel = 32

type SkipList[T any] struct {
	cmp       jit.Comparator[T]
	head      *skipListNode[T]
	currLevel int
	size      int
}

func NewSkipList[T any](cmp jit.Comparator[T]) *SkipList[T] {
	return &SkipList[T]{
		cmp: cmp,
		head: &skipListNode[T]{
			next: make([]*skipListNode[T], MaxLevel),
			span: make([]int, MaxLevel),
		},
		currLevel: 1,
		size:      0,
	}
}

func SkipListOf[T any](cmp jit.Comparator[T], slice []T) *SkipList[T] {
	sl := NewSkipList(cmp)
	for _, v := range slice {
		sl.Insert(v)
	}
	return sl
}

// Insert inserts a value into the skip list.
func (sl *SkipList[T]) Insert(val T) {
	update := make([]*skipListNode[T], MaxLevel)
	// record the cumulative span of the current node in each level
	rank := make([]int, MaxLevel)

	currN := sl.head
	// find the place to insert
	for i := sl.currLevel - 1; i >= 0; i-- {
		if i == sl.currLevel-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// update the higher level node to head node
		for currN.next[i] != nil && sl.cmp(currN.next[i].val, val) < 0 {
			// cumulative span
			rank[i] += currN.span[i]
			currN = currN.next[i]
		}
		update[i] = currN
	}

	// generate level of new node
	level := sl.randLevel()
	if level > sl.currLevel {
		for i := sl.currLevel; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].span[i] = sl.size
		}
		sl.currLevel = level
	}

	// create a new node
	newN := &skipListNode[T]{
		val:  val,
		next: make([]*skipListNode[T], level),
		span: make([]int, level),
	}

	// update index for all levels
	for i := 0; i < level; i++ {
		newN.next[i] = update[i].next[i]
		update[i].next[i] = newN

		newN.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].span[i] = (rank[0] - rank[i]) + 1
	}

	for i := level; i < sl.currLevel; i++ {
		update[i].span[i]++
	}

	sl.size++
}

// Delete remove target from the skip list.
func (sl *SkipList[T]) Delete(target T) bool {
	update := make([]*skipListNode[T], MaxLevel)
	currN := sl.head

	// find the precursor node
	for i := sl.currLevel - 1; i >= 0; i-- {
		for currN.next[i] != nil && sl.cmp(currN.next[i].val, target) < 0 {
			currN = currN.next[i]
		}
		update[i] = currN
	}

	// locate the target node
	currN = currN.next[0]
	if currN == nil || sl.cmp(currN.val, target) != 0 {
		return false
	}

	// update index for all levels
	for i := 0; i < len(currN.next); i++ {
		if update[i].next[i] == currN {
			update[i].span[i] += currN.span[i] - 1
			update[i].next[i] = currN.next[i]
			continue
		}
		update[i].span[i]--
	}

	for i := len(currN.next); i < sl.currLevel; i++ {
		update[i].span[i]--
	}

	// update current level of skip list
	for sl.currLevel > 1 && sl.head.next[sl.currLevel-1] == nil {
		sl.currLevel--
	}

	sl.size--
	return true
}

// Exists checks if the target element exists in the skip list and returns true if found, otherwise returns false.
func (sl *SkipList[T]) Exists(target T) bool {
	currN := sl.head

	for i := sl.currLevel - 1; i >= 0; i-- {
		for currN.next[i] != nil && sl.cmp(currN.next[i].val, target) < 0 {
			currN = currN.next[i]
		}
	}

	currN = currN.next[0]
	return currN != nil && sl.cmp(currN.val, target) == 0
}

func (sl *SkipList[T]) randLevel() int {
	level := 1
	//nolint:gosec // G404: math/rand is sufficient for skip list balancing
	for rand.Float64() < 0.25 && level < MaxLevel {
		level++
	}
	return level
}

// GetByIndex retrieves the value at the specified index in the skip list.
// The index starts from 0.
func (sl *SkipList[T]) GetByIndex(index int) (T, bool) {
	if index < 0 || index >= sl.size {
		var zero T
		return zero, false
	}

	currN := sl.head
	pos := -1

	for i := sl.currLevel - 1; i >= 0; i-- {
		for currN.next[i] != nil && (pos+currN.span[i]) <= index {
			pos += currN.span[i]
			currN = currN.next[i]
		}
	}

	if pos == index && currN != nil {
		return currN.val, true
	}

	var zero T
	return zero, false
}

// Peek returns the first element in the skip list without removing it.
func (sl *SkipList[T]) Peek() (T, bool) {
	if sl.head.next[0] == nil {
		var zero T
		return zero, false
	}
	return sl.head.next[0].val, true
}

func (sl *SkipList[T]) ToSlice() []T {
	res := make([]T, 0, sl.size)
	currN := sl.head

	for currN.next[0] != nil {
		res = append(res, currN.next[0].val)
		currN = currN.next[0]
	}
	return res
}

func (sl *SkipList[T]) Len() int {
	return sl.size
}

type skipListNode[T any] struct {
	val  T
	next []*skipListNode[T]
	span []int
}
