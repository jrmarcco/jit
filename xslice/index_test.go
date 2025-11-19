package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    "found",
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := IndexFunc(tc.slice, func(t int) bool { return t == tc.elem })
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    "found",
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Index(tc.slice, tc.elem)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestLastIndexFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    "found",
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
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

func TestLastIndex(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		elem    int
		wantRes int
	}{
		{
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: -1,
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: -1,
		}, {
			name:    "found",
			slice:   []int{1, 2, 3},
			elem:    2,
			wantRes: 1,
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
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: []int{},
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: []int{},
		}, {
			name:    "found",
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
			name:    "empty slice",
			slice:   []int{},
			elem:    1,
			wantRes: []int{},
		}, {
			name:    "not found",
			slice:   []int{1, 2, 3},
			elem:    4,
			wantRes: []int{},
		}, {
			name:    "found",
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

func ExampleIndexFunc() {
	slice := []int{1, 2, 3, 2, 4}
	res := IndexFunc(slice, func(t int) bool { return t == 2 })
	fmt.Println(res)
	// Output: 1
}

func ExampleIndex() {
	slice := []int{1, 2, 3, 2, 4}
	res := Index(slice, 2)
	fmt.Println(res)
	// Output: 1
}

func ExampleLastIndexFunc() {
	slice := []int{1, 2, 3, 2, 4}
	res := LastIndexFunc(slice, func(t int) bool { return t == 2 })
	fmt.Println(res)
	// Output: 3
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
