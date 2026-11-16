package xlist

import (
	"github.com/jrmarcco/jit/internal/errs"
	"github.com/jrmarcco/jit/internal/slice"
)

var _ List[any] = (*ArrayList[any])(nil)

// ArrayList 为基于切片实现的列表。
type ArrayList[T any] struct {
	vals []T
}

// Insert 在 index 位置插入元素，越界时返回错误。
func (al *ArrayList[T]) Insert(index int, val T) error {
	vals, err := slice.Add(al.vals, index, val)
	if err != nil {
		return err
	}
	al.vals = vals
	return nil
}

// Append 在尾部追加元素。
func (al *ArrayList[T]) Append(vals ...T) error {
	al.vals = append(al.vals, vals...)
	return nil
}

// Del 删除 index 位置的元素并按需缩容，越界时返回错误。
func (al *ArrayList[T]) Del(index int) error {
	vals, err := slice.Del(al.vals, index)
	if err != nil {
		return err
	}

	al.vals = slice.Shrink(vals)
	return nil
}

// Set 覆盖 index 位置的元素，越界时返回错误。
func (al *ArrayList[T]) Set(index int, val T) error {
	if index < 0 || index >= len(al.vals) {
		return errs.NewErrIndexOutOfBounds(len(al.vals), index)
	}
	al.vals[index] = val
	return nil
}

// Get 返回 index 位置的元素，越界时返回错误。
func (al *ArrayList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(al.vals) {
		var zero T
		return zero, errs.NewErrIndexOutOfBounds(len(al.vals), index)
	}
	return al.vals[index], nil
}

// Iter 遍历元素，visitFunc 返回非 nil error 时提前终止并返回该 error。
func (al *ArrayList[T]) Iter(visitFunc func(idx int, val T) error) error {
	for i, va := range al.vals {
		if err := visitFunc(i, va); err != nil {
			return err
		}
	}
	return nil
}

// ToSlice 返回元素的副本切片。
func (al *ArrayList[T]) ToSlice() []T {
	res := make([]T, len(al.vals))
	copy(res, al.vals)
	return res
}

// Cap 返回底层数组容量。
func (al *ArrayList[T]) Cap() int {
	return cap(al.vals)
}

// Len 返回元素个数。
func (al *ArrayList[T]) Len() int {
	return len(al.vals)
}

// NewArrayList 创建 ArrayList，size 为底层数组的预估容量。
func NewArrayList[T any](size int) *ArrayList[T] {
	return &ArrayList[T]{
		vals: make([]T, 0, size),
	}
}

// ArrayListOf 由 vals 的副本创建 ArrayList，
// 调用方之后可以继续使用 vals，列表不会读写其底层数组。
func ArrayListOf[T any](vals []T) *ArrayList[T] {
	res := make([]T, len(vals))
	copy(res, vals)

	return &ArrayList[T]{
		vals: res,
	}
}
