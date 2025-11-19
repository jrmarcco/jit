package xslice

import (
	"testing"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestDel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		index   int
		wantRes []int
		wantErr error
	}{
		{
			name:    "delete from non-empty slice at index 0",
			slice:   []int{1, 2, 3},
			index:   0,
			wantRes: []int{2, 3},
			wantErr: nil,
		}, {
			name:    "delete from non-empty slice at index negative",
			slice:   []int{1, 2, 3},
			index:   -1,
			wantRes: []int{1, 2, 3},
			wantErr: errs.ErrIndexOutOfBounds(3, -1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := Del(tc.slice, tc.index)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestFilterDel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		filter  func(idx int, elem int) bool
		wantRes []int
	}{
		{
			name:    "delete from empty slice",
			slice:   []int{},
			filter:  func(_ int, _ int) bool { return true },
			wantRes: []int{},
		}, {
			name:    "delete nothing",
			slice:   []int{1, 2, 3},
			filter:  func(_ int, _ int) bool { return false },
			wantRes: []int{1, 2, 3},
		}, {
			name:    "delete all",
			slice:   []int{1, 2, 3},
			filter:  func(_ int, _ int) bool { return true },
			wantRes: []int{},
		}, {
			name:    "delete first element",
			slice:   []int{1, 2, 3},
			filter:  func(idx int, _ int) bool { return idx == 0 },
			wantRes: []int{2, 3},
		}, {
			name:    "delete last element",
			slice:   []int{1, 2, 3},
			filter:  func(idx int, _ int) bool { return idx == 2 },
			wantRes: []int{1, 2},
		}, {
			name:    "delete first and last element",
			slice:   []int{1, 2, 3},
			filter:  func(idx int, _ int) bool { return idx == 0 || idx == 2 },
			wantRes: []int{2},
		}, {
			name:    "delete middle element",
			slice:   []int{1, 2, 3},
			filter:  func(idx int, _ int) bool { return idx == 1 },
			wantRes: []int{1, 3},
		}, {
			name:    "delete all elements that are even",
			slice:   []int{1, 2, 3, 4, 5, 6},
			filter:  func(_ int, elem int) bool { return elem%2 == 0 },
			wantRes: []int{1, 3, 5},
		}, {
			name:    "delete first three elements",
			slice:   []int{1, 2, 3, 4, 5, 6},
			filter:  func(idx int, _ int) bool { return idx < 3 },
			wantRes: []int{4, 5, 6},
		}, {
			name:    "delete last four elements",
			slice:   []int{1, 2, 3, 4, 5, 6},
			filter:  func(idx int, _ int) bool { return idx >= 3 },
			wantRes: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := FilterDel(tc.slice, tc.filter)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
