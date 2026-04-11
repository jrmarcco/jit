package xpool

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jrmarcco/jit/xbean/option"
)

const (
	panicBuffLen = 2048

	defaultMaxIdleTime      = 10 * time.Second
	defaultSubmitTimeout    = 15 * time.Second
	defaultErrHandleTimeout = 3 * time.Second

	defaultTaskExecTimeout = 0

	defaultSubmitBackoff = 100 * time.Microsecond
	maxSubmitBackoff     = 2 * time.Millisecond

	defaultErrQueueSize  = 128
	defaultErrWorkerCnt  = 1
	defaultStateChanSize = 1
)

const (
	stateCreated int32 = iota
	stateRunning
	stateClosing
	stateClosed
	stateLocked
)

var (
	errTaskRunningPanic = fmt.Errorf("[xpool] panic when running task, stack")

	errInvalidParam = fmt.Errorf("[xpool] invalid param")
	errInvalidTask  = fmt.Errorf("[xpool] invalid task")

	errPoolIsNotRunning = fmt.Errorf("[xpool] task pool is not running")
	errPoolIsRunning    = fmt.Errorf("[xpool] task pool is running")
	errPoolIsClosing    = fmt.Errorf("[xpool] task pool is closing")
	errPoolIsClosed     = fmt.Errorf("[xpool] task pool is closed")
	errPoolIsLocked     = fmt.Errorf("[xpool] task pool is locked")
)

type taskExecTimeoutCtxKey struct{}

// WithTaskExecTimeoutInContext 在 Submit 调用级别设置任务执行超时。
// 取值为 0 时表示该次提交的任务不启用执行超时。
func WithTaskExecTimeoutInContext(ctx context.Context, timeout time.Duration) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, taskExecTimeoutCtxKey{}, timeout)
}

var _ Task = (*TaskFunc)(nil)

type TaskFunc func(ctx context.Context) error

func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

type taskWrapper struct {
	task    Task
	timeout time.Duration
}

func (t *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]

			slog.Error(
				"[xpool] panic when running task",
				"panic", r,
				"stack", string(buf),
			)

			err = fmt.Errorf("%w: %+v", errTaskRunningPanic, r)
		}
	}()
	return t.task.Run(ctx)
}

// timeoutGroup 超时组。
// 管理任务池中超时的 goroutine id。
type timeoutGroup struct {
	mu sync.RWMutex

	cnt   int32
	idMap map[int32]struct{}
}

func (g *timeoutGroup) in(id int32) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.idMap[id]
	return ok
}

func (g *timeoutGroup) add(id int32) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.idMap[id]; !ok {
		g.idMap[id] = struct{}{}
		g.cnt++
	}
}

func (g *timeoutGroup) del(id int32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.idMap[id]; ok {
		delete(g.idMap, id)
		g.cnt--
	}
}

func (g *timeoutGroup) size() int32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.cnt
}

var _ TaskPool = (*BlockTaskPool)(nil)

