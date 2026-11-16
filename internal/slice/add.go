package slice

import (
	"github.com/jrmarcco/jit/internal/errs"
)

// Add 在 slice 的 index 位置插入 item。
// 与 slices.Insert 语义一致 ( 容量足够时原地搬移 )，越界时返回错误而不是 panic。
func Add[T any](slice []T, index int, item T) ([]T, error) {
	length := len(slice)

	if index < 0 || index > length {
		return nil, errs.NewErrIndexOutOfBounds(length, index)
	}

	if index == length {
		return append(slice, item), nil
	}

	// 扩容一个位置，长度 + 1。
	var zeroVal T
	slice = append(slice, zeroVal)

	// 注意此时 length 是扩容后的长度，无需再减 1。
	// i 从 length 起且恒大于 index ( >= 0 )，因此 i-1 恒为合法下标。
	for i := length; i > index; i-- {
		slice[i] = slice[i-1]
	}

	slice[index] = item

	return slice, nil
}
