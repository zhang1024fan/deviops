// middleware/logger.go （可选保留，也可以删除）

package middleware

import (
	"dodevops-api/pkg/log"
	"github.com/gin-gonic/gin"
)

// Logger 返回标准 Gin 请求日志中间件
func Logger() gin.HandlerFunc {
	return log.CustomGinLogger()
}
