package controller

import (
	"strconv"

	"dodevops-api/api/app/model"
	"dodevops-api/api/app/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplicationController 应用控制器
type ApplicationController struct {
	appService service.IApplicationService
}

// NewApplicationController 创建应用控制器
func NewApplicationController(db *gorm.DB) *ApplicationController {
	return &ApplicationController{
		appService: service.NewApplicationService(db),
	}
}

// CreateApplication 创建应用
// @Summary 创建应用
// @Description 创建新的应用，应用编码(code)为可选参数，不提供则根据应用名称自动生成
// @Tags Application
// @Accept json
// @Produce json
// @Param request body model.CreateApplicationRequest true "创建应用请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.Application}
// @Router /apps [post]
func (ac *ApplicationController) CreateApplication(c *gin.Context) {
	var req model.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.CreateApplication(c, &req)
}

// GetApplicationList 获取应用列表
// @Summary 获取应用列表
// @Description 获取应用列表，支持分页和筛选
// @Tags Application
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param name query string false "应用名称(模糊查询)"
// @Param code query string false "应用编码(模糊查询)"
// @Param app_type query string false "应用类型"
// @Param status query int false "状态"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ApplicationListResponse}
// @Router /apps [get]
func (ac *ApplicationController) GetApplicationList(c *gin.Context) {
	var req model.ApplicationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.GetApplicationList(c, &req)
}

// GetApplicationDetail 获取应用详情
// @Summary 获取应用详情
// @Description 根据应用ID获取应用详细信息
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.Application}
// @Router /apps/{id} [get]
func (ac *ApplicationController) GetApplicationDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	ac.appService.GetApplicationDetail(c, uint(id))
}

// UpdateApplication 更新应用
// @Summary 更新应用
// @Description 更新应用信息
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Param request body model.UpdateApplicationRequest true "更新应用请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.Application}
// @Router /apps/{id} [put]
func (ac *ApplicationController) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	var req model.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.UpdateApplication(c, uint(id), &req)
}

// DeleteApplication 删除应用
// @Summary 删除应用
// @Description 删除指定的应用
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/{id} [delete]
func (ac *ApplicationController) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	ac.appService.DeleteApplication(c, uint(id))
}

// AddAppJenkinsEnv 为应用添加Jenkins环境配置
// @Summary 为应用添加Jenkins环境配置
// @Description 为指定应用添加新的Jenkins环境配置
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Param request body model.CreateJenkinsEnvRequest true "Jenkins环境配置请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsEnv}
// @Router /apps/{id}/jenkins-envs [post]
func (ac *ApplicationController) AddAppJenkinsEnv(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	var req model.CreateJenkinsEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 确保应用ID匹配
	req.AppID = uint(appID)
	ac.appService.CreateJenkinsEnv(c, &req)
}

// UpdateAppJenkinsEnv 更新应用的Jenkins环境配置
// @Summary 更新应用的Jenkins环境配置
// @Description 更新指定应用的Jenkins环境配置
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Param env_id path int true "Jenkins环境配置ID"
// @Param request body model.UpdateJenkinsEnvRequest true "更新Jenkins环境配置请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsEnv}
// @Router /apps/{id}/jenkins-envs/{env_id} [put]
func (ac *ApplicationController) UpdateAppJenkinsEnv(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	envIDStr := c.Param("env_id")
	envID, err := strconv.ParseUint(envIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的Jenkins环境配置ID")
		return
	}

	var req model.UpdateJenkinsEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.UpdateAppJenkinsEnv(c, uint(appID), uint(envID), &req)
}

// DeleteAppJenkinsEnv 删除应用的Jenkins环境配置
// @Summary 删除应用的Jenkins环境配置
// @Description 删除指定应用的Jenkins环境配置
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Param env_id path int true "Jenkins环境配置ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/{id}/jenkins-envs/{env_id} [delete]
func (ac *ApplicationController) DeleteAppJenkinsEnv(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	envIDStr := c.Param("env_id")
	envID, err := strconv.ParseUint(envIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的Jenkins环境配置ID")
		return
	}

	ac.appService.DeleteAppJenkinsEnv(c, uint(appID), uint(envID))
}

// GetAppJenkinsEnvs 获取应用的所有Jenkins环境配置
// @Summary 获取应用的所有Jenkins环境配置
// @Description 根据应用ID获取所有Jenkins环境配置
// @Tags JenkinsEnv
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]model.JenkinsEnv}
// @Router /apps/{id}/jenkins-envs [get]
func (ac *ApplicationController) GetAppJenkinsEnvs(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的应用ID")
		return
	}

	ac.appService.GetAppJenkinsEnvs(c, uint(appID))
}


// GetJenkinsServers 获取Jenkins服务器列表
// @Summary 获取Jenkins服务器列表
// @Description 获取所有类型为Jenkins(type=4)的服务器配置信息，用于Jenkins环境配置选择
// @Tags Application
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]model.JenkinsServerOption}
// @Router /apps/jenkins-servers [get]
func (ac *ApplicationController) GetJenkinsServers(c *gin.Context) {
	ac.appService.GetJenkinsServers(c)
}

// ValidateJenkinsJob 验证Jenkins任务是否存在
// @Summary 验证Jenkins任务是否存在
// @Description 验证指定Jenkins服务器中是否存在指定的任务名称
// @Tags Application
// @Accept json
// @Produce json
// @Param request body model.ValidateJenkinsJobRequest true "验证Jenkins任务请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ValidateJenkinsJobResponse}
// @Router /apps/jenkins-job/validate [post]
func (ac *ApplicationController) ValidateJenkinsJob(c *gin.Context) {
	var req model.ValidateJenkinsJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.ValidateJenkinsJob(c, &req)
}

