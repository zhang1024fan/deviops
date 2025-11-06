package model

import (
	"dodevops-api/common/util"
)

type AccountAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Alias     string    `gorm:"size:128;not null" json:"alias"`      // 别名
	Host      string    `gorm:"size:128;not null" json:"host"`       // IP地址
	Port      int       `gorm:"not null" json:"port"`                // 端口
	Name      string    `gorm:"size:128;not null" json:"name"`       // 用户名
	Password  string    `gorm:"type:text;not null" json:"password"`  // 密码(加密存储)
	Type      int       `gorm:"not null" json:"type"`                // 类型
	Remark    string    `gorm:"type:text" json:"remark"`             // 备注
	CreatedAt util.HTime `json:"createdAt"`                          // 创建时间
	UpdatedAt util.HTime `json:"updatedAt"`                          // 更新时间
}
// 表名
func (AccountAuth) TableName() string {
	return "config_account"
}

// EncryptPassword 加密密码
func (a *AccountAuth) EncryptPassword() error {
	encrypted, err := util.AESEncrypt(a.Password)
	if err != nil {
		return err
	}
	a.Password = encrypted
	return nil
}

// DecryptPassword 解密密码
func (a *AccountAuth) DecryptPassword() (string, error) {
	return util.AESDecrypt(a.Password)
}
