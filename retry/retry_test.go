package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JrMarcco/jit/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Parallel()

	bizErr := errors.New("biz error")

	tcs := []struct {
		name     string
		bizFunc  func() error
		strategy Strategy
		wantErr  error
	}{
		{
			name: "successful first execution",
			bizFunc: func() error {
				t.Logf("a business logic")
				return nil
			},
			strategy: func() *FixedIntervalStrategy {
				s, err := NewFixedIntervalStrategy(time.Second, 5)
				require.NoError(t, err)
				return s
			}(),
		}, {
			name: "failed after max retry times",
			bizFunc: func() error {
				t.Logf("a business logic")
				return bizErr
			},
			strategy: func() *FixedIntervalStrategy {
				s, err := NewFixedIntervalStrategy(time.Second, 2)
				require.NoError(t, err)
				return s
			}(),
			wantErr: errs.ErrRetryTimeExhausted(bizErr),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			err := Retry(ctx, tc.strategy, tc.bizFunc)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func ExampleRetry() {
	bizErr := errors.New("biz error")
	bizFunc := func() error {
		fmt.Println("a business logic")
		return bizErr
	}

	s, err := NewExponentialBackoffStrategy(time.Second, 5*time.Second, 1)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = Retry(ctx, s, bizFunc)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// a business logic
	// a business logic
	// [jit] retry time exhausted, the latest error: biz error
}
