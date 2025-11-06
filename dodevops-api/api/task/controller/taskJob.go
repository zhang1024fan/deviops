package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"dodevops-api/api/task/model"
	"dodevops-api/api/task/service"
	"dodevops-api/common"
	"dodevops-api/common/result"
)

var (
	taskJobService service.TaskJobService
)

// getTaskJobService 获取任务服务（延迟初始化）
func getTaskJobService() service.TaskJobService {
	if taskJobService == nil {
		taskJobService = service.NewTaskService(common.GetDB())
	}
	return taskJobService
}

// CreateTaskRequest 创建任务请求体
type CreateTaskRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	Shell    string `json:"shell" binding:"required"`
	HostIDs  string `json:"host_ids" binding:"required"`
	CronExpr string `json:"cron_expr,omitempty"`
	Remark   string `json:"remark,omitempty"`
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建新的任务
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param body body CreateTaskRequest true "任务信息"
// @Success 200 {object} result.Result
// @Router /api/v1/task/add [post]
// @Security ApiKeyAuth
func CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	// 验证主机ID格式(必须是逗号分隔的数字)
	if _, err := strconv.Atoi(strings.ReplaceAll(req.HostIDs, ",", "")); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "主机ID必须是数字或逗号分隔的数字")
		return
	}

	// 验证任务名称是否已存在
	if exists, err := getTaskJobService().TaskNameExists(req.Name); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "检查任务名称失败: "+err.Error())
		return
	} else if exists {
		result.Failed(ctx, http.StatusBadRequest, "任务名称已存在")
		return
	}

	// 验证定时任务必须包含CronExpr
	if req.Type == 2 && req.CronExpr == "" {
		result.Failed(ctx, http.StatusBadRequest, "定时任务必须包含定时表达式")
		return
	}

	// 计算任务数量(根据Shell中逗号分隔的模板ID数量)
	taskCount := 1
	if strings.Contains(req.Shell, ",") {
		taskCount = len(strings.Split(req.Shell, ","))
	}

	task := model.Task{
		Name:         req.Name,
		Type:         req.Type,
		Shell:        req.Shell,
		HostIDs:      req.HostIDs,
		CronExpr:     req.CronExpr,
		Remark:       req.Remark,
		Status:       1, // 默认等待中
		TaskCount:    taskCount,
		ExecuteCount: 0, // 初始执行次数为0
	}

	// 如果是定时任务，计算下次执行时间
	if req.Type == 2 && req.CronExpr != "" {
		if nextTime, err := calculateNextRunTime(req.CronExpr); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "定时表达式无效: "+err.Error())
			return
		} else {
			task.NextRunTime = nextTime
		}
	}

	err := getTaskJobService().CreateTask(&task)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "创建任务失败: "+err.Error())
		return
	}

	// 重新获取任务以确保数据一致性
	createdTask, err := getTaskJobService().GetTask(task.ID)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取创建的任务失败: "+err.Error())
		return
	}

	result.Success(ctx, createdTask)
}


// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param body body model.Task true "需要更新的任务字段(必须包含ID)"
// @Success 200 {object} result.Result
// @Router /api/v1/task/update [put]
// @Security ApiKeyAuth
func UpdateTask(ctx *gin.Context) {
	var task model.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	if task.ID == 0 {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务ID")
		return
	}

	// 验证定时任务必须包含CronExpr
	if task.Type == 2 && task.CronExpr == "" {
		result.Failed(ctx, http.StatusBadRequest, "定时任务必须包含定时表达式")
		return
	}

	if err := getTaskJobService().UpdateTask(task.ID, &task); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "更新任务失败: "+err.Error())
		return
	}

	result.Success(ctx, task)
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除任务
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param body body model.TaskIDRequest true "任务ID请求"
// @Success 200 {object} result.Result
// @Router /api/v1/task/delete [delete]
// @Security ApiKeyAuth
func DeleteTask(ctx *gin.Context) {
	var req model.TaskIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	if err := getTaskJobService().DeleteTask(req.ID); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "删除任务失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// GetTaskByID 根据ID查询任务
// @Summary 根据ID查询任务
// @Description 根据任务ID查询任务详情
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id query uint true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/task/get [get]
// @Security ApiKeyAuth
func GetTaskByID(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少任务ID参数")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	task, err := getTaskJobService().GetTask(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusNotFound, "任务不存在")
		return
	}

	result.Success(ctx, task)
}

