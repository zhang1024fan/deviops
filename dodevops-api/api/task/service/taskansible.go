package service

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	cmdbmodel "dodevops-api/api/cmdb/model"
	configcentermodel "dodevops-api/api/configcenter/model"
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	taskmodel "dodevops-api/api/task/model"
	"dodevops-api/common"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// RealTimeLogWriter 实时日志写入器，支持立即刷新到磁盘
type RealTimeLogWriter struct {
	file *os.File
}

// Write 实现io.Writer接口，每次写入后立即刷新到磁盘
func (w *RealTimeLogWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.Write(p)
	if err != nil {
		return n, err
	}
	// 立即刷新到磁盘，确保SSE能实时读取
	w.file.Sync()
	return n, nil
}

// WriteWithTime 带时间戳的写入
func (w *RealTimeLogWriter) WriteWithTime(content string) error {
	_, err := w.Write([]byte(content))
	return err
}

// ITaskAnsibleService 定义Ansible任务服务接口
type ITaskAnsibleService interface {
	CreateTask(c *gin.Context, req *CreateTaskRequest)     // 创建任务
	CreateK8sTask(c *gin.Context, req *CreateK8sTaskRequest) // 创建K8s任务
	List(c *gin.Context, page, size int)                   // 获取任务列表
	StartJob(c *gin.Context, taskID uint)                  // 启动任务
	StopJob(c *gin.Context, taskID, workID uint)           // 停止任务
	GetJobLog(c *gin.Context, taskID, workID uint)         // 实时获取任务日志(SSE)
	GetJobStatus(c *gin.Context, taskID, workID uint)      // 获取任务状态
	GetTaskDetail(c *gin.Context, taskID uint)             // 获取任务详情
	GetWorkByID(taskID, workID uint) (*model.TaskAnsibleWork, error) // 获取子任务详情
	DeleteTask(c *gin.Context, taskID uint)                // 删除任务
	GetTasksByName(c *gin.Context, name string)            // 根据名称模糊查询任务
	GetTasksByType(c *gin.Context, taskType int)           // 根据类型查询任务
}

// CreateTaskRequest 创建任务请求参数
type CreateTaskRequest struct {
	TaskType         int               `json:"taskType"`
	Name             string            `json:"name"`
	HostGroups       map[string][]uint `json:"hostGroups"`
	GitRepo          string            `json:"gitRepo"`
	RolesContent     []byte            `json:"rolesContent"`
	PlaybookContents [][]byte          `json:"playbookContents"`
	Variables        map[string]string `json:"variables"`
}

// CreateK8sTaskRequest 创建K8s任务请求参数
type CreateK8sTaskRequest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	ClusterName       string   `json:"cluster_name"`
	ClusterVersion    string   `json:"cluster_version"`
	DeploymentMode    int      `json:"deployment_mode"`
	MasterHostIDs     []uint   `json:"master_host_ids"`
	WorkerHostIDs     []uint   `json:"worker_host_ids"`
	EtcdHostIDs       []uint   `json:"etcd_host_ids"`
	EnabledComponents []string `json:"enabled_components"`
	PrivateRegistry   string         `json:"private_registry"`
	RegistryUsername  string         `json:"registry_username"`
	RegistryPassword  string         `json:"registry_password"`
	RegistryConfig    *RegistryConfig `json:"registry_config"` // 新的嵌套配置格式
}

// RegistryConfig 镜像仓库配置
type RegistryConfig struct {
	PrivateRegistry    string `json:"privateRegistry"`    // 私有镜像仓库地址
	RegistryUsername   string `json:"registryUsername"`   // 镜像仓库用户名
	RegistryPassword   string `json:"registryPassword"`   // 镜像仓库密码
	UsePrivateRegistry bool   `json:"usePrivateRegistry"` // 是否使用私有仓库
}

// K8sNodeInfo K8s节点信息
type K8sNodeInfo struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// K8sConfigJSON K8s配置文件结构
type K8sConfigJSON struct {
	Cluster struct {
		Name           string `json:"name"`
		Version        string `json:"version"`
		DeploymentMode int    `json:"deployment_mode"`
	} `json:"cluster"`
	Nodes struct {
		Masters []K8sNodeInfo `json:"masters"`
		Workers []K8sNodeInfo `json:"workers"`
		Etcd    []K8sNodeInfo `json:"etcd"`
	} `json:"nodes"`
	Components struct {
		Enabled []string `json:"enabled"`
	} `json:"components"`
	Registry struct {
		PrivateRegistry string `json:"private_registry,omitempty"`
		Username        string `json:"username,omitempty"`
		Password        string `json:"password,omitempty"`
	} `json:"registry,omitempty"`
}

// TaskAnsibleServiceImpl 实现Ansible任务服务
type TaskAnsibleServiceImpl struct {
	dao *dao.TaskAnsibleDao
}

func NewTaskAnsibleService(db *gorm.DB) ITaskAnsibleService {
	return &TaskAnsibleServiceImpl{
		dao: dao.NewTaskAnsibleDao(db),
	}
}

// HostSSHInfo 主机SSH连接信息
type HostSSHInfo struct {
	ID       uint
	IP       string
	Port     int
	User     string
	Password string
	Key      string
	AuthType int    // 认证类型：1-密码，2-私钥，3-公钥免认证
}

// HostSSHInfoCollection 主机信息集合
type HostSSHInfoCollection struct {
	Groups    map[string][]HostSSHInfo
	HostInfos map[uint]HostSSHInfo
}