// GetApplicationsForDeployment 获取可发布的应用列表
// @Summary 获取可发布的应用列表
// @Description 根据业务组、部门和环境获取可发布的应用列表
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param business_group_id query int true "业务组ID"
// @Param business_dept_id query int true "业务部门ID"
// @Param environment query string true "环境名称"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]model.ApplicationForDeployment}
// @Router /apps/deployment/applications [get]
func (ac *ApplicationController) GetApplicationsForDeployment(c *gin.Context) {
	var req model.GetApplicationsForDeploymentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.GetApplicationsForDeployment(c, &req)
}

// CreateQuickDeployment 创建快速发布
// @Summary 创建快速发布
// @Description 创建快速发布流程，包含多个应用的发布任务
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param request body model.CreateQuickDeploymentRequest true "创建快速发布请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.QuickDeployment}
// @Router /apps/deployment/quick [post]
func (ac *ApplicationController) CreateQuickDeployment(c *gin.Context) {
	var req model.CreateQuickDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.CreateQuickDeployment(c, &req)
}

// ExecuteQuickDeployment 执行快速发布
// @Summary 执行快速发布
// @Description 启动快速发布流程，支持串行或并行执行模式
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param request body model.ExecuteQuickDeploymentRequest true "执行快速发布请求，支持选择执行模式"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/deployment/execute [post]
func (ac *ApplicationController) ExecuteQuickDeployment(c *gin.Context) {
	var req model.ExecuteQuickDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.ExecuteQuickDeployment(c, &req)
}

// GetQuickDeploymentList 获取快速发布列表
// @Summary 获取快速发布列表
// @Description 获取快速发布列表，支持分页和筛选
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param business_group_id query int false "业务组ID"
// @Param business_dept_id query int false "业务部门ID"
// @Param environment query string false "环境名称"
// @Param status query int false "状态"
// @Param creator_id query int false "创建人ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.QuickDeploymentListResponse}
// @Router /apps/deployment/list [get]
func (ac *ApplicationController) GetQuickDeploymentList(c *gin.Context) {
	var req model.QuickDeploymentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.GetQuickDeploymentList(c, &req)
}

// GetQuickDeploymentDetail 获取快速发布详情
// @Summary 获取快速发布详情
// @Description 获取快速发布详情，包含所有任务信息
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param id path int true "发布ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.QuickDeployment}
// @Router /apps/deployment/{id} [get]
func (ac *ApplicationController) GetQuickDeploymentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的发布ID")
		return
	}

	ac.appService.GetQuickDeploymentDetail(c, uint(id))
}

// GetServiceTree 获取业务线服务树
// @Summary 获取业务线服务树
// @Description 根据业务线查询服务，按业务线重新组装排序服务，类似于服务树结构
// @Tags ServiceTree
// @Accept json
// @Produce json
// @Param business_group_ids query []uint false "业务组ID列表，为空则查询所有"
// @Param status query int false "应用状态筛选，为空则查询所有状态"
// @Param environment query string false "环境名称筛选，为空则不筛选环境配置"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]model.BusinessLineServiceTree}
// @Router /apps/service-tree [get]
func (ac *ApplicationController) GetServiceTree(c *gin.Context) {
	var req model.GetServiceTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.GetServiceTree(c, &req)
}

// GetTaskBuildLog 获取任务构建日志
// @Summary 获取任务构建日志
// @Description 获取快速发布任务的Jenkins构建日志
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Param start query int false "日志起始位置" default(0)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Router /apps/deployment/tasks/{task_id}/log [get]
func (ac *ApplicationController) GetTaskBuildLog(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的任务ID")
		return
	}

	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))

	ac.appService.GetTaskBuildLog(c, uint(taskID), start)
}

// GetBusinessGroupOptions 获取业务组选项（连级选择器）
// @Summary 获取业务组选项
// @Description 获取业务组和业务部门的树形选择器数据，支持二级分组
// @Tags ServiceTree
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]map[string]interface{}}
// @Router /apps/business-group-options [get]
func (ac *ApplicationController) GetBusinessGroupOptions(c *gin.Context) {
	ac.appService.GetBusinessGroupOptions(c)
}

// GetAppEnvironment 获取单个应用的环境配置
// @Summary 获取应用环境配置
// @Description 获取指定应用在指定环境的详细配置信息，包括Jenkins配置状态
// @Tags Application
// @Accept json
// @Produce json
// @Param app_id query int true "应用ID"
// @Param environment query string true "环境名称"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.AppEnvironmentResponse}
// @Router /apps/environment [get]
func (ac *ApplicationController) GetAppEnvironment(c *gin.Context) {
	var req model.GetAppEnvironmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ac.appService.GetAppEnvironment(c, &req)
}

// GetTaskStatus 获取任务状态
// @Summary 获取任务状态
// @Description 获取快速发布任务的实时状态信息
// @Tags QuickDeployment
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.TaskStatusResponse}
// @Router /apps/deployment/tasks/{task_id}/status [get]
func (ac *ApplicationController) GetTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的任务ID")
		return
	}

	ac.appService.GetTaskStatus(c, uint(taskID))
}

// DeleteQuickDeployment 删除快速发布（级联删除所有子任务）
func (ac *ApplicationController) DeleteQuickDeployment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的发布ID")
		return
	}

	ac.appService.DeleteQuickDeployment(c, uint(id))
}