// ListTasks 获取任务列表
// @Summary 获取任务列表
// @Description 获取任务列表，支持分页和条件查询
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param name query string false "任务名称"
// @Param status query int false "任务状态(1=等待中,2=运行中,3=成功,4=异常)"
// @Success 200 {object} result.Result
// @Router /api/v1/task/list [get]
// @Security ApiKeyAuth
func ListTasks(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	name := ctx.Query("name")
	statusStr := ctx.Query("status")

	var status int
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	tasks, total, err := getTaskJobService().ListTasksWithDetails(page, pageSize, name, status)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取任务列表失败: "+err.Error())
		return
	}

	result.SuccessWithPage(ctx, tasks, total, page, pageSize)
}

// ListTasksWithDetails 获取任务列表（包含关联信息）
// @Summary 获取任务列表（包含关联信息）
// @Description 获取任务列表，支持分页和条件查询，包含模板和主机的详细信息
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param name query string false "任务名称"
// @Param status query int false "任务状态(1=等待中,2=运行中,3=成功,4=异常)"
// @Success 200 {object} result.Result
// @Router /api/v1/task/list-with-details [get]
// @Security ApiKeyAuth
func ListTasksWithDetails(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	name := ctx.Query("name")
	statusStr := ctx.Query("status")

	var status int
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	tasks, total, err := getTaskJobService().ListTasksWithDetails(page, pageSize, name, status)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取任务列表失败: "+err.Error())
		return
	}

	result.SuccessWithPage(ctx, tasks, total, page, pageSize)
}

// GetTasksByName 根据名称查询任务
// @Summary 根据名称查询任务
// @Description 根据任务名称模糊查询任务
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param name query string true "任务名称"
// @Success 200 {object} result.Result
// @Router /api/v1/task/query/name [get]
// @Security ApiKeyAuth
func GetTasksByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务名称")
		return
	}

	tasks, err := getTaskJobService().GetTasksByName(name)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询任务失败: "+err.Error())
		return
	}

	result.Success(ctx, tasks)
}

// GetTasksByType 根据类型查询任务
// @Summary 根据类型查询任务
// @Description 根据任务类型查询任务
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param type query int true "任务类型(1=普通任务,2=定时任务,3=ansible任务,4=工作作业)"
// @Success 200 {object} result.Result
// @Router /api/v1/task/query/type [get]
// @Security ApiKeyAuth
func GetTasksByType(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	if typeStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务类型")
		return
	}

	taskType, err := strconv.Atoi(typeStr)
	if err != nil || taskType < 1 || taskType > 4 {
		result.Failed(ctx, http.StatusBadRequest, "任务类型必须是1-4之间的数字")
		return
	}

	tasks, err := getTaskJobService().GetTasksByType(taskType)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询任务失败: "+err.Error())
		return
	}

	result.Success(ctx, tasks)
}

// GetTasksByStatus 根据状态查询任务
// @Summary 根据状态查询任务
// @Description 根据任务状态查询任务
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param status query int true "任务状态(1=等待中,2=运行中,3=成功,4=异常)"
// @Success 200 {object} result.Result
// @Router /api/v1/task/query/status [get]
// @Security ApiKeyAuth
func GetTasksByStatus(ctx *gin.Context) {
	statusStr := ctx.Query("status")
	if statusStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务状态")
		return
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务状态")
		return
	}

	tasks, err := getTaskJobService().GetTasksByStatus(status)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询任务失败: "+err.Error())
		return
	}

	result.Success(ctx, tasks)
}

