package xretry

import (
	"math/bits"
	"sync/atomic"
	"time"
)

var _ Strategy = (*AdaptiveTimeoutStrategy)(nil)

type AdaptiveTimeoutStrategy struct {
	strategy  Strategy // basic retry strategy
	threshold int      // timeout threshold

	bufferSize int      // size of the slide window
	ringBuffer []uint64 // using as a slide window to store timeout information

	totalBit uint64
}

func NewAdaptiveTimeoutStrategy(strategy Strategy, bufferSize, threshold int) *AdaptiveTimeoutStrategy {
	const defaultTotalBit = 64
	return &AdaptiveTimeoutStrategy{
		strategy:   strategy,
		threshold:  threshold,
		bufferSize: bufferSize,
		ringBuffer: make([]uint64, bufferSize),
		totalBit:   uint64(defaultTotalBit) & uint64(bufferSize),
	}
}

func (a *AdaptiveTimeoutStrategy) Next() (time.Duration, bool) {
	failureCnt := a.getFailureCnt()
	if failureCnt >= a.threshold {
		return 0, false
	}
	return a.strategy.Next()
}

func (a *AdaptiveTimeoutStrategy) NextWithRetried(retriedTimes int32) (time.Duration, bool) {
	failureCnt := a.getFailureCnt()
	if failureCnt >= a.threshold {
		return 0, false
	}
	return a.strategy.NextWithRetried(retriedTimes)
}

func (a *AdaptiveTimeoutStrategy) Report(err error) Strategy {
	if err == nil {
		a.markAsSuccess()
		return a
	}

	a.markAsFailure()
	return a
}

func (a *AdaptiveTimeoutStrategy) markAsSuccess() {}

func (a *AdaptiveTimeoutStrategy) markAsFailure() {}

func (a *AdaptiveTimeoutStrategy) getFailureCnt() int {
	var cnt int
	for i := 0; i < a.bufferSize; i++ {
		val := atomic.LoadUint64(&a.ringBuffer[i])
		cnt += bits.OnesCount64(val)
	}

	return cnt
}
