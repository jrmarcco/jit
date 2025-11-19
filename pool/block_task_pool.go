package pool

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/JrMarcco/jit/bean/option"
)

const (
	panicBuffLen            = 2048
	defaultMaxIdleTime      = 10 * time.Second
	defaultSubmitTimeout    = 15 * time.Second
	defaultErrHandleTimeout = 3 * time.Second
)

const (
	stateCreated int32 = iota
	stateRunning
	stateClosing
	stateClosed
	stateLocked
)

var (
	errTaskRunningPanic = fmt.Errorf("[jit] panic when running task, stack")

	errInvalidParam = fmt.Errorf("[jit] invalid param")
	errInvalidTask  = fmt.Errorf("[jit] invalid task")

	errPoolIsNotRunning = fmt.Errorf("[jit] task pool is not running")
	errPoolIsRunning    = fmt.Errorf("[jit] task pool is running")
	errPoolIsClosing    = fmt.Errorf("[jit] task pool is closing")
	errPoolIsClosed     = fmt.Errorf("[jit] task pool is closed")
	errPoolIsLocked     = fmt.Errorf("[jit] task pool is locked")
)

var _ Task = (*TaskFunc)(nil)

type TaskFunc func(ctx context.Context) error

func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

type taskWrapper struct {
	task Task
}

func (t *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]

			slog.Error(
				"[jit] panic when running task",
				"panic", r,
				"stack", string(buf),
			)

			err = fmt.Errorf("%w: %+v", errTaskRunningPanic, r)
		}
	}()
	return t.task.Run(ctx)
}

// timeoutGoroutine 超时组，管理任务池中超时的 goroutine id。
type timeoutGoroutine struct {
	mu sync.RWMutex

	cnt   int32
	idMap map[int32]struct{}
}

func (g *timeoutGoroutine) in(id int32) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.idMap[id]
	return ok
}

func (g *timeoutGoroutine) add(id int32) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.idMap[id]; !ok {
		g.idMap[id] = struct{}{}
		g.cnt++
	}
}

func (g *timeoutGoroutine) del(id int32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.idMap[id]; ok {
		delete(g.idMap, id)
		g.cnt--
	}
}

func (g *timeoutGoroutine) size() int32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.cnt
}

var _ TaskPool = (*BlockTaskPool)(nil)

// BlockTaskPool 并发阻塞的任务池。
// 任务池会动态控制 goroutine 数量，按需创建 goroutine 执行任务。
type BlockTaskPool struct {
	mu sync.RWMutex

	id int32 // goroutine id

	maxIdleTime   time.Duration
	submitTimeout time.Duration

	state         int32 // 内部状态
	totalG        int32 // goroutine 总数
	totalRunningG int32 // 正在执行任务的 goroutine 总数

	// 参数 initG / coreG / maxG 的作用是分层管理 goroutine。
	// 三个参数将 goroutine 分为 3 个区间:
	//
	//		[1	  , initG]:	永久 goroutine。
	//		(initG, coreG]: 核心 goroutine（带超时机制的 goroutine）
	//						当 goroutine 处于这个区间，在退出前（maxIdleTime）尝试拿任务，
	//						拿到任务则继续执行，没拿到则超时退出。
	//		(coreG, maxG ]: 临时 goroutine
	//						当 goroutine 处于这个区间且当前对立没有可执行任务，则快速退出。
	initG int32 // 初始 goroutine 数量
	coreG int32 // 核心 goroutine 数量
	maxG  int32 // 最大 goroutine 数量

	queue            chan Task // 任务队列
	queueBacklogRate float64   // 任务队列积压率

	timeoutG *timeoutGoroutine // 超时的 goroutine

	interruptCtx        context.Context
	interruptCancelFunc context.CancelFunc

	errHandler       func(ctx context.Context, err error) // 错误处理器
	errHandleTimeout time.Duration
}

// Submit 提交一个任务。
// 在队列已满的情况下，调用者会被阻塞。
// 在 Start 方法被调用后仍然可以调用 Submit 方法。
func (p *BlockTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return errInvalidTask
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.submitTimeout)
		defer cancel()
	}

	for {
		if atomic.LoadInt32(&p.state) == stateClosing {
			return errPoolIsClosing
		}
		if atomic.LoadInt32(&p.state) == stateClosed {
			return errPoolIsClosed
		}

		tw := &taskWrapper{
			task: task,
		}

		ok, err := p.trySubmit(ctx, tw, stateCreated)
		if ok || err != nil {
			return err
		}

		ok, err = p.trySubmit(ctx, tw, stateRunning)
		if ok || err != nil {
			return err
		}
	}
}

