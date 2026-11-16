package xslice

import "slices"

// Map 通过 fn 把切片映射为新切片。
func Map[Src any, Dst any](src []Src, fn func(idx int, src Src) Dst) []Dst {
	dst := make([]Dst, len(src))

	for i := range src {
		dst[i] = fn(i, src[i])
	}

	return dst
}

// MapPtr 与 Map 相同，但 fn 接收源元素的指针，
// 避免迭代大结构体时的复制开销。
func MapPtr[Src any, Dst any](src []Src, fn func(idx int, src *Src) Dst) []Dst {
	dst := make([]Dst, len(src))

	for i := range src {
		dst[i] = fn(i, &src[i])
	}

	return dst
}

// FilterMap 同时完成过滤与映射:
// fn 返回 ok 为 true 时才把结果计入返回切片。
func FilterMap[Src any, Dst any](src []Src, filter func(idx int, src Src) (Dst, bool)) []Dst {
	dst := make([]Dst, 0, len(src))

	for i := range src {
		if d, ok := filter(i, src[i]); ok {
			dst = append(dst, d)
		}
	}

	return dst
}

// FilterMapPtr 与 FilterMap 相同，但 filter 接收源元素的指针，
// 避免迭代大结构体时的复制开销。
func FilterMapPtr[Src any, Dst any](src []Src, filter func(idx int, src *Src) (Dst, bool)) []Dst {
	dst := make([]Dst, 0, len(src))

	for i := range src {
		if d, ok := filter(i, &src[i]); ok {
			dst = append(dst, d)
		}
	}

	return dst
}

// ToMap 把切片转换为 map，key 由 fn 生成，value 为元素本身。
// 重复 key 保留最后一个元素。
func ToMap[K comparable, V any](slice []V, fn func(elem V) K) map[K]V {
	res := make(map[K]V, len(slice))
	for i := range slice {
		res[fn(slice[i])] = slice[i]
	}
	return res
}

// ToMapWithVal 把切片转换为 map，key 与 value 均由 fn 生成。
// 重复 key 保留最后一对。
func ToMapWithVal[E any, K comparable, V any](slice []E, fn func(elem E) (K, V)) map[K]V {
	res := make(map[K]V, len(slice))
	for i := range slice {
		k, v := fn(slice[i])
		res[k] = v
	}
	return res
}

// deDuplicateFunc 返回去重后的切片 ( 相等性由 eq 定义，保序 )。时间复杂度 O(n²)。
func deDuplicateFunc[T comparable](slice []T, eq eqFunc[T]) []T {
	res := make([]T, 0, len(slice))
	for i := range slice {
		if !slices.ContainsFunc(slice[i+1:], func(t T) bool { return eq(slice[i], t) }) {
			res = append(res, slice[i])
		}
	}
	return res
}

// deDuplicate 返回去重后的切片 ( 不保序 )。
func deDuplicate[T comparable](slice []T) []T {
	m := toMap(slice)
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

// toMap 把切片转换为集合 ( 仅含 key 的 map )。
func toMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{}, len(slice))
	for i := range slice {
		m[slice[i]] = struct{}{}
	}
	return m
}
