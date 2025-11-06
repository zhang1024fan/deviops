package dao

import (
	"dodevops-api/api/k8s/model"
	"gorm.io/gorm"
)

type KubeClusterDao struct {
	DB *gorm.DB
}

func NewKubeClusterDao(db *gorm.DB) *KubeClusterDao {
	return &KubeClusterDao{DB: db}
}

// Create 创建集群
func (d *KubeClusterDao) Create(cluster *model.KubeCluster) error {
	return d.DB.Create(cluster).Error
}

// GetByID 根据ID获取集群
func (d *KubeClusterDao) GetByID(id uint) (*model.KubeCluster, error) {
	var cluster model.KubeCluster
	err := d.DB.Where("id = ?", id).First(&cluster).Error
	return &cluster, err
}

// GetByName 根据名称获取集群
func (d *KubeClusterDao) GetByName(name string) (*model.KubeCluster, error) {
	var cluster model.KubeCluster
	err := d.DB.Where("name = ?", name).First(&cluster).Error
	return &cluster, err
}

// List 获取集群列表
func (d *KubeClusterDao) List(page, size int) ([]model.KubeCluster, int64, error) {
	var clusters []model.KubeCluster
	var total int64

	// 计算总数
	if err := d.DB.Model(&model.KubeCluster{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := d.DB.Offset(offset).
		Limit(size).
		Order("created_at DESC").
		Find(&clusters).Error

	return clusters, total, err
}

// Update 更新集群
func (d *KubeClusterDao) Update(id uint, updates map[string]interface{}) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除集群
func (d *KubeClusterDao) Delete(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.KubeCluster{}).Error
}

// UpdateStatus 更新集群状态
func (d *KubeClusterDao) UpdateStatus(id uint, status int) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateCredential 更新集群凭证
func (d *KubeClusterDao) UpdateCredential(id uint, credential string) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Update("credential", credential).Error
}

// IsClusterNameExists 检查集群名称是否存在
func (d *KubeClusterDao) IsClusterNameExists(name string) (bool, error) {
	var count int64
	err := d.DB.Model(&model.KubeCluster{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// GetClusterCountByStatus 根据状态统计集群数量
func (d *KubeClusterDao) GetClusterCountByStatus(status int) (int64, error) {
	var count int64
	err := d.DB.Model(&model.KubeCluster{}).Where("status = ?", status).Count(&count).Error
	return count, err
}