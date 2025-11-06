package model

import (
	"time"
)

// KubeCluster K8s集群表
type KubeCluster struct {
	ID          uint      `gorm:"primaryKey;comment:'主键ID'" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex;comment:'集群名称'" json:"name"`
	Version     string    `gorm:"size:50;not null;comment:'集群版本'" json:"version"`
	Status      int       `gorm:"not null;default:1;comment:'集群状态:1-创建中,2-运行中,3-离线'" json:"status"`
	Credential  string    `gorm:"type:text;comment:'集群凭证(kubeconfig)'" json:"credential"`
	Description string    `gorm:"type:text;comment:'集群描述'" json:"description"`
	ClusterType int       `gorm:"not null;default:1;comment:'集群类型:1-自建,2-导入'" json:"clusterType"`
	NodeCount   int       `gorm:"default:0;comment:'节点数量'" json:"nodeCount"`
	ReadyNodes  int       `gorm:"default:0;comment:'就绪节点数'" json:"readyNodes"`
	MasterNodes int       `gorm:"default:0;comment:'Master节点数'" json:"masterNodes"`
	WorkerNodes int       `gorm:"default:0;comment:'Worker节点数'" json:"workerNodes"`
	LastSyncAt  *time.Time `gorm:"comment:'最后同步时间'" json:"lastSyncAt"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// KubeClusterStatus 集群状态枚举
const (
	ClusterStatusCreating = 1 // 创建中
	ClusterStatusRunning  = 2 // 运行中
	ClusterStatusStopped  = 3 // 离线
)

// KubeClusterType 集群类型枚举
const (
	ClusterTypeSelfBuilt = 1 // 自建集群
	ClusterTypeImported  = 2 // 导入集群
)

// TableName 指定表名
func (KubeCluster) TableName() string {
	return "k8s_cluster"
}

// GetStatusText 获取集群状态文本
func (c *KubeCluster) GetStatusText() string {
	switch c.Status {
	case ClusterStatusCreating:
		return "创建中"
	case ClusterStatusRunning:
		return "运行中"
	case ClusterStatusStopped:
		return "离线"
	default:
		return "未知"
	}
}

// NodeInfo 从K8s API获取的节点信息
type NodeInfo struct {
	Name         string            `json:"name"`
	Role         string            `json:"role"`         // master/worker/etcd
	Status       string            `json:"status"`       // Ready/NotReady
	InternalIP   string            `json:"internalIP"`
	ExternalIP   string            `json:"externalIP"`
	Version      string            `json:"version"`
	OS           string            `json:"os"`
	Capacity     map[string]string `json:"capacity"`     // CPU、内存等资源容量
	Allocatable  map[string]string `json:"allocatable"`  // 可分配资源
	Conditions   []NodeCondition   `json:"conditions"`   // 节点状态条件
}

// NodeCondition 节点状态条件
type NodeCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// ClusterDetailResponse 集群详情响应
type ClusterDetailResponse struct {
	Cluster     KubeCluster      `json:"cluster"`
	Nodes       []NodeInfo       `json:"nodes"`
	Summary     ClusterSummary   `json:"summary"`
	Components  []ComponentInfo  `json:"components"`  // 安装的组件
	Workloads   WorkloadSummary  `json:"workloads"`   // 工作负载统计
	Network     NetworkInfo      `json:"network"`     // 网络配置
	Events      []ClusterEvent   `json:"events"`      // 集群事件
	Monitoring  MonitoringInfo   `json:"monitoring"`  // 监控信息
	Runtime     RuntimeSummary   `json:"runtime"`     // 运行时信息
}

// ClusterSummary 集群概要信息
type ClusterSummary struct {
	TotalNodes    int `json:"totalNodes"`
	ReadyNodes    int `json:"readyNodes"`
	MasterNodes   int `json:"masterNodes"`
	WorkerNodes   int `json:"workerNodes"`
}

// ComponentInfo 组件信息
type ComponentInfo struct {
	Name      string `json:"name"`      // 组件名称
	Namespace string `json:"namespace"` // 命名空间
	Status    string `json:"status"`    // 状态
	Version   string `json:"version"`   // 版本
	Type      string `json:"type"`      // 组件类型 (system/addon)
}

// WorkloadSummary 工作负载统计
type WorkloadSummary struct {
	TotalDeployments  int `json:"totalDeployments"`  // Deployment总数
	TotalStatefulSets int `json:"totalStatefulSets"` // StatefulSet总数
	TotalDaemonSets   int `json:"totalDaemonSets"`   // DaemonSet总数
	TotalJobs         int `json:"totalJobs"`         // Job总数
	TotalCronJobs     int `json:"totalCronJobs"`     // CronJob总数
	TotalPods         int `json:"totalPods"`         // Pod总数
	RunningPods       int `json:"runningPods"`       // 运行中的Pod数
}

// NetworkInfo 网络配置信息
type NetworkInfo struct {
	ServiceCIDR     string `json:"serviceCIDR"`     // Service CIDR
	PodCIDR         string `json:"podCIDR"`         // Pod CIDR
	DNSService      string `json:"dnsService"`      // DNS服务
	NetworkPlugin   string `json:"networkPlugin"`   // 网络插件
	ProxyMode       string `json:"proxyMode"`       // 服务转发模式
	APIServerEndpoint string `json:"apiServerEndpoint"` // API Server内网端点
}

// ClusterEvent 集群事件
type ClusterEvent struct {
	Type          string `json:"type"`          // 事件类型
	Reason        string `json:"reason"`        // 原因
	Message       string `json:"message"`       // 消息
	Source        string `json:"source"`        // 事件源
	FirstTime     string `json:"firstTime"`     // 首次时间
	LastTime      string `json:"lastTime"`      // 最后时间
	Count         int32  `json:"count"`         // 发生次数
	InvolvedObject string `json:"involvedObject"` // 相关对象
}

// MonitoringInfo 监控信息
type MonitoringInfo struct {
	CPU     ClusterResourceMetrics `json:"cpu"`     // CPU监控
	Memory  ClusterResourceMetrics `json:"memory"`  // 内存监控
	Network NetworkMetrics         `json:"network"` // 网络监控
	Storage StorageMetrics         `json:"storage"` // 存储监控
}

// ClusterResourceMetrics 集群资源监控指标
type ClusterResourceMetrics struct {
	Total       string  `json:"total"`       // 总量
	Used        string  `json:"used"`        // 已使用
	Available   string  `json:"available"`   // 可用量
	UsageRate   float64 `json:"usageRate"`   // 使用率 (0-100)
	RequestRate float64 `json:"requestRate"` // 请求率 (0-100)
}

// NetworkMetrics 网络监控指标
type NetworkMetrics struct {
	InboundTraffic  string `json:"inboundTraffic"`  // 入站流量
	OutboundTraffic string `json:"outboundTraffic"` // 出站流量
	PacketsIn       int64  `json:"packetsIn"`       // 入站包数
	PacketsOut      int64  `json:"packetsOut"`      // 出站包数
}

// StorageMetrics 存储监控指标
type StorageMetrics struct {
	TotalPVs      int    `json:"totalPVs"`      // PV总数
	BoundPVs      int    `json:"boundPVs"`      // 已绑定PV数
	TotalPVCs     int    `json:"totalPVCs"`     // PVC总数
	StorageClasses []string `json:"storageClasses"` // 存储类列表
}

// RuntimeSummary 运行时概要信息
type RuntimeSummary struct {
	KubernetesVersion    string `json:"kubernetesVersion"`    // Kubernetes版本
	ContainerRuntime     string `json:"containerRuntime"`     // 容器运行时
	APIServerVersion     string `json:"apiServerVersion"`     // API Server版本
	EtcdVersion          string `json:"etcdVersion"`          // etcd版本
	CoreDNSVersion       string `json:"coreDNSVersion"`       // CoreDNS版本
	KubeProxyVersion     string `json:"kubeProxyVersion"`     // kube-proxy版本
	UpTime               string `json:"upTime"`               // 集群运行时间
}

// CreateKubeClusterRequest 创建K8s集群请求
type CreateKubeClusterRequest struct {
	Name        string `json:"name" binding:"required"`        // 集群名称
	Description string `json:"description"`                    // 集群描述
	ClusterType int    `json:"clusterType"`                    // 集群类型:1-自建,2-导入(默认为自建)
	
	// 自建集群参数
	Version           string     `json:"version"`           // K8s版本
	NodeConfig        *NodeConfig `json:"nodeConfig"`      // 节点配置
	AutoDeploy        bool       `json:"autoDeploy"`       // 是否自动部署
	DeploymentMode    int        `json:"deploymentMode"`   // 部署模式:1-单Master,2-多Master
	EnabledComponents []string       `json:"enabledComponents"` // 启用组件
	PrivateRegistry   string         `json:"privateRegistry"`   // 私有镜像仓库地址（兼容旧版本）
	RegistryUsername  string         `json:"registryUsername"`  // 镜像仓库用户名（兼容旧版本）
	RegistryPassword  string         `json:"registryPassword"`  // 镜像仓库密码（兼容旧版本）
	RegistryConfig    *RegistryConfig `json:"registryConfig"`    // 镜像仓库配置（新版本）
	TaskName          string         `json:"taskName"`          // 任务名称
	TaskDescription   string         `json:"taskDescription"`   // 任务描述
	
	// 导入集群参数
	Kubeconfig string `json:"kubeconfig"` // K8s凭证(kubeconfig内容)
}

// NodeConfig 节点配置
type NodeConfig struct {
	MasterHostIDs []uint `json:"masterHostIds" binding:"required"` // Master节点主机ID
	WorkerHostIDs []uint `json:"workerHostIds"`                    // Worker节点主机ID  
	EtcdHostIDs   []uint `json:"etcdHostIds" binding:"required"`   // ETCD节点主机ID
}

// RegistryConfig 镜像仓库配置
type RegistryConfig struct {
	PrivateRegistry    string `json:"privateRegistry"`    // 私有镜像仓库地址
	RegistryUsername   string `json:"registryUsername"`   // 镜像仓库用户名
	RegistryPassword   string `json:"registryPassword"`   // 镜像仓库密码
	UsePrivateRegistry bool   `json:"usePrivateRegistry"` // 是否使用私有仓库
}

// KubeClusterListResponse 集群列表响应
type KubeClusterListResponse struct {
	List  []KubeCluster `json:"list"`
	Total int64         `json:"total"`
}

// UpdateKubeClusterRequest 更新K8s集群请求
type UpdateKubeClusterRequest struct {
	Name        string `json:"name"`        // 集群名称
	Description string `json:"description"` // 集群描述
	Credential  string `json:"credential"`  // K8s凭证(kubeconfig内容)
	Version     string `json:"version"`     // 集群版本(可选，同步时会自动更新)
}

// K8sNode K8s节点信息（扩展版本）
type K8sNode struct {
	Name          string                `json:"name"`          // 节点名称
	InternalIP    string               `json:"internalIP"`    // 内部IP地址
	ExternalIP    string               `json:"externalIP"`    // 外部IP地址
	Status        string               `json:"status"`        // 节点状态 Ready/NotReady
	Roles         string               `json:"roles"`         // 节点角色 control-plane,master 或 worker
	Conditions    []NodeCondition      `json:"conditions"`    // 节点状态详细条件
	Configuration NodeConfiguration    `json:"configuration"` // 节点配置信息
	PodMetrics    PodMetrics          `json:"podMetrics"`    // 容器组统计
	Resources     NodeResources       `json:"resources"`     // CPU和内存资源
	Runtime       RuntimeInfo         `json:"runtime"`       // 运行时和版本信息
	Scheduling    NodeSchedulingInfo  `json:"scheduling"`    // 调度相关信息（污点、是否可调度等）
	CreatedAt     string              `json:"createdAt"`     // 创建时间
}

// NodeConfiguration 节点配置信息
type NodeConfiguration struct {
	Role           string `json:"role"`           // 节点角色 master/worker
	Architecture   string `json:"architecture"`   // 系统架构
	KernelVersion  string `json:"kernelVersion"`  // 内核版本
	OSImage        string `json:"osImage"`        // 操作系统镜像
	Labels         map[string]string `json:"labels"`    // 节点标签
	Annotations    map[string]string `json:"annotations"` // 节点注释
}

// PodMetrics 容器组统计信息
type PodMetrics struct {
	Allocated int `json:"allocated"` // 已分配的Pod数量
	Total     int `json:"total"`     // 总的Pod容量
}

// NodeResources 节点资源信息
type NodeResources struct {
	CPU    ResourceInfo `json:"cpu"`    // CPU资源
	Memory ResourceInfo `json:"memory"` // 内存资源
}

// ResourceInfo 资源详细信息
type ResourceInfo struct {
	Requests    string `json:"requests"`    // 请求量
	Usage       string `json:"usage"`       // 使用量
	Capacity    string `json:"capacity"`    // 总容量
	Allocatable string `json:"allocatable"` // 可分配量
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	KubeletVersion       string `json:"kubeletVersion"`       // Kubelet版本
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"` // 容器运行时版本
	KubeProxyVersion     string `json:"kubeProxyVersion"`     // KubeProxy版本
	OperatingSystem      string `json:"operatingSystem"`      // 操作系统
	OSImage              string `json:"osImage"`              // 操作系统镜像
}

