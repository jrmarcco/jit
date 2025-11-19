package xsync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	t.Parallel()

	cnt := 0
	p := NewPool(func() []byte {
		cnt++
		res := make([]byte, 1, 16)
		res[0] = 'A'
		return res
	})

	res := p.Get()
	assert.Equal(t, []byte{'A'}, res)

	res = append(res, 'B')
	p.Put(res)

	res = p.Get()
	if cnt == 1 {
		assert.Equal(t, []byte{'A', 'B'}, res)
	} else {
		assert.Equal(t, []byte{'A'}, res)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/JrMarcco/easy_kit/plus_sync
// cpu: 13th Gen Intel(R) Core(TM) i7-13700KF
// === RUN   BenchmarkPool_Get
// BenchmarkPool_Get
// === RUN   BenchmarkPool_Get/pool
// BenchmarkPool_Get/pool
// BenchmarkPool_Get/pool-24               28404036                42.74 ns/op            0 B/op          0 allocs/op
// === RUN   BenchmarkPool_Get/sync.Pool
// BenchmarkPool_Get/sync.Pool
// BenchmarkPool_Get/sync.Pool-24          27822126                42.71 ns/op            0 B/op          0 allocs/op
func BenchmarkPool_Get(b *testing.B) {
	p := NewPool(func() string {
		return ""
	})

	sp := &sync.Pool{
		New: func() any {
			return ""
		},
	}

	b.Run("pool", func(b *testing.B) {
		for b.Loop() {
			p.Get()
		}
	})

	b.Run("sync.Pool", func(b *testing.B) {
		for b.Loop() {
			sp.Get()
		}
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/JrMarcco/easy_kit/plus_sync
// cpu: 13th Gen Intel(R) Core(TM) i7-13700KF
// === RUN   BenchmarkPool_Get_struct
// BenchmarkPool_Get_struct
// === RUN   BenchmarkPool_Get_struct/pool
// BenchmarkPool_Get_struct/pool
// BenchmarkPool_Get_struct/pool-24                21106546                56.58 ns/op           48 B/op          1 allocs/op
// === RUN   BenchmarkPool_Get_struct/sync.Pool
// BenchmarkPool_Get_struct/sync.Pool
// BenchmarkPool_Get_struct/sync.Pool-24           22635163                52.60 ns/op           48 B/op          1 allocs/op
func BenchmarkPool_Get_struct(b *testing.B) {
	type test struct {
		ID   int
		Name string
		Addr string
	}

	p := NewPool(func() test {
		return test{ID: 1, Name: "test", Addr: "127.0.0.1:8080"}
	})

	sp := &sync.Pool{
		New: func() any {
			return test{ID: 1, Name: "test", Addr: "127.0.0.1:8080"}
		},
	}

	b.Run("pool", func(b *testing.B) {
		for b.Loop() {
			p.Get()
		}
	})

	b.Run("sync.Pool", func(b *testing.B) {
		for b.Loop() {
			sp.Get()
		}
	})
}
