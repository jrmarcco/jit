package xslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 本文件内复用的用例名常量 ( goconst )。
const (
	tcContains    = "contains"
	tcNotContains = "not contains"
	tcNilElems    = "nil elems"
)

func TestContainsAnyFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elems   []int
		wantRes bool
	}{
		{
			name:    tcContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 6},
			wantRes: true,
		}, {
			name:    tcNotContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    tcEmptySlice,
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilSlice,
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilElems,
			slice:   []int{1, 2, 3},
			elems:   nil,
			wantRes: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ContainsAnyFunc(tc.slice, tc.elems, func(src, dst int) bool { return src == dst })
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestContainsAny(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elems   []int
		wantRes bool
	}{
		{
			name:    tcContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 6},
			wantRes: true,
		}, {
			name:    tcNotContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    tcEmptySlice,
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilSlice,
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilElems,
			slice:   []int{1, 2, 3},
			elems:   nil,
			wantRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ContainsAny(tc.slice, tc.elems)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestContainsAllFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elems   []int
		wantRes bool
	}{
		{
			name:    tcContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 5},
			wantRes: true,
		}, {
			name:    tcNotContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    tcEmptySlice,
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilSlice,
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:  tcNilElems,
			slice: []int{1, 2, 3},
			elems: nil,
			// 空元素集真空包含。
			wantRes: true,
		}, {
			name:    "empty elems",
			slice:   []int{},
			elems:   []int{},
			wantRes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ContainsAllFunc(tc.slice, tc.elems, func(src, dst int) bool { return src == dst })
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestContainsAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elems   []int
		wantRes bool
	}{
		{
			name:    tcContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 5},
			wantRes: true,
		}, {
			name:    tcNotContains,
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    tcEmptySlice,
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    tcNilSlice,
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:  tcNilElems,
			slice: []int{1, 2, 3},
			elems: nil,
			// 空元素集真空包含。
			wantRes: true,
		}, {
			name:    "empty elems",
			slice:   []int{},
			elems:   []int{},
			wantRes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ContainsAll(tc.slice, tc.elems)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
