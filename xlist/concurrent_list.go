package xlist

import "sync"

var _ List[any] = (*ConcurrentList[any])(nil)

// ConcurrentList 通过 RWMutex 包装任意 List，
// 使其全部操作并发安全:
// 读操作 ( Get/Iter/ToSlice/Cap/Len ) 持读锁，其余操作持写锁。
type ConcurrentList[T any] struct {
	List[T]
	mu sync.RWMutex
}

// NewConcurrentList 创建 l 的并发安全包装。
// l 为 nil 时首次使用会 panic，请传入已构造的 List。
func NewConcurrentList[T any](l List[T]) *ConcurrentList[T] {
	return &ConcurrentList[T]{List: l}
}

// Insert 在 index 位置插入元素，越界时返回错误。
func (cl *ConcurrentList[T]) Insert(index int, val T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Insert(index, val)
}

// Append 在尾部追加元素。
func (cl *ConcurrentList[T]) Append(vals ...T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Append(vals...)
}

// Del 删除 index 位置的元素，越界时返回错误。
func (cl *ConcurrentList[T]) Del(index int) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Del(index)
}

// Set 覆盖 index 位置的元素，越界时返回错误。
func (cl *ConcurrentList[T]) Set(index int, val T) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.List.Set(index, val)
}

// Get 返回 index 位置的元素，越界时返回错误。
func (cl *ConcurrentList[T]) Get(index int) (T, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Get(index)
}

// Iter 遍历元素 ( 持读锁，遍历期间阻塞写操作 )，
// visitFunc 返回非 nil error 时提前终止并返回该 error。
func (cl *ConcurrentList[T]) Iter(visitFunc func(idx int, val T) error) error {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Iter(visitFunc)
}

// ToSlice 返回元素的副本切片。
func (cl *ConcurrentList[T]) ToSlice() []T {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.ToSlice()
}

// Cap 返回容量。
func (cl *ConcurrentList[T]) Cap() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Cap()
}

// Len 返回元素个数。
func (cl *ConcurrentList[T]) Len() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.List.Len()
}