// GetHostSSHInfo 获取主机SSH信息
func (s *TaskAnsibleServiceImpl) GetHostSSHInfo(hostGroups map[string][]uint) (*HostSSHInfoCollection, error) {
	// 获取所有唯一主机ID
	allHostIDs := make([]uint, 0)
	idMap := make(map[uint]bool)
	for _, ids := range hostGroups {
		for _, id := range ids {
			if id > 0 && !idMap[id] { // 确保ID有效且不重复
				idMap[id] = true
				allHostIDs = append(allHostIDs, id)
			}
		}
	}

	// 如果没有获取到有效主机ID，返回错误
	if len(allHostIDs) == 0 {
		return nil, fmt.Errorf("没有获取到有效的主机ID")
	}

	// 查询主机信息
	var hosts []cmdbmodel.CmdbHost
	if err := s.dao.DB.Where("id IN ?", allHostIDs).Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	// 构建返回结果
	collection := &HostSSHInfoCollection{
		Groups:    make(map[string][]HostSSHInfo),
		HostInfos: make(map[uint]HostSSHInfo),
	}

	// 构建主机信息映射
	for _, host := range hosts {
		// 获取SSH认证信息
		var ecsAuth configcentermodel.EcsAuth
		if err := common.GetDB().First(&ecsAuth, host.SSHKeyID).Error; err != nil {
			return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
		}

		info := HostSSHInfo{
			ID:       host.ID,
			IP:       host.SSHIP,
			Port:     host.SSHPort,
			User:     host.SSHName,
			AuthType: ecsAuth.Type,
		}

		// 根据认证类型设置相应字段
		switch ecsAuth.Type {
		case 1: // 密码认证
			info.Password = ecsAuth.Password
		case 2: // 私钥认证
			info.Key = ecsAuth.PublicKey
		case 3: // 公钥免认证
			// 不需要设置额外信息
		}
		collection.HostInfos[host.ID] = info
	}

	// 构建分组信息
	for groupName, ids := range hostGroups {
		for _, id := range ids {
			if info, ok := collection.HostInfos[id]; ok {
				collection.Groups[groupName] = append(collection.Groups[groupName], info)
			}
		}
	}

	return collection, nil
}

// GetAllHostIDs 获取所有主机ID
func (c *HostSSHInfoCollection) GetAllHostIDs() []uint {
	ids := make([]uint, 0, len(c.HostInfos))
	for id := range c.HostInfos {
		ids = append(ids, id)
	}
	return ids
}

// GenerateInventory 生成Ansible inventory文件内容
func (c *HostSSHInfoCollection) GenerateInventory() string {
	var builder strings.Builder

	// 写入分组信息
	for groupName, hosts := range c.Groups {
		builder.WriteString(fmt.Sprintf("[%s]\n", groupName))
		for _, host := range hosts {
			builder.WriteString(fmt.Sprintf("%s ansible_ssh_port=%d ansible_ssh_user=%s",
				host.IP, host.Port, host.User))

			switch host.AuthType {
			case 1: // 密码认证
				if host.Password != "" {
					builder.WriteString(fmt.Sprintf(" ansible_ssh_pass=%s", host.Password))
				}
			case 2: // 私钥认证
				if host.Key != "" {
					builder.WriteString(fmt.Sprintf(" ansible_ssh_private_key_file=%s", host.Key))
				}
			case 3: // 公钥免认证
				// 不添加额外的认证参数，使用系统默认SSH配置
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// List 获取任务列表
func (s *TaskAnsibleServiceImpl) List(c *gin.Context, page, size int) {
	tasks, total, err := s.dao.List(page, size)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("获取任务列表失败: %v", err)})
		return
	}
	c.JSON(200, gin.H{"data": tasks, "total": total})
}

// DeleteTask 删除任务
func (s *TaskAnsibleServiceImpl) DeleteTask(c *gin.Context, taskID uint) {
	// 1. 首先获取任务信息（用于删除相关文件目录）
	task, err := s.dao.GetTaskDetail(taskID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "任务不存在"})
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("获取任务信息失败: %v", err)})
		}
		return
	}

	// 2. 检查任务状态，不允许删除正在运行的任务
	if task.Status == 2 {
		c.JSON(400, gin.H{"error": "不能删除正在运行中的任务，请先停止任务"})
		return
	}

	// 3. 从数据库删除任务和子任务（级联删除）
	if err := s.dao.Delete(taskID); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("删除任务失败: %v", err)})
		return
	}

	// 4. 删除任务相关的文件目录（异步处理，避免影响响应速度）
	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		// 删除任务目录: task/{taskID}/{taskName}
		taskDir := fmt.Sprintf("task/%d/%s", taskID, task.Name)
		if _, err := os.Stat(taskDir); err == nil {
			os.RemoveAll(taskDir)
		}

		// 删除日志目录: logs/ansible/{taskID}
		logDir := fmt.Sprintf("logs/ansible/%d", taskID)
		if _, err := os.Stat(logDir); err == nil {
			os.RemoveAll(logDir)
		}
	}()
	c.JSON(200, gin.H{
		"message": "任务删除成功",
		"data": gin.H{
			"deleted_task_id":   taskID,
			"deleted_task_name": task.Name,
			"deleted_sub_tasks": len(task.Works),
		},
	})
}

