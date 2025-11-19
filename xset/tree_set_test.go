package xset

import (
	"testing"

	"github.com/JrMarcco/jit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cmp = jit.Comparator[int](func(a, b int) int {
	return a - b
})

func TestTreeSet_Add(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		cmp       jit.Comparator[int]
		initElems []int
		wantElems []int
		wantErr   error
	}{
		{
			name:    "nil comparator",
			cmp:     nil,
			wantErr: ErrNilComparator,
		}, {
			name:      "add single element",
			cmp:       cmp,
			initElems: []int{1},
			wantElems: []int{1},
		}, {
			name:      "add multiple elements",
			cmp:       cmp,
			initElems: []int{1, 2, 3, 4},
			wantElems: []int{1, 2, 3, 4},
		}, {
			name:      "add duplicate elements",
			cmp:       cmp,
			initElems: []int{1, 2, 3, 4, 1, 2, 3, 4},
			wantElems: []int{1, 2, 3, 4},
		}, {
			name:      "add elements in disorder",
			cmp:       cmp,
			initElems: []int{4, 3, 2, 1},
			wantElems: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s, err := NewTreeSet(tc.cmp)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			for _, elem := range tc.initElems {
				s.Add(elem)
			}

			assert.ElementsMatch(t, tc.wantElems, s.Elems())
		})
	}
}

func TestTreeSet_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		initElems []int
		rmElems   []int
		wantElems []int
		wantErr   error
	}{
		{
			name:      "del non-exist element",
			initElems: []int{1, 2, 3, 4},
			rmElems:   []int{5},
			wantElems: []int{1, 2, 3, 4},
		}, {
			name:      "del exist element",
			initElems: []int{1, 2, 3, 4},
			rmElems:   []int{1},
			wantElems: []int{2, 3, 4},
		}, {
			name:      "del multiple elements",
			initElems: []int{1, 2, 3, 4},
			rmElems:   []int{1, 2},
			wantElems: []int{3, 4},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s, err := NewTreeSet(cmp)
			require.NoError(t, err)

			for _, elem := range tc.initElems {
				s.Add(elem)
			}

			for _, elem := range tc.rmElems {
				s.Del(elem)
			}

			assert.ElementsMatch(t, tc.wantElems, s.Elems())
		})
	}
}

func TestTreeSet_Exist(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		initElems []int
		elem      int
		wantRes   bool
	}{
		{
			name:      "exist element",
			initElems: []int{1, 2, 3, 4},
			elem:      1,
			wantRes:   true,
		}, {
			name:      "not exist element",
			initElems: []int{1, 2, 3, 4},
			elem:      5,
			wantRes:   false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s, err := NewTreeSet(cmp)
			require.NoError(t, err)

			for _, elem := range tc.initElems {
				s.Add(elem)
			}

			assert.Equal(t, tc.wantRes, s.Exist(tc.elem))
		})
	}
}
