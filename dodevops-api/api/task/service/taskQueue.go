package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"dodevops-api/api/task/model"
	"dodevops-api/common"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/semaphore"
)

// TaskQueue Redis任务队列服务
type TaskQueue struct {
	client       *redis.Client
	ctx          context.Context
	cancel       context.CancelFunc
	workerPool   *WorkerPool
	metrics      *QueueMetrics
	mu           sync.RWMutex
	running      bool
}

// WorkerPool 工作池
type WorkerPool struct {
	workers      int
	semaphore    *semaphore.Weighted
	hostLimits   map[uint]*semaphore.Weighted // 每台主机的并发限制
	globalLimit  *semaphore.Weighted          // 全局并发限制
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	taskQueue    *TaskQueue
	mu           sync.RWMutex
}

// QueueMetrics 队列指标
type QueueMetrics struct {
	EnqueuedTotal    int64 `json:"enqueued_total"`    // 总入队数
	ProcessedTotal   int64 `json:"processed_total"`   // 总处理数
	FailedTotal      int64 `json:"failed_total"`      // 总失败数
	ActiveWorkers    int64 `json:"active_workers"`    // 活跃工作者数
	QueueLength      int64 `json:"queue_length"`      // 队列长度
	LastProcessTime  time.Time `json:"last_process_time"` // 最后处理时间
	mu               sync.RWMutex
}

