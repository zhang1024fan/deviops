package dao

import (
	"dodevops-api/api/task/model"
	"gorm.io/gorm"
)

type TaskTemplateDaoInterface interface {
	CreateTemplate(template *model.TaskTemplate) error
	GetAllTemplates() ([]model.TaskTemplate, error)
	UpdateTemplate(template *model.TaskTemplate) error
	DeleteTemplate(id uint) error
	GetTemplateByID(id uint) (*model.TaskTemplate, error)
	GetTemplatesByName(name string) ([]model.TaskTemplate, error)
	GetTemplatesByType(templateType int) ([]model.TaskTemplate, error)
}

type TaskTemplateDao struct {
	db *gorm.DB
}

func NewTaskTemplateDao(db *gorm.DB) *TaskTemplateDao {
	return &TaskTemplateDao{db: db}
}

func (d *TaskTemplateDao) CreateTemplate(template *model.TaskTemplate) error {
	return d.db.Create(template).Error
}

func (d *TaskTemplateDao) GetAllTemplates() ([]model.TaskTemplate, error) {
	var templates []model.TaskTemplate
	if err := d.db.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (d *TaskTemplateDao) UpdateTemplate(template *model.TaskTemplate) error {
	return d.db.Model(template).Select("name", "type", "content", "remark", "updated_at", "updated_by").Updates(template).Error
}

func (d *TaskTemplateDao) DeleteTemplate(id uint) error {
	return d.db.Delete(&model.TaskTemplate{}, id).Error
}

func (d *TaskTemplateDao) GetTemplateByID(id uint) (*model.TaskTemplate, error) {
	var template model.TaskTemplate
	if err := d.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (d *TaskTemplateDao) GetTemplatesByName(name string) ([]model.TaskTemplate, error) {
	var templates []model.TaskTemplate
	if err := d.db.Where("name LIKE ?", "%"+name+"%").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (d *TaskTemplateDao) GetTemplatesByType(templateType int) ([]model.TaskTemplate, error) {
	var templates []model.TaskTemplate
	if err := d.db.Where("type = ?", templateType).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