// GetJobLog 实时获取任务日志(SSE实现) - 优化版本
func (s *TaskAnsibleServiceImpl) GetJobLog(c *gin.Context, taskID, workID uint) {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	
	// 检查认证状态（调试信息）
	token := c.Query("token")
	if token == "" || token == "null" {
		// 对于已完成的任务，即使token为空也允许读取日志
		// sendSSEError(c, "认证失败：token为空")
		// return
	}

	// 获取任务记录（仅查询一次）
	work, err := s.dao.GetWorkByID(taskID, workID)
	if err != nil {
		sendSSEError(c, fmt.Sprintf("获取任务记录失败: %v", err))
		return
	}

	// 检查日志路径
	if work.LogPath == "" {
		sendSSEError(c, "日志文件路径不存在")
		return
	}

	// 确保使用绝对路径读取日志文件
	var logPath string
	if filepath.IsAbs(work.LogPath) {
		logPath = work.LogPath
	} else {
		// 如果是相对路径，转换为绝对路径
		// 获取当前工作目录
		cwd, _ := os.Getwd()
		// 检查是否在任务子目录中，如果是则返回到项目根目录
		if strings.Contains(cwd, "/task/") {
			// 切换到项目根目录计算绝对路径
			projectRoot := strings.Split(cwd, "/task/")[0]
			logPath = filepath.Join(projectRoot, work.LogPath)
		} else {
			// 已经在项目根目录
			logPath = filepath.Join(cwd, work.LogPath)
		}
	}

	// 缓存初始状态，避免重复查询
	initialStatus := work.Status
	isCompleted := initialStatus == 3 || initialStatus == 4

	// 等待日志文件创建（优化：如果任务已完成且无文件，直接返回）
	var file *os.File
	maxWaitTime := 30
	if isCompleted {
		maxWaitTime = 3 // 已完成任务最多等待3秒
	}

	for i := 0; i < maxWaitTime; i++ {
		file, err = os.Open(logPath)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			sendSSEError(c, fmt.Sprintf("打开日志文件失败: %v", err))
			return
		}
		// 如果任务已完成但没有日志文件，提前退出
		if isCompleted && i >= 2 {
			sendSSEError(c, "任务已完成但日志文件不存在")
			return
		}
		time.Sleep(1 * time.Second)
	}
	if file == nil {
		sendSSEError(c, "日志文件不存在或创建超时")
		return
	}
	defer file.Close()

	// 从文件开头开始读取，跟踪文件变化
	var lastPos int64 = 0
	reader := bufio.NewReader(file)

	// 读取完整的日志文件内容
	lineCount := 0
	batchSize := 10 // 每10行flush一次，平衡性能和实时性
	
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			sendSSEError(c, fmt.Sprintf("读取日志失败: %v", err))
			return
		}
		lineCount++
		
		// 发送日志内容 (确保非空行才发送)
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			fmt.Fprintf(c.Writer, "data: %s\n\n", trimmed)
		} else {
			// 发送空行
			fmt.Fprintf(c.Writer, "data: \n\n")
		}
		
		// 批量flush，减少网络开销
		if lineCount%batchSize == 0 {
			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}

			// 检查客户端连接状态
			if c.Request.Context().Err() != nil {
				return
			}
		}
		
		lastPos, _ = file.Seek(0, io.SeekCurrent)
	}
	
	// 最后flush剩余数据
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	// 如果任务已完成，发送完成信号并退出
	if isCompleted {
		fmt.Fprintf(c.Writer, "event: complete\ndata: Task completed with status %d, total lines: %d\n\n", initialStatus, lineCount)
		c.Writer.Flush()
		return
	}

	// 实时监控文件变化（仅对运行中的任务）
	ticker := time.NewTicker(200 * time.Millisecond) // 提高检查频率到200ms，增强实时性
	defer ticker.Stop()

	// 状态检查计数器，减少数据库查询频率
	statusCheckCounter := 0
	statusCheckInterval := 10 // 每10次文件检查才查询一次数据库状态(即每2秒检查一次状态)

	// 添加心跳机制，确保连接活跃
	heartbeatCounter := 0
	heartbeatInterval := 25 // 每25次检查发送一次心跳(即每5秒)

	for {
		select {
		case <-c.Done():
			return
		case <-ticker.C:
			// 发送心跳，保持连接活跃
			heartbeatCounter++
			if heartbeatCounter >= heartbeatInterval {
				heartbeatCounter = 0
				fmt.Fprintf(c.Writer, ": heartbeat\n\n") // SSE注释格式的心跳
				c.Writer.Flush()
			}

			// 检查文件大小是否变化
			stat, err := file.Stat()
			if err != nil {
				continue
			}

			// 如果文件增大了，读取新内容
			if stat.Size() > lastPos {
				// 移动到上次读取的位置
				file.Seek(lastPos, io.SeekStart)
				reader = bufio.NewReader(file)

				// 读取新内容
				for {
					line, err := reader.ReadString('\n')
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}
					// 发送新日志内容
					fmt.Fprintf(c.Writer, "data: %s\n\n", strings.TrimSpace(line))
					c.Writer.Flush()
				}
				// 更新最后读取位置
				lastPos, _ = file.Seek(0, io.SeekCurrent)
			}

			// 减少数据库查询频率：每10次文件检查才查询一次状态，且仅对运行中任务
			statusCheckCounter++
			if statusCheckCounter >= statusCheckInterval && !isCompleted {
				statusCheckCounter = 0

				// 检查任务是否已经完成 (使用轻量级查询提升性能)
				currentStatus, err := s.dao.GetWorkStatus(taskID, workID)
				if err == nil && (currentStatus == 3 || currentStatus == 4) {
					// 任务已完成，再读取一次最终内容后退出
					isCompleted = true // 标记为已完成，避免后续重复查询

					time.Sleep(300 * time.Millisecond) // 等待最后的日志写入
					stat, _ := file.Stat()
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
							fmt.Fprintf(c.Writer, "data: %s\n\n", strings.TrimSpace(line))
							c.Writer.Flush()
						}
					}
					// 发送完成信号
					fmt.Fprintf(c.Writer, "event: complete\ndata: Task completed with status %d\n\n", currentStatus)
					c.Writer.Flush()
					return
				}
			}
		}
	}
}

func sendSSEError(c *gin.Context, msg string) {
	fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", msg)
	c.Writer.Flush()
}

// GetJobStatus 获取任务状态
func (s *TaskAnsibleServiceImpl) GetJobStatus(c *gin.Context, taskID, workID uint) {
	status, err := s.dao.GetJobStatus(taskID, workID)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("获取任务状态失败: %v", err)})
		return
	}
	c.JSON(200, gin.H{"data": status})
}

// GetTaskDetail 获取任务详情
// 功能：根据任务ID获取完整任务信息，包括：
// - 任务基本信息（名称、类型、状态等）
// - 关联的子任务(works)列表
// - 全局变量配置
// - 主机分组信息
// - 执行记录等
func (s *TaskAnsibleServiceImpl) GetTaskDetail(c *gin.Context, taskID uint) {
	task, err := s.dao.GetTaskDetail(taskID)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取任务详情失败: %v", err))
		return
	}

	// 构建精简的子任务列表
	type WorkInfo struct {
		WorkID        uint   `json:"workid"`        // 子任务ID
		EntryFileName string `json:"EntryFileName"` // 子任务名称
		Status        int    `json:"status"`        // 子任务状态
		Duration      int    `json:"Duration"`      // 子任务执行耗时
	}

	works := make([]WorkInfo, len(task.Works))
	for i, work := range task.Works {
		works[i] = WorkInfo{
			WorkID:        work.ID,
			EntryFileName: work.EntryFileName,
			Status:        work.Status,
			Duration:      work.Duration,
		}
	}

	// 构建精简的任务信息
	taskInfo := gin.H{
		"ID":        task.ID,        // 父任务ID
		"Name":      task.Name,      // 父任务名称
		"TaskCount": task.TaskCount, // 子任务数量
		"Works":     works,          // 子任务列表
	}

	result.Success(c, gin.H{
		"task_info": taskInfo,
	})
}

