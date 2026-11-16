package xslice

import "slices"

// SymmDiffSetFunc 返回两个切片的对称差 ( 相等性由 eq 定义 )。
// 结果保持 src 在前、dst 独有元素在后的顺序。
func SymmDiffSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := []T{}

	// 找出 src 中不在 dst 内的元素。
	for _, v := range src {
		if !slices.ContainsFunc(dst, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}

	// 找出 dst 中不在 src 内的元素。
	for _, v := range dst {
		if !slices.ContainsFunc(src, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}

	return deDuplicateFunc(res, eq)
}

// SymmDiffSet 返回两个切片的对称差。
// 返回元素的顺序不做保证。
func SymmDiffSet[T comparable](src, dst []T) []T {
	srcMap, dstMap := toMap(src), toMap(dst)

	for k := range dstMap {
		if _, ok := srcMap[k]; ok {
			delete(srcMap, k)
			continue
		}

		srcMap[k] = struct{}{}
	}

	res := make([]T, 0, len(srcMap))
	for k := range srcMap {
		res = append(res, k)
	}

	return res
}
