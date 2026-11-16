package xmap

import (
	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/tree"
)

var _ imap[any, any] = (*TreeMap[any, any])(nil)

// TreeMap 为基于红黑树实现的 map。
// Keys/Vals/Iter/KeyVals 返回的元素按 key 有序。
//
// TreeMap 非并发安全，需要外部同步 ( 如 sync.RWMutex )。
type TreeMap[K any, V any] struct {
	tree *tree.RBTree[K, V]
}

// NewTreeMap 创建 TreeMap，cmp 为 key 比较器，不可为 nil。
func NewTreeMap[K any, V any](cmp jit.Comparator[K]) (*TreeMap[K, V], error) {
	if cmp == nil {
		return nil, ErrNilComparator
	}
	return &TreeMap[K, V]{tree: tree.NewRBTree[K, V](cmp)}, nil
}

// NewTreeMapWithMap 创建 TreeMap 并写入 m 中的全部键值对。
func NewTreeMapWithMap[K comparable, V any](cmp jit.Comparator[K], m map[K]V) (*TreeMap[K, V], error) {
	treeMap, err := NewTreeMap[K, V](cmp)
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		if err := treeMap.Put(k, v); err != nil {
			return nil, err
		}
	}

	return treeMap, nil
}

// Size 返回元素个数。
func (tm *TreeMap[K, V]) Size() int64 {
	return tm.tree.Size()
}

// Keys 返回按 key 升序的全部 key。
func (tm *TreeMap[K, V]) Keys() []K {
	return tm.tree.Keys()
}

// Vals 返回按 key 升序对应的全部 value。
func (tm *TreeMap[K, V]) Vals() []V {
	return tm.tree.Vals()
}

// Put 插入或覆盖键值对 ( upsert 语义: key 已存在时覆盖 )。
// 当前恒返回 nil，保留 error 以为未来扩展留余地。
func (tm *TreeMap[K, V]) Put(key K, val V) error {
	tm.tree.Upsert(key, val)
	return nil
}

// Del 删除 key 并返回其值，key 不存在时返回零值与 false。
func (tm *TreeMap[K, V]) Del(key K) (V, bool) {
	v, err := tm.tree.Del(key)
	return v, err == nil
}

// Get 返回 key 对应的值，key 不存在时返回零值与 false。
func (tm *TreeMap[K, V]) Get(key K) (V, bool) {
	v, err := tm.tree.Get(key)
	return v, err == nil
}

// Iter 按 key 升序遍历，visitFunc 返回 false 时提前终止。
func (tm *TreeMap[K, V]) Iter(visitFunc func(key K, val V) bool) {
	tm.tree.Iter(visitFunc)
}

// KeyVals 返回按 key 升序且下标对齐的 keys 与 vals。
func (tm *TreeMap[K, V]) KeyVals() (keys []K, vals []V) {
	keys, vals = tm.tree.Kvs()
	return keys, vals
}
