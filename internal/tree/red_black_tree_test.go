package tree

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jrmarcco/jit"
	"github.com/jrmarcco/jit/internal/errs"
)

//nolint:unparam // 测试方法
func benchmarkShuffledKeys(size int) []int {
	keys := make([]int, size)
	for i := range size {
		keys[i] = i + 1
	}

	//nolint:gosec // 测试用例，不需要使用更安全的随机数生成器。
	r := rand.New(rand.NewSource(42))
	r.Shuffle(size, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}

var testCmp = func() jit.Comparator[int] {
	return func(a, b int) int { return a - b }
}()

func validRBTree[K any, V any](root *rbNode[K, V]) bool {
	if root.getColor() != black {
		return false
	}

	// count the number of black nodes on the path from the root to the farthest leaf
	cnt := 0
	num := 0
	node := root

	// count the black nodes on the path to the leftmost leaf
	for node != nil {
		if node.getColor() == black {
			cnt++
		}
		node = node.left
	}

	return validRBNode(root, cnt, num)
}

func validRBNode[K, V any](node *rbNode[K, V], cnt, num int) bool {
	if node == nil {
		return true
	}

	// red node with red parent is invalid
	if node.getColor() == red && node.parent.getColor() == red {
		return false
	}

	if node.getColor() == black {
		num++
	}

	if node.left == nil && node.right == nil {
		// leaf node
		if num != cnt {
			return false
		}
	}

	return validRBNode(node.left, cnt, num) && validRBNode(node.right, cnt, num)
}

func TestNewRBTree(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		cmp     jit.Comparator[int]
		wantRes bool
	}{
		{
			name:    "int cmp",
			cmp:     testCmp,
			wantRes: true,
		}, {
			name:    "nil cmp",
			cmp:     nil,
			wantRes: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rbt := NewRBTree[int, string](tc.cmp)
			assert.Equal(t, tc.wantRes, validRBTree(rbt.root))
		})
	}
}

func TestRBTree_ValidateRBTree(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		node    *rbNode[int, int]
		wantRes bool
	}{
		{
			name:    "nil",
			node:    nil,
			wantRes: true,
		}, {
			name:    "root with color black",
			node:    &rbNode[int, int]{left: nil, right: nil, color: black},
			wantRes: true,
		}, {
			name:    "root with color red",
			node:    &rbNode[int, int]{left: nil, right: nil, color: red},
			wantRes: false,
		}, {
			name: "root with one child",
			node: &rbNode[int, int]{
				left: &rbNode[int, int]{
					right: nil,
					left:  nil,
					color: red,
				},
				right: nil,
				color: black,
			},
			wantRes: true,
		}, {
			name: "root with two children",
			node: &rbNode[int, int]{
				left: &rbNode[int, int]{
					right: nil,
					left:  nil,
					color: red,
				},
				right: &rbNode[int, int]{
					right: nil,
					left:  nil,
					color: black,
				},
				color: black,
			},
			wantRes: false,
		}, {
			name: "root with grandson (single red node child)",
			node: &rbNode[int, int]{
				left: &rbNode[int, int]{
					right: &rbNode[int, int]{
						right: nil,
						left:  nil,
						color: red,
					},
					left:  nil,
					color: black,
				},
				right: &rbNode[int, int]{
					right: nil,
					left: &rbNode[int, int]{
						right: nil,
						left:  nil,
						color: red,
					},
					color: black,
				},
				color: black,
			},
			wantRes: true,
		}, {
			name: "root with grandson (full red node children)",
			node: &rbNode[int, int]{
				parent: nil,
				key:    7,
				left: &rbNode[int, int]{
					key:   5,
					color: black,
					left: &rbNode[int, int]{
						key:   4,
						color: red,
					},
					right: &rbNode[int, int]{
						key:   6,
						color: red,
					},
				},
				right: &rbNode[int, int]{
					key:   10,
					color: red,
					left: &rbNode[int, int]{
						key:   9,
						color: black,
						left: &rbNode[int, int]{
							key:   8,
							color: red,
						},
					},
					right: &rbNode[int, int]{
						key:   12,
						color: black,
						left: &rbNode[int, int]{
							key:   11,
							color: red,
						},
					},
				},
				color: black,
			},
			wantRes: true,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantRes, validRBTree(tc.node))
		})
	}
}

