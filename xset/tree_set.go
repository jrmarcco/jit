package xset

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/xmap"
)

var _ Set[any] = (*TreeSet[any])(nil)

type TreeSet[T any] struct {
	tm *xmap.TreeMap[T, struct{}]
}

func NewTreeSet[T any](cmp jit.Comparator[T]) (*TreeSet[T], error) {
	tm, err := xmap.NewTreeMap[T, struct{}](cmp)
	if err != nil {
		return nil, err
	}

	return &TreeSet[T]{tm: tm}, nil
}

func (s *TreeSet[T]) Add(elem T) {
	_ = s.tm.Put(elem, struct{}{})
}

func (s *TreeSet[T]) Del(elem T) {
	_, _ = s.tm.Del(elem)
}

func (s *TreeSet[T]) Exist(elem T) bool {
	_, ok := s.tm.Get(elem)
	return ok
}

func (s *TreeSet[T]) Elems() []T {
	return s.tm.Keys()
}
