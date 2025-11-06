package dao

import (
	"dodevops-api/api/cmdb/model"
	"gorm.io/gorm"
)

type CmdbSQLDao struct {
	db *gorm.DB
}

func NewCmdbSQLDao(db *gorm.DB) *CmdbSQLDao {
	return &CmdbSQLDao{db: db}
}

// Create 创建数据库记录
func (d *CmdbSQLDao) Create(db *model.CmdbSQL) error {
	return d.db.Create(db).Error
}

// Update 更新数据库记录
func (d *CmdbSQLDao) Update(db *model.CmdbSQL) error {
	return d.db.Save(db).Error
}

// Delete 删除数据库记录
func (d *CmdbSQLDao) Delete(id uint) error {
	return d.db.Delete(&model.CmdbSQL{}, id).Error
}

// GetByID 根据ID获取数据库记录
func (d *CmdbSQLDao) GetByID(id uint) (*model.CmdbSQL, error) {
	var db model.CmdbSQL
	err := d.db.First(&db, id).Error
	return &db, err
}

// List 分页查询数据库列表
func (d *CmdbSQLDao) List(page, pageSize int) ([]model.CmdbSQL, int64, error) {
	var dbs []model.CmdbSQL
	var count int64

	if err := d.db.Model(&model.CmdbSQL{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := d.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&dbs).Error
	return dbs, count, err
}

// GetByAccountID 根据账号ID查询数据库
func (d *CmdbSQLDao) GetByAccountID(accountID uint) ([]model.CmdbSQL, error) {
	var dbs []model.CmdbSQL
	err := d.db.Where("account_id = ?", accountID).Find(&dbs).Error
	return dbs, err
}

// GetByGroupID 根据业务组ID查询数据库
func (d *CmdbSQLDao) GetByGroupID(groupID uint) ([]model.CmdbSQL, error) {
	var dbs []model.CmdbSQL
	err := d.db.Where("group_id = ?", groupID).Find(&dbs).Error
	return dbs, err
}

// GetByName 根据名称查询数据库
func (d *CmdbSQLDao) GetByName(name string) ([]model.CmdbSQL, error) {
	var dbs []model.CmdbSQL
	err := d.db.Where("name = ?", name).Find(&dbs).Error
	return dbs, err
}

// GetByType 根据类型查询数据库
func (d *CmdbSQLDao) GetByType(dbType int) ([]model.CmdbSQL, error) {
	var dbs []model.CmdbSQL
	err := d.db.Where("type = ?", dbType).Find(&dbs).Error
	return dbs, err
}
