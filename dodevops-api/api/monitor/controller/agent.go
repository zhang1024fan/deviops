package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/config"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AgentController struct {
	agentService service.AgentServiceInterface
}

func NewAgentController() *AgentController {
	return &AgentController{
		agentService: service.NewAgentService(),
	}
}

// DeployAgent 部署agent到指定主机(支持单个或多个)
// @Summary 部署agent到指定主机(支持单个或多个)
// @Description 自动编译agent二进制文件，拷贝到目标主机并启动服务，单个主机传[hostId]，多个主机传[hostId1,hostId2,hostId3]
// @Tags 监控
// @Accept json
// @Produce json
// @Param request body model.BatchDeployAgentDto true "部署参数"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/deploy [post]
// @Security ApiKeyAuth
func (c *AgentController) DeployAgent(ctx *gin.Context) {
	var dto model.BatchDeployAgentDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 如果只有一个主机ID，调用单个部署方法
	if len(dto.HostIDs) == 1 {
		c.agentService.DeployAgent(ctx, dto.HostIDs[0])
	} else {
		// 多个主机ID，调用批量部署方法
		c.agentService.BatchDeployAgent(ctx, &dto)
	}
}


// UninstallAgent 卸载指定主机的agent(支持单个或多个)
// @Summary 卸载指定主机的agent(支持单个或多个)
// @Description 停止agent服务并删除相关文件，单个主机传[hostId]，多个主机传[hostId1,hostId2,hostId3]
// @Tags 监控
// @Accept json
// @Produce json
// @Param request body model.BatchDeployAgentDto true "卸载参数(只需hostIds字段)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/uninstall [delete]
// @Security ApiKeyAuth
func (c *AgentController) UninstallAgent(ctx *gin.Context) {
	var dto model.BatchDeployAgentDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 如果只有一个主机ID，调用单个卸载方法
	if len(dto.HostIDs) == 1 {
		c.agentService.UninstallAgent(ctx, dto.HostIDs[0])
	} else {
		// 多个主机ID，调用批量卸载方法
		c.agentService.BatchUninstallAgent(ctx, &dto)
	}
}

// GetAgentStatus 获取agent状态
// @Summary 获取agent状态
// @Description 获取指定主机上agent的运行状态、版本信息等
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/status/{id} [get]
// @Security ApiKeyAuth
func (c *AgentController) GetAgentStatus(ctx *gin.Context) {
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	c.agentService.GetAgentStatus(ctx, uint(hostID))
}

// RestartAgent 重启agent
// @Summary 重启agent
// @Description 重启指定主机上的agent服务
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/restart/{id} [post]
// @Security ApiKeyAuth
func (c *AgentController) RestartAgent(ctx *gin.Context) {
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	c.agentService.RestartAgent(ctx, uint(hostID))
}

// GetAgentList 获取Agent列表
// @Summary 获取Agent列表
// @Description 获取所有Agent的列表信息，支持分页和筛选
// @Tags 监控
// @Accept json
// @Produce json
// @Param hostId query uint false "主机ID"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param pageSize query int false "页大小"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/list [get]
// @Security ApiKeyAuth
func (c *AgentController) GetAgentList(ctx *gin.Context) {
	// 解析查询参数
	dto := &model.AgentListDto{}

	if hostIDStr := ctx.Query("hostId"); hostIDStr != "" {
		if hostID, err := strconv.ParseUint(hostIDStr, 10, 32); err == nil {
			dto.HostID = uint(hostID)
		}
	}

	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			dto.Status = status
		}
	}

	// 注意：Platform字段已移除，因为Agent只支持Linux平台

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			dto.Page = page
		}
	}

	if pageSizeStr := ctx.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			dto.PageSize = pageSize
		}
	}

	// 设置默认值
	if dto.Page <= 0 {
		dto.Page = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	c.agentService.GetAgentList(ctx, dto)
}

// UpdateHeartbeat 更新Agent心跳
// @Summary 更新Agent心跳
// @Description Agent主动上报心跳信息，通过IP自动识别主机
// @Tags 监控
// @Accept json
// @Produce json
// @Param heartbeat body model.AgentHeartbeatDto true "心跳数据"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/heartbeat [post]
func (c *AgentController) UpdateHeartbeat(ctx *gin.Context) {
	var heartbeat model.AgentHeartbeatDto
	if err := ctx.ShouldBindJSON(&heartbeat); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// Token验证
	expectedToken := config.Config.Monitor.Agent.HeartbeatToken
	if expectedToken == "" {
		expectedToken = "agent-heartbeat-token-2024" // 默认token
	}
	if heartbeat.Token != expectedToken {
		result.Failed(ctx, http.StatusUnauthorized, "Token验证失败")
		return
	}

	// 如果IP为空，从请求中获取客户端IP
	if heartbeat.IP == "" {
		clientIP := ctx.ClientIP()
		if clientIP == "" {
			result.Failed(ctx, http.StatusBadRequest, "无法获取客户端IP")
			return
		}
		heartbeat.IP = clientIP
	}

	c.agentService.UpdateHeartbeatByIP(ctx, &heartbeat)
}

// GetAgentStatistics 获取Agent统计信息
// @Summary 获取Agent统计信息
// @Description 获取Agent的统计信息，包括各状态数量、平台分布等
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/statistics [get]
// @Security ApiKeyAuth
func (c *AgentController) GetAgentStatistics(ctx *gin.Context) {
	c.agentService.GetAgentStatistics(ctx)
}

// DeleteAgent 删除agent数据
// @Summary 删除agent数据
// @Description 删除指定的agent数据，用于服务器离线无法正常卸载的情况
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "Agent ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/agent/delete/{id} [delete]
// @Security ApiKeyAuth
func (c *AgentController) DeleteAgent(ctx *gin.Context) {
	agentIDStr := ctx.Param("id")
	agentID, err := strconv.ParseUint(agentIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的Agent ID")
		return
	}

	c.agentService.DeleteAgentByID(ctx, uint(agentID))
}