// GetWorkByID 获取子任务详情
func (s *TaskAnsibleServiceImpl) GetWorkByID(taskID, workID uint) (*model.TaskAnsibleWork, error) {
	return s.dao.GetWorkByID(taskID, workID)
}

// StartJob 启动任务
func (s *TaskAnsibleServiceImpl) StartJob(c *gin.Context, taskID uint) {
	// 1. 获取任务详情（包含子任务）
	task, err := s.dao.GetTaskDetail(taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("获取任务失败: %v", err)})
		return
	}

	// 检查任务是否存在子任务
	if len(task.Works) == 0 {
		c.JSON(400, gin.H{"error": "任务没有子任务，无法执行"})
		return
	}

	// 2. 更新任务状态为运行中（状态=2）
	if err := s.dao.DB.Model(&model.TaskAnsible{}).Where("id = ?", taskID).Update("status", 2).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("更新任务状态失败: %v", err)})
		return
	}

	// 3. 异步执行Ansible任务（优化版本 - 直接执行，无需重复查询）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.updateTaskErrorStatus(taskID, fmt.Errorf("任务执行异常: %v", r))
			}
		}()

		// 构建任务目录路径：task/{taskID}/{taskName}
		taskDir := fmt.Sprintf("task/%d/%s", taskID, task.Name)

		// 检查任务目录是否存在
		if _, err := os.Stat(taskDir); os.IsNotExist(err) {
			s.updateTaskErrorStatus(taskID, fmt.Errorf("任务目录不存在: %s", taskDir))
			return
		}

		// 获取当前工作目录
		originalDir, _ := os.Getwd()

		// 执行每个子任务
		allSuccess := true
		for _, work := range task.Works {

			// 创建日志目录（使用绝对路径）
			absLogDir := filepath.Join(originalDir, fmt.Sprintf("logs/ansible/%d/%d", taskID, work.ID))
			if err := os.MkdirAll(absLogDir, 0755); err != nil {
				s.updateTaskErrorStatus(taskID, fmt.Errorf("创建日志目录失败: %v", err))
				return
			}

			// 生成日志文件路径（使用绝对路径）
			// 对于K8s任务，使用脚本名而不是完整路径
			var logFileName string
			if task.Type == 3 {
				logFileName = "deploy-simple.sh"
			} else {
				logFileName = work.EntryFileName
			}
			absLogPath := filepath.Join(absLogDir, fmt.Sprintf("%s.log", logFileName))
			// 用于数据库存储的相对路径
			relativeLogPath := fmt.Sprintf("logs/ansible/%d/%d/%s.log", taskID, work.ID, logFileName)

			// 更新子任务状态为运行中，记录开始时间和日志路径
			workStartTime := time.Now()
			s.dao.DB.Model(&model.TaskAnsibleWork{}).
				Where("id = ?", work.ID).
				Updates(map[string]interface{}{
					"status":     2, // 运行中
					"start_time": workStartTime,
					"log_path":   relativeLogPath, // 使用相对路径存储到数据库
				})

			// 切换到任务目录
			if err := os.Chdir(taskDir); err != nil {
				s.updateWorkErrorStatus(work.ID, fmt.Errorf("切换到任务目录失败: %v", err))
				allSuccess = false
				continue
			}

			// 检查playbook文件是否存在（K8s任务跳过此检查）
			var playbookPath string
			if task.Type != 3 {
				playbookPath = work.EntryFileName
				if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
					os.Chdir(originalDir)
					s.updateWorkErrorStatus(work.ID, fmt.Errorf("Playbook文件不存在: %s", playbookPath))
					allSuccess = false
					continue
				}
			} else {
				playbookPath = work.EntryFileName // K8s任务也需要这个变量，但不检查文件
			}

			// 根据任务类型构建不同的执行命令
			var cmdArgs []string
			if task.Type == 3 { // K8s任务
				// 检查config.json文件是否存在
				if _, err := os.Stat("config.json"); os.IsNotExist(err) {
					os.Chdir(originalDir)
					s.updateWorkErrorStatus(work.ID, fmt.Errorf("config.json文件不存在，任务创建可能有问题"))
					allSuccess = false
					continue
				}

				// 检查部署脚本是否存在
				scriptPath := filepath.Join("scripts", "deploy-simple.sh")
				if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
					os.Chdir(originalDir)
					s.updateWorkErrorStatus(work.ID, fmt.Errorf("K8s部署脚本不存在: %s", scriptPath))
					allSuccess = false
					continue
				}

				// 构建K8s部署命令
				cmdArgs = []string{"bash", scriptPath}
			} else {
				// Ansible任务
				// 检查hosts文件是否存在（创建任务时已生成）
				if _, err := os.Stat("hosts"); os.IsNotExist(err) {
					os.Chdir(originalDir)
					s.updateWorkErrorStatus(work.ID, fmt.Errorf("hosts文件不存在，任务创建可能有问题"))
					allSuccess = false
					continue
				}

				// 构建Ansible命令
				cmdArgs = []string{"ansible-playbook", "-i", "hosts", playbookPath, "-v"}
			}

			// 执行命令
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

			// 创建日志文件用于实时写入（使用绝对路径）
			logFile, err := os.OpenFile(absLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				os.Chdir(originalDir)
				s.updateWorkErrorStatus(work.ID, fmt.Errorf("创建日志文件失败: %v", err))
				allSuccess = false
				continue
			}

			// 写入命令信息到日志文件
			logFile.WriteString(fmt.Sprintf("[%s] 开始执行任务\n", time.Now().Format("2006-01-02 15:04:05")))
			logFile.WriteString(fmt.Sprintf("命令: %s\n", strings.Join(cmdArgs, " ")))
			logFile.WriteString(fmt.Sprintf("工作目录: %s\n", taskDir))
			logFile.WriteString("==========================================\n")
			logFile.Sync() // 立即刷新到磁盘

			// 创建实时日志写入器
			logWriter := &RealTimeLogWriter{
				file: logFile,
			}

			// 设置命令输出到实时日志写入器（用于SSE实时读取）
			cmd.Stdout = logWriter
			cmd.Stderr = logWriter

			// 执行命令并记录结果
			err = cmd.Run()

			// 记录执行完成
			logWriter.WriteWithTime(fmt.Sprintf("\n[%s] 任务执行完成\n", time.Now().Format("2006-01-02 15:04:05")))
			if err != nil {
				logWriter.WriteWithTime(fmt.Sprintf("执行错误: %v\n", err))
			} else {
				logWriter.WriteWithTime("执行成功\n")
			}
			logFile.Close()

			// 切换回原目录
			os.Chdir(originalDir)

			// 计算执行耗时
			workEndTime := time.Now()
			duration := int(workEndTime.Sub(workStartTime).Seconds())

			// 确定子任务状态
			workStatus := 3 // 成功
			exitCode := 0
			if err != nil {
				workStatus = 4 // 失败
				allSuccess = false
				if cmd.ProcessState != nil {
					exitCode = cmd.ProcessState.ExitCode()
				} else {
					exitCode = -1
				}
			}

			// 更新子任务完成状态
			s.dao.DB.Model(&model.TaskAnsibleWork{}).
				Where("id = ?", work.ID).
				Updates(map[string]interface{}{
					"status":    workStatus,
					"end_time":  workEndTime,
					"duration":  duration,
					"exit_code": exitCode,
					"error_msg": func() string {
						if err != nil {
							return err.Error()
						}
						return ""
					}(),
				})
		}

		// 最后更新父任务状态、总耗时和更新时间
		finalStatus := 3 // 成功
		if !allSuccess {
			finalStatus = 4 // 失败
		}

		// 从数据库重新查询最新的子任务数据来计算总耗时
		var works []model.TaskAnsibleWork
		if err := s.dao.DB.Where("task_id = ?", taskID).Find(&works).Error; err == nil {
			// 计算总耗时（所有子任务耗时的总和）
			var totalDuration int64
			for _, work := range works {
				totalDuration += int64(work.Duration)
			}

			s.dao.DB.Model(&model.TaskAnsible{}).
				Where("id = ?", taskID).
				Updates(map[string]interface{}{
					"status":         finalStatus,
					"total_duration": totalDuration,
					"updated_at":     time.Now(),
				})
		}
	}()

	c.JSON(200, gin.H{"message": "任务已开始执行"})
}

