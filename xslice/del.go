package xslice

import "github.com/jrmarcco/jit/internal/slice"

// Del 删除 slice 中 index 位置的元素:
// 尾部元素左移一位，腾出的尾部槽位清零
// ( 指针元素不残留引用，与 slices.Delete 对齐 )，
// 返回共享底层数组的缩短切片。越界时返回错误。
func Del[T any](src []T, index int) ([]T, error) {
	return slice.Del(src, index)
}

// FilterDel 删除 slice 中满足 filter 条件的元素，
// 返回压缩后的前缀切片，腾出的尾部槽位已清零
// ( 指针元素不残留引用，与 slices.DeleteFunc 对齐 )。
func FilterDel[T any](src []T, filter func(idx int, elem T) bool) []T {
	index := 0

	for idx := range src {
		if filter(idx, src[idx]) {
			continue
		}

		src[index] = src[idx]
		index++
	}

	// 清零腾出的尾部以释放引用。
	clear(src[index:])

	return src[:index]
}
