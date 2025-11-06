package service

import (
	"fmt"
	"dodevops-api/api/dashboard/model1"
	"dodevops-api/common"
	"math"
	"net"
	"sync"
	"time"
	"gorm.io/gorm"
)

type IDashboardService interface {
	GetDashboardStats() (*model.DashboardStats, error)
	GetBusinessDistributionStats() (*model.BusinessDistributionStats, error)
	GetAssetStats() (*model.AssetStats, error)
}

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService() IDashboardService {
	return &DashboardService{
		db: common.GetDB(),
	}
}

// GetDashboardStats 获取看板统计数据
func (s *DashboardService) GetDashboardStats() (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	// 1. 获取主机统计
	hostStats, err := s.getHostStats()
	if err != nil {
		return nil, err
	}
	stats.HostStats = *hostStats

	// 2. 获取K8s集群统计
	k8sStats, err := s.getK8sClusterStats()
	if err != nil {
		return nil, err
	}
	stats.K8sClusterStats = *k8sStats

	// 3. 获取发布统计 (Jenkins)
	deployStats, err := s.getDeploymentStats()
	if err != nil {
		return nil, err
	}
	stats.DeploymentStats = *deployStats

	// 4. 获取任务统计
	taskStats, err := s.getTaskStats()
	if err != nil {
		return nil, err
	}
	stats.TaskStats = *taskStats

	// 5. 获取服务统计
	serviceStats, err := s.getServiceStats()
	if err != nil {
		return nil, err
	}
	stats.ServiceStats = *serviceStats

	// 6. 获取数据库统计
	databaseStats, err := s.getDatabaseStats()
	if err != nil {
		return nil, err
	}
	stats.DatabaseStats = *databaseStats

	return stats, nil
}

// getHostStats 获取主机统计
func (s *DashboardService) getHostStats() (*model.HostStats, error) {
	var total int64

	// 统计总主机数 (从CMDB主机表)
	err := s.db.Table("cmdb_host").Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 获取所有主机信息，通过TCP检测在线状态 (与监控API保持一致)
	var hosts []struct {
		ID      uint   `gorm:"column:id"`
		SSHIP   string `gorm:"column:ssh_ip"`
		SSHPort int    `gorm:"column:ssh_port"`
	}

	err = s.db.Table("cmdb_host").
		Select("id, ssh_ip, ssh_port").
		Where("ssh_ip != '' AND ssh_port > 0").
		Find(&hosts).Error
	if err != nil {
		return nil, err
	}

	// 并发检测主机在线状态 (与监控API逻辑一致)
	online := s.checkHostsOnlineStatus(hosts)
	offline := max(0, total-int64(online))

	return &model.HostStats{
		Total:   int(total),
		Online:  online,
		Offline: int(offline),
	}, nil
}

// getK8sClusterStats 获取K8s集群统计
func (s *DashboardService) getK8sClusterStats() (*model.K8sClusterStats, error) {
	var total int64
	var healthy int64

	// 统计总集群数
	err := s.db.Table("k8s_cluster").Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 统计健康集群数 (status=2表示运行中)
	err = s.db.Table("k8s_cluster").
		Where("status = ?", 2).  // ClusterStatusRunning = 2
		Count(&healthy).Error
	if err != nil {
		return nil, err
	}

	offline := max(0, total-healthy)

	return &model.K8sClusterStats{
		Total:   int(total),
		Healthy: int(healthy),
		Offline: int(offline),
	}, nil
}

// getDeploymentStats 获取发布统计 (快速发布任务记录)
func (s *DashboardService) getDeploymentStats() (*model.DeploymentStats, error) {
	var total int64
	var success int64

	// 检查快速发布任务表是否存在
	if !s.db.Migrator().HasTable("quick_deployment_tasks") {
		return &model.DeploymentStats{
			Total:       0,
			Success:     0,
			Failed:      0,
			SuccessRate: 0.0,
		}, nil
	}

	// 统计总发布次数
	err := s.db.Table("quick_deployment_tasks").Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 统计成功发布次数 (status=3表示成功)
	err = s.db.Table("quick_deployment_tasks").
		Where("status = ?", 3).
		Count(&success).Error
	if err != nil {
		return nil, err
	}

	failed := max(0, total-success)

	var successRate float64 = 0
	if total > 0 {
		successRate = math.Round((float64(success)/float64(total))*100*100) / 100 // 保留2位小数
	}

	return &model.DeploymentStats{
		Total:       int(total),
		Success:     int(success),
		Failed:      int(failed),
		SuccessRate: successRate,
	}, nil
}

