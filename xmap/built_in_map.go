package xmap

var _ imap[any, any] = (*builtInMap[any, any])(nil)

type builtInMap[K comparable, V any] struct {
	data map[K]V
}

func newBuiltInMap[K comparable, V any](data map[K]V) *builtInMap[K, V] {
	return &builtInMap[K, V]{data: data}
}

func (m *builtInMap[K, V]) Size() int64 {
	return int64(len(m.data))
}

func (m *builtInMap[K, V]) Keys() []K {
	return Keys(m.data)
}

func (m *builtInMap[K, V]) Vals() []V {
	return Vals(m.data)
}

func (m *builtInMap[K, V]) Put(key K, val V) error {
	m.data[key] = val
	return nil
}

func (m *builtInMap[K, V]) Del(key K) (V, bool) {
	val, ok := m.data[key]
	delete(m.data, key)
	return val, ok
}

func (m *builtInMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *builtInMap[K, V]) Iter(visitFunc func(key K, val V) bool) {
	for k, v := range m.data {
		if !visitFunc(k, v) {
			break
		}
	}
}
