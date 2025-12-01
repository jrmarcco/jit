package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jrmarcco/jit/bean/option"
	"github.com/stretchr/testify/assert"
)

func poolInternalState(p *BlockTaskPool) int32 {
	for {
		state := atomic.LoadInt32(&p.state)
		if state != stateLocked {
			return state
		}
	}
}

func runningPool(t *testing.T, initG, queueSize int32, opts ...option.Opt[BlockTaskPool]) *BlockTaskPool {
	t.Helper()

	p, err := NewBlockTaskPool(initG, queueSize, opts...)
	assert.NoError(t, err)

	assert.Equal(t, poolInternalState(p), stateCreated)
	assert.NoError(t, p.Start())
	assert.Equal(t, poolInternalState(p), stateRunning)
	return p
}

func runningPoolWithFilledQueue(t *testing.T, initG, queueSize int32) (p *BlockTaskPool, wait chan struct{}) {
	t.Helper()

	p = runningPool(t, initG, queueSize)
	wait = make(chan struct{})

	for i := int32(0); i < initG+queueSize; i++ {
		err := p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
			<-wait
			return nil
		}))
		assert.NoError(t, err)
	}
	return
}

func TestNewBlockTaskPool(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		initG     int32
		queueSize int32
		wantErr   error
	}{
		{
			name:      "basic",
			initG:     1,
			queueSize: 1,
			wantErr:   nil,
		}, {
			name:      "queue size is negative",
			initG:     1,
			queueSize: -1,
			wantErr:   errInvalidParam,
		}, {
			name:      "queue size is 0",
			initG:     1,
			queueSize: 0,
			wantErr:   nil,
		}, {
			name:      "queue size greater than 0",
			initG:     1,
			queueSize: 1,
			wantErr:   nil,
		}, {
			name:      "init goroutines is negative",
			initG:     -1,
			queueSize: 1,
			wantErr:   errInvalidParam,
		}, {
			name:      "init goroutines is 0",
			initG:     0,
			queueSize: 1,
			wantErr:   errInvalidParam,
		}, {
			name:      "init goroutines greater than 0",
			initG:     1,
			queueSize: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p, err := NewBlockTaskPool(tc.initG, tc.queueSize)
			assert.ErrorIs(t, err, tc.wantErr)

			if err == nil {
				assert.NotNil(t, p)
				assert.Equal(t, stateCreated, poolInternalState(p))

				assert.Equal(t, tc.initG, p.initG)
				assert.Equal(t, int(tc.queueSize), cap(p.queue))
			}
		})
	}
}

func TestNewBlockTaskPoolWithOption(t *testing.T) {
	t.Parallel()

	t.Run("with max idle time", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithMaxIdleTime(time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(1), p.coreG)
		assert.Equal(t, int32(1), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, time.Second, p.maxIdleTime)
	})

	t.Run("with submit timeout", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithSubmitTimeout(time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(1), p.coreG)
		assert.Equal(t, int32(1), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, time.Second, p.submitTimeout)
	})

	t.Run("with core goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with max goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithMaxG(2))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with core and max goroutine", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2), WithMaxG(4))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(4), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, int32(4), p.maxG)
	})

	t.Run("with core != init and max == init", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(2), WithMaxG(1))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(2), p.coreG)
		assert.Equal(t, int32(2), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with core == init and max != init", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithCoreG(1), WithMaxG(4))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, int32(4), p.coreG)
		assert.Equal(t, int32(4), p.maxG)
		assert.Equal(t, 3, cap(p.queue))
	})

	t.Run("with queue backlog rate", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(0.5))
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, stateCreated, poolInternalState(p))
		assert.Equal(t, int32(1), p.initG)
		assert.Equal(t, 3, cap(p.queue))
		assert.Equal(t, 0.5, p.queueBacklogRate)
	})

	t.Run("queue backlog rate is negative", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(-1))
		assert.ErrorIs(t, err, errInvalidParam)
		assert.Nil(t, p)
	})

	t.Run("queue backlog rate is greater than 1", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3, WithQueueBacklogRate(1.1))
		assert.ErrorIs(t, err, errInvalidParam)
		assert.Nil(t, p)
	})
}

