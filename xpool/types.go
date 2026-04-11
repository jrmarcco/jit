package xpool

import (
	"context"
	"time"
)

// Task 任务
type Task interface {
	Run(ctx context.Context) error
}

// TaskPool 任务池
type TaskPool interface {
	// Submit 提交任务。
	Submit(ctx context.Context, task Task) error

	// Start 开始调度任务。
	Start() error

	// Shutdown 关闭任务池，优雅退出实现。
	// 调用 Shutdown 后任务池会拒绝新任务提交，并继续执行剩下的任务。
	// 所有任务执行完成后通过 chan 通知调用者。
	// 注意：
	//	通知发出后需要关闭 chan。
	Shutdown() (<-chan struct{}, error)
	// ShutdownNow 立即关闭任务池，并返回剩余的任务。
	// 当前正在执行的任务是否会被中断取决于 TaskPool 和 Task 的具体的实现。
	ShutdownNow() ([]Task, error)

	// State 暴露 Pool 生命周期内的运行状态。
	// interval 为检查状态的时间间隔。
	// state chan 创建失败时返回 error。
	State(ctx context.Context, interval time.Duration) (<-chan State, error)
}

type State struct {
	QueueSize int32

	GoroutineCnt int32

	WaitingCnt int32
	RunningCnt int32

	SubmitRetryCnt int64
	TaskTimeoutCnt int64
	ErrDropCnt     int64
	StateDropCnt   int64

	PoolState int32

	Timestamp int64
}
