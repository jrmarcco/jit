package xretry

import (
	"sync/atomic"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

var _ Strategy = (*ExponentialBackoffStrategy)(nil)

// ExponentialBackoffStrategy 为指数退避重试策略:
// 间隔从 initialInterval 起，每次重试翻倍，封顶 maxInterval。
// maxRetryTime <= 0 表示不限次数。
type ExponentialBackoffStrategy struct {
	initInterval time.Duration
	maxInterval  time.Duration

	// reachMaxInterval 标记间隔已达上限，后续直接返回 maxInterval。
	reachMaxInterval atomic.Bool

	maxTimes     int32
	retriedTimes int32
}

// NewExponentialBackoffStrategy 创建指数退避策略。
// initialInterval 必须大于 0 且不超过 maxInterval。
func NewExponentialBackoffStrategy(initialInterval, maxInterval time.Duration, maxRetryTime int32) (*ExponentialBackoffStrategy, error) {
	if initialInterval <= 0 {
		return nil, errs.NewErrInvalidInterval(initialInterval)
	}

	if initialInterval > maxInterval {
		return nil, errs.NewErrInvalidMaxInterval(initialInterval, maxInterval)
	}

	return &ExponentialBackoffStrategy{
		initInterval: initialInterval,
		maxInterval:  maxInterval,
		maxTimes:     maxRetryTime,
	}, nil
}

// Next 返回本次重试的等待间隔，超过 maxTimes 次后不再允许重试。
func (e *ExponentialBackoffStrategy) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&e.retriedTimes, 1)
	if e.maxTimes <= 0 || retriedTimes <= e.maxTimes {
		if e.reachMaxInterval.Load() {
			return e.maxInterval, true
		}

		// interval = initialInterval * 2^(retriedTimes-1)。
		// 移位超出 int64 范围时间隔必然超过 maxInterval，
		// 直接封顶，不依赖溢出行为。
		const maxShift = 62
		if retriedTimes-1 > maxShift {
			e.reachMaxInterval.Store(true)
			return e.maxInterval, true
		}
		interval := e.initInterval << (retriedTimes - 1)

		// interval <= 0 覆盖移位后溢出为负数的情况。
		if interval <= 0 || interval > e.maxInterval {
			e.reachMaxInterval.Store(true)
			return e.maxInterval, true
		}
		return interval, true
	}

	return 0, false
}