// GetNextExecutionTime 计算下次执行时间
// @Summary 计算下次执行时间
// @Description 根据cron表达式计算下次执行时间
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param cron query string true "cron表达式"
// @Success 200 {object} result.Result
// @Router /api/v1/task/next-execution [get]
// @Security ApiKeyAuth
func GetNextExecutionTime(ctx *gin.Context) {
	cronExpr := ctx.Query("cron")
	if cronExpr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供cron表达式")
		return
	}

	// 首先验证cron表达式格式
	if _, err := cron.ParseStandard(cronExpr); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的cron表达式格式: "+err.Error())
		return
	}

	nextTime, err := getTaskJobService().GetNextExecutionTime(cronExpr)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的cron表达式: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"next_execution_time": nextTime.Format("2006-01-02 15:04:05"),
	})
}

// GetTaskTemplatesWithStatus 获取任务模板及状态
// @Summary 获取任务模板及状态
// @Description 根据任务ID获取关联模板信息及状态
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id query uint true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/task/templates [get]
// @Security ApiKeyAuth
func GetTaskTemplatesWithStatus(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务ID")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	templates, err := getTaskJobService().GetTaskTemplatesWithStatus(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询模板信息失败: "+err.Error())
		return
	}

	result.Success(ctx, templates)
}

// GetTaskExecutionInfo 获取任务执行信息（包含执行次数和下次执行时间）
// @Summary 获取任务执行信息
// @Description 获取任务的执行次数和下次执行时间
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id query uint true "任务ID"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Router /api/v1/task/execution-info [get]
// @Security ApiKeyAuth
func GetTaskExecutionInfo(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务ID")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	task, err := getTaskJobService().GetTask(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusNotFound, "任务不存在")
		return
	}

	// 构造返回信息
	executionInfo := map[string]interface{}{
		"id":            task.ID,
		"name":          task.Name,
		"type":          task.Type,
		"execute_count": task.ExecuteCount,
		"next_run_time": task.NextRunTime,
		"cron_expr":     task.CronExpr,
		"status":        task.Status,
		"last_end_time": task.EndTime,
	}

	// 如果是定时任务且有cron表达式，计算下次执行时间
	if task.Type == 2 && task.CronExpr != "" {
		if nextTime, err := calculateNextRunTime(task.CronExpr); err == nil {
			executionInfo["calculated_next_run_time"] = nextTime
		}
	}

	result.Success(ctx, executionInfo)
}

// calculateNextRunTime 计算下次执行时间
func calculateNextRunTime(cronExpr string) (*time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, err
	}

	nextTime := schedule.Next(time.Now())
	return &nextTime, nil
}

// IncrementExecuteCount 增加执行次数并更新下次执行时间
func IncrementExecuteCount(taskID uint) error {
	service := getTaskJobService()

	// 获取任务信息
	task, err := service.GetTask(taskID)
	if err != nil {
		return err
	}

	// 增加执行次数
	task.ExecuteCount++

	// 如果是定时任务，更新下次执行时间
	if task.Type == 2 && task.CronExpr != "" {
		if nextTime, err := calculateNextRunTime(task.CronExpr); err == nil {
			task.NextRunTime = nextTime
		}
	}

	// 更新任务
	return service.UpdateTask(task.ID, task)
}

// UpdateTaskAfterExecution 任务执行完成后更新相关字段
func UpdateTaskAfterExecution(taskID uint, status int, duration int, endTime time.Time, tasklog string) error {
	service := getTaskJobService()

	// 获取任务信息
	task, err := service.GetTask(taskID)
	if err != nil {
		return err
	}

	// 增加执行次数
	task.ExecuteCount++
	task.Status = status
	task.Duration = duration
	task.EndTime = &endTime
	task.Tasklog = tasklog

	// 如果是定时任务，更新下次执行时间
	if task.Type == 2 && task.CronExpr != "" {
		if nextTime, err := calculateNextRunTime(task.CronExpr); err == nil {
			task.NextRunTime = nextTime
		}
	}

	// 更新任务
	return service.UpdateTask(task.ID, task)
}
