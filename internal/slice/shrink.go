// Package slice 实现公共包共享的原地切片操作 ( 插入、删除、缩容 )。
package slice

// 缩容阈值与缩放系数 ( 按容量分级触发 )。
const (
	hugeCapacityThreshold = 4096
	hugeRatioThreshold    = 2.0
	hugeShrinkFactor      = 1.5

	largeCapacityThreshold = 1024
	largeRatioThreshold    = 2.0

	mediumCapacityThreshold = 256
	mediumRatioThreshold    = 2.5
	mediumShrinkFactor      = 0.625
	halfShrinkFactor        = 0.5

	smallRatioThreshold = 3.0
)

// shrink 计算缩容后的容量。
// 返回值为新容量与是否发生了缩容。
func shrink(capacity, length int) (int, bool) {
	if length == 0 || capacity == length {
		return capacity, false
	}

	// 计算容量与长度的比值。
	ratio := float32(capacity) / float32(length)

	switch {
	// 超大容量: 比值 >= 2 时，缩至新长度的 1.5 倍。
	case capacity > hugeCapacityThreshold && ratio >= hugeRatioThreshold:
		return int(float32(length) * hugeShrinkFactor), true
	// 大容量: 比值 >= 2 时，缩至原容量的 50%。
	case capacity > largeCapacityThreshold && ratio >= largeRatioThreshold:
		return int(float32(capacity) * halfShrinkFactor), true
	// 中容量: 比值 >= 2.5 时，缩至原容量的 62.5%。
	case capacity > mediumCapacityThreshold && ratio >= mediumRatioThreshold:
		return int(float32(capacity) * mediumShrinkFactor), true
	// 小容量: 比值 >= 3 时，缩至原容量的 50%。
	case ratio >= smallRatioThreshold:
		return int(float32(capacity) * halfShrinkFactor), true
	}

	return capacity, false
}

// Shrink 按需缩容:
// 容量与长度比值达到阈值时重新分配底层数组，否则原样返回。
func Shrink[T any](slice []T) []T {
	c, length := cap(slice), len(slice)

	newCap, shrunken := shrink(c, length)
	if !shrunken {
		return slice
	}

	res := make([]T, 0, newCap)
	res = append(res, slice...)

	return res
}
