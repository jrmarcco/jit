package xlist

import (
	"github.com/jrmarcco/jit/internal/errs"
	"github.com/jrmarcco/jit/internal/slice"
)

var _ List[any] = (*ArrayList[any])(nil)

type ArrayList[T any] struct {
	vals []T
}

func (al *ArrayList[T]) Insert(index int, val T) error {
	vals, err := slice.Add(al.vals, index, val)
	if err != nil {
		return err
	}
	al.vals = vals
	return nil
}

func (al *ArrayList[T]) Append(vals ...T) error {
	al.vals = append(al.vals, vals...)
	return nil
}

func (al *ArrayList[T]) Del(index int) error {
	vals, err := slice.Del(al.vals, index)
	if err != nil {
		return err
	}

	al.vals = slice.Shrink(vals)
	return nil
}

func (al *ArrayList[T]) Set(index int, val T) error {
	if index < 0 || index >= len(al.vals) {
		return errs.ErrIndexOutOfBounds(len(al.vals), index)
	}
	al.vals[index] = val
	return nil
}

func (al *ArrayList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(al.vals) {
		var zero T
		return zero, errs.ErrIndexOutOfBounds(len(al.vals), index)
	}
	return al.vals[index], nil
}

func (al *ArrayList[T]) Iter(visitFunc func(idx int, val T) error) error {
	for i, va := range al.vals {
		if err := visitFunc(i, va); err != nil {
			return err
		}
	}
	return nil
}

func (al *ArrayList[T]) ToSlice() []T {
	res := make([]T, len(al.vals))
	copy(res, al.vals)
	return res
}

func (al *ArrayList[T]) Cap() int {
	return cap(al.vals)
}

func (al *ArrayList[T]) Len() int {
	return len(al.vals)
}

func NewArrayList[T any](size int) *ArrayList[T] {
	return &ArrayList[T]{
		vals: make([]T, 0, size),
	}
}

func ArrayListOf[T any](slice []T) *ArrayList[T] {
	return &ArrayList[T]{
		vals: slice,
	}
}
