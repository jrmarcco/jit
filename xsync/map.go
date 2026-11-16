package xsync

import "sync"

// Map 为 sync.Map 的泛型封装。
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load 返回 key 对应的值，key 不存在时返回零值与 false。
func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := m.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}

	val, _ := v.(V)
	return val, true
}

// Store 写入键值对。
func (m *Map[K, V]) Store(key K, val V) {
	m.m.Store(key, val)
}

// LoadOrStore 返回 key 已有的值;
// key 不存在时写入 val 并返回 val 与 false。
func (m *Map[K, V]) LoadOrStore(key K, val V) (V, bool) {
	v, ok := m.m.LoadOrStore(key, val)
	if v != nil {
		val, _ := v.(V)
		return val, ok
	}

	var zero V
	return zero, ok
}

// LoadAndDelete 取出并删除 key 对应的值，key 不存在时返回零值与 false。
func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	v, ok := m.m.LoadAndDelete(key)
	if v != nil {
		val, _ := v.(V)
		return val, ok
	}

	var zero V
	return zero, ok
}

// Delete 删除 key ( 不存在时不产生变化 )。
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Range 遍历所有键值对，fn 返回 false 时提前终止。
func (m *Map[K, V]) Range(fn func(key K, val V) bool) {
	m.m.Range(func(key, val any) bool {
		var (
			k K
			v V
		)

		if val != nil {
			v, _ = val.(V)
		}
		if key != nil {
			k, _ = key.(K)
		}

		return fn(k, v)
	})
}
