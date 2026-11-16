// Package xmap 提供泛型 map 工具与 map 实现:
// 自定义哈希的 HashMap、红黑树 TreeMap 与一 key 多值的 MultiMap。
package xmap

// imap 为 xmap 内部各 map 实现的公共接口。
type imap[K any, V any] interface {
	Size() int64
	Keys() []K
	Vals() []V
	Put(key K, val V) error
	Del(key K) (V, bool)
	Get(key K) (V, bool)
	Iter(visitFunc func(key K, val V) bool)
}
