package xset

var _ Set[any] = (*MapSet[any])(nil)

// MapSet 为基于内建 map 实现的集合。
type MapSet[T comparable] struct {
	m map[T]struct{}
}

// NewMapSet 创建 MapSet，size 为底层 map 的预估容量。
func NewMapSet[T comparable](size int) *MapSet[T] {
	return &MapSet[T]{
		m: make(map[T]struct{}, size),
	}
}

// Add 添加元素 ( 已存在时不产生变化 )。
func (s *MapSet[T]) Add(key T) {
	s.m[key] = struct{}{}
}

// Del 删除元素 ( 不存在时不产生变化 )。
func (s *MapSet[T]) Del(key T) {
	delete(s.m, key)
}

// Exist 判断元素是否存在。
func (s *MapSet[T]) Exist(key T) bool {
	_, ok := s.m[key]
	return ok
}

// Elems 返回所有元素组成的切片。
func (s *MapSet[T]) Elems() []T {
	keys := make([]T, 0, len(s.m))
	for key := range s.m {
		keys = append(keys, key)
	}
	return keys
}

// Size 返回元素个数。
func (s *MapSet[T]) Size() int {
	return len(s.m)
}
