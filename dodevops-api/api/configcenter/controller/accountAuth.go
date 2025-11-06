package controller

import (
	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type AccountAuthController struct {
	service *service.AccountAuthService
}

func NewAccountAuthController() *AccountAuthController {
	return &AccountAuthController{
		service: service.NewAccountAuthService(),
	}
}

// Create 创建账号
// @Summary 创建账号认证信息
// @Tags Config配置中心
// @Param account body model.AccountAuth true "账号认证信息"
// @Success 200 {object} result.Result{data=model.AccountAuth}
// @Router /api/v1/config/accountauth [post]
// @Security ApiKeyAuth
func (c *AccountAuthController) Create(ctx *gin.Context) {
	var account model.AccountAuth
	if err := ctx.ShouldBindJSON(&account); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.Create(&account); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	result.Success(ctx, account)
}

// Update 更新账号
// @Summary 更新账号认证信息
// @Tags Config配置中心
// @Param account body model.AccountAuth true "账号认证信息"
// @Success 200 {object} result.Result{data=model.AccountAuth}
// @Router /api/v1/config/accountauth [put]
// @Security ApiKeyAuth
func (c *AccountAuthController) Update(ctx *gin.Context) {
	var account model.AccountAuth
	if err := ctx.ShouldBindJSON(&account); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.Update(&account); err != nil {
		result.Failed(ctx, 500, "更新失败: "+err.Error())
		return
	}

	result.Success(ctx, account)
}

// Delete 删除账号
// @Summary 删除账号认证信息
// @Tags Config配置中心
// @Param id query uint true "账号ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/accountauth [delete]
// @Security ApiKeyAuth
func (c *AccountAuthController) Delete(ctx *gin.Context) {
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

// GetByID 根据ID获取账号详情
// @Summary 根据ID获取账号详情
// @Tags Config配置中心
// @Param id query uint true "账号ID"
// @Success 200 {object} result.Result{data=model.AccountAuth}
// @Router /api/v1/config/accountauth [get]
// @Security ApiKeyAuth
func (c *AccountAuthController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	account, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, account)
}

// List 获取账号列表（分页）
// @Summary 获取账号列表（分页）
// @Tags Config配置中心
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.AccountAuth}}
// @Router /api/v1/config/accountauth/list [get]
// @Security ApiKeyAuth
func (c *AccountAuthController) List(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	accounts, total, err := c.service.ListWithPage(params.Page, params.PageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	pageResult := result.PageResult{
		List:     accounts,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	result.Success(ctx, pageResult)
}

// DecryptPassword 解密密码
// @Summary 解密密码
// @Tags Config配置中心
// @Param id query uint true "账号ID"
// @Success 200 {object} result.Result{data=string}
// @Router /api/v1/config/accountauth/decrypt [post]
// @Security ApiKeyAuth
func (c *AccountAuthController) DecryptPassword(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	password, err := c.service.DecryptPassword(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "解密失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"password": password,
	})
}

// GetByType 根据类型查询账号
// @Summary 根据类型查询账号
// @Tags Config配置中心
// @Param type query string true "账号类型"
// @Success 200 {object} result.Result{data=[]model.AccountAuth}
// @Router /api/v1/config/accountauth/type [get]
// @Security ApiKeyAuth
func (c *AccountAuthController) GetByType(ctx *gin.Context) {
	accountType := ctx.Query("type")
	accounts, err := c.service.GetByType(accountType)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}
	result.Success(ctx, accounts)
}

// GetByAlias 根据别名查询账号
// @Summary 根据别名查询账号
// @Tags Config配置中心
// @Param alias query string true "账号别名"
// @Success 200 {object} result.Result{data=model.AccountAuth}
// @Router /api/v1/config/accountauth/alias [get]
// @Security ApiKeyAuth
func (c *AccountAuthController) GetByAlias(ctx *gin.Context) {
	alias := ctx.Query("alias")
	account, err := c.service.GetByAlias(alias)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}
	result.Success(ctx, account)
}