// K8sNodeDetail 节点详细信息
type K8sNodeDetail struct {
	K8sNode
	Pods     []PodInfo `json:"pods"`     // 节点上运行的Pod列表
	Events   []EventInfo `json:"events"`   // 相关事件
	Metrics  NodeMetrics `json:"metrics"`  // 详细监控指标
}

// NodeDetailResponse 节点详细信息响应
type NodeDetailResponse struct {
	// 基本信息
	Name         string            `json:"name"`         // 节点名称
	CreatedAt    string            `json:"createdAt"`    // 创建时间
	UID          string            `json:"uid"`          // UID
	ProviderID   string            `json:"providerID"`   // 提供者ID

	// IP地址信息
	InternalIP   string            `json:"internalIP"`   // 内部IP
	ExternalIP   string            `json:"externalIP"`   // 外部IP
	Hostname     string            `json:"hostname"`     // 主机名

	// 系统信息
	OSImage      string            `json:"osImage"`      // 系统镜像
	KernelVersion string           `json:"kernelVersion"` // 内核版本
	Architecture string            `json:"architecture"` // 架构
	OperatingSystem string         `json:"operatingSystem"` // 操作系统
	MachineID    string            `json:"machineID"`    // 机器ID
	SystemUUID   string            `json:"systemUUID"`   // 系统UUID
	BootID       string            `json:"bootID"`       // 启动ID

	// K8s组件版本
	KubeletVersion       string `json:"kubeletVersion"`       // Kubelet版本
	KubeProxyVersion     string `json:"kubeProxyVersion"`     // Kube-Proxy版本
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"` // 容器运行时版本

	// 资源信息
	Capacity     map[string]string `json:"capacity"`     // 总容量
	Allocatable  map[string]string `json:"allocatable"`  // 可分配资源

	// 调度信息
	Unschedulable bool              `json:"unschedulable"` // 是否不可调度
	Taints       []NodeTaint        `json:"taints"`        // 污点列表
	Labels       map[string]string  `json:"labels"`        // 标签
	Annotations  map[string]string  `json:"annotations"`   // 注释

	// 状态信息
	Status       string             `json:"status"`       // 节点状态
	Conditions   []NodeCondition    `json:"conditions"`   // 状态条件

	// 监控信息
	Monitoring   NodeMonitoringInfo `json:"monitoring"`   // 监控信息

	// Pod信息
	PodInfo      NodePodInfo        `json:"podInfo"`      // Pod统计信息
	PodList      []NodePodDetail    `json:"podList"`      // 节点上的Pod列表

	// CIDR信息
	PodCIDR      string             `json:"podCIDR"`      // 容器组CIDR
	PodCIDRs     []string           `json:"podCIDRs"`     // 容器组CIDR列表
}

// NodeMonitoringInfo 节点监控信息
type NodeMonitoringInfo struct {
	CPU     NodeResourceUsage `json:"cpu"`     // CPU使用情况
	Memory  NodeResourceUsage `json:"memory"`  // 内存使用情况
	Storage NodeResourceUsage `json:"storage"` // 存储使用情况
	Network NodeNetworkUsage  `json:"network"` // 网络使用情况
}

// NodeResourceUsage 节点资源使用情况
type NodeResourceUsage struct {
	Total       string  `json:"total"`       // 总量
	Used        string  `json:"used"`        // 已使用
	Available   string  `json:"available"`   // 可用量
	Requests    string  `json:"requests"`    // 请求量
	Limits      string  `json:"limits"`      // 限制量
	UsageRate   float64 `json:"usageRate"`   // 使用率 (0-100)
	RequestRate float64 `json:"requestRate"` // 请求率 (0-100)
}

// NodeNetworkUsage 节点网络使用情况
type NodeNetworkUsage struct {
	InboundBytes   int64  `json:"inboundBytes"`   // 入站字节数
	OutboundBytes  int64  `json:"outboundBytes"`  // 出站字节数
	InboundPackets int64  `json:"inboundPackets"` // 入站包数
	OutboundPackets int64 `json:"outboundPackets"` // 出站包数
}

// NodePodInfo Pod统计信息
type NodePodInfo struct {
	TotalPods      int `json:"totalPods"`      // Pod总数
	RunningPods    int `json:"runningPods"`    // 运行中的Pod数
	PendingPods    int `json:"pendingPods"`    // 等待中的Pod数
	FailedPods     int `json:"failedPods"`     // 失败的Pod数
	SucceededPods  int `json:"succeededPods"`  // 成功的Pod数
}

// NodePodDetail 节点上Pod的详细信息
type NodePodDetail struct {
	Name          string            `json:"name"`          // Pod名称
	Namespace     string            `json:"namespace"`     // 命名空间
	Status        string            `json:"status"`        // 状态
	Phase         string            `json:"phase"`         // 阶段
	RestartCount  int32             `json:"restartCount"`  // 重启次数
	CreatedAt     string            `json:"createdAt"`     // 创建时间
	Labels        map[string]string `json:"labels"`        // 标签
	CPURequests   string            `json:"cpuRequests"`   // CPU请求
	CPULimits     string            `json:"cpuLimits"`     // CPU限制
	MemoryRequests string           `json:"memoryRequests"` // 内存请求
	MemoryLimits  string            `json:"memoryLimits"`  // 内存限制
	Containers    []ContainerStatus `json:"containers"`    // 容器状态
}

// ContainerStatus 容器状态
type ContainerStatus struct {
	Name         string `json:"name"`         // 容器名称
	Image        string `json:"image"`        // 镜像
	State        string `json:"state"`        // 状态
	Ready        bool   `json:"ready"`        // 就绪状态
	RestartCount int32  `json:"restartCount"` // 重启次数
}

// PodInfo Pod信息
type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	CPUUsage  string `json:"cpuUsage"`
	MemUsage  string `json:"memUsage"`
}

// EventInfo 事件信息
type EventInfo struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	FirstTime string `json:"firstTime"`
	LastTime  string `json:"lastTime"`
	Count     int    `json:"count"`
}

// NodeMetrics 节点监控指标
type NodeMetrics struct {
	CPUUsagePercentage    float64 `json:"cpuUsagePercentage"`
	MemoryUsagePercentage float64 `json:"memoryUsagePercentage"`
	DiskUsagePercentage   float64 `json:"diskUsagePercentage"`
	NetworkInBytes        int64   `json:"networkInBytes"`
	NetworkOutBytes       int64   `json:"networkOutBytes"`
}

// NodeTaint 节点污点
type NodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"` // NoSchedule, PreferNoSchedule, NoExecute
}

// NodeSchedulingInfo 节点调度信息
type NodeSchedulingInfo struct {
	Unschedulable bool        `json:"unschedulable"` // 是否不可调度
	Taints        []NodeTaint `json:"taints"`        // 节点污点
}

// AddTaintRequest 添加污点请求
type AddTaintRequest struct {
	Key    string `json:"key" binding:"required"`
	Value  string `json:"value"`
	Effect string `json:"effect" binding:"required,oneof=NoSchedule PreferNoSchedule NoExecute"`
}

