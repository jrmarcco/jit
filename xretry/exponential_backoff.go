package xretry

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

var _ Strategy = (*ExponentialBackoffStrategy)(nil)

type ExponentialBackoffStrategy struct {
	initInterval time.Duration
	maxInterval  time.Duration

	reachMaxInterval atomic.Bool

	maxTimes     int32
	retriedTimes int32
}

func NewExponentialBackoffStrategy(initialInterval, maxInterval time.Duration, maxRetryTime int32) (*ExponentialBackoffStrategy, error) {
	if initialInterval <= 0 {
		return nil, errs.ErrInvalidInterval(initialInterval)
	}

	if initialInterval > maxInterval {
		return nil, errs.ErrInvalidMaxInterval(initialInterval)
	}

	return &ExponentialBackoffStrategy{
		initInterval: initialInterval,
		maxInterval:  maxInterval,
		maxTimes:     maxRetryTime,
	}, nil
}

func (e *ExponentialBackoffStrategy) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&e.retriedTimes, 1)
	return e.nextRetry(retriedTimes)
}

func (e *ExponentialBackoffStrategy) NextWithRetried(retriedTimes int32) (time.Duration, bool) {
	return e.nextRetry(retriedTimes)
}

func (e *ExponentialBackoffStrategy) nextRetry(retriedTimes int32) (time.Duration, bool) {
	if e.maxTimes <= 0 || retriedTimes <= e.maxTimes {
		if e.reachMaxInterval.Load() {
			return e.maxInterval, true
		}

		const two = 2
		interval := e.initInterval * time.Duration(math.Pow(two, float64(retriedTimes-1)))

		// interval = 0 prevents an input interval = 0 when create strategy.
		// interval < 0 means the interval is over max int32 value after math.Pow.
		if interval <= 0 || interval > e.maxInterval {
			e.reachMaxInterval.Store(true)
			return e.maxInterval, true
		}
		return interval, true
	}

	return 0, false
}

func (e *ExponentialBackoffStrategy) Report(_ error) Strategy {
	return e
}
