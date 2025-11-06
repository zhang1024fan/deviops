package controller

import (
	"strconv"

	"dodevops-api/api/app/model"
	"dodevops-api/api/app/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JenkinsController Jenkins控制器
type JenkinsController struct {
	jenkinsService service.IJenkinsService
}

// NewJenkinsController 创建Jenkins控制器
func NewJenkinsController(db *gorm.DB) *JenkinsController {
	return &JenkinsController{
		jenkinsService: service.NewJenkinsService(db),
	}
}

// GetJenkinsServers 获取Jenkins服务器列表
// @Summary 获取Jenkins服务器列表
// @Description 获取所有配置的Jenkins服务器
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsServerListResponse}
// @Router /jenkins/servers [get]
func (jc *JenkinsController) GetJenkinsServers(c *gin.Context) {
	jc.jenkinsService.GetJenkinsServers(c)
}

// GetJenkinsServerDetail 获取Jenkins服务器详情
// @Summary 获取Jenkins服务器详情
// @Description 根据服务器ID获取Jenkins服务器详情
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param id path int true "服务器ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsServerInfo}
// @Router /jenkins/servers/{id} [get]
func (jc *JenkinsController) GetJenkinsServerDetail(c *gin.Context) {
	serverIDStr := c.Param("id")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jc.jenkinsService.GetJenkinsServerDetail(c, uint(serverID))
}

// TestJenkinsConnection 测试Jenkins连接
// @Summary 测试Jenkins连接
// @Description 测试Jenkins服务器连接是否正常
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param request body model.TestJenkinsConnectionRequest true "连接测试请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.TestJenkinsConnectionResponse}
// @Router /jenkins/test-connection [post]
func (jc *JenkinsController) TestJenkinsConnection(c *gin.Context) {
	var req model.TestJenkinsConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	jc.jenkinsService.TestJenkinsConnection(c, &req)
}

// GetJobs 获取Jenkins任务列表
// @Summary 获取Jenkins任务列表
// @Description 获取指定Jenkins服务器的所有任务
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsJobListResponse}
// @Router /jenkins/{serverId}/jobs [get]
func (jc *JenkinsController) GetJobs(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jc.jenkinsService.GetJobs(c, uint(serverID))
}

// SearchJobs 搜索Jenkins任务
// @Summary 搜索Jenkins任务
// @Description 根据关键词模糊搜索指定Jenkins服务器的任务
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param keyword query string true "搜索关键词"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsJobListResponse}
// @Router /jenkins/{serverId}/jobs/search [get]
func (jc *JenkinsController) SearchJobs(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		result.Failed(c, 400, "搜索关键词不能为空")
		return
	}

	jc.jenkinsService.SearchJobs(c, uint(serverID), keyword)
}

// GetJobDetail 获取Jenkins任务详情
// @Summary 获取Jenkins任务详情
// @Description 获取指定任务的详细信息和构建历史
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param jobName path string true "任务名称"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsJobDetailResponse}
// @Router /jenkins/{serverId}/jobs/{jobName} [get]
func (jc *JenkinsController) GetJobDetail(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jobName := c.Param("jobName")
	if jobName == "" {
		result.Failed(c, 400, "任务名称不能为空")
		return
	}

	jc.jenkinsService.GetJobDetail(c, uint(serverID), jobName)
}

// StartJob 启动Jenkins任务
// @Summary 启动Jenkins任务
// @Description 启动指定的Jenkins任务，支持带参数构建
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param jobName path string true "任务名称"
// @Param request body model.StartJobRequest false "启动任务请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.StartJobResponse}
// @Router /jenkins/{serverId}/jobs/{jobName}/start [post]
func (jc *JenkinsController) StartJob(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jobName := c.Param("jobName")
	if jobName == "" {
		result.Failed(c, 400, "任务名称不能为空")
		return
	}

	var req model.StartJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	jc.jenkinsService.StartJob(c, uint(serverID), jobName, &req)
}

