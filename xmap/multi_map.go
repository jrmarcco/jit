package xmap

import "github.com/jrmarcco/jit"

type MultiMap[K any, V any] struct {
	m imap[K, []V]
}

func (m *MultiMap[K, V]) Size() int64 {
	return m.m.Size()
}

func (m *MultiMap[K, V]) Keys() []K {
	return m.m.Keys()
}

func (m *MultiMap[K, V]) Vals() [][]V {
	vals := m.m.Vals()
	res := make([][]V, 0, len(vals))
	for i := range vals {
		res = append(res, append([]V{}, vals[i]...))
	}
	return res
}

func (m *MultiMap[K, V]) Put(key K, val V) error {
	return m.PuyMany(key, val)
}

func (m *MultiMap[K, V]) PuyMany(key K, vals ...V) error {
	val, ok := m.Get(key)
	if !ok {
		val = make([]V, 0, len(vals))
	}

	val = append(val, vals...)
	return m.m.Put(key, val)
}

func (m *MultiMap[K, V]) Del(key K) ([]V, bool) {
	return m.m.Del(key)
}

func (m *MultiMap[K, V]) Get(key K) ([]V, bool) {
	if val, ok := m.m.Get(key); ok {
		return append([]V{}, val...), true
	}
	return nil, false
}

func (m *MultiMap[K, V]) Iter(visitFunc func(key K, val V) bool) {
	m.m.Iter(func(key K, val []V) bool {
		for _, v := range val {
			if !visitFunc(key, v) {
				return false
			}
		}
		return true
	})
}

func NewMultiTreeMap[K comparable, V any](cmp jit.Comparator[K]) (*MultiMap[K, V], error) {
	treeMap, err := NewTreeMap[K, []V](cmp)
	if err != nil {
		return nil, err
	}
	return &MultiMap[K, V]{
		m: treeMap,
	}, nil
}

func NewMultiHashMap[K Hashable, V any](size int) (*MultiMap[K, V], error) {
	var m imap[K, []V] = NewHashMap[K, []V](size)
	return &MultiMap[K, V]{
		m: m,
	}, nil
}