// updateTaskErrorStatus 更新任务为错误状态
func (s *TaskAnsibleServiceImpl) updateTaskErrorStatus(taskID uint, err error) {
	s.dao.DB.Model(&model.TaskAnsible{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status":     4, // 异常
			"error_msg":  err.Error(),
			"updated_at": time.Now(),
		})
}

// updateWorkErrorStatus 更新子任务为错误状态
func (s *TaskAnsibleServiceImpl) updateWorkErrorStatus(workID uint, err error) {
	s.dao.DB.Model(&model.TaskAnsibleWork{}).
		Where("id = ?", workID).
		Updates(map[string]interface{}{
			"status":    4, // 异常
			"end_time":  time.Now(),
			"error_msg": err.Error(),
		})
}

// StopJob 停止任务
func (s *TaskAnsibleServiceImpl) StopJob(c *gin.Context, taskID, workID uint) {
	if err := s.dao.StopJob(taskID, workID); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("停止任务失败: %v", err)})
		return
	}
	c.JSON(200, gin.H{"message": "任务停止成功"})
}

// CreateTask 创建Ansible任务
func (s *TaskAnsibleServiceImpl) CreateTask(c *gin.Context, req *CreateTaskRequest) {
	// 从请求中提取参数
	taskType := req.TaskType
	name := req.Name
	hostGroups := req.HostGroups
	gitRepo := req.GitRepo
	rolesContent := req.RolesContent
	variables := req.Variables

	// 检查任务名称是否已存在
	var count int64
	if err := s.dao.DB.Model(&taskmodel.TaskAnsible{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		result.Failed(c, 500, fmt.Sprintf("检查任务名称失败: %v", err))
		return
	}
	if count > 0 {
		result.Failed(c, 400, "任务名称已存在")
		return
	}

	// 获取主机信息
	hostInfos, err := s.GetHostSSHInfo(hostGroups)
	if err != nil {
		result.Failed(c, 500, err.Error())
		return
	}

	// 获取所有主机ID
	allHostIDs := make([]uint, 0)
	for _, ids := range hostGroups {
		for _, id := range ids {
			if id > 0 { // 确保ID有效
				allHostIDs = append(allHostIDs, id)
			}
		}
	}

	// 创建任务记录
	task := &taskmodel.TaskAnsible{
		Name:       name,
		Type:       taskType, // 1=手动，2=Git导入
		HostGroups: toJSON(hostGroups),
		AllHostIDs: toJSON(allHostIDs),
		Status:     1, // 1表示等待中
	}

	// 如果是Git任务，设置仓库地址
	if taskType == 2 {
		task.GitRepo = gitRepo
	}

	// 保存到数据库
	if err := s.dao.Create(task); err != nil {
		result.Failed(c, 500, fmt.Sprintf("创建任务失败: %v", err))
		return
	}

	// 创建项目目录 (task/ansible任务ID/ansible任务名称/)
	projectDir := fmt.Sprintf("task/%d/%s", task.ID, name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		result.Failed(c, 500, fmt.Sprintf("创建项目目录失败: %v", err))
		return
	}

	// 根据任务类型处理不同的逻辑
	if taskType == 1 {
		// Type=1: 手动创建任务
		if err := s.handleManualTask(c, task, projectDir, hostInfos, rolesContent, variables); err != nil {
			result.Failed(c, 500, err.Error())
			return
		}
	} else if taskType == 2 {
		// Type=2: Git导入任务
		if err := s.handleGitTask(c, task, projectDir, hostInfos, gitRepo, variables); err != nil {
			result.Failed(c, 500, err.Error())
			return
		}
	} else {
		result.Failed(c, 400, "不支持的任务类型")
		return
	}

	// 重新查询任务数据以获取最新的TaskCount和Works
	updatedTask, err := s.dao.GetTaskDetail(task.ID)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取任务详情失败: %v", err))
		return
	}

	result.Success(c, updatedTask)
}

