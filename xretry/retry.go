package xretry

import (
	"context"
	"time"

	"github.com/jrmarcco/jit/internal/errs"
)

// Retry 以 strategy 控制的间隔反复执行 bizFunc，直到成功、重试耗尽或 ctx 结束。
//   - bizFunc 返回 nil: 立即返回 nil。
//   - strategy 不再允许重试: 返回包装最近一次错误的 errs.ErrRetryTimeExhausted。
//   - ctx 超时或取消: 返回 ctx.Err()。
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
			return errs.NewErrRetryTimeExhausted(err)
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
