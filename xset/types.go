// Package xset 提供集合实现:
// 基于内建 map 的 MapSet 与基于红黑树的 TreeSet。
package xset

// Set 为元素唯一的集合抽象。
// 除非实现另行说明，各实现均非并发安全。
type Set[T comparable] interface {
	Add(key T)
	Del(key T)
	Exist(key T) bool
	Elems() []T
	Size() int
}