// RemoveTaintRequest 移除污点请求
type RemoveTaintRequest struct {
	Key    string `json:"key" binding:"required"`
	Effect string `json:"effect"`
}

// AddLabelRequest 添加标签请求
type AddLabelRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// RemoveLabelRequest 移除标签请求
type RemoveLabelRequest struct {
	Key string `json:"key" binding:"required"`
}

// CordonNodeRequest 封锁节点请求
type CordonNodeRequest struct {
	Unschedulable bool   `json:"unschedulable"`
	Reason        string `json:"reason"`
}

// DrainNodeRequest 驱逐节点请求
type DrainNodeRequest struct {
	Force              bool `json:"force"`               // 强制驱逐
	DeleteLocalData    bool `json:"deleteLocalData"`     // 删除本地数据
	IgnoreDaemonSets   bool `json:"ignoreDaemonSets"`    // 忽略DaemonSet
	GracePeriodSeconds int  `json:"gracePeriodSeconds"`  // 优雅终止时间
}

// NodeResourceAllocation 节点资源分配详情
type NodeResourceAllocation struct {
	NodeName string                    `json:"nodeName"`
	Capacity map[string]string         `json:"capacity"`     // 节点总容量
	Allocatable map[string]string      `json:"allocatable"`  // 可分配资源
	Allocated map[string]string        `json:"allocated"`    // 已分配资源
	PodList []PodResourceInfo         `json:"podList"`      // Pod资源使用详情
}

// PodResourceInfo Pod资源信息
type PodResourceInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Requests  map[string]string `json:"requests"`
	Limits    map[string]string `json:"limits"`
}

// K8sNamespace K8s命名空间信息
type K8sNamespace struct {
	Name          string             `json:"name"`          // 命名空间名称
	Status        string             `json:"status"`        // 状态 Active/Terminating
	Labels        map[string]string  `json:"labels"`        // 标签
	Annotations   map[string]string  `json:"annotations"`   // 注释
	ResourceQuotas []ResourceQuotaDetail `json:"resourceQuotas"` // 资源配额列表
	LimitRanges   []LimitRangeDetail    `json:"limitRanges"`    // 默认资源限制列表
	ResourceCount NamespaceResourceCount `json:"resourceCount"` // 资源统计
	CreatedAt     string             `json:"createdAt"`     // 创建时间
}

// NamespaceResourceCount 命名空间资源统计
type NamespaceResourceCount struct {
	PodCount       int `json:"podCount"`       // Pod数量
	ServiceCount   int `json:"serviceCount"`   // Service数量
	SecretCount    int `json:"secretCount"`    // Secret数量
	ConfigMapCount int `json:"configMapCount"` // ConfigMap数量
}

// ResourceQuotaDetail 详细资源配额信息
type ResourceQuotaDetail struct {
	Name        string            `json:"name"`        // ResourceQuota名称
	Hard        map[string]string `json:"hard"`        // 硬限制
	Used        map[string]string `json:"used"`        // 已使用
	CPUQuota    QuotaInfo         `json:"cpuQuota"`    // CPU配额详情
	MemoryQuota QuotaInfo         `json:"memoryQuota"` // 内存配额详情
	StorageQuota QuotaInfo        `json:"storageQuota"` // 存储配额详情
	CreatedAt   string            `json:"createdAt"`   // 创建时间
}

// QuotaInfo 配额详细信息
type QuotaInfo struct {
	Hard string `json:"hard"` // 限制值
	Used string `json:"used"` // 已使用值
}

// LimitRangeDetail 默认资源限制详情
type LimitRangeDetail struct {
	Name      string            `json:"name"`      // LimitRange名称
	Limits    []LimitRangeItem  `json:"limits"`    // 限制项列表
	CreatedAt string            `json:"createdAt"` // 创建时间
}

// LimitRangeItem 限制项
type LimitRangeItem struct {
	Type                 string            `json:"type"`                 // 限制类型 Container/Pod/PersistentVolumeClaim
	Max                  map[string]string `json:"max,omitempty"`        // 最大限制
	Min                  map[string]string `json:"min,omitempty"`        // 最小限制
	Default              map[string]string `json:"default,omitempty"`    // 默认值
	DefaultRequest       map[string]string `json:"defaultRequest,omitempty"` // 默认请求值
	MaxLimitRequestRatio map[string]string `json:"maxLimitRequestRatio,omitempty"` // 最大限制与请求比率
}

// CreateNamespaceRequest 创建命名空间请求
type CreateNamespaceRequest struct {
	Name        string            `json:"name" binding:"required"`        // 命名空间名称
	Labels      map[string]string `json:"labels"`                         // 标签
	Annotations map[string]string `json:"annotations"`                    // 注释
}

// UpdateNamespaceRequest 更新命名空间请求
type UpdateNamespaceRequest struct {
	Labels      map[string]string `json:"labels"`      // 标签
	Annotations map[string]string `json:"annotations"` // 注释
}

// NamespaceListResponse 命名空间列表响应
type NamespaceListResponse struct {
	Namespaces []K8sNamespace `json:"namespaces"`
	Total      int            `json:"total"`
}

// CreateResourceQuotaRequest 创建资源配额请求
type CreateResourceQuotaRequest struct {
	Name        string            `json:"name" binding:"required"`        // ResourceQuota名称
	Hard        map[string]string `json:"hard" binding:"required"`        // 硬限制
	Scopes      []string          `json:"scopes,omitempty"`               // 作用域
	ScopeSelector map[string]interface{} `json:"scopeSelector,omitempty"`   // 作用域选择器
}

// UpdateResourceQuotaRequest 更新资源配额请求
type UpdateResourceQuotaRequest struct {
	Hard        map[string]string `json:"hard" binding:"required"`        // 硬限制
	Scopes      []string          `json:"scopes,omitempty"`               // 作用域
	ScopeSelector map[string]interface{} `json:"scopeSelector,omitempty"`   // 作用域选择器
}

// CreateLimitRangeRequest 创建限制范围请求
type CreateLimitRangeRequest struct {
	Name string                   `json:"name" binding:"required"` // LimitRange名称
	Spec LimitRangeRequestSpec    `json:"spec" binding:"required"` // LimitRange规格
}

// LimitRangeRequestSpec LimitRange规格
type LimitRangeRequestSpec struct {
	Limits []LimitRangeItem `json:"limits" binding:"required"` // 限制项列表
}

// UpdateLimitRangeRequest 更新限制范围请求
type UpdateLimitRangeRequest struct {
	Spec LimitRangeRequestSpec `json:"spec" binding:"required"` // LimitRange规格
}

// ResourceQuotaListResponse 资源配额列表响应
type ResourceQuotaListResponse struct {
	ResourceQuotas []ResourceQuotaDetail `json:"resourceQuotas"`
	Total          int                   `json:"total"`
}

// LimitRangeListResponse 限制范围列表响应
type LimitRangeListResponse struct {
	LimitRanges []LimitRangeDetail `json:"limitRanges"`
	Total       int                `json:"total"`
}

// ================== 工作负载相关结构体 ==================

// WorkloadType 工作负载类型
type WorkloadType string

const (
	WorkloadTypeDeployment  WorkloadType = "Deployment"
	WorkloadTypeStatefulSet WorkloadType = "StatefulSet"
	WorkloadTypeDaemonSet   WorkloadType = "DaemonSet"
	WorkloadTypeJob         WorkloadType = "Job"
	WorkloadTypeCronJob     WorkloadType = "CronJob"
	WorkloadTypePod         WorkloadType = "Pod"
)

// K8sWorkload K8s工作负载基础信息
type K8sWorkload struct {
	Name         string                 `json:"name"`         // 名称
	Type         WorkloadType          `json:"type"`         // 类型
	Namespace    string                 `json:"namespace"`    // 命名空间
	Labels       map[string]string      `json:"labels"`       // 标签
	Replicas     int32                  `json:"replicas"`     // 副本数
	ReadyReplicas int32                 `json:"readyReplicas"` // 就绪副本数
	Resources    WorkloadResources      `json:"resources"`    // 资源配置
	Images       []string              `json:"images"`       // 镜像列表
	Status       string                `json:"status"`       // 状态
	CreatedAt    string                `json:"createdAt"`    // 创建时间
	UpdatedAt    string                `json:"updatedAt"`    // 更新时间
}

// WorkloadResources 工作负载资源配置
type WorkloadResources struct {
	Requests ResourceSpec `json:"requests"` // 资源请求
	Limits   ResourceSpec `json:"limits"`   // 资源限制
}

// ResourceSpec 资源规格
type ResourceSpec struct {
	CPU    string `json:"cpu"`    // CPU
	Memory string `json:"memory"` // 内存
}

// K8sWorkloadDetail 工作负载详情
type K8sWorkloadDetail struct {
	K8sWorkload
	Pods       []K8sPodInfo     `json:"pods"`       // Pod列表
	Conditions []WorkloadCondition `json:"conditions"` // 状态条件
	Events     []K8sEvent       `json:"events"`     // 相关事件
	Spec       interface{}      `json:"spec"`       // 完整规格配置
}

// K8sPodInfo Pod基础信息
type K8sPodInfo struct {
	Name          string            `json:"name"`          // 实例名称
	Status        string            `json:"status"`        // 状态
	Phase         string            `json:"phase"`         // 阶段
	RestartCount  int32             `json:"restartCount"`  // 重启次数
	NodeName      string            `json:"nodeName"`      // 节点名称
	PodIP         string            `json:"podIP"`         // Pod IP
	HostIP        string            `json:"hostIP"`        // 主机IP
	Resources     WorkloadResources `json:"resources"`     // 资源配置
	RunningTime   string            `json:"runningTime"`   // 运行时间
	CreatedAt     string            `json:"createdAt"`     // 创建时间
	Labels        map[string]string `json:"labels"`        // 标签
	Containers    []ContainerInfo   `json:"containers"`    // 容器信息
}