// trySubmit 尝试提交一个任务。
func (p *BlockTaskPool) trySubmit(ctx context.Context, task Task, state int32) (bool, error) {
	// 锁定 task pool。
	if atomic.CompareAndSwapInt32(&p.state, state, stateLocked) {
		// 当 trySubmit 成功返回时解除锁定 task pool。
		defer atomic.CompareAndSwapInt32(&p.state, stateLocked, state)

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case p.queue <- task:
			if state == stateRunning && p.allowToCreateG() {
				// 任务池处于运行状态且允许创建新 goroutine 执行任务。
				p.increaseG(1)
				id := atomic.AddInt32(&p.id, 1)
				go p.newG(id)

				slog.Info("[jit] create new goroutine", "id", id)
			}

			// 任务池还未运行 或 当前不允许创建 goroutine，直接成功提交。
			return true, nil
		default:
			return false, nil
		}
	}
	return false, nil
}

// allowToCreateG 判断当前是否允许创建新的 goroutine。
// 当满足以下条件时：
//
//	1、goroutine 总数小于最大 goroutine 数量；
//	2、队列存在待运行的 task 且队列积压率达到阈值。
//
// 此时允许创建新的 goroutine 执行任务。
func (p *BlockTaskPool) allowToCreateG() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.totalG >= p.maxG {
		return false
	}

	// 计算队列占用率
	rate := float64(len(p.queue)) / float64(cap(p.queue))

	// 队列存在待运行的 task 且队列积压率达到阈值
	return rate != 0 && rate >= p.queueBacklogRate
}

// newG 创建新的 goroutine， 参数 id 用来表示新创建的 goroutine。
func (p *BlockTaskPool) newG(id int32) {
	// 创建一个持续时间为 0 的 Timer（假超时），该 Timer 会立即过期并向其 channel 发送信号。
	//
	// 注意：
	// 	timer 只保证在等待 x 时间后才发送信号，而不是在 x 时间内发送信号。
	//
	// 这里假超时的目的是保证除任务池退出的情况外， goroutine 至少执行一个任务。
	// 同时这里不能使用 nil timer，会导致 for 循环内的 case <-idleTImer.C 发生 panic。
	idleTimer := time.NewTimer(0)
	if !idleTimer.Stop() {
		// 从 channel 中读取并丢弃该信号，避免假超时导致 goroutine 退出
		<-idleTimer.C
	}

	for {
		select {
		case <-p.interruptCtx.Done():
			// 收到整个 task pool 的中断信号
			p.decreaseG(1)
			return

		case <-idleTimer.C:
			// 空闲时收到超时信号，即 goroutine 在 maxIdleTime 时间内没获取到可执行任务
			p.handleIdleTimeout(id)
			return

		case task, ok := <-p.queue:
			if !p.processTask(id, task, ok, idleTimer) {
				return
			}
		}
	}
}

// handleIdleTimeout 处理空闲超时。
func (p *BlockTaskPool) handleIdleTimeout(id int32) {
	p.mu.Lock()
	// 任务池 goroutine 总数 -1
	p.totalG--
	// 从超时组移除当前 goroutine id
	p.timeoutG.del(id)
	p.mu.Unlock()
}

// processTask 处理任务的核心流程。
func (p *BlockTaskPool) processTask(id int32, task Task, ok bool, idleTimer *time.Timer) bool {
	// 当 p.queue 被 close，这里会收到零值 task 和 false
	if p.timeoutG.in(id) {
		// 当前 goroutine 在超时组中，且在超时前成功拿到任务执行
		p.timeoutG.del(id)
		if !idleTimer.Stop() {
			<-idleTimer.C
		}
	}

	if !ok {
		// !ok 意味着任务队列被关闭
		p.decreaseG(1)
		// 任务池中没有 goroutine
		if p.countG() == 0 {
			// 因 shutdown 导致的 goroutine 退出，
			// 最后一个退出的 goroutine 需要负责状态迁移，并通知外部调用者。
			if atomic.CompareAndSwapInt32(&p.state, stateClosing, stateClosed) {
				// 调用 context.CancelFunc 通知外部调用者
				p.interruptCancelFunc()
			}
		}
		return false
	}

	// 成功获取可执行任务
	atomic.AddInt32(&p.totalRunningG, 1)
	err := task.Run(p.interruptCtx)
	atomic.AddInt32(&p.totalRunningG, -1)

	// 处理任务执行错误
	if err != nil && p.errHandler != nil {
		p.handleTaskError(err)
	}

	return p.shouldContinue(id, idleTimer)
}

