// 快捷导航工具 Service层
// author xiaoRui
package service

import (
	"dodevops-api/api/tool/dao1"
	"dodevops-api/api/tool/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// IToolService 接口定义
type IToolService interface {
	CreateTool(c *gin.Context, dto model.AddToolDto)      // 创建导航工具
	GetToolByID(c *gin.Context, id uint)                  // 根据ID获取导航工具
	UpdateTool(c *gin.Context, dto model.UpdateToolDto)   // 更新导航工具
	DeleteTool(c *gin.Context, id uint)                   // 删除导航工具
	GetToolList(c *gin.Context, dto model.ToolQueryDto)   // 获取导航工具列表（分页）
	GetAllTools(c *gin.Context)                           // 获取所有启用的导航工具
}

type ToolServiceImpl struct{}

// CreateTool 创建导航工具
func (s ToolServiceImpl) CreateTool(c *gin.Context, dto model.AddToolDto) {
	// 参数校验
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	// 创建导航工具
	err = dao.CreateTool(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetToolByID 根据ID获取导航工具
func (s ToolServiceImpl) GetToolByID(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	tool, err := dao.GetToolByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	result.Success(c, tool)
}

// UpdateTool 更新导航工具
func (s ToolServiceImpl) UpdateTool(c *gin.Context, dto model.UpdateToolDto) {
	// 参数校验
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	// 检查是否存在
	_, err = dao.GetToolByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "导航工具不存在")
		return
	}

	// 更新导航工具
	err = dao.UpdateTool(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// DeleteTool 删除导航工具
func (s ToolServiceImpl) DeleteTool(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	// 检查是否存在
	_, err := dao.GetToolByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "导航工具不存在")
		return
	}

	// 删除导航工具
	err = dao.DeleteTool(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetToolList 获取导航工具列表（分页）
func (s ToolServiceImpl) GetToolList(c *gin.Context, dto model.ToolQueryDto) {
	// 设置默认分页参数
	if dto.PageNum <= 0 {
		dto.PageNum = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	tools, total, err := dao.GetToolList(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	// 返回分页数据
	data := map[string]interface{}{
		"list":     tools,
		"total":    total,
		"pageNum":  dto.PageNum,
		"pageSize": dto.PageSize,
	}

	result.Success(c, data)
}

// GetAllTools 获取所有导航工具
func (s ToolServiceImpl) GetAllTools(c *gin.Context) {
	tools, err := dao.GetAllTools()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	result.Success(c, tools)
}

var toolService = ToolServiceImpl{}

// ToolService 获取服务实例
func ToolService() IToolService {
	return &toolService
}
