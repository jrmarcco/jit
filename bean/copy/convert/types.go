package convert

type Convertor[S any, D any] interface {
	Convert(s S) (D, error)
}

type ConvertFunc[S any, D any] func(s S) (D, error)

func (cf ConvertFunc[S, D]) Convert(s S) (D, error) {
	return cf(s)
}
