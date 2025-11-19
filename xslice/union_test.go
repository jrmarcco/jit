package xslice

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnionSetFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "empty src",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty dst",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		}, {
			name: "no intersection",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		}, {
			name: "intersection",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := UnionSetFunc(tc.src, tc.dst, func(src, dst int) bool { return src == dst })
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func TestUnionSet(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "empty src",
			src:  []int{},
			dst:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty dst",
			src:  []int{1, 2, 3},
			dst:  []int{},
			want: []int{1, 2, 3},
		}, {
			name: "no intersection",
			src:  []int{1, 2, 3},
			dst:  []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		}, {
			name: "intersection",
			src:  []int{1, 2, 3, 4, 5},
			dst:  []int{4, 5, 6, 7, 8},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := UnionSet(tc.src, tc.dst)
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func ExampleUnionSetFunc() {
	src := []int{1, 2, 3}
	dst := []int{4, 5, 6}
	res := UnionSetFunc(src, dst, func(src, dst int) bool { return src == dst })
	fmt.Println(res)
	// Output: [1 2 3 4 5 6]
}

func ExampleUnionSet() {
	src := []int{1, 2, 3}
	dst := []int{4, 5, 6}
	res := UnionSet(src, dst)
	sort.Ints(res)
	fmt.Println(res)
	// Output: [1 2 3 4 5 6]
}
