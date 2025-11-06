package controller

import (
	"fmt"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type KeyManageController struct {
	service *service.KeyManageService
}

func NewKeyManageController() *KeyManageController {
	return &KeyManageController{
		service: service.NewKeyManageService(),
	}
}

// Create 创建密钥
// @Summary 创建密钥管理信息
// @Tags Config配置中心
// @Param keyManage body model.CreateKeyManageDto true "密钥管理信息"
// @Success 200 {object} result.Result{data=model.KeyManage}
// @Router /api/v1/config/keymanage [post]
// @Security ApiKeyAuth
func (c *KeyManageController) Create(ctx *gin.Context) {
	var dto model.CreateKeyManageDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	keyManage := model.KeyManage{
		KeyType:   dto.KeyType,
		KeyID:     dto.KeyID,
		KeySecret: dto.KeySecret,
		Remark:    dto.Remark,
	}

	if err := c.service.Create(&keyManage); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	result.Success(ctx, keyManage)
}

// Update 更新密钥
// @Summary 更新密钥管理信息
// @Tags Config配置中心
// @Param keyManage body model.UpdateKeyManageDto true "密钥管理信息"
// @Success 200 {object} result.Result{data=model.KeyManage}
// @Router /api/v1/config/keymanage [put]
// @Security ApiKeyAuth
func (c *KeyManageController) Update(ctx *gin.Context) {
	var dto model.UpdateKeyManageDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	keyManage := model.KeyManage{
		ID:        dto.ID,
		KeyType:   dto.KeyType,
		KeyID:     dto.KeyID,
		KeySecret: dto.KeySecret,
		Remark:    dto.Remark,
	}

	if err := c.service.Update(&keyManage); err != nil {
		result.Failed(ctx, 500, "更新失败: "+err.Error())
		return
	}

	result.Success(ctx, keyManage)
}

// Delete 删除密钥
// @Summary 删除密钥管理信息
// @Tags Config配置中心
// @Param id query uint true "密钥ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/keymanage [delete]
// @Security ApiKeyAuth
func (c *KeyManageController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	if err := c.service.Delete(req.ID); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// GetByID 根据ID获取密钥详情
// @Summary 根据ID获取密钥详情
// @Tags Config配置中心
// @Param id query uint true "密钥ID"
// @Success 200 {object} result.Result{data=model.KeyManage}
// @Router /api/v1/config/keymanage [get]
// @Security ApiKeyAuth
func (c *KeyManageController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	keyManage, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, keyManage)
}

// List 获取密钥列表（分页）
// @Summary 获取密钥列表（分页）
// @Tags Config配置中心
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.KeyManage}}
// @Router /api/v1/config/keymanage/list [get]
// @Security ApiKeyAuth
func (c *KeyManageController) List(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	keyManages, total, err := c.service.ListWithPage(params.Page, params.PageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	pageResult := result.PageResult{
		List:     keyManages,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	result.Success(ctx, pageResult)
}

// DecryptKeys 解密密钥信息
// @Summary 解密密钥信息
// @Tags Config配置中心
// @Param id query uint true "密钥ID"
// @Success 200 {object} result.Result{data=map[string]string}
// @Router /api/v1/config/keymanage/decrypt [post]
// @Security ApiKeyAuth
func (c *KeyManageController) DecryptKeys(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	keyID, keySecret, err := c.service.DecryptKeys(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "解密失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"keyId":     keyID,
		"keySecret": keySecret,
	})
}

// GetByType 根据云厂商类型查询密钥
// @Summary 根据云厂商类型查询密钥
// @Tags Config配置中心
// @Param type query int true "云厂商类型"
// @Success 200 {object} result.Result{data=[]model.KeyManage}
// @Router /api/v1/config/keymanage/type [get]
// @Security ApiKeyAuth
func (c *KeyManageController) GetByType(ctx *gin.Context) {
	var req struct {
		Type int `form:"type" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	keyManages, err := c.service.GetByType(req.Type)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}
	result.Success(ctx, keyManages)
}

// SyncCloudHosts 同步云主机（统一接口）
// @Summary 根据密钥类型自动同步对应云厂商的主机
// @Tags Config配置中心
// @Param request body object{keyId=uint,keyType=int} true "同步参数，keyType: 1=阿里云,2=腾讯云,3=百度云"
// @Success 200 {object} result.Result
// @Router /api/v1/config/keymanage/sync [post]
// @Security ApiKeyAuth
func (c *KeyManageController) SyncCloudHosts(ctx *gin.Context) {
	var req struct {
		KeyID   uint `json:"keyId" binding:"required"`
		KeyType int  `json:"keyType" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 固定使用默认分组ID为1
	groupID := uint(1)
	// 固定使用全区域扫描
	region := "all"

	// 验证密钥是否存在
	keyManage, err := c.service.GetByID(req.KeyID)
	if err != nil {
		result.Failed(ctx, 500, "密钥不存在: "+err.Error())
		return
	}

	// 验证前端传递的keyType与数据库中的keyType是否一致
	if keyManage.KeyType != req.KeyType {
		result.Failed(ctx, 400, fmt.Sprintf("密钥类型不匹配，数据库中为%d，传递的为%d", keyManage.KeyType, req.KeyType))
		return
	}

	// 根据keyType调用不同的同步方法
	switch req.KeyType {
	case 1: // 阿里云
		c.service.SyncAliyunHosts(ctx, req.KeyID, groupID, region)
	case 2: // 腾讯云
		c.service.SyncTencentHosts(ctx, req.KeyID, groupID)
	case 3: // 百度云
		c.service.SyncBaiduHosts(ctx, req.KeyID, groupID)
	default:
		result.Failed(ctx, 400, "不支持的云厂商类型，目前支持：1=阿里云，2=腾讯云，3=百度云")
	}
}