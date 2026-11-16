package xset

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/xmap"
)

var _ Set[any] = (*TreeSet[any])(nil)

// TreeSet 为基于红黑树实现的有序集合，Elems 返回元素按比较器升序。
type TreeSet[T any] struct {
	tm *xmap.TreeMap[T, struct{}]
}

// NewTreeSet 创建 TreeSet，cmp 为元素比较器，不可为 nil。
func NewTreeSet[T any](cmp jit.Comparator[T]) (*TreeSet[T], error) {
	tm, err := xmap.NewTreeMap[T, struct{}](cmp)
	if err != nil {
		return nil, err
	}

	return &TreeSet[T]{tm: tm}, nil
}

// Add 添加元素 ( 已存在时不产生变化 )。
func (s *TreeSet[T]) Add(elem T) {
	_ = s.tm.Put(elem, struct{}{})
}

// Del 删除元素 ( 不存在时不产生变化 )。
func (s *TreeSet[T]) Del(elem T) {
	_, _ = s.tm.Del(elem)
}

// Exist 判断元素是否存在。
func (s *TreeSet[T]) Exist(elem T) bool {
	_, ok := s.tm.Get(elem)
	return ok
}

// Elems 返回按比较器升序的全部元素。
func (s *TreeSet[T]) Elems() []T {
	return s.tm.Keys()
}

// Size 返回元素个数。
func (s *TreeSet[T]) Size() int {
	return int(s.tm.Size())
}