// getTaskStats 获取任务统计
func (s *DashboardService) getTaskStats() (*model.TaskStats, error) {
	var total int64
	var success int64

	// 检查任务记录表是否存在
	if !s.db.Migrator().HasTable("task_work") {
		return &model.TaskStats{
			Total:       0,
			Success:     0,
			Failed:      0,
			SuccessRate: 0.0,
		}, nil
	}

	// 统计总任务执行次数
	err := s.db.Table("task_work").Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 统计成功任务次数 (假设status=3表示成功)
	err = s.db.Table("task_work").
		Where("status = ?", 3).  // 假设3表示任务成功
		Count(&success).Error
	if err != nil {
		return nil, err
	}

	failed := max(0, total-success)

	var successRate float64 = 0
	if total > 0 {
		successRate = math.Round((float64(success)/float64(total))*100*100) / 100 // 保留2位小数
	}

	return &model.TaskStats{
		Total:       int(total),
		Success:     int(success),
		Failed:      int(failed),
		SuccessRate: successRate,
	}, nil
}

// getServiceStats 获取服务统计
func (s *DashboardService) getServiceStats() (*model.ServiceStats, error) {
	var total int64
	var businessLines int64

	// 统计服务总数 (从应用管理表)
	if s.db.Migrator().HasTable("app_application") {
		err := s.db.Table("app_application").Count(&total).Error
		if err != nil {
			return nil, err
		}
	}

	// 统计业务线数量 (从主机分组表或应用分组)
	if s.db.Migrator().HasTable("cmdb_host_group") {
		err := s.db.Table("cmdb_host_group").Count(&businessLines).Error
		if err != nil {
			return nil, err
		}
	}

	return &model.ServiceStats{
		Total:         int(total),
		BusinessLines: int(businessLines),
	}, nil
}

// checkHostsOnlineStatus 并发检测主机在线状态 (与监控API保持一致)
func (s *DashboardService) checkHostsOnlineStatus(hosts []struct {
	ID      uint   `gorm:"column:id"`
	SSHIP   string `gorm:"column:ssh_ip"`
	SSHPort int    `gorm:"column:ssh_port"`
}) int {
	if len(hosts) == 0 {
		return 0
	}

	// 限制并发数量，防止过多并发连接
	maxConcurrent := min(20, len(hosts))
	semaphore := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	var mu sync.Mutex
	onlineCount := 0

	// 并发检测每个主机的在线状态
	for _, host := range hosts {
		wg.Add(1)
		go func(h struct {
			ID      uint   `gorm:"column:id"`
			SSHIP   string `gorm:"column:ssh_ip"`
			SSHPort int    `gorm:"column:ssh_port"`
		}) {
			defer wg.Done()

			// 获取信号量，限制并发数
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// TCP连接检测 (500ms超时，与监控API一致)
			if h.SSHIP != "" && h.SSHPort > 0 {
				conn, err := net.DialTimeout("tcp",
					fmt.Sprintf("%s:%d", h.SSHIP, h.SSHPort),
					500*time.Millisecond)
				if err == nil {
					conn.Close()
					mu.Lock()
					onlineCount++
					mu.Unlock()
				}
			}
		}(host)
	}

	// 等待所有检测完成
	wg.Wait()
	return onlineCount
}

// getDatabaseStats 获取数据库统计
func (s *DashboardService) getDatabaseStats() (*model.DatabaseStats, error) {
	var total int64

	// 检查数据库表是否存在
	if !s.db.Migrator().HasTable("cmdb_sql") {
		return &model.DatabaseStats{
			Total:  0,
			ByType: make(map[string]int),
		}, nil
	}

	// 统计总数据库数量
	err := s.db.Table("cmdb_sql").Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 统计各类型数据库数量
	var typeStats []struct {
		Type  int `json:"type"`
		Count int `json:"count"`
	}

	err = s.db.Table("cmdb_sql").
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeStats).Error
	if err != nil {
		return nil, err
	}

	// 数据库类型映射 (根据cmdbSQL.go:11的定义)
	typeMap := map[int]string{
		1: "MySQL",
		2: "PostgreSQL",
		3: "Redis",
		4: "MongoDB",
		5: "Elasticsearch",
	}

	// 构建按类型统计的map
	byType := make(map[string]int)
	for _, stat := range typeStats {
		if typeName, exists := typeMap[stat.Type]; exists {
			byType[typeName] = stat.Count
		} else {
			byType["Other"] = stat.Count
		}
	}

	return &model.DatabaseStats{
		Total:  int(total),
		ByType: byType,
	}, nil
}

