// 服务部署 DAO层
// author xiaoRui
package dao

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/pkg/db"
)

// CreateServiceDeploy 创建服务部署记录
func CreateServiceDeploy(deploy *model.ServiceDeploy) error {
	return db.Db.Create(deploy).Error
}

// GetServiceDeployByID 根据ID获取部署记录
func GetServiceDeployByID(id uint) (*model.ServiceDeploy, error) {
	var deploy model.ServiceDeploy
	err := db.Db.Where("id = ?", id).First(&deploy).Error
	return &deploy, err
}

// UpdateServiceDeploy 更新部署记录
func UpdateServiceDeploy(deploy *model.ServiceDeploy) error {
	return db.Db.Save(deploy).Error
}

// UpdateServiceDeployStatus 更新部署状态
func UpdateServiceDeployStatus(id uint, status int, log string) error {
	updates := map[string]interface{}{
		"status":     status,
		"deploy_log": log,
	}
	return db.Db.Model(&model.ServiceDeploy{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteServiceDeploy 删除部署记录
func DeleteServiceDeploy(id uint) error {
	return db.Db.Where("id = ?", id).Delete(&model.ServiceDeploy{}).Error
}

// GetServiceDeployList 获取部署记录列表（分页）
func GetServiceDeployList(dto model.DeployQueryDto) ([]model.ServiceDeploy, int64, error) {
	var deploys []model.ServiceDeploy
	var total int64

	query := db.Db.Model(&model.ServiceDeploy{})

	// 条件筛选
	if dto.ServiceName != "" {
		query = query.Where("service_name LIKE ?", "%"+dto.ServiceName+"%")
	}
	if dto.HostID > 0 {
		query = query.Where("host_id = ?", dto.HostID)
	}
	// Status: nil表示查询全部，非nil表示查询指定状态
	if dto.Status != nil {
		query = query.Where("status = ?", *dto.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (dto.PageNum - 1) * dto.PageSize
	err := query.Order("create_time DESC").Limit(dto.PageSize).Offset(offset).Find(&deploys).Error

	return deploys, total, err
}

// GetAllServiceDeploys 获取所有部署记录
func GetAllServiceDeploys() ([]model.ServiceDeploy, error) {
	var deploys []model.ServiceDeploy
	err := db.Db.Order("create_time DESC").Find(&deploys).Error
	return deploys, err
}

// CheckServiceDeployExists 检查服务是否已在主机上部署
func CheckServiceDeployExists(hostID uint, serviceID string) (bool, error) {
	var count int64
	err := db.Db.Model(&model.ServiceDeploy{}).
		Where("host_id = ? AND service_id = ? AND status IN (0, 1)", hostID, serviceID).
		Count(&count).Error
	return count > 0, err
}

// GetServiceDeployByHostAndService 根据主机和服务获取部署记录
func GetServiceDeployByHostAndService(hostID uint, serviceID string) (*model.ServiceDeploy, error) {
	var deploy model.ServiceDeploy
	err := db.Db.Where("host_id = ? AND service_id = ?", hostID, serviceID).
		Order("create_time DESC").
		First(&deploy).Error
	return &deploy, err
}
