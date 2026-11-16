package xretry

import (
	"sync/atomic"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

var _ Strategy = (*FixedIntervalStrategy)(nil)

// FixedIntervalStrategy 为固定间隔重试策略。
// maxTimes <= 0 表示不限次数。
type FixedIntervalStrategy struct {
	interval     time.Duration
	maxTimes     int32
	retriedTimes int32
}

// NewFixedIntervalStrategy 创建固定间隔策略，interval 必须大于 0。
func NewFixedIntervalStrategy(interval time.Duration, maxTimes int32) (*FixedIntervalStrategy, error) {
	if interval <= 0 {
		return nil, errs.NewErrInvalidInterval(interval)
	}
	return &FixedIntervalStrategy{
		interval:     interval,
		maxTimes:     maxTimes,
		retriedTimes: 0,
	}, nil
}

// Next 返回固定的重试间隔，超过 maxTimes 次后不再允许重试。
func (f *FixedIntervalStrategy) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&f.retriedTimes, 1)
	if f.maxTimes <= 0 || retriedTimes <= f.maxTimes {
		return f.interval, true
	}
	return 0, false
}
