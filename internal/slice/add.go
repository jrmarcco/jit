package slice

import (
	"github.com/jrmarcco/jit/internal/errs"
)

func Add[T any](slice []T, index int, item T) ([]T, error) {
	length := len(slice)

	if index < 0 || index > length {
		return nil, errs.ErrIndexOutOfBounds(length, index)
	}

	if index == length {
		return append(slice, item), nil
	}

	// expand one position, length + 1
	var zeroVal T
	slice = append(slice, zeroVal)

	// note that length is the length after expansion, so no need to subtract 1
	for i := length; i > index; i-- {
		if i-1 >= 0 {
			slice[i] = slice[i-1]
		}
	}

	slice[index] = item

	return slice, nil
}
