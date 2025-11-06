package dao

import (
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

// AgentDao Agent数据访问对象接口 (Linux主机专用)
type AgentDao interface {
	Create(agent *model.Agent) error
	GetByID(id uint) (*model.Agent, error)
	GetByHostID(hostID uint) (*model.Agent, error)
	Update(agent *model.Agent) error
	UpdateStatus(id uint, status int, errorMsg string) error
	UpdateHeartbeat(hostID uint, heartbeat *model.AgentHeartbeatDto) error
	UpdateHeartbeatByIP(hostID uint, heartbeat *model.AgentHeartbeatDto) error // 通过IP更新心跳
	UpdateByHostID(hostID uint, updates map[string]interface{}) error
	UpdateInstallProgress(hostID uint, progress int) error // 新增：更新安装进度
	Delete(id uint) error
	DeleteByHostID(hostID uint) error // 通过主机ID删除Agent记录
	GetList(dto *model.AgentListDto) ([]model.Agent, int64, error)
	GetAllRunning() ([]model.Agent, error)
	GetOfflineAgents(duration time.Duration) ([]model.Agent, error)
	CreateOrUpdate(agent *model.Agent) error
	GetAgentStatistics() (map[string]int64, error)
}

// AgentDaoImpl Agent数据访问对象实现
type AgentDaoImpl struct {
	db *gorm.DB
}

// NewAgentDao 创建Agent DAO实例
func NewAgentDao() AgentDao {
	return &AgentDaoImpl{
		db: common.GetDB(),
	}
}

// Create 创建Agent记录
func (d *AgentDaoImpl) Create(agent *model.Agent) error {
	agent.CreateTime.Time = time.Now()
	agent.UpdateTime.Time = time.Now()
	return d.db.Create(agent).Error
}

// GetByID 根据ID获取Agent
func (d *AgentDaoImpl) GetByID(id uint) (*model.Agent, error) {
	var agent model.Agent
	err := d.db.Where("id = ?", id).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByHostID 根据主机ID获取Agent
func (d *AgentDaoImpl) GetByHostID(hostID uint) (*model.Agent, error) {
	var agent model.Agent
	err := d.db.Where("host_id = ?", hostID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// Update 更新Agent
func (d *AgentDaoImpl) Update(agent *model.Agent) error {
	agent.UpdateTime.Time = time.Now()
	return d.db.Save(agent).Error
}

// UpdateStatus 更新Agent状态
func (d *AgentDaoImpl) UpdateStatus(id uint, status int, errorMsg string) error {
	updates := map[string]interface{}{
		"status":      status,
		"error_msg":   errorMsg,
		"update_time": time.Now(),
	}
	return d.db.Model(&model.Agent{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateHeartbeat 更新Agent心跳信息 (Linux主机专用)
func (d *AgentDaoImpl) UpdateHeartbeat(hostID uint, heartbeat *model.AgentHeartbeatDto) error {
	updates := map[string]interface{}{
		"pid":            heartbeat.PID,
		"last_heartbeat": time.Now(),
		"update_time":    time.Now(),
		"status":         model.AgentStatusRunning,
	}
	return d.db.Model(&model.Agent{}).Where("host_id = ?", hostID).Updates(updates).Error
}

// UpdateHeartbeatByIP 通过IP更新Agent心跳信息 (Linux主机专用)
func (d *AgentDaoImpl) UpdateHeartbeatByIP(hostID uint, heartbeat *model.AgentHeartbeatDto) error {
	updates := map[string]interface{}{
		"pid":            heartbeat.PID,
		"port":           heartbeat.Port, // 更新端口信息
		"last_heartbeat": time.Now(),
		"update_time":    time.Now(),
		"status":         model.AgentStatusRunning,
	}
	return d.db.Model(&model.Agent{}).Where("host_id = ?", hostID).Updates(updates).Error
}

// UpdateInstallProgress 更新Agent安装进度
func (d *AgentDaoImpl) UpdateInstallProgress(hostID uint, progress int) error {
	updates := map[string]interface{}{
		"install_progress": progress,
		"update_time":      time.Now(),
	}
	return d.db.Model(&model.Agent{}).Where("host_id = ?", hostID).Updates(updates).Error
}

// Delete 删除Agent
func (d *AgentDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.Agent{}, id).Error
}

// DeleteByHostID 通过主机ID删除Agent记录
func (d *AgentDaoImpl) DeleteByHostID(hostID uint) error {
	return d.db.Where("host_id = ?", hostID).Delete(&model.Agent{}).Error
}

// GetList 获取Agent列表 (Linux主机专用)
func (d *AgentDaoImpl) GetList(dto *model.AgentListDto) ([]model.Agent, int64, error) {
	var agents []model.Agent
	var total int64

	query := d.db.Model(&model.Agent{})

	// 条件筛选
	if dto.HostID > 0 {
		query = query.Where("host_id = ?", dto.HostID)
	}
	if dto.Status > 0 {
		query = query.Where("status = ?", dto.Status)
	}
	// 所有Agent都是Linux主机，无需按platform筛选

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	page := dto.Page
	if page <= 0 {
		page = 1
	}
	pageSize := dto.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&agents).Error

	return agents, total, err
}

// GetAllRunning 获取所有运行中的Agent
func (d *AgentDaoImpl) GetAllRunning() ([]model.Agent, error) {
	var agents []model.Agent
	err := d.db.Where("status = ?", model.AgentStatusRunning).Find(&agents).Error
	return agents, err
}

// GetOfflineAgents 获取离线的Agent（超过指定时间没有心跳）
func (d *AgentDaoImpl) GetOfflineAgents(duration time.Duration) ([]model.Agent, error) {
	var agents []model.Agent
	threshold := time.Now().Add(-duration)
	err := d.db.Where("status = ? AND last_heartbeat < ?",
		model.AgentStatusRunning, threshold).Find(&agents).Error
	return agents, err
}

// CreateOrUpdate 创建或更新Agent（根据HostID）
func (d *AgentDaoImpl) CreateOrUpdate(agent *model.Agent) error {
	var existing model.Agent
	err := d.db.Where("host_id = ?", agent.HostID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		return d.Create(agent)
	} else if err != nil {
		// 其他错误
		return err
	} else {
		// 存在，更新记录
		agent.ID = existing.ID
		agent.CreateTime = existing.CreateTime
		return d.Update(agent)
	}
}

// UpdateByHostID 根据主机ID更新Agent
func (d *AgentDaoImpl) UpdateByHostID(hostID uint, updates map[string]interface{}) error {
	updates["update_time"] = time.Now()
	return d.db.Model(&model.Agent{}).Where("host_id = ?", hostID).Updates(updates).Error
}

// GetAgentStatistics 获取Agent统计信息 (Linux主机专用)
func (d *AgentDaoImpl) GetAgentStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总数
	var total int64
	d.db.Model(&model.Agent{}).Count(&total)
	stats["total"] = total

	// 各状态统计
	var running, startError, failed, deploying int64
	d.db.Model(&model.Agent{}).Where("status = ?", model.AgentStatusRunning).Count(&running)
	d.db.Model(&model.Agent{}).Where("status = ?", model.AgentStatusStartError).Count(&startError)
	d.db.Model(&model.Agent{}).Where("status = ?", model.AgentStatusDeployFailed).Count(&failed)
	d.db.Model(&model.Agent{}).Where("status = ?", model.AgentStatusDeploying).Count(&deploying)

	stats["running"] = running
	stats["startError"] = startError
	stats["failed"] = failed
	stats["deploying"] = deploying

	// Linux主机统计(全部都是Linux)
	stats["linux"] = total
	stats["windows"] = 0
	stats["darwin"] = 0

	return stats, nil
}
