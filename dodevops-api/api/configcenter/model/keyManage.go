package model

import (
	"dodevops-api/common/util"
)

type KeyManage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	KeyType   int       `gorm:"not null" json:"keyType"`             // 云厂商类型：1=阿里云，2=腾讯云，3=百度云，4=华为云，5=AWS云
	KeyID     string    `gorm:"type:text;not null" json:"keyId"`     // 密钥ID(加密存储)
	KeySecret string    `gorm:"type:text;not null" json:"keySecret"` // 密钥Secret(加密存储)
	Remark    string    `gorm:"type:text" json:"remark"`             // 备注信息
	CreatedAt util.HTime `json:"createdAt"`                          // 创建时间
	UpdatedAt util.HTime `json:"updatedAt"`                          // 更新时间
}

// TableName 表名
func (KeyManage) TableName() string {
	return "config_keymanage"
}

// EncryptKeys 加密密钥信息
func (k *KeyManage) EncryptKeys() error {
	encryptedID, err := util.AESEncrypt(k.KeyID)
	if err != nil {
		return err
	}
	encryptedSecret, err := util.AESEncrypt(k.KeySecret)
	if err != nil {
		return err
	}
	k.KeyID = encryptedID
	k.KeySecret = encryptedSecret
	return nil
}

// DecryptKeys 解密密钥信息
func (k *KeyManage) DecryptKeys() (string, string, error) {
	keyID, err := util.AESDecrypt(k.KeyID)
	if err != nil {
		return "", "", err
	}
	keySecret, err := util.AESDecrypt(k.KeySecret)
	if err != nil {
		return "", "", err
	}
	return keyID, keySecret, nil
}

// CreateKeyManageDto 创建密钥DTO
type CreateKeyManageDto struct {
	KeyType   int    `json:"keyType" binding:"required"`   // 云厂商类型：1=阿里云，2=腾讯云，3=百度云，4=华为云，5=AWS云
	KeyID     string `json:"keyId" binding:"required"`     // 密钥ID
	KeySecret string `json:"keySecret" binding:"required"` // 密钥Secret
	Remark    string `json:"remark"`                       // 备注信息
}

// UpdateKeyManageDto 更新密钥DTO
type UpdateKeyManageDto struct {
	ID        uint   `json:"id" binding:"required"`        // 密钥ID
	KeyType   int    `json:"keyType" binding:"required"`   // 云厂商类型：1=阿里云，2=腾讯云，3=百度云，4=华为云，5=AWS云
	KeyID     string `json:"keyId" binding:"required"`     // 密钥ID
	KeySecret string `json:"keySecret" binding:"required"` // 密钥Secret
	Remark    string `json:"remark"`                       // 备注信息
}