// K8sPodDetail Pod详情
type K8sPodDetail struct {
	K8sPodInfo
	Conditions []PodCondition    `json:"conditions"`    // Pod状态条件
	Events     []K8sEvent        `json:"events"`        // 相关事件
	Volumes    []VolumeInfo      `json:"volumes"`       // 挂载卷信息
	Spec       interface{}       `json:"spec"`          // 完整规格配置
}

// ContainerInfo 容器信息
type ContainerInfo struct {
	Name         string        `json:"name"`         // 容器名称
	Image        string        `json:"image"`        // 镜像
	State        string        `json:"state"`        // 状态
	Ready        bool          `json:"ready"`        // 就绪状态
	RestartCount int32         `json:"restartCount"` // 重启次数
	Resources    WorkloadResources `json:"resources"` // 资源配置
	Ports        []ContainerPort   `json:"ports"`     // 端口配置
	Env          []EnvVar      `json:"env"`          // 环境变量
}

// ContainerPort 容器端口
type ContainerPort struct {
	Name          string `json:"name"`          // 端口名称
	ContainerPort int32  `json:"containerPort"` // 容器端口
	Protocol      string `json:"protocol"`      // 协议
}

// EnvVar 环境变量
type EnvVar struct {
	Name  string `json:"name"`  // 变量名
	Value string `json:"value"` // 变量值
}

// VolumeInfo 存储卷信息
type VolumeInfo struct {
	Name      string `json:"name"`      // 卷名称
	Type      string `json:"type"`      // 卷类型
	MountPath string `json:"mountPath"` // 挂载路径
	ReadOnly  bool   `json:"readOnly"`  // 只读状态
}

// WorkloadCondition 工作负载状态条件
type WorkloadCondition struct {
	Type               string `json:"type"`               // 条件类型
	Status             string `json:"status"`             // 状态
	LastTransitionTime string `json:"lastTransitionTime"` // 最后转换时间
	Reason             string `json:"reason"`             // 原因
	Message            string `json:"message"`            // 消息
}

// PodCondition Pod状态条件
type PodCondition struct {
	Type               string `json:"type"`               // 条件类型
	Status             string `json:"status"`             // 状态
	LastTransitionTime string `json:"lastTransitionTime"` // 最后转换时间
	Reason             string `json:"reason"`             // 原因
	Message            string `json:"message"`            // 消息
}

// K8sEvent K8s事件
type K8sEvent struct {
	Type      string `json:"type"`      // 事件类型
	Reason    string `json:"reason"`    // 原因
	Message   string `json:"message"`   // 消息
	Source    string `json:"source"`    // 事件源
	Count     int32  `json:"count"`     // 发生次数
	FirstTime string `json:"firstTime"` // 首次时间
	LastTime  string `json:"lastTime"`  // 最后时间
}

// ================== 请求和响应结构体 ==================

// WorkloadListResponse 工作负载列表响应
type WorkloadListResponse struct {
	Workloads []K8sWorkload `json:"workloads"`
	Total     int           `json:"total"`
}

// EventListResponse 事件列表响应
type EventListResponse struct {
	Events    []K8sEvent        `json:"events"`    // 事件列表
	Total     int               `json:"total"`     // 事件总数
	Namespace string            `json:"namespace,omitempty"` // 命名空间（如果是命名空间级别的查询）
	Filter    map[string]string `json:"filter"`    // 过滤条件
}

// CreateDeploymentRequest 创建Deployment请求
type CreateDeploymentRequest struct {
	Name      string                 `json:"name" binding:"required"`      // Deployment名称
	Replicas  int32                  `json:"replicas"`                    // 副本数，默认1
	Labels    map[string]string      `json:"labels"`                      // 标签
	Template  PodTemplateSpec        `json:"template" binding:"required"` // Pod模板
	Strategy  DeploymentStrategy     `json:"strategy"`                    // 部署策略
}

// CreateStatefulSetRequest 创建StatefulSet请求
type CreateStatefulSetRequest struct {
	Name        string                 `json:"name" binding:"required"`        // StatefulSet名称
	Replicas    int32                  `json:"replicas"`                      // 副本数
	Labels      map[string]string      `json:"labels"`                        // 标签
	Template    PodTemplateSpec        `json:"template" binding:"required"`   // Pod模板
	ServiceName string                 `json:"serviceName" binding:"required"` // 关联Service名称
}

// CreateDaemonSetRequest 创建DaemonSet请求
type CreateDaemonSetRequest struct {
	Name     string                 `json:"name" binding:"required"`     // DaemonSet名称
	Labels   map[string]string      `json:"labels"`                     // 标签
	Template PodTemplateSpec        `json:"template" binding:"required"` // Pod模板
}

// CreateJobRequest 创建Job请求
type CreateJobRequest struct {
	Name        string                 `json:"name" binding:"required"`        // Job名称
	Labels      map[string]string      `json:"labels"`                        // 标签
	Template    PodTemplateSpec        `json:"template" binding:"required"`   // Pod模板
	Completions int32                  `json:"completions"`                   // 完成数量
	Parallelism int32                  `json:"parallelism"`                   // 并行数量
}

// CreateCronJobRequest 创建CronJob请求
type CreateCronJobRequest struct {
	Name     string                 `json:"name" binding:"required"`     // CronJob名称
	Labels   map[string]string      `json:"labels"`                     // 标签
	Schedule string                 `json:"schedule" binding:"required"` // Cron表达式
	JobTemplate JobTemplateSpec     `json:"jobTemplate" binding:"required"` // Job模板
}

// PodTemplateSpec Pod模板规格
type PodTemplateSpec struct {
	Labels     map[string]string `json:"labels"`     // Pod标签
	Containers []ContainerSpec   `json:"containers" binding:"required"` // 容器规格
	Volumes    []VolumeSpec      `json:"volumes"`    // 存储卷规格
	NodeSelector map[string]string `json:"nodeSelector"` // 节点选择器
	Tolerations []Toleration     `json:"tolerations"` // 容忍度
}

// ContainerSpec 容器规格
type ContainerSpec struct {
	Name         string            `json:"name" binding:"required"`  // 容器名称
	Image        string            `json:"image" binding:"required"` // 镜像
	Ports        []ContainerPort   `json:"ports"`                   // 端口配置
	Env          []EnvVar          `json:"env"`                     // 环境变量
	Resources    WorkloadResources `json:"resources"`               // 资源配置
	VolumeMounts []VolumeMount     `json:"volumeMounts"`            // 存储卷挂载
	Command      []string          `json:"command"`                 // 启动命令
	Args         []string          `json:"args"`                    // 启动参数
}

// VolumeSpec 存储卷规格
type VolumeSpec struct {
	Name   string                 `json:"name" binding:"required"` // 卷名称
	Type   string                 `json:"type" binding:"required"` // 卷类型
	Config map[string]interface{} `json:"config"`                 // 卷配置
}

// VolumeMount 存储卷挂载
type VolumeMount struct {
	Name      string `json:"name" binding:"required"`      // 卷名称
	MountPath string `json:"mountPath" binding:"required"` // 挂载路径
	ReadOnly  bool   `json:"readOnly"`                    // 只读模式
}

// DeploymentStrategy 部署策略
type DeploymentStrategy struct {
	Type          string                     `json:"type"`          // 策略类型
	RollingUpdate RollingUpdateDeployment   `json:"rollingUpdate"` // 滚动更新配置
}

// RollingUpdateDeployment 滚动更新部署配置
type RollingUpdateDeployment struct {
	MaxUnavailable string `json:"maxUnavailable"` // 最大不可用数量
	MaxSurge       string `json:"maxSurge"`       // 最大激增数量
}

// JobTemplateSpec Job模板规格
type JobTemplateSpec struct {
	Labels      map[string]string `json:"labels"`      // Job标签
	Template    PodTemplateSpec   `json:"template" binding:"required"` // Pod模板
	Completions int32             `json:"completions"` // 完成数量
	Parallelism int32             `json:"parallelism"` // 并行数量
}

// Toleration 容忍度
type Toleration struct {
	Key      string `json:"key"`      // 键
	Operator string `json:"operator"` // 操作符
	Value    string `json:"value"`    // 值
	Effect   string `json:"effect"`   // 效果
}

// ================== 操作相关请求结构体 ==================

// ScaleWorkloadRequest 伸缩工作负载请求
type ScaleWorkloadRequest struct {
	Replicas int32 `json:"replicas" binding:"required"` // 目标副本数
}

// UpdateWorkloadRequest 更新工作负载请求
type UpdateWorkloadRequest struct {
	Labels    map[string]string `json:"labels"`    // 标签
	Template  PodTemplateSpec   `json:"template"`  // Pod模板
	Strategy  interface{}       `json:"strategy"`  // 部署策略
}

// RestartWorkloadRequest 重启工作负载请求
type RestartWorkloadRequest struct {
	RestartedAt string `json:"restartedAt"` // 重启时间戳
}

// ================== 监控相关结构体 ==================

// PodMetricsInfo Pod监控信息
type PodMetricsInfo struct {
	PodName      string                       `json:"podName"`      // Pod名称
	Namespace    string                       `json:"namespace"`    // 命名空间
	NodeName     string                       `json:"nodeName"`     // 节点名称
	Timestamp    string                       `json:"timestamp"`    // 采集时间
	Containers   []ContainerMetricsInfo       `json:"containers"`   // 容器监控信息列表
	TotalUsage   ResourceUsage                `json:"totalUsage"`   // 总使用量
	ResourceQuota PodResourceQuota            `json:"resourceQuota"` // 资源配额信息
	UsageRate    ResourceUsageRate            `json:"usageRate"`    // 使用率信息
}

// ContainerMetricsInfo 容器监控信息
type ContainerMetricsInfo struct {
	Name         string            `json:"name"`         // 容器名称
	Usage        ResourceUsage     `json:"usage"`        // 资源使用量
	Requests     ResourceUsage     `json:"requests"`     // 资源请求量
	Limits       ResourceUsage     `json:"limits"`       // 资源限制量  
	UsageRate    ResourceUsageRate `json:"usageRate"`    // 使用率
	State        string            `json:"state"`        // 容器状态
	RestartCount int32             `json:"restartCount"` // 重启次数
}

