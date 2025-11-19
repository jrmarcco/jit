package xslice

import "slices"

// ContainsFunc checks if the slice contains an element that satisfies the given function.
func ContainsFunc[T comparable](slice []T, eq func(t T) bool) bool {
	return slices.ContainsFunc(slice, eq)
}

// Contains checks if the slice contains the given element.
func Contains[T comparable](slice []T, elem T) bool {
	return ContainsFunc(slice, func(t T) bool { return t == elem })
}

// ContainsAnyFunc checks if the slice contains any of the given elements.
func ContainsAnyFunc[T comparable](slice, elems []T, eq eqFunc[T]) bool {
	for _, e := range elems {
		for _, v := range slice {
			if eq(v, e) {
				return true
			}
		}
	}
	return false
}

// ContainsAny checks if the slice contains any of the given elements.
func ContainsAny[T comparable](slice, elems []T) bool {
	srcMap := toMap(slice)
	for _, e := range elems {
		if _, ok := srcMap[e]; ok {
			return true
		}
	}
	return false
}

// ContainsAllFunc checks if the slice contains all of the given elements.
func ContainsAllFunc[T comparable](slice, elems []T, eq eqFunc[T]) bool {
	if slice == nil || elems == nil {
		return false
	}

	for _, e := range elems {
		if !ContainsFunc(slice, func(t T) bool { return eq(e, t) }) {
			return false
		}
	}
	return true
}

// ContainsAll checks if the slice contains all of the given elements.
func ContainsAll[T comparable](slice, elems []T) bool {
	if slice == nil || elems == nil {
		return false
	}

	srcMap := toMap(slice)
	for _, e := range elems {
		if _, ok := srcMap[e]; !ok {
			return false
		}
	}
	return true
}
