package xlist

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jrmarcco/jit"
)

var testCmp = func() jit.Comparator[int] {
	return func(a, b int) int { return a - b }
}()

func TestSkipList_Insert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		list      *SkipList[int]
		val       int
		wantSlice []int
		wantSize  int
	}{
		{
			name:      "basic",
			list:      NewSkipList(testCmp),
			val:       1,
			wantSlice: []int{1},
			wantSize:  1,
		}, {
			name: "insert exists value",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       2,
			wantSlice: []int{1, 2, 2, 3},
			wantSize:  4,
		}, {
			name: "insert to head",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val: 0,
			wantSlice: []int{
				0, 1, 2, 3,
			},
			wantSize: 4,
		}, {
			name: "insert to tail",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       4,
			wantSlice: []int{1, 2, 3, 4},
			wantSize:  4,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.list.Insert(tc.val)
			assert.Equal(t, tc.list.Len(), tc.wantSize)
			assert.Equal(t, tc.list.ToSlice(), tc.wantSlice)
		})
	}
}

func TestSkipList_Delete(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		list      *SkipList[int]
		val       int
		wantSlice []int
		wantSize  int
		wantRes   bool
	}{
		{
			name: "basic",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       2,
			wantSlice: []int{1, 3},
			wantSize:  2,
			wantRes:   true,
		}, {
			name: "delete non-exist value",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       4,
			wantSlice: []int{1, 2, 3},
			wantSize:  3,
			wantRes:   false,
		}, {
			name: "delete head",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       1,
			wantSlice: []int{2, 3},
			wantSize:  2,
			wantRes:   true,
		}, {
			name: "delete tail",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			val:       3,
			wantSlice: []int{1, 2},
			wantSize:  2,
			wantRes:   true,
		}, {
			name:      "delete from empty list",
			list:      NewSkipList(testCmp),
			val:       1,
			wantSlice: []int{},
			wantSize:  0,
			wantRes:   false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := tc.list.Delete(tc.val)
			assert.Equal(t, tc.list.Len(), tc.wantSize)
			assert.Equal(t, tc.list.ToSlice(), tc.wantSlice)
			assert.Equal(t, res, tc.wantRes)
		})
	}
}

func TestSkipList_Exist(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		list    *SkipList[int]
		target  int
		wantRes bool
	}{
		{
			name: "basic",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			target:  2,
			wantRes: true,
		}, {
			name: "not exist",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			target:  4,
			wantRes: false,
		}, {
			name:    "empty list",
			list:    NewSkipList(testCmp),
			target:  1,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			exists := tc.list.Exists(tc.target)
			assert.Equal(t, exists, tc.wantRes)
		})
	}
}

func TestSkipList_GetByIndex(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		list    *SkipList[int]
		idx     int
		wantVal int
		wantRes bool
	}{
		{
			name: "basic",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     1,
			wantVal: 2,
			wantRes: true,
		}, {
			name: "not exist",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     4,
			wantVal: 0,
			wantRes: false,
		}, {
			name:    "empty list",
			list:    NewSkipList(testCmp),
			idx:     0,
			wantVal: 0,
			wantRes: false,
		}, {
			name: "head",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     0,
			wantVal: 1,
			wantRes: true,
		}, {
			name: "tail",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     2,
			wantVal: 3,
			wantRes: true,
		}, {
			name: "out of range",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     6,
			wantRes: false,
		}, {
			name: "negative index",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			idx:     -1,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, ok := tc.list.Get(tc.idx)
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
			}
		})
	}
}

func TestSkipList_Peek(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		list    *SkipList[int]
		wantVal int
		wantRes bool
	}{
		{
			name: "basic",
			list: func() *SkipList[int] {
				sl := NewSkipList(testCmp)
				sl.Insert(1)
				sl.Insert(2)
				sl.Insert(3)
				return sl
			}(),
			wantVal: 1,
			wantRes: true,
		}, {
			name:    "empty list",
			list:    NewSkipList(testCmp),
			wantVal: 0,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, ok := tc.list.Peek()
			assert.Equal(t, ok, tc.wantRes)

			if ok {
				assert.Equal(t, val, tc.wantVal)
			}
		})
	}
}