// NodeMetricsInfo 节点监控信息
type NodeMetricsInfo struct {
	NodeName      string                     `json:"nodeName"`      // 节点名称
	Timestamp     string                     `json:"timestamp"`     // 采集时间
	Usage         ResourceUsage              `json:"usage"`         // 资源使用量
	Capacity      ResourceUsage              `json:"capacity"`      // 总容量
	Allocatable   ResourceUsage              `json:"allocatable"`   // 可分配量
	UsageRate     ResourceUsageRate          `json:"usageRate"`     // 使用率
	PodCount      int                        `json:"podCount"`      // Pod数量
	PodMetrics    []PodMetricsSummary        `json:"podMetrics"`    // Pod监控摘要
	SystemInfo    NodeSystemInfo             `json:"systemInfo"`    // 系统信息
}

// PodMetricsSummary Pod监控摘要
type PodMetricsSummary struct {
	PodName   string            `json:"podName"`   // Pod名称
	Namespace string            `json:"namespace"` // 命名空间
	Usage     ResourceUsage     `json:"usage"`     // 资源使用量
	UsageRate ResourceUsageRate `json:"usageRate"` // 使用率
}

// NodeSystemInfo 节点系统信息
type NodeSystemInfo struct {
	KernelVersion           string `json:"kernelVersion"`           // 内核版本
	OSImage                 string `json:"osImage"`                 // 操作系统镜像
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"` // 容器运行时版本
	KubeletVersion          string `json:"kubeletVersion"`          // Kubelet版本
	KubeProxyVersion        string `json:"kubeProxyVersion"`        // KubeProxy版本
	Architecture            string `json:"architecture"`            // 系统架构
}

// NamespaceMetricsInfo 命名空间监控信息
type NamespaceMetricsInfo struct {
	Namespace     string                 `json:"namespace"`     // 命名空间名称
	Timestamp     string                 `json:"timestamp"`     // 采集时间
	PodCount      int                    `json:"podCount"`      // Pod总数
	RunningPods   int                    `json:"runningPods"`   // 运行中的Pod数
	TotalUsage    ResourceUsage          `json:"totalUsage"`    // 总资源使用量
	ResourceQuota NamespaceResourceQuota `json:"resourceQuota"` // 命名空间资源配额
	UsageRate     ResourceUsageRate      `json:"usageRate"`     // 资源使用率
	PodMetrics    []PodMetricsSummary    `json:"podMetrics"`    // Pod监控列表
}

// ResourceUsage 资源使用量
type ResourceUsage struct {
	CPU    string `json:"cpu"`    // CPU使用量 (如: "100m", "1.5")
	Memory string `json:"memory"` // 内存使用量 (如: "128Mi", "1Gi")
}

// ResourceUsageRate 资源使用率
type ResourceUsageRate struct {
	CPURate    float64 `json:"cpuRate"`    // CPU使用率 (百分比: 0-100)
	MemoryRate float64 `json:"memoryRate"` // 内存使用率 (百分比: 0-100)
}

// PodResourceQuota Pod资源配额信息
type PodResourceQuota struct {
	Requests ResourceUsage `json:"requests"` // 资源请求量
	Limits   ResourceUsage `json:"limits"`   // 资源限制量
}

// NamespaceResourceQuota 命名空间资源配额
type NamespaceResourceQuota struct {
	Hard ResourceUsage `json:"hard"` // 硬限制
	Used ResourceUsage `json:"used"` // 已使用
}

// MetricsListResponse 监控数据列表响应
type MetricsListResponse struct {
	Items interface{} `json:"items"` // 监控数据列表，可以是Pod、Node或Namespace监控信息
	Total int         `json:"total"` // 总数
}

// MetricsQueryRequest 监控查询请求
type MetricsQueryRequest struct {
	ResourceType string   `json:"resourceType"` // 资源类型: pod/node/namespace
	Names        []string `json:"names"`        // 资源名称列表，为空则获取全部
	StartTime    string   `json:"startTime"`    // 开始时间 (可选，用于未来支持历史数据)
	EndTime      string   `json:"endTime"`      // 结束时间 (可选，用于未来支持历史数据)
}

// ================== Pod YAML管理相关结构体 ==================

// CreatePodFromYAMLRequest 通过YAML创建Pod请求
type CreatePodFromYAMLRequest struct {
	YAMLContent  string `json:"yamlContent" binding:"required"`  // YAML内容
	DryRun       bool   `json:"dryRun"`                          // 是否只进行校验不实际创建
	ValidateOnly bool   `json:"validateOnly"`                    // 是否只校验YAML格式
}

// ValidateYAMLRequest YAML校验请求
type ValidateYAMLRequest struct {
	YAMLContent  string `json:"yamlContent" binding:"required"`  // YAML内容
	ResourceType string `json:"resourceType"`                    // 资源类型，如pod、deployment等
}

// ValidateYAMLResponse YAML校验响应
type ValidateYAMLResponse struct {
	Valid        bool                   `json:"valid"`        // 是否有效
	Errors       []string               `json:"errors"`       // 错误列表
	Warnings     []string               `json:"warnings"`     // 警告列表
	Suggestions  []string               `json:"suggestions"`  // 建议列表
	ParsedObject map[string]interface{} `json:"parsedObject"` // 解析后的对象
}

// CreatePodFromYAMLResponse 通过YAML创建Pod响应
type CreatePodFromYAMLResponse struct {
	Success      bool                   `json:"success"`      // 是否创建成功
	PodName      string                 `json:"podName"`      // 创建的Pod名称
	Namespace    string                 `json:"namespace"`    // Pod所在的命名空间
	Message      string                 `json:"message"`      // 响应消息
	ValidationResult *ValidateYAMLResponse `json:"validationResult,omitempty"` // 校验结果（DryRun时返回）
	ParsedObject map[string]interface{} `json:"parsedObject,omitempty"` // 解析的对象信息
}

// UpdatePodYAMLRequest 通过YAML更新Pod请求
type UpdatePodYAMLRequest struct {
	YAMLContent  string `json:"yamlContent" binding:"required"`  // YAML内容
	DryRun       bool   `json:"dryRun"`                          // 是否只进行校验不实际更新
	ValidateOnly bool   `json:"validateOnly"`                    // 是否只校验YAML格式
	Force        bool   `json:"force"`                           // 是否强制更新（删除重建）
}

// UpdatePodYAMLResponse 通过YAML更新Pod响应
type UpdatePodYAMLResponse struct {
	Success          bool                   `json:"success"`          // 是否更新成功
	PodName          string                 `json:"podName"`          // 更新的Pod名称
	Namespace        string                 `json:"namespace"`        // Pod所在的命名空间
	Message          string                 `json:"message"`          // 响应消息
	UpdateStrategy   string                 `json:"updateStrategy"`   // 更新策略 (patch/recreate)
	ValidationResult *ValidateYAMLResponse  `json:"validationResult,omitempty"` // 校验结果（DryRun时返回）
	Changes          []string               `json:"changes,omitempty"`          // 变更说明
	Warnings         []string               `json:"warnings,omitempty"`         // 警告信息
}

// ================== 通用工作负载YAML管理相关结构体 ==================

// GetWorkloadYAMLRequest 获取工作负载YAML请求
type GetWorkloadYAMLRequest struct {
	WorkloadType string `json:"workloadType" binding:"required"` // 工作负载类型: deployment,statefulset,daemonset,job,cronjob
	WorkloadName string `json:"workloadName" binding:"required"` // 工作负载名称
}

// UpdateWorkloadYAMLRequest 通用工作负载YAML更新请求
type UpdateWorkloadYAMLRequest struct {
	WorkloadType string `json:"workloadType" binding:"required"` // 工作负载类型: deployment,statefulset,daemonset,job,cronjob
	WorkloadName string `json:"workloadName" binding:"required"` // 工作负载名称
	YAMLContent  string `json:"yamlContent" binding:"required"`  // YAML内容
	DryRun       bool   `json:"dryRun"`                          // 是否只进行校验不实际更新
	ValidateOnly bool   `json:"validateOnly"`                    // 是否只校验YAML格式
	Force        bool   `json:"force"`                           // 是否强制更新
}

// UpdateWorkloadYAMLResponse 通用工作负载YAML更新响应
type UpdateWorkloadYAMLResponse struct {
	Success          bool                   `json:"success"`          // 是否更新成功
	WorkloadType     string                 `json:"workloadType"`     // 工作负载类型
	WorkloadName     string                 `json:"workloadName"`     // 工作负载名称
	Namespace        string                 `json:"namespace"`        // 命名空间
	Message          string                 `json:"message"`          // 响应消息
	UpdateStrategy   string                 `json:"updateStrategy"`   // 更新策略 (patch/update/rolling)
	ValidationResult *ValidateYAMLResponse  `json:"validationResult,omitempty"` // 校验结果（DryRun时返回）
	Changes          []string               `json:"changes,omitempty"`          // 变更说明
	Warnings         []string               `json:"warnings,omitempty"`         // 警告信息
	AppliedAt        string                 `json:"appliedAt,omitempty"`        // 应用时间
}

// GetWorkloadYAMLResponse 获取工作负载YAML响应
type GetWorkloadYAMLResponse struct {
	Success      bool   `json:"success"`      // 是否获取成功
	WorkloadType string `json:"workloadType"` // 工作负载类型
	WorkloadName string `json:"workloadName"` // 工作负载名称
	Namespace    string `json:"namespace"`    // 命名空间
	YAMLContent  string `json:"yamlContent"`  // YAML内容
	Message      string `json:"message"`      // 响应消息
}

// ================== Deployment 回滚相关结构体 ==================

// DeploymentRevision Deployment版本信息
type DeploymentRevision struct {
	Revision       int64             `json:"revision"`       // 版本号
	CreationTime   string            `json:"creationTime"`   // 创建时间
	ChangeReason   string            `json:"changeReason"`   // 变更原因
	Images         []string          `json:"images"`         // 镜像列表
	Labels         map[string]string `json:"labels"`         // 标签
	Annotations    map[string]string `json:"annotations"`    // 注释
	Status         string            `json:"status"`         // 版本状态 (current/historical)
	ReplicasSummary ReplicasSummary   `json:"replicasSummary"` // 副本统计
}

