// 鉴权中间件
// author xiaoRui

package middleware

import (
	"encoding/json"
	"fmt"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 检查是否为websocket升级请求
		if c.GetHeader("Upgrade") == "websocket" {
			// 从URL参数获取token
			token := c.Query("token")
			if token != "" {
				if _, err := jwt.ValidateToken(token); err == nil {
					c.Next()
					return
				}
			}

			// 尝试从WebSocket消息中获取认证信息
			conn, buf, err := c.Writer.(http.Hijacker).Hijack()
			if err != nil {
				result.Failed(c, int(result.ApiCode.WEBSOCKETERROR), result.ApiCode.GetMessage(result.ApiCode.WEBSOCKETERROR))
				c.Abort()
				return
			}
			defer conn.Close()
			defer buf.Flush()

			// 读取第一条消息作为认证消息
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				result.Failed(c, int(result.ApiCode.WEBSOCKETERROR), result.ApiCode.GetMessage(result.ApiCode.WEBSOCKETERROR))
				c.Abort()
				return
			}
			if op != ws.OpText && op != ws.OpBinary {
				result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
				c.Abort()
				return
			}

			// 解析JSON格式的认证消息
			var authMsg struct {
				Type  string `json:"type"`
				Token string `json:"token"`
			}
			if err := json.Unmarshal(msg, &authMsg); err != nil || authMsg.Type != "auth" {
				result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
				c.Abort()
				return
			}

			if _, err := jwt.ValidateToken(authMsg.Token); err == nil {
				c.Next()
				return
			} else {
				result.Failed(c, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
				c.Abort()
				return
			}
		}

		// 常规HTTP请求鉴权 - 支持SSE的token query参数
		authHeader := c.Request.Header.Get("Authorization")
		
		// 如果没有Authorization头，检查是否为SSE连接并从query参数获取token
		if authHeader == "" {
			token := c.Query("token")
			if token != "" {
				// 验证token
				mc, err := jwt.ValidateToken(token)
				if err == nil {
					c.Set(constant.ContextKeyUserObj, mc)
					c.Next()
					return
				} else {
					// 添加调试信息
					fmt.Printf("[DEBUG] Token验证失败: %v, token长度: %d, Secret: %s\n", err, len(token), string(jwt.Secret))
				}
			}
			
			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
			c.Abort()
			return
		}
		mc, err := jwt.ValidateToken(parts[1])
		if err != nil {
			result.Failed(c, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
			c.Abort()
			return
		}
		c.Set(constant.ContextKeyUserObj, mc)
		c.Next()
	}
}