func TestRBTree_Insert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		putNodes []*rbNode[int, int]
		wantRes  bool
		wantErr  error
		wantSize int64
		wantKeys []int
		wantVals []int
	}{
		{
			name:     "insert one node(insert root node)",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 1,
			wantKeys: []int{1},
			wantVals: []int{1},
		}, {
			name:     "insert two nodes(insert to black parent node)",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 2,
			wantKeys: []int{1, 2},
			wantVals: []int{1, 2},
		}, {
			name:     "insert multi nodes",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 5,
			wantKeys: []int{1, 2, 3, 4, 5},
			wantVals: []int{1, 2, 3, 4, 5},
		}, {
			name:     "insert multi desc order nodes",
			putNodes: []*rbNode[int, int]{{key: 5, val: 5}, {key: 4, val: 4}, {key: 3, val: 3}, {key: 2, val: 2}, {key: 1, val: 1}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 5,
			wantKeys: []int{1, 2, 3, 4, 5},
			wantVals: []int{1, 2, 3, 4, 5},
		}, {
			name:     "insert multi disorder nodes",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 3, val: 3}, {key: 2, val: 2}, {key: 4, val: 4}, {key: 5, val: 5}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 5,
			wantKeys: []int{1, 2, 3, 4, 5},
			wantVals: []int{1, 2, 3, 4, 5},
		}, {
			name:     "insert same key",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 1, val: 2}},
			wantErr:  errs.ErrSameRBNode,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rbt := NewRBTree[int, int](testCmp)
			for _, node := range tc.putNodes {
				err := rbt.Put(node.key, node.val)
				if err != nil {
					assert.Equal(t, tc.wantErr, err)
					return
				}
			}

			assert.Equal(t, tc.wantRes, validRBTree(rbt.root))
			assert.Equal(t, tc.wantSize, rbt.Size())

			keys, vals := rbt.Kvs()
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.wantVals, vals)
		})
	}
}

func TestRBTree_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		putNodes []*rbNode[int, int]
		delNodes []*rbNode[int, int]
		wantRes  bool
		wantErr  error
		wantSize int64
		wantVals []int
	}{
		{
			name:     "del first node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 1}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 4,
			wantVals: []int{2, 3, 4, 5},
		}, {
			name:     "del last node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 5}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 4,
			wantVals: []int{1, 2, 3, 4},
		}, {
			name:     "del root node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 2}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 4,
			wantVals: []int{1, 3, 4, 5},
		}, {
			name:     "del middle node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 3}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 4,
			wantVals: []int{1, 2, 4, 5},
		}, {
			name:     "del multi nodes",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 2}, {key: 3}, {key: 5}},
			wantRes:  true,
			wantErr:  nil,
			wantSize: 2,
			wantVals: []int{1, 4},
		}, {
			name:     "del non-existent node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			delNodes: []*rbNode[int, int]{{key: 6}},
			wantErr:  errs.ErrNodeNotFound,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rbt := NewRBTree[int, int](testCmp)
			for _, node := range tc.putNodes {
				_ = rbt.Put(node.key, node.val)
			}

			for _, node := range tc.delNodes {
				_, err := rbt.Del(node.key)
				if err != nil {
					assert.Equal(t, tc.wantErr, err)
					return
				}

				assert.Equal(t, tc.wantRes, validRBTree(rbt.root))
			}

			vals := rbt.Vals()
			assert.Equal(t, tc.wantVals, vals)
			assert.Equal(t, tc.wantSize, rbt.Size())
		})
	}
}

func TestRBTree_Set(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		putNodes []*rbNode[int, int]
		setNodes []*rbNode[int, int]
		wantVals []int
		wantErr  error
	}{
		{
			name:     "set one node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}},
			setNodes: []*rbNode[int, int]{{key: 1, val: 2}},
			wantVals: []int{2},
			wantErr:  nil,
		}, {
			name:     "set multi nodes",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			setNodes: []*rbNode[int, int]{{key: 1, val: 2}, {key: 2, val: 3}},
			wantVals: []int{2, 3, 3, 4, 5},
			wantErr:  nil,
		}, {
			name:     "set non-existent node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}, {key: 3, val: 3}, {key: 4, val: 4}, {key: 5, val: 5}},
			setNodes: []*rbNode[int, int]{{key: 6, val: 6}},
			wantErr:  errs.ErrNodeNotFound,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rbt := NewRBTree[int, int](testCmp)
			for _, node := range tc.putNodes {
				_ = rbt.Put(node.key, node.val)
			}

			for _, node := range tc.setNodes {
				err := rbt.Set(node.key, node.val)
				if err != nil {
					assert.Equal(t, tc.wantErr, err)
					return
				}
			}

			vals := rbt.Vals()
			assert.Equal(t, tc.wantVals, vals)
		})
	}
}

