package service

import (
	"fmt"
	cmdbmodel "dodevops-api/api/cmdb/model"
	configcentermodel "dodevops-api/api/configcenter/model"
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type ITaskWorkService interface {
	StartJob(taskID uint) error  // 启动任务
	StopJob(taskID, templateID uint) error  // 停止任务
	GetJobLog(taskID, templateID uint) (string, error) // 获取任务日志
	GetJobStatus(taskID, templateID uint) (*model.TaskWorkStatus, error) // 获取任务状态
	ScheduleJob(job *model.TaskWork) error // 添加任务
}

type TaskWorkServiceImpl struct {
	dao dao.TaskWorkDaoInterface
	taskJobService TaskJobService
}

func (s *TaskWorkServiceImpl) StartJob(taskID uint) error {
	// 获取任务信息
	jobs, err := s.dao.GetByTaskID(taskID)
	if err != nil {
		return fmt.Errorf("获取任务失败: %v", err)
	}

	// 检查是否有任务正在运行
	for _, job := range jobs {
		if job.Status == 2 { // 2表示运行中
			return fmt.Errorf("任务正在运行中")
		}
	}

	// 获取任务队列服务
	taskQueue := GetTaskQueue()
	if taskQueue == nil {
		return fmt.Errorf("任务队列服务未初始化")
	}

	// 将所有任务提交到 Redis 队列
	submittedCount := 0
	for _, job := range jobs {
		if job.Status != 1 { // 非等待状态跳过
			continue
		}

		// 更新状态为运行中
		if err := s.dao.UpdateStatus(job.ID, 2); err != nil {
			continue
		}

		// 获取任务优先级（可根据任务类型或其他条件决定）
		priority := s.getTaskPriority(&job)

		// 提交到 Redis 队列
		if err := taskQueue.Enqueue(&job, priority); err != nil {
			// 回滚状态
			s.dao.UpdateStatus(job.ID, 1)
			continue
		}

		submittedCount++
	}

	if submittedCount == 0 {
		return fmt.Errorf("没有有效的任务可以执行")
	}

	return nil
}

// getTaskPriority 获取任务优先级
func (s *TaskWorkServiceImpl) getTaskPriority(job *model.TaskWork) string {
	// 获取父任务信息
	parentTask, err := s.dao.GetParentTask(job.TaskID)
	if err != nil {
		return "normal" // 默认优先级
	}

	// 根据任务类型确定优先级
	switch parentTask.Type {
	case model.TaskTypeImmediate: // 立即执行任务 - 高优先级
		return "high"
	case model.TaskTypeScheduled: // 定时任务 - 普通优先级
		return "normal"
	case model.TaskTypeAnsible: // Ansible任务 - 低优先级（可能耗时较长）
		return "low"
	default:
		return "normal"
	}
}

func (s *TaskWorkServiceImpl) executeJob(job *model.TaskWork) error {
	startTime := time.Now()

	// 检查父任务状态，如果已停止或暂停则不再执行
	parentTask, err := s.dao.GetParentTask(job.TaskID)
	if err != nil {
		return fmt.Errorf("获取父任务状态失败: %v", err)
	}

	// 检查任务是否被停止
	if parentTask.Status == model.TaskStatusFailed {
		// 更新子任务状态为异常以匹配父任务状态
		if err := s.dao.UpdateStatus(job.ID, model.TaskStatusFailed); err != nil {
			return fmt.Errorf("更新子任务状态失败: %v", err)
		}
		return fmt.Errorf("父任务已被停止，子任务标记为异常")
	}

	// 检查任务是否被暂停（仅对定时任务有效）
	if parentTask.Status == model.TaskStatusPaused && parentTask.Type == model.TaskTypeScheduled {
		// 暂停状态的任务不执行，但不标记为异常
		return fmt.Errorf("定时任务已被暂停，跳过执行")
	}

	// 获取模板和主机信息
	template, host, err := s.getJobDetails(job)
	if err != nil {
		// 仅当获取关键信息失败时才标记为异常
		s.dao.UpdateStatus(job.ID, 4)
		return fmt.Errorf("获取任务详情失败: %v", err)
	}

	// 创建日志目录
	logDir := fmt.Sprintf("logs/task_%d", job.TaskID)
	os.MkdirAll(logDir, 0755)

	// 执行SSH任务（带重试和严格错误处理）
	var logContent string
	var sshErr error
	maxRetries := 3               // 增加重试次数
	retryDelay := 5 * time.Second // 增加重试间隔

	for i := 0; i <= maxRetries; i++ {
		logContent, sshErr = s.executeSSHTask(job, template, host)
		if sshErr == nil {
			break
		}

		if isSSHConnectionError(sshErr) {
			if i < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			// 最终连接失败，必须标记为异常
			if updateErr := s.dao.UpdateStatus(job.ID, 4); updateErr != nil {
				return fmt.Errorf("SSH连接失败且状态更新失败: %v (原错误: %v)", updateErr, sshErr)
			}
			// 同时更新任务日志记录失败原因
			s.dao.UpdateLog(job.ID, fmt.Sprintf("SSH连接失败: %v", sshErr))
			return fmt.Errorf("SSH连接失败: %v", sshErr)
		} else {
			// 命令执行问题不标记为异常，但记录详细日志
			logContent = fmt.Sprintf("命令执行出错: %v\n%s", sshErr, logContent)
			break
		}
	}

	// 写入日志文件
	logPath := filepath.Join(logDir, fmt.Sprintf("task_%d_template_%d.log", job.TaskID, job.TemplateID))
	if err := os.WriteFile(logPath, []byte(logContent), 0644); err != nil {
		return fmt.Errorf("写入日志文件失败: %v", err)
	}

	// 记录任务结束时间和耗时
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	if err := s.dao.UpdateTiming(job.ID, &startTime, &endTime, duration); err != nil {
		return fmt.Errorf("记录任务时间失败: %v", err)
	}

	// 更新任务状态和日志 - 确保这些更新不覆盖时间字段
	if err := s.dao.UpdateLog(job.ID, logContent); err != nil {
		return err
	}
	if err := s.dao.UpdateLogPath(job.ID, logPath); err != nil {
		return err
	}

	// 更新任务为完成状态
	if err := s.dao.UpdateStatus(job.ID, 3); err != nil {
		return err
	}

	// 获取该任务的所有子任务
	jobs, err := s.dao.GetByTaskID(job.TaskID)
	if err != nil {
		return fmt.Errorf("获取子任务失败: %v", err)
	}

	// 获取父任务类型
	taskDao := dao.NewTaskDao(common.GetDB())
	parentTask, err = taskDao.GetById(job.TaskID)
	if err != nil {
		return fmt.Errorf("获取父任务信息失败: %v", err)
	}

	// 计算总耗时和确定最终状态
	var totalDuration int
	finalStatus := 3 // 默认成功
	allCompleted := true
	now := time.Now()

	for _, j := range jobs {
		if j.Duration > 0 {
			totalDuration += j.Duration
		}
		if j.Status == 4 { // 如果有任何子任务失败
			finalStatus = 4
		}
		if j.Status == 1 || j.Status == 2 { // 如果有子任务未完成
			allCompleted = false
		}
	}

	// 只有所有子任务都完成时才更新父任务
	if allCompleted {
		// 对于定时任务，父任务状态保持为运行中(2)
		if parentTask.Type == 2 { // 2表示定时任务
			finalStatus = 2

			// 将所有子任务状态重置为等待中(1)，以便下次定时触发时可以再次执行
			for _, j := range jobs {
				if j.Status == 3 || j.Status == 4 { // 成功或失败状态
					s.dao.UpdateStatus(j.ID, 1)
				}
			}
		}

		// 更新父任务状态、耗时和结束时间（使用Select只更新必要字段）
		task := &model.Task{
			ID:       job.TaskID,
			Status:   finalStatus,
			Duration: totalDuration,
			EndTime:  &now,
		}

		// 使用Select明确指定要更新的字段，避免created_at被错误更新
		if err := common.GetDB().Model(task).
			Select("status", "duration", "end_time").
			Where("id = ?", job.TaskID).
			Updates(task).Error; err != nil {
			return fmt.Errorf("更新父任务失败: %v", err)
		}

		// 更新execute_count和next_run_time（对所有类型的任务）
		s.taskJobService.UpdateTaskAfterExecution(job.TaskID, parentTask.CronExpr)
	}

	return nil
}

func (s *TaskWorkServiceImpl) StopJob(taskID, templateID uint) error {
	// 1. 获取任务详情并检查状态
	job, err := s.dao.GetByTaskAndTemplateID(taskID, templateID)
	if err != nil {
		return fmt.Errorf("获取任务失败: %v", err)
	}

	// 2. 如果任务不是运行中状态(2)，直接返回错误
	if job.Status != 2 {
		return fmt.Errorf("任务当前状态为%d，无法停止非运行中的任务", job.Status)
	}

	// 3. 执行远程停止操作
	s.killRemoteProcess(job)

	// 4. 更新状态为异常(4)并记录结束时间
	now := time.Now()
	if err := s.dao.UpdateStatus(job.ID, 4); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	// 5. 更新结束时间和耗时
	s.dao.UpdateTiming(job.ID, nil, &now, 0)

	// 6. 停止所有关联的子任务
	jobs, err := s.dao.GetByTaskID(taskID)
	if err == nil {
		// 首先停止当前任务
		s.dao.UpdateStatus(job.ID, 4)

		// 停止所有子任务
		for _, j := range jobs {
			if j.ID == job.ID {
				continue // 跳过当前任务，已经处理过
			}

			// 只停止运行中或等待中的任务
			if j.Status == 1 || j.Status == 2 {
				s.dao.UpdateStatus(j.ID, 4)
			}
		}

		// 更新父任务状态为已停止
		taskDao := dao.NewTaskDao(common.GetDB())
		taskDao.UpdateStatus(taskID, 4)
	}

	// 7. 记录操作日志
	logContent := fmt.Sprintf("任务被手动停止于 %s", now.Format(time.RFC3339))
	s.dao.UpdateLog(job.ID, logContent)

	return nil
}

func (s *TaskWorkServiceImpl) GetJobLog(taskID, templateID uint) (string, error) {
	job, err := s.dao.GetByTaskAndTemplateID(taskID, templateID)
	if err != nil {
		return "", fmt.Errorf("获取任务失败: %v", err)
	}

	if job.LogPath == "" {
		return "", fmt.Errorf("任务日志路径为空")
	}

	content, err := os.ReadFile(job.LogPath)
	if err != nil {
		return "", fmt.Errorf("读取日志文件失败: %v", err)
	}

	return string(content), nil
}

func (s *TaskWorkServiceImpl) GetJobStatus(taskID, templateID uint) (*model.TaskWorkStatus, error) {
	job, err := s.dao.GetByTaskAndTemplateID(taskID, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %v", err)
	}

	return &model.TaskWorkStatus{
		Status:    job.Status,
		Log:       job.Log,
		StartTime: job.StartTime.Format(time.RFC3339),
		EndTime:   job.EndTime.Format(time.RFC3339),
		Duration:  job.Duration,
	}, nil
}

func (s *TaskWorkServiceImpl) ScheduleJob(job *model.TaskWork) error {
	// 获取关联的TaskJob记录
	taskJob, err := s.dao.GetTaskJobByID(job.TaskID)
	if err != nil {
		return fmt.Errorf("获取任务信息失败: %v", err)
	}

	// 立即执行任务 (普通任务)
	if taskJob.Type == 1 {
		// 更新状态为运行中
		if err := s.dao.UpdateStatus(job.ID, 2); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		// 获取任务队列服务
		taskQueue := GetTaskQueue()
		if taskQueue == nil {
			// 如果队列服务不可用，回退到原始方式
			go func() {
				s.executeJob(job)
			}()
			return nil
		}

		// 提交到 Redis 队列
		priority := s.getTaskPriority(job)
		if err := taskQueue.Enqueue(job, priority); err != nil {
			return fmt.Errorf("提交任务到队列失败: %v", err)
		}

		return nil
	}

	// 定时任务处理
	if taskJob.Type == 2 {
		// 从task_job获取cron表达式并计算下次执行时间
		if taskJob.CronExpr == "" {
			return fmt.Errorf("定时任务必须包含cron表达式")
		}

		// 使用cron库解析cron表达式（标准5位格式：分 时 日 月 周）
		schedule, err := cron.ParseStandard(taskJob.CronExpr)
		if err != nil {
			return fmt.Errorf("解析cron表达式失败: %v", err)
		}

		// 计算下次执行时间
		nextTime := schedule.Next(time.Now())
		if nextTime.IsZero() {
			return fmt.Errorf("无法计算有效的下次执行时间")
		}

		job.ScheduledTime = &nextTime

		// 验证执行时间格式
		if job.ScheduledTime.Before(time.Now()) {
			return fmt.Errorf("定时任务执行时间(%v)不能早于当前时间(%v)", 
				job.ScheduledTime.Format("2006-01-02 15:04:05"),
				time.Now().Format("2006-01-02 15:04:05"))
		}

		// 计算延迟时间
		delay := time.Until(*job.ScheduledTime)
		if delay <= 0 {
			return fmt.Errorf("定时任务执行时间已过期，请设置未来的时间")
		}

		// 先注册到全局 cron 调度器（在数据库操作之前，避免死锁）
		scheduler := GetGlobalScheduler()
		if scheduler == nil {
			return fmt.Errorf("全局调度器未初始化")
		}

		// 检查任务是否已在调度器中（防止重复注册）
		entries := scheduler.GetEntries()
		if _, exists := entries[job.TaskID]; exists {
			// 任务已注册，只更新数据库状态即可
			taskDao := dao.NewTaskDao(common.GetDB())
			if err := taskDao.UpdateStatus(job.TaskID, 2); err != nil {
				return fmt.Errorf("更新父任务状态失败: %v", err)
			}
			return nil
		}

		// 注册定时任务（使用统一的回调函数）
		if err := scheduler.AddScheduledTask(job.TaskID, taskJob.CronExpr, createScheduledTaskHandler(job.TaskID)); err != nil {
			return fmt.Errorf("注册定时任务失败: %v", err)
		}

		// 注册成功后，再更新数据库状态（避免在持有锁的情况下操作数据库）
		taskDao := dao.NewTaskDao(common.GetDB())
		if err := taskDao.UpdateStatus(job.TaskID, 2); err != nil {
			// 如果状态更新失败，移除已注册的任务
			scheduler.RemoveScheduledTask(job.TaskID)
			return fmt.Errorf("更新父任务状态失败: %v", err)
		}

		// 更新子任务的计划执行时间
		if err := s.dao.UpdateScheduledJob(job.ID, 1, job.ScheduledTime); err != nil {
			return fmt.Errorf("更新子任务状态失败: %v", err)
		}

		return nil
	}

	return fmt.Errorf("不支持的任务类型: %d (有效类型: 1=普通任务,2=定时任务,3=ansible任务)", taskJob.Type)
}

// Helper methods
func (s *TaskWorkServiceImpl) getJobDetails(job *model.TaskWork) (*model.TaskTemplate, *cmdbmodel.CmdbHost, error) {
	// 获取模板信息
	var template model.TaskTemplate
	if err := common.GetDB().First(&template, job.TemplateID).Error; err != nil {
		return nil, nil, fmt.Errorf("获取模板失败: %v", err)
	}

	// 获取主机信息
	var host cmdbmodel.CmdbHost
	if err := common.GetDB().First(&host, job.HostID).Error; err != nil {
		return nil, nil, fmt.Errorf("获取主机失败: %v", err)
	}

	return &template, &host, nil
}
// 
func isSSHConnectionError(err error) bool {
	// 判断是否为SSH连接错误(包括超时)
	return err != nil && (strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "authentication") ||
		strings.Contains(err.Error(), "timeout"))
}
// 
func (s *TaskWorkServiceImpl) executeSSHTask(job *model.TaskWork, template *model.TaskTemplate, host *cmdbmodel.CmdbHost) (string, error) {
	// 1. 获取认证凭证
	var ecsAuth configcentermodel.EcsAuth
	if err := common.GetDB().First(&ecsAuth, host.SSHKeyID).Error; err != nil {
		return "", fmt.Errorf("获取认证凭证失败: %v", err)
	}

	// 2. 初始化SSH配置
	sshUtil := util.NewSSHUtil()
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,   // SSH连接IP
		Port:     host.SSHPort, // SSH端口
		Username: host.SSHName, // SSH用户名
		Type:     ecsAuth.Type, // 认证类型
	}

	// 3. 根据认证类型配置
	switch ecsAuth.Type {
	case 1: // 密码认证
		sshConfig.Password = ecsAuth.Password
	case 2: // 密钥认证
		sshConfig.PublicKey = ecsAuth.PublicKey
	case 3: // 公钥免认证 - 自动查找本地私钥
		// 不需要设置额外信息，SSHUtil会自动查找本地私钥
	default:
		return "", fmt.Errorf("不支持的认证类型: %d", ecsAuth.Type)
	}

	// 4. 构造PID文件路径
	pidFile := fmt.Sprintf("/tmp/task_%d_%d.pid", job.TaskID, job.TemplateID)

	// 5. 创建临时脚本文件并执行
	scriptFile := fmt.Sprintf("/tmp/task_%d_%d.sh", job.TaskID, job.TemplateID)
	cmd := fmt.Sprintf(`
		# 创建脚本文件
		cat > %s << 'EOF'
#!/bin/bash
echo "[$(date '+%%F-%%H:%%M:%%S')] 任务开始"
%s
echo "[$(date '+%%F-%%H:%%M:%%S')] 任务完成"
EOF
		# 设置执行权限
		chmod +x %s
		# 验证脚本文件内容
		if ! grep -q "#!/bin/bash" %s; then
			echo "脚本文件内容验证失败"
			exit 1
		fi
		# 执行并记录PID
		nohup %s > /tmp/task_%d_%d.log 2>&1 & echo $! > %s
		# 验证PID文件
		sleep 0.5
		if [ ! -f %s ]; then
			echo "PID文件创建失败"
			exit 1
		fi
	`, scriptFile, template.Content, scriptFile, scriptFile,
		scriptFile, job.TaskID, job.TemplateID, pidFile, pidFile)

	// 6. 执行远程命令
	_, err := sshUtil.ExecuteRemoteCommand(sshConfig, cmd)
	if err != nil {
		return "", fmt.Errorf("SSH执行失败: %v", err)
	}

	// 8. 验证PID文件是否创建 (带重试)
	maxRetries := 3
	retryDelay := 1 * time.Second
	var pid string

	for i := 0; i < maxRetries; i++ {
		checkCmd := fmt.Sprintf(`
			if [ -f %s ]; then
				cat %s
				exit 0
			else
				exit 1
			fi
		`, pidFile, pidFile)

		checkOutput, err := sshUtil.ExecuteRemoteCommand(sshConfig, checkCmd)
		if err == nil && checkOutput != "" {
			pid = strings.TrimSpace(checkOutput)
			break
		}
		time.Sleep(retryDelay)
	}

	if pid == "" {
		return "", fmt.Errorf("无法获取有效的进程PID")
	}

	// 9. 验证进程是否存在 (宽松检查，即使进程已结束也继续执行)
	verifyCmd := fmt.Sprintf(`
		if ps -p %s > /dev/null 2>&1; then
			echo "进程%s验证成功"
		else
			echo "进程%s已结束，继续执行脚本"
		fi
		exit 0
	`, pid, pid, pid)

	sshUtil.ExecuteRemoteCommand(sshConfig, verifyCmd)

	// 10. 等待任务完成并获取完整日志
	waitCmd := fmt.Sprintf(`
		while ps -p %s > /dev/null 2>&1; do
			sleep 1
		done
		echo "任务完成"
	`, pid)

	if _, err := sshUtil.ExecuteRemoteCommand(sshConfig, waitCmd); err != nil {
		return "", fmt.Errorf("等待任务完成失败: %v", err)
	}

	// 获取完整执行日志
	finalLogCmd := fmt.Sprintf("cat /tmp/task_%d_%d.log", job.TaskID, job.TemplateID)
	finalOutput, err := sshUtil.ExecuteRemoteCommand(sshConfig, finalLogCmd)
	if err != nil {
		return "", fmt.Errorf("获取完整日志失败: %v", err)
	}

	// 清理临时文件
	cleanCmd := fmt.Sprintf("rm -f /tmp/task_%d_%d.sh /tmp/task_%d_%d.pid",
		job.TaskID, job.TemplateID, job.TaskID, job.TemplateID)
	sshUtil.ExecuteRemoteCommand(sshConfig, cleanCmd)

	return finalOutput, nil
}
// 杀死远程进程
func (s *TaskWorkServiceImpl) killRemoteProcess(job *model.TaskWork) error {
	// 1. 获取主机信息和SSH认证
	_, host, err := s.getJobDetails(job)
	if err != nil {
		return fmt.Errorf("获取主机信息失败: %v", err)
	}

	// 2. 获取SSH认证凭证
	var ecsAuth configcentermodel.EcsAuth
	if err := common.GetDB().First(&ecsAuth, host.SSHKeyID).Error; err != nil {
		return fmt.Errorf("获取SSH认证凭证失败: %v", err)
	}

	// 3. 验证SSH认证信息
	if ecsAuth.Type < 1 || ecsAuth.Type > 3 {
		return fmt.Errorf("不支持的SSH认证类型: %d (支持1-密码、2-密钥、3-公钥免认证)", ecsAuth.Type)
	}
	if ecsAuth.Type == 1 && ecsAuth.Password == "" {
		return fmt.Errorf("密码认证类型但密码为空")
	}
	if ecsAuth.Type == 2 && ecsAuth.PublicKey == "" {
		return fmt.Errorf("密钥认证类型但公钥为空")
	}

	// 4. 初始化SSH配置 (与executeSSHTask保持一致)
	sshUtil := util.NewSSHUtil()
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,   // SSH连接IP
		Port:     host.SSHPort, // SSH端口
		Username: host.SSHName, // SSH用户名
		Type:     ecsAuth.Type, // 认证类型
	}

	// 根据认证类型配置
	switch ecsAuth.Type {
	case 1: // 密码认证
		sshConfig.Password = ecsAuth.Password
	case 2: // 密钥认证
		sshConfig.PublicKey = ecsAuth.PublicKey
	case 3: // 公钥免认证 - 自动查找本地私钥
		// 不需要设置额外信息，SSHUtil会自动查找本地私钥
	default:
		return fmt.Errorf("不支持的认证类型: %d", ecsAuth.Type)
	}

	// 5. 构造PID文件路径和终止命令
	pidFile := fmt.Sprintf("/tmp/task_%d_%d.pid", job.TaskID, job.TemplateID)

	// 更健壮的终止命令:
	cmd := fmt.Sprintf(`
		# 检查PID文件是否存在
		if [ ! -f %s ]; then
			echo "PID文件 %s 不存在";
			exit 0;
		fi
		
		# 读取PID
		PID=$(cat %s);
		
		# 检查进程是否存在
		if ! ps -p $PID > /dev/null 2>&1; then
			echo "进程 $PID 不存在";
			rm -f %s;
			exit 0;
		fi
		
		# 终止进程
		echo "正在终止进程 $PID";
		if kill -9 $PID; then
			# 等待进程终止
			for i in {1..5}; do
				if ! ps -p $PID > /dev/null 2>&1; then
					echo "成功终止进程 $PID";
					rm -f %s;
					exit 0;
				fi
				sleep 1;
			done
			
			# 如果进程仍然存在
			if ps -p $PID > /dev/null 2>&1; then
				echo "警告: 进程 $PID 仍然运行";
				exit 1;
			fi
		else
			echo "终止进程 $PID 失败";
			exit 1;
		fi`,
		pidFile, pidFile, pidFile, pidFile, pidFile)

	// 6. 执行终止命令
	output, err := sshUtil.ExecuteRemoteCommand(sshConfig, cmd)
	if err != nil {

		// 检查是否是连接/认证错误
		if strings.Contains(err.Error(), "authentication") ||
			strings.Contains(err.Error(), "connection") {
			return fmt.Errorf("SSH连接/认证失败: %v", err)
		}

		// 检查输出中是否有成功终止的迹象
		if strings.Contains(output, "成功终止进程") ||
			strings.Contains(output, "进程不存在") ||
			strings.Contains(output, "PID文件不存在") {
			return nil
		}
		return fmt.Errorf("终止进程失败: %v\n输出:\n%s", err, output)
	}

	return nil
}

// 辅助函数：获取任务类型名称
func getJobType(jobType int) string {
	switch jobType {
	case 1:
		return "立即执行"
	case 2:
		return "定时任务"
	default:
		return "未知类型"
	}
}

// 创建任务服务实例
func NewTaskWorkService() ITaskWorkService {
	return &TaskWorkServiceImpl{
		dao: dao.TaskWorkDao(),
		taskJobService: NewTaskService(common.GetDB()),
	}
}