// BlockTaskPool 并发阻塞的任务池。
// 任务池会动态控制 goroutine 数量，按需创建 goroutine 执行任务。
//
// 建议使用 WithTaskExecTimeout 传入兜底的任务执行超时时间。
// 避免满请求堆死 worker。
type BlockTaskPool struct {
	// 用于 goroutine 缩容/超时这类决策路径的短临界区。
	// 提交与计数热点使用原子操作，减少整体锁竞争。
	mu sync.Mutex

	id int32 // goroutine id

	maxIdleTime     time.Duration // goroutine 空闲超时时间
	submitTimeout   time.Duration // 提交超时时间
	taskExecTimeout time.Duration // 任务执行超时时间

	state          int32 // 内部状态
	totalG         int32 // goroutine 总数
	totalRunningG  int32 // 正在执行任务的 goroutine 总数
	submitRetryCnt int64
	taskTimeoutCnt int64
	errDropCnt     int64
	stateDropCnt   int64

	// 参数 initG / coreG / maxG 的作用是分层管理 goroutine ( worker )。
	// 三个参数将 goroutine 分为 3 个区间:
	//
	//		[1	  , initG]:	常驻 goroutine ( 不因空闲退出 )。
	//		(initG, coreG]: 核心弹性 goroutine ( 空闲 maxIdleTime 后退出 )。
	//						当 goroutine 处于这个区间，会在退出前 ( maxIdleTime ) 尝试获取任务。
	//						拿到任务则继续执行，没拿到则超时退出。
	//		(coreG, maxG ]: 临时 goroutine
	//						当 goroutine 处于这个区间且当前对立没有可执行任务，则快速退出。
	initG int32 // 初始 goroutine 数量
	coreG int32 // 核心 goroutine 数量
	maxG  int32 // 最大 goroutine 数量

	queue            chan Task // 任务队列
	queueBacklogRate float64   // 任务队列积压率

	timeoutG *timeoutGroup // 超时的 goroutine id 组

	interruptCtx        context.Context
	interruptCancelFunc context.CancelFunc

	errHandler       func(ctx context.Context, err error) // 错误处理器
	errHandleTimeout time.Duration

	errQueueSize int32
	errQueue     chan error

	errWorkerCnt int32
	errWorkers   sync.Once
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

// WithTaskExecTimeout 配置任务执行超时时间。
// 取值为 0 时表示不启用任务执行超时。
func WithTaskExecTimeout(taskExecTimeout time.Duration) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.taskExecTimeout = taskExecTimeout
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

func WithErrQueueSize(errQueueSize int32) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.errQueueSize = errQueueSize
	}
}

func WithErrWorkerCnt(errWorkerCnt int32) option.Opt[BlockTaskPool] {
	return func(p *BlockTaskPool) {
		p.errWorkerCnt = errWorkerCnt
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
		queue: make(chan Task, queueSize),

		initG: initG,
		coreG: initG,
		maxG:  initG,

		maxIdleTime:      defaultMaxIdleTime,
		submitTimeout:    defaultSubmitTimeout,
		taskExecTimeout:  defaultTaskExecTimeout,
		errHandleTimeout: defaultErrHandleTimeout,

		errQueueSize: defaultErrQueueSize,
		errWorkerCnt: defaultErrWorkerCnt,
	}

	ctx := context.Background()
	p.interruptCtx, p.interruptCancelFunc = context.WithCancel(ctx)
	atomic.StoreInt32(&p.state, stateCreated)

	option.Apply(p, opts...)

	// 默认情况 coreG == maxG。
	// 当 coreG == maxG 时，goroutine 的分层会简化为两层:
	//
	//		[1	  , initG]:	常驻 goroutine ( 不因空闲退出 )。
	//		(initG, coreG]: 核心弹性 goroutine ( 空闲 maxIdleTime 后退出 )。
	//
	// 这样使 goroutine 的数量更平滑 ( initG -> maxG 渐进式增长 )，资源更可预测。
	// 所有核心 goroutine 都有相同的生命周期管理。
	if p.coreG != p.initG && p.maxG == p.initG {
		// 只使用 option WithCoreG 的情况，保证 maxG 至少等于 coreG。
		p.maxG = p.coreG
	} else if p.coreG == p.initG && p.maxG != p.initG {
		// 只使用 option WithMaxG 的情况。
		p.coreG = p.maxG
	}

	if !(p.initG <= p.coreG && p.coreG <= p.maxG) {
		return nil, fmt.Errorf("%w: goroutine required to satisfy [ init <= core <= max ]", errInvalidParam)
	}

	p.timeoutG = &timeoutGroup{
		idMap: make(map[int32]struct{}),
	}

	if p.queueBacklogRate < float64(0) || p.queueBacklogRate > float64(1) {
		return nil, fmt.Errorf("%w: queue backlog rate should be in [0, 1]", errInvalidParam)
	}
	if p.maxIdleTime <= 0 {
		return nil, fmt.Errorf("%w: max idle time should be greater than 0", errInvalidParam)
	}
	if p.submitTimeout <= 0 {
		return nil, fmt.Errorf("%w: submit timeout should be greater than 0", errInvalidParam)
	}
	if p.taskExecTimeout < 0 {
		return nil, fmt.Errorf("%w: task exec timeout should be greater or equal to 0", errInvalidParam)
	}
	if p.errHandleTimeout <= 0 {
		return nil, fmt.Errorf("%w: err handle timeout should be greater than 0", errInvalidParam)
	}
	if p.errQueueSize <= 0 {
		return nil, fmt.Errorf("%w: error queue size should be greater than 0", errInvalidParam)
	}
	if p.errWorkerCnt <= 0 {
		return nil, fmt.Errorf("%w: error worker count should be greater than 0", errInvalidParam)
	}

	p.errQueue = make(chan error, p.errQueueSize)
	return p, nil
}