// GetBusinessDistributionStats 获取业务分布统计
func (s *DashboardService) GetBusinessDistributionStats() (*model.BusinessDistributionStats, error) {
	// 检查应用表是否存在
	if !s.db.Migrator().HasTable("app_application") {
		return &model.BusinessDistributionStats{
			TotalServices: 0,
			BusinessLines: []model.BusinessLineStats{},
		}, nil
	}

	// 先检查应用表总记录数
	var totalApps int64
	if err := s.db.Table("app_application").Count(&totalApps).Error; err != nil {
		return nil, err
	}

	// 如果没有应用数据，直接返回空结果
	if totalApps == 0 {
		return &model.BusinessDistributionStats{
			TotalServices: 0,
			BusinessLines: []model.BusinessLineStats{},
		}, nil
	}


	// 获取业务分布统计数据
	var businessStats []struct {
		BusinessGroupID uint   `json:"business_group_id"`
		GroupName       string `json:"group_name"`
		ServiceCount    int    `json:"service_count"`
	}

	// 联查应用表和业务组表，统计各业务线的服务数量
	// 先尝试不过滤状态，看看是否有数据
	err := s.db.Table("app_application a").
		Select("a.business_group_id, COALESCE(g.name, '未知业务线') as group_name, COUNT(*) as service_count").
		Joins("LEFT JOIN cmdb_group g ON a.business_group_id = g.id").
		// 暂时去掉状态过滤，查看所有应用
		// Where("a.status = 1").
		Group("a.business_group_id, g.name").
		Order("service_count DESC").
		Find(&businessStats).Error

	// 如果JOIN查询失败，尝试简单查询
	if err != nil {
		// 回退到简单查询，只统计应用表
		var simpleStats []struct {
			BusinessGroupID uint `json:"business_group_id"`
			Count           int  `json:"count"`
		}

		if err2 := s.db.Table("app_application").
			Select("business_group_id, COUNT(*) as count").
			Group("business_group_id").
			Find(&simpleStats).Error; err2 != nil {
			return nil, err // 返回原始错误
		}

		// 转换为标准格式
		businessStats = []struct {
			BusinessGroupID uint   `json:"business_group_id"`
			GroupName       string `json:"group_name"`
			ServiceCount    int    `json:"service_count"`
		}{}

		for _, stat := range simpleStats {
			businessStats = append(businessStats, struct {
				BusinessGroupID uint   `json:"business_group_id"`
				GroupName       string `json:"group_name"`
				ServiceCount    int    `json:"service_count"`
			}{
				BusinessGroupID: stat.BusinessGroupID,
				GroupName:       fmt.Sprintf("业务组%d", stat.BusinessGroupID),
				ServiceCount:    stat.Count,
			})
		}
	}

	// 计算总服务数量
	totalServices := 0
	for _, stat := range businessStats {
		totalServices += stat.ServiceCount
	}

	// 构建业务线统计数据
	businessLines := make([]model.BusinessLineStats, 0, len(businessStats))
	for _, stat := range businessStats {
		percentage := 0.0
		if totalServices > 0 {
			percentage = math.Round((float64(stat.ServiceCount)/float64(totalServices))*100*100) / 100 // 保留2位小数
		}

		businessLines = append(businessLines, model.BusinessLineStats{
			ID:           stat.BusinessGroupID,
			Name:         stat.GroupName,
			ServiceCount: stat.ServiceCount,
			Percentage:   percentage,
		})
	}

	return &model.BusinessDistributionStats{
		TotalServices: totalServices,
		BusinessLines: businessLines,
	}, nil
}

