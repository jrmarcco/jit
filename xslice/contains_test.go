package xslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elem    int
		wantRes bool
	}{
		{
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elem:    3,
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elem:    6,
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elem:    1,
			wantRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := ContainsFunc(tc.slice, func(t int) bool { return t == tc.elem })
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elem    int
		wantRes bool
	}{
		{
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elem:    3,
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elem:    6,
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elem:    1,
			wantRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Contains(tc.slice, tc.elem)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestContainsAnyFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slice   []int
		elems   []int
		wantRes bool
	}{
		{
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 6},
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil elems",
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
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 6},
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil elems",
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
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 5},
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil elems",
			slice:   []int{1, 2, 3},
			elems:   nil,
			wantRes: false,
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
			name:    "contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{3, 5},
			wantRes: true,
		}, {
			name:    "not contains",
			slice:   []int{1, 2, 3, 4, 5},
			elems:   []int{6, 7},
			wantRes: false,
		}, {
			name:    "empty slice",
			slice:   []int{},
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil slice",
			slice:   nil,
			elems:   []int{1, 2, 3},
			wantRes: false,
		}, {
			name:    "nil elems",
			slice:   []int{1, 2, 3},
			elems:   nil,
			wantRes: false,
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
