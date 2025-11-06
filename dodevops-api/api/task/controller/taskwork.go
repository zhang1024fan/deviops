package controller

import (
	"net/http"
	"strconv"
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/service"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type ITaskWorkController interface {
	StartJob(c *gin.Context)                  // 启动指定任务
	GetJobLog(c *gin.Context)                 // 获取任务日志
	StopJob(c *gin.Context)                   // 停止单个任务
	GetJobStatus(c *gin.Context)              // 获取任务状态
}

type TaskWorkController struct {
	service service.ITaskWorkService
}

func TaskWork() ITaskWorkController {
	return &TaskWorkController{
		service: service.NewTaskWorkService(),
	}
}

// StartJob 启动任务
// @Summary 启动任务
// @Description 根据任务ID启动任务
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param taskId query int true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/taskjob/start [post]
// @Security ApiKeyAuth
func (t *TaskWorkController) StartJob(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	if taskIdStr == "" {
		result.Failed(c, http.StatusBadRequest, "缺少任务ID参数")
		return
	}

	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	// 获取任务详情
	task, err := service.NewTaskService(common.GetDB()).GetTask(uint(taskId))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取任务失败: "+err.Error())
		return
	}

	// 根据任务类型处理
	if task.Type == 2 { // 定时任务
		// 获取所有关联的task_work记录
		taskWorks, err := dao.TaskWorkDao().GetByTaskID(uint(taskId))
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取任务工作项失败: "+err.Error())
			return
		}

		// 为每个task_work记录调度定时任务
		for _, work := range taskWorks {
			if work.Status == 1 { // 只调度等待中的任务
				if err := t.service.ScheduleJob(&work); err != nil {
					result.Failed(c, http.StatusInternalServerError, "调度定时任务失败: "+err.Error())
					return
				}
			}
		}
	} else { // 普通任务
		err = t.service.StartJob(uint(taskId))
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	result.Success(c, gin.H{
		"message": "任务已开始" + map[bool]string{true: "调度", false: "执行"}[task.Type == 2],
	})
}

// GetJobLog 获取任务日志
// @Summary 获取任务日志
// @Description 根据任务ID和模板ID获取日志
// @Tags 任务作业
// @Param taskId query int true "任务ID"
// @Param templateId query int true "模板ID"
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/taskjob/log [get]
// @Security ApiKeyAuth
func (t *TaskWorkController) GetJobLog(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	templateIdStr := c.Query("templateId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	templateId, err := strconv.ParseUint(templateIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "模板ID格式错误")
		return
	}

	logContent, err := t.service.GetJobLog(uint(taskId), uint(templateId))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, logContent)
}

// StopJob 停止单个任务
// @Summary 停止单个任务
// @Description 根据任务ID和模板ID停止任务
// @Tags 任务作业
// @Param taskId query int true "任务ID"
// @Param templateId query int true "模板ID"
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/taskjob/stop [post]
// @Security ApiKeyAuth
func (t *TaskWorkController) StopJob(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	templateIdStr := c.Query("templateId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	templateId, err := strconv.ParseUint(templateIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "模板ID格式错误")
		return
	}

	err = t.service.StopJob(uint(taskId), uint(templateId))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, gin.H{
		"message": "任务已停止",
	})
}

// GetJobStatus 获取任务状态
// @Summary 获取任务状态
// @Description 根据任务ID和模板ID获取任务状态
// @Tags 任务作业
// @Param taskId query int true "任务ID"
// @Param templateId query int true "模板ID"
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/taskjob/status [get]
// @Security ApiKeyAuth
func (t *TaskWorkController) GetJobStatus(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	templateIdStr := c.Query("templateId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "任务ID格式错误")
		return
	}

	templateId, err := strconv.ParseUint(templateIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "模板ID格式错误")
		return
	}

	status, err := t.service.GetJobStatus(uint(taskId), uint(templateId))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, status)
}
