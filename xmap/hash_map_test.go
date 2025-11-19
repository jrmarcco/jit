package xmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Hashable = (*testKey)(nil)

type testKey struct {
	id uint64
}

func (t testKey) Hash() uint64 {
	return t.id % 10
}

func (t testKey) Equals(other any) bool {
	val, ok := other.(testKey)
	if !ok {
		return false
	}
	return t.id == val.id
}

func TestHashMap(t *testing.T) {
	t.Parallel()

	initData := []struct {
		key testKey
		val string
	}{
		{
			key: testKey{id: 1},
			val: "1",
		}, {
			key: testKey{id: 2},
			val: "2",
		}, {
			key: testKey{id: 3},
			val: "3",
		}, {
			key: testKey{id: 4},
			val: "4",
		}, {
			key: testKey{id: 11},
			val: "11",
		}, {
			key: testKey{id: 12},
			val: "12",
		}, {
			key: testKey{id: 3},
			val: "13",
		}, {
			key: testKey{id: 4},
			val: "14",
		},
	}

	hm := NewHashMap[testKey, string](16)
	for _, d := range initData {
		_ = hm.Put(d.key, d.val)
	}

	wantHm := NewHashMap[testKey, string](16)
	wantHm.m = map[uint64]*node[testKey, string]{
		1: {
			key: testKey{id: 1},
			val: "1",
			next: &node[testKey, string]{
				key: testKey{id: 11},
				val: "11",
			},
		},
		2: {
			key: testKey{id: 2},
			val: "2",
			next: &node[testKey, string]{
				key: testKey{id: 12},
				val: "12",
			},
		},
		3: {key: testKey{id: 3}, val: "13"},
		4: {key: testKey{id: 4}, val: "14"},
	}

	assert.Equal(t, hm.m, wantHm.m)
	tcs := []struct {
		name    string
		key     testKey
		wantVal string
		wantRes bool
	}{
		{
			name:    "basic get",
			key:     testKey{id: 1},
			wantVal: "1",
			wantRes: true,
		}, {
			name:    "same hash get",
			key:     testKey{id: 12},
			wantVal: "12",
			wantRes: true,
		}, {
			name:    "key not exists",
			key:     testKey{id: 5},
			wantRes: false,
		}, {
			name:    "val not exists",
			key:     testKey{id: 31},
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, ok := hm.Get(tc.key)
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
			}
		})
	}
}

func TestHashMap_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		key     testKey
		hm      *HashMap[testKey, string]
		wantHm  *HashMap[testKey, string]
		wantVal string
		wantRes bool
	}{
		{
			name: "basic",
			key:  testKey{id: 1},
			hm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				_ = hm.Put(testKey{id: 1}, "1")
				return hm
			}(),
			wantHm:  NewHashMap[testKey, string](16),
			wantVal: "1",
			wantRes: true,
		}, {
			name:    "key not exists",
			key:     testKey{id: 5},
			hm:      NewHashMap[testKey, string](16),
			wantHm:  NewHashMap[testKey, string](16),
			wantVal: "",
			wantRes: false,
		}, {
			name: "bucket with only one node",
			key:  testKey{id: 1},
			hm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				_ = hm.Put(testKey{id: 1}, "101")
				_ = hm.Put(testKey{id: 12}, "12")
				return hm
			}(),
			wantHm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				hm.m = map[uint64]*node[testKey, string]{
					2: {
						key: testKey{id: 12},
						val: "12",
					},
				}
				return hm
			}(),
			wantVal: "101",
			wantRes: true,
		}, {
			name: "bucket with multi nodes",
			key:  testKey{id: 11},
			hm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				_ = hm.Put(testKey{id: 1}, "1")
				_ = hm.Put(testKey{id: 11}, "11")
				_ = hm.Put(testKey{id: 111}, "111")
				_ = hm.Put(testKey{id: 12}, "12")
				return hm
			}(),
			wantHm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				hm.m = map[uint64]*node[testKey, string]{
					1: {
						key: testKey{id: 1},
						val: "1",
						next: &node[testKey, string]{
							key: testKey{id: 111},
							val: "111",
						},
					},
					2: {
						key: testKey{id: 12},
						val: "12",
					},
				}
				return hm
			}(),
			wantVal: "11",
			wantRes: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, ok := tc.hm.Del(tc.key)
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
				assert.Equal(t, tc.hm.m, tc.wantHm.m)
			}
		})
	}
}

func TestHashMap_Kvs(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		hm       *HashMap[testKey, string]
		wantKeys []testKey
		wantVals []string
		wantSize int64
	}{
		{
			name:     "empty",
			hm:       NewHashMap[testKey, string](16),
			wantKeys: []testKey{},
			wantVals: []string{},
			wantSize: int64(0),
		}, {
			name: "basic",
			hm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				_ = hm.Put(testKey{id: 1}, "1")
				_ = hm.Put(testKey{id: 2}, "2")
				_ = hm.Put(testKey{id: 3}, "3")
				_ = hm.Put(testKey{id: 4}, "4")
				return hm
			}(),
			wantKeys: []testKey{
				{id: uint64(1)},
				{id: uint64(2)},
				{id: uint64(3)},
				{id: uint64(4)},
			},
			wantVals: []string{
				"1",
				"2",
				"3",
				"4",
			},
			wantSize: int64(4),
		}, {
			name: "bucket with multi nodes",
			hm: func() *HashMap[testKey, string] {
				hm := NewHashMap[testKey, string](16)
				_ = hm.Put(testKey{id: 1}, "1")
				_ = hm.Put(testKey{id: 11}, "11")
				_ = hm.Put(testKey{id: 111}, "111")
				_ = hm.Put(testKey{id: 12}, "12")
				return hm
			}(),
			wantKeys: []testKey{
				{id: uint64(1)},
				{id: uint64(11)},
				{id: uint64(111)},
				{id: uint64(12)},
			},
			wantVals: []string{
				"1",
				"11",
				"111",
				"12",
			},
			wantSize: int64(2),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keys := tc.hm.Keys()
			assert.Equal(t, len(tc.wantKeys), len(keys))
			assert.ElementsMatch(t, tc.wantKeys, keys)

			vals := tc.hm.Vals()
			assert.Equal(t, len(tc.wantVals), len(vals))
			assert.ElementsMatch(t, tc.wantVals, vals)

			assert.Equal(t, tc.hm.Size(), tc.wantSize)
		})
	}
}
