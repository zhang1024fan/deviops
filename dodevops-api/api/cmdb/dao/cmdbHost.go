package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"
	"time"

	"gorm.io/gorm"
)

const (
	hostCachePrefix = "cmdb:host:"
	cacheExpireTime = 5 * time.Minute
)

type CmdbHostDao struct {
	db *gorm.DB
}

func NewCmdbHostDao() CmdbHostDao {
	return CmdbHostDao{
		db: common.GetDB(),
	}
}

// 获取主机列表(分页)
func (d *CmdbHostDao) GetCmdbHostListWithPage(page, pageSize int) ([]model.CmdbHost, int64) {
	var list []model.CmdbHost
	var total int64
	
	// 使用预加载Group信息并优化查询
	db := d.db.Preload("Group")
	db.Model(&model.CmdbHost{}).Count(&total)
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	
	return list, total
}

// 获取主机列表
func (d *CmdbHostDao) GetCmdbHostList() []model.CmdbHost {
	var list []model.CmdbHost
	d.db.Find(&list)
	return list
}

// 根据ID获取主机
func (d *CmdbHostDao) GetCmdbHostById(id uint) (model.CmdbHost, error) {
	var host model.CmdbHost
	err := d.db.Where("id = ?", id).First(&host).Error
	return host, err
}

// 根据名称获取主机
func (d *CmdbHostDao) GetCmdbHostByName(name string) (model.CmdbHost, error) {
	var host model.CmdbHost
	err := d.db.Where("host_name = ?", name).First(&host).Error
	return host, err
}

// 根据IP获取单个主机
func (d *CmdbHostDao) GetCmdbHostByIP(ip string) (model.CmdbHost, error) {
	var host model.CmdbHost
	err := d.db.Preload("Group").Where("private_ip = ? OR public_ip = ? OR ssh_ip = ?", ip, ip, ip).First(&host).Error
	return host, err
}

// 检查主机名称是否存在
func (d *CmdbHostDao) CheckNameExists(name string) bool {
	var count int64
	d.db.Model(&model.CmdbHost{}).Where("host_name = ?", name).Count(&count)
	return count > 0
}

// 创建主机
func (d *CmdbHostDao) CreateCmdbHost(host *model.CmdbHost) error {
	return d.db.Create(host).Error
}

// 更新主机
func (d *CmdbHostDao) UpdateCmdbHost(id uint, host *model.CmdbHost) error {
	return d.db.Model(&model.CmdbHost{}).Where("id = ?", id).Updates(host).Error
}

// 删除主机
func (d *CmdbHostDao) DeleteCmdbHost(id uint) error {
	return d.db.Delete(&model.CmdbHost{}, id).Error
}

// 根据实例ID获取主机
func (d *CmdbHostDao) GetCmdbHostByInstanceID(instanceID string) (*model.CmdbHost, error) {
	var host model.CmdbHost
	err := d.db.Where("instance_id = ?", instanceID).First(&host).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

// 根据SSH IP获取主机
func (d *CmdbHostDao) GetCmdbHostBySSHIP(sshIP string) *model.CmdbHost {
	var host model.CmdbHost
	err := d.db.Where("ssh_ip = ?", sshIP).First(&host).Error
	if err != nil {
		return nil
	}
	return &host
}

// 递归获取所有子分组ID（包括自身）
func (d *CmdbHostDao) getAllChildGroupIds(groupId uint) []uint {
	var groupIds []uint
	groupIds = append(groupIds, groupId) // 包含自身

	var childGroups []model.CmdbGroup
	d.db.Where("parent_id = ?", groupId).Find(&childGroups)

	for _, child := range childGroups {
		childIds := d.getAllChildGroupIds(child.ID)
		groupIds = append(groupIds, childIds...)
	}

	return groupIds
}

// 根据分组ID获取主机列表（包括所有子分组的主机，支持分页）
func (d *CmdbHostDao) GetCmdbHostsByGroupId(groupId uint) []model.CmdbHost {
	var list []model.CmdbHost
	// 获取当前分组及所有子分组的ID
	groupIds := d.getAllChildGroupIds(groupId)
	d.db.Preload("Group").Where("group_id IN ?", groupIds).Find(&list)
	return list
}

// 根据分组ID获取主机列表（包括所有子分组的主机，支持分页）
func (d *CmdbHostDao) GetCmdbHostsByGroupIdWithPage(groupId uint, page, pageSize int) ([]model.CmdbHost, int64) {
	var list []model.CmdbHost
	var count int64
	
	// 获取当前分组及所有子分组的ID
	groupIds := d.getAllChildGroupIds(groupId)
	
	// 计算总数
	d.db.Model(&model.CmdbHost{}).Where("group_id IN ?", groupIds).Count(&count)
	
	// 分页查询
	offset := (page - 1) * pageSize
	d.db.Preload("Group").Where("group_id IN ?", groupIds).Offset(offset).Limit(pageSize).Find(&list)
	
	return list, count
}

// 根据主机名称模糊查询
func (d *CmdbHostDao) GetCmdbHostsByHostNameLike(name string) []model.CmdbHost {
	var list []model.CmdbHost
	d.db.Preload("Group").Where("host_name LIKE ?", "%"+name+"%").Find(&list)
	return list
}

// 根据IP查询(匹配内网IP、公网IP或SSH IP)
func (d *CmdbHostDao) GetCmdbHostsByIP(ip string) []model.CmdbHost {
	var list []model.CmdbHost
	d.db.Preload("Group").Where("private_ip = ? OR public_ip = ? OR ssh_ip = ?", ip, ip, ip).Find(&list)
	return list
}

// 根据状态查询
func (d *CmdbHostDao) GetCmdbHostsByStatus(status int) []model.CmdbHost {
	var list []model.CmdbHost
	d.db.Preload("Group").Where("status = ?", status).Find(&list)
	return list
}

// 根据ID列表批量获取主机
func (d *CmdbHostDao) GetCmdbHostsByIds(ids []uint) ([]model.CmdbHost, error) {
	var hosts []model.CmdbHost
	err := d.db.Preload("Group").Where("id IN ?", ids).Find(&hosts).Error
	return hosts, err
}
