package service

import (
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	"dodevops-api/common/result"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IConfigAnsibleService interface {
	Create(c *gin.Context, req *CreateConfigRequest)
	Update(c *gin.Context, id uint, req *UpdateConfigRequest)
	Delete(c *gin.Context, id uint)
	Get(c *gin.Context, id uint)
	List(c *gin.Context, page, size int, name string, configType int)
}

type ConfigAnsibleServiceImpl struct {
	dao *dao.ConfigAnsibleDao
}

func NewConfigAnsibleService(db *gorm.DB) IConfigAnsibleService {
	return &ConfigAnsibleServiceImpl{
		dao: dao.NewConfigAnsibleDao(db),
	}
}

type CreateConfigRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    int    `json:"type" binding:"required,oneof=1 2 3 4"` // 1-inventory 2-global_vars 3-extra_vars 4-cli_args
	Content string `json:"content" binding:"required"`
	Remark  string `json:"remark"`
}

type UpdateConfigRequest struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Remark  string `json:"remark"`
}

func (s *ConfigAnsibleServiceImpl) Create(c *gin.Context, req *CreateConfigRequest) {
	// 检查名称是否存在
	existing, _ := s.dao.List(1, 1, req.Name, 0)
	if existing.Total > 0 {
		for _, item := range existing.List {
			if item.Name == req.Name {
				result.Failed(c, 400, "配置名称已存在")
				return
			}
		}
	}

	// 获取当前用户（假设从中间件获取，这里模拟）
	username := c.GetString("username")
	if username == "" {
		username = "system"
	}

	config := &model.ConfigAnsible{
		Name:      req.Name,
		Type:      req.Type,
		Content:   req.Content,
		Remark:    req.Remark,
		CreatedBy: username,
		UpdatedBy: username,
	}

	if err := s.dao.Create(config); err != nil {
		result.Failed(c, 500, "创建配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) Update(c *gin.Context, id uint, req *UpdateConfigRequest) {
	config, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "配置不存在")
		return
	}

	if req.Name != "" {
		config.Name = req.Name
	}
	if req.Type > 0 {
		config.Type = req.Type
	}
	if req.Content != "" {
		config.Content = req.Content
	}
	if req.Remark != "" {
		config.Remark = req.Remark
	}

	username := c.GetString("username")
	if username == "" {
		username = "system"
	}
	config.UpdatedBy = username
	config.UpdatedAt = time.Now()

	if err := s.dao.Update(config); err != nil {
		result.Failed(c, 500, "更新配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) Delete(c *gin.Context, id uint) {
	if err := s.dao.Delete(id); err != nil {
		result.Failed(c, 500, "删除配置失败: "+err.Error())
		return
	}
	result.Success(c, nil)
}

func (s *ConfigAnsibleServiceImpl) Get(c *gin.Context, id uint) {
	config, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "配置不存在")
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) List(c *gin.Context, page, size int, name string, configType int) {
	data, err := s.dao.List(page, size, name, configType)
	if err != nil {
		result.Failed(c, 500, "获取配置列表失败: "+err.Error())
		return
	}
	result.Success(c, data)
}