// StopBuild 停止Jenkins构建
// @Summary 停止Jenkins构建
// @Description 停止指定的Jenkins构建任务
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param jobName path string true "任务名称"
// @Param buildNumber path int true "构建编号"
// @Param request body model.StopBuildRequest false "停止构建请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.StopBuildResponse}
// @Router /jenkins/{serverId}/jobs/{jobName}/builds/{buildNumber}/stop [post]
func (jc *JenkinsController) StopBuild(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jobName := c.Param("jobName")
	if jobName == "" {
		result.Failed(c, 400, "任务名称不能为空")
		return
	}

	buildNumberStr := c.Param("buildNumber")
	buildNumber, err := strconv.Atoi(buildNumberStr)
	if err != nil {
		result.Failed(c, 400, "无效的构建编号")
		return
	}

	var req model.StopBuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	jc.jenkinsService.StopBuild(c, uint(serverID), jobName, buildNumber, &req)
}

// GetBuildDetail 获取Jenkins构建详情
// @Summary 获取Jenkins构建详情
// @Description 获取指定构建的详细信息
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param jobName path string true "任务名称"
// @Param buildNumber path int true "构建编号"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsBuildDetailResponse}
// @Router /jenkins/{serverId}/jobs/{jobName}/builds/{buildNumber} [get]
func (jc *JenkinsController) GetBuildDetail(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jobName := c.Param("jobName")
	if jobName == "" {
		result.Failed(c, 400, "任务名称不能为空")
		return
	}

	buildNumberStr := c.Param("buildNumber")
	buildNumber, err := strconv.Atoi(buildNumberStr)
	if err != nil {
		result.Failed(c, 400, "无效的构建编号")
		return
	}

	jc.jenkinsService.GetBuildDetail(c, uint(serverID), jobName, buildNumber)
}

// GetBuildLog 获取Jenkins构建日志
// @Summary 获取Jenkins构建日志
// @Description 获取指定构建的日志信息，支持分页获取
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Param jobName path string true "任务名称"
// @Param buildNumber path int true "构建编号"
// @Param start query int false "开始位置" default(0)
// @Param html query bool false "是否返回HTML格式" default(false)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.GetBuildLogResponse}
// @Router /jenkins/{serverId}/jobs/{jobName}/builds/{buildNumber}/log [get]
func (jc *JenkinsController) GetBuildLog(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jobName := c.Param("jobName")
	if jobName == "" {
		result.Failed(c, 400, "任务名称不能为空")
		return
	}

	buildNumberStr := c.Param("buildNumber")
	buildNumber, err := strconv.Atoi(buildNumberStr)
	if err != nil {
		result.Failed(c, 400, "无效的构建编号")
		return
	}

	// 获取查询参数
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	html := c.DefaultQuery("html", "false") == "true"

	req := &model.GetBuildLogRequest{
		Start: start,
		Html:  html,
	}

	jc.jenkinsService.GetBuildLog(c, uint(serverID), jobName, buildNumber, req)
}

// GetSystemInfo 获取Jenkins系统信息
// @Summary 获取Jenkins系统信息
// @Description 获取指定Jenkins服务器的系统信息
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsSystemInfo}
// @Router /jenkins/{serverId}/system-info [get]
func (jc *JenkinsController) GetSystemInfo(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jc.jenkinsService.GetSystemInfo(c, uint(serverID))
}

// GetQueueInfo 获取Jenkins队列信息
// @Summary 获取Jenkins队列信息
// @Description 获取指定Jenkins服务器的构建队列信息
// @Tags Jenkins
// @Accept json
// @Produce json
// @Param serverId path int true "服务器ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.JenkinsQueue}
// @Router /jenkins/{serverId}/queue [get]
func (jc *JenkinsController) GetQueueInfo(c *gin.Context) {
	serverIDStr := c.Param("serverId")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务器ID")
		return
	}

	jc.jenkinsService.GetQueueInfo(c, uint(serverID))
}