package slice

import (
	"testing"

	"github.com/jrmarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		index   int
		item    int
		wantRes []int
		wantErr error
	}{
		{
			name:    "add to non-empty slice at index out of bounds",
			slice:   []int{1, 2, 3},
			index:   4,
			item:    0,
			wantRes: nil,
			wantErr: errs.ErrIndexOutOfBounds(3, 4),
		}, {
			name:    "add to non-empty slice at index negative",
			slice:   []int{1, 2, 3},
			index:   -1,
			item:    0,
			wantRes: nil,
			wantErr: errs.ErrIndexOutOfBounds(3, -1),
		}, {
			name:    "add to empty slice",
			slice:   []int{},
			index:   0,
			item:    1,
			wantRes: []int{1},
			wantErr: nil,
		}, {
			name:    "add to non-empty slice at index start",
			slice:   []int{1, 2, 3},
			index:   0,
			item:    0,
			wantRes: []int{0, 1, 2, 3},
			wantErr: nil,
		}, {
			name:    "add to non-empty slice at index middle",
			slice:   []int{1, 2, 3},
			index:   1,
			item:    0,
			wantRes: []int{1, 0, 2, 3},
			wantErr: nil,
		}, {
			name:    "add to non-empty slice at index end",
			slice:   []int{1, 2, 3},
			index:   3,
			item:    0,
			wantRes: []int{1, 2, 3, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := Add(tc.slice, tc.index, tc.item)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}
