package xretry

import (
	"context"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

func Retry(ctx context.Context, strategy Strategy, bizFunc func() error) error {
	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()

	for {
		err := bizFunc()
		if err == nil {
			return nil
		}

		next, ok := strategy.Next()
		if !ok {
			return errs.ErrRetryTimeExhausted(err)
		}

		if ticker == nil {
			ticker = time.NewTicker(next)
		}
		ticker.Reset(next)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
