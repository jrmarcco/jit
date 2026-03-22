package xretry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExponentialBackoffRetry_Next(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		strategy *ExponentialBackoffStrategy
		wantDur  time.Duration
		wantRes  bool
	}{
		{
			name: "basic",
			strategy: &ExponentialBackoffStrategy{
				initInterval: time.Second,
				maxInterval:  time.Minute,
				maxTimes:     3,
				retriedTimes: 0,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "over max retry time",
			strategy: &ExponentialBackoffStrategy{
				initInterval: time.Second,
				maxInterval:  time.Minute,
				maxTimes:     3,
				retriedTimes: 3,
			},
			wantDur: 0,
			wantRes: false,
		}, {
			name: "initial interval is 0",
			strategy: &ExponentialBackoffStrategy{
				initInterval: 0,
				maxInterval:  time.Minute,
				maxTimes:     3,
				retriedTimes: 0,
			},
			wantDur: time.Minute,
			wantRes: true,
		}, {
			name: "over max interval",
			strategy: &ExponentialBackoffStrategy{
				initInterval: time.Second,
				maxInterval:  time.Minute,
				maxTimes:     10,
				retriedTimes: 8,
			},
			wantDur: time.Minute,
			wantRes: true,
		}, {
			name: "in max interval",
			strategy: &ExponentialBackoffStrategy{
				initInterval: time.Second,
				maxInterval:  time.Minute,
				maxTimes:     10,
				retriedTimes: 2,
			},
			wantDur: 4 * time.Second,
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

func TestExponentialBackoffStrategy_NextWithRetried(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     *ExponentialBackoffStrategy
		retriedTimes int32
		wantDur      time.Duration
		wantRes      bool
	}{
		{
			name: "basic",
			strategy: func() *ExponentialBackoffStrategy {
				s, err := NewExponentialBackoffStrategy(time.Second, time.Minute, 3)
				require.NoError(t, err)
				return s
			}(),
			retriedTimes: 1,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name: "equal max retry time",
			strategy: func() *ExponentialBackoffStrategy {
				s, err := NewExponentialBackoffStrategy(time.Second, time.Minute, 3)
				require.NoError(t, err)
				return s
			}(),
			retriedTimes: 3,
			wantDur:      4 * time.Second,
			wantRes:      true,
		}, {
			name: "over max retry time",
			strategy: func() *ExponentialBackoffStrategy {
				s, err := NewExponentialBackoffStrategy(time.Second, time.Minute, 3)
				require.NoError(t, err)
				return s
			}(),
			retriedTimes: 4,
			wantDur:      0,
			wantRes:      false,
		}, {
			name: "over max interval",
			strategy: func() *ExponentialBackoffStrategy {
				s, err := NewExponentialBackoffStrategy(time.Second, time.Minute, 10)
				require.NoError(t, err)
				return s
			}(),
			retriedTimes: 8,
			wantDur:      time.Minute,
			wantRes:      true,
		}, {
			name: "in max interval",
			strategy: func() *ExponentialBackoffStrategy {
				s, err := NewExponentialBackoffStrategy(time.Second, time.Minute, 10)
				require.NoError(t, err)
				return s
			}(),
			retriedTimes: 3,
			wantDur:      4 * time.Second,
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
