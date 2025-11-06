package controller

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/task/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketController WebSocket控制器
type WebSocketController struct {
	service service.ITaskAnsibleService
}

func NewWebSocketController(service service.ITaskAnsibleService) *WebSocketController {
	return &WebSocketController{service: service}
}

// 配置WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源（生产环境应该限制）
		return true
	},
}

// LogMessage WebSocket日志消息结构
type LogMessage struct {
	Type      string `json:"type"`      // "log" | "status" | "complete" | "error"
	Content   string `json:"content"`   // 日志内容
	LineNum   int    `json:"line_num"`  // 行号
	Timestamp string `json:"timestamp"` // 时间戳
	TaskID    uint   `json:"task_id"`   // 任务ID
	WorkID    uint   `json:"work_id"`   // 子任务ID
}

// GetJobLogWS 通过WebSocket获取任务日志
// @Summary 通过WebSocket实时获取Ansible任务日志
// @Description 建立WebSocket连接实时推送任务执行日志
// @Tags 任务作业
// @Param id path int true "任务ID"
// @Param work_id path int true "子任务ID"  
// @Param token query string false "认证token"
// @Router /api/v1/ws/task/ansible/{id}/log/{work_id} [get]
// @Security ApiKeyAuth
func (c *WebSocketController) GetJobLogWS(ctx *gin.Context) {
	// WebSocket认证检查
	token := ctx.Query("token")
	if token == "" {
		// 尝试从Authorization header获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// 验证token
	if token != "" {
		_, err := jwt.ValidateToken(token)
		if err != nil {
			result.Failed(ctx, http.StatusUnauthorized, "认证失败")
			return
		}
	}

	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	workID, err := strconv.ParseUint(ctx.Param("work_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的子任务ID")
		return
	}

	// 升级为WebSocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 开始处理日志
	c.handleLogStreaming(conn, uint(taskID), uint(workID))
}

// handleLogStreaming 处理日志流传输
func (c *WebSocketController) handleLogStreaming(conn *websocket.Conn, taskID, workID uint) {
	// 获取任务信息
	work, err := c.getWorkInfo(taskID, workID)
	if err != nil {
		c.sendMessage(conn, LogMessage{
			Type:    "error",
			Content: fmt.Sprintf("获取任务信息失败: %v", err),
			TaskID:  taskID,
			WorkID:  workID,
		})
		return
	}

	// 检查日志文件
	logPath := c.getLogFilePath(work.LogPath)
	if logPath == "" {
		c.sendMessage(conn, LogMessage{
			Type:    "error",
			Content: "日志文件路径不存在",
			TaskID:  taskID,
			WorkID:  workID,
		})
		return
	}

	// 发送状态信息
	c.sendMessage(conn, LogMessage{
		Type:      "status",
		Content:   fmt.Sprintf("任务状态: %d", work.Status),
		TaskID:    taskID,
		WorkID:    workID,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	})

	// 根据任务状态选择处理方式
	if work.Status == 3 || work.Status == 4 {
		// 已完成任务：发送完整日志
		c.sendCompleteLog(conn, logPath, taskID, workID, work.Status)
	} else {
		// 运行中任务：实时流式推送
		c.sendRealtimeLog(conn, logPath, taskID, workID)
	}
}

// sendCompleteLog 发送完整日志
func (c *WebSocketController) sendCompleteLog(conn *websocket.Conn, logPath string, taskID, workID uint, status int) {
	file, err := os.Open(logPath)
	if err != nil {
		c.sendMessage(conn, LogMessage{
			Type:    "error",
			Content: fmt.Sprintf("打开日志文件失败: %v", err),
			TaskID:  taskID,
			WorkID:  workID,
		})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 发送日志行
		c.sendMessage(conn, LogMessage{
			Type:      "log",
			Content:   line,
			LineNum:   lineNum,
			TaskID:    taskID,
			WorkID:    workID,
			Timestamp: time.Now().Format("15:04:05"),
		})

		// 每10行检查一次连接状态
		if lineNum%10 == 0 {
			if c.isConnClosed(conn) {
				return
			}
		}
	}

	// 发送完成信号
	c.sendMessage(conn, LogMessage{
		Type:      "complete",
		Content:   fmt.Sprintf("任务完成，状态: %d，总共 %d 行日志", status, lineNum),
		LineNum:   lineNum,
		TaskID:    taskID,
		WorkID:    workID,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	})
}

// sendRealtimeLog 发送实时日志
func (c *WebSocketController) sendRealtimeLog(conn *websocket.Conn, logPath string, taskID, workID uint) {
	var file *os.File
	var err error

	// 等待日志文件创建
	for range 30 {
		file, err = os.Open(logPath)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			c.sendMessage(conn, LogMessage{
				Type:    "error",
				Content: fmt.Sprintf("打开日志文件失败: %v", err),
				TaskID:  taskID,
				WorkID:  workID,
			})
			return
		}
		time.Sleep(1 * time.Second)
	}

	if file == nil {
		c.sendMessage(conn, LogMessage{
			Type:    "error",
			Content: "等待日志文件创建超时",
			TaskID:  taskID,
			WorkID:  workID,
		})
		return
	}
	defer file.Close()

	// 实时读取并推送日志
	reader := bufio.NewReader(file)
	var lastPos int64 = 0
	lineNum := 0
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// 使用range遍历ticker.C通道
	for range ticker.C {
		// 检查连接状态
		if c.isConnClosed(conn) {
			return
		}

		// 检查文件大小变化
		stat, err := file.Stat()
		if err != nil {
			continue
		}

		if stat.Size() > lastPos {
			// 读取新内容
			file.Seek(lastPos, io.SeekStart)
			reader = bufio.NewReader(file)

			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}

				lineNum++
				c.sendMessage(conn, LogMessage{
					Type:      "log",
					Content:   strings.TrimSpace(line),
					LineNum:   lineNum,
					TaskID:    taskID,
					WorkID:    workID,
					Timestamp: time.Now().Format("15:04:05"),
				})
			}

			lastPos, _ = file.Seek(0, io.SeekCurrent)
		}

		// 检查任务是否完成
		currentWork, err := c.getWorkInfo(taskID, workID)
		if err == nil && (currentWork.Status == 3 || currentWork.Status == 4) {
			// 最后读取一次，然后退出
			time.Sleep(500 * time.Millisecond)
			stat, _ = file.Stat()
			if stat.Size() > lastPos {
				file.Seek(lastPos, io.SeekStart)
				reader = bufio.NewReader(file)
				for {
					line, err := reader.ReadString('\n')
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}
					lineNum++
					c.sendMessage(conn, LogMessage{
						Type:      "log",
						Content:   strings.TrimSpace(line),
						LineNum:   lineNum,
						TaskID:    taskID,
						WorkID:    workID,
						Timestamp: time.Now().Format("15:04:05"),
					})
				}
			}

			// 发送完成信号
			c.sendMessage(conn, LogMessage{
				Type:      "complete",
				Content:   fmt.Sprintf("任务完成，状态: %d，总共 %d 行日志", currentWork.Status, lineNum),
				LineNum:   lineNum,
				TaskID:    taskID,
				WorkID:    workID,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			})

			return
		}
	}
}

