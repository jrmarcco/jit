package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mf matchFunc[int] = func(v int) bool {
	return v%2 == 0
}

func TestFind(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		match   matchFunc[int]
		wantRes int
		wantOk  bool
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			match:   mf,
			wantRes: 0,
			wantOk:  false,
		}, {
			name:    "no match",
			slice:   []int{1, 3, 5},
			match:   mf,
			wantRes: 0,
			wantOk:  false,
		}, {
			name:    "single match",
			slice:   []int{2, 4, 6},
			match:   mf,
			wantRes: 2,
			wantOk:  true,
		}, {
			name:    "multiple matches",
			slice:   []int{1, 2, 3, 4, 5},
			match:   mf,
			wantRes: 2,
			wantOk:  true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, ok := Find(tc.slice, tc.match)
			assert.Equal(t, tc.wantOk, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestFindAll(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		match   matchFunc[int]
		wantRes []int
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			match:   mf,
			wantRes: []int{},
		}, {
			name:    "no match",
			slice:   []int{1, 3, 5},
			match:   mf,
			wantRes: []int{},
		}, {
			name:    "single match",
			slice:   []int{2, 4, 6},
			match:   mf,
			wantRes: []int{2, 4, 6},
		}, {
			name:    "multiple matches",
			slice:   []int{1, 2, 3, 4, 5},
			match:   mf,
			wantRes: []int{2, 4},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := FindAll(tc.slice, tc.match)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func ExampleFind() {
	slice := []int{1, 2, 3, 4, 5}
	res, ok := Find(slice, mf)
	fmt.Println(res, ok)
	// Output: 2 true
}

func ExampleFindAll() {
	slice := []int{1, 2, 3, 4, 5}
	res := FindAll(slice, mf)
	fmt.Println(res)
	// Output: [2 4]
}
