package xslice

// UnionSetFunc 返回两个切片的并集 ( 相等性由 eq 定义 )。
// 结果保持 src 在前、dst 独有元素在后的顺序。
func UnionSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(src)+len(dst))

	res = append(res, src...)
	res = append(res, dst...)

	return deDuplicateFunc(res, eq)
}

// UnionSet 返回两个切片的并集。
// 返回元素的顺序不做保证。
func UnionSet[T comparable](src, dst []T) []T {
	srcMap, dstMap := toMap(src), toMap(dst)
	for k := range srcMap {
		dstMap[k] = struct{}{}
	}

	ret := make([]T, 0, len(dstMap))
	for k := range dstMap {
		ret = append(ret, k)
	}

	return ret
}
