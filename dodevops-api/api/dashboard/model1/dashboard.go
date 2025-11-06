package model

// DashboardStats 看板统计数据
type DashboardStats struct {
	HostStats        HostStats        `json:"hostStats"`        // 主机统计
	K8sClusterStats  K8sClusterStats  `json:"k8sClusterStats"`  // K8s集群统计
	DeploymentStats  DeploymentStats  `json:"deploymentStats"`  // 发布统计
	TaskStats        TaskStats        `json:"taskStats"`        // 任务统计
	ServiceStats     ServiceStats     `json:"serviceStats"`     // 服务统计
	DatabaseStats    DatabaseStats    `json:"databaseStats"`    // 数据库统计
}

// HostStats 主机统计
type HostStats struct {
	Total   int `json:"total"`   // 主机总数
	Online  int `json:"online"`  // 在线数量
	Offline int `json:"offline"` // 离线数量
}

// K8sClusterStats K8s集群统计
type K8sClusterStats struct {
	Total   int `json:"total"`   // 集群总数
	Healthy int `json:"healthy"` // 健康数量
	Offline int `json:"offline"` // 离线数量
}

// DeploymentStats 发布统计
type DeploymentStats struct {
	Total       int     `json:"total"`       // 发布总次数
	Success     int     `json:"success"`     // 成功次数
	Failed      int     `json:"failed"`      // 失败次数
	SuccessRate float64 `json:"successRate"` // 成功率
}

// TaskStats 任务统计
type TaskStats struct {
	Total       int     `json:"total"`       // 任务执行总次数
	Success     int     `json:"success"`     // 成功次数
	Failed      int     `json:"failed"`      // 失败次数
	SuccessRate float64 `json:"successRate"` // 成功率
}

// ServiceStats 服务统计
type ServiceStats struct {
	Total         int `json:"total"`         // 服务总数
	BusinessLines int `json:"businessLines"` // 业务线数量
}

// DatabaseStats 数据库统计
type DatabaseStats struct {
	Total   int            `json:"total"`   // 数据库总数
	ByType  map[string]int `json:"byType"`  // 按类型统计
}

// BusinessDistributionStats 业务分布统计
type BusinessDistributionStats struct {
	TotalServices    int                    `json:"totalServices"`    // 总服务数量
	BusinessLines    []BusinessLineStats    `json:"businessLines"`    // 业务线列表
}

// BusinessLineStats 业务线统计
type BusinessLineStats struct {
	ID           uint   `json:"id"`           // 业务组ID
	Name         string `json:"name"`         // 业务线名称
	ServiceCount int    `json:"serviceCount"` // 服务数量
	Percentage   float64 `json:"percentage"`  // 占比
}

// AssetStats 资产统计
type AssetStats struct {
	TotalAssets int                `json:"totalAssets"` // 总资产数量
	Categories  []AssetCategoryStats `json:"categories"`  // 资产分类统计
}

// AssetCategoryStats 资产分类统计
type AssetCategoryStats struct {
	Category string           `json:"category"` // 资产分类名称 (主机/数据库/K8s集群)
	Total    int              `json:"total"`    // 该分类总数
	Items    []AssetItemStats `json:"items"`    // 具体资产项统计
}

// AssetItemStats 资产项统计
type AssetItemStats struct {
	Name  string `json:"name"`  // 资产项名称 (自建主机/阿里云/MySQL等)
	Count int    `json:"count"` // 数量
}