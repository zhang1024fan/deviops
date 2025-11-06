// 验证码 控制层
// author xiaoRui

package controller

import (
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)
// @Tags System系统管理
// @Summary 验证码接口
// @Produce json
// @Description 验证码接口
// @Success 200 {object} result.Result
// @router /api/v1/captcha [get]
func Captcha(c *gin.Context) {
	id, base64Image := service.CaptMake()
	result.Success(c, map[string]interface{}{"idKey": id, "image": base64Image})
}
