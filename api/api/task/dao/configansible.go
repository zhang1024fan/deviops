package dao

import (
	"dodevops-api/api/task/model"

	"gorm.io/gorm"
)

type ConfigAnsibleDao struct {
	DB *gorm.DB
}

func NewConfigAnsibleDao(db *gorm.DB) *ConfigAnsibleDao {
	return &ConfigAnsibleDao{DB: db}
}

// Create 创建配置
func (d *ConfigAnsibleDao) Create(config *model.ConfigAnsible) error {
	return d.DB.Create(config).Error
}

// Update 更新配置
func (d *ConfigAnsibleDao) Update(config *model.ConfigAnsible) error {
	return d.DB.Save(config).Error
}

// Delete 删除配置
func (d *ConfigAnsibleDao) Delete(id uint) error {
	return d.DB.Delete(&model.ConfigAnsible{}, id).Error
}

// GetByID 根据ID获取配置
func (d *ConfigAnsibleDao) GetByID(id uint) (*model.ConfigAnsible, error) {
	var config model.ConfigAnsible
	err := d.DB.First(&config, id).Error
	return &config, err
}

// ListResponse 列表响应结构
type ListResponse struct {
	List  []model.ConfigAnsible `json:"list"`
	Total int64                 `json:"total"`
}

// List 获取配置列表（支持多条件查询）
func (d *ConfigAnsibleDao) List(page, size int, name string, configType int) (*ListResponse, error) {
	var configs []model.ConfigAnsible
	var total int64
	db := d.DB.Model(&model.ConfigAnsible{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if configType > 0 {
		db = db.Where("type = ?", configType)
	}

	err := db.Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Order("id desc").
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	return &ListResponse{List: configs, Total: total}, nil
}