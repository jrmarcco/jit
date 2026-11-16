package xlist

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/list"
)

// SkipList 为基于跳表实现的有序列表。
//
// 注意:
//   - Insert 允许重复值 ( 多重集语义 )，Delete 每次删除其中一个。
//   - SkipList 非并发安全，需要外部同步。
type SkipList[T any] struct {
	skipList *list.SkipList[T]
}

// NewSkipList 创建 SkipList，cmp 为元素比较器。
func NewSkipList[T any](cmp jit.Comparator[T]) *SkipList[T] {
	return &SkipList[T]{
		skipList: list.NewSkipList(cmp),
	}
}

// Insert 插入一个值 ( 允许重复 )。
func (sl *SkipList[T]) Insert(val T) {
	sl.skipList.Insert(val)
}

// Delete 删除一个目标值，成功返回 true，目标不存在返回 false。
func (sl *SkipList[T]) Delete(target T) bool {
	return sl.skipList.Delete(target)
}

// Exists 判断目标值是否存在。
func (sl *SkipList[T]) Exists(target T) bool {
	return sl.skipList.Exists(target)
}

// Get 返回第 index 个值 ( 按升序，index 从 0 起 )，越界返回 false。
func (sl *SkipList[T]) Get(index int) (T, bool) {
	return sl.skipList.GetByIndex(index)
}

// Peek 返回最小的值 ( 不删除 )，列表为空返回 false。
func (sl *SkipList[T]) Peek() (T, bool) {
	return sl.skipList.Peek()
}

// Len 返回元素个数。
func (sl *SkipList[T]) Len() int {
	return sl.skipList.Len()
}

// ToSlice 返回按升序排列的全部元素。
func (sl *SkipList[T]) ToSlice() []T {
	return sl.skipList.ToSlice()
}