// 辅助方法
func (c *WebSocketController) getWorkInfo(taskID, workID uint) (*WorkInfo, error) {
	// 调用service获取work信息
	work, err := c.service.GetWorkByID(taskID, workID)
	if err != nil {
		return nil, err
	}
	
	return &WorkInfo{
		ID:      work.ID,
		TaskID:  work.TaskID,
		Status:  work.Status,
		LogPath: work.LogPath,
	}, nil
}

type WorkInfo struct {
	ID      uint
	TaskID  uint
	Status  int
	LogPath string
}

func (c *WebSocketController) getLogFilePath(logPath string) string {
	if logPath == "" {
		return ""
	}

	if filepath.IsAbs(logPath) {
		return logPath
	}

	// 转换为绝对路径
	cwd, _ := os.Getwd()
	if strings.Contains(cwd, "/task/") {
		projectRoot := strings.Split(cwd, "/task/")[0]
		return filepath.Join(projectRoot, logPath)
	}
	return filepath.Join(cwd, logPath)
}

func (c *WebSocketController) sendMessage(conn *websocket.Conn, msg LogMessage) error {
	return conn.WriteJSON(msg)
}

func (c *WebSocketController) isConnClosed(conn *websocket.Conn) bool {
	// 设置短暂的写超时来检测连接状态
	conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	err := conn.WriteMessage(websocket.PingMessage, nil)
	conn.SetWriteDeadline(time.Time{})
	return err != nil
}