// unzipContent 解压zip内容到目录
func unzipContent(data []byte, dest string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// isZipFile 检查是否为有效的zip文件
func isZipFile(content []byte) bool {
	if len(content) < 4 {
		return false
	}
	// 检查zip文件头签名
	if content[0] != 0x50 || content[1] != 0x4B || content[2] != 0x03 || content[3] != 0x04 {
		return false
	}

	r := bytes.NewReader(content)
	_, err := zip.NewReader(r, r.Size())
	return err == nil
}

// toJSON 转换为JSON字符串
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// handleManualTask 处理手动创建的任务(Type=1)
func (s *TaskAnsibleServiceImpl) handleManualTask(c *gin.Context, task *taskmodel.TaskAnsible, projectDir string, hostInfos *HostSSHInfoCollection, rolesContent []byte, variables map[string]string) error {
	// 创建子目录结构
	subDirs := []string{"vars", "roles"}
	for _, dir := range subDirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("创建%s目录失败: %v", dir, err)
		}
	}

	// 处理roles目录
	if len(rolesContent) > 0 {
		// 验证是否为有效的zip文件
		if !isZipFile(rolesContent) {
			return fmt.Errorf("roles文件必须是有效的zip格式，请确认上传的是正确的zip文件")
		}

		rolesDir := filepath.Join(projectDir, "roles")
		if err := os.MkdirAll(rolesDir, 0755); err != nil {
			return fmt.Errorf("创建roles目录失败: %v", err)
		}
		// 解压roles内容到目录
		if err := unzipContent(rolesContent, rolesDir); err != nil {
			return fmt.Errorf("解压roles失败: %v", err)
		}

		// 检查并修复多出来的roles目录层级
		rolesSubDir := filepath.Join(rolesDir, "roles")
		if _, err := os.Stat(rolesSubDir); err == nil {
			// 如果存在roles/roles子目录，则移动内容到上层目录
			files, err := os.ReadDir(rolesSubDir)
			if err != nil {
				return fmt.Errorf("读取roles子目录失败: %v", err)
			}

			for _, file := range files {
				oldPath := filepath.Join(rolesSubDir, file.Name())
				newPath := filepath.Join(rolesDir, file.Name())
				if err := os.Rename(oldPath, newPath); err != nil {
					return fmt.Errorf("移动roles文件失败: %v", err)
				}
			}

			// 删除空的roles子目录
			if err := os.Remove(rolesSubDir); err != nil {
				return fmt.Errorf("删除空roles子目录失败: %v", err)
			}
		}
	}

	// 处理playbook文件
	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("获取表单数据失败: %v", err)
	}

	files := form.File["playbooks"]
	taskCount := 0 // 记录任务数量

	for _, file := range files {
		// 使用上传文件的原始文件名
		fileName := file.Filename
		filePath := filepath.Join(projectDir, fileName)

		// 确保文件名以.yml结尾
		if !strings.HasSuffix(fileName, ".yml") {
			fileName = fileName + ".yml"
			filePath = filePath + ".yml"
		}

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			return fmt.Errorf("保存playbook文件失败: %v", err)
		}
		// 创建子任务记录
		subTask := &taskmodel.TaskAnsibleWork{
			TaskID:        task.ID,
			EntryFileName: fileName,
			EntryFilePath: filePath,
			Status:        1, // 等待中
		}
		if err := s.dao.DB.Create(subTask).Error; err != nil {
			return fmt.Errorf("创建子任务失败: %v", err)
		}
		taskCount++
	}

	// 更新任务数量到数据库
	if err := s.dao.DB.Model(&taskmodel.TaskAnsible{}).Where("id = ?", task.ID).Update("task_count", taskCount).Error; err != nil {
		return fmt.Errorf("更新任务数量失败: %v", err)
	}

	// 保存全局变量
	if err := s.saveGlobalVariables(c, task, projectDir, variables); err != nil {
		return err
	}

	// 生成inventory文件
	if err := s.generateInventoryFile(c, projectDir, hostInfos); err != nil {
		return err
	}

	return nil
}

// handleGitTask 处理Git导入的任务(Type=2)
func (s *TaskAnsibleServiceImpl) handleGitTask(c *gin.Context, task *taskmodel.TaskAnsible, projectDir string, hostInfos *HostSSHInfoCollection, gitRepo string, variables map[string]string) error {

	// 1. 下载Git仓库
	if err := s.cloneGitRepository(gitRepo, projectDir); err != nil {
		return fmt.Errorf("下载Git仓库失败: %v", err)
	}

	// 2. 解析仓库目录结构，识别playbook文件
	playbookFiles, err := s.parseGitRepository(projectDir)
	if err != nil {
		return fmt.Errorf("解析Git仓库失败: %v", err)
	}

	// 3. 创建子任务记录
	if err := s.createSubTasksFromPlaybooks(task.ID, projectDir, playbookFiles); err != nil {
		return fmt.Errorf("创建子任务失败: %v", err)
	}

	// 4. 处理全局变量(如果有提供)
	if len(variables) > 0 {
		if err := s.saveGlobalVariables(c, task, projectDir, variables); err != nil {
			return err
		}
	}

	// 5. 强制重写hosts文件（根据传入的主机id生成）
	// 注意：不使用仓库中的hosts文件，强制重写
	if err := s.generateInventoryFile(c, projectDir, hostInfos); err != nil {
		return err
	}

	return nil
}

// cloneGitRepository 克隆Git仓库
func (s *TaskAnsibleServiceImpl) cloneGitRepository(gitRepo, projectDir string) error {
	// 使用git命令克隆仓库
	cmd := exec.Command("git", "clone", gitRepo, projectDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone失败: %v", err)
	}

	return nil
}

