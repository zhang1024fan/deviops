package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"time"

	"github.com/gin-gonic/gin"
)

// 接口名称修改为 CmdbGroupServiceInterface 避免冲突
type CmdbGroupServiceInterface interface {
	CreateCmdbGroup(c *gin.Context, group model.CmdbGroup) // 创建分组
	GetAllCmdbGroups(c *gin.Context)                       // 获取所有分组
	GetAllCmdbGroupsWithHosts(c *gin.Context)              // 获取所有分组及关联主机
	UpdateCmdbGroup(c *gin.Context, group model.CmdbGroup) // 更新分组
	DeleteCmdbGroup(c *gin.Context, id uint)               // 删除分组
	GetCmdbGroupByName(c *gin.Context, name string)        // 根据名称查询分组
}

type CmdbGroupServiceImpl struct{}

// 新增分组
func (s CmdbGroupServiceImpl) CreateCmdbGroup(c *gin.Context, group model.CmdbGroup) {
	dao := dao.NewCmdbGroupDao()
	if dao.CheckNameExists(group.Name) {
		result.FailedWithCode(c, constant.GROUP_EXIST, "分组已存在无法创建")
		return
	}
	group.CreateTime = util.HTime{Time: time.Now()}
	err := dao.CreateCmdbGroup(&group)
	if err != nil {
		result.Failed(c, constant.GROUP_EXIST, "创建分组失败")
		return
	}
	result.Success(c, true)
}

// 查询所有分组并返回树形结构
func (s CmdbGroupServiceImpl) GetAllCmdbGroups(c *gin.Context) {
	groupDao := dao.NewCmdbGroupDao()
	hostDao := dao.NewCmdbHostDao()
	
	// 获取所有分组
	groups := groupDao.GetCmdbGroupList()
	// 获取所有主机
	hosts := hostDao.GetCmdbHostList()
	
	result.Success(c, model.BuildCmdbGroupTreeWithHostCount(groups, hosts))
}

// 查询所有分组及关联主机(树形结构)
func (s CmdbGroupServiceImpl) GetAllCmdbGroupsWithHosts(c *gin.Context) {
	groupDao := dao.NewCmdbGroupDao()
	hostDao := dao.NewCmdbHostDao()

	// 获取所有分组
	groups := groupDao.GetCmdbGroupList()
	// 获取所有主机并转换为VO
	hosts := hostDao.GetCmdbHostList()
	var hostVos []model.CmdbHostVo
	for _, host := range hosts {
		hostVos = append(hostVos, model.CmdbHostVo{
			ID:         host.ID,
			HostName:   host.HostName,
			Name:       host.Name,
			GroupID:    host.GroupID,
			PrivateIP:  host.PrivateIP,
			PublicIP:   host.PublicIP,
			SSHIP:      host.SSHIP,
			SSHName:    host.SSHName,
			SSHKeyID:   host.SSHKeyID,
			SSHPort:    host.SSHPort,
			Remark:     host.Remark,
			Vendor:     getVendorName(host.Vendor),
			Region:     host.Region,
			InstanceID: host.InstanceID,
			OS:         host.OS,
			Status:     host.Status,
			CPU:        host.CPU,
			Memory:     host.Memory,
			Disk:       host.Disk,
			BillingType: host.BillingType,
			CreateTime: host.CreateTime,
			ExpireTime: host.ExpireTime,
			UpdateTime: host.UpdateTime,
		})
	}

	result.Success(c, model.BuildCmdbGroupHostTree(groups, hostVos))
}

// 获取云厂商名称
func getVendorName(vendor int) string {
	switch vendor {
	case 1:
		return "自建"
	case 2:
		return "阿里云"
	case 3:
		return "腾讯云"
	default:
		return "未知"
	}
}

// 更新分组
func (s CmdbGroupServiceImpl) UpdateCmdbGroup(c *gin.Context, group model.CmdbGroup) {
	dao := dao.NewCmdbGroupDao()
	err := dao.UpdateCmdbGroup(group.ID, &group)
	if err != nil {
		result.FailedWithCode(c, constant.GROUP_EXIST, err.Error())
		return
	}
	result.Success(c, true)
}

// 删除分组
func (s CmdbGroupServiceImpl) DeleteCmdbGroup(c *gin.Context, id uint) {
	dao := dao.NewCmdbGroupDao()
	err := dao.DeleteCmdbGroup(id)
	if err != nil {
		result.FailedWithCode(c, constant.GROUP_EXIST, err.Error())
		return
	}
	result.Success(c, true)
}

// 根据名称查询分组
func (s CmdbGroupServiceImpl) GetCmdbGroupByName(c *gin.Context, name string) {
	dao := dao.NewCmdbGroupDao()
	group, err := dao.GetCmdbGroupByName(name)
	if err != nil {
		result.Failed(c, constant.GROUP_EXIST, "查询分组失败")
		return
	}
	result.Success(c, group)
}

// 全局服务调用方法
func GetCmdbGroupService() CmdbGroupServiceInterface {
	return &CmdbGroupServiceImpl{}
}
