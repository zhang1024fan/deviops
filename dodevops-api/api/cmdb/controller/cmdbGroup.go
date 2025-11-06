package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

var sysDept model.CmdbGroup

// @Summary 新增资产分组接口
// @Produce json
// @Tags CMDB资产管理
// @Description 新增资产分组接口
// @Param data body model.CmdbGroup true "data"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/groupadd [post]
// @Security ApiKeyAuth
func CreateCmdbGroup(c *gin.Context) {
	_ = c.BindJSON(&sysDept)
	service.GetCmdbGroupService().CreateCmdbGroup(c, sysDept)
}

// @Summary 查询所有资产分组（树形结构）
// @Produce json
// @Tags CMDB资产管理
// @Description 查询所有资产分组，并以树形结构返回
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/grouplist [get]
// @Security ApiKeyAuth
func GetAllCmdbGroups(c *gin.Context) {
	service.GetCmdbGroupService().GetAllCmdbGroups(c)
}

// @Summary 查询所有资产分组及关联主机（树形结构）
// @Produce json
// @Tags CMDB资产管理
// @Description 查询所有资产分组及关联主机，并以树形结构返回
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/grouplistwithhosts [get]
// @Security ApiKeyAuth
func GetAllCmdbGroupsWithHosts(c *gin.Context) {
	service.GetCmdbGroupService().GetAllCmdbGroupsWithHosts(c)
}

// @Summary 更新资产分组接口
// @Produce json
// @Tags CMDB资产管理
// @Description 更新资产分组接口
// @Param data body model.CmdbGroup true "data"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/groupupdate [put]
// @Security ApiKeyAuth
func UpdateCmdbGroup(c *gin.Context) {
	var group model.CmdbGroup
	_ = c.BindJSON(&group)
	service.GetCmdbGroupService().UpdateCmdbGroup(c, group)
}

// @Summary 删除资产分组接口
// @Produce json
// @Tags CMDB资产管理
// @Description 删除资产分组接口
// @Param data body model.CmdbGroupIdDto true "分组ID"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/groupdelete [delete]
// @Security ApiKeyAuth
func DeleteCmdbGroup(c *gin.Context) {
	var dto model.CmdbGroupIdDto
	if err := c.BindJSON(&dto); err != nil {
		result.Failed(c, constant.GROUP_EXIST, "参数错误")
		return
	}
	service.GetCmdbGroupService().DeleteCmdbGroup(c, dto.Id)
}

// @Summary 根据名称查询资产分组
// @Produce json
// @Tags CMDB资产管理
// @Description 根据名称查询资产分组
// @Param name query string true "分组名称"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/groupbyname [get]
// @Security ApiKeyAuth
func GetCmdbGroupByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		result.Failed(c, constant.GROUP_EXIST, "分组名称不能为空")
		return
	}
	service.GetCmdbGroupService().GetCmdbGroupByName(c, name)
}
