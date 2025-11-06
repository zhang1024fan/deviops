package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"dodevops-api/api/task/model"
	"dodevops-api/api/task/service"
	"dodevops-api/common"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// TaskMonitorController 任务监控控制器
type TaskMonitorController struct{}

// NewTaskMonitorController 创建任务监控控制器
func NewTaskMonitorController() *TaskMonitorController {
	return &TaskMonitorController{}
}

// GetQueueMetrics 获取队列指标
// @Summary 获取任务队列指标
// @Description 获取任务队列的运行指标，包括队列长度、处理统计等
// @Tags 任务监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/queue/metrics [get]
func (tc *TaskMonitorController) GetQueueMetrics(c *gin.Context) {
	taskQueue := service.GetTaskQueue()
	if taskQueue == nil {
		result.Failed(c, http.StatusInternalServerError, "任务队列服务未初始化")
		return
	}

	metrics := taskQueue.GetMetrics()
	result.Success(c, metrics)
}

// GetSchedulerStats 获取调度器统计
// @Summary 获取调度器统计信息
// @Description 获取全局调度器的统计信息，包括活跃任务数、下次运行时间等
// @Tags 任务监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/scheduler/stats [get]
func (tc *TaskMonitorController) GetSchedulerStats(c *gin.Context) {
	scheduler := service.GetGlobalScheduler()
	if scheduler == nil {
		result.Failed(c, http.StatusInternalServerError, "全局调度器未初始化")
		return
	}

	stats := scheduler.GetStats()

	// 添加已注册任务ID列表
	entries := scheduler.GetEntries()
	taskIDs := make([]uint, 0, len(entries))
	for taskID := range entries {
		taskIDs = append(taskIDs, taskID)
	}
	stats["registered_task_ids"] = taskIDs

	result.Success(c, stats)
}

// GetSystemStatus 获取系统状态
// @Summary 获取任务系统整体状态
// @Description 获取任务队列和调度器的整体运行状态
// @Tags 任务监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/system/status [get]
func (tc *TaskMonitorController) GetSystemStatus(c *gin.Context) {
	status := make(map[string]interface{})

	// 获取队列状态
	taskQueue := service.GetTaskQueue()
	if taskQueue != nil {
		status["queue"] = map[string]interface{}{
			"running": taskQueue.IsRunning(),
			"metrics": taskQueue.GetMetrics(),
		}
	} else {
		status["queue"] = map[string]interface{}{
			"running": false,
			"error":   "任务队列服务未初始化",
		}
	}

	// 获取调度器状态
	scheduler := service.GetGlobalScheduler()
	if scheduler != nil {
		status["scheduler"] = scheduler.GetStats()
	} else {
		status["scheduler"] = map[string]interface{}{
			"running": false,
			"error":   "全局调度器未初始化",
		}
	}

	// 获取Redis连接状态
	redisClient := common.GetRedisClient()
	if redisClient != nil {
		_, err := redisClient.Ping(redisClient.Context()).Result()
		status["redis"] = map[string]interface{}{
			"connected": err == nil,
			"error":     err,
		}
	} else {
		status["redis"] = map[string]interface{}{
			"connected": false,
			"error":     "Redis客户端未初始化",
		}
	}

	result.Success(c, status)
}

// GetQueueDetails 获取队列详情
// @Summary 获取队列详细信息
// @Description 获取各个优先级队列的详细信息
// @Tags 任务监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/queue/details [get]
func (tc *TaskMonitorController) GetQueueDetails(c *gin.Context) {
	redisClient := common.GetRedisClient()
	if redisClient == nil {
		result.Failed(c, http.StatusInternalServerError, "Redis客户端未初始化")
		return
	}

	details := make(map[string]interface{})

	// 获取各个队列的长度
	queues := map[string]string{
		"high":   service.QueueHigh,
		"normal": service.QueueNormal,
		"low":    service.QueueLow,
		"retry":  service.QueueRetry,
		"failed": service.QueueFailed,
	}

	for name, key := range queues {
		length, err := redisClient.LLen(redisClient.Context(), key).Result()
		if err != nil {
			details[name] = map[string]interface{}{
				"length": 0,
				"error":  err.Error(),
			}
		} else {
			details[name] = map[string]interface{}{
				"length": length,
				"key":    key,
			}
		}
	}

	result.Success(c, details)
}

