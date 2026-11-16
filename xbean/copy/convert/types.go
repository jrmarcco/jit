// Package convert 定义拷贝器使用的类型转换抽象。
package convert

// Convertor 为类型转换器: 把 S 类型的值转换为 D 类型。
type Convertor[S any, D any] interface {
	Convert(s S) (D, error)
}

// ConvertFunc 为 Convertor 的函数适配器。
type ConvertFunc[S any, D any] func(s S) (D, error)

// Convert 执行转换。
func (cf ConvertFunc[S, D]) Convert(s S) (D, error) {
	return cf(s)
}
