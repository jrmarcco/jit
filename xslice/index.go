package xslice

// LastIndexFunc 返回最后一个满足条件元素的下标，不存在时返回 -1。
// 标准库 slices 包没有等价函数。
func LastIndexFunc[T any](slice []T, match matchFunc[T]) int {
	for k := len(slice) - 1; k >= 0; k-- {
		if match(slice[k]) {
			return k
		}
	}
	return -1
}

// LastIndex 返回 elem 最后一次出现的下标，不存在时返回 -1。
func LastIndex[T comparable](slice []T, elem T) int {
	return LastIndexFunc(slice, func(t T) bool { return t == elem })
}

// IndexAllFunc 返回所有满足条件元素的下标。
func IndexAllFunc[T any](slice []T, match matchFunc[T]) []int {
	res := make([]int, 0, len(slice)>>2+1)
	for k, v := range slice {
		if match(v) {
			res = append(res, k)
		}
	}
	return res
}

// IndexAll 返回所有等于 elem 的元素下标。
func IndexAll[T comparable](slice []T, elem T) []int {
	return IndexAllFunc(slice, func(t T) bool { return t == elem })
}
