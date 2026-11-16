// Package xlist 提供列表实现:
// 数组列表、双向链表、写时复制的数组列表、并发安全包装器与跳表。
package xlist

// List 为列表抽象。
// 除 ConcurrentList 外，各实现均非并发安全。
type List[T any] interface {
	// Insert 在 index 位置插入元素。
	Insert(index int, val T) error

	// Append 在尾部追加元素。
	Append(vals ...T) error

	// Del 删除 index 位置的元素。
	Del(index int) error

	// Set 覆盖 index 位置的元素。
	Set(index int, val T) error

	// Get 返回 index 位置的元素。
	Get(index int) (T, error)

	// Iter 遍历元素，visitFunc 返回非 nil error 时提前终止并返回该 error。
	Iter(visitFunc func(idx int, val T) error) error

	// ToSlice 返回元素的副本切片。
	ToSlice() []T

	// Cap 返回容量 ( 链表实现返回与 Len 相同的值 )。
	Cap() int

	// Len 返回元素个数。
	Len() int
}
