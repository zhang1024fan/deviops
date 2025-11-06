package service

import (
	"errors"
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
)

type CmdbSQLService struct {
	dao *dao.CmdbSQLDao
}

func NewCmdbSQLService(dao *dao.CmdbSQLDao) *CmdbSQLService {
	return &CmdbSQLService{dao: dao}
}

// CreateDatabase 创建数据库记录
func (s *CmdbSQLService) CreateDatabase(db *model.CmdbSQL) error {
	// 验证类型值
	if db.Type < 1 || db.Type > 5 {
		return errors.New("无效的数据库类型")
	}
	return s.dao.Create(db)
}

// UpdateDatabase 更新数据库记录
func (s *CmdbSQLService) UpdateDatabase(db *model.CmdbSQL) error {
	// 验证类型值
	if db.Type < 1 || db.Type > 5 {
		return errors.New("无效的数据库类型")
	}
	return s.dao.Update(db)
}

// DeleteDatabase 删除数据库记录
func (s *CmdbSQLService) DeleteDatabase(id uint) error {
	return s.dao.Delete(id)
}

// GetDatabase 获取单个数据库详情
func (s *CmdbSQLService) GetDatabase(id uint) (*model.CmdbSQL, error) {
	return s.dao.GetByID(id)
}

// ListDatabases 分页查询数据库列表
func (s *CmdbSQLService) ListDatabases(page, pageSize int) ([]model.CmdbSQL, int64, error) {
	return s.dao.List(page, pageSize)
}

// GetDatabasesByAccount 根据账号查询数据库
func (s *CmdbSQLService) GetDatabasesByAccount(accountID uint) ([]model.CmdbSQL, error) {
	return s.dao.GetByAccountID(accountID)
}

// GetDatabasesByGroup 根据业务组查询数据库
func (s *CmdbSQLService) GetDatabasesByGroup(groupID uint) ([]model.CmdbSQL, error) {
	return s.dao.GetByGroupID(groupID)
}

// GetDatabasesByName 根据名称查询数据库
func (s *CmdbSQLService) GetDatabasesByName(name string) ([]model.CmdbSQL, error) {
	return s.dao.GetByName(name)
}

// GetDatabasesByType 根据类型查询数据库
func (s *CmdbSQLService) GetDatabasesByType(dbType int) ([]model.CmdbSQL, error) {
	if dbType < 1 || dbType > 5 {
		return nil, errors.New("无效的数据库类型")
	}
	return s.dao.GetByType(dbType)
}
