package jit

// Ptr 返回指向 v 的指针。
func Ptr[T any](v T) *T {
	return &v
}

// DePtr 解引用指针 p。
// p 为 nil 时 panic。
func DePtr[T any](p *T) T {
	return *p
}
