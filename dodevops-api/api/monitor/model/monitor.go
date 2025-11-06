package model

type HostMetrics struct {
	CPUUsage     float64 `json:"cpuUsage"`     // CPU使用率百分比
	MemoryUsage  float64 `json:"memoryUsage"`  // 内存使用率百分比
	DiskUsage    float64 `json:"diskUsage"`    // 磁盘使用率百分比
	OnlineStatus int     `json:"onlineStatus"` // 在线状态(0:在线,1:离线)
}

type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"` // 用于matrix类型
			Value  []interface{}     `json:"value"`  // 用于vector类型
		} `json:"result"`
	} `json:"data"`
}

// HostMetricHistory 主机指标历史数据
type HostMetricHistory struct {
	HostID   uint              `json:"hostId"`   // 主机ID
	Metric   string            `json:"metric"`   // 指标类型(cpu/memory/disk)
	TimeData []MetricDataPoint `json:"timeData"` // 时间序列数据
}

// MetricDataPoint 指标数据点
type MetricDataPoint struct {
	Timestamp int64   `json:"timestamp"` // 时间戳(秒)
	Value     float64 `json:"value"`     // 指标值
}

// AllMetricsHistory 主机所有指标历史数据
type AllMetricsHistory struct {
	HostID         uint              `json:"hostId"`         // 主机ID
	CPU            []MetricDataPoint `json:"cpu"`            // CPU使用率历史数据
	Memory         []MetricDataPoint `json:"memory"`         // 内存使用率历史数据
	Disk           []MetricDataPoint `json:"disk"`           // 磁盘使用率历史数据
	DiskReadKB     []MetricDataPoint `json:"diskReadKB"`     // 磁盘读取速率(KB/s)
	DiskWriteKB    []MetricDataPoint `json:"diskWriteKB"`    // 磁盘写入速率(KB/s)
	NetworkReceive []MetricDataPoint `json:"networkReceive"` // 网络接收历史数据
	NetworkSend    []MetricDataPoint `json:"networkSend"`    // 网络发送历史数据
	Load1min       []MetricDataPoint `json:"load1min"`       // 1分钟系统负载
	Load5min       []MetricDataPoint `json:"load5min"`       // 5分钟系统负载
	Load15min      []MetricDataPoint `json:"load15min"`      // 15分钟系统负载
	TotalProcesses []MetricDataPoint `json:"totalProcesses"` // 系统进程总数
}

// ProcessMetrics 进程级别监控指标
type ProcessMetrics struct {
	PID        uint              `json:"pid"`        // 进程ID
	Name       string            `json:"name"`       // 进程名称
	CPUPercent []MetricDataPoint `json:"cpuPercent"` // CPU使用百分比历史数据
	MemPercent []MetricDataPoint `json:"memPercent"` // 内存使用百分比历史数据
	Host       string            `json:"host"`       // 主机名
}

// ProcessInfo 进程基本信息
type ProcessInfo struct {
	PID        uint    `json:"pid"`        // 进程ID
	Name       string  `json:"name"`       // 进程名称
	CPUPercent float64 `json:"cpuPercent"` // CPU使用百分比
	MemPercent float64 `json:"memPercent"` // 内存使用百分比
	Host       string  `json:"host"`       // 主机名
}

// TopProcessesResult 进程TOP排行结果
type TopProcessesResult struct {
	HostID     uint          `json:"hostId"`     // 主机ID
	HostName   string        `json:"hostName"`   // 主机名称
	TopCPU     []ProcessInfo `json:"topCPU"`     // CPU使用率前5的进程
	TopMemory  []ProcessInfo `json:"topMemory"`  // 内存使用率前5的进程
	UpdateTime int64         `json:"updateTime"` // 更新时间戳
}

// PortInfo TCP端口信息
type PortInfo struct {
	Port       string  `json:"port"`       // 端口号
	PID        string  `json:"pid"`        // 进程ID
	Service    string  `json:"service"`    // 服务名称
	Status     int     `json:"status"`     // 监听状态(1:监听中,0:未监听)
	CPUUsage   float64 `json:"cpuUsage"`   // 进程CPU使用率百分比
	MemUsage   float64 `json:"memUsage"`   // 进程内存使用率百分比
}

// HostPortsResult 主机端口信息结果
type HostPortsResult struct {
	HostID     uint       `json:"hostId"`     // 主机ID
	HostName   string     `json:"hostName"`   // 主机名称
	Ports      []PortInfo `json:"ports"`      // 端口列表
	Total      int        `json:"total"`      // 端口总数
	UpdateTime int64      `json:"updateTime"` // 更新时间戳
}
