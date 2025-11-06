package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type CmdbHostCloudController struct {
	service service.CmdbHostCloudServiceInterface
}

func NewCmdbHostCloudController() *CmdbHostCloudController {
	return &CmdbHostCloudController{
		service: service.GetCmdbHostCloudService(),
	}
}

// 创建阿里云主机
// @Summary 创建阿里云主机
// @Description 创建阿里云主机(通过阿里云API获取主机信息)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CreateCmdbHostCloudDto true "阿里云主机信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostcloudcreatealiyun [post]
// @Security ApiKeyAuth
func (c *CmdbHostCloudController) CreateAliyunHost(ctx *gin.Context) {
	var dto model.CreateCmdbHostCloudDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.CreateAliyunHost(ctx, &dto)
}

// 创建腾讯云主机
// @Summary 创建腾讯云主机
// @Description 创建腾讯云主机(通过腾讯云API获取主机信息)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CreateCmdbHostCloudDto true "腾讯云主机信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostcloudcreatetencent [post]
// @Security ApiKeyAuth
func (c *CmdbHostCloudController) CreateTencentHost(ctx *gin.Context) {
	var dto model.CreateCmdbHostCloudDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.CreateTencentHost(ctx, &dto)
}

// 创建百度云主机
// @Summary 创建百度云主机
// @Description 创建百度云主机(通过百度云API自动扫描所有区域并获取主机信息)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CreateCmdbHostCloudDto true "百度云主机信息(AccessKey和SecretKey)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostcloudcreatebaidu [post]
// @Security ApiKeyAuth
func (c *CmdbHostCloudController) CreateBaiduHost(ctx *gin.Context) {
	var dto model.CreateCmdbHostCloudDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.CreateBaiduHost(ctx, &dto)
}

