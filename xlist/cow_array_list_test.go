package xlist

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jrmarcco/jit/internal/errs"
)

func TestCowArrayList_Insert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		al      *CowArrayList[int]
		index   int
		value   int
		wantRes []int
		wantErr error
		wantLen int
	}{
		{
			name:    tcBasic,
			al:      NewCowArrayList[int](4),
			index:   0,
			value:   1,
			wantLen: 1,
			wantRes: []int{1},
		}, {
			name:    tcBasicOutOfBounds,
			al:      NewCowArrayList[int](4),
			index:   1,
			value:   1,
			wantErr: errs.NewErrIndexOutOfBounds(0, 1),
			wantLen: 0,
		}, {
			name:    tcInsertToMiddleOfEmpty,
			al:      NewCowArrayList[int](4),
			index:   2,
			value:   3,
			wantErr: errs.NewErrIndexOutOfBounds(0, 2),
			wantLen: 0,
		}, {
			name:    tcInsertToHead,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			value:   0,
			wantRes: []int{0, 1, 2, 3, 4, 5},
			wantLen: 6,
		}, {
			name:    tcInsertToTail,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			value:   6,
			wantRes: []int{1, 2, 3, 4, 5, 6},
			wantLen: 6,
		}, {
			name:    tcInsertToMiddle,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			value:   3,
			wantRes: []int{1, 2, 3, 3, 4, 5},
			wantLen: 6,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.al.Insert(tc.index, tc.value)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.Len(), tc.wantLen)

			if err == nil {
				assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
			}
		})
	}
}

func TestCowArrayList_Append(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		al      *CowArrayList[int]
		values  []int
		wantRes []int
		wantLen int
	}{
		{
			name:    tcBasic,
			al:      NewCowArrayList[int](4),
			values:  []int{1, 2, 3, 4, 5},
			wantRes: []int{1, 2, 3, 4, 5},
			wantLen: 5,
		}, {
			name:    tcAppendToNonEmptyList,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			values:  []int{6, 7, 8, 9, 10},
			wantRes: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantLen: 10,
		}, {
			name:    tcAppendSingleItem,
			al:      NewCowArrayList[int](4),
			values:  []int{1},
			wantRes: []int{1},
			wantLen: 1,
		}, {
			name:    tcAppendSingleItemToNonEmpty,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			values:  []int{6},
			wantRes: []int{1, 2, 3, 4, 5, 6},
			wantLen: 6,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.al.Append(tc.values...)
			assert.NoError(t, err)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
			assert.Equal(t, tc.al.Len(), tc.wantLen)
		})
	}
}

func TestCowArrayList_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		al      *CowArrayList[int]
		index   int
		wantRes []int
		wantErr error
		wantLen int
	}{
		{
			name:    tcBasic,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			wantRes: []int{1, 2, 4, 5},
			wantLen: 4,
		}, {
			name:    tcOutOfBounds,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			wantErr: errs.NewErrIndexOutOfBounds(5, 5),
			wantRes: []int{1, 2, 3, 4, 5},
			wantLen: 5,
		}, {
			name:    tcDeleteFromHead,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantRes: []int{2, 3, 4, 5},
			wantLen: 4,
		}, {
			name:    tcDeleteFromTail,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			wantRes: []int{1, 2, 3, 4},
			wantLen: 4,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.al.Del(tc.index)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
		})
	}
}

func TestCowArrayList_Set(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		al      *CowArrayList[int]
		index   int
		value   int
		wantRes []int
		wantErr error
	}{
		{
			name:    tcBasic,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			value:   6,
			wantRes: []int{1, 2, 6, 4, 5},
		}, {
			name:    tcOutOfBounds,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			value:   6,
			wantErr: errs.NewErrIndexOutOfBounds(5, 5),
			wantRes: []int{1, 2, 3, 4, 5},
		}, {
			name:    tcSetFromHead,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			value:   6,
			wantRes: []int{6, 2, 3, 4, 5},
		}, {
			name:    tcSetFromTail,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			value:   6,
			wantRes: []int{1, 2, 3, 4, 6},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.al.Set(tc.index, tc.value)

			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, tc.al.ToSlice(), tc.wantRes)
		})
	}
}

func TestCowArrayList_Get(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		al      *CowArrayList[int]
		index   int
		wantRes int
		wantErr error
	}{
		{
			name:    tcBasic,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   2,
			wantRes: 3,
		}, {
			name:    tcOutOfBounds,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   5,
			wantErr: errs.NewErrIndexOutOfBounds(5, 5),
		}, {
			name:    tcEmptyList,
			al:      NewCowArrayList[int](4),
			index:   0,
			wantErr: errs.NewErrIndexOutOfBounds(0, 0),
		}, {
			name:    tcGetFromHead,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   0,
			wantRes: 1,
		}, {
			name:    tcGetFromTail,
			al:      CowArrayListOf([]int{1, 2, 3, 4, 5}),
			index:   4,
			wantRes: 5,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := tc.al.Get(tc.index)
			assert.Equal(t, err, tc.wantErr)

			if err == nil {
				assert.Equal(t, res, tc.wantRes)
			}
		})
	}
}

func ExampleCowArrayList_Iter() {
	al := ArrayListOf([]int{1, 2, 3, 4, 5})
	_ = al.Iter(func(idx int, val int) error {
		fmt.Printf("%d: %d\n", idx, 2*val)
		return nil
	})
	// Output:
	// 0: 2
	// 1: 4
	// 2: 6
	// 3: 8
	// 4: 10
}

// TestCowArrayList_ConcurrentAccess 验证写时复制列表在读操作
// ( Get/Iter/ToSlice/Len/Cap ) 与写操作 ( Append/Del/Set/Insert )
// 并发执行时无数据竞争。
//
// 仅在以 -race 构建测试二进制时报告竞争。
func TestCowArrayList_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	cal := CowArrayListOf([]int{1, 2, 3, 4, 5})

	const (
		writerCnt = 4
		readerCnt = 4
		ops       = 200
	)

	var wg sync.WaitGroup

	for i := range writerCnt {
		writerIdx := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < ops; j++ {
				switch (writerIdx + j) % 3 {
				case 0:
					_ = cal.Append(j)
				case 1:
					if cal.Len() > 0 {
						_ = cal.Set(cal.Len()-1, j)
					}
				default:
					if cal.Len() > 1 {
						_ = cal.Del(0)
					}
				}
			}
		}()
	}

	for i := range readerCnt {
		readerIdx := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < ops; j++ {
				switch (readerIdx + j) % 4 {
				case 0:
					if cal.Len() > 0 {
						_, _ = cal.Get(cal.Len() - 1)
					}
				case 1:
					_ = cal.Iter(func(idx int, val int) error {
						_ = idx
						_ = val
						return nil
					})
				case 2:
					_ = cal.ToSlice()
				default:
					_ = cal.Cap()
				}
			}
		}()
	}

	wg.Wait()
}