// ReplicasSummary 副本统计信息
type ReplicasSummary struct {
	Desired   int32 `json:"desired"`   // 期望副本数
	Current   int32 `json:"current"`   // 当前副本数
	Updated   int32 `json:"updated"`   // 已更新副本数
	Ready     int32 `json:"ready"`     // 就绪副本数
	Available int32 `json:"available"` // 可用副本数
}

// DeploymentRolloutHistoryResponse Deployment版本历史响应
type DeploymentRolloutHistoryResponse struct {
	DeploymentName string               `json:"deploymentName"` // Deployment名称
	Namespace      string               `json:"namespace"`      // 命名空间
	CurrentRevision int64               `json:"currentRevision"` // 当前版本号
	TotalRevisions  int                 `json:"totalRevisions"` // 总版本数
	Revisions      []DeploymentRevision `json:"revisions"`      // 版本列表
}

// GetDeploymentRevisionRequest 获取Deployment版本详情请求
type GetDeploymentRevisionRequest struct {
	Revision int64 `json:"revision" binding:"required"` // 版本号
}

// DeploymentRevisionDetail Deployment版本详情
type DeploymentRevisionDetail struct {
	DeploymentRevision
	PodTemplate    PodTemplateSpec       `json:"podTemplate"`    // Pod模板
	Strategy       DeploymentStrategy    `json:"strategy"`       // 部署策略
	Conditions     []WorkloadCondition   `json:"conditions"`     // 状态条件
	Events         []K8sEvent           `json:"events"`         // 相关事件
	ReplicaSets    []ReplicaSetInfo     `json:"replicaSets"`    // 关联的ReplicaSet信息
}

// ReplicaSetInfo ReplicaSet信息
type ReplicaSetInfo struct {
	Name         string  `json:"name"`         // ReplicaSet名称
	Revision     int64   `json:"revision"`     // 版本号
	Replicas     int32   `json:"replicas"`     // 副本数
	ReadyReplicas int32  `json:"readyReplicas"` // 就绪副本数
	Status       string  `json:"status"`       // 状态
	CreatedAt    string  `json:"createdAt"`    // 创建时间
}

// RollbackDeploymentRequest Deployment回滚请求
type RollbackDeploymentRequest struct {
	ToRevision int64 `json:"toRevision"` // 目标版本号，0表示回滚到上一版本
}

// RollbackDeploymentResponse Deployment回滚响应
type RollbackDeploymentResponse struct {
	Success         bool   `json:"success"`         // 是否回滚成功
	Message         string `json:"message"`         // 响应消息
	FromRevision    int64  `json:"fromRevision"`    // 回滚前版本号
	ToRevision      int64  `json:"toRevision"`      // 回滚后版本号
	DeploymentName  string `json:"deploymentName"`  // Deployment名称
	Namespace       string `json:"namespace"`       // 命名空间
	RolloutStatus   string `json:"rolloutStatus"`   // 滚动发布状态
}

// PauseDeploymentResponse 暂停Deployment响应
type PauseDeploymentResponse struct {
	Success        bool   `json:"success"`        // 是否暂停成功
	Message        string `json:"message"`        // 响应消息
	DeploymentName string `json:"deploymentName"` // Deployment名称
	Namespace      string `json:"namespace"`      // 命名空间
	Status         string `json:"status"`         // 当前状态
}

// ResumeDeploymentResponse 恢复Deployment响应
type ResumeDeploymentResponse struct {
	Success        bool   `json:"success"`        // 是否恢复成功
	Message        string `json:"message"`        // 响应消息
	DeploymentName string `json:"deploymentName"` // Deployment名称
	Namespace      string `json:"namespace"`      // 命名空间
	Status         string `json:"status"`         // 当前状态
}

// DeploymentRolloutStatusResponse Deployment滚动发布状态响应
type DeploymentRolloutStatusResponse struct {
	DeploymentName  string               `json:"deploymentName"`  // Deployment名称
	Namespace       string               `json:"namespace"`       // 命名空间
	CurrentRevision int64                `json:"currentRevision"` // 当前版本号
	UpdatedReplicas int32                `json:"updatedReplicas"` // 已更新副本数
	ReadyReplicas   int32                `json:"readyReplicas"`   // 就绪副本数
	AvailableReplicas int32              `json:"availableReplicas"` // 可用副本数
	ObservedGeneration int64             `json:"observedGeneration"` // 观察到的代数
	Conditions      []WorkloadCondition  `json:"conditions"`      // 状态条件
	Strategy        DeploymentStrategy   `json:"strategy"`        // 部署策略
	Paused          bool                 `json:"paused"`          // 是否已暂停
	ProgressDeadline int32               `json:"progressDeadline"` // 进度截止时间
	RolloutComplete  bool                `json:"rolloutComplete"` // 是否滚动发布完成
	Status          string               `json:"status"`          // 总体状态 (Progressing/Complete/Failed/Paused)
}

// ================== Service 相关结构体 ==================

// K8sService K8s Service信息
type K8sService struct {
	Name         string                 `json:"name"`         // 服务名称
	Namespace    string                 `json:"namespace"`    // 命名空间
	Labels       map[string]string      `json:"labels"`       // 标签
	Type         string                 `json:"type"`         // 服务类型 ClusterIP/NodePort/LoadBalancer/ExternalName
	Selector     map[string]string      `json:"selector"`     // 选择器
	ClusterIP    string                 `json:"clusterIP"`    // 集群IP
	ExternalIPs  []string               `json:"externalIPs"`  // 外部IP列表
	Ports        []ServicePort          `json:"ports"`        // 端口配置
	Endpoints    []ServiceEndpoint      `json:"endpoints"`    // 端点信息
	CreatedAt    string                 `json:"createdAt"`    // 创建时间
	Status       string                 `json:"status"`       // 状态
}

// ServicePort 服务端口配置
type ServicePort struct {
	Name       string `json:"name"`       // 端口名称
	Protocol   string `json:"protocol"`   // 协议 TCP/UDP
	Port       int32  `json:"port"`       // 服务端口
	TargetPort string `json:"targetPort"` // 目标端口
	NodePort   int32  `json:"nodePort"`   // 节点端口(NodePort类型时使用)
}

// ServiceEndpoint 服务端点
type ServiceEndpoint struct {
	IP       string          `json:"ip"`       // 端点IP
	Hostname string          `json:"hostname"` // 主机名
	NodeName string          `json:"nodeName"` // 节点名称
	Ready    bool            `json:"ready"`    // 就绪状态
	Ports    []EndpointPort  `json:"ports"`    // 端口信息
}

// EndpointPort 端点端口
type EndpointPort struct {
	Name     string `json:"name"`     // 端口名称
	Port     int32  `json:"port"`     // 端口号
	Protocol string `json:"protocol"` // 协议
}

// ServiceDetail Service详情
type ServiceDetail struct {
	K8sService
	Events     []K8sEvent       `json:"events"`     // 相关事件
	Pods       []K8sPodInfo     `json:"pods"`       // 关联的Pod列表
	Spec       interface{}      `json:"spec"`       // 完整规格配置
}

// CreateServiceRequest 创建Service请求
type CreateServiceRequest struct {
	Name        string            `json:"name" binding:"required"`        // 服务名称
	Labels      map[string]string `json:"labels"`                         // 标签
	Type        string            `json:"type" binding:"required"`        // 服务类型
	Selector    map[string]string `json:"selector" binding:"required"`    // 选择器
	Ports       []ServicePortSpec `json:"ports" binding:"required"`       // 端口配置
	ExternalIPs []string          `json:"externalIPs"`                    // 外部IP列表
}

// ServicePortSpec 服务端口规格
type ServicePortSpec struct {
	Name       string `json:"name"`                        // 端口名称
	Protocol   string `json:"protocol" binding:"required"` // 协议
	Port       int32  `json:"port" binding:"required"`     // 服务端口
	TargetPort string `json:"targetPort"`                  // 目标端口
	NodePort   int32  `json:"nodePort"`                    // 节点端口(NodePort类型时使用)
}

// UpdateServiceRequest 更新Service请求
type UpdateServiceRequest struct {
	Labels      map[string]string `json:"labels"`      // 标签
	Type        string            `json:"type"`        // 服务类型
	Selector    map[string]string `json:"selector"`    // 选择器
	Ports       []ServicePortSpec `json:"ports"`       // 端口配置
	ExternalIPs []string          `json:"externalIPs"` // 外部IP列表
}

// ServiceListResponse 服务列表响应
type ServiceListResponse struct {
	Services []K8sService `json:"services"`
	Total    int          `json:"total"`
}

// ================== Ingress 相关结构体 ==================

// K8sIngress K8s Ingress信息
type K8sIngress struct {
	Name              string                 `json:"name"`              // Ingress名称
	Namespace         string                 `json:"namespace"`         // 命名空间
	Labels            map[string]string      `json:"labels"`            // 标签
	Class             string                 `json:"class"`             // Ingress类/控制器类型
	ControllerName    string                 `json:"controllerName"`    // Controller名称
	ControllerVersion string                 `json:"controllerVersion"` // Controller版本
	Type              string                 `json:"type"`              // Ingress类型：公网Nginx/内网Nginx等
	Rules             []IngressRule          `json:"rules"`             // 路由规则
	TLS               []IngressTLS           `json:"tls"`               // TLS配置
	LoadBalancer      IngressLoadBalancer    `json:"loadBalancer"`      // 负载均衡器信息
	Endpoints         []string               `json:"endpoints"`         // 访问端点
	CreatedAt         string                 `json:"createdAt"`         // 创建时间
	Status            string                 `json:"status"`            // 状态
}

// IngressRule Ingress路由规则
type IngressRule struct {
	Host string           `json:"host"` // 主机名
	HTTP IngressRuleValue `json:"http"` // HTTP规则
}

// IngressRuleValue Ingress规则值
type IngressRuleValue struct {
	Paths []IngressPath `json:"paths"` // 路径规则
}

// IngressPath Ingress路径规则
type IngressPath struct {
	Path     string             `json:"path"`     // 路径
	PathType string             `json:"pathType"` // 路径类型 Exact/Prefix/ImplementationSpecific
	Backend  IngressBackend     `json:"backend"`  // 后端服务
}

