package xslice

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
			name:    "add to non-empty slice at index 0",
			slice:   []int{1, 2, 3},
			index:   0,
			item:    0,
			wantRes: []int{0, 1, 2, 3},
			wantErr: nil,
		}, {
			name:    "add to non-empty slice at index negative",
			slice:   []int{1, 2, 3},
			index:   -1,
			item:    0,
			wantRes: []int{1, 2, 3, 0},
			wantErr: errs.ErrIndexOutOfBounds(3, -1),
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
