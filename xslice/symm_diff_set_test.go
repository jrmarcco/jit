package xslice

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymmDiffSetFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "basic",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},
			want: []int{1, 2, 3, 6, 7, 8},
		}, {
			name: "empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		}, {
			name: "no common elements",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		}, {
			name: "all elements are common",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: []int{},
		}, {
			name: "src is nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		}, {
			name: "dst is nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := SymmDiffSetFunc(tc.src, tc.dst, func(a, b int) bool { return a == b })
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func TestSymmDiffSet(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "basic",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},
			want: []int{1, 2, 3, 6, 7, 8},
		}, {
			name: "empty",
			src:  []int{},
			dst:  []int{},
			want: []int{},
		}, {
			name: "no common elements",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		}, {
			name: "all elements are common",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: []int{},
		}, {
			name: "src is nil",
			src:  nil,
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		}, {
			name: "dst is nil",
			src:  []int{1, 2, 3},
			dst:  nil,
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := SymmDiffSet(tc.src, tc.dst)
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func ExampleSymmDiffSetFunc() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{4, 5, 6, 7, 8}
	res := SymmDiffSetFunc(src, dst, func(a, b int) bool { return a == b })
	fmt.Println(res)
	// Output: [1 2 3 6 7 8]
}

func ExampleSymmDiffSet() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{4, 5, 6, 7, 8}
	res := SymmDiffSet(src, dst)
	sort.Ints(res)
	fmt.Println(res)
	// Output: [1 2 3 6 7 8]
}
