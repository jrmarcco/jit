package xmap

import "github.com/jrmarcco/jit/xsync"

type Hashable interface {
	Hash() uint64
	Equals(key any) bool
}

type node[K Hashable, V any] struct {
	key  K
	val  V
	next *node[K, V]
}

func (n *node[K, V]) reset() {
	var key K
	var val V

	n.key = key
	n.val = val
	n.next = nil
}

var _ imap[Hashable, any] = (*HashMap[Hashable, any])(nil)

type HashMap[K Hashable, V any] struct {
	m    map[uint64]*node[K, V] // each element is a hash bucket(node chain).
	pool *xsync.Pool[*node[K, V]]
}

func (h *HashMap[K, V]) newNode(key K, val V) *node[K, V] {
	n := h.pool.Get()
	n.val = val
	n.key = key

	return n
}

func (h *HashMap[K, V]) Size() int64 {
	return int64(len(h.m))
}

func (h *HashMap[K, V]) Keys() []K {
	res := make([]K, 0)
	for _, headN := range h.m {
		currN := headN
		for currN != nil {
			res = append(res, currN.key)
			currN = currN.next
		}
	}
	return res
}

func (h *HashMap[K, V]) Vals() []V {
	res := make([]V, 0)
	for _, headN := range h.m {
		currN := headN
		for currN != nil {
			res = append(res, currN.val)
			currN = currN.next
		}
	}
	return res
}

func (h *HashMap[K, V]) Put(key K, val V) error {
	hash := key.Hash()
	headN, ok := h.m[hash]
	if !ok {
		newNode := h.newNode(key, val)
		h.m[hash] = newNode
		return nil
	}

	preN := headN
	for headN != nil {
		if headN.key.Equals(key) {
			headN.val = val
			return nil
		}
		// find the next node in the hash bucket
		preN = headN
		headN = headN.next
	}

	newNode := h.newNode(key, val)
	preN.next = newNode
	return nil
}

func (h *HashMap[K, V]) Del(key K) (V, bool) {
	headN, ok := h.m[key.Hash()]
	if !ok {
		var zero V
		return zero, false
	}

	preN := headN
	level := 0
	for headN != nil {
		if headN.key.Equals(key) {
			switch {
			//  level == 0 means that the node is the root node of the hash bucket
			case level == 0 && headN.next == nil:
				// the hash bucket only has one node(root node)
				delete(h.m, key.Hash())
			case level == 0 && headN.next != nil:
				h.m[key.Hash()] = headN.next
			default:
				preN.next = headN.next
			}

			v := headN.val
			headN.reset()
			h.pool.Put(headN)
			return v, true
		}

		// find the next node in the hash bucket
		level++
		preN = headN
		headN = headN.next
	}

	var zero V
	return zero, false
}

func (h *HashMap[K, V]) Get(key K) (V, bool) {
	hash := key.Hash()
	if headN, ok := h.m[hash]; ok {
		for headN != nil {
			if headN.key.Equals(key) {
				return headN.val, true
			}
			// find the next node in the hash bucket
			headN = headN.next
		}
	}

	var zero V
	return zero, false
}

func (h *HashMap[K, V]) Iter(visitFunc func(key K, val V) bool) {
	for _, headN := range h.m {
		currN := headN
		for ; currN != nil; currN = currN.next {
			if !visitFunc(currN.key, currN.val) {
				return
			}
		}
	}
}

func NewHashMap[K Hashable, V any](size int) *HashMap[K, V] {
	return &HashMap[K, V]{
		m:    make(map[uint64]*node[K, V], size),
		pool: xsync.NewPool[*node[K, V]](func() *node[K, V] { return &node[K, V]{} }),
	}
}
