package xlist

import (
	"sync"
	"sync/atomic"

	"github.com/jrmarcco/jit/internal/errs"
	"github.com/jrmarcco/jit/internal/slice"
)

var _ List[any] = (*CowArrayList[any])(nil)

// CowArrayList 为基于切片的写时复制 ( copy-on-write ) 列表:
// 写操作持锁整体替换切片，读操作无锁。
// 适合读多写少的场景。
//
// 内部切片经 atomic.Pointer 发布:
// 读侧单次 Load 即得到一致快照，写侧在互斥锁内整体替换。
type CowArrayList[T any] struct {
	mu   sync.Mutex
	vals atomic.Pointer[[]T]
}

// NewCowArrayList 创建 CowArrayList，size 为底层数组的预估容量。
func NewCowArrayList[T any](size int) *CowArrayList[T] {
	vals := make([]T, 0, size)
	cal := &CowArrayList[T]{}
	cal.vals.Store(&vals)
	return cal
}

// CowArrayListOf 由 vals 的副本创建 CowArrayList，
// 调用方之后可以继续使用 vals。
func CowArrayListOf[T any](vals []T) *CowArrayList[T] {
	res := make([]T, len(vals))
	copy(res, vals)

	cal := &CowArrayList[T]{}
	cal.vals.Store(&res)
	return cal
}

// Insert 在 index 位置插入元素 ( 复制后写入新切片 )，越界时返回错误。
func (cal *CowArrayList[T]) Insert(index int, val T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	old := *cal.vals.Load()
	length := len(old)
	newVals := make([]T, length, length+1)
	copy(newVals, old)

	var err error
	newVals, err = slice.Add(newVals, index, val)
	if err != nil {
		return err
	}

	cal.vals.Store(&newVals)
	return nil
}

// Append 在尾部追加元素 ( 复制后写入新切片 )。
func (cal *CowArrayList[T]) Append(vals ...T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	old := *cal.vals.Load()
	length := len(old)
	newVals := make([]T, length, length+len(vals))
	copy(newVals, old)

	newVals = append(newVals, vals...)
	cal.vals.Store(&newVals)
	return nil
}

// Del 删除 index 位置的元素 ( 复制后写入新切片 )，越界时返回错误。
func (cal *CowArrayList[T]) Del(index int) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	old := *cal.vals.Load()
	length := len(old)
	if index < 0 || index >= length {
		return errs.NewErrIndexOutOfBounds(length, index)
	}

	newVals := make([]T, length-1)
	newIndex := 0
	for i, v := range old {
		if i == index {
			continue
		}
		newVals[newIndex] = v
		newIndex++
	}
	cal.vals.Store(&newVals)
	return nil
}

// Set 覆盖 index 位置的元素 ( 复制后写入新切片 )，越界时返回错误。
func (cal *CowArrayList[T]) Set(index int, val T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	old := *cal.vals.Load()
	length := len(old)
	if index < 0 || index >= length {
		return errs.NewErrIndexOutOfBounds(length, index)
	}

	newVals := make([]T, length)
	copy(newVals, old)
	newVals[index] = val
	cal.vals.Store(&newVals)
	return nil
}

// Get 返回 index 位置的元素 ( 无锁读取快照 )，越界时返回错误。
func (cal *CowArrayList[T]) Get(index int) (T, error) {
	vals := *cal.vals.Load()
	length := len(vals)
	if index < 0 || index >= length {
		var zero T
		return zero, errs.NewErrIndexOutOfBounds(length, index)
	}
	return vals[index], nil
}

// Iter 无锁遍历单次 Load 取得的快照，
// visitFunc 返回非 nil error 时提前终止并返回该 error。
func (cal *CowArrayList[T]) Iter(visitFunc func(idx int, val T) error) error {
	vals := *cal.vals.Load()
	for i, va := range vals {
		if err := visitFunc(i, va); err != nil {
			return err
		}
	}
	return nil
}

// ToSlice 返回快照的副本切片。
func (cal *CowArrayList[T]) ToSlice() []T {
	vals := *cal.vals.Load()
	res := make([]T, len(vals))
	copy(res, vals)
	return res
}

// Cap 返回当前快照的底层数组容量。
func (cal *CowArrayList[T]) Cap() int {
	return cap(*cal.vals.Load())
}

// Len 返回当前快照的元素个数。
func (cal *CowArrayList[T]) Len() int {
	return len(*cal.vals.Load())
}
