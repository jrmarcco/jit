package xslice

// Find 返回第一个满足条件的元素，不存在时返回零值与 false。
func Find[T any](slice []T, match matchFunc[T]) (T, bool) {
	for _, v := range slice {
		if match(v) {
			return v, true
		}
	}

	var t T
	return t, false
}

// FindAll 返回所有满足条件的元素。
func FindAll[T any](slice []T, match matchFunc[T]) []T {
	// 预估结果切片容量: 切片长度的 25% 再加一。
	res := make([]T, 0, len(slice)>>2+1)
	for _, v := range slice {
		if match(v) {
			res = append(res, v)
		}
	}
	return res
}
