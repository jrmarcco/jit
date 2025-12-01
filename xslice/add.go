package xslice

import "github.com/jrmarcco/jit/internal/slice"

// Add insert an item at the specified index in the slice.
func Add[T any](src []T, index int, item T) ([]T, error) {
	return slice.Add(src, index, item)
}
