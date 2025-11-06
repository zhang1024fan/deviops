package dao

import (
	"time"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
)

type KeyManageDao struct{}

func NewKeyManageDao() *KeyManageDao {
	return &KeyManageDao{}
}

// Create 创建密钥
func (d *KeyManageDao) Create(keyManage *model.KeyManage) error {
	if err := keyManage.EncryptKeys(); err != nil {
		return err
	}
	return common.GetDB().Create(keyManage).Error
}

// Update 更新密钥
func (d *KeyManageDao) Update(keyManage *model.KeyManage) error {
	if err := keyManage.EncryptKeys(); err != nil {
		return err
	}
	// 设置当前时间为更新时间
	keyManage.UpdatedAt = util.HTime{Time: time.Now()}
	
	return common.GetDB().Model(keyManage).Updates(map[string]interface{}{
		"key_type":   keyManage.KeyType,
		"key_id":     keyManage.KeyID,
		"key_secret": keyManage.KeySecret,
		"remark":     keyManage.Remark,
		"updated_at": keyManage.UpdatedAt,
	}).Error
}

// Delete 删除密钥
func (d *KeyManageDao) Delete(id uint) error {
	return common.GetDB().Delete(&model.KeyManage{}, id).Error
}

// GetByID 根据ID查询密钥
func (d *KeyManageDao) GetByID(id uint) (*model.KeyManage, error) {
	var keyManage model.KeyManage
	err := common.GetDB().First(&keyManage, id).Error
	return &keyManage, err
}

// List 获取密钥列表
func (d *KeyManageDao) List() ([]model.KeyManage, error) {
	var keyManages []model.KeyManage
	err := common.GetDB().Find(&keyManages).Error
	return keyManages, err
}

// ListWithPage 获取密钥列表（分页）
func (d *KeyManageDao) ListWithPage(page, pageSize int) ([]model.KeyManage, int64, error) {
	var keyManages []model.KeyManage
	var total int64
	
	offset := (page - 1) * pageSize
	
	db := common.GetDB()
	if err := db.Model(&model.KeyManage{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := db.Offset(offset).Limit(pageSize).Find(&keyManages).Error
	return keyManages, total, err
}

// GetByType 根据云厂商类型查询密钥
func (d *KeyManageDao) GetByType(keyType int) ([]model.KeyManage, error) {
	var keyManages []model.KeyManage
	err := common.GetDB().Where("key_type = ?", keyType).Find(&keyManages).Error
	return keyManages, err
}