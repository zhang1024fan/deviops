// 图片上传 服务层
// author xiaoRui

package service

import (
	"fmt"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/log"
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"time"
)

type IUploadService interface {
	Upload(c *gin.Context)
}

type UploadServiceImpl struct{}

// 图片上传
func (u UploadServiceImpl) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), result.ApiCode.GetMessage(result.ApiCode.FILEUPLOADERROR))
		return
	}

	now := time.Now()
	ext := path.Ext(file.Filename)
	fileName := strconv.Itoa(now.Nanosecond()) + ext

	// 文件保存路径（相对于项目根目录）
	filePath := fmt.Sprintf("%s%s/%s",
		config.Config.ImageSettings.UploadDir,
		now.Format("20060102"),
		fileName)

	// 创建目录
	dir := fmt.Sprintf("%s%s",
		config.Config.ImageSettings.UploadDir,
		now.Format("20060102"))

	if err := util.CreateDir(dir); err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), "创建目录失败: "+err.Error())
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		result.Failed(c, int(result.ApiCode.FILEUPLOADERROR), "保存文件失败: "+err.Error())
		return
	}

	// 构造浏览器可访问的 URL
	relativeURL := "/api/v1/upload/" + now.Format("20060102") + "/" + fileName

	// 优先使用配置的 ImageHost，如果为空则使用当前请求的 Host
	imageHost := config.Config.ImageSettings.ImageHost
	log.Log().Infof("Upload - Configured imageHost: [%s]", imageHost)
	if imageHost == "" || imageHost == "https://www.deviops.cn" {
		// 动态获取当前请求的协议和Host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		// 优先使用 X-Forwarded-Host，否则使用 Host header
		// 注意：X-Forwarded-Host 通常不包含端口，需要从 X-Forwarded-Port 或原始 Host 获取
		host := c.Request.Header.Get("X-Forwarded-Host")
		if host == "" {
			// 直接使用 Host header，它包含端口（如果有的话）
			host = c.Request.Host
		} else {
			// 如果使用了 X-Forwarded-Host，检查是否有 X-Forwarded-Port
			port := c.Request.Header.Get("X-Forwarded-Port")
			if port != "" && port != "80" && port != "443" {
				host = host + ":" + port
			}
		}
		imageHost = scheme + "://" + host
	}

	finalURL := imageHost + relativeURL
	log.Log().Infof("Upload - Final URL: %s", finalURL)

	// 返回结果
	result.Success(c, finalURL)
}

var uploadService = UploadServiceImpl{}

func UploadService() IUploadService {
	return &uploadService
}
