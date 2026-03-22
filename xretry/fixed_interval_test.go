package xretry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedIntervalRetry_Next(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		strategy *FixedIntervalStrategy
		wantDur  time.Duration
		wantRes  bool
	}{
		{
			name: "basic",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     3,
				retriedTimes: 0,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "over max retry time",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     3,
				retriedTimes: 3,
			},
			wantDur: 0,
			wantRes: false,
		}, {
			name: "max retry time is 0",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     0,
				retriedTimes: 3,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "max retry time is negative",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     -1,
				retriedTimes: 3,
			},
			wantDur: time.Second,
			wantRes: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dur, ok := tc.strategy.Next()
			assert.Equal(t, tc.wantRes, ok)
			if ok {
				assert.Equal(t, tc.wantDur, dur)
			}
		})
	}
}

func TestFixedIntervalRetry_NextWithRetried(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     *FixedIntervalStrategy
		retriedTimes int32
		wantDur      time.Duration
		wantRes      bool
	}{
		{
			name: "basic",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     3,
				retriedTimes: 0,
			},
			retriedTimes: 0,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name: "equal max retry time",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     3,
				retriedTimes: 0,
			},
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name: "over max retry time",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     3,
				retriedTimes: 0,
			},
			retriedTimes: 4,
			wantDur:      0,
			wantRes:      false,
		}, {
			name: "max retry time is 0",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     0,
				retriedTimes: 0,
			},
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name: "max retry time is negative",
			strategy: &FixedIntervalStrategy{
				interval:     time.Second,
				maxTimes:     -1,
				retriedTimes: 0,
			},
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dur, ok := tc.strategy.NextWithRetried(tc.retriedTimes)
			assert.Equal(t, tc.wantRes, ok)
			if ok {
				assert.Equal(t, tc.wantDur, dur)
			}
		})
	}
}