func TestRBTree_Get(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		putNodes []*rbNode[int, int]
		key      int
		wantVal  int
		wantErr  error
	}{
		{
			name:     "basic",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}},
			key:      1,
			wantVal:  1,
			wantErr:  nil,
		}, {
			name:     "non-existent node",
			putNodes: []*rbNode[int, int]{{key: 1, val: 1}, {key: 2, val: 2}},
			key:      3,
			wantVal:  0,
			wantErr:  errs.ErrNodeNotFound,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rbt := NewRBTree[int, int](testCmp)
			for _, node := range tc.putNodes {
				_ = rbt.Put(node.key, node.val)
			}

			val, err := rbt.Get(tc.key)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestRBTree_Upsert(t *testing.T) {
	t.Parallel()

	rbt := NewRBTree[int, int](testCmp)
	rbt.Upsert(1, 1)
	rbt.Upsert(1, 2)
	rbt.Upsert(2, 3)

	assert.Equal(t, int64(2), rbt.Size())
	keys, vals := rbt.Kvs()
	assert.Equal(t, []int{1, 2}, keys)
	assert.Equal(t, []int{2, 3}, vals)
	assert.True(t, validRBTree(rbt.root))
}

func TestRBTree_RandomOperations(t *testing.T) {
	t.Parallel()

	//nolint:gosec // 测试用例，不需要使用更安全的随机数生成器。
	r := rand.New(rand.NewSource(12345))
	rbt := NewRBTree[int, int](testCmp)
	expected := make(map[int]int)

	const (
		opPut = iota
		opSet
		opUpsert
		opGet
	)

	for i := 0; i < 5_000; i++ {
		key := r.Intn(300)
		val := r.Intn(10_000)

		switch r.Intn(4) {
		case opPut:
			err := rbt.Put(key, val)
			_, exists := expected[key]
			if exists {
				assert.Equal(t, errs.ErrSameRBNode, err)
			} else {
				assert.NoError(t, err)
				expected[key] = val
			}
		case opSet:
			err := rbt.Set(key, val)
			_, exists := expected[key]
			if exists {
				assert.NoError(t, err)
				expected[key] = val
			} else {
				assert.Equal(t, errs.ErrNodeNotFound, err)
			}
		case opUpsert:
			rbt.Upsert(key, val)
			expected[key] = val
		case opGet:
			got, err := rbt.Get(key)
			expVal, exists := expected[key]
			if exists {
				assert.NoError(t, err)
				assert.Equal(t, expVal, got)
			} else {
				assert.Equal(t, errs.ErrNodeNotFound, err)
			}
		}

		assert.Equal(t, int64(len(expected)), rbt.Size())
	}

	keys := make([]int, 0, len(expected))
	for k := range expected {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	wantVals := make([]int, 0, len(keys))
	for _, k := range keys {
		wantVals = append(wantVals, expected[k])
	}

	gotKeys, gotVals := rbt.Kvs()
	assert.Equal(t, keys, gotKeys)
	assert.Equal(t, wantVals, gotVals)
}

func BenchmarkRBTree_Put(b *testing.B) {
	keys := benchmarkShuffledKeys(10_000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbt := NewRBTree[int, int](testCmp)
		for _, key := range keys {
			_ = rbt.Put(key, key)
		}
	}
}

func BenchmarkRBTree_PutOrdered(b *testing.B) {
	keys := make([]int, 10_000)
	for i := range keys {
		keys[i] = i + 1
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbt := NewRBTree[int, int](testCmp)
		for _, key := range keys {
			_ = rbt.Put(key, key)
		}
	}
}

func BenchmarkRBTree_PutReverseOrdered(b *testing.B) {
	keys := make([]int, 10_000)
	for i := range keys {
		keys[i] = len(keys) - i
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbt := NewRBTree[int, int](testCmp)
		for _, key := range keys {
			_ = rbt.Put(key, key)
		}
	}
}

func BenchmarkRBTree_Get(b *testing.B) {
	keys := benchmarkShuffledKeys(10_000)
	rbt := NewRBTree[int, int](testCmp)
	for _, key := range keys {
		_ = rbt.Put(key, key)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		_, _ = rbt.Get(key)
	}
}

func BenchmarkRBTree_Del(b *testing.B) {
	keys := benchmarkShuffledKeys(10_000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbt := NewRBTree[int, int](testCmp)
		for _, key := range keys {
			_ = rbt.Put(key, key)
		}
		for _, key := range keys {
			_, _ = rbt.Del(key)
		}
	}
}

func BenchmarkRBTree_UpsertExisting(b *testing.B) {
	keys := benchmarkShuffledKeys(10_000)
	rbt := NewRBTree[int, int](testCmp)
	for _, key := range keys {
		_ = rbt.Put(key, key)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		rbt.Upsert(key, i)
	}
}