// parseGitRepository 解析Git仓库，查找playbook文件
func (s *TaskAnsibleServiceImpl) parseGitRepository(projectDir string) ([]string, error) {
	var playbookFiles []string

	// 递归查找.yml和.yaml文件作为playbook
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过.git目录
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// 查找playbook文件(根目录下的.yml/.yaml文件，排除vars和roles目录中的文件)
		if !info.IsDir() {
			relPath, _ := filepath.Rel(projectDir, path)
			// 只考虑根目录下的yml/yaml文件作为playbook
			if !strings.Contains(relPath, "/") && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
				// 排除常见的非playbook文件
				if info.Name() != "hosts" && info.Name() != "ansible.cfg" &&
					!strings.HasPrefix(info.Name(), "group_vars") &&
					!strings.HasPrefix(info.Name(), "host_vars") &&
					info.Name() != "vars.yml" && info.Name() != "vars.yaml" {
					playbookFiles = append(playbookFiles, info.Name())
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析目录失败: %v", err)
	}

	if len(playbookFiles) == 0 {
		return nil, fmt.Errorf("未在仓库根目录找到有效的playbook文件(.yml/.yaml)")
	}

	// 按文件名排序，确保执行顺序的可预测性
	for i := 0; i < len(playbookFiles); i++ {
		for j := i + 1; j < len(playbookFiles); j++ {
			if playbookFiles[i] > playbookFiles[j] {
				playbookFiles[i], playbookFiles[j] = playbookFiles[j], playbookFiles[i]
			}
		}
	}

	return playbookFiles, nil
}

// createSubTasksFromPlaybooks 根据playbook文件创建子任务
func (s *TaskAnsibleServiceImpl) createSubTasksFromPlaybooks(taskID uint, projectDir string, playbookFiles []string) error {
	for _, fileName := range playbookFiles {
		filePath := filepath.Join(projectDir, fileName)

		// 创建子任务记录
		subTask := &taskmodel.TaskAnsibleWork{
			TaskID:        taskID,
			EntryFileName: fileName,
			EntryFilePath: filePath,
			Status:        1, // 等待中
		}

		if err := s.dao.DB.Create(subTask).Error; err != nil {
			return fmt.Errorf("创建子任务失败: %v", err)
		}
	}

	// 更新任务数量到数据库
	if err := s.dao.DB.Model(&taskmodel.TaskAnsible{}).Where("id = ?", taskID).Update("task_count", len(playbookFiles)).Error; err != nil {
		return fmt.Errorf("更新任务数量失败: %v", err)
	}

	return nil
}

// saveGlobalVariables 保存全局变量
func (s *TaskAnsibleServiceImpl) saveGlobalVariables(c *gin.Context, task *taskmodel.TaskAnsible, projectDir string, variables map[string]string) error {
	if len(variables) == 0 {
		return nil
	}

	varsFile := filepath.Join(projectDir, "vars/all.yml")
	if err := os.MkdirAll(filepath.Dir(varsFile), 0755); err != nil {
		return fmt.Errorf("创建变量目录失败: %v", err)
	}

	// 构建YAML内容
	var content strings.Builder
	for k, v := range variables {
		content.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	// 写入文件
	if err := os.WriteFile(varsFile, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("保存全局变量失败: %v", err)
	}

	// 保存到数据库
	task.GlobalVars = toJSON(variables)
	if err := s.dao.DB.Model(&taskmodel.TaskAnsible{}).Where("id = ?", task.ID).Update("global_vars", task.GlobalVars).Error; err != nil {
		return fmt.Errorf("更新全局变量失败: %v", err)
	}

	return nil
}

// generateInventoryFile 生成inventory文件
func (s *TaskAnsibleServiceImpl) generateInventoryFile(c *gin.Context, projectDir string, hostInfos *HostSSHInfoCollection) error {
	inventory := hostInfos.GenerateInventory()
	inventoryPath := filepath.Join(projectDir, "hosts")
	if err := os.WriteFile(inventoryPath, []byte(inventory), 0644); err != nil {
		return fmt.Errorf("生成inventory文件失败: %v", err)
	}
	return nil
}

// GetTasksByName 根据名称模糊查询任务
func (s *TaskAnsibleServiceImpl) GetTasksByName(c *gin.Context, name string) {
	tasks, err := s.dao.GetByName(name)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("查询任务失败: %v", err))
		return
	}
	result.Success(c, gin.H{"data": tasks, "total": len(tasks)})
}

// GetTasksByType 根据类型查询任务
func (s *TaskAnsibleServiceImpl) GetTasksByType(c *gin.Context, taskType int) {
	tasks, err := s.dao.GetByType(taskType)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("查询任务失败: %v", err))
		return
	}
	result.Success(c, gin.H{"data": tasks, "total": len(tasks)})
}

// CreateK8sTask 创建K8s部署任务
func (s *TaskAnsibleServiceImpl) CreateK8sTask(c *gin.Context, req *CreateK8sTaskRequest) {
	// 1. 验证主机ID并获取主机信息
	hostInfos, err := s.getK8sHostInfo(req)
	if err != nil {
		result.Failed(c, 500, fmt.Sprintf("获取主机信息失败: %v", err))
		return
	}

	// 2. 创建任务主记录
	task := &model.TaskAnsible{
		Name:        req.Name,
		Description: req.Description,
		Type:        3, // K8s任务类型
		GitRepo:     "git@gitee.com:zhang_fan1024/zf-k8s.git", // 固定的K8s Git仓库
		HostGroups:  s.buildK8sHostGroups(req),
		AllHostIDs:  s.buildK8sAllHostIDs(req),
		GlobalVars:  s.buildK8sGlobalVars(req),
		Status:      1, // 等待中
		TaskCount:   1, // K8s任务固定为1个
	}

	// 3. 保存任务到数据库
	if err := s.dao.Create(task); err != nil {
		result.Failed(c, 500, fmt.Sprintf("创建任务失败: %v", err))
		return
	}

	// 4. 创建项目目录
	projectDir := fmt.Sprintf("./task/%d/%s", task.ID, task.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		s.updateTaskErrorStatus(task.ID, err)
		result.Failed(c, 500, fmt.Sprintf("创建项目目录失败: %v", err))
		return
	}

	// 5. 克隆K8s Git仓库
	if err := s.cloneGitRepository(task.GitRepo, projectDir); err != nil {
		s.updateTaskErrorStatus(task.ID, err)
		result.Failed(c, 500, fmt.Sprintf("克隆K8s仓库失败: %v", err))
		return
	}

	// 6. 生成config.json配置文件
	if err := s.generateK8sConfig(projectDir, req, hostInfos); err != nil {
		s.updateTaskErrorStatus(task.ID, err)
		result.Failed(c, 500, fmt.Sprintf("生成K8s配置失败: %v", err))
		return
	}

	// 7. 创建子任务（K8s部署脚本）
	if err := s.createK8sSubTask(task.ID, projectDir); err != nil {
		s.updateTaskErrorStatus(task.ID, err)
		result.Failed(c, 500, fmt.Sprintf("创建K8s子任务失败: %v", err))
		return
	}

	// 8. 更新任务状态
	s.dao.DB.Model(task).Update("status", 1) // 等待中

	result.Success(c, gin.H{"data": task})
}

