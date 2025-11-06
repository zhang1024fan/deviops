package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type CmdbHostCloudDao struct {
	db *gorm.DB
}

func NewCmdbHostCloudDao() CmdbHostCloudDao {
	return CmdbHostCloudDao{
		db: common.GetDB(),
	}
}

// 创建云主机
func (d *CmdbHostCloudDao) CreateCmdbHost(host *model.CmdbHost) error {
	return d.db.Create(host).Error
}

// 根据ID获取云主机
func (d *CmdbHostCloudDao) GetCmdbHostById(id uint) (model.CmdbHost, error) {
	var host model.CmdbHost
	err := d.db.Where("id = ?", id).First(&host).Error
	return host, err
}
