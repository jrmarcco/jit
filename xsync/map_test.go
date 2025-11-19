package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type param struct {
	name string
	val  int64
}

func TestMap_Load(t *testing.T) {
	t.Parallel()

	p := &param{name: "first", val: 1}
	empty := &param{}

	tcs := []struct {
		name    string
		key     string
		m       func() *Map[string, *param]
		wantVal *param
		wantRes bool
	}{
		{
			name: "basic",
			key:  "first",
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store("first", p)
				return &m
			},
			wantVal: p,
			wantRes: true,
		}, {
			name: "not exist",
			key:  "second",
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store("first", p)
				return &m
			},
			wantVal: nil,
			wantRes: false,
		}, {
			name: "nil item",
			key:  "first",
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store("first", nil)
				return &m
			},
			wantVal: nil,
			wantRes: true,
		}, {
			name: "empty item",
			key:  "first",
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store("first", empty)
				return &m
			},
			wantVal: empty,
			wantRes: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := tc.m()
			val, ok := m.Load(tc.key)
			assert.Equal(t, tc.wantRes, ok)
			assert.Same(t, tc.wantVal, val)
		})
	}
}

func TestMap_LoadOrStore(t *testing.T) {
	t.Parallel()

	t.Run("store not nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p := &param{name: "first", val: 1}

		val, ok := m.LoadOrStore("first", p)
		assert.False(t, ok)
		assert.Same(t, p, val)
	})

	t.Run("load not nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p := &param{name: "first", val: 1}

		val, ok := m.LoadOrStore("first", p)
		assert.False(t, ok)
		assert.Same(t, p, val)

		val, ok = m.LoadOrStore("first", &param{name: "second", val: 2})
		assert.True(t, ok)
		assert.Same(t, p, val)
	})

	t.Run("store nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}

		val, ok := m.LoadOrStore("first", nil)
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("load nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}

		val, ok := m.LoadOrStore("first", nil)
		assert.False(t, ok)
		assert.Nil(t, val)

		p := &param{name: "first", val: 1}
		val, ok = m.LoadOrStore("first", p)
		assert.True(t, ok)
		assert.Nil(t, val)
	})
}

func TestMap_LoadAndDelete(t *testing.T) {
	t.Parallel()

	t.Run("not nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p := &param{name: "first", val: 1}

		m.Store("first", p)

		val, ok := m.LoadAndDelete("first")
		assert.True(t, ok)
		assert.Same(t, p, val)

		val, ok = m.Load("first")
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		m.Store("first", nil)

		val, ok := m.LoadAndDelete("first")
		assert.True(t, ok)
		assert.Nil(t, val)

		val, ok = m.Load("first")
		assert.False(t, ok)
		assert.Nil(t, val)
	})
}

func TestMap_Delete(t *testing.T) {
	t.Parallel()

	m := Map[string, *param]{}
	p1 := &param{name: "first", val: 1}
	p2 := &param{name: "second", val: 2}

	m.Store("first", p1)
	m.Store("second", p2)

	val, ok := m.Load("first")
	assert.True(t, ok)
	assert.Same(t, p1, val)

	m.Delete("first")
	val, ok = m.Load("first")
	assert.False(t, ok)
	assert.Nil(t, val)

	val, ok = m.Load("second")
	assert.True(t, ok)
	assert.Same(t, p2, val)
}

func TestMap_Range(t *testing.T) {
	t.Parallel()

	t.Run("pointer item", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p1 := &param{name: "first", val: 1}
		p2 := &param{name: "second", val: 2}
		empty := &param{}

		m.Store("first", p1)
		m.Store("second", p2)
		m.Store("third", nil)
		m.Store("fourth", empty)

		wantRes := map[string]*param{
			"first":  p1,
			"second": p2,
			"third":  nil,
			"fourth": empty,
		}

		res := make(map[string]*param, 4)

		m.Range(func(key string, val *param) bool {
			res[key] = val
			return true
		})

		assert.Equal(t, wantRes, res)
	})

	t.Run("struct item", func(t *testing.T) {
		t.Parallel()

		m := Map[string, param]{}
		p1 := param{name: "first", val: 1}
		p2 := param{name: "second", val: 2}
		empty := param{}

		m.Store("first", p1)
		m.Store("second", p2)
		m.Store("third", param{})
		m.Store("fourth", empty)

		wantRes := map[string]param{
			"first":  p1,
			"second": p2,
			"third":  {},
			"fourth": empty,
		}

		res := make(map[string]param, 4)

		m.Range(func(key string, val param) bool {
			res[key] = val
			return true
		})

		assert.Equal(t, wantRes, res)
	})
}
