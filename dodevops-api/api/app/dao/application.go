package dao

import (
	"dodevops-api/api/app/model"
	"gorm.io/gorm"
)

// IApplicationDao 应用DAO接口
type IApplicationDao interface {
	// 应用管理
	CreateApplication(app *model.Application) error
	GetApplicationByID(id uint) (*model.Application, error)
	GetApplicationByCode(code string) (*model.Application, error)
	UpdateApplication(id uint, updates map[string]interface{}) error
	DeleteApplication(id uint) error
	GetApplicationList(req *model.ApplicationListRequest) ([]model.Application, int64, error)

	// Jenkins环境配置管理
	CreateJenkinsEnv(env *model.JenkinsEnv) error
	GetJenkinsEnvByID(id uint) (*model.JenkinsEnv, error)
	GetJenkinsEnvByAppAndEnv(appID uint, envName string) (*model.JenkinsEnv, error)
	UpdateJenkinsEnv(id uint, updates map[string]interface{}) error
	DeleteJenkinsEnv(id uint) error
	GetJenkinsEnvList(req *model.JenkinsEnvListRequest) ([]model.JenkinsEnv, int64, error)
	GetJenkinsEnvsByAppID(appID uint) ([]model.JenkinsEnv, error)
}

// ApplicationDao 应用DAO实现
type ApplicationDao struct {
	db *gorm.DB
}

// NewApplicationDao 创建应用DAO
func NewApplicationDao(db *gorm.DB) IApplicationDao {
	return &ApplicationDao{db: db}
}

// CreateApplication 创建应用
func (d *ApplicationDao) CreateApplication(app *model.Application) error {
	return d.db.Create(app).Error
}

// GetApplicationByID 根据ID获取应用
func (d *ApplicationDao) GetApplicationByID(id uint) (*model.Application, error) {
	var app model.Application
	err := d.db.Preload("JenkinsEnvs").First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationByCode 根据编码获取应用
func (d *ApplicationDao) GetApplicationByCode(code string) (*model.Application, error) {
	var app model.Application
	err := d.db.Where("code = ?", code).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// UpdateApplication 更新应用
func (d *ApplicationDao) UpdateApplication(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.Application{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteApplication 删除应用（物理删除）
func (d *ApplicationDao) DeleteApplication(id uint) error {
	return d.db.Unscoped().Delete(&model.Application{}, id).Error
}

// GetApplicationList 获取应用列表
func (d *ApplicationDao) GetApplicationList(req *model.ApplicationListRequest) ([]model.Application, int64, error) {
	var apps []model.Application
	var total int64

	query := d.db.Model(&model.Application{})

	// 构建查询条件
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.ProgrammingLang != "" {
		query = query.Where("programming_lang = ?", req.ProgrammingLang)
	}
	if req.BusinessGroupID != nil {
		query = query.Where("business_group_id = ?", *req.BusinessGroupID)
	}
	if req.BusinessDeptID != nil {
		query = query.Where("business_dept_id = ?", *req.BusinessDeptID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// CreateJenkinsEnv 创建Jenkins环境配置
func (d *ApplicationDao) CreateJenkinsEnv(env *model.JenkinsEnv) error {
	return d.db.Create(env).Error
}

// GetJenkinsEnvByID 根据ID获取Jenkins环境配置
func (d *ApplicationDao) GetJenkinsEnvByID(id uint) (*model.JenkinsEnv, error) {
	var env model.JenkinsEnv
	err := d.db.Preload("Application").First(&env, id).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// GetJenkinsEnvByAppAndEnv 根据应用ID和环境名称获取配置
func (d *ApplicationDao) GetJenkinsEnvByAppAndEnv(appID uint, envName string) (*model.JenkinsEnv, error) {
	var env model.JenkinsEnv
	err := d.db.Where("app_id = ? AND env_name = ?", appID, envName).First(&env).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// UpdateJenkinsEnv 更新Jenkins环境配置
func (d *ApplicationDao) UpdateJenkinsEnv(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.JenkinsEnv{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteJenkinsEnv 删除Jenkins环境配置（物理删除）
func (d *ApplicationDao) DeleteJenkinsEnv(id uint) error {
	return d.db.Unscoped().Delete(&model.JenkinsEnv{}, id).Error
}

// GetJenkinsEnvList 获取Jenkins环境配置列表
func (d *ApplicationDao) GetJenkinsEnvList(req *model.JenkinsEnvListRequest) ([]model.JenkinsEnv, int64, error) {
	var envs []model.JenkinsEnv
	var total int64

	query := d.db.Model(&model.JenkinsEnv{}).Preload("Application")

	// 构建查询条件
	if req.AppID != 0 {
		query = query.Where("app_id = ?", req.AppID)
	}
	if req.EnvName != "" {
		query = query.Where("env_name = ?", req.EnvName)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	if err := query.Order("created_at DESC").Find(&envs).Error; err != nil {
		return nil, 0, err
	}

	return envs, total, nil
}

// GetJenkinsEnvsByAppID 根据应用ID获取所有环境配置
func (d *ApplicationDao) GetJenkinsEnvsByAppID(appID uint) ([]model.JenkinsEnv, error) {
	var envs []model.JenkinsEnv
	err := d.db.Where("app_id = ?", appID).Order("env_name").Find(&envs).Error
	return envs, err
}