// getK8sHostInfo 获取K8s节点主机信息
func (s *TaskAnsibleServiceImpl) getK8sHostInfo(req *CreateK8sTaskRequest) (map[string][]K8sNodeInfo, error) {
	hostInfos := map[string][]K8sNodeInfo{
		"masters": {},
		"workers": {},
		"etcd":    {},
	}

	// 获取所有唯一主机ID
	allHostIDs := make(map[uint]bool)
	for _, id := range req.MasterHostIDs {
		allHostIDs[id] = true
	}
	for _, id := range req.WorkerHostIDs {
		allHostIDs[id] = true
	}
	for _, id := range req.EtcdHostIDs {
		allHostIDs[id] = true
	}

	// 转换为切片
	var hostIDs []uint
	for id := range allHostIDs {
		hostIDs = append(hostIDs, id)
	}

	// 从数据库查询主机信息
	var hosts []cmdbmodel.CmdbHost
	if err := s.dao.DB.Where("id IN ?", hostIDs).Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("查询主机信息失败: %v", err)
	}

	// 构建主机信息映射
	hostMap := make(map[uint]K8sNodeInfo)
	for _, host := range hosts {
		// 获取SSH认证信息
		var password string
		if host.SSHKeyID != 0 {
			var ecsAuth configcentermodel.EcsAuth
			s.dao.DB.Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&ecsAuth)
			// 只有密码认证时才设置password，其他类型保持空字符串
			if ecsAuth.Type == 1 {
				password = ecsAuth.Password
			}
			// 注意：K8s任务当前只支持密码认证，type=2和type=3需要其他处理方式
		}

		hostMap[host.ID] = K8sNodeInfo{
			Name:     host.HostName,
			IP:       host.SSHIP, // 使用SSH IP
			User:     host.SSHName,
			Password: password,
		}
	}

	// 分配主机到不同角色
	for _, id := range req.MasterHostIDs {
		if info, exists := hostMap[id]; exists {
			hostInfos["masters"] = append(hostInfos["masters"], info)
		}
	}
	for _, id := range req.WorkerHostIDs {
		if info, exists := hostMap[id]; exists {
			hostInfos["workers"] = append(hostInfos["workers"], info)
		}
	}
	for _, id := range req.EtcdHostIDs {
		if info, exists := hostMap[id]; exists {
			hostInfos["etcd"] = append(hostInfos["etcd"], info)
		}
	}

	return hostInfos, nil
}

// buildK8sHostGroups 构建K8s主机分组JSON
func (s *TaskAnsibleServiceImpl) buildK8sHostGroups(req *CreateK8sTaskRequest) string {
	hostGroups := map[string][]uint{
		"masters": req.MasterHostIDs,
		"workers": req.WorkerHostIDs,
		"etcd":    req.EtcdHostIDs,
	}
	data, _ := json.Marshal(hostGroups)
	return string(data)
}

// buildK8sAllHostIDs 构建所有主机ID的JSON数组
func (s *TaskAnsibleServiceImpl) buildK8sAllHostIDs(req *CreateK8sTaskRequest) string {
	allIDs := make(map[uint]bool)
	for _, id := range req.MasterHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.WorkerHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.EtcdHostIDs {
		allIDs[id] = true
	}

	var ids []uint
	for id := range allIDs {
		ids = append(ids, id)
	}

	data, _ := json.Marshal(ids)
	return string(data)
}

// buildK8sGlobalVars 构建K8s全局变量JSON
func (s *TaskAnsibleServiceImpl) buildK8sGlobalVars(req *CreateK8sTaskRequest) string {
	vars := map[string]string{
		"cluster_name":       req.ClusterName,
		"cluster_version":    req.ClusterVersion,
		"deployment_mode":    fmt.Sprintf("%d", req.DeploymentMode),
		"enabled_components": strings.Join(req.EnabledComponents, ","),
	}

	if req.PrivateRegistry != "" {
		vars["private_registry"] = req.PrivateRegistry
		vars["registry_username"] = req.RegistryUsername
		vars["registry_password"] = req.RegistryPassword
	}

	data, _ := json.Marshal(vars)
	return string(data)
}

// generateK8sConfig 生成K8s配置文件(config.json)
func (s *TaskAnsibleServiceImpl) generateK8sConfig(projectDir string, req *CreateK8sTaskRequest, hostInfos map[string][]K8sNodeInfo) error {
	config := K8sConfigJSON{}

	// 集群配置
	config.Cluster.Name = req.ClusterName
	config.Cluster.Version = req.ClusterVersion
	config.Cluster.DeploymentMode = req.DeploymentMode

	// 节点配置
	config.Nodes.Masters = hostInfos["masters"]
	config.Nodes.Workers = hostInfos["workers"]
	config.Nodes.Etcd = hostInfos["etcd"]

	// 组件配置
	if len(req.EnabledComponents) > 0 {
		config.Components.Enabled = req.EnabledComponents
	} else {
		config.Components.Enabled = []string{"calico", "coredns"}
	}

	// 仓库配置
	if req.PrivateRegistry != "" {
		config.Registry.PrivateRegistry = req.PrivateRegistry
		config.Registry.Username = req.RegistryUsername
		config.Registry.Password = req.RegistryPassword
	}

	// 写入config.json文件
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	configPath := filepath.Join(projectDir, "config.json")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// createK8sSubTask 创建K8s子任务
func (s *TaskAnsibleServiceImpl) createK8sSubTask(taskID uint, projectDir string) error {
	// K8s任务只有一个子任务：执行部署脚本
	work := &model.TaskAnsibleWork{
		TaskID:        taskID,
		EntryFileName: "deploy-simple.sh",
		EntryFilePath: filepath.Join(projectDir, "scripts", "deploy-simple.sh"),
		LogPath:       filepath.Join("./logs/ansible", fmt.Sprintf("%d", taskID), "deploy-k8s.log"),
		Status:        1, // 等待中
	}

	// 确保日志目录存在
	logDir := filepath.Dir(work.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 保存子任务
	if err := s.dao.DB.Create(work).Error; err != nil {
		return fmt.Errorf("创建子任务失败: %v", err)
	}

	return nil
}
