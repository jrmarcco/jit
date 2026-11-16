// Package xslice 提供泛型切片工具:
// 元素查找、集合运算 ( 并/交/差/对称差 )、映射与聚合。
package xslice

// eqFunc 比较两个元素是否相等的函数。
type eqFunc[T any] func(src, dst T) bool

// matchFunc 判断元素是否满足条件的函数。
type matchFunc[T any] func(T) bool
