package xset

var _ Set[any] = (*MapSet[any])(nil)

type MapSet[T comparable] struct {
	m map[T]struct{}
}

func NewMapSet[T comparable](size int) *MapSet[T] {
	return &MapSet[T]{
		m: make(map[T]struct{}, size),
	}
}

func (s *MapSet[T]) Add(key T) {
	s.m[key] = struct{}{}
}

func (s *MapSet[T]) Del(key T) {
	delete(s.m, key)
}

func (s *MapSet[T]) Exist(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *MapSet[T]) Elems() []T {
	keys := make([]T, 0, len(s.m))
	for key := range s.m {
		keys = append(keys, key)
	}
	return keys
}

func (s *MapSet[T]) Size() int {
	return len(s.m)
}
