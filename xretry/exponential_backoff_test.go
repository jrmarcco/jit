package xretry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
