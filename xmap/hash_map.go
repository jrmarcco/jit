package xmap

import "github.com/jrmarcco/jit/xsync"

// Hashable 为 HashMap 的 key 约束。
// 实现方必须保证相等的 key 产生相同的 hash
// ( 相同的 hash 不要求 key 相等 )。
type Hashable interface {
	Hash() uint64
	Equals(key any) bool
}

// node 为哈希桶链表节点。
type node[K Hashable, V any] struct {
	key  K
	val  V
	next *node[K, V]
}

// reset 清零节点字段，归还对象池前调用以释放引用。
func (n *node[K, V]) reset() {
	var key K
	var val V

	n.key = key
	n.val = val
	n.next = nil
}

var _ imap[Hashable, any] = (*HashMap[Hashable, any])(nil)

// HashMap 为基于哈希桶 ( 节点链 ) 实现的 map。
//
// 注意:
//   - HashMap 非并发安全，需要外部同步 ( 如 sync.RWMutex )。
//   - 相同 hash 的 key 链在同一个桶内，
//     内部 map 的长度是桶数而不是元素数。
type HashMap[K Hashable, V any] struct {
	m    map[uint64]*node[K, V] // 每个元素为一个哈希桶 ( 节点链 )。
	size int                    // 元素总数 ( 含链上节点 )。
	pool *xsync.Pool[*node[K, V]]
}

// NewHashMap 创建 HashMap，size 为内部 map 的预估容量。
func NewHashMap[K Hashable, V any](size int) *HashMap[K, V] {
	return &HashMap[K, V]{
		m:    make(map[uint64]*node[K, V], size),
		pool: xsync.NewPool(func() *node[K, V] { return &node[K, V]{} }),
	}
}

// newNode 从对象池取出节点并填充键值。
func (h *HashMap[K, V]) newNode(key K, val V) *node[K, V] {
	n := h.pool.Get()
	n.val = val
	n.key = key

	return n
}

// Size 返回 map 中的元素总数。
func (h *HashMap[K, V]) Size() int64 {
	return int64(h.size)
}

// Keys 返回所有 key 组成的切片。
func (h *HashMap[K, V]) Keys() []K {
	res := make([]K, 0, h.size)
	for _, headN := range h.m {
		currN := headN
		for currN != nil {
			res = append(res, currN.key)
			currN = currN.next
		}
	}
	return res
}

// Vals 返回所有 value 组成的切片。
func (h *HashMap[K, V]) Vals() []V {
	res := make([]V, 0, h.size)
	for _, headN := range h.m {
		currN := headN
		for currN != nil {
			res = append(res, currN.val)
			currN = currN.next
		}
	}
	return res
}

// Put 插入或覆盖键值对: key 已存在时覆盖，否则插入。
// 当前恒返回 nil，保留 error 以为未来扩展留余地。
func (h *HashMap[K, V]) Put(key K, val V) error {
	hash := key.Hash()
	headN, ok := h.m[hash]
	if !ok {
		newNode := h.newNode(key, val)
		h.m[hash] = newNode
		h.size++
		return nil
	}

	preN := headN
	for headN != nil {
		if headN.key.Equals(key) {
			headN.val = val
			return nil
		}
		// 在哈希桶内继续向后寻找。
		preN = headN
		headN = headN.next
	}

	newNode := h.newNode(key, val)
	preN.next = newNode
	h.size++
	return nil
}

// Del 删除 key 并返回其值，key 不存在时返回零值与 false。
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
			// level == 0 表示当前节点是哈希桶的根节点。
			case level == 0 && headN.next == nil:
				// 哈希桶内只有一个节点 ( 根节点 )。
				delete(h.m, key.Hash())
			case level == 0 && headN.next != nil:
				h.m[key.Hash()] = headN.next
			default:
				preN.next = headN.next
			}

			v := headN.val
			headN.reset()
			h.pool.Put(headN)
			h.size--
			return v, true
		}

		// 在哈希桶内继续向后寻找。
		level++
		preN = headN
		headN = headN.next
	}

	var zero V
	return zero, false
}

// Get 返回 key 对应的值，key 不存在时返回零值与 false。
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	hash := key.Hash()
	if headN, ok := h.m[hash]; ok {
		for headN != nil {
			if headN.key.Equals(key) {
				return headN.val, true
			}
			// 在哈希桶内继续向后寻找。
			headN = headN.next
		}
	}

	var zero V
	return zero, false
}

// Iter 遍历所有键值对，visitFunc 返回 false 时提前终止。
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