// Submit 提交一个任务。
//
// 在队列已满的情况下，调用者会被阻塞。
// 在 Start 方法被调用后仍然可以调用 Submit 方法。
//
// 注意：
//
//	强烈推荐在调用 Submit 时使用 WithTaskExecTimeout 配置任务执行超时时间。
func (p *BlockTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return errInvalidTask
	}

	// 解析任务执行超时时间。
	taskExecTimeout, err := p.resolveTaskExecTimeout(ctx)
	if err != nil {
		return err
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.submitTimeout)
		defer cancel()
	}

	backoff := defaultSubmitBackoff
	for {
		if atomic.LoadInt32(&p.state) == stateClosing {
			return errPoolIsClosing
		}
		if atomic.LoadInt32(&p.state) == stateClosed {
			return errPoolIsClosed
		}

		tw := &taskWrapper{
			task:    task,
			timeout: taskExecTimeout,
		}

		// 尝试 created 路径 ( p.state == stateCreated )。
		// 此时 pool 还没 Start，允许先把任务提交到队列 ( 只入队不扩容 )。
		// trySubmit 会尝试把 pool state 临时 CAS 成 locked 状态，保证安全写队列。
		var ok bool
		ok, err = p.trySubmit(ctx, tw, stateCreated)
		if ok || err != nil {
			return err
		}

		// 尝试 running 路径 ( p.state == stateRunning )。
		// pool 已运行，提交后可能出发扩容逻辑 ( 创建新 goroutine )。
		ok, err = p.trySubmit(ctx, tw, stateRunning)
		if ok || err != nil {
			return err
		}

		// 队列已满或状态切换中。
		// 短暂退避后重试 ( 默认退避间隔 100us ~ 2ms )，避免队列满/状态竞争时的空转自旋。
		if err := p.waitSubmitRetry(ctx, backoff); err != nil {
			return err
		}
		// 记录提交重试次数，便于观测提交压力。
		atomic.AddInt64(&p.submitRetryCnt, 1)
		if backoff < maxSubmitBackoff {
			backoff *= 2
			if backoff > maxSubmitBackoff {
				backoff = maxSubmitBackoff
			}
		}
	}
}

// waitSubmitRetry 等待提交重试 ( 默认退避间隔 100us ~ 2ms )。
func (p *BlockTaskPool) waitSubmitRetry(ctx context.Context, backoff time.Duration) error {
	timer := time.NewTimer(backoff)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
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

				//nolint:contextcheck // newG 方法在执行任务时候会为每个 task 创建独立的 context 隔离控制。
				go p.newG(id)

				slog.Debug("[xpool] new goroutine created", "id", id)
			}

			// 任务池还未运行 或 当前不允许创建新的 goroutine，直接成功提交。
			return true, nil
		default:
			return false, nil
		}
	}
	return false, nil
}

