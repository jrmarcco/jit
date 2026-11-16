package xmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jrmarcco/jit"
)

// 表驱动断言间复用的样例值 ( goconst )。
const (
	val1a = "1a"
	val1b = "1b"
	val1c = "1c"
	val1d = "1d"
	val2a = "2a"
	val2b = "2b"
)

func TestMultiHashMap(t *testing.T) {
	t.Parallel()

	mm, err := NewMultiHashMap[testKey, string](8)
	require.NoError(t, err)

	// 空 map。
	_, ok := mm.Get(testKey{id: 1})
	assert.False(t, ok)
	assert.Equal(t, int64(0), mm.Size())

	// 向同一个 key 追加多个值。
	assert.NoError(t, mm.Put(testKey{id: 1}, val1a))
	assert.NoError(t, mm.Put(testKey{id: 1}, val1b))
	assert.NoError(t, mm.PutMany(testKey{id: 1}, val1c, val1d))
	assert.NoError(t, mm.Put(testKey{id: 2}, val2a))

	vals, ok := mm.Get(testKey{id: 1})
	assert.True(t, ok)
	assert.Equal(t, []string{val1a, val1b, val1c, val1d}, vals)

	vals, ok = mm.Get(testKey{id: 2})
	assert.True(t, ok)
	assert.Equal(t, []string{val2a}, vals)

	assert.Equal(t, int64(2), mm.Size())
	assert.ElementsMatch(t, []testKey{{id: 1}, {id: 2}}, mm.Keys())

	// Vals 返回内部切片的副本。
	allVals := mm.Vals()
	assert.ElementsMatch(t, [][]string{{val1a, val1b, val1c, val1d}, {val2a}}, allVals)
	allVals[0][0] = "mutated"
	got, _ := mm.Get(testKey{id: 1})
	assert.Equal(t, val1a, got[0], "Vals must return copies")

	// Get 必须返回副本: 修改返回值不影响 map 内部。
	got[0] = "mutated"
	gotAgain, _ := mm.Get(testKey{id: 1})
	assert.Equal(t, val1a, gotAgain[0])

	// Del 删除 key 及其全部值。
	deleted, ok := mm.Del(testKey{id: 1})
	assert.True(t, ok)
	assert.Equal(t, []string{val1a, val1b, val1c, val1d}, deleted)
	assert.Equal(t, int64(1), mm.Size())

	_, ok = mm.Del(testKey{id: 1})
	assert.False(t, ok)
}

func TestMultiHashMap_Iter(t *testing.T) {
	t.Parallel()

	mm, err := NewMultiHashMap[testKey, string](8)
	require.NoError(t, err)

	assert.NoError(t, mm.PutMany(testKey{id: 1}, val1a, val1b))
	assert.NoError(t, mm.PutMany(testKey{id: 2}, val2a))

	got := make(map[testKey][]string)
	mm.Iter(func(key testKey, val string) bool {
		got[key] = append(got[key], val)
		return true
	})

	assert.Equal(t, map[testKey][]string{
		{id: 1}: {val1a, val1b},
		{id: 2}: {val2a},
	}, got)
}

func TestMultiHashMap_IterBreak(t *testing.T) {
	t.Parallel()

	mm, err := NewMultiHashMap[testKey, string](8)
	require.NoError(t, err)

	assert.NoError(t, mm.PutMany(testKey{id: 1}, val1a, val1b, val1c))

	cnt := 0
	mm.Iter(func(_ testKey, _ string) bool {
		cnt++
		return cnt < 2
	})
	assert.Equal(t, 2, cnt, "iteration should stop when visitFunc returns false")
}

func TestMultiTreeMap(t *testing.T) {
	t.Parallel()

	mm, err := NewMultiTreeMap[int, string](jit.Comparator[int](func(src, dst int) int {
		switch {
		case src < dst:
			return -1
		case src > dst:
			return 1
		default:
			return 0
		}
	}))
	require.NoError(t, err)

	assert.NoError(t, mm.PutMany(2, val2a, val2b))
	assert.NoError(t, mm.PutMany(1, val1a))

	// TreeMap 的 key 有序。
	assert.Equal(t, []int{1, 2}, mm.Keys())

	vals, ok := mm.Get(2)
	assert.True(t, ok)
	assert.Equal(t, []string{val2a, val2b}, vals)
	assert.Equal(t, int64(2), mm.Size())

	_, ok = mm.Del(1)
	assert.True(t, ok)
	assert.Equal(t, []int{2}, mm.Keys())
}

// TestMultiMap_PuyManyDeprecated 验证废弃别名
// 仍委托给 PutMany。
func TestMultiMap_PuyManyDeprecated(t *testing.T) {
	t.Parallel()

	mm, err := NewMultiHashMap[testKey, string](8)
	require.NoError(t, err)

	assert.NoError(t, mm.PuyMany(testKey{id: 1}, val1a))
	vals, ok := mm.Get(testKey{id: 1})
	assert.True(t, ok)
	assert.Equal(t, []string{val1a}, vals)
}
