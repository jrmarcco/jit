package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jrmarcco/jit/internal/errs"
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
			name:    "delete from empty slice",
			slice:   []int{},
			index:   0,
			wantRes: []int{},
			wantErr: errs.ErrIndexOutOfBounds(0, 0),
		}, {
			name:    "delete from non-empty slice at index out of bounds",
			slice:   []int{1, 2, 3},
			index:   4,
			wantRes: []int{1, 2, 3},
			wantErr: errs.ErrIndexOutOfBounds(3, 4),
		}, {
			name:    "delete from non-empty slice at index negative",
			slice:   []int{1, 2, 3},
			index:   -1,
			wantRes: []int{1, 2, 3},
			wantErr: errs.ErrIndexOutOfBounds(3, -1),
		}, {
			name:    "delete from non-empty slice at index start",
			slice:   []int{1, 2, 3},
			index:   0,
			wantRes: []int{2, 3},
			wantErr: nil,
		}, {
			name:    "delete from non-empty slice at index middle",
			slice:   []int{1, 2, 3},
			index:   1,
			wantRes: []int{1, 3},
			wantErr: nil,
		}, {
			name:    "delete from non-empty slice at index end",
			slice:   []int{1, 2, 3},
			index:   2,
			wantRes: []int{1, 2},
			wantErr: nil,
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
