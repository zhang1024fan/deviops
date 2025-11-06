package controller

import (
	"net/http"
	"strconv"
	cmdbService "dodevops-api/api/cmdb/service"
	"dodevops-api/api/task/model"
	"dodevops-api/api/task/service"
	"dodevops-api/api/task/dao"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

var (
	hostSSHService = func() cmdbService.CmdbHostSSHServiceInterface { return nil }
	taskTemplateDao = func() dao.TaskTemplateDaoInterface { return nil }
	taskTemplateService *service.TaskService
)

func initServices() {
	if taskTemplateService == nil {
		hostSSHSvc := cmdbService.GetCmdbHostSSHService()
		taskTemplateDAO := dao.NewTaskTemplateDao(common.GetDB())
		taskTemplateService = service.NewTaskServiceImpl(hostSSHSvc, taskTemplateDAO)
	}
}

// CreateTemplate 创建任务模板
// @Summary 创建任务模板
// @Description 创建新的任务模板
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param body body model.TaskTemplate true "任务模板信息"
// @Success 200 {object} result.Result
// @Router /api/v1/template/add [post]
// @Security ApiKeyAuth
func CreateTemplate(ctx *gin.Context) {
	initServices()
	var template model.TaskTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	taskTemplateService.CreateTemplate(ctx, &template)
}

// GetAllTemplates 获取所有模板
// @Summary 获取所有模板
// @Description 获取所有任务模板列表
// @Tags 任务中心
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/template/list [get]
// @Security ApiKeyAuth
func GetAllTemplates(ctx *gin.Context) {
	initServices()
	taskTemplateService.GetAllTemplates(ctx)
}

// UpdateTemplate 更新模板
// @Summary 更新模板
// @Description 更新任务模板
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id query uint true "模板ID"
// @Param body body model.TaskTemplate true "需要更新的模板字段"
// @Success 200 {object} result.Result
// @Router /api/v1/template/update [put]
// @Security ApiKeyAuth
func UpdateTemplate(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少模板ID参数")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var template model.TaskTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}
	template.ID = uint(id)

	taskTemplateService.UpdateTemplate(ctx, &template)
}

// DeleteTemplate 删除模板
// @Summary 删除模板
// @Description 删除指定ID的任务模板
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id query uint true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/template/delete [delete]
// @Security ApiKeyAuth
func DeleteTemplate(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少模板ID参数")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的模板ID")
		return
	}
	taskTemplateService.DeleteTemplate(ctx, uint(id))
}

// GetTemplateByID 根据ID获取模板
// @Summary 根据ID获取模板
// @Description 获取指定ID的任务模板
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param id path uint true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/template/info/{id} [get]
// @Security ApiKeyAuth
func GetTemplateByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少模板ID参数")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的模板ID")
		return
	}
	taskTemplateService.GetTemplateByID(ctx, uint(id))
}

// GetTemplateContent 根据ID获取脚本内容
// @Summary 获取脚本内容
// @Description 获取指定ID的任务模板脚本内容
// @Tags 任务中心
// @Accept json
// @Produce text/plain
// @Param id path uint true "模板ID"
// @Success 200 {string} string "脚本内容"
// @Router /api/v1/template/content/{id} [get]
// @Security ApiKeyAuth
func GetTemplateContent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.String(http.StatusBadRequest, "缺少模板ID参数")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, "无效的模板ID")
		return
	}
	taskTemplateService.GetTemplateContent(ctx, uint(id))
}

// GetTemplatesByName 根据名称模糊查询模板
// @Summary 根据名称模糊查询模板
// @Description 获取名称包含指定字符串的任务模板列表
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param name query string true "模板名称"
// @Success 200 {object} result.Result
// @Router /api/v1/template/query/name [get]
// @Security ApiKeyAuth
func GetTemplatesByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少名称参数")
		return
	}
	taskTemplateService.GetTemplatesByName(ctx, name)
}

// GetTemplatesByType 根据类型查询模板
// @Summary 根据类型查询模板
// @Description 获取指定类型的任务模板列表
// @Tags 任务中心
// @Accept json
// @Produce json
// @Param type query int true "模板类型(1=shell, 2=python, 3=ansible)"
// @Success 200 {object} result.Result
// @Router /api/v1/template/query/type [get]
// @Security ApiKeyAuth
func GetTemplatesByType(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	if typeStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少类型参数")
		return
	}
	templateType, err := strconv.Atoi(typeStr)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的模板类型")
		return
	}
	taskTemplateService.GetTemplatesByType(ctx, templateType)
}