func TestBlockTaskPool_State(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 3)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		ch, err := p.State(t.Context(), time.Millisecond)
		assert.NoError(t, err)
		assert.NotZero(t, ch)

		done, err := p.Shutdown()
		assert.NoError(t, err)
		<-done
	})

	t.Run("call after context cancel", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err = p.State(ctx, time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("call after shutdown", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		done, err := p.Shutdown()
		assert.NoError(t, err)

		<-done

		_, err = p.State(t.Context(), time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("call after shutdown now", func(t *testing.T) {
		t.Parallel()

		p, err := NewBlockTaskPool(1, 1)
		assert.NoError(t, err)

		err = p.Start()
		assert.NoError(t, err)

		_, err = p.ShutdownNow()
		assert.NoError(t, err)

		_, err = p.State(t.Context(), time.Millisecond)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("close chan after context timeout", func(t *testing.T) {
		t.Parallel()

		initG, queueSize := int32(1), int32(3)
		p, waitCh := runningPoolWithFilledQueue(t, initG, queueSize)

		ctx, cancel := context.WithTimeout(t.Context(), 3*time.Millisecond)
		stateCh, err := p.State(ctx, time.Millisecond)
		assert.NoError(t, err)

		go func() {
			// 模拟 context 超时
			<-time.After(3 * time.Millisecond)
			cancel()
		}()

		for {
			state, ok := <-stateCh
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}

		close(waitCh)
		_, err = p.Shutdown()
		assert.NoError(t, err)
	})

	t.Run("close chan after context canceled", func(t *testing.T) {
		t.Parallel()

		initG, queueSize := int32(1), int32(3)
		p, waitCh := runningPoolWithFilledQueue(t, initG, queueSize)

		ctx, cancel := context.WithCancel(t.Context())
		stateCh, err := p.State(ctx, time.Millisecond)
		assert.NoError(t, err)

		go func() {
			cancel()
		}()

		for {
			state, ok := <-stateCh
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}

		close(waitCh)
		_, err = p.Shutdown()
		assert.NoError(t, err)
	})

	t.Run("close chan after shutdown", func(t *testing.T) {
		t.Parallel()

		p := runningPool(t, 1, 3)

		ch, err := p.State(t.Context(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err = p.Shutdown()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})

	t.Run("close chan after shutdown now", func(t *testing.T) {
		t.Parallel()

		p := runningPool(t, 1, 3)

		ch, err := p.State(t.Context(), time.Millisecond)
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err = p.ShutdownNow()
			assert.NoError(t, err)
		}()

		for {
			state, ok := <-ch
			if !ok {
				break
			}
			assert.NotZero(t, state)
		}
	})
}

func TestBlockTaskPool_Submit(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name       string
		poolFunc   func(t *testing.T) *BlockTaskPool
		submitFunc func(t *testing.T, p *BlockTaskPool)
		wantErr    error
	}{
		{
			name: "basic",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				t.Helper()

				p, err := NewBlockTaskPool(1, 3)
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				t.Helper()

				var err error
				for range 3 {
					err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
						return nil
					}))
					assert.NoError(t, err)
				}
			},
		}, {
			name: "nil task",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				t.Helper()
				p, err := NewBlockTaskPool(1, 1)
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				t.Helper()

				err := p.Submit(t.Context(), nil)
				assert.ErrorIs(t, err, errInvalidTask)
			},
		}, {
			name: "submit timeout",
			poolFunc: func(t *testing.T) *BlockTaskPool {
				t.Helper()

				p, err := NewBlockTaskPool(1, 1, WithSubmitTimeout(time.Millisecond))
				assert.NoError(t, err)
				assert.NotNil(t, p)
				return p
			},
			submitFunc: func(t *testing.T, p *BlockTaskPool) {
				t.Helper()

				done := make(chan struct{})
				err := p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
					<-done
					return nil
				}))
				assert.NoError(t, err)

				ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
				err = p.Submit(ctx, TaskFunc(func(_ context.Context) error {
					<-done
					return nil
				}))
				cancel()
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				close(done)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := tc.poolFunc(t)
			assert.Equal(t, poolInternalState(p), stateCreated)

			tc.submitFunc(t, p)
		})
	}
}