// GetAssetStats 获取资产统计
func (s *DashboardService) GetAssetStats() (*model.AssetStats, error) {
	var categories []model.AssetCategoryStats
	totalAssets := 0

	// 1. 统计主机资产 (按云平台分类)
	hostCategory, err := s.getHostAssetStats()
	if err != nil {
		return nil, err
	}
	if hostCategory.Total > 0 {
		categories = append(categories, *hostCategory)
		totalAssets += hostCategory.Total
	}

	// 2. 统计数据库资产 (按数据库类型分类)
	dbCategory, err := s.getDatabaseAssetStats()
	if err != nil {
		return nil, err
	}
	if dbCategory.Total > 0 {
		categories = append(categories, *dbCategory)
		totalAssets += dbCategory.Total
	}

	// 3. 统计K8s集群资产
	k8sCategory, err := s.getK8sAssetStats()
	if err != nil {
		return nil, err
	}
	if k8sCategory.Total > 0 {
		categories = append(categories, *k8sCategory)
		totalAssets += k8sCategory.Total
	}

	return &model.AssetStats{
		TotalAssets: totalAssets,
		Categories:  categories,
	}, nil
}

// getHostAssetStats 获取主机资产统计
func (s *DashboardService) getHostAssetStats() (*model.AssetCategoryStats, error) {
	// 检查主机表是否存在
	if !s.db.Migrator().HasTable("cmdb_host") {
		return &model.AssetCategoryStats{
			Category: "主机",
			Total:    0,
			Items:    []model.AssetItemStats{},
		}, nil
	}

	// 统计各云平台主机数量
	var hostStats []struct {
		Vendor int `json:"vendor"`
		Count  int `json:"count"`
	}

	err := s.db.Table("cmdb_host").
		Select("vendor, COUNT(*) as count").
		Group("vendor").
		Find(&hostStats).Error
	if err != nil {
		return nil, err
	}

	// 云平台映射 (根据cmdbHost.go:20的定义)
	vendorMap := map[int]string{
		1: "自建主机",
		2: "阿里云",
		3: "腾讯云",
	}

	var items []model.AssetItemStats
	total := 0
	for _, stat := range hostStats {
		if vendorName, exists := vendorMap[stat.Vendor]; exists {
			items = append(items, model.AssetItemStats{
				Name:  vendorName,
				Count: stat.Count,
			})
		} else {
			items = append(items, model.AssetItemStats{
				Name:  "其他云平台",
				Count: stat.Count,
			})
		}
		total += stat.Count
	}

	return &model.AssetCategoryStats{
		Category: "主机",
		Total:    total,
		Items:    items,
	}, nil
}

// getDatabaseAssetStats 获取数据库资产统计
func (s *DashboardService) getDatabaseAssetStats() (*model.AssetCategoryStats, error) {
	// 检查数据库表是否存在
	if !s.db.Migrator().HasTable("cmdb_sql") {
		return &model.AssetCategoryStats{
			Category: "数据库",
			Total:    0,
			Items:    []model.AssetItemStats{},
		}, nil
	}

	// 统计各类型数据库数量
	var dbStats []struct {
		Type  int `json:"type"`
		Count int `json:"count"`
	}

	err := s.db.Table("cmdb_sql").
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&dbStats).Error
	if err != nil {
		return nil, err
	}

	// 数据库类型映射 (根据cmdbSQL.go:11的定义)
	typeMap := map[int]string{
		1: "MySQL",
		2: "PostgreSQL",
		3: "Redis",
		4: "MongoDB",
		5: "Elasticsearch",
	}

	var items []model.AssetItemStats
	total := 0
	for _, stat := range dbStats {
		if typeName, exists := typeMap[stat.Type]; exists {
			items = append(items, model.AssetItemStats{
				Name:  typeName,
				Count: stat.Count,
			})
		} else {
			items = append(items, model.AssetItemStats{
				Name:  "其他数据库",
				Count: stat.Count,
			})
		}
		total += stat.Count
	}

	return &model.AssetCategoryStats{
		Category: "数据库",
		Total:    total,
		Items:    items,
	}, nil
}

// getK8sAssetStats 获取K8s集群资产统计
func (s *DashboardService) getK8sAssetStats() (*model.AssetCategoryStats, error) {
	// 检查K8s集群表是否存在
	if !s.db.Migrator().HasTable("k8s_cluster") {
		return &model.AssetCategoryStats{
			Category: "K8s集群",
			Total:    0,
			Items:    []model.AssetItemStats{},
		}, nil
	}

	// 统计K8s集群总数
	var total int64
	err := s.db.Table("k8s_cluster").Count(&total).Error
	if err != nil {
		return nil, err
	}

	var items []model.AssetItemStats
	if total > 0 {
		items = append(items, model.AssetItemStats{
			Name:  "K8s集群",
			Count: int(total),
		})
	}

	return &model.AssetCategoryStats{
		Category: "K8s集群",
		Total:    int(total),
		Items:    items,
	}, nil
}