package xslice

import "github.com/jrmarcco/jit/internal/slice"

// Del removes an element at the specified index from the slice.
func Del[T any](src []T, index int) ([]T, error) {
	return slice.Del(src, index)
}

// FilterDel removes elements from the slice that match the filter function.
func FilterDel[T any](src []T, filter func(idx int, elem T) bool) []T {
	index := 0

	for idx := range src {
		if filter(idx, src[idx]) {
			continue
		}

		src[index] = src[idx]
		index++
	}

	return src[:index]
}
