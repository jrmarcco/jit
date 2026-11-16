package xslice

// Reverse 返回元素逆序的新切片，不修改原切片。
func Reverse[T any](slice []T) []T {
	length := len(slice)
	reversed := make([]T, length)
	for i, v := range slice {
		reversed[length-i-1] = v
	}
	return reversed
}