func TestBlockTaskPool_Shutdown(t *testing.T) {
	t.Parallel()

	p1, err := NewBlockTaskPool(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, p1)
	assert.Equal(t, poolInternalState(p1), stateCreated)

	done, err := p1.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, errPoolIsNotRunning)
	assert.Equal(t, poolInternalState(p1), stateCreated)

	err = p1.Start()
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p1), stateRunning)

	done, err = p1.Shutdown()
	assert.NotNil(t, done)
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p1), stateClosing)
	<-done
	assert.Equal(t, poolInternalState(p1), stateClosed)

	done, err = p1.Shutdown()
	assert.Nil(t, done)
	assert.ErrorIs(t, err, errPoolIsClosed)

	p2, err := NewBlockTaskPool(1, 3)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
	assert.Equal(t, poolInternalState(p2), stateCreated)

	tasks, err := p2.ShutdownNow()
	assert.Equal(t, 0, len(tasks))
	assert.ErrorIs(t, err, errPoolIsNotRunning)
	assert.Equal(t, poolInternalState(p2), stateCreated)

	err = p2.Start()
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p2), stateRunning)

	tasks, err = p2.ShutdownNow()
	assert.Equal(t, 0, len(tasks))
	assert.NoError(t, err)
	assert.Equal(t, poolInternalState(p2), stateClosed)

	_, err = p2.ShutdownNow()
	assert.ErrorIs(t, err, errPoolIsClosed)
}

