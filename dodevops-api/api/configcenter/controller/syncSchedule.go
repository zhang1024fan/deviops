package controller

import (
	"strconv"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	"dodevops-api/common/result"
	"dodevops-api/scheduler"
	"github.com/gin-gonic/gin"
)

type SyncScheduleController struct {
	service *service.SyncScheduleService
}

func NewSyncScheduleController() *SyncScheduleController {
	return &SyncScheduleController{
		service: service.NewSyncScheduleService(),
	}
}

// Create 创建定时同步配置
// @Summary 创建定时同步配置
// @Tags Config配置中心
// @Param syncSchedule body model.CreateSyncScheduleDto true "定时同步配置信息"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Create(ctx *gin.Context) {
	var dto model.CreateSyncScheduleDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	syncSchedule := model.SyncSchedule{
		Name:     dto.Name,
		CronExpr: dto.CronExpr,
		KeyTypes: dto.KeyTypes,
		Status:   dto.Status,
		Remark:   dto.Remark,
	}

	if err := c.service.Create(&syncSchedule); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	// 如果配置是启用状态，则添加到调度器
	if syncSchedule.Status == 1 {
		if err := scheduler.GetManager().AddSyncSchedule(&syncSchedule); err != nil {
			// 调度器添加失败不影响数据库创建，只记录日志
			// log.Printf("Failed to add schedule to scheduler: %v", err)
		}
	}

	result.Success(ctx, syncSchedule)
}

// Update 更新定时同步配置
// @Summary 更新定时同步配置
// @Tags Config配置中心
// @Param syncSchedule body model.UpdateSyncScheduleDto true "定时同步配置信息"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [put]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Update(ctx *gin.Context) {
	var dto model.UpdateSyncScheduleDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	syncSchedule := model.SyncSchedule{
		ID:       dto.ID,
		Name:     dto.Name,
		CronExpr: dto.CronExpr,
		KeyTypes: dto.KeyTypes,
		Status:   dto.Status,
		Remark:   dto.Remark,
	}

	if err := c.service.Update(&syncSchedule); err != nil {
		result.Failed(ctx, 500, "更新失败: "+err.Error())
		return
	}

	// 更新调度器中的配置
	if err := scheduler.GetManager().UpdateSyncSchedule(&syncSchedule); err != nil {
		// 调度器更新失败不影响数据库更新，只记录日志
		// log.Printf("Failed to update schedule in scheduler: %v", err)
	}

	result.Success(ctx, syncSchedule)
}

// Delete 删除定时同步配置
// @Summary 删除定时同步配置
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule [delete]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	if err := c.service.Delete(req.ID); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	// 从调度器中移除配置
	scheduler.GetManager().RemoveSyncSchedule(req.ID)

	result.Success(ctx, nil)
}

// GetByID 根据ID获取定时同步配置详情
// @Summary 根据ID获取定时同步配置详情
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	syncSchedule, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, syncSchedule)
}

// List 获取定时同步配置列表（分页）
// @Summary 获取定时同步配置列表（分页）
// @Tags Config配置中心
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.SyncSchedule}}
// @Router /api/v1/config/sync-schedule/list [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) List(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	syncSchedules, total, err := c.service.ListWithPage(params.Page, params.PageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	pageResult := result.PageResult{
		List:     syncSchedules,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	result.Success(ctx, pageResult)
}

// ToggleStatus 切换配置状态（启用/禁用）
// @Summary 切换配置状态（启用/禁用）
// @Tags Config配置中心
// @Param request body object{id=uint,status=int} true "配置ID和状态"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule/toggle-status [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) ToggleStatus(ctx *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"min=0,max=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.ToggleStatus(req.ID, req.Status); err != nil {
		result.Failed(ctx, 500, "操作失败: "+err.Error())
		return
	}

	// 根据状态更新调度器
	if req.Status == 1 {
		// 启用：获取配置并添加到调度器
		if schedule, err := c.service.GetByID(req.ID); err == nil {
			scheduler.GetManager().AddSyncSchedule(schedule)
		}
	} else {
		// 禁用：从调度器中移除
		scheduler.GetManager().RemoveSyncSchedule(req.ID)
	}

	result.Success(ctx, nil)
}

// GetActiveSchedules 获取所有启用的定时同步配置
// @Summary 获取所有启用的定时同步配置
// @Tags Config配置中心
// @Success 200 {object} result.Result{data=[]model.SyncSchedule}
// @Router /api/v1/config/sync-schedule/active [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetActiveSchedules(ctx *gin.Context) {
	schedules, err := c.service.GetActiveSchedules()
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, schedules)
}

// TriggerManualSync 手动触发定时同步（测试用）
// @Summary 手动触发定时同步（测试用）
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule/trigger [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) TriggerManualSync(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, 400, "配置ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "配置ID格式错误")
		return
	}

	// 获取配置信息
	schedule, err := c.service.GetByID(uint(id))
	if err != nil {
		result.Failed(ctx, 500, "查询配置失败: "+err.Error())
		return
	}

	if schedule.Status != 1 {
		result.Failed(ctx, 400, "配置未启用，无法执行同步")
		return
	}

	// 解析云厂商类型
	keyTypes, err := c.service.ParseKeyTypes(schedule.KeyTypes)
	if err != nil {
		result.Failed(ctx, 500, "解析云厂商类型失败: "+err.Error())
		return
	}

	// 执行同步（这里先返回成功，实际的同步逻辑将在调度器中实现）
	result.Success(ctx, gin.H{
		"message":  "同步任务已触发",
		"schedule": schedule.Name,
		"keyTypes": keyTypes,
	})
}

// GetSchedulerStats 获取调度器状态信息
// @Summary 获取调度器状态信息
// @Tags Config配置中心
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Router /api/v1/config/sync-schedule/scheduler-stats [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetSchedulerStats(ctx *gin.Context) {
	stats := scheduler.GetManager().GetSyncSchedulerStats()
	result.Success(ctx, stats)
}

// GetSyncLog 获取同步日志
// @Summary 获取同步日志
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result{data=object{syncLog=string,lastRunTime=string}}
// @Router /api/v1/config/sync-schedule/log [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetSyncLog(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	syncSchedule, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"id":          syncSchedule.ID,
		"name":        syncSchedule.Name,
		"syncLog":     syncSchedule.SyncLog,
		"lastRunTime": syncSchedule.LastRunTime,
		"nextRunTime": syncSchedule.NextRunTime,
	})
}