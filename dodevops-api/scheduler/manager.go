package scheduler

import (
	"log"
	"sync"

	"dodevops-api/api/configcenter/model"
)

// Manager 调度器管理器
type Manager struct {
	syncScheduler *SyncScheduler
	mutex         sync.RWMutex
	running       bool
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取调度器管理器实例（单例模式）
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			syncScheduler: NewSyncScheduler(),
			running:       false,
		}
	})
	return instance
}

// Start 启动所有调度器
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.running {
		log.Println("调度器管理器已经在运行中")
		return nil
	}

	log.Println("启动调度器管理器...")

	// 启动定时同步调度器
	if err := m.syncScheduler.Start(); err != nil {
		return err
	}

	m.running = true
	log.Println("调度器管理器启动成功")
	return nil
}

// Stop 停止所有调度器
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running {
		log.Println("调度器管理器未在运行")
		return
	}

	log.Println("停止调度器管理器...")

	// 停止定时同步调度器
	m.syncScheduler.Stop()

	m.running = false
	log.Println("调度器管理器已停止")
}

// IsRunning 检查调度器管理器是否在运行
func (m *Manager) IsRunning() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.running
}

// AddSyncSchedule 添加定时同步配置
func (m *Manager) AddSyncSchedule(schedule *model.SyncSchedule) error {
	if !m.running {
		log.Println("调度器管理器未运行，无法添加定时同步配置")
		return nil
	}

	return m.syncScheduler.AddSchedule(schedule)
}

// RemoveSyncSchedule 移除定时同步配置
func (m *Manager) RemoveSyncSchedule(scheduleID uint) {
	if !m.running {
		log.Println("调度器管理器未运行")
		return
	}

	m.syncScheduler.RemoveSchedule(scheduleID)
}

// UpdateSyncSchedule 更新定时同步配置
func (m *Manager) UpdateSyncSchedule(schedule *model.SyncSchedule) error {
	if !m.running {
		log.Println("调度器管理器未运行，无法更新定时同步配置")
		return nil
	}

	return m.syncScheduler.UpdateSchedule(schedule)
}

// ReloadSyncSchedules 重新加载所有定时同步配置
func (m *Manager) ReloadSyncSchedules() error {
	if !m.running {
		log.Println("调度器管理器未运行，无法重新加载配置")
		return nil
	}

	return m.syncScheduler.LoadActiveSchedules()
}

// GetSyncSchedulerStats 获取定时同步调度器状态
func (m *Manager) GetSyncSchedulerStats() map[string]interface{} {
	if !m.running {
		return map[string]interface{}{
			"status": "stopped",
		}
	}

	stats := m.syncScheduler.GetJobStats()
	stats["status"] = "running"
	return stats
}