package xslice

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/errs"
)

func zeroVal[T jit.RealNumber]() T {
	var zero T
	return zero
}

// Max returns the maximum value in the slice.
func Max[T jit.RealNumber](slice []T) (T, error) {
	if len(slice) == 0 {
		return zeroVal[T](), errs.ErrEmptySlice()
	}
	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > res {
			res = slice[i]
		}
	}
	return res, nil
}

// Min returns the minimum value in the slice.
func Min[T jit.RealNumber](slice []T) (T, error) {
	if len(slice) == 0 {
		return zeroVal[T](), errs.ErrEmptySlice()
	}

	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] < res {
			res = slice[i]
		}
	}
	return res, nil
}

// Sum returns the sum of the slice.
func Sum[T jit.RealNumber](slice []T) T {
	ret := zeroVal[T]()

	if len(slice) == 0 {
		return ret
	}

	for _, v := range slice {
		ret += v
	}

	return ret
}
