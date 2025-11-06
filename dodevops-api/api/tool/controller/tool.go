// 快捷导航工具 Controller层
// author xiaoRui
package controller

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/api/tool/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTool 创建导航工具
// @Tags Tool导航工具
// @Summary 创建导航工具
// @Description 创建导航工具接口
// @Accept json
// @Produce json
// @Param data body model.AddToolDto true "导航工具信息"
// @Success 200 {object} result.Result
// @Router /api/v1/tool [post]
// @Security ApiKeyAuth
func CreateTool(c *gin.Context) {
	var dto model.AddToolDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.ToolService().CreateTool(c, dto)
}

// GetToolByID 根据ID获取导航工具
// @Tags Tool导航工具
// @Summary 根据ID获取导航工具
// @Description 根据ID获取导航工具详情
// @Produce json
// @Param id path int true "工具ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/{id} [get]
// @Security ApiKeyAuth
func GetToolByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.ToolService().GetToolByID(c, uint(id))
}

// UpdateTool 更新导航工具
// @Tags Tool导航工具
// @Summary 更新导航工具
// @Description 更新导航工具接口
// @Accept json
// @Produce json
// @Param data body model.UpdateToolDto true "导航工具信息"
// @Success 200 {object} result.Result
// @Router /api/v1/tool [put]
// @Security ApiKeyAuth
func UpdateTool(c *gin.Context) {
	var dto model.UpdateToolDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.ToolService().UpdateTool(c, dto)
}

// DeleteTool 删除导航工具
// @Tags Tool导航工具
// @Summary 删除导航工具
// @Description 删除导航工具接口
// @Produce json
// @Param id path int true "工具ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/{id} [delete]
// @Security ApiKeyAuth
func DeleteTool(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.ToolService().DeleteTool(c, uint(id))
}

// GetToolList 获取导航工具列表（分页）
// @Tags Tool导航工具
// @Summary 获取导航工具列表
// @Description 获取导航工具列表（分页）
// @Produce json
// @Param title query string false "标题（模糊查询）"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/list [get]
// @Security ApiKeyAuth
func GetToolList(c *gin.Context) {
	var dto model.ToolQueryDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.ToolService().GetToolList(c, dto)
}

// GetAllTools 获取所有导航工具
// @Tags Tool导航工具
// @Summary 获取所有导航工具
// @Description 获取所有导航工具（不分页）
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/tool/all [get]
// @Security ApiKeyAuth
func GetAllTools(c *gin.Context) {
	service.ToolService().GetAllTools(c)
}
