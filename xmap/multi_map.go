package xmap

import "github.com/jrmarcco/jit"

// MultiMap 为一 key 多值的 map，内部由 imap 实现承载。
type MultiMap[K any, V any] struct {
	m imap[K, []V]
}

// NewMultiTreeMap 创建基于 TreeMap 的 MultiMap，cmp 为 key 比较器，不可为 nil。
func NewMultiTreeMap[K comparable, V any](cmp jit.Comparator[K]) (*MultiMap[K, V], error) {
	treeMap, err := NewTreeMap[K, []V](cmp)
	if err != nil {
		return nil, err
	}
	return &MultiMap[K, V]{
		m: treeMap,
	}, nil
}

// NewMultiHashMap 创建基于 HashMap 的 MultiMap，size 为内部 map 的预估容量。
// 返回的 error 当前恒为 nil，保留在签名中以为未来扩展留余地。
func NewMultiHashMap[K Hashable, V any](size int) (*MultiMap[K, V], error) {
	var m imap[K, []V] = NewHashMap[K, []V](size)
	return &MultiMap[K, V]{
		m: m,
	}, nil
}

// Size 返回 key 的个数 ( 不展开多值 )。
func (m *MultiMap[K, V]) Size() int64 {
	return m.m.Size()
}

// Keys 返回所有 key 组成的切片。
func (m *MultiMap[K, V]) Keys() []K {
	return m.m.Keys()
}

// Vals 返回每个 key 对应的值列表，内层切片为副本。
func (m *MultiMap[K, V]) Vals() [][]V {
	vals := m.m.Vals()
	res := make([][]V, 0, len(vals))
	for i := range vals {
		res = append(res, append([]V{}, vals[i]...))
	}
	return res
}

// Put 向 key 追加一个值。
func (m *MultiMap[K, V]) Put(key K, val V) error {
	return m.PutMany(key, val)
}

// PutMany 向 key 追加一个或多个值。
func (m *MultiMap[K, V]) PutMany(key K, vals ...V) error {
	val, ok := m.Get(key)
	if !ok {
		val = make([]V, 0, len(vals))
	}

	val = append(val, vals...)
	return m.m.Put(key, val)
}

// PuyMany 向 key 追加一个或多个值。
//
// Deprecated: PuyMany 为 PutMany 的拼写错误，请使用 PutMany，
// 将在未来版本移除。
func (m *MultiMap[K, V]) PuyMany(key K, vals ...V) error {
	return m.PutMany(key, vals...)
}

// Del 删除 key 及其全部值并返回值列表，key 不存在时返回 nil 与 false。
func (m *MultiMap[K, V]) Del(key K) ([]V, bool) {
	return m.m.Del(key)
}

// Get 返回 key 对应的值列表 ( 副本 )，key 不存在时返回 nil 与 false。
func (m *MultiMap[K, V]) Get(key K) ([]V, bool) {
	if val, ok := m.m.Get(key); ok {
		return append([]V{}, val...), true
	}
	return nil, false
}

// Iter 逐值遍历 ( 同一 key 的多个值按顺序连续访问 )，
// visitFunc 返回 false 时提前终止。
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