// ClearFailedQueue 清空失败队列
// @Summary 清空失败队列
// @Description 清空失败任务队列中的所有任务
// @Tags 任务监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/queue/clear-failed [post]
func (tc *TaskMonitorController) ClearFailedQueue(c *gin.Context) {
	redisClient := common.GetRedisClient()
	if redisClient == nil {
		result.Failed(c, http.StatusInternalServerError, "Redis客户端未初始化")
		return
	}

	// 获取失败队列长度
	beforeCount, err := redisClient.LLen(redisClient.Context(), service.QueueFailed).Result()
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取队列长度失败: "+err.Error())
		return
	}

	// 清空失败队列
	err = redisClient.Del(redisClient.Context(), service.QueueFailed).Err()
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "清空失败队列失败: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"cleared_count": beforeCount,
		"queue":         service.QueueFailed,
		"message":       "失败队列已清空",
	}

	result.Success(c, data)
}

// RetryFailedTasks 重试失败任务
// @Summary 重试失败队列中的任务
// @Description 将失败队列中的任务重新提交到正常队列
// @Tags 任务监控
// @Accept json
// @Produce json
// @Param limit query int false "重试任务数量限制，默认10，最大100"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/queue/retry-failed [post]
func (tc *TaskMonitorController) RetryFailedTasks(c *gin.Context) {
	redisClient := common.GetRedisClient()
	if redisClient == nil {
		result.Failed(c, http.StatusInternalServerError, "Redis客户端未初始化")
		return
	}

	// 获取重试数量限制
	limit := 10 // 默认重试10个
	if l := c.Query("limit"); l != "" {
		if _, err := strconv.Atoi(l); err == nil {
			limit, _ = strconv.Atoi(l)
			if limit > 100 {
				limit = 100 // 最大100个
			}
		}
	}

	retryCount := 0
	errors := []string{}

	// 从失败队列中取出任务并重新提交
	for i := 0; i < limit; i++ {
		result, err := redisClient.BRPop(redisClient.Context(), 0, service.QueueFailed).Result()
		if err != nil {
			break // 队列为空或其他错误
		}

		if len(result) != 2 {
			continue
		}

		taskData := result[1]

		// 解析任务消息
		var message service.TaskMessage
		if err := json.Unmarshal([]byte(taskData), &message); err != nil {
			errors = append(errors, "解析任务失败: "+err.Error())
			continue
		}

		// 重置重试次数
		message.RetryCount = 0

		// 重新序列化
		newData, err := json.Marshal(message)
		if err != nil {
			errors = append(errors, "序列化任务失败: "+err.Error())
			continue
		}

		// 根据优先级提交到对应队列
		queueName := service.QueueNormal
		if message.Priority == "high" {
			queueName = service.QueueHigh
		} else if message.Priority == "low" {
			queueName = service.QueueLow
		}

		// 提交到队列
		if err := redisClient.LPush(redisClient.Context(), queueName, newData).Err(); err != nil {
			errors = append(errors, "提交任务失败: "+err.Error())
			continue
		}

		retryCount++
	}

	data := map[string]interface{}{
		"retry_count": retryCount,
		"errors":      errors,
		"message":     "已重试失败任务",
	}

	result.Success(c, data)
}