// IngressBackend Ingress后端服务
type IngressBackend struct {
	Service IngressServiceBackend `json:"service"` // 服务后端
}

// IngressServiceBackend Ingress服务后端
type IngressServiceBackend struct {
	Name string             `json:"name"` // 服务名称
	Port IngressServicePort `json:"port"` // 服务端口
}

// IngressServicePort Ingress服务端口
type IngressServicePort struct {
	Number int32  `json:"number"` // 端口号
	Name   string `json:"name"`   // 端口名称
}

// IngressTLS Ingress TLS配置
type IngressTLS struct {
	Hosts      []string `json:"hosts"`      // 主机列表
	SecretName string   `json:"secretName"` // 证书Secret名称
}

// IngressLoadBalancer Ingress负载均衡器信息
type IngressLoadBalancer struct {
	Ingress []IngressLoadBalancerIngress `json:"ingress"` // 入口信息
}

// IngressLoadBalancerIngress 负载均衡器入口信息
type IngressLoadBalancerIngress struct {
	IP       string                       `json:"ip"`       // IP地址
	Hostname string                       `json:"hostname"` // 主机名
	Ports    []IngressLoadBalancerPort    `json:"ports"`    // 端口信息
}

// IngressLoadBalancerPort 负载均衡器端口信息
type IngressLoadBalancerPort struct {
	Port     int32  `json:"port"`     // 端口号
	Protocol string `json:"protocol"` // 协议
}

