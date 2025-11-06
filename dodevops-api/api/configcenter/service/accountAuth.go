package service

import (
	"fmt"
	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
)

type AccountAuthService struct {
	dao *dao.AccountAuthDao
}

func NewAccountAuthService() *AccountAuthService {
	return &AccountAuthService{
		dao: dao.NewAccountAuthDao(),
	}
}

// Create 创建账号
func (s *AccountAuthService) Create(account *model.AccountAuth) error {
	// 检查alias是否已存在
	existing, err := s.GetByAlias(account.Alias)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("alias '%s' 已存在", account.Alias)
	}
	return s.dao.Create(account)
}

// Update 更新账号
func (s *AccountAuthService) Update(account *model.AccountAuth) error {
	return s.dao.Update(account)
}

// Delete 删除账号
func (s *AccountAuthService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID 根据ID查询账号
func (s *AccountAuthService) GetByID(id uint) (*model.AccountAuth, error) {
	return s.dao.GetByID(id)
}

// List 获取账号列表
func (s *AccountAuthService) List() ([]model.AccountAuth, error) {
	return s.dao.List()
}

// ListWithPage 获取账号列表（分页）
func (s *AccountAuthService) ListWithPage(page, pageSize int) ([]model.AccountAuth, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// DecryptPassword 解密密码
func (s *AccountAuthService) DecryptPassword(id uint) (string, error) {
	account, err := s.dao.GetByID(id)
	if err != nil {
		return "", err
	}
	return account.DecryptPassword()
}

// GetByType 根据类型查询账号
func (s *AccountAuthService) GetByType(accountType string) ([]model.AccountAuth, error) {
	return s.dao.GetByType(accountType)
}

// GetByAlias 根据别名查询账号
func (s *AccountAuthService) GetByAlias(alias string) (*model.AccountAuth, error) {
	return s.dao.GetByAlias(alias)
}
