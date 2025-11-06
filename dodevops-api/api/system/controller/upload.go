// 图片上传 控制层
// author xiaoRui

package controller

import (
	"dodevops-api/api/system/service"

	"github.com/gin-gonic/gin"
)

// @Tags System系统管理
// 单图片上传
// @Summary 单图片上传接口
// @Description 单图片上传接口
// @Produce json
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {object} result.Result
// @Router /api/v1/upload [post]
// @Security ApiKeyAuth
func Upload(c *gin.Context) {
	service.UploadService().Upload(c)
}
