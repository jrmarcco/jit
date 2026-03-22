package xretry

import (
	"sync/atomic"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

var _ Strategy = (*FixedIntervalStrategy)(nil)

type FixedIntervalStrategy struct {
	interval     time.Duration
	maxTimes     int32
	retriedTimes int32
}

func NewFixedIntervalStrategy(interval time.Duration, maxTimes int32) (*FixedIntervalStrategy, error) {
	if interval <= 0 {
		return nil, errs.ErrInvalidInterval(interval)
	}
	return &FixedIntervalStrategy{
		interval:     interval,
		maxTimes:     maxTimes,
		retriedTimes: 0,
	}, nil
}

func (f *FixedIntervalStrategy) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&f.retriedTimes, 1)
	return f.nextRetry(retriedTimes)
}

func (f *FixedIntervalStrategy) NextWithRetried(retriedTimes int32) (time.Duration, bool) {
	return f.nextRetry(retriedTimes)
}

func (f *FixedIntervalStrategy) nextRetry(retriedTimes int32) (time.Duration, bool) {
	if f.maxTimes <= 0 || retriedTimes <= f.maxTimes {
		return f.interval, true
	}
	return 0, false
}

func (f *FixedIntervalStrategy) Report(_ error) Strategy {
	return f
}
