// Package errs 定义各内部包共享的哨兵错误与错误构造函数。
//
// 所有构造函数均通过 %w 包装哨兵错误，调用方可以用 errors.Is 匹配。
package errs

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrSameRBNode 为向红黑树插入重复 key 时返回的错误。
	ErrSameRBNode = errors.New("[jit] cannot insert same red-black tree node")

	// ErrNodeNotFound 为红黑树节点不存在时返回的错误。
	ErrNodeNotFound = errors.New("[jit] cannot find node in red-black tree")

	// ErrIndexOutOfBounds 为切片下标越界时返回的错误。
	ErrIndexOutOfBounds = errors.New("[jit] index out of bounds")

	// ErrEmptySlice 为操作要求非空切片时返回的错误。
	ErrEmptySlice = errors.New("[jit] slice is empty")

	// ErrInvalidKeyValLen 为 keys 与 vals 长度不一致时返回的错误。
	ErrInvalidKeyValLen = errors.New("[jit] keys and vals have different lengths")

	// ErrInvalidInterval 为重试间隔非正数时返回的错误。
	ErrInvalidInterval = errors.New("[jit] invalid interval")

	// ErrInvalidMaxInterval 为初始间隔大于最大间隔时返回的错误。
	ErrInvalidMaxInterval = errors.New("[jit] invalid max interval")

	// ErrRetryTimeExhausted 为重试次数耗尽时返回的错误。
	ErrRetryTimeExhausted = errors.New("[jit] retry time exhausted")
)

// NewErrIndexOutOfBounds 返回携带细节的越界错误。
func NewErrIndexOutOfBounds(length, index int) error {
	return fmt.Errorf("%w: index %d out of bounds for length %d", ErrIndexOutOfBounds, index, length)
}

// NewErrEmptySlice 返回空切片错误。
func NewErrEmptySlice() error {
	return ErrEmptySlice
}

// NewErrInvalidKeyValLen 返回 keys/vals 长度不一致错误。
func NewErrInvalidKeyValLen() error {
	return ErrInvalidKeyValLen
}

// NewErrInvalidInterval 返回间隔非法 ( 应大于 0 ) 的错误。
func NewErrInvalidInterval(interval time.Duration) error {
	return fmt.Errorf("%w: %v, expected interval value should be greater than 0", ErrInvalidInterval, interval)
}

// NewErrInvalidMaxInterval 返回初始间隔大于最大间隔的错误。
func NewErrInvalidMaxInterval(initialInterval, maxInterval time.Duration) error {
	return fmt.Errorf(
		"%w: initial interval %v is greater than max interval %v",
		ErrInvalidMaxInterval, initialInterval, maxInterval,
	)
}

// NewErrRetryTimeExhausted 返回重试次数耗尽的错误，
// 并以 %w 包装业务函数最近一次返回的错误。
func NewErrRetryTimeExhausted(latestErr error) error {
	return fmt.Errorf("%w, the latest error: %w", ErrRetryTimeExhausted, latestErr)
}
