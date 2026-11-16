package xslice

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/errs"
)

// Max 返回切片中的最大值。
// 切片为空时返回包装了 errs.ErrEmptySlice 的错误。
func Max[T jit.RealNumber](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, errs.NewErrEmptySlice()
	}
	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > res {
			res = slice[i]
		}
	}
	return res, nil
}

// Min 返回切片中的最小值。
// 切片为空时返回包装了 errs.ErrEmptySlice 的错误。
func Min[T jit.RealNumber](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, errs.NewErrEmptySlice()
	}

	res := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] < res {
			res = slice[i]
		}
	}
	return res, nil
}

// Sum 返回切片元素之和。
// 空切片的和为 T 的零值。
func Sum[T jit.RealNumber](slice []T) T {
	var ret T

	if len(slice) == 0 {
		return ret
	}

	for _, v := range slice {
		ret += v
	}

	return ret
}
