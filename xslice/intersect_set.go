package xslice

import "slices"

// IntersectSetFunc 返回两个切片的交集 ( 相等性由 eq 定义 )。
// 结果保持 src 的顺序。
func IntersectSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(src))
	for _, v := range src {
		if slices.ContainsFunc(dst, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}
	return deDuplicateFunc(res, eq)
}

// IntersectSet 返回两个切片的交集。
// 返回元素的顺序不做保证。
func IntersectSet[T comparable](src, dst []T) []T {
	m := toMap(src)
	res := make([]T, 0, len(m))
	for _, v := range dst {
		if _, ok := m[v]; ok {
			res = append(res, v)
		}
	}
	return deDuplicate(res)
}
