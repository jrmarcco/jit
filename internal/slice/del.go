package slice

import "github.com/jrmarcco/jit/internal/errs"

// Del 删除 index 位置的元素:
// 尾部元素左移一位，并将腾出的尾部槽位清零
// ( 指针元素不残留引用，与 slices.Delete 对齐 )，
// 返回共享底层数组的缩短切片。越界时返回错误。
func Del[T any](slice []T, index int) ([]T, error) {
	length := len(slice)

	if index < 0 || index >= length {
		return nil, errs.NewErrIndexOutOfBounds(length, index)
	}

	for i := index; i < length-1; i++ {
		slice[i] = slice[i+1]
	}

	// 清零腾出的尾部槽位以释放引用。
	var zero T
	slice[length-1] = zero

	return slice[:length-1], nil
}
