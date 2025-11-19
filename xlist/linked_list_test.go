package xlist

import (
	"fmt"
	"testing"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkedList(t *testing.T) {
	t.Parallel()

	ll := LinkedListOf([]int{1, 2, 3})

	var err error
	err = ll.Append(6, 5, 4)
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3, 6, 5, 4}, ll.ToSlice())
	assert.Equal(t, 6, ll.Len())

	err = ll.Insert(0, -1)
	require.NoError(t, err)

	err = ll.Insert(0, -2)
	require.NoError(t, err)

	assert.Equal(t, []int{-2, -1, 1, 2, 3, 6, 5, 4}, ll.ToSlice())
}

func TestLinkedList_Insert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		ll       *LinkedList[int]
		index    int
		value    int
		wantRes  []int
		wantErr  error
		wantSize int
	}{
		{
			name:     "basic",
			ll:       NewLinkedList[int](),
			index:    0,
			value:    1,
			wantRes:  []int{1},
			wantSize: 1,
		}, {
			name:     "insert to head",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    0,
			value:    -1,
			wantRes:  []int{-1, 1, 2, 3},
			wantSize: 4,
		}, {
			name:     "insert to tail",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    3,
			value:    4,
			wantRes:  []int{1, 2, 3, 4},
			wantSize: 4,
		}, {
			name:     "insert to middle",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    1,
			value:    -1,
			wantRes:  []int{1, -1, 2, 3},
			wantSize: 4,
		}, {
			name:     "out of bounds",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    10,
			value:    10,
			wantErr:  errs.ErrIndexOutOfBounds(3, 10),
			wantRes:  []int{1, 2, 3},
			wantSize: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.ll.Insert(tc.index, tc.value)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSize, tc.ll.Len())
			assert.Equal(t, tc.wantRes, tc.ll.ToSlice())

			if err == nil {
				val, err := tc.ll.Get(tc.index)
				assert.NoError(t, err)
				assert.Equal(t, tc.value, val)
			}
		})
	}
}

func TestLinkedList_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		ll       *LinkedList[int]
		index    int
		wantRes  []int
		wantErr  error
		wantSize int
	}{
		{
			name:     "basic",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    1,
			wantRes:  []int{1, 3},
			wantSize: 2,
		}, {
			name:     "delete head",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    0,
			wantRes:  []int{2, 3},
			wantSize: 2,
		}, {
			name:     "delete tail",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    2,
			wantRes:  []int{1, 2},
			wantSize: 2,
		}, {
			name:     "out of bounds",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    10,
			wantErr:  errs.ErrIndexOutOfBounds(3, 10),
			wantRes:  []int{1, 2, 3},
			wantSize: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.ll.Del(tc.index)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSize, tc.ll.Len())
			assert.Equal(t, tc.wantRes, tc.ll.ToSlice())
		})
	}
}

func TestLinkedList_Set(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		ll       *LinkedList[int]
		index    int
		value    int
		wantRes  []int
		wantErr  error
		wantSize int
	}{
		{
			name:     "basic",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    1,
			value:    4,
			wantRes:  []int{1, 4, 3},
			wantSize: 3,
		}, {
			name:     "set head",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    0,
			value:    4,
			wantRes:  []int{4, 2, 3},
			wantSize: 3,
		}, {
			name:     "set tail",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    2,
			value:    4,
			wantRes:  []int{1, 2, 4},
			wantSize: 3,
		}, {
			name:     "out of bounds",
			ll:       LinkedListOf([]int{1, 2, 3}),
			index:    10,
			value:    10,
			wantErr:  errs.ErrIndexOutOfBounds(3, 10),
			wantRes:  []int{1, 2, 3},
			wantSize: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.ll.Set(tc.index, tc.value)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSize, tc.ll.Len())
			assert.Equal(t, tc.wantRes, tc.ll.ToSlice())

			if err == nil {
				val, err := tc.ll.Get(tc.index)
				assert.NoError(t, err)
				assert.Equal(t, tc.value, val)
			}
		})
	}
}

func TestLinkedList_Get(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		ll      *LinkedList[int]
		index   int
		wantRes int
		wantErr error
	}{
		{
			name:    "basic",
			ll:      LinkedListOf([]int{1, 2, 3}),
			index:   1,
			wantRes: 2,
		}, {
			name:    "head",
			ll:      LinkedListOf([]int{1, 2, 3}),
			index:   0,
			wantRes: 1,
		}, {
			name:    "tail",
			ll:      LinkedListOf([]int{1, 2, 3}),
			index:   2,
			wantRes: 3,
		}, {
			name:    "out of bounds",
			ll:      LinkedListOf([]int{1, 2, 3}),
			index:   10,
			wantErr: errs.ErrIndexOutOfBounds(3, 10),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val, err := tc.ll.Get(tc.index)
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, val)
			}
		})
	}
}

func ExampleLinkedList_Iter() {
	ll := LinkedListOf([]int{1, 2, 3, 4, 5})
	_ = ll.Iter(func(idx int, val int) error {
		fmt.Printf("%d: %d\n", idx+1, val*val)
		return nil
	})
	// Output:
	// 1: 1
	// 2: 4
	// 3: 9
	// 4: 16
	// 5: 25
}
