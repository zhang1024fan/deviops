package controller

import (
	"dodevops-api/api/task/service"
	"dodevops-api/common/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConfigAnsibleController struct {
	service service.IConfigAnsibleService
}

func NewConfigAnsibleController(service service.IConfigAnsibleService) *ConfigAnsibleController {
	return &ConfigAnsibleController{service: service}
}

// Create 创建配置
// @Summary 创建Ansible配置
// @Description 创建Inventory/Vars/Args等配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param request body service.CreateConfigRequest true "创建配置请求"
// @Success 200 {object} result.Result{data=model.ConfigAnsible}
// @Router /api/v1/config/ansible [post]
// @Security ApiKeyAuth
func (c *ConfigAnsibleController) Create(ctx *gin.Context) {
	var req service.CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Create(ctx, &req)
}

// Update 更新配置
// @Summary 更新Ansible配置
// @Description 更新配置内容
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Param request body service.UpdateConfigRequest true "更新配置请求"
// @Success 200 {object} result.Result{data=model.ConfigAnsible}
// @Router /api/v1/config/ansible/{id} [put]
// @Security ApiKeyAuth
func (c *ConfigAnsibleController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的ID")
		return
	}
	var req service.UpdateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Update(ctx, uint(id), &req)
}

// Delete 删除配置
// @Summary 删除Ansible配置
// @Description 删除指定的配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/ansible/{id} [delete]
// @Security ApiKeyAuth
func (c *ConfigAnsibleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的ID")
		return
	}
	c.service.Delete(ctx, uint(id))
}

// Get 获取配置详情
// @Summary 获取Ansible配置详情
// @Description 根据ID获取配置详情
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} result.Result{data=model.ConfigAnsible}
// @Router /api/v1/config/ansible/{id} [get]
// @Security ApiKeyAuth
func (c *ConfigAnsibleController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的ID")
		return
	}
	c.service.Get(ctx, uint(id))
}

// List 获取配置列表
// @Summary 获取Ansible配置列表
// @Description 分页获取配置列表，支持按名称和类型过滤
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "配置名称（模糊查询）"
// @Param type query int false "配置类型(1-inventory 2-global_vars 3-extra_vars 4-cli_args)"
// @Success 200 {object} result.Result{data=dao.ListResponse}
// @Router /api/v1/config/ansible [get]
// @Security ApiKeyAuth
func (c *ConfigAnsibleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	name := ctx.Query("name")
	configType, _ := strconv.Atoi(ctx.Query("type"))

	c.service.List(ctx, page, size, name, configType)
}