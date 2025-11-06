package dao

import (
	"errors"
	cmdbModel "dodevops-api/api/cmdb/model"
	configModel "dodevops-api/api/configcenter/model"
	"dodevops-api/common"
)

type CmdbHostSSHDao struct{}

func NewCmdbHostSSHDao() *CmdbHostSSHDao {
	return &CmdbHostSSHDao{}
}

func (d *CmdbHostSSHDao) GetHostSSHInfo(hostID uint) (*cmdbModel.CmdbHost, error) {
	var host cmdbModel.CmdbHost
	if err := common.GetDB().Where("id = ?", hostID).First(&host).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

func (d *CmdbHostSSHDao) GetSSHCredentials(keyID uint) (string, error) {
	var auth configModel.EcsAuth
	if err := common.GetDB().Table("config_ecsauth").Where("id = ?", keyID).First(&auth).Error; err != nil {
		return "", err
	}

	switch auth.Type {
	case 1: // 密码认证
		return auth.Password, nil
	case 2: // 密钥认证
		return auth.PublicKey, nil
	case 3: // 公钥免认证
		return "", nil // Type 3 不需要返回凭据，由系统自动查找本地私钥
	default:
		return "", errors.New("不支持的认证方式")
	}
}