// TaskMessage 任务消息结构
type TaskMessage struct {
	TaskWork   *model.TaskWork `json:"task_work"`
	Priority   string          `json:"priority"`
	RetryCount int             `json:"retry_count"`
	MaxRetries int             `json:"max_retries"`
	EnqueuedAt time.Time       `json:"enqueued_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// 队列名称常量
const (
	QueueHigh     = "dodevops:task_queue:high"
	QueueNormal   = "dodevops:task_queue:normal"
	QueueLow      = "dodevops:task_queue:low"
	QueueRetry    = "dodevops:task_queue:retry"
	QueueFailed   = "dodevops:task_queue:failed"
	MetricsKey    = "dodevops:task_queue:metrics"
	WorkerKey     = "dodevops:task_queue:workers"
)

// TaskQueueConfig 配置
type TaskQueueConfig struct {
	MaxWorkers        int           `json:"max_workers"`
	MaxConcurrent     int           `json:"max_concurrent"`
	HostMaxConcurrent int           `json:"host_max_concurrent"`
	MaxRetries        int           `json:"max_retries"`
	RetryDelay        time.Duration `json:"retry_delay"`
	ConsumerTimeout   time.Duration `json:"consumer_timeout"`
	EnableMetrics     bool          `json:"enable_metrics"`
}

// DefaultConfig 默认配置
func DefaultConfig() *TaskQueueConfig {
	return &TaskQueueConfig{
		MaxWorkers:        runtime.NumCPU() * 2, // CPU核心数 * 2
		MaxConcurrent:     50,                   // 全局最大并发
		HostMaxConcurrent: 5,                    // 每台主机最大并发
		MaxRetries:        3,                    // 最大重试次数
		RetryDelay:        30 * time.Second,     // 重试延迟
		ConsumerTimeout:   5 * time.Second,      // 消费者超时
		EnableMetrics:     true,                 // 启用指标
	}
}

// NewTaskQueue 创建任务队列
func NewTaskQueue(config *TaskQueueConfig) (*TaskQueue, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := common.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}

	ctx, cancel := context.WithCancel(context.Background())

	queue := &TaskQueue{
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		metrics: &QueueMetrics{LastProcessTime: time.Now()},
		running: false,
	}

	// 创建工作池
	pool, err := NewWorkerPool(config, queue)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("创建工作池失败: %v", err)
	}

	queue.workerPool = pool
	return queue, nil
}

// NewWorkerPool 创建工作池
func NewWorkerPool(config *TaskQueueConfig, taskQueue *TaskQueue) (*WorkerPool, error) {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers:     config.MaxWorkers,
		semaphore:   semaphore.NewWeighted(int64(config.MaxWorkers)),
		hostLimits:  make(map[uint]*semaphore.Weighted),
		globalLimit: semaphore.NewWeighted(int64(config.MaxConcurrent)),
		ctx:         ctx,
		cancel:      cancel,
		taskQueue:   taskQueue,
	}

	return pool, nil
}

// Start 启动任务队列
func (tq *TaskQueue) Start() error {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.running {
		return fmt.Errorf("任务队列已在运行")
	}

	log.Printf("启动任务队列，工作者数量: %d", tq.workerPool.workers)

	// 启动工作池
	if err := tq.workerPool.Start(); err != nil {
		return fmt.Errorf("启动工作池失败: %v", err)
	}

	tq.running = true
	log.Println("任务队列启动成功")
	return nil
}

// Stop 停止任务队列
func (tq *TaskQueue) Stop() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if !tq.running {
		return
	}

	log.Println("停止任务队列...")

	// 停止工作池
	tq.workerPool.Stop()

	// 取消上下文
	tq.cancel()

	tq.running = false
	log.Println("任务队列已停止")
}

// Enqueue 将任务加入队列
func (tq *TaskQueue) Enqueue(taskWork *model.TaskWork, priority string) error {
	return tq.EnqueueWithOptions(taskWork, priority, nil)
}

// EnqueueWithOptions 将任务加入队列（带选项）
func (tq *TaskQueue) EnqueueWithOptions(taskWork *model.TaskWork, priority string, metadata map[string]interface{}) error {
	if !tq.IsRunning() {
		return fmt.Errorf("任务队列未运行")
	}

	message := &TaskMessage{
		TaskWork:   taskWork,
		Priority:   priority,
		RetryCount: 0,
		MaxRetries: 3,
		EnqueuedAt: time.Now(),
		Metadata:   metadata,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %v", err)
	}

	// 根据优先级选择队列
	queueName := tq.getQueueName(priority)

	// 加入队列
	if err := tq.client.LPush(tq.ctx, queueName, data).Err(); err != nil {
		return fmt.Errorf("加入队列失败: %v", err)
	}

	// 更新指标
	tq.metrics.IncEnqueued()

	log.Printf("任务已加入队列: TaskID=%d, Priority=%s, Queue=%s",
		taskWork.TaskID, priority, queueName)

	return nil
}

// getQueueName 根据优先级获取队列名称
func (tq *TaskQueue) getQueueName(priority string) string {
	switch priority {
	case "high":
		return QueueHigh
	case "low":
		return QueueLow
	default:
		return QueueNormal
	}
}

// IsRunning 检查是否运行中
func (tq *TaskQueue) IsRunning() bool {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	return tq.running
}

// GetMetrics 获取队列指标
func (tq *TaskQueue) GetMetrics() *QueueMetrics {
	// 获取队列长度
	tq.updateQueueLength()

	tq.metrics.mu.RLock()
	defer tq.metrics.mu.RUnlock()

	// 复制一份指标返回
	return &QueueMetrics{
		EnqueuedTotal:   tq.metrics.EnqueuedTotal,
		ProcessedTotal:  tq.metrics.ProcessedTotal,
		FailedTotal:     tq.metrics.FailedTotal,
		ActiveWorkers:   tq.metrics.ActiveWorkers,
		QueueLength:     tq.metrics.QueueLength,
		LastProcessTime: tq.metrics.LastProcessTime,
	}
}

// updateQueueLength 更新队列长度
func (tq *TaskQueue) updateQueueLength() {
	totalLength := int64(0)

	queues := []string{QueueHigh, QueueNormal, QueueLow, QueueRetry}
	for _, queue := range queues {
		length, err := tq.client.LLen(tq.ctx, queue).Result()
		if err == nil {
			totalLength += length
		}
	}

	tq.metrics.mu.Lock()
	tq.metrics.QueueLength = totalLength
	tq.metrics.mu.Unlock()
}

// Start 启动工作池
func (wp *WorkerPool) Start() error {
	log.Printf("启动工作池，工作者数量: %d", wp.workers)

	// 启动工作者
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// 启动重试处理器
	wp.wg.Add(1)
	go wp.retryProcessor()

	return nil
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	log.Println("停止工作池...")

	// 取消上下文
	wp.cancel()

	// 等待所有工作者退出
	wp.wg.Wait()

	log.Println("工作池已停止")
}

// worker 工作者
func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	log.Printf("工作者 %d 启动", workerID)
	defer log.Printf("工作者 %d 退出", workerID)

	// 优先级队列列表
	queues := []string{QueueHigh, QueueNormal, QueueLow}

	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
			// 从队列中获取任务
			result, err := wp.taskQueue.client.BRPop(
				wp.ctx,
				wp.taskQueue.workerPool.getTimeout(),
				queues...,
			).Result()

			if err != nil {
				if err == redis.Nil {
					// 队列为空，继续循环
					continue
				}
				if err == context.Canceled {
					return
				}
				log.Printf("工作者 %d 获取任务失败: %v", workerID, err)
				time.Sleep(time.Second)
				continue
			}

			if len(result) != 2 {
				continue
			}

			queueName := result[0]
			taskData := result[1]

			// 解析任务
			var message TaskMessage
			if err := json.Unmarshal([]byte(taskData), &message); err != nil {
				log.Printf("工作者 %d 解析任务失败: %v", workerID, err)
				continue
			}

			// 处理任务
			wp.processTask(workerID, queueName, &message)
		}
	}
}

// processTask 处理任务
func (wp *WorkerPool) processTask(workerID int, _ string, message *TaskMessage) {
	taskWork := message.TaskWork

	log.Printf("工作者 %d 开始处理任务: TaskID=%d, TemplateID=%d",
		workerID, taskWork.TaskID, taskWork.TemplateID)

	// 获取资源锁
	if err := wp.acquireResources(taskWork.HostID); err != nil {
		log.Printf("工作者 %d 获取资源失败: %v", workerID, err)
		wp.retryTask(message)
		return
	}

	// 释放资源锁
	defer wp.releaseResources(taskWork.HostID)

	// 更新活跃工作者数
	wp.taskQueue.metrics.IncActiveWorkers()
	defer wp.taskQueue.metrics.DecActiveWorkers()

	// 执行任务
	taskService := NewTaskWorkService()
	if err := taskService.(*TaskWorkServiceImpl).executeJob(taskWork); err != nil {
		log.Printf("工作者 %d 执行任务失败: %v", workerID, err)

		// 更新失败指标
		wp.taskQueue.metrics.IncFailed()

		// 重试逻辑
		wp.retryTask(message)
		return
	}

	// 更新成功指标
	wp.taskQueue.metrics.IncProcessed()
	wp.taskQueue.metrics.UpdateLastProcessTime()

	log.Printf("工作者 %d 完成任务: TaskID=%d", workerID, taskWork.TaskID)
}

// acquireResources 获取资源锁
func (wp *WorkerPool) acquireResources(hostID uint) error {
	// 获取全局资源锁
	if err := wp.globalLimit.Acquire(wp.ctx, 1); err != nil {
		return fmt.Errorf("获取全局资源锁失败: %v", err)
	}

	// 获取主机资源锁
	wp.mu.Lock()
	hostLimit, exists := wp.hostLimits[hostID]
	if !exists {
		hostLimit = semaphore.NewWeighted(5) // 每台主机最多5个并发
		wp.hostLimits[hostID] = hostLimit
	}
	wp.mu.Unlock()

	if err := hostLimit.Acquire(wp.ctx, 1); err != nil {
		wp.globalLimit.Release(1)
		return fmt.Errorf("获取主机资源锁失败: %v", err)
	}

	return nil
}

// releaseResources 释放资源锁
func (wp *WorkerPool) releaseResources(hostID uint) {
	wp.mu.RLock()
	if hostLimit, exists := wp.hostLimits[hostID]; exists {
		hostLimit.Release(1)
	}
	wp.mu.RUnlock()

	wp.globalLimit.Release(1)
}

// retryTask 重试任务
func (wp *WorkerPool) retryTask(message *TaskMessage) {
	message.RetryCount++

	if message.RetryCount >= message.MaxRetries {
		// 超过最大重试次数，移到失败队列
		wp.moveToFailedQueue(message)
		return
	}

	// 延迟重试
	message.Metadata = map[string]interface{}{
		"retry_at": time.Now().Add(30 * time.Second).Unix(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化重试任务失败: %v", err)
		return
	}

	// 加入重试队列
	if err := wp.taskQueue.client.LPush(wp.ctx, QueueRetry, data).Err(); err != nil {
		log.Printf("加入重试队列失败: %v", err)
	}
}

// moveToFailedQueue 移动到失败队列
func (wp *WorkerPool) moveToFailedQueue(message *TaskMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化失败任务失败: %v", err)
		return
	}

	if err := wp.taskQueue.client.LPush(wp.ctx, QueueFailed, data).Err(); err != nil {
		log.Printf("加入失败队列失败: %v", err)
	}

	log.Printf("任务已移至失败队列: TaskID=%d, RetryCount=%d",
		message.TaskWork.TaskID, message.RetryCount)
}

// retryProcessor 重试处理器
func (wp *WorkerPool) retryProcessor() {
	defer wp.wg.Done()

	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-ticker.C:
			wp.processRetryQueue()
		}
	}
}

// processRetryQueue 处理重试队列
func (wp *WorkerPool) processRetryQueue() {
	for {
		// 从重试队列获取任务
		result, err := wp.taskQueue.client.BRPop(wp.ctx, time.Second, QueueRetry).Result()
		if err != nil {
			if err == redis.Nil {
				// 队列为空
				break
			}
			continue
		}

		if len(result) != 2 {
			continue
		}

		var message TaskMessage
		if err := json.Unmarshal([]byte(result[1]), &message); err != nil {
			continue
		}

		// 检查是否到了重试时间
		if retryAt, exists := message.Metadata["retry_at"]; exists {
			if retryTime, ok := retryAt.(float64); ok {
				if time.Now().Unix() < int64(retryTime) {
					// 还没到重试时间，重新放回队列
					data, _ := json.Marshal(message)
					wp.taskQueue.client.LPush(wp.ctx, QueueRetry, data)
					continue
				}
			}
		}

		// 重新加入正常队列
		queueName := wp.taskQueue.getQueueName(message.Priority)
		data, _ := json.Marshal(message)
		wp.taskQueue.client.LPush(wp.ctx, queueName, data)
	}
}

// getTimeout 获取超时时间
func (wp *WorkerPool) getTimeout() time.Duration {
	return 5 * time.Second
}

// 指标方法
func (m *QueueMetrics) IncEnqueued() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EnqueuedTotal++
}

func (m *QueueMetrics) IncProcessed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ProcessedTotal++
}

func (m *QueueMetrics) IncFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FailedTotal++
}

func (m *QueueMetrics) IncActiveWorkers() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveWorkers++
}

func (m *QueueMetrics) DecActiveWorkers() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ActiveWorkers > 0 {
		m.ActiveWorkers--
	}
}

func (m *QueueMetrics) UpdateLastProcessTime() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LastProcessTime = time.Now()
}

// 全局任务队列实例
var globalTaskQueue *TaskQueue

// InitTaskQueue 初始化全局任务队列
func InitTaskQueue(config *TaskQueueConfig) error {
	var err error
	globalTaskQueue, err = NewTaskQueue(config)
	if err != nil {
		return err
	}

	return globalTaskQueue.Start()
}

// GetTaskQueue 获取全局任务队列
func GetTaskQueue() *TaskQueue {
	return globalTaskQueue
}

// StopTaskQueue 停止全局任务队列
func StopTaskQueue() {
	if globalTaskQueue != nil {
		globalTaskQueue.Stop()
	}
}