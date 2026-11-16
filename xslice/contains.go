package xslice

import "slices"

// ContainsAnyFunc 判断 slice 是否包含 elems 中的任意元素，相等性由 eq 定义。
func ContainsAnyFunc[T any](slice, elems []T, eq eqFunc[T]) bool {
	for _, e := range elems {
		if slices.ContainsFunc(slice, func(t T) bool { return eq(t, e) }) {
			return true
		}
	}
	return false
}

// ContainsAny 判断 slice 是否包含 elems 中的任意元素。
func ContainsAny[T comparable](slice, elems []T) bool {
	srcMap := toMap(slice)
	for _, e := range elems {
		if _, ok := srcMap[e]; ok {
			return true
		}
	}
	return false
}

// ContainsAllFunc 判断 slice 是否包含 elems 中的全部元素，相等性由 eq 定义。
// elems 为空 ( 含 nil ) 时恒为 true ( 空集真空包含 )。
func ContainsAllFunc[T any](slice, elems []T, eq eqFunc[T]) bool {
	for _, e := range elems {
		if !slices.ContainsFunc(slice, func(t T) bool { return eq(e, t) }) {
			return false
		}
	}
	return true
}

// ContainsAll 判断 slice 是否包含 elems 中的全部元素。
// elems 为空 ( 含 nil ) 时恒为 true ( 空集真空包含 )。
func ContainsAll[T comparable](slice, elems []T) bool {
	srcMap := toMap(slice)
	for _, e := range elems {
		if _, ok := srcMap[e]; !ok {
			return false
		}
	}
	return true
}
