package xretry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdaptiveTimeoutRetry_Next(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     Strategy
		retriedTimes int32
		wantNext     time.Duration
		wantRes      bool
	}{
		{
			name: "basic",
			strategy: func() *AdaptiveTimeoutStrategy {
				fis, err := NewFixedIntervalStrategy(time.Second, 5)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(fis, 8, 100)
			}(),
			retriedTimes: 0,
			wantNext:     1 * time.Second,
			wantRes:      true,
		}, {
			name: "after once retry",
			strategy: func() *AdaptiveTimeoutStrategy {
				ebs, err := NewExponentialBackoffStrategy(time.Second, 5*time.Second, 10)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(ebs, 8, 100)
			}(),
			retriedTimes: 1,
			wantNext:     2 * time.Second,
			wantRes:      true,
		}, {
			name: "after 4 times retry",
			strategy: func() *AdaptiveTimeoutStrategy {
				ebs, err := NewExponentialBackoffStrategy(time.Second, 5*time.Second, 10)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(ebs, 8, 100)
			}(),
			retriedTimes: 4,
			wantNext:     5 * time.Second,
			wantRes:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for range tc.retriedTimes {
				_, _ = tc.strategy.Next()
			}

			next, ok := tc.strategy.Next()
			assert.Equal(t, tc.wantRes, ok)

			if ok {
				assert.Equal(t, tc.wantNext, next)
			}
		})
	}
}

func TestAdaptiveTimeoutRetry_NextWithRetried(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     *AdaptiveTimeoutStrategy
		retriedTimes int32
		wantNext     time.Duration
		wantRes      bool
	}{
		{
			name: "fixed interval",
			strategy: func() *AdaptiveTimeoutStrategy {
				fis, err := NewFixedIntervalStrategy(time.Second, 5)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(fis, 8, 100)
			}(),
			retriedTimes: 3,
			wantNext:     1 * time.Second,
			wantRes:      true,
		},
		{
			name: "exponential backoff in range",
			strategy: func() *AdaptiveTimeoutStrategy {
				ebs, err := NewExponentialBackoffStrategy(time.Second, 5*time.Second, 10)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(ebs, 8, 100)
			}(),
			retriedTimes: 3,
			wantNext:     4 * time.Second,
			wantRes:      true,
		},
		{
			name: "exponential backoff over max interval",
			strategy: func() *AdaptiveTimeoutStrategy {
				ebs, err := NewExponentialBackoffStrategy(time.Second, 5*time.Second, 10)
				require.NoError(t, err)
				return NewAdaptiveTimeoutStrategy(ebs, 8, 100)
			}(),
			retriedTimes: 10,
			wantNext:     5 * time.Second,
			wantRes:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			next, ok := tc.strategy.NextWithRetried(tc.retriedTimes)
			assert.Equal(t, tc.wantRes, ok)
			if ok {
				assert.Equal(t, tc.wantNext, next)
			}
		})
	}
}

func ExampleAdaptiveTimeoutStrategy() {
	ebs, err := NewExponentialBackoffStrategy(time.Second, 30*time.Second, 10)
	if err != nil {
		panic(err)
	}
	s := NewAdaptiveTimeoutStrategy(ebs, 8, 100)

	next, ok := s.Next()
	for ok {
		fmt.Println(next)
		next, ok = s.Next()
	}
	// Output:
	// 1s
	// 2s
	// 4s
	// 8s
	// 16s
	// 30s
	// 30s
	// 30s
	// 30s
	// 30s
}
