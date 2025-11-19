package errs

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrSameRBNode   = errors.New("[jit] cannot insert same red-black tree node")
	ErrNodeNotFound = errors.New("[jit] cannot find node in red-black tree")
)

func NilErr(name string) error {
	return fmt.Errorf("[jit] %s is nil", name)
}

func ErrIndexOutOfBounds(length, index int) error {
	return fmt.Errorf("[jit] index %d out of bounds for length %d", index, length)
}

func ErrEmptySlice() error {
	return fmt.Errorf("[jit] slice is empty")
}

func ErrInvalidKeyValLen() error {
	return fmt.Errorf("[jit] keys and vals have different lengths")
}

func ErrInvalidInterval(interval time.Duration) error {
	return fmt.Errorf("[jit] invalid interval: %v, expected interval value should greater than 0", interval)
}

func ErrInvalidMaxInterval(maxInterval time.Duration) error {
	return fmt.Errorf(
		"[jit] invalid max interval: %v, expected max interval value should greater than init interval",
		maxInterval,
	)
}

func ErrRetryTimeExhausted(latestErr error) error {
	return fmt.Errorf("[jit] retry time exhausted, the latest error: %w", latestErr)
}
