// Package list 实现 xlist.SkipList 底层的跳表。
package list

import (
	"math/rand/v2"

	"github.com/jrmarcco/jit"
)

// MaxLevel 为跳表节点的最大层数。
const MaxLevel = 32

// SkipList 为基于跳表实现的有序列表
// ( Redis zskiplist 风格，span 支持按排名查询 )。
//
// 注意:
//   - Insert 允许重复值: 相等的值存于多个节点 ( 多重集语义 )，
//     Delete 每次删除其中一个。
//   - SkipList 非并发安全，需要外部同步 ( 如 sync.RWMutex )。
type SkipList[T any] struct {
	cmp       jit.Comparator[T]
	head      *skipListNode[T]
	currLevel int
	size      int
}

// NewSkipList 创建空跳表，cmp 为元素比较器。
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

// SkipListOf 由切片创建跳表 ( 按序逐个插入 )。
func SkipListOf[T any](cmp jit.Comparator[T], slice []T) *SkipList[T] {
	sl := NewSkipList(cmp)
	for _, v := range slice {
		sl.Insert(v)
	}
	return sl
}

// Insert 插入一个值 ( 允许重复 )。
func (sl *SkipList[T]) Insert(val T) {
	update := make([]*skipListNode[T], MaxLevel)
	// 记录当前节点在各层的累计 span ( 排名 )。
	rank := make([]int, MaxLevel)

	currN := sl.head
	// 寻找插入位置。
	for i := sl.currLevel - 1; i >= 0; i-- {
		if i == sl.currLevel-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 把高层的前驱节点更新为 head。
		for currN.next[i] != nil && sl.cmp(currN.next[i].val, val) < 0 {
			// 累计 span。
			rank[i] += currN.span[i]
			currN = currN.next[i]
		}
		update[i] = currN
	}

	// 生成新节点的层数。
	level := sl.randLevel()
	if level > sl.currLevel {
		for i := sl.currLevel; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].span[i] = sl.size
		}
		sl.currLevel = level
	}

	// 创建新节点。
	newN := &skipListNode[T]{
		val:  val,
		next: make([]*skipListNode[T], level),
		span: make([]int, level),
	}

	// 更新各层索引。
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

// Delete 删除一个目标值，成功返回 true，目标不存在返回 false。
func (sl *SkipList[T]) Delete(target T) bool {
	update := make([]*skipListNode[T], MaxLevel)
	currN := sl.head

	// 寻找目标节点的前驱节点。
	for i := sl.currLevel - 1; i >= 0; i-- {
		for currN.next[i] != nil && sl.cmp(currN.next[i].val, target) < 0 {
			currN = currN.next[i]
		}
		update[i] = currN
	}

	// 定位目标节点。
	currN = currN.next[0]
	if currN == nil || sl.cmp(currN.val, target) != 0 {
		return false
	}

	// 更新各层索引。
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

	// 回收跳表当前的最高层。
	for sl.currLevel > 1 && sl.head.next[sl.currLevel-1] == nil {
		sl.currLevel--
	}

	sl.size--
	return true
}

// Exists 判断目标元素是否存在。
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

// randLevel 按概率 p=0.25 随机生成节点层数，上限 MaxLevel。
func (sl *SkipList[T]) randLevel() int {
	level := 1
	//nolint:gosec // G404: 跳表均衡只需要非安全随机数。
	for rand.Float64() < 0.25 && level < MaxLevel {
		level++
	}
	return level
}

// GetByIndex 返回第 index 个值 ( 按升序，index 从 0 起 )，越界返回 false。
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

// Peek 返回最小的值 ( 不删除 )，列表为空返回 false。
func (sl *SkipList[T]) Peek() (T, bool) {
	if sl.head.next[0] == nil {
		var zero T
		return zero, false
	}
	return sl.head.next[0].val, true
}

// ToSlice 返回按升序排列的全部元素。
func (sl *SkipList[T]) ToSlice() []T {
	res := make([]T, 0, sl.size)
	currN := sl.head

	for currN.next[0] != nil {
		res = append(res, currN.next[0].val)
		currN = currN.next[0]
	}
	return res
}

// Len 返回元素个数。
func (sl *SkipList[T]) Len() int {
	return sl.size
}

// skipListNode 为跳表节点。
type skipListNode[T any] struct {
	val  T
	next []*skipListNode[T]
	span []int
}
