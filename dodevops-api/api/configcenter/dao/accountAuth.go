package dao

import (
	"time"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
)

type AccountAuthDao struct{}

func NewAccountAuthDao() *AccountAuthDao {
	return &AccountAuthDao{}
}

// Create 创建账号
func (d *AccountAuthDao) Create(account *model.AccountAuth) error {
	if err := account.EncryptPassword(); err != nil {
		return err
	}
	return common.GetDB().Create(account).Error
}

// Update 更新账号
func (d *AccountAuthDao) Update(account *model.AccountAuth) error {
	if err := account.EncryptPassword(); err != nil {
		return err
	}
	// 设置当前时间为更新时间
	account.UpdatedAt = util.HTime{Time: time.Now()}
	
	return common.GetDB().Model(account).Updates(map[string]interface{}{
		"alias":     account.Alias,
		"host":      account.Host,
		"port":      account.Port,
		"name":      account.Name,
		"password":  account.Password,
		"type":      account.Type,
		"remark":    account.Remark,
		"updated_at": account.UpdatedAt,
	}).Error
}

// Delete 删除账号
func (d *AccountAuthDao) Delete(id uint) error {
	return common.GetDB().Delete(&model.AccountAuth{}, id).Error
}

// GetByID 根据ID查询账号
func (d *AccountAuthDao) GetByID(id uint) (*model.AccountAuth, error) {
	var account model.AccountAuth
	err := common.GetDB().First(&account, id).Error
	return &account, err
}
// List 获取账号列表
func (d *AccountAuthDao) List() ([]model.AccountAuth, error) {
	var accounts []model.AccountAuth
	err := common.GetDB().Find(&accounts).Error
	return accounts, err
}

// ListWithPage 获取账号列表（分页）
func (d *AccountAuthDao) ListWithPage(page, pageSize int) ([]model.AccountAuth, int64, error) {
	var accounts []model.AccountAuth
	var total int64
	
	offset := (page - 1) * pageSize
	
	db := common.GetDB()
	if err := db.Model(&model.AccountAuth{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := db.Offset(offset).Limit(pageSize).Find(&accounts).Error
	return accounts, total, err
}
// GetByType 根据类型查询账号
func (d *AccountAuthDao) GetByType(accountType string) ([]model.AccountAuth, error) {
	var accounts []model.AccountAuth
	err := common.GetDB().Where("type = ?", accountType).Find(&accounts).Error
	return accounts, err
}

// GetByAlias 根据别名查询账号
func (d *AccountAuthDao) GetByAlias(alias string) (*model.AccountAuth, error) {
	var account model.AccountAuth
	err := common.GetDB().Where("alias = ?", alias).First(&account).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// GetAccountsByType 根据类型获取账号列表（分页）
func (d *AccountAuthDao) GetAccountsByType(accountType int, page, pageSize int) ([]model.AccountAuth, int64, error) {
	var accounts []model.AccountAuth
	var total int64

	offset := (page - 1) * pageSize

	db := common.GetDB()
	if err := db.Model(&model.AccountAuth{}).Where("type = ?", accountType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Where("type = ?", accountType).Offset(offset).Limit(pageSize).Find(&accounts).Error
	return accounts, total, err
}
