package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 本文件内复用的用例名常量 ( goconst )。
const (
	tcFound    = "found"
	tcNotFound = "not found"
)

func TestLastIndexFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    tcNotFound,
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    tcFound,
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
		}, {
			name:    "found last occurrence",
			slice:   []int{1, 2, 3, 2, 4},
			elem:    2,
			wantRes: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := LastIndexFunc(tc.slice, func(t int) bool { return t == tc.elem })
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func ExampleLastIndexFunc() {
	slice := []int{1, 2, 3, 2, 4}
	res := LastIndexFunc(slice, func(t int) bool { return t == 2 })
	fmt.Println(res)
	// Output: 3
}

func TestLastIndex(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    tcNotFound,
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    tcFound,
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
		}, {
			name:    "found last occurrence",
			slice:   []int{1, 2, 3, 2, 4},
			elem:    2,
			wantRes: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := LastIndex(tc.slice, tc.elem)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestIndexAllFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes []int
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			elem:    1,
			wantRes: []int{},
		}, {
			name:    tcNotFound,
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: []int{},
		}, {
			name:    tcFound,
			slice:   []int{1, 2, 3, 2, 4},
			elem:    2,
			wantRes: []int{1, 3},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := IndexAllFunc(tc.slice, func(t int) bool { return t == tc.elem })
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}

func TestIndexAll(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes []int
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			elem:    1,
			wantRes: []int{},
		}, {
			name:    tcNotFound,
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: []int{},
		}, {
			name:    tcFound,
			slice:   []int{1, 2, 3, 2, 4},
			elem:    2,
			wantRes: []int{1, 3},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := IndexAll(tc.slice, tc.elem)
			assert.ElementsMatch(t, tc.wantRes, res)
		})
	}
}

func ExampleLastIndex() {
	slice := []int{1, 2, 3, 2, 4}
	res := LastIndex(slice, 2)
	fmt.Println(res)
	// Output: 3
}

func ExampleIndexAllFunc() {
	slice := []int{1, 2, 3, 2, 4}
	res := IndexAllFunc(slice, func(t int) bool { return t == 2 })
	fmt.Println(res)
	// Output: [1 3]
}

func ExampleIndexAll() {
	slice := []int{1, 2, 3, 2, 4}
	res := IndexAll(slice, 2)
	fmt.Println(res)
	// Output: [1 3]
}
