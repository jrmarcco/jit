package jit

// RealNumber 为数值约束，匹配任意整数与浮点数类型。
type RealNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// Number 为数值约束，在 RealNumber 基础上追加复数类型。
type Number interface {
	RealNumber | ~complex64 | ~complex128
}
