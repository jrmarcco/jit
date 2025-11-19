package xmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltInMap_Put(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		key  int
		val  int
		data map[int]int
	}{
		{
			name: "basic",
			key:  4,
			val:  4,
			data: map[int]int{1: 1, 2: 2, 3: 3},
		}, {
			name: "exist",
			key:  1,
			val:  2,
			data: map[int]int{1: 1, 2: 2, 3: 3},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			err := m.Put(tc.key, tc.val)
			assert.Nil(t, err)

			if err == nil {
				val, ok := m.Get(tc.key)
				assert.True(t, ok)
				assert.Equal(t, tc.val, val)
			}
		})
	}
}

func TestBuiltInMap_Del(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		key     int
		data    map[int]int
		wantVal int
		wantRes bool
	}{
		{
			name:    "basic",
			key:     1,
			data:    map[int]int{1: 1, 2: 2, 3: 3},
			wantVal: 1,
			wantRes: true,
		}, {
			name:    "not exist",
			key:     4,
			data:    map[int]int{1: 1, 2: 2, 3: 3},
			wantVal: 0,
			wantRes: false,
		}, {
			name:    "empty",
			key:     1,
			data:    map[int]int{},
			wantVal: 0,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			gotVal, gotRes := m.Del(tc.key)
			assert.Equal(t, tc.wantRes, gotRes)
			if gotRes {
				assert.Equal(t, tc.wantVal, gotVal)
				return
			}

			_, gotRes = m.Get(tc.key)
			assert.False(t, gotRes)
		})
	}
}

func TestBuiltInMap_Get(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		key     int
		data    map[int]int
		wantVal int
		wantRes bool
	}{
		{
			name:    "basic",
			key:     1,
			data:    map[int]int{1: 1, 2: 2, 3: 3},
			wantVal: 1,
			wantRes: true,
		}, {
			name:    "not exist",
			key:     4,
			data:    map[int]int{1: 1, 2: 2, 3: 3},
			wantVal: 0,
			wantRes: false,
		}, {
			name:    "empty",
			key:     1,
			data:    map[int]int{},
			wantVal: 0,
			wantRes: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			gotVal, gotRes := m.Get(tc.key)
			assert.Equal(t, tc.wantRes, gotRes)

			if gotRes {
				assert.Equal(t, tc.wantVal, gotVal)
				return
			}
		})
	}
}

func TestBuiltInMap_Keys(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		data    map[string]int
		wantRes []string
	}{
		{
			name: "basic",
			data: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantRes: []string{"a", "b", "c"},
		}, {
			name:    "empty",
			data:    map[string]int{},
			wantRes: []string{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			assert.ElementsMatch(t, tc.wantRes, m.Keys())
		})
	}
}

func TestBuiltInMap_Vals(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		data    map[string]int
		wantRes []int
	}{
		{
			name: "basic",
			data: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantRes: []int{1, 2, 3},
		}, {
			name:    "empty",
			data:    map[string]int{},
			wantRes: []int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			assert.ElementsMatch(t, tc.wantRes, m.Vals())
		})
	}
}

func TestBuiltInMap_Size(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		data    map[string]int
		wantRes int64
	}{
		{
			name: "basic",
			data: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantRes: 3,
		}, {
			name:    "empty",
			data:    map[string]int{},
			wantRes: 0,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newBuiltInMap(tc.data)
			assert.Equal(t, tc.wantRes, m.Size())
		})
	}
}
