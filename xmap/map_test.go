package xmap

import (
	"testing"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		m    map[int]int
		want []int
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []int{},
		}, {
			name: "nil",
			m:    nil,
			want: []int{},
		},
	}

	for _, tc := range tcs {
		res := Keys(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestVals(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		m    map[int]int
		want []int
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []int{1, 2, 3},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []int{},
		}, {
			name: "nil",
			m:    nil,
			want: []int{},
		},
	}

	for _, tc := range tcs {
		res := Vals(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestKeysVals(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		m    map[int]int
		want []MapKV[int, int]
	}{
		{
			name: "basic",
			m:    map[int]int{1: 1, 2: 2, 3: 3},
			want: []MapKV[int, int]{{Key: 1, Val: 1}, {Key: 2, Val: 2}, {Key: 3, Val: 3}},
		}, {
			name: "empty",
			m:    map[int]int{},
			want: []MapKV[int, int]{},
		}, {
			name: "nil",
			m:    nil,
			want: []MapKV[int, int]{},
		},
	}

	for _, tc := range tcs {
		res := KeysVals(tc.m)
		assert.ElementsMatch(t, tc.want, res)
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		keys    []int
		vals    []int
		wantRes map[int]int
		wantErr error
	}{
		{
			name:    "basic",
			keys:    []int{1, 2, 3},
			vals:    []int{1, 2, 3},
			wantRes: map[int]int{1: 1, 2: 2, 3: 3},
			wantErr: nil,
		}, {
			name:    "nil keys",
			keys:    nil,
			vals:    []int{1, 2, 3},
			wantRes: nil,
			wantErr: errs.NilErr("keys or vals"),
		}, {
			name:    "nil vals",
			keys:    []int{1, 2, 3},
			vals:    nil,
			wantRes: nil,
			wantErr: errs.NilErr("keys or vals"),
		}, {
			name:    "different lengths",
			keys:    []int{1, 2, 3},
			vals:    []int{1, 2},
			wantRes: nil,
			wantErr: errs.ErrInvalidKeyValLen(),
		}, {
			name:    "empty keys",
			keys:    []int{},
			vals:    []int{1, 2, 3},
			wantRes: map[int]int{},
			wantErr: nil,
		}, {
			name:    "empty vals",
			keys:    []int{1, 2, 3},
			vals:    []int{},
			wantRes: map[int]int{},
			wantErr: nil,
		},
	}

	for _, tc := range tcs {
		res, err := ToMap(tc.keys, tc.vals)
		assert.Equal(t, tc.wantErr, err)

		if err != nil {
			return
		}

		assert.Equal(t, tc.wantRes, res)
	}
}

func TestMerge(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		maps    []map[int]int
		wantRes map[int]int
	}{
		{
			name: "basic",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{4: 4, 5: 5, 6: 6},
			},
			wantRes: map[int]int{
				1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6,
			},
		}, {
			name: "with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{1: 10, 2: 20, 3: 30},
			},
			wantRes: map[int]int{
				1: 10, 2: 20, 3: 30,
			},
		}, {
			name: "diff len with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{1: 10, 2: 20, 3: 30, 4: 40},
			},
			wantRes: map[int]int{
				1: 10, 2: 20, 3: 30, 4: 40,
			},
		}, {
			name: "over two maps with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7},
				{1: 10, 2: 20, 3: 30, 4: 40, 5: 50, 6: 60},
				{1: 10, 2: 20, 3: 30},
				{1: 100, 2: 200, 3: 300},
				{1: 1000, 2: 2000, 3: 3000},
			},
			wantRes: map[int]int{
				1: 1000, 2: 2000, 3: 3000, 4: 40, 5: 50, 6: 60, 7: 7,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := Merge(tc.maps...)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestMergeFunc(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		maps      []map[int]int
		mergeFunc func(int, int) int
		wantRes   map[int]int
	}{
		{
			name: "basic",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{4: 4, 5: 5, 6: 6},
			},
			mergeFunc: func(a, b int) int {
				return a + b
			},
			wantRes: map[int]int{
				1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6,
			},
		}, {
			name: "with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{1: 10, 2: 20, 3: 30},
			},
			mergeFunc: func(a, b int) int {
				return a + b
			},
			wantRes: map[int]int{
				1: 11, 2: 22, 3: 33,
			},
		}, {
			name: "diff len with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3},
				{1: 10, 2: 20, 3: 30, 4: 40},
			},
			mergeFunc: func(a, b int) int {
				return a * b
			},
			wantRes: map[int]int{
				1: 10, 2: 40, 3: 90, 4: 40,
			},
		}, {
			name: "over two maps with key conflict",
			maps: []map[int]int{
				{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7},
				{1: 10, 2: 20, 3: 30, 4: 40, 5: 50, 6: 60},
				{1: 10, 2: 20, 3: 30},
				{1: 100, 2: 200, 3: 300},
				{1: 1000, 2: 2000, 3: 3000},
			},
			mergeFunc: func(a, b int) int {
				return a + b
			},
			wantRes: map[int]int{
				1: 1121, 2: 2242, 3: 3363, 4: 44, 5: 55, 6: 66, 7: 7,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := MergeFunc(tc.mergeFunc, tc.maps...)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
