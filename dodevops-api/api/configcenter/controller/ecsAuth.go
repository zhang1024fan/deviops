package controller

import (
	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	"dodevops-api/common/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EcsAuthController struct {
	service service.EcsAuthServiceInterface
}

// 分页参数
type PageParams struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

func NewEcsAuthController() *EcsAuthController {
	return &EcsAuthController{
		service: service.GetEcsAuthService(),
	}
}

// GetEcsAuthList 获取所有凭据（分页）
// @Summary 获取所有凭据（分页）
// @Tags Config配置中心
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.EcsAuthVo}}
// @Router /api/v1/config/ecsauthlist [get]
// @Security ApiKeyAuth
func (c *EcsAuthController) GetEcsAuthList(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	c.service.GetEcsAuthListWithPage(ctx, params.Page, params.PageSize)
}

// CreateEcsAuth 创建凭据
// @Summary 创建凭据
// @Param data body model.CreateEcsPasswordAuthDto true "凭据信息"
// @Success 200 {object} result.Result
// @Router /api/v1/config/ecsauthadd [post]
// @Tags Config配置中心
// @Security ApiKeyAuth
func (c *EcsAuthController) CreateEcsAuth(ctx *gin.Context) {
	var dto model.CreateEcsPasswordAuthDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	c.service.CreateEcsAuth(ctx, &dto)
}

// GetEcsAuthByName 根据名称获取凭据
// @Summary 根据名称获取凭据
// @Tags Config配置中心
// @Param name query string true "凭据名称"
// @Success 200 {object} result.Result{data=model.EcsAuthVo}
// @Router /api/v1/config/ecsauthinfo [get]
// @Security ApiKeyAuth
func (c *EcsAuthController) GetEcsAuthByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		result.Failed(ctx, int(result.ApiCode.FAILED), "name参数不能为空")
		return
	}
	c.service.GetEcsAuthByName(ctx, name)
}

// UpdateEcsAuth 更新凭据
// @Summary 更新凭据
// @Tags Config配置中心
// @Param data body model.UpdateEcsAuthDto true "凭据信息"
// @Success 200 {object} result.Result
// @Router /api/v1/config/ecsauthupdate [put]
// @Security ApiKeyAuth
func (c *EcsAuthController) UpdateEcsAuth(ctx *gin.Context) {
	var dto model.UpdateEcsAuthDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	c.service.UpdateEcsAuth(ctx, dto.Id, &dto.CreateEcsPasswordAuthDto)
}

// DeleteEcsAuth 删除凭据
// @Summary 删除凭据
// @Tags Config配置中心
// @Param data body model.EcsAuthIdDto true "凭据ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/ecsauthdelete [delete]
// @Security ApiKeyAuth
func (c *EcsAuthController) DeleteEcsAuth(ctx *gin.Context) {
	var idDto model.EcsAuthIdDto
	if err := ctx.ShouldBindJSON(&idDto); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	c.service.DeleteEcsAuth(ctx, idDto.Id)
}

// GetEcsAuthById 根据ID获取凭据详情
// @Summary 根据ID获取凭据详情
// @Tags Config配置中心
// @Param id query int true "凭据ID"
// @Success 200 {object} result.Result{data=model.EcsAuthVo}
// @Router /api/v1/config/ecsauthdetail [get]
// @Security ApiKeyAuth
func (c *EcsAuthController) GetEcsAuthById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, int(result.ApiCode.FAILED), "id参数不能为空")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "id参数格式错误")
		return
	}
	c.service.GetEcsAuthById(ctx, uint(id))
}