// allowToCreateG 判断当前是否允许创建新的 goroutine。
// 当满足以下条件时允许创建新的 goroutine：
//
//  1. goroutine 总数小于最大 goroutine 数量。
//  2. 队列存在待运行的 task 且队列积压率达到阈值。
func (p *BlockTaskPool) allowToCreateG() bool {
	// 与 shouldContinue/handleIdleTimeout 共用短临界区。
	// 避免扩缩容判定抖动。
	p.mu.Lock()
	defer p.mu.Unlock()

	// 使用原子快照读取当前 goroutine 数。
	currentTotalG := p.countG()
	if currentTotalG >= p.maxG {
		return false
	}

	if cap(p.queue) == 0 {
		// 无缓冲队列没有积压长度概念，使用 worker 饱和度作为扩容信号。
		return atomic.LoadInt32(&p.totalRunningG) >= currentTotalG
	}

	// 计算队列占用率
	rate := float64(len(p.queue)) / float64(cap(p.queue))

	// 队列存在待运行的 task 且队列积压率达到阈值
	return rate != 0 && rate >= p.queueBacklogRate
}

// newG 创建新的 goroutine。
// 参数 id 用来标识新创建的 goroutine。
func (p *BlockTaskPool) newG(id int32) {
	// 创建一个持续时间为 0 的 timer ( 假超时 )。
	// 该 timer 会立即过期并向其 channel 发送信号。
	//
	// 这里假超时的目的是创建一个可安全引用的 timer 占位。
	// 创建后会立即清空初始到期信号，避免新创建的 goroutine 误出发超时退出。
	// 保证除任务池退出的情况 goroutine 至少执行一个任务。
	//
	// 注意：
	//  1. timer 只保证在等待 x 时间后才发送信号，而不是在 x 时间内发送信号。
	//  2. 不能使用 nil timer，否则会导致 for 循环内的 case <-idleTimer.C 发生 panic。
	idleTimer := time.NewTimer(0)
	if !idleTimer.Stop() {
		// 从 channel 中读取并丢弃该信号，避免假超时导致 goroutine 退出。
		<-idleTimer.C
	}

	for {
		select {
		case <-p.interruptCtx.Done():
			// 收到整个 task pool 的中断信号。
			p.decreaseG(1)
			return
		case <-idleTimer.C:
			// 空闲时收到超时信号 ( 即 goroutine 在 maxIdleTime 时间内没获取到可执行任务 )。
			p.handleIdleTimeout(id)
			return
		case task, ok := <-p.queue:
			// 获取到可执行任务。
			if !p.processTask(id, task, ok, idleTimer) {
				return
			}
		}
	}
}

// handleIdleTimeout 处理空闲超时。
func (p *BlockTaskPool) handleIdleTimeout(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 任务池 goroutine 总数 -1
	p.decreaseG(1)
	// 从超时组移除当前 goroutine id
	p.timeoutG.del(id)
}

// processTask 处理任务的核心流程。
func (p *BlockTaskPool) processTask(id int32, task Task, ok bool, idleTimer *time.Timer) bool {
	// 当 p.queue 被 close 时这里会收到零值 task 和 false。
	if p.timeoutG.in(id) {
		// 当前 goroutine 在超时组中且在超时前成功拿到任务执行。
		p.timeoutG.del(id)
		if !idleTimer.Stop() {
			<-idleTimer.C
		}
	}

	if !ok {
		// !ok 意味着任务队列被关闭。
		p.decreaseG(1)
		// 任务池中没有 goroutine。
		if p.countG() == 0 {
			// 因 shutdown 导致的 goroutine 退出。
			// 最后一个退出的 goroutine 需要负责状态迁移，并通知外部调用者。
			if atomic.CompareAndSwapInt32(&p.state, stateClosing, stateClosed) {
				// 调用 context.CancelFunc 通知外部调用者。
				p.interruptCancelFunc()
			}
		}
		return false
	}

	// 成功获取可执行任务。
	atomic.AddInt32(&p.totalRunningG, 1)

	runCtx, cancel := p.resolveTaskRunContext(task)
	err := task.Run(runCtx)
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		atomic.AddInt64(&p.taskTimeoutCnt, 1)
	}
	cancel()

	atomic.AddInt32(&p.totalRunningG, -1)

	// 处理任务执行错误。
	if err != nil && p.errHandler != nil {
		p.handleTaskError(err)
	}
	return p.shouldContinue(id, idleTimer)
}

