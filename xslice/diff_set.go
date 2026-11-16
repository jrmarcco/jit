package xslice

import "slices"

// DiffSet 返回 src 中不在 dst 内的元素 ( 差集 )。
// 返回元素的顺序不做保证。
func DiffSet[T comparable](src, dst []T) []T {
	srcMap := toMap(src)
	for _, v := range dst {
		// 从 srcMap 中移除 v。
		delete(srcMap, v)
	}

	ret := make([]T, 0, len(srcMap))
	for k := range srcMap {
		ret = append(ret, k)
	}

	return ret
}

// DiffSetFunc 返回 src 中不在 dst 内的元素 ( 相等性由 eq 定义 )。
// 结果保持 src 的顺序。
func DiffSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	slice := make([]T, 0, len(src))
	for _, v := range src {
		// 判断 v 是否在 dst 中。
		if !slices.ContainsFunc(dst, func(t T) bool { return eq(v, t) }) {
			slice = append(slice, v)
		}
	}

	ret := make([]T, 0, len(slice))
	// 去重。
	for idx, val := range slice {
		if !slices.ContainsFunc(slice[idx+1:], func(t T) bool { return eq(val, t) }) {
			ret = append(ret, val)
		}
	}

	return ret
}