// PauseScheduledTask 暂停定时任务
// @Summary 暂停定时任务
// @Description 暂停正在运行的定时任务
// @Tags 任务监控
// @Accept json
// @Produce json
// @Param task_id query int true "任务ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/scheduled/pause [post]
func (tc *TaskMonitorController) PauseScheduledTask(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	// 获取全局调度器
	scheduler := service.GetGlobalScheduler()
	if scheduler == nil {
		result.Failed(c, http.StatusInternalServerError, "全局调度器未初始化")
		return
	}

	// 暂停任务
	if err := scheduler.PauseScheduledTask(uint(taskID)); err != nil {
		result.Failed(c, http.StatusInternalServerError, "暂停任务失败: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"task_id": taskID,
		"message": "定时任务已暂停",
		"status":  "paused",
	}

	result.Success(c, data)
}

// ResumeScheduledTask 恢复定时任务
// @Summary 恢复定时任务
// @Description 恢复已暂停的定时任务
// @Tags 任务监控
// @Accept json
// @Produce json
// @Param task_id query int true "任务ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/scheduled/resume [post]
func (tc *TaskMonitorController) ResumeScheduledTask(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	// 获取全局调度器
	scheduler := service.GetGlobalScheduler()
	if scheduler == nil {
		result.Failed(c, http.StatusInternalServerError, "全局调度器未初始化")
		return
	}

	// 恢复任务
	if err := scheduler.ResumeScheduledTask(uint(taskID)); err != nil {
		result.Failed(c, http.StatusInternalServerError, "恢复任务失败: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"task_id": taskID,
		"message": "定时任务已恢复",
		"status":  "running",
	}

	result.Success(c, data)
}

// ResetScheduledTaskStatus 重置定时任务状态
// @Summary 重置定时任务子任务状态
// @Description 将定时任务的所有子任务状态重置为等待中(1)，用于修复状态异常的情况
// @Tags 任务监控
// @Accept json
// @Produce json
// @Param task_id query int true "任务ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/scheduled/reset [post]
func (tc *TaskMonitorController) ResetScheduledTaskStatus(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	// 获取任务信息
	taskService := service.NewTaskService(common.GetDB())
	task, err := taskService.GetTask(uint(taskID))
	if err != nil {
		result.Failed(c, http.StatusNotFound, "任务不存在: "+err.Error())
		return
	}

	// 检查是否为定时任务
	if task.Type != model.TaskTypeScheduled {
		result.Failed(c, http.StatusBadRequest, "只有定时任务可以重置状态")
		return
	}

	// 重置所有子任务状态为等待中
	db := common.GetDB()
	dbResult := db.Exec("UPDATE task_work SET status = 1 WHERE task_id = ? AND status IN (2, 3, 4)", taskID)
	if dbResult.Error != nil {
		result.Failed(c, http.StatusInternalServerError, "重置任务状态失败: "+dbResult.Error.Error())
		return
	}

	data := map[string]interface{}{
		"task_id":       taskID,
		"affected_rows": dbResult.RowsAffected,
		"message":       "定时任务状态已重置",
	}

	result.Success(c, data)
}

// GetTaskStatus 获取任务状态详情
// @Summary 获取任务状态详情
// @Description 获取任务的详细状态信息，包括可执行的操作
// @Tags 任务监控
// @Accept json
// @Produce json
// @Param task_id query int true "任务ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /task/monitor/task/status [get]
func (tc *TaskMonitorController) GetTaskStatus(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	// 获取任务信息
	taskService := service.NewTaskService(common.GetDB())
	task, err := taskService.GetTask(uint(taskID))
	if err != nil {
		result.Failed(c, http.StatusNotFound, "任务不存在: "+err.Error())
		return
	}

	// 检查调度器中的状态
	scheduler := service.GetGlobalScheduler()
	var schedulerStatus string
	if scheduler != nil {
		entries := scheduler.GetEntries()
		if _, exists := entries[uint(taskID)]; exists {
			schedulerStatus = "scheduled"
		} else {
			schedulerStatus = "not_scheduled"
		}
	} else {
		schedulerStatus = "scheduler_unavailable"
	}

	// 构建响应数据
	data := map[string]interface{}{
		"task_id":          task.ID,
		"name":             task.Name,
		"type":             task.Type,
		"type_name":        model.GetTypeName(task.Type),
		"status":           task.Status,
		"status_name":      model.GetStatusName(task.Status),
		"cron_expr":        task.CronExpr,
		"next_run_time":    task.NextRunTime,
		"execute_count":    task.ExecuteCount,
		"created_at":       task.CreatedAt,
		"scheduler_status": schedulerStatus,
		"operations": map[string]bool{
			"can_pause":  task.CanPause(),
			"can_resume": task.CanResume(),
			"can_stop":   task.CanStop(),
		},
	}

	result.Success(c, data)
}