func TestBlockTaskPool_state_machine(t *testing.T) {
	t.Parallel()

	t.Run("without submit after start", func(t *testing.T) {
		t.Parallel()

		t.Run("need goroutine le allow goroutine", func(t *testing.T) {
			t.Parallel()

			initG, coreG := int32(1), int32(3)
			queueSize := coreG

			p, err := NewBlockTaskPool(initG, queueSize, WithCoreG(coreG))
			assert.NoError(t, err)
			assert.NotNil(t, p)

			assert.Equal(t, stateCreated, poolInternalState(p))

			done := make(chan struct{}, queueSize)
			for i := 0; i < int(queueSize); i++ {
				err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
					// 阻塞 task
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}
			assert.Equal(t, int32(0), p.countG())

			assert.NoError(t, p.Start())
			assert.Equal(t, stateRunning, poolInternalState(p))
			assert.Equal(t, p.maxG, p.countG())

			// 释放 task
			close(done)
		})

		t.Run("need goroutine gt allow goroutine", func(t *testing.T) {
			t.Parallel()

			initG, maxG, queueSize := int32(5), int32(10), int32(15)
			p, err := NewBlockTaskPool(initG, queueSize, WithMaxG(maxG))
			assert.NoError(t, err)
			assert.NotNil(t, p)

			assert.Equal(t, stateCreated, poolInternalState(p))

			done := make(chan struct{}, queueSize)
			for i := 0; i < int(queueSize); i++ {
				err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
					// 阻塞 task
					<-done
					return nil
				}))
				assert.NoError(t, err)
			}
			// 这个时候 queue 里有 15 个任务
			assert.Equal(t, int32(0), p.countG())

			// 启动后只会启动 10 个 goroutine 去执行任务（maxG = 10）
			assert.NoError(t, p.Start())
			assert.Equal(t, stateRunning, poolInternalState(p))
			assert.Equal(t, maxG, p.countG())

			// 释放 task
			close(done)
		})
	})

	t.Run("start in concurrency with submit", func(t *testing.T) {
		t.Parallel()

		initG, maxG, queueSize := int32(5), int32(10), int32(15)
		p, err := NewBlockTaskPool(initG, queueSize, WithMaxG(maxG))
		assert.NoError(t, err)
		assert.NotNil(t, p)

		assert.Equal(t, stateCreated, poolInternalState(p))

		errChan := make(chan error)
		go func() {
			// 模拟并发 start
			time.Sleep(10 * time.Millisecond)
			errChan <- p.Start()
		}()

		done := make(chan struct{}, queueSize)
		for i := 0; i < int(queueSize); i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				// 阻塞 task
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.NoError(t, <-errChan)
		assert.Equal(t, stateRunning, poolInternalState(p))
		assert.Equal(t, maxG, p.countG())

		// 释放 task
		close(done)
	})

	t.Run("goroutine keep in init", func(t *testing.T) {
		t.Parallel()

		initG, queueSize := int32(1), int32(3)

		p := runningPool(t, initG, queueSize)
		assert.Equal(t, initG, p.countG())

		var err error
		done := make(chan struct{}, queueSize)
		for i := 0; i < int(queueSize); i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.Equal(t, initG, p.countG())
		for i := 0; i < int(initG); i++ {
			done <- struct{}{}
		}

		for i := 0; i < int(queueSize-initG); i++ {
			assert.Equal(t, initG, p.countG())
			done <- struct{}{}
		}
	})

	t.Run("goroutine from init to core", func(t *testing.T) {
		t.Parallel()

		initG, coreG := int32(8), int32(16)
		queueSize := coreG

		p := runningPool(t, initG, queueSize, WithCoreG(coreG), WithMaxIdleTime(2*time.Millisecond))
		assert.Equal(t, initG, p.countG())

		var err error
		done := make(chan struct{})

		// 这里要注意并发竞争：
		// 提交 task 时，task 还在 p.queue 内还未被永久 goroutine 区间的 goroutine 取走就进行 p.allowToCreateG() 判断，
		// 此时 p.queue 的 task 数会大于等于 1，也即是队列积压率大于 0，
		// 就会创建核心 goroutine 去立即执行当前 task。
		for i := 0; i < int(initG); i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.Equal(t, 0, len(p.queue))
		// 至少有 initG
		// 这里用 LessOrEqual 判断是为了兼容并发竞争导致的核心 goroutine 创建
		assert.LessOrEqual(t, initG, p.countG())

		// 增加任务数 init -> core ( 8 -> 16 )
		for i := initG; i < coreG; i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}
		assert.Equal(t, coreG, p.countG())

		// 释放所有任务
		close(done)

		// 等待空闲后 goroutine 数恢复到 initG
		for p.countG() > initG {
		}
		assert.Equal(t, initG, p.countG())
	})

	t.Run("goroutine from core to max", func(t *testing.T) {
		t.Parallel()

		initG, coreG, maxG := int32(8), int32(16), int32(32)
		queueSize := coreG

		p := runningPool(t, initG, queueSize, WithCoreG(coreG), WithMaxG(maxG), WithMaxIdleTime(2*time.Millisecond))
		assert.Equal(t, initG, p.countG())

		var err error
		done := make(chan struct{})

		// 这里要注意并发竞争：
		// 提交 task 时，task 还在 p.queue 内还未被永久 goroutine 区间的 goroutine 取走就进行 p.allowToCreateG() 判断，
		// 此时 p.queue 的 task 数会大于等于 1，也即是队列积压率大于 0，
		// 就会创建核心 goroutine 去立即执行当前 task。
		for i := 0; i < int(initG); i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}

		assert.Equal(t, 0, len(p.queue))
		// 至少有 initG
		// 这里用 LessOrEqual 判断是为了兼容并发竞争导致的核心 goroutine 创建
		assert.LessOrEqual(t, initG, p.countG())

		// 增加任务数 init -> core ( 8 -> 16 )
		for i := initG; i < coreG; i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}
		assert.Equal(t, coreG, p.countG())

		// 增加任务 core -> max ( 16 -> 32 )
		for i := coreG; i < maxG; i++ {
			err = p.Submit(t.Context(), TaskFunc(func(_ context.Context) error {
				<-done
				return nil
			}))
			assert.NoError(t, err)
		}
		assert.Equal(t, maxG, p.countG())

		// 释放所有任务
		close(done)

		// 等待空闲后 goroutine 数恢复到 initG
		for p.countG() > initG {
		}
		assert.Equal(t, initG, p.countG())
	})
}