// handleTaskError 处理任务执行错误。
// 在独立的 goroutine 中调用错误处理器，避免错误处理器发生 panic 影响任务池的运行。
func (p *BlockTaskPool) handleTaskError(err error) {
	// 在独立的 goroutine 中调用错误 errHandler，
	// 避免 errHandler 发生 panic 影响任务池的运行。
	go func(ctx context.Context, err error) {
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		// 超时控制，避免 goroutine 泄露
		ctx, cancel := context.WithTimeout(ctx, p.errHandleTimeout)
		p.errHandler(ctx, err)
		cancel()
	}(p.interruptCtx, err)
}

// shouldContinue 封装了任务执行后判断是否继续运行的策略（包含缩容和重置超时）。
func (p *BlockTaskPool) shouldContinue(id int32, idleTimer *time.Timer) bool {
	// 任务执行完成后的判断。
	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查队列中是否还有任务需要执行。
	noTaskToExec := len(p.queue) == 0 || int32(len(p.queue)) < p.totalG
	// 临时 goroutine 的快速退出策略
	if noTaskToExec && p.coreG < p.totalG && p.totalG <= p.maxG {
		// 当前 goroutine 处于 (coreG, maxG] 区间（即临时 goroutine），直接退出 goroutine。
		p.totalG--
		return false
	}

	// p.totalG-p.timeoutG.size() -> 当前活跃的 goroutine 数
	// 核心 goroutine 的超时管理，为属于 (initG, coreG] 区间的 goroutine 设置超时器。
	if p.initG < p.totalG-p.timeoutG.size() {
		// 核心 goroutine 不立即退出能保证在一定时间（maxIdleTime）由任务提交带来的扩容，保持核心处理能力。
		idleTimer.Reset(p.maxIdleTime)
		p.timeoutG.add(id)
	}
	return true
}

func (p *BlockTaskPool) increaseG(delta int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalG += delta
}

func (p *BlockTaskPool) decreaseG(delta int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalG -= delta
}

// countG 查看任务池中有多少 goroutine。
func (p *BlockTaskPool) countG() int32 {
	var cnt int32
	p.mu.RLock()
	cnt = p.totalG
	p.mu.RUnlock()
	return cnt
}

// Start 开始调度执行。
func (p *BlockTaskPool) Start() error {
	for {
		if atomic.LoadInt32(&p.state) == stateClosing {
			return errPoolIsClosing
		}
		if atomic.LoadInt32(&p.state) == stateClosed {
			return errPoolIsClosed
		}
		if atomic.LoadInt32(&p.state) == stateRunning {
			return errPoolIsRunning
		}
		if atomic.LoadInt32(&p.state) == stateLocked {
			return errPoolIsLocked
		}

		if atomic.CompareAndSwapInt32(&p.state, stateCreated, stateLocked) {
			// 计算允许创建的 goroutine 数量。
			cntG := p.initG

			// 需求的 goroutine 数 = 队列中任务数 - 初始 goroutine 数。
			needG := int32(len(p.queue)) - p.initG
			if needG > 0 {
				// 允许创建的最大 goroutine 数
				allowMaxG := p.maxG - p.initG
				if needG <= allowMaxG {
					cntG += needG
				} else {
					cntG += allowMaxG
				}
			}

			p.increaseG(cntG)
			for i := int32(0); i < cntG; i++ {
				go p.newG(atomic.AddInt32(&p.id, 1))
			}
			atomic.CompareAndSwapInt32(&p.state, stateLocked, stateRunning)
			return nil
		}
	}
}

// Shutdown 关闭任务池。
// 调用后将拒绝 Submit 调用，但会继续执行队列中剩下的任务。
// 所有任务执行完成后会发送信号到 chan 并负责关闭 chan。
func (p *BlockTaskPool) Shutdown() (<-chan struct{}, error) {
	for {
		if atomic.LoadInt32(&p.state) == stateCreated {
			return nil, errPoolIsNotRunning
		}
		if atomic.LoadInt32(&p.state) == stateClosed {
			return nil, errPoolIsClosed
		}
		if atomic.LoadInt32(&p.state) == stateClosing {
			return nil, errPoolIsClosing
		}

		if atomic.CompareAndSwapInt32(&p.state, stateRunning, stateClosing) {
			// 关闭任务队列，拒绝新任务提交。
			// 注意：
			//  close(p.queue) 只是把 chan 标记为“关闭”状态，
			//	此时工作 goroutine 还可以读取 chan 的剩余数据，知道 chan 的数据全被取走。
			close(p.queue)
			return p.interruptCtx.Done(), nil
		}
	}
}

