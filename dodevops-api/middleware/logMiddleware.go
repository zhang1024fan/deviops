// 操作日志中间件
// author xiaoRui

package middleware

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行下一个处理器
		c.Next()

		// 请求处理完成后记录操作日志
		method := strings.ToLower(c.Request.Method)
		sysAdmin, err := jwt.GetAdmin(c)

		// 只记录非GET请求，且用户认证成功的操作
		if method != "get" && err == nil && sysAdmin != nil {
			url := c.Request.URL.Path
			log := model.SysOperationLog{
				AdminId:     sysAdmin.ID,
				Username:    sysAdmin.Username,
				Method:      method,
				Ip:          c.ClientIP(),
				Url:         url,
				Description: GetAPIDescription(url, method), // 自动匹配API描述
				CreateTime:  util.HTime{Time: time.Now()},
			}
			dao.CreateSysOperationLog(log)
		}
	}
}
