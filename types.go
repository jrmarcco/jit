// Package jit 为自用的小型泛型工具集:
// 指针工具、数值约束与比较器。
package jit

// Comparator 比较 src 与 dst。
// src < dst 时返回负值，src == dst 时返回 0，src > dst 时返回正值。
type Comparator[T any] func(src, dst T) int
