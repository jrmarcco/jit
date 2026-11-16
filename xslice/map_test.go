package xslice

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 本文件内复用的用例名常量 ( goconst )。
const tcNil = "nil"

func TestMap(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		fn   func(idx int, src int) int
		want []int
	}{
		{
			name: tcBasic,
			src:  []int{1, 2, 3},
			fn:   func(_ int, src int) int { return src * 2 },
			want: []int{2, 4, 6},
		}, {
			name: tcEmpty,
			src:  []int{},
			fn:   func(_ int, src int) int { return src * 2 },
			want: []int{},
		}, {
			name: tcNil,
			src:  nil,
			fn:   func(_ int, src int) int { return src * 2 },
			want: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Map(tc.src, tc.fn)
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func TestFilterMap(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		fn   func(idx int, src int) (int, bool)
		want []int
	}{
		{
			name: tcBasic,
			src:  []int{1, 2, 3, 4, 5},
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: []int{2, 6, 10},
		}, {
			name: tcEmpty,
			src:  []int{},
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: []int{},
		}, {
			name: tcNil,
			src:  nil,
			fn:   func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 },
			want: nil,
		}, {
			name: tcNoMatch,
			src:  []int{1, 2, 3, 4, 5},
			fn:   func(_ int, src int) (int, bool) { return src * 2, src > 10 },
			want: []int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := FilterMap(tc.src, tc.fn)
			assert.ElementsMatch(t, tc.want, res)
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		fn   func(elem int) string
		want map[string]int
	}{
		{
			name: tcBasic,
			src:  []int{1, 2, 3, 4, 5},
			fn:   strconv.Itoa,
			want: map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5},
		}, {
			name: tcEmpty,
			src:  []int{},
			fn:   strconv.Itoa,
			want: map[string]int{},
		}, {
			name: tcNil,
			src:  nil,
			fn:   strconv.Itoa,
			want: map[string]int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ToMap(tc.src, tc.fn)
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestToMapWithVal(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		src  []int
		fn   func(elem int) (string, int)
		want map[string]int
	}{
		{
			name: tcBasic,
			src:  []int{1, 2, 3, 4, 5},
			fn:   func(elem int) (string, int) { return strconv.Itoa(elem), elem },
			want: map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5},
		}, {
			name: tcEmpty,
			src:  []int{},
			fn:   func(elem int) (string, int) { return strconv.Itoa(elem), elem },
			want: map[string]int{},
		}, {
			name: tcNil,
			src:  nil,
			fn:   func(elem int) (string, int) { return strconv.Itoa(elem), elem },
			want: map[string]int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ToMapWithVal(tc.src, tc.fn)
			assert.Equal(t, tc.want, res)
		})
	}
}

func ExampleMap() {
	src := []int{1, 2, 3, 4, 5}
	res := Map(src, func(_ int, src int) int { return src * 2 })
	fmt.Println(res)
	// Output: [2 4 6 8 10]
}

func ExampleFilterMap() {
	src := []int{1, 2, 3, 4, 5}
	res := FilterMap(src, func(idx int, src int) (int, bool) { return src * 2, idx%2 == 0 })
	fmt.Println(res)
	// Output: [2 6 10]
}
