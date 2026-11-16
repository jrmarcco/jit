package xslice

import "github.com/jrmarcco/jit/internal/slice"

// Add 在 slice 的 index 位置插入 item，返回插入后的切片。
// 容量足够时原地搬移 ( 复用调用方底层数组 )，越界时返回错误。
func Add[T any](src []T, index int, item T) ([]T, error) {
	return slice.Add(src, index, item)
}
