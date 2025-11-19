package slice

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

func shrink(capacity, length int) (int, bool) {
	if length == 0 || capacity == length {
		return capacity, false
	}

	// calculate the ratio of capacity to length
	ratio := float32(capacity) / float32(length)

	switch {
	// huge capacity: when the ratio >= 2, shrink to 1.5 times of the original capacity
	case capacity > hugeCapacityThreshold && ratio >= hugeRatioThreshold:
		return int(float32(length) * hugeShrinkFactor), true
	// large capacity: when the ratio >= 2, shrink to 50% of the original capacity
	case capacity > largeCapacityThreshold && ratio >= largeRatioThreshold:
		return int(float32(capacity) * halfShrinkFactor), true
	// medium capacity: when the ratio >= 2.5, shrink to 62.5% of the original capacity
	case capacity > mediumCapacityThreshold && ratio >= mediumRatioThreshold:
		return int(float32(capacity) * mediumShrinkFactor), true
	// small capacity: when the ratio >= 3, shrink to 50% of the original capacity
	case ratio >= smallRatioThreshold:
		return int(float32(capacity) * halfShrinkFactor), true
	}

	return capacity, false
}

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
