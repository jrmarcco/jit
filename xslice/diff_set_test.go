package xslice

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffSet(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		exp  []int
	}{
		{
			name: "empty src",
			src:  []int{},
			dst:  []int{1, 2, 3},
			exp:  []int{},
		}, {
			name: "empty dst",
			src:  []int{1, 2, 3},
			dst:  []int{},
			exp:  []int{1, 2, 3},
		}, {
			name: "no diff",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			exp:  []int{1, 2, 3},
		}, {
			name: "diff",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			exp:  []int{1},
		}, {
			name: "diff with duplicates",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{2, 3, 4, 6, 7},
			exp:  []int{1, 5},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSet(tc.src, tc.dst)
			assert.ElementsMatch(t, got, tc.exp)
		})
	}
}

func TestDiffSetFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		exp  []int
	}{
		{
			name: "empty src",
			src:  []int{},
			dst:  []int{1, 2, 3},
			exp:  []int{},
		}, {
			name: "empty dst",
			src:  []int{1, 2, 3},
			dst:  []int{},
			exp:  []int{1, 2, 3},
		}, {
			name: "no diff",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			exp:  []int{1, 2, 3},
		}, {
			name: "diff",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			exp:  []int{1},
		}, {
			name: "diff with duplicates",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{2, 3, 4, 6, 7},
			exp:  []int{1, 5},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSetFunc(tc.src, tc.dst, func(src, dst int) bool { return src == dst })
			assert.ElementsMatch(t, got, tc.exp)
		})
	}
}

func ExampleDiffSet() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{2, 3, 4, 6, 7}
	diff := DiffSet(src, dst)
	sort.Ints(diff)
	fmt.Println(diff)
	// Output: [1 5]
}

func ExampleDiffSetFunc() {
	src := []int{1, 2, 3, 4, 5}
	dst := []int{2, 3, 4, 6, 7}
	eq := func(a, b int) bool { return a == b }
	diff := DiffSetFunc(src, dst, eq)
	fmt.Println(diff)
	// Output: [1 5]
}