// IngressDetail Ingress详情
type IngressDetail struct {
	K8sIngress
	Events []K8sEvent  `json:"events"` // 相关事件
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreateIngressRequest 创建Ingress请求
type CreateIngressRequest struct {
	Name        string                 `json:"name" binding:"required"`        // Ingress名称
	Labels      map[string]string      `json:"labels"`                         // 标签
	Class       string                 `json:"class"`                          // Ingress类
	Rules       []IngressRuleSpec      `json:"rules" binding:"required"`       // 路由规则
	TLS         []IngressTLSSpec       `json:"tls"`                            // TLS配置
	Annotations map[string]string      `json:"annotations"`                    // 注释
}

// IngressRuleSpec Ingress规则规格
type IngressRuleSpec struct {
	Host  string                  `json:"host" binding:"required"` // 主机名
	Paths []IngressPathSpec       `json:"paths" binding:"required"` // 路径规则
}

// IngressPathSpec Ingress路径规格
type IngressPathSpec struct {
	Path        string                    `json:"path" binding:"required"`     // 路径
	PathType    string                    `json:"pathType" binding:"required"` // 路径类型
	ServiceName string                    `json:"serviceName" binding:"required"` // 服务名称
	ServicePort int32                     `json:"servicePort" binding:"required"` // 服务端口
}

// IngressTLSSpec Ingress TLS规格
type IngressTLSSpec struct {
	Hosts      []string `json:"hosts" binding:"required"`      // 主机列表
	SecretName string   `json:"secretName" binding:"required"` // 证书Secret名称
}

// UpdateIngressRequest 更新Ingress请求
type UpdateIngressRequest struct {
	Labels      map[string]string      `json:"labels"`      // 标签
	Class       string                 `json:"class"`       // Ingress类
	Rules       []IngressRuleSpec      `json:"rules"`       // 路由规则
	TLS         []IngressTLSSpec       `json:"tls"`         // TLS配置
	Annotations map[string]string      `json:"annotations"` // 注释
}

// IngressListResponse Ingress列表响应
type IngressListResponse struct {
	Ingresses []K8sIngress `json:"ingresses"`
	Total     int          `json:"total"`
}

// ================== 存储管理相关结构体 ==================

// K8sPersistentVolumeClaim PVC持久卷声明信息
type K8sPersistentVolumeClaim struct {
	Name         string            `json:"name"`         // PVC名称
	Namespace    string            `json:"namespace"`    // 命名空间
	Labels       map[string]string `json:"labels"`       // 标签
	Capacity     string            `json:"capacity"`     // 总量
	AccessModes  []string          `json:"accessModes"`  // 访问模式
	Status       string            `json:"status"`       // 状态 Pending/Bound/Lost
	StorageClass string            `json:"storageClass"` // 存储类型
	VolumeMode   string            `json:"volumeMode"`   // 卷模式 Filesystem/Block
	VolumeName   string            `json:"volumeName"`   // 关联的存储卷
	CreatedAt    string            `json:"createdAt"`    // 创建时间
}

// PVCDetail PVC详情
type PVCDetail struct {
	K8sPersistentVolumeClaim
	Events []K8sEvent  `json:"events"` // 相关事件
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreatePVCRequest 创建PVC请求
type CreatePVCRequest struct {
	Name         string            `json:"name" binding:"required"`         // PVC名称
	Labels       map[string]string `json:"labels"`                          // 标签
	StorageClass string            `json:"storageClass"`                    // 存储类名称
	AccessModes  []string          `json:"accessModes" binding:"required"`  // 访问模式
	Resources    PVCResourcesSpec  `json:"resources" binding:"required"`    // 资源配置
	VolumeMode   string            `json:"volumeMode"`                      // 卷模式
	Selector     *PVCSelectorSpec  `json:"selector"`                        // 标签选择器
}

// PVCResourcesSpec PVC资源规格
type PVCResourcesSpec struct {
	Requests map[string]string `json:"requests" binding:"required"` // 资源请求
	Limits   map[string]string `json:"limits"`                      // 资源限制
}

// PVCSelectorSpec PVC选择器规格
type PVCSelectorSpec struct {
	MatchLabels      map[string]string `json:"matchLabels"`      // 匹配标签
	MatchExpressions []PVCMatchExp     `json:"matchExpressions"` // 匹配表达式
}

// PVCMatchExp PVC匹配表达式
type PVCMatchExp struct {
	Key      string   `json:"key"`      // 键
	Operator string   `json:"operator"` // 操作符
	Values   []string `json:"values"`   // 值
}

// UpdatePVCRequest 更新PVC请求
type UpdatePVCRequest struct {
	Labels    map[string]string `json:"labels"`    // 标签
	Resources PVCResourcesSpec  `json:"resources"` // 资源配置
}

// PVCListResponse PVC列表响应
type PVCListResponse struct {
	PVCs  []K8sPersistentVolumeClaim `json:"pvcs"`
	Total int                        `json:"total"`
}

// K8sPersistentVolume PV持久卷信息
type K8sPersistentVolume struct {
	Name                     string            `json:"name"`                     // PV名称
	Labels                   map[string]string `json:"labels"`                   // 标签
	Capacity                 string            `json:"capacity"`                 // 总量
	AccessModes              []string          `json:"accessModes"`              // 访问模式
	ReclaimPolicy            string            `json:"reclaimPolicy"`            // 回收策略 Retain/Delete/Recycle
	Status                   string            `json:"status"`                   // 状态 Available/Bound/Released/Failed
	StorageClass             string            `json:"storageClass"`             // 存储类型
	VolumeMode               string            `json:"volumeMode"`               // 卷模式
	ClaimRef                 *PVClaimRef       `json:"claimRef"`                 // 绑定存储声明
	PersistentVolumeSource   PVSource          `json:"persistentVolumeSource"`   // 存储源
	NodeAffinity             *PVNodeAffinity   `json:"nodeAffinity"`             // 节点亲和性
	MountOptions             []string          `json:"mountOptions"`             // 挂载选项
	CreatedAt                string            `json:"createdAt"`                // 创建时间
}

// PVClaimRef PV声明引用
type PVClaimRef struct {
	Kind       string `json:"kind"`       // 类型
	Namespace  string `json:"namespace"`  // 命名空间
	Name       string `json:"name"`       // 名称
	UID        string `json:"uid"`        // UID
	APIVersion string `json:"apiVersion"` // API版本
}

// PVSource PV存储源
type PVSource struct {
	HostPath          *PVHostPathVolumeSource          `json:"hostPath,omitempty"`          // HostPath存储源
	NFS               *PVNFSVolumeSource               `json:"nfs,omitempty"`               // NFS存储源
	ISCSI             *PVISCSIVolumeSource             `json:"iscsi,omitempty"`             // iSCSI存储源
	CSI               *PVCSIVolumeSource               `json:"csi,omitempty"`               // CSI存储源
	Local             *PVLocalVolumeSource             `json:"local,omitempty"`             // Local存储源
	AWSElasticBlockStore *PVAWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty"` // AWS EBS
	// 可以根据需要添加更多存储源类型
}

// PVHostPathVolumeSource HostPath存储源
type PVHostPathVolumeSource struct {
	Path string `json:"path"` // 主机路径
	Type string `json:"type"` // 类型
}

// PVNFSVolumeSource NFS存储源
type PVNFSVolumeSource struct {
	Server   string `json:"server"`   // NFS服务器
	Path     string `json:"path"`     // NFS路径
	ReadOnly bool   `json:"readOnly"` // 只读
}

// PVISCSIVolumeSource iSCSI存储源
type PVISCSIVolumeSource struct {
	TargetPortal string   `json:"targetPortal"` // 目标门户
	IQN          string   `json:"iqn"`          // IQN
	Lun          int32    `json:"lun"`          // LUN
	ISCSIInterface string `json:"iscsiInterface"` // iSCSI接口
	FSType       string   `json:"fsType"`       // 文件系统类型
	ReadOnly     bool     `json:"readOnly"`     // 只读
	Portals      []string `json:"portals"`      // 门户列表
}

// PVCSIVolumeSource CSI存储源
type PVCSIVolumeSource struct {
	Driver       string            `json:"driver"`       // 驱动名称
	VolumeHandle string            `json:"volumeHandle"` // 卷句柄
	ReadOnly     bool              `json:"readOnly"`     // 只读
	FSType       string            `json:"fsType"`       // 文件系统类型
	VolumeAttributes map[string]string `json:"volumeAttributes"` // 卷属性
}

// PVLocalVolumeSource Local存储源
type PVLocalVolumeSource struct {
	Path string `json:"path"` // 本地路径
	FSType string `json:"fsType"` // 文件系统类型
}

// PVAWSElasticBlockStoreVolumeSource AWS EBS存储源
type PVAWSElasticBlockStoreVolumeSource struct {
	VolumeID  string `json:"volumeID"`  // 卷ID
	FSType    string `json:"fsType"`    // 文件系统类型
	Partition int32  `json:"partition"` // 分区
	ReadOnly  bool   `json:"readOnly"`  // 只读
}

// PVNodeAffinity PV节点亲和性
type PVNodeAffinity struct {
	Required *PVNodeSelector `json:"required"` // 必需的节点选择器
}

// PVNodeSelector PV节点选择器
type PVNodeSelector struct {
	NodeSelectorTerms []PVNodeSelectorTerm `json:"nodeSelectorTerms"` // 节点选择器条件
}

// PVNodeSelectorTerm PV节点选择器条件
type PVNodeSelectorTerm struct {
	MatchExpressions []PVNodeSelectorRequirement `json:"matchExpressions"` // 匹配表达式
	MatchFields      []PVNodeSelectorRequirement `json:"matchFields"`      // 匹配字段
}

// PVNodeSelectorRequirement PV节点选择器需求
type PVNodeSelectorRequirement struct {
	Key      string   `json:"key"`      // 键
	Operator string   `json:"operator"` // 操作符
	Values   []string `json:"values"`   // 值
}

// PVDetail PV详情
type PVDetail struct {
	K8sPersistentVolume
	Events []K8sEvent  `json:"events"` // 相关事件
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreatePVRequest 创建PV请求
type CreatePVRequest struct {
	Name                   string            `json:"name" binding:"required"`                   // PV名称
	Labels                 map[string]string `json:"labels"`                                    // 标签
	Capacity               map[string]string `json:"capacity" binding:"required"`               // 容量
	AccessModes            []string          `json:"accessModes" binding:"required"`            // 访问模式
	ReclaimPolicy          string            `json:"reclaimPolicy"`                             // 回收策略
	StorageClassName       string            `json:"storageClassName"`                          // 存储类名称
	VolumeMode             string            `json:"volumeMode"`                                // 卷模式
	PersistentVolumeSource PVSourceSpec      `json:"persistentVolumeSource" binding:"required"` // 存储源
	NodeAffinity           *PVNodeAffinity   `json:"nodeAffinity"`                              // 节点亲和性
	MountOptions           []string          `json:"mountOptions"`                              // 挂载选项
}

// PVSourceSpec PV存储源规格
type PVSourceSpec struct {
	HostPath *PVHostPathVolumeSource `json:"hostPath,omitempty"` // HostPath存储源
	NFS      *PVNFSVolumeSource      `json:"nfs,omitempty"`      // NFS存储源
	ISCSI    *PVISCSIVolumeSource    `json:"iscsi,omitempty"`    // iSCSI存储源
	CSI      *PVCSIVolumeSource      `json:"csi,omitempty"`      // CSI存储源
	Local    *PVLocalVolumeSource    `json:"local,omitempty"`    // Local存储源
	AWSElasticBlockStore *PVAWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty"` // AWS EBS
}

// UpdatePVRequest 更新PV请求
type UpdatePVRequest struct {
	Labels        map[string]string `json:"labels"`        // 标签
	ReclaimPolicy string            `json:"reclaimPolicy"` // 回收策略
	MountOptions  []string          `json:"mountOptions"`  // 挂载选项
}

// PVListResponse PV列表响应
type PVListResponse struct {
	PVs   []K8sPersistentVolume `json:"pvs"`
	Total int                   `json:"total"`
}

// K8sStorageClass 存储类信息
type K8sStorageClass struct {
	Name                 string                 `json:"name"`                 // 存储类名称
	Labels               map[string]string      `json:"labels"`               // 标签
	Provisioner          string                 `json:"provisioner"`          // 提供者
	Parameters           map[string]string      `json:"parameters"`           // 参数
	ReclaimPolicy        string                 `json:"reclaimPolicy"`        // 回收策略
	VolumeBindingMode    string                 `json:"volumeBindingMode"`    // 卷绑定模式
	AllowVolumeExpansion bool                   `json:"allowVolumeExpansion"` // 允许卷扩展
	MountOptions         []string               `json:"mountOptions"`         // 挂载选项
	AllowedTopologies    []StorageClassTopology `json:"allowedTopologies"`    // 允许的拓扑
	CreatedAt            string                 `json:"createdAt"`            // 创建时间
}

// StorageClassTopology 存储类拓扑
type StorageClassTopology struct {
	MatchLabelExpressions []StorageClassTopologyExp `json:"matchLabelExpressions"` // 匹配标签表达式
}

// StorageClassTopologyExp 存储类拓扑表达式
type StorageClassTopologyExp struct {
	Key    string   `json:"key"`    // 键
	Values []string `json:"values"` // 值
}

// StorageClassDetail 存储类详情
type StorageClassDetail struct {
	K8sStorageClass
	Events []K8sEvent  `json:"events"` // 相关事件
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreateStorageClassRequest 创建存储类请求
type CreateStorageClassRequest struct {
	Name                 string                    `json:"name" binding:"required"`                 // 存储类名称
	Labels               map[string]string         `json:"labels"`                                  // 标签
	Provisioner          string                    `json:"provisioner" binding:"required"`          // 提供者
	Parameters           map[string]string         `json:"parameters"`                              // 参数
	ReclaimPolicy        string                    `json:"reclaimPolicy"`                           // 回收策略
	VolumeBindingMode    string                    `json:"volumeBindingMode"`                       // 卷绑定模式
	AllowVolumeExpansion bool                      `json:"allowVolumeExpansion"`                    // 允许卷扩展
	MountOptions         []string                  `json:"mountOptions"`                            // 挂载选项
	AllowedTopologies    []StorageClassTopology    `json:"allowedTopologies"`                       // 允许的拓扑
}

// UpdateStorageClassRequest 更新存储类请求
type UpdateStorageClassRequest struct {
	Labels               map[string]string      `json:"labels"`               // 标签
	Parameters           map[string]string      `json:"parameters"`           // 参数
	AllowVolumeExpansion bool                   `json:"allowVolumeExpansion"` // 允许卷扩展
	MountOptions         []string               `json:"mountOptions"`         // 挂载选项
	AllowedTopologies    []StorageClassTopology `json:"allowedTopologies"`    // 允许的拓扑
}

// StorageClassListResponse 存储类列表响应
type StorageClassListResponse struct {
	StorageClasses []K8sStorageClass `json:"storageClasses"`
	Total          int               `json:"total"`
}

// ===================== 配置管理数据模型 =====================

// K8sConfigMap ConfigMap 数据模型
type K8sConfigMap struct {
	Name        string            `json:"name"`        // ConfigMap名称
	Namespace   string            `json:"namespace"`   // 命名空间
	Labels      map[string]string `json:"labels"`      // 标签
	Data        map[string]string `json:"data"`        // 数据
	BinaryData  map[string][]byte `json:"binaryData"`  // 二进制数据
	Immutable   bool              `json:"immutable"`   // 是否不可变
	CreatedTime string            `json:"createdTime"` // 创建时间
}

// ConfigMapDetail ConfigMap详情
type ConfigMapDetail struct {
	K8sConfigMap
	Events []K8sEvent  `json:"events"` // 相关事件
	Usage  []string    `json:"usage"`  // 使用情况（哪些Pod在使用）
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreateConfigMapRequest 创建ConfigMap请求
type CreateConfigMapRequest struct {
	Name       string            `json:"name" binding:"required"` // ConfigMap名称
	Labels     map[string]string `json:"labels"`                  // 标签
	Data       map[string]string `json:"data"`                    // 数据
	BinaryData map[string][]byte `json:"binaryData"`              // 二进制数据
	Immutable  bool              `json:"immutable"`               // 是否不可变
}

// UpdateConfigMapRequest 更新ConfigMap请求
type UpdateConfigMapRequest struct {
	Labels     map[string]string `json:"labels"`     // 标签
	Data       map[string]string `json:"data"`       // 数据
	BinaryData map[string][]byte `json:"binaryData"` // 二进制数据
}

// ConfigMapListResponse ConfigMap列表响应
type ConfigMapListResponse struct {
	ConfigMaps []K8sConfigMap `json:"configMaps"`
	Total      int            `json:"total"`
}

// K8sSecret Secret 数据模型
type K8sSecret struct {
	Name        string            `json:"name"`        // Secret名称
	Namespace   string            `json:"namespace"`   // 命名空间
	Labels      map[string]string `json:"labels"`      // 标签
	Type        string            `json:"type"`        // Secret类型
	Data        map[string][]byte `json:"data"`        // 数据(base64编码)
	StringData  map[string]string `json:"stringData"`  // 字符串数据
	Immutable   bool              `json:"immutable"`   // 是否不可变
	CreatedTime string            `json:"createdTime"` // 创建时间
}

// SecretDetail Secret详情
type SecretDetail struct {
	K8sSecret
	Events []K8sEvent  `json:"events"` // 相关事件
	Usage  []string    `json:"usage"`  // 使用情况（哪些Pod在使用）
	Spec   interface{} `json:"spec"`   // 完整规格配置
}

// CreateSecretRequest 创建Secret请求
type CreateSecretRequest struct {
	Name       string            `json:"name" binding:"required"` // Secret名称
	Labels     map[string]string `json:"labels"`                  // 标签
	Type       string            `json:"type"`                    // Secret类型
	Data       map[string][]byte `json:"data"`                    // 数据(base64编码)
	StringData map[string]string `json:"stringData"`              // 字符串数据
	Immutable  bool              `json:"immutable"`               // 是否不可变
}

// UpdateSecretRequest 更新Secret请求
type UpdateSecretRequest struct {
	Labels     map[string]string `json:"labels"`     // 标签
	Data       map[string][]byte `json:"data"`       // 数据(base64编码)
	StringData map[string]string `json:"stringData"` // 字符串数据
}

// SecretListResponse Secret列表响应
type SecretListResponse struct {
	Secrets []K8sSecret `json:"secrets"`
	Total   int         `json:"total"`
}