func (p *BlockTaskPool) resolveTaskExecTimeout(ctx context.Context) (time.Duration, error) {
	taskExecTimeout := p.taskExecTimeout
	if v := ctx.Value(taskExecTimeoutCtxKey{}); v != nil {
		overrideTaskExecTimeout, ok := v.(time.Duration)
		if !ok {
			return 0, fmt.Errorf("%w: invalid task exec timeout in context", errInvalidParam)
		}
		taskExecTimeout = overrideTaskExecTimeout
	}

	if taskExecTimeout < 0 {
		return 0, fmt.Errorf("%w: task exec timeout should be greater or equal to 0", errInvalidParam)
	}
	return taskExecTimeout, nil
}

func (p *BlockTaskPool) resolveTaskRunContext(task Task) (context.Context, context.CancelFunc) {
	if tw, ok := task.(*taskWrapper); ok && tw.timeout > 0 {
		return context.WithTimeout(p.interruptCtx, tw.timeout)
	}
	return p.interruptCtx, func() {}
}

// handleTaskError 处理任务执行错误。
// 这里只将错误放入错误队列，由专门的错误处理 goroutine 处理。
func (p *BlockTaskPool) handleTaskError(err error) {
	select {
	case p.errQueue <- err:
	default:
		atomic.AddInt64(&p.errDropCnt, 1)
		slog.Warn("[xpool] drop task error: error queue is full")
	}
}

// startErrorWorkers 启动错误处理 goroutine。
func (p *BlockTaskPool) startErrorWorkers() {
	if p.errHandler == nil {
		return
	}

	p.errWorkers.Do(func() {
		for i := int32(0); i < p.errWorkerCnt; i++ {
			go p.errorWorker()
		}
	})
}

// errorWorker 错误处理 goroutine。
func (p *BlockTaskPool) errorWorker() {
	for {
		select {
		case <-p.interruptCtx.Done():
			return
		case err := <-p.errQueue:
			p.runErrorHandler(err)
		}
	}
}

// runErrorHandler 运行错误处理器。
func (p *BlockTaskPool) runErrorHandler(err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]

			slog.Error(
				"[xpool] panic when running error handler",
				"panic", r,
				"source_error", err,
				"stack", string(buf),
			)
		}
	}()

	// 超时控制，避免错误处理长期阻塞。
	ctx, cancel := context.WithTimeout(p.interruptCtx, p.errHandleTimeout)
	defer cancel()
	p.errHandler(ctx, err)
}

// shouldContinue 封装了任务执行后判断是否继续运行的策略（包含缩容和重置超时）。
func (p *BlockTaskPool) shouldContinue(id int32, idleTimer *time.Timer) bool {
	// 任务执行完成后的判断。
	// 这里使用短临界区保证缩容判定稳定，避免并发抖动导致过度扩缩容。
	p.mu.Lock()
	defer p.mu.Unlock()

	currentTotalG := p.countG()

	// 检查队列中是否还有任务需要执行。
	noTaskToExec := len(p.queue) == 0 || int32(len(p.queue)) < currentTotalG
	// 临时 goroutine 的快速退出策略
	if noTaskToExec && p.coreG < currentTotalG && currentTotalG <= p.maxG {
		// 当前 goroutine 处于 (coreG, maxG] 区间（即临时 goroutine），直接退出 goroutine。
		p.decreaseG(1)
		return false
	}

	// 核心 goroutine 的超时管理，
	// 为属于 (initG, coreG] 区间的 goroutine 设置超时器。
	//
	// currentTotalG - p.timeoutG.size() -> 当前活跃的 goroutine 数。
	if p.initG < currentTotalG-p.timeoutG.size() {
		// 核心 goroutine 不立即退出能保证在一定时间 ( maxIdleTime ) 内由任务提交带来的扩容，保持核心处理能力。
		idleTimer.Reset(p.maxIdleTime)
		p.timeoutG.add(id)
	}
	return true
}

// increaseG 增加 goroutine 数量。
func (p *BlockTaskPool) increaseG(delta int32) {
	atomic.AddInt32(&p.totalG, delta)
}

// decreaseG 减少 goroutine 数量。
//
//nolint:unparam // delta 为预留参数，用于后续扩展。
func (p *BlockTaskPool) decreaseG(delta int32) {
	atomic.AddInt32(&p.totalG, -delta)
}

// countG 查看任务池中有多少 goroutine。
func (p *BlockTaskPool) countG() int32 {
	return atomic.LoadInt32(&p.totalG)
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
			p.startErrorWorkers()
			atomic.CompareAndSwapInt32(&p.state, stateLocked, stateRunning)
			return nil
		}
	}
}

// Shutdown 关闭任务池。
//
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
			//	此时工作 goroutine 还可以读取 chan 的剩余数据，直到 chan 的数据被全部取走。
			close(p.queue)
			return p.interruptCtx.Done(), nil
		}
	}
}

// ShutdownNow 立即关闭任务池，并返回剩余未执行的任务 ( 不包含执行中的任务 )。
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
				// Submit 时会包装为 taskWrapper。
				// ShutdownNow 返回原始任务更符合调用方直觉。
				if tw, ok := task.(*taskWrapper); ok {
					tasks = append(tasks, tw.task)
					continue
				}
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
	if interval <= 0 {
		return nil, fmt.Errorf("%w: interval should be greater than 0", errInvalidParam)
	}

	stateChan := make(chan State, defaultStateChanSize)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case timestamp := <-ticker.C:
				p.sendState(stateChan, timestamp.UnixMilli())
			case <-ctx.Done():
				p.sendStateFinal(stateChan, time.Now().UnixMilli())
				close(stateChan)
				return
			case <-p.interruptCtx.Done():
				p.sendStateFinal(stateChan, time.Now().UnixMilli())
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
		// 记录状态采样丢弃次数，便于判断消费端是否滞后。
		atomic.AddInt64(&p.stateDropCnt, 1)
	}
}

func (p *BlockTaskPool) sendStateFinal(ch chan State, timestamp int64) {
	select {
	case ch <- p.getState(timestamp):
		return
	default:
	}

	// 用最终状态替换通道中旧状态，尽量保证调用方收到最后一次采样。
	select {
	case <-ch:
		atomic.AddInt64(&p.stateDropCnt, 1)
	default:
	}
	select {
	case ch <- p.getState(timestamp):
	default:
		atomic.AddInt64(&p.stateDropCnt, 1)
	}
}

func (p *BlockTaskPool) getState(timestamp int64) State {
	return State{
		QueueSize:    int32(cap(p.queue)),
		GoroutineCnt: p.countG(),

		WaitingCnt: int32(len(p.queue)),
		RunningCnt: atomic.LoadInt32(&p.totalRunningG),

		SubmitRetryCnt: atomic.LoadInt64(&p.submitRetryCnt),
		TaskTimeoutCnt: atomic.LoadInt64(&p.taskTimeoutCnt),
		ErrDropCnt:     atomic.LoadInt64(&p.errDropCnt),
		StateDropCnt:   atomic.LoadInt64(&p.stateDropCnt),

		PoolState: atomic.LoadInt32(&p.state),
		Timestamp: timestamp,
	}
}
