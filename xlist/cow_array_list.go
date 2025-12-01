package xlist

import (
	"sync"

	"github.com/jrmarcco/jit/internal/errs"
	"github.com/jrmarcco/jit/internal/slice"
)

var _ List[any] = (*CowArrayList[any])(nil)

// CowArrayList a copy-on-write array list implementation base on slice.
// Lock on writing and no lock on reading.
// Suitable for read-many-write-few.
type CowArrayList[T any] struct {
	mu   sync.Mutex
	vals []T
}

func (cal *CowArrayList[T]) Insert(index int, val T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	length := len(cal.vals)
	newVals := make([]T, length, length+1)
	// copy slice
	copy(newVals, cal.vals)

	var err error
	newVals, err = slice.Add(newVals, index, val)
	if err != nil {
		return err
	}

	cal.vals = newVals
	return nil
}

func (cal *CowArrayList[T]) Append(vals ...T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	length := len(cal.vals)
	newVals := make([]T, length, length+len(vals))
	copy(newVals, cal.vals)

	newVals = append(newVals, vals...)
	cal.vals = newVals
	return nil
}

func (cal *CowArrayList[T]) Del(index int) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	length := len(cal.vals)
	if index < 0 || index >= length {
		return errs.ErrIndexOutOfBounds(length, index)
	}

	newVals := make([]T, length-1)
	newIndex := 0
	for i, v := range cal.vals {
		if i == index {
			continue
		}
		newVals[newIndex] = v
		newIndex++
	}
	cal.vals = newVals
	return nil
}

func (cal *CowArrayList[T]) Set(index int, val T) error {
	cal.mu.Lock()
	defer cal.mu.Unlock()

	length := len(cal.vals)
	if index < 0 || index >= length {
		return errs.ErrIndexOutOfBounds(length, index)
	}

	newVals := make([]T, length)
	copy(newVals, cal.vals)
	newVals[index] = val
	cal.vals = newVals
	return nil
}

func (cal *CowArrayList[T]) Get(index int) (T, error) {
	length := len(cal.vals)
	if index < 0 || index >= length {
		var zero T
		return zero, errs.ErrIndexOutOfBounds(length, index)
	}
	return cal.vals[index], nil
}

func (cal *CowArrayList[T]) Iter(visitFunc func(idx int, val T) error) error {
	for i, va := range cal.vals {
		if err := visitFunc(i, va); err != nil {
			return err
		}
	}
	return nil
}

func (cal *CowArrayList[T]) ToSlice() []T {
	res := make([]T, len(cal.vals))
	copy(res, cal.vals)
	return res
}

func (cal *CowArrayList[T]) Cap() int {
	return cap(cal.vals)
}

func (cal *CowArrayList[T]) Len() int {
	return len(cal.vals)
}

func NewCowArrayList[T any](size int) *CowArrayList[T] {
	return &CowArrayList[T]{
		vals: make([]T, 0, size),
	}
}

func CowArrayListOf[T any](slice []T) *CowArrayList[T] {
	vals := make([]T, len(slice))
	copy(vals, slice)

	return &CowArrayList[T]{
		vals: vals,
	}
}
