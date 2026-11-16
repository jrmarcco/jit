// Package xsync 提供并发原语的泛型封装:
// Pool、Map 以及支持 context 的 Cond。
package xsync

import "sync"

// Pool 为 sync.Pool 的泛型封装。
type Pool[T any] struct {
	p sync.Pool
}

// NewPool 以工厂函数创建 Pool。
// factory 不能返回 nil。
func NewPool[T any](factory func() T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{
			New: func() any {
				return factory()
			},
		},
	}
}

// Get 从池中取出一个对象 ( 池为空时调用 factory 创建 )。
func (p *Pool[T]) Get() T {
	v, _ := p.p.Get().(T)
	return v
}

// Put 把对象放回池中。
func (p *Pool[T]) Put(t T) {
	p.p.Put(t)
}
