package xslice

// SymmDiffSetFunc symmetric difference of two slices
func SymmDiffSetFunc[T comparable](src, dst []T, eq eqFunc[T]) []T {
	res := []T{}

	// find elements not in src
	for _, v := range src {
		if !ContainsFunc(dst, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}

	// find elements not in dst
	for _, v := range dst {
		if !ContainsFunc(src, func(t T) bool { return eq(v, t) }) {
			res = append(res, v)
		}
	}

	return deDuplicateFunc(res, eq)
}

// SymmDiffSet symmetric difference of two slices
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
