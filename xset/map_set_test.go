package xset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// overcapacity
func TestMapSet_Add(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		size     int
		keys     []int
		wantSize int
		wantRes  map[int]struct{}
	}{
		{
			name:    "add vals",
			size:    10,
			keys:    []int{1, 2, 3, 4},
			wantRes: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewMapSet[int](tc.size)
			for _, key := range tc.keys {
				s.Add(key)
			}

			assert.Equal(t, tc.wantRes, s.m)
		})
	}
}

func TestMapSet_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		delVal  int
		m       map[int]struct{}
		wantSet map[int]struct{}
	}{
		{
			name:    "del val",
			delVal:  1,
			m:       map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			wantSet: map[int]struct{}{2: {}, 3: {}, 4: {}},
		}, {
			name:    "del val not exist",
			delVal:  5,
			m:       map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			wantSet: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewMapSet[int](len(tc.m))
			s.m = tc.m
			s.Del(tc.delVal)
			assert.Equal(t, tc.wantSet, s.m)
		})
	}
}

func TestMapSet_Exist(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		val     int
		m       map[int]struct{}
		wantRes bool
	}{
		{
			name:    "exist val",
			val:     1,
			m:       map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			wantRes: true,
		}, {
			name:    "exist val not exist",
			val:     5,
			m:       map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewMapSet[int](len(tc.m))
			s.m = tc.m
			assert.Equal(t, tc.wantRes, s.Exist(tc.val))
		})
	}
}

func TestMapSet_Keys(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		m       map[int]struct{}
		wantRes []int
	}{
		{
			name:    "keys",
			m:       map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			wantRes: []int{1, 2, 3, 4},
		}, {
			name:    "disorder keys",
			m:       map[int]struct{}{4: {}, 3: {}, 2: {}, 1: {}},
			wantRes: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewMapSet[int](len(tc.m))
			s.m = tc.m
			assert.ElementsMatch(t, tc.wantRes, s.Elems())
		})
	}
}