// ShutdownNow 立即关闭任务池，并返回剩余的任务（不包含执行中的任务也）。
func (p *BlockTaskPool) ShutdownNow() ([]Task, error) {
	for {
		if atomic.LoadInt32(&p.state) == stateCreated {
			return nil, errPoolIsNotRunning
		}
		if atomic.LoadInt32(&p.state) == stateClosed {
			return nil, errPoolIsClosed
		}
		if atomic.LoadInt32(&p.state) == stateClosing {
			return nil, errPoolIsClosing
		}

		if atomic.CompareAndSwapInt32(&p.state, stateRunning, stateClosed) {
			close(p.queue)
			p.interruptCancelFunc()

			tasks := make([]Task, 0, len(p.queue))
			for task := range p.queue {
				tasks = append(tasks, task)
			}
			return tasks, nil
		}
	}
}

// State 查询任务池内部状态。
func (p *BlockTaskPool) State(ctx context.Context, interval time.Duration) (<-chan State, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if p.interruptCtx.Err() != nil {
		return nil, p.interruptCtx.Err()
	}

	stateChan := make(chan State)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case timestamp := <-ticker.C:
				p.sendState(stateChan, timestamp.UnixMilli())
			case <-ctx.Done():
				p.sendState(stateChan, time.Now().UnixMilli())
				close(stateChan)
				return
			case <-p.interruptCtx.Done():
				p.sendState(stateChan, time.Now().UnixMilli())
				close(stateChan)
				return
			}
		}
	}()
	return stateChan, nil
}

func (p *BlockTaskPool) sendState(ch chan<- State, timestamp int64) {
	select {
	case ch <- p.getState(timestamp):
	default:
		// 发送失败直接丢弃。
	}
}

func (p *BlockTaskPool) getState(timestamp int64) State {
	return State{
		QueueSize:    int32(cap(p.queue)),
		GoroutineCnt: p.countG(),
		WaitingCnt:   int32(len(p.queue)),
		RunningCnt:   atomic.LoadInt32(&p.totalRunningG),
		PoolState:    atomic.LoadInt32(&p.state),
		Timestamp:    timestamp,
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.maxIdleTime = maxIdleTime
	}
}

func WithSubmitTimeout(submitTimeout time.Duration) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.submitTimeout = submitTimeout
	}
}

func WithCoreG(coreG int32) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.coreG = coreG
	}
}

func WithMaxG(maxG int32) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.maxG = maxG
	}
}

func WithQueueBacklogRate(queueBacklogRate float64) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.queueBacklogRate = queueBacklogRate
	}
}

func WithErrorHandler(errHandler func(ctx context.Context, err error)) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.errHandler = errHandler
	}
}

func WithErrHandleTimeout(errHandleTimeout time.Duration) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.errHandleTimeout = errHandleTimeout
	}
}

// NewBlockTaskPool 创建任务池。
func NewBlockTaskPool(initG, queueSize int32, opts ...option.Opt[BlockTaskPool]) (*BlockTaskPool, error) {
	if initG <= 0 {
		return nil, fmt.Errorf("%w: init goroutine should be greater than 0", errInvalidParam)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("%w: queue size should be greater or equal to 0", errInvalidParam)
	}

	p := &BlockTaskPool{
		queue:            make(chan Task, queueSize),
		initG:            initG,
		coreG:            initG,
		maxG:             initG,
		maxIdleTime:      defaultMaxIdleTime,
		submitTimeout:    defaultSubmitTimeout,
		errHandleTimeout: defaultErrHandleTimeout,
	}

	ctx := context.Background()
	p.interruptCtx, p.interruptCancelFunc = context.WithCancel(ctx)
	atomic.StoreInt32(&p.state, stateCreated)

	option.Apply(p, opts...)

	// 默认情况 coreG == maxG。
	// 当 coreG == maxG 时，goroutine 的分层会简化为两层:
	//
	//		[1	 , initG]: 永久 goroutine
	//		(initG, maxG]: 带超时机制的核心 goroutine
	//
	// 这样使 goroutine 的数量更平滑（只有 initG -> maxG 的渐进式增长），资源更可预测。
	// 所有核心 goroutine 都有相同的生命周期管理。
	if p.coreG != p.initG && p.maxG == p.initG {
		// 只使用 option WithCoreG 的情况，保证 maxG 至少等于 coreG。
		p.maxG = p.coreG
	} else if p.coreG == p.initG && p.maxG != p.initG {
		// 只使用 option WithMaxG 的情况
		p.coreG = p.maxG
	}

	if !(p.initG <= p.coreG && p.coreG <= p.maxG) {
		return nil, fmt.Errorf("%w: goroutine required to satisfy [ init <= core <= max ]", errInvalidParam)
	}

	p.timeoutG = &timeoutGoroutine{
		idMap: make(map[int32]struct{}),
	}
	if p.queueBacklogRate < float64(0) || p.queueBacklogRate > float64(1) {
		return nil, fmt.Errorf("%w: queue backlog rate should be in [0, 1]", errInvalidParam)
	}
	return p, nil
}
