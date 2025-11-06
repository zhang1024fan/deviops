// pkg/log/logger.go

package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// 设置默认日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 创建 logs 目录（如果不存在）
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		logger.Errorf("Failed to create logs directory: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// 同时输出到 stdout 和文件
		mw := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(mw)
	} else {
		logger.Info("Failed to open log file, using stdout only")
	}

	// 🔥 强制设置为 TextFormatter，并启用颜色输出
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 - 15:04:05",
		ForceColors:     true, // 强制启用颜色，即使输出到文件
	})

	// 输出一条 debug 日志确认当前 formatter 是 TextFormatter
	logger.Debug("Logger initialized with TextFormatter.")
}

// Log 返回全局 Logger 实例，用于业务日志记录
func Log() *logrus.Logger {
	return logger
}

// CustomGinLogger 自定义 Gin 中间件：只输出简洁的文本日志
// pkg/log/logger.go

// CustomGinLogger 自定义 Gin 中间件：输出简洁的文本日志（与 GORM 风格一致）

const (
	bgGreen    = "\x1b[42m"
	bgRed      = "\x1b[41m"
	bgYellow   = "\x1b[43m"
	bgBlue     = "\x1b[44m"
	bgMagenta  = "\x1b[45m"
	bgCyan     = "\x1b[46m"
	colorReset = "\x1b[0m"
)

func getBackgroundColorForStatusCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return bgGreen
	case code >= 400 && code < 500:
		return bgRed
	case code >= 500:
		return bgYellow
	default:
		return colorReset
	}
}

func getBackgroundColorForMethod(method string) string {
	switch method {
	case "GET":
		return bgBlue
	case "POST":
		return bgMagenta
	case "PUT":
		return bgCyan
	case "DELETE":
		return bgRed
	default:
		return colorReset
	}
}

func CustomGinLogger() gin.HandlerFunc {
	// 检测是否是终端环境（用于决定是否启用颜色）
	isTerminal := false
	if fileInfo, _ := os.Stdout.Stat(); fileInfo != nil {
		isTerminal = (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 根据是否是终端决定是否使用颜色
		var logMessage string
		if isTerminal {
			// 终端环境：使用颜色
			statusBgColor := getBackgroundColorForStatusCode(statusCode)
			methodBgColor := getBackgroundColorForMethod(method)
			logMessage = fmt.Sprintf("[GIN] %s | %s%3d%s | %13v | %15s | %s%-6s%s %q",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusBgColor, statusCode, colorReset,
				latency,
				ip,
				methodBgColor, method, colorReset,
				path,
			)
		} else {
			// 文件日志：不使用颜色
			logMessage = fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-6s %q",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				ip,
				method,
				path,
			)
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			logger.Error(logMessage)
		} else if statusCode >= 400 {
			logger.Warn(logMessage)
		} else {
			logger.Info(logMessage)
		}
	}
}

// Setup 初始化日志（为了兼容migrate.go的调用）
func Setup() {
	// 日志已经在init()方法中初始化了，这里只是提供一个兼容性方法
	if logger == nil {
		panic("Logger initialization failed")
	}
}
