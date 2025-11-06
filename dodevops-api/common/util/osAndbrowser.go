// os和browser和工具类
// @author xiaoRui

package util

import (
	"github.com/gin-gonic/gin"
	useragent "github.com/wenlng/go-user-agent"
)

// GetOs 获取os
func GetOs(c *gin.Context) string {
	userAgent := c.Request.Header.Get("User-Agent")
	os := useragent.GetOsName(userAgent)
	return os
}

// GetBrowser 获取browser
func GetBrowser(c *gin.Context) string {
	userAgent := c.Request.Header.Get("User-Agent")
	browser := useragent.GetBrowserName(userAgent)
	return browser
}
