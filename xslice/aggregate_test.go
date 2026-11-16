package xslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jrmarcco/jit/internal/errs"
)

func TestMax(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		wantRes int
		wantErr error
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			wantRes: 0,
			wantErr: errs.NewErrEmptySlice(),
		}, {
			name:    tcNilSlice,
			slice:   nil,
			wantRes: 0,
			wantErr: errs.NewErrEmptySlice(),
		}, {
			name:    tcSingleElement,
			slice:   []int{1},
			wantRes: 1,
			wantErr: nil,
		}, {
			name:    tcMultipleElements,
			slice:   []int{1, 2, 3, 4, 5},
			wantRes: 5,
			wantErr: nil,
		}, {
			name:    tcNegativeElements,
			slice:   []int{-1, -2, -3, -4, -5},
			wantRes: -1,
			wantErr: nil,
		}, {
			name:    tcMixedElements,
			slice:   []int{1, -2, 3, -4, 5},
			wantRes: 5,
			wantErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := Max(tc.slice)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestMin(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		wantRes int
		wantErr error
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			wantRes: 0,
			wantErr: errs.NewErrEmptySlice(),
		}, {
			name:    tcNilSlice,
			slice:   nil,
			wantRes: 0,
			wantErr: errs.NewErrEmptySlice(),
		}, {
			name:    tcSingleElement,
			slice:   []int{1},
			wantRes: 1,
			wantErr: nil,
		}, {
			name:    tcMultipleElements,
			slice:   []int{1, 2, 3, 4, 5},
			wantRes: 1,
			wantErr: nil,
		}, {
			name:    tcNegativeElements,
			slice:   []int{-1, -2, -3, -4, -5},
			wantRes: -5,
			wantErr: nil,
		}, {
			name:    tcMixedElements,
			slice:   []int{1, -2, 3, -4, 5},
			wantRes: -4,
			wantErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := Min(tc.slice)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSum(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		slice   []int
		wantRes int
	}{
		{
			name:    tcEmptySlice,
			slice:   []int{},
			wantRes: 0,
		}, {
			name:    tcNilSlice,
			slice:   nil,
			wantRes: 0,
		}, {
			name:    tcSingleElement,
			slice:   []int{1},
			wantRes: 1,
		}, {
			name:    tcMultipleElements,
			slice:   []int{1, 2, 3, 4, 5},
			wantRes: 15,
		}, {
			name:    tcNegativeElements,
			slice:   []int{-1, -2, -3, -4, -5},
			wantRes: -15,
		}, {
			name:    tcMixedElements,
			slice:   []int{1, -2, 3, -4, 5},
			wantRes: 3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Sum(tc.slice)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func ExampleMax() {
	slice := []int{1, 2, 3, 4, 5}
	m, _ := Max(slice)
	fmt.Println(m)
	// Output: 5
}

func ExampleMin() {
	slice := []int{1, 2, 3, 4, 5}
	m, _ := Min(slice)
	fmt.Println(m)
	// Output: 1
}

func ExampleSum() {
	slice := []int{1, 2, 3, 4, 5}
	sum := Sum(slice)
	fmt.Println(sum)
	// Output: 15
}
