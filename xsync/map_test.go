package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type param struct {
	name string
	val  int64
}

// 表驱动用例间复用的 map key ( goconst )。
const (
	keyFirst  = "first"
	keySecond = "second"
)

func TestMap_Load(t *testing.T) {
	t.Parallel()

	p := &param{name: keyFirst, val: 1}
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
			key:  keyFirst,
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store(keyFirst, p)
				return &m
			},
			wantVal: p,
			wantRes: true,
		}, {
			name: "not exist",
			key:  keySecond,
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store(keyFirst, p)
				return &m
			},
			wantVal: nil,
			wantRes: false,
		}, {
			name: "nil item",
			key:  keyFirst,
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store(keyFirst, nil)
				return &m
			},
			wantVal: nil,
			wantRes: true,
		}, {
			name: "empty item",
			key:  keyFirst,
			m: func() *Map[string, *param] {
				m := Map[string, *param]{}
				m.Store(keyFirst, empty)
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
		p := &param{name: keyFirst, val: 1}

		val, ok := m.LoadOrStore(keyFirst, p)
		assert.False(t, ok)
		assert.Same(t, p, val)
	})

	t.Run("load not nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p := &param{name: keyFirst, val: 1}

		val, ok := m.LoadOrStore(keyFirst, p)
		assert.False(t, ok)
		assert.Same(t, p, val)

		val, ok = m.LoadOrStore(keyFirst, &param{name: keySecond, val: 2})
		assert.True(t, ok)
		assert.Same(t, p, val)
	})

	t.Run("store nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}

		val, ok := m.LoadOrStore(keyFirst, nil)
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("load nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}

		val, ok := m.LoadOrStore(keyFirst, nil)
		assert.False(t, ok)
		assert.Nil(t, val)

		p := &param{name: keyFirst, val: 1}
		val, ok = m.LoadOrStore(keyFirst, p)
		assert.True(t, ok)
		assert.Nil(t, val)
	})
}

func TestMap_LoadAndDelete(t *testing.T) {
	t.Parallel()

	t.Run("not nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p := &param{name: keyFirst, val: 1}

		m.Store(keyFirst, p)

		val, ok := m.LoadAndDelete(keyFirst)
		assert.True(t, ok)
		assert.Same(t, p, val)

		val, ok = m.Load(keyFirst)
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		m.Store(keyFirst, nil)

		val, ok := m.LoadAndDelete(keyFirst)
		assert.True(t, ok)
		assert.Nil(t, val)

		val, ok = m.Load(keyFirst)
		assert.False(t, ok)
		assert.Nil(t, val)
	})
}

func TestMap_Delete(t *testing.T) {
	t.Parallel()

	m := Map[string, *param]{}
	p1 := &param{name: keyFirst, val: 1}
	p2 := &param{name: keySecond, val: 2}

	m.Store(keyFirst, p1)
	m.Store(keySecond, p2)

	val, ok := m.Load(keyFirst)
	assert.True(t, ok)
	assert.Same(t, p1, val)

	m.Delete(keyFirst)
	val, ok = m.Load(keyFirst)
	assert.False(t, ok)
	assert.Nil(t, val)

	val, ok = m.Load(keySecond)
	assert.True(t, ok)
	assert.Same(t, p2, val)
}

func TestMap_Range(t *testing.T) {
	t.Parallel()

	t.Run("pointer item", func(t *testing.T) {
		t.Parallel()

		m := Map[string, *param]{}
		p1 := &param{name: keyFirst, val: 1}
		p2 := &param{name: keySecond, val: 2}
		empty := &param{}

		m.Store(keyFirst, p1)
		m.Store(keySecond, p2)
		m.Store("third", nil)
		m.Store("fourth", empty)

		wantRes := map[string]*param{
			keyFirst:  p1,
			keySecond: p2,
			"third":   nil,
			"fourth":  empty,
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
		p1 := param{name: keyFirst, val: 1}
		p2 := param{name: keySecond, val: 2}
		empty := param{}

		m.Store(keyFirst, p1)
		m.Store(keySecond, p2)
		m.Store("third", param{})
		m.Store("fourth", empty)

		wantRes := map[string]param{
			keyFirst:  p1,
			keySecond: p2,
			"third":   {},
			"fourth":  empty,
		}

		res := make(map[string]param, 4)

		m.Range(func(key string, val param) bool {
			res[key] = val
			return true
		})

		assert.Equal(t, wantRes, res)
	})
}
