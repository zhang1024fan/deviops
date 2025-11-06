package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)


// Config 结构体定义
type Config struct {
	Processes []struct {
		Name string `yaml:"name"`
	} `yaml:"processes"`
	Endpoints []struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
	} `yaml:"endpoints"`
	Pushgateway struct {
		URL          string `yaml:"url"`
		PushInterval string `yaml:"push_interval"`
		Job          string `yaml:"job"`
		Instance     string `yaml:"instance"`
	} `yaml:"pushgateway"`
	Heartbeat struct {
		ServerURL string `yaml:"server_url"`
		Interval  string `yaml:"interval"`
		Token     string `yaml:"token"`
	} `yaml:"heartbeat"`
	CollectInterval string `yaml:"collect_interval"`
}

// 默认配置
var defaultConfig = Config{
	Processes: []struct {
		Name string `yaml:"name"`
	}{
		{Name: "nginx"},
		{Name: "mysql"},
		{Name: "redis"},
		{Name: "go"},
		{Name: "bash"},
	},
	Endpoints: []struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
	}{
		{URL: "http://localhost:9100/metrics", Method: "GET"},
	},
	Pushgateway: struct {
		URL          string `yaml:"url"`
		PushInterval string `yaml:"push_interval"`
		Job          string `yaml:"job"`
		Instance     string `yaml:"instance"`
	}{
		URL:          "", // 从配置文件或环境变量读取
		PushInterval: "30s",
		Job:          "devops-agent-%s",
		Instance:     "$HOSTNAME",
	},
	Heartbeat: struct {
		ServerURL string `yaml:"server_url"`
		Interval  string `yaml:"interval"`
		Token     string `yaml:"token"`
	}{
		ServerURL: "", // 从配置文件或环境变量读取
		Interval:  "30s",
		Token:     "agent-heartbeat-token-2024", // 默认token
	},
	CollectInterval: "15s",
}

var (
	hostname string

	// 保存上次IO计数器值
	lastDiskIO map[string]disk.IOCountersStat
	lastNetIO  psnet.IOCountersStat

	// 新增指标
	diskReadRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_kb_per_second",
			Help: "磁盘读取速率(KB/s)",
		},
		[]string{"instance", "device"},
	)

	diskWriteRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_kb_per_second",
			Help: "磁盘写入速率(KB/s)",
		},
		[]string{"instance", "device"},
	)

	totalProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_total_processes",
			Help: "系统进程总数",
		},
		[]string{"instance"},
	)

	// Pushgateway状态指标
	pushSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushgateway_success_total",
			Help: "成功推送到Pushgateway的次数",
		},
	)

	pushFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushgateway_failures_total",
			Help: "推送到Pushgateway失败的次数",
		},
	)

	// 系统指标
	cpuTotalUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "CPU总使用百分比",
		},
		[]string{"instance"},
	)

	systemLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average",
			Help: "系统负载平均值",
		},
		[]string{"instance", "period"},
	)

	memUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_percent",
			Help: "系统内存使用百分比",
		},
		[]string{"instance"},
	)

	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_usage_percent",
			Help: "磁盘使用率百分比",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	netRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_receive_kb_per_second",
			Help: "网络接收速率(KB/s)",
		},
		[]string{"instance"},
	)

	netSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_send_kb_per_second",
			Help: "网络发送速率(KB/s)",
		},
		[]string{"instance"},
	)

	// 进程CPU占用监控
	processCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_percent",
			Help: "TOP进程CPU使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// 进程内存占用监控
	processMem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_percent",
			Help: "TOP进程内存使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// 端点响应时间监控
	endpointResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_response_time_ms",
			Help: "HTTP端点响应时间(毫秒)",
		},
		[]string{"host", "url", "method"},
	)

	// 端点存活监控
	endpointStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_status",
			Help: "HTTP端点存活状态(1=正常,0=异常)",
		},
		[]string{"host", "url", "method"},
	)

	// Ping响应时间监控
	pingResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_response_time_ms",
			Help: "Ping网络连通性响应时间(毫秒)",
		},
		[]string{"host", "target"},
	)

	// Ping存活监控
	pingStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_status",
			Help: "Ping网络连通性存活状态(1=成功,0=失败)",
		},
		[]string{"host", "target"},
	)

	// TCP端口服务监控
	tcpPortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tcp_port_listening",
			Help: "TCP端口监听状态(1=监听中,0=未监听)",
		},
		[]string{"instance", "port", "service", "pid"},
	)
)

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	// 尝试连接外部地址来获取本机IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Error getting local IP: %v", err)
		return ""
	}
	defer conn.Close()
	
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// HeartbeatData 心跳数据结构
type HeartbeatData struct {
	PID   int    `json:"pid"`
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

// sendHeartbeat 发送心跳到服务端
func sendHeartbeat(serverURL, token string) {
	// 获取当前进程ID
	pid := os.Getpid()
	
	// 获取本机IP
	localIP := getLocalIP()
	if localIP == "" {
		log.Println("Warning: Could not determine local IP address")
		return
	}
	
	// 构造心跳数据
	heartbeat := HeartbeatData{
		PID:   pid,
		IP:    localIP,
		Port:  9100, // Agent默认端口
		Token: token,
	}
	
	// 序列化为JSON
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		log.Printf("Error marshaling heartbeat data: %v", err)
		return
	}
	
	// 发送HTTP POST请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending heartbeat: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		log.Printf("Heartbeat sent successfully - PID: %d, IP: %s, Port: %d, Token: %s", pid, localIP, 9100, token)
	} else {
		log.Printf("Heartbeat failed with status: %d", resp.StatusCode)
	}
}

func init() {
	hostname = getHostname()
	if hostname == "" {
		hostname = "unknown"
	}

	// 指标已在全局变量中初始化，这里不需要重复创建
}


func monitorSystemMetrics(intervalSeconds float64) {
	// 采集CPU指标
	if percents, err := cpu.Percent(0, false); err == nil { // false表示获取总CPU使用率
		if len(percents) > 0 {
			cpuTotalUsage.WithLabelValues(hostname).Set(percents[0])
		}
	}

	// 采集系统负载
	if load, err := load.Avg(); err == nil {
		systemLoad.WithLabelValues(hostname, "1min").Set(load.Load1)
		systemLoad.WithLabelValues(hostname, "5min").Set(load.Load5)
		systemLoad.WithLabelValues(hostname, "15min").Set(load.Load15)
	}

	// 采集内存指标
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memUsage.WithLabelValues(hostname).Set(memInfo.UsedPercent)
	}

	// 采集磁盘使用率(只采集根分区)
	if partitions, err := disk.Partitions(false); err == nil {
		for _, part := range partitions {
			if part.Mountpoint == "/" { // 只处理根分区
				if usage, err := disk.Usage(part.Mountpoint); err == nil {
					diskUsage.WithLabelValues(
						hostname,
						part.Device,
						part.Mountpoint,
					).Set(usage.UsedPercent)
				}
			}
		}
	}

	// 采集磁盘IO指标
	if diskIO, err := disk.IOCounters(); err == nil {
		if lastDiskIO != nil {
			for device, stats := range diskIO {
				if lastStats, ok := lastDiskIO[device]; ok {
					// 计算间隔内的速率(KB/s)
					readRate := float64(stats.ReadBytes-lastStats.ReadBytes) / 1024 / intervalSeconds
					writeRate := float64(stats.WriteBytes-lastStats.WriteBytes) / 1024 / intervalSeconds
					diskReadRate.WithLabelValues(hostname, device).Set(readRate)
					diskWriteRate.WithLabelValues(hostname, device).Set(writeRate)
				}
			}
		}
		lastDiskIO = diskIO
	}

	// 采集进程总数
	if processes, err := process.Processes(); err == nil {
		totalProcesses.WithLabelValues(hostname).Set(float64(len(processes)))
	}

	// 采集网络指标
	if netStats, err := psnet.IOCounters(true); err == nil { // true表示获取所有接口
		if len(netStats) > 0 {
			// 找到活动的物理网络接口(排除docker、虚拟网桥等)
			var mainInterface *psnet.IOCountersStat
			for i := range netStats {
				name := netStats[i].Name
				// 跳过虚拟接口和docker网桥
				if strings.HasPrefix(name, "docker") ||
					strings.HasPrefix(name, "br-") ||
					strings.HasPrefix(name, "veth") ||
					strings.HasPrefix(name, "lo") {
					continue
				}
				// 优先选择有流量的接口
				if netStats[i].BytesRecv > 0 || netStats[i].BytesSent > 0 {
					mainInterface = &netStats[i]
					break
				}
			}
			// 如果没有找到活动接口，使用第一个非虚拟接口
			if mainInterface == nil {
				for i := range netStats {
					name := netStats[i].Name
					if !strings.HasPrefix(name, "docker") &&
						!strings.HasPrefix(name, "br-") &&
						!strings.HasPrefix(name, "veth") &&
						!strings.HasPrefix(name, "lo") {
						mainInterface = &netStats[i]
						break
					}
				}
			}
			// 如果还是没有找到，使用第一个接口
			if mainInterface == nil && len(netStats) > 0 {
				mainInterface = &netStats[0]
			}

			// 记录详细统计
			log.Printf("Selected interface %s - BytesRecv: %d, BytesSent: %d",
				mainInterface.Name, mainInterface.BytesRecv, mainInterface.BytesSent)

			if lastNetIO.Name != "" {
				recvRate := float64(mainInterface.BytesRecv-lastNetIO.BytesRecv) / 1024 / intervalSeconds
				sentRate := float64(mainInterface.BytesSent-lastNetIO.BytesSent) / 1024 / intervalSeconds
				log.Printf("Network rates - Receive: %.2f KB/s, Send: %.2f KB/s", recvRate, sentRate)
				netRecv.WithLabelValues(hostname).Set(recvRate)
				netSent.WithLabelValues(hostname).Set(sentRate)
			}
			lastNetIO = *mainInterface
		} else {
			log.Println("No network interfaces found")
		}
	} else {
		log.Printf("Failed to get network stats: %v", err)
	}
}

func monitorProcesses(config *Config) {
	processes, _ := process.Processes()

	// 收集所有进程的CPU和内存使用率
	type procStats struct {
		name string
		pid  int32
		cpu  float64
		mem  float32
	}
	var allProcs []procStats

	for _, p := range processes {
		name, _ := p.Name()
		pid := p.Pid
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		allProcs = append(allProcs, procStats{
			name: name,
			pid:  pid,
			cpu:  cpuPercent,
			mem:  memPercent,
		})
	}

	// 按CPU使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].cpu > allProcs[j].cpu
	})

	// 清理旧的CPU指标
	processCPU.Reset()

	// 记录CPU使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processCPU.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(p.cpu)
	}

	// 按内存使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].mem > allProcs[j].mem
	})

	// 清理旧的内存指标
	processMem.Reset()

	// 记录内存使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processMem.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(float64(p.mem))
	}
}

func monitorEndpoints(config *Config) {
	for _, endpoint := range config.Endpoints {
		start := time.Now()
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(endpoint.URL)
		duration := time.Since(start).Milliseconds()

		status := 0
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = 1
			}
			resp.Body.Close()
		}

		endpointResponseTime.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(duration))
		endpointStatus.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(status))
		log.Printf("Endpoint %s %s - Status: %d, Response: %dms",
			endpoint.Method, endpoint.URL, status, duration)
	}
}

func monitorPing(target string) {
	start := time.Now()

	// 使用系统ping命令
	cmd := exec.Command("ping", "-c", "1", target)
	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	status := 0
	if err == nil {
		status = 1
	}

	pingResponseTime.WithLabelValues(hostname, target).Set(float64(duration))
	pingStatus.WithLabelValues(hostname, target).Set(float64(status))
	log.Printf("Ping %s - Status: %d, Response: %dms", target, status, duration)
}

// monitorTCPPorts TCP端口监控
func monitorTCPPorts() {
	// 执行netstat命令获取TCP监听端口
	cmd := exec.Command("netstat", "-ntlp")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to execute netstat: %%v", err)
		return
	}

	// 解析netstat输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[2:] { // 跳过标题行
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 7 && fields[5] == "LISTEN" {
			// 解析本地地址和端口
			localAddr := fields[3]
			if strings.Contains(localAddr, ":") {
				parts := strings.Split(localAddr, ":")
				if len(parts) >= 2 {
					port := parts[len(parts)-1]

					// 解析进程信息
					processInfo := fields[6]
					pid := ""
					service := ""

					if processInfo != "-" && strings.Contains(processInfo, "/") {
						processParts := strings.Split(processInfo, "/")
						if len(processParts) >= 2 {
							pid = processParts[0]
							service = processParts[1]
						}
					}

					// 如果无法获取服务名，使用端口号作为服务名
					if service == "" {
						service = "port-" + port
					}

					// 记录TCP端口状态
					tcpPortStatus.WithLabelValues(hostname, port, service, pid).Set(1)

					log.Printf("TCP Port Monitor: %s:%s - Service: %s, PID: %s", hostname, port, service, pid)
				}
			}
		}
	}
}

// StartAgent 启动代理服务
func StartAgent() error {
	return StartAgentWithConfig(&defaultConfig)
}

// StartAgentWithConfig 使用指定配置启动代理服务
func StartAgentWithConfig(config *Config) error {
	// 设置默认值（如果配置为空）
	if config.Pushgateway.URL == "" {
		config.Pushgateway.URL = "http://8.130.14.34:9091"
	}
	if config.Heartbeat.ServerURL == "" {
		config.Heartbeat.ServerURL = "http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat"
	}
	
	// 支持通过环境变量覆盖配置
	if url := os.Getenv("PUSHGATEWAY_URL"); url != "" {
		config.Pushgateway.URL = url
	}
	if job := os.Getenv("PUSHGATEWAY_JOB"); job != "" {
		config.Pushgateway.Job = job
	}
	if interval := os.Getenv("PUSHGATEWAY_INTERVAL"); interval != "" {
		config.Pushgateway.PushInterval = interval
	}
	if collectInterval := os.Getenv("COLLECT_INTERVAL"); collectInterval != "" {
		config.CollectInterval = collectInterval
	}
	// 心跳相关环境变量
	if heartbeatURL := os.Getenv("HEARTBEAT_SERVER_URL"); heartbeatURL != "" {
		config.Heartbeat.ServerURL = heartbeatURL
	}
	if heartbeatInterval := os.Getenv("HEARTBEAT_INTERVAL"); heartbeatInterval != "" {
		config.Heartbeat.Interval = heartbeatInterval
	}
	if heartbeatToken := os.Getenv("HEARTBEAT_TOKEN"); heartbeatToken != "" {
		config.Heartbeat.Token = heartbeatToken
	}

	// 定时采集
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in monitor goroutine: %v", r)
			}
		}()

		// 解析采集间隔
		collectInterval, err := time.ParseDuration(config.CollectInterval)
		if err != nil {
			collectInterval = 15 * time.Second // 默认值
			log.Printf("Invalid collect interval, using default 15s")
		}
		log.Printf("Collect interval set to: %v", collectInterval)

		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in monitor function: %v", r)
					}
				}()

				monitorSystemMetrics(collectInterval.Seconds())
				monitorProcesses(config)
				monitorEndpoints(config)
				monitorPing("8.8.8.8") // 默认ping Google DNS
				monitorTCPPorts() // TCP端口监控
			}()
			time.Sleep(collectInterval)
		}
	}()

	// 启动HTTP服务
	// 添加查询示例路由
	http.HandleFunc("/query-examples", func(w http.ResponseWriter, r *http.Request) {
		examples := `
Prometheus查询示例：

1. 查询所有核心的CPU使用率:
system_cpu_usage_percent{}

2. 查询特定核心(core0)的CPU使用率:
system_cpu_usage_percent{core="core0"}

3. 查询多个核心的CPU使用率:
system_cpu_usage_percent{core=~"core0|core1"}

4. 查询内存使用率大于70%的指标:
system_memory_usage_percent > 70

5. 查询特定设备的磁盘读取速度:
system_disk_read_kb_per_second{device="disk0"}

6. 查询网络接收速率:
system_network_receive_kb_per_second

7. 查询进程CPU使用率(按进程名):
process_cpu_percent{name="go"}

8. 查询端点响应时间大于10ms的:
endpoint_response_time_ms > 10

8. 查询状态异常的端点:
endpoint_status == 0

9. 5分钟内的CPU使用率变化:
rate(system_cpu_usage_percent[5m])

10. 各核心CPU使用率平均值:
avg by (core) (system_cpu_usage_percent)
`
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(examples))
	})

	// 创建自定义注册表，不包含Go运行时指标
	registry := prometheus.NewRegistry()

	// 注册所有自定义指标
	registry.MustRegister(
		cpuTotalUsage,
		systemLoad,
		memUsage,
		diskUsage,
		diskReadRate,
		diskWriteRate,
		totalProcesses,
		netRecv,
		netSent, // 确保网络发送指标被注册
		processCPU,
		processMem,
		endpointResponseTime,
		endpointStatus,
		pingResponseTime,
		pingStatus,
		tcpPortStatus,
	)

	// 注册推送状态指标
	registry.MustRegister(pushSuccess, pushFailures)

	// 使用自定义注册表
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// 启动Pushgateway推送goroutine
	if config.Pushgateway.URL != "" {
		pushInterval, err := time.ParseDuration(config.Pushgateway.PushInterval)
		if err != nil {
			pushInterval = 30 * time.Second
		}

		instance := strings.ReplaceAll(config.Pushgateway.Instance, "$HOSTNAME", hostname)

		go func() {
			for {
				pusher := push.New(config.Pushgateway.URL, config.Pushgateway.Job)
				pusher = pusher.Grouping("instance", instance)
				err := pusher.Gatherer(registry).Push()

				if err != nil {
					log.Printf("Push to Pushgateway failed: %v", err)
					pushFailures.Inc()
				} else {
					pushSuccess.Inc()
				}

				time.Sleep(pushInterval)
			}
		}()
	}

	// 启动心跳goroutine
	if config.Heartbeat.ServerURL != "" {
		heartbeatInterval, err := time.ParseDuration(config.Heartbeat.Interval)
		if err != nil {
			heartbeatInterval = 30 * time.Second // 默认30秒
			log.Printf("Invalid heartbeat interval, using default 30s")
		}
		log.Printf("Heartbeat interval set to: %v, server: %s", heartbeatInterval, config.Heartbeat.ServerURL)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in heartbeat goroutine: %v", r)
				}
			}()
			
			// 立即发送一次心跳
			sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)
			
			// 定时发送心跳
			ticker := time.NewTicker(heartbeatInterval)
			defer ticker.Stop()
			
			for range ticker.C {
				sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)
			}
		}()
	} else {
		log.Println("Heartbeat disabled: no server URL configured")
	}

	// 持续运行HTTP服务
	server := &http.Server{
		Addr:    ":9100",
		Handler: nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Agent started successfully on port 9100")
	return nil
}


// findProjectRoot 查找包含go.mod的项目根目录
func findProjectRoot() string {
	// 获取当前执行文件的路径
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "/tmp" // 如果都失败了，返回临时目录
	}

	// 从执行文件所在目录向上查找go.mod
	dir := execPath
	for i := 0; i < 10; i++ { // 最多向上查找10级目录
		dir = dir[:strings.LastIndex(dir, "/")]
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}
		if dir == "/" || dir == "" {
			break
		}
	}

	// 如果还是找不到，返回当前目录
	if pwd, err := os.Getwd(); err == nil {
		return pwd
	}
	return "/tmp"
}

// CreateAgentMainFile 创建独立可执行的agent main文件（完全独立，不依赖任何外部包）
func CreateAgentMainFile(heartbeatURL, heartbeatToken, pushgatewayURL string) string {
	log.Printf("CreateAgentMainFile: 开始生成agent代码，heartbeatURL=%s", heartbeatURL)

	// 设置默认值
	if heartbeatURL == "" {
		heartbeatURL = "http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat"
	}
	if heartbeatToken == "" {
		heartbeatToken = "agent-heartbeat-token-2024"
	}
	if pushgatewayURL == "" {
		pushgatewayURL = "http://8.130.14.34:9091"
	}

	// 直接构建完整的agent代码
	mainCode := fmt.Sprintf(`package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Config 结构体定义
type Config struct {
	Processes []struct {
		Name string
	}
	Endpoints []struct {
		URL    string
		Method string
	}
	Pushgateway struct {
		URL          string
		PushInterval string
		Job          string
		Instance     string
	}
	Heartbeat struct {
		ServerURL string
		Interval  string
		Token     string
	}
	CollectInterval string
}

// 声明全局变量
var (
	// CPU使用率监控
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "CPU使用率百分比",
		},
		[]string{"instance", "cpu"},
	)

	// 内存使用率监控
	memUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_percent",
			Help: "内存使用率百分比",
		},
		[]string{"instance"},
	)

	// 总内存容量监控
	memTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_total_bytes",
			Help: "系统总内存容量(字节)",
		},
		[]string{"instance"},
	)

	// 已用内存监控
	memUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_used_bytes",
			Help: "已使用内存容量(字节)",
		},
		[]string{"instance"},
	)

	// 可用内存监控
	memFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_free_bytes",
			Help: "可用内存容量(字节)",
		},
		[]string{"instance"},
	)

	// 磁盘使用率监控
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_usage_percent",
			Help: "磁盘使用率百分比",
		},
		[]string{"instance", "mountpoint"},
	)

	// 系统负载监控
	systemLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average",
			Help: "系统负载平均值(1分钟/5分钟/15分钟)",
		},
		[]string{"instance", "period"},
	)

	// 磁盘读取速度监控
	diskRead = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_kb_per_second",
			Help: "磁盘读取速度(KB/秒)",
		},
		[]string{"instance", "device"},
	)

	// 磁盘写入速度监控
	diskWrite = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_kb_per_second",
			Help: "磁盘写入速度(KB/秒)",
		},
		[]string{"instance", "device"},
	)

	// 进程总数监控
	totalProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_total_processes",
			Help: "系统进程总数",
		},
		[]string{"instance"},
	)

	// 网络接收流量监控
	netRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_receive_kb_per_second",
			Help: "网络接收流量速度(KB/秒)",
		},
		[]string{"instance"},
	)

	// 网络发送流量监控
	netSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_send_kb_per_second",
			Help: "网络发送流量速度(KB/秒)",
		},
		[]string{"instance"},
	)

	// 进程CPU占用监控
	processCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_percent",
			Help: "TOP进程CPU使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// 进程内存占用监控
	processMem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_percent",
			Help: "TOP进程内存使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// TCP端口服务监控
	tcpPortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tcp_port_listening",
			Help: "TCP端口监听状态(1=监听中,0=未监听)",
		},
		[]string{"instance", "port", "service", "pid"},
	)

	// 进程CPU使用率监控(按PID)
	processCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage_percent",
			Help: "进程CPU使用率百分比(按PID)",
		},
		[]string{"instance", "pid", "name"},
	)

	// 进程内存使用率监控(按PID)
	processMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_usage_percent",
			Help: "进程内存使用率百分比(按PID)",
		},
		[]string{"instance", "pid", "name"},
	)

	// 端点响应时间监控
	endpointResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_response_time_ms",
			Help: "HTTP端点响应时间(毫秒)",
		},
		[]string{"host", "url", "method"},
	)

	// 端点存活监控
	endpointStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_status",
			Help: "HTTP端点存活状态(1=正常,0=异常)",
		},
		[]string{"host", "url", "method"},
	)

	// Ping响应时间监控
	pingResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_response_time_ms",
			Help: "Ping网络连通性响应时间(毫秒)",
		},
		[]string{"host", "target"},
	)

	// Ping存活监控
	pingStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_status",
			Help: "Ping网络连通性存活状态(1=成功,0=失败)",
		},
		[]string{"host", "target"},
	)

	// 上一次读取的磁盘IO数据，用于计算速率
	lastDiskIO disk.IOCountersStat
	lastNetIO  psnet.IOCountersStat
)

// HeartbeatData 心跳数据结构
type HeartbeatData struct {
	PID      int    `+"`json:\"pid\"`"+`
	IP       string `+"`json:\"ip\"`"+`
	Hostname string `+"`json:\"hostname\"`"+`
	Port     int    `+"`json:\"port\"`"+`
	Token    string `+"`json:\"token\"`"+`
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Error getting local IP: %%v", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func sendHeartbeat(serverURL, token string) {
	pid := os.Getpid()
	localIP := getLocalIP()
	if localIP == "" {
		log.Println("Warning: Could not determine local IP address")
		return
	}

	// 获取hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	heartbeat := HeartbeatData{
		PID:      pid,
		IP:       localIP,
		Hostname: hostname,
		Port:     9100,
		Token:    token,
	}

	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		log.Printf("Error marshaling heartbeat data: %%v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending heartbeat: %%v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("Heartbeat sent successfully - PID: %%d, IP: %%s, Port: %%d, Token: %%s", pid, localIP, 9100, token)
	} else {
		log.Printf("Heartbeat failed with status: %%d", resp.StatusCode)
	}
}

func collectSystemMetrics() {
	// 使用hostname作为instance标识（用于Prometheus查询区分主机）
	instanceID, err := os.Hostname()
	if err != nil || instanceID == "" {
		// 如果获取hostname失败，使用IP作为备选
		instanceID = getLocalIP()
		if instanceID == "" {
			instanceID = "unknown"
		}
	}

	// 采集CPU指标
	if cpuPercents, err := cpu.Percent(0, true); err == nil {
		for i, cpuPercent := range cpuPercents {
			cpuUsage.WithLabelValues(instanceID, fmt.Sprintf("cpu%%d", i)).Set(cpuPercent)
		}
	}

	// 采集内存指标
	if memStats, err := mem.VirtualMemory(); err == nil {
		memUsage.WithLabelValues(instanceID).Set(memStats.UsedPercent)
		memTotal.WithLabelValues(instanceID).Set(float64(memStats.Total))
		memUsed.WithLabelValues(instanceID).Set(float64(memStats.Used))
		memFree.WithLabelValues(instanceID).Set(float64(memStats.Available))
	}

	// 采集磁盘指标
	if diskStats, err := disk.Partitions(false); err == nil {
		for _, partition := range diskStats {
			if usage, err := disk.Usage(partition.Mountpoint); err == nil {
				diskUsage.WithLabelValues(instanceID, partition.Mountpoint).Set(usage.UsedPercent)
			}
		}
	}

	// 采集负载平均值
	if loadStats, err := load.Avg(); err == nil {
		systemLoad.WithLabelValues(instanceID, "1min").Set(loadStats.Load1)
		systemLoad.WithLabelValues(instanceID, "5min").Set(loadStats.Load5)
		systemLoad.WithLabelValues(instanceID, "15min").Set(loadStats.Load15)
	}

	// 采集磁盘IO指标
	const intervalSeconds = 10.0
	if diskIOStats, err := disk.IOCounters(); err == nil {
		var totalReadBytes, totalWriteBytes uint64
		for _, diskIO := range diskIOStats {
			totalReadBytes += diskIO.ReadBytes
			totalWriteBytes += diskIO.WriteBytes
		}
		diskIO := disk.IOCountersStat{
			ReadBytes:  totalReadBytes,
			WriteBytes: totalWriteBytes,
		}

		if lastDiskIO.ReadBytes != 0 {
			readRate := float64(diskIO.ReadBytes-lastDiskIO.ReadBytes) / 1024 / intervalSeconds
			writeRate := float64(diskIO.WriteBytes-lastDiskIO.WriteBytes) / 1024 / intervalSeconds
			diskRead.WithLabelValues(instanceID, "vda").Set(readRate)
			diskWrite.WithLabelValues(instanceID, "vda").Set(writeRate)
		}
		lastDiskIO = diskIO
	}

	// 采集进程总数
	if processes, err := process.Processes(); err == nil {
		totalProcesses.WithLabelValues(instanceID).Set(float64(len(processes)))
	}

	// 采集网络指标
	if netStats, err := psnet.IOCounters(true); err == nil {
		if len(netStats) > 0 {
			var mainInterface *psnet.IOCountersStat
			for i := range netStats {
				name := netStats[i].Name
				if strings.HasPrefix(name, "docker") ||
					strings.HasPrefix(name, "br-") ||
					strings.HasPrefix(name, "veth") ||
					strings.HasPrefix(name, "lo") {
					continue
				}
				if netStats[i].BytesRecv > 0 || netStats[i].BytesSent > 0 {
					mainInterface = &netStats[i]
					break
				}
			}
			if mainInterface == nil {
				for i := range netStats {
					name := netStats[i].Name
					if !strings.HasPrefix(name, "docker") &&
						!strings.HasPrefix(name, "br-") &&
						!strings.HasPrefix(name, "veth") &&
						!strings.HasPrefix(name, "lo") {
						mainInterface = &netStats[i]
						break
					}
				}
			}
			if mainInterface == nil && len(netStats) > 0 {
				mainInterface = &netStats[0]
			}

			if lastNetIO.Name != "" {
				const intervalSeconds = 10.0
				recvRate := float64(mainInterface.BytesRecv-lastNetIO.BytesRecv) / 1024 / intervalSeconds
				sentRate := float64(mainInterface.BytesSent-lastNetIO.BytesSent) / 1024 / intervalSeconds
				netRecv.WithLabelValues(instanceID).Set(recvRate)
				netSent.WithLabelValues(instanceID).Set(sentRate)
			}
			lastNetIO = *mainInterface
		}
	}
}

func monitorProcesses(config *Config) {
	// 动态获取主机名，而不是使用编译时写死的名字
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	processes, _ := process.Processes()

	type procStats struct {
		name string
		pid  int32
		cpu  float64
		mem  float32
	}
	var allProcs []procStats

	for _, p := range processes {
		name, _ := p.Name()
		pid := p.Pid
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		allProcs = append(allProcs, procStats{
			name: name,
			pid:  pid,
			cpu:  cpuPercent,
			mem:  memPercent,
		})
	}

	// 按CPU使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].cpu > allProcs[j].cpu
	})

	processCPU.Reset()

	// 记录CPU使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processCPU.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(p.cpu)
	}

	// 按内存使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].mem > allProcs[j].mem
	})

	processMem.Reset()

	// 记录内存使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processMem.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(float64(p.mem))
	}
}

func monitorEndpoints(config *Config) {
	// 动态获取主机名，而不是使用编译时写死的名字
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	for _, endpoint := range config.Endpoints {
		start := time.Now()
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(endpoint.URL)
		duration := time.Since(start).Milliseconds()

		status := 0
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = 1
			}
			resp.Body.Close()
		}

		endpointResponseTime.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(duration))
		endpointStatus.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(status))
	}
}

func monitorPing(target string) {
	// 动态获取主机名，而不是使用编译时写死的名字
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	start := time.Now()

	cmd := exec.Command("ping", "-c", "1", target)
	err = cmd.Run()
	duration := time.Since(start).Milliseconds()

	status := 0
	if err == nil {
		status = 1
	}

	pingResponseTime.WithLabelValues(hostname, target).Set(float64(duration))
	pingStatus.WithLabelValues(hostname, target).Set(float64(status))
}

func monitorTCPPorts() {
	// 动态获取主机名，而不是使用编译时写死的名字
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	cmd := exec.Command("netstat", "-ntlp")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to execute netstat: %%v", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[2:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 7 && fields[5] == "LISTEN" {
			localAddr := fields[3]
			if strings.Contains(localAddr, ":") {
				parts := strings.Split(localAddr, ":")
				if len(parts) >= 2 {
					port := parts[len(parts)-1]

					processInfo := fields[6]
					pid := ""
					service := ""

					if processInfo != "-" && strings.Contains(processInfo, "/") {
						processParts := strings.Split(processInfo, "/")
						if len(processParts) >= 2 {
							pid = processParts[0]
							service = processParts[1]
						}
					}

					if service == "" {
						service = "port-" + port
					}

					tcpPortStatus.WithLabelValues(hostname, port, service, pid).Set(1)

					// 收集该进程的CPU和内存使用率
					if pid != "" {
						if pidInt, err := strconv.Atoi(pid); err == nil {
							if p, err := process.NewProcess(int32(pidInt)); err == nil {
								if cpuPercent, err := p.CPUPercent(); err == nil {
									processCPUUsage.WithLabelValues(hostname, pid, service).Set(cpuPercent)
								}
								if memPercent, err := p.MemoryPercent(); err == nil {
									processMemoryUsage.WithLabelValues(hostname, pid, service).Set(float64(memPercent))
								}
							}
						}
					}
				}
			}
		}
	}
}

func registerMetrics(registry *prometheus.Registry) {
	registry.MustRegister(cpuUsage)
	registry.MustRegister(memUsage)
	registry.MustRegister(memTotal)
	registry.MustRegister(memUsed)
	registry.MustRegister(memFree)
	registry.MustRegister(diskUsage)
	registry.MustRegister(systemLoad)
	registry.MustRegister(diskRead)
	registry.MustRegister(diskWrite)
	registry.MustRegister(totalProcesses)
	registry.MustRegister(netRecv)
	registry.MustRegister(netSent)
	registry.MustRegister(processCPU)
	registry.MustRegister(processMem)
	registry.MustRegister(tcpPortStatus)
	registry.MustRegister(processCPUUsage)
	registry.MustRegister(processMemoryUsage)
	registry.MustRegister(endpointResponseTime)
	registry.MustRegister(endpointStatus)
	registry.MustRegister(pingResponseTime)
	registry.MustRegister(pingStatus)
}

func getEnvVar(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func StartAgentWithConfig(config *Config) error {
	// 创建自定义注册表，不包含Go运行时指标
	registry := prometheus.NewRegistry()
	registerMetrics(registry)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	go func() {
		log.Println("Agent starting on :9100")
		if err := http.ListenAndServe(":9100", nil); err != nil {
			log.Printf("HTTP server failed: %%v", err)
		}
	}()

	// 定期采集系统指标
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// 定期推送到Pushgateway
	pushTicker := time.NewTicker(30 * time.Second)
	defer pushTicker.Stop()

	// 心跳
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ticker.C:
			collectSystemMetrics()
			monitorProcesses(config)
			monitorEndpoints(config)
			monitorPing("8.8.8.8")
			monitorTCPPorts()
		case <-pushTicker.C:
			if config.Pushgateway.URL != "" {
				pusher := push.New(config.Pushgateway.URL, config.Pushgateway.Job)
				// 不使用.Grouping()因为指标中已经包含instance标签

				// 使用registry推送所有指标（不包含Go运行时指标）
				if err := pusher.Gatherer(registry).Push(); err != nil {
					log.Printf("Could not push to gateway: %%v", err)
				} else {
					log.Printf("Successfully pushed metrics to gateway: %%s", config.Pushgateway.URL)
				}
			}
		case <-heartbeatTicker.C:
			if config.Heartbeat.ServerURL != "" {
				sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)
			}
		}
	}
}

func main() {
	log.Println("Starting dodevops-agent...")

	// 动态获取hostname作为Job名称，确保每台主机在Pushgateway中有独立的job
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
		log.Printf("Warning: could not get hostname, using default: %%v", err)
	}
	log.Printf("Agent hostname: %%s", hostname)

	config := &Config{
		Processes: []struct {
			Name string
		}{
			{Name: "nginx"},
			{Name: "mysql"},
			{Name: "redis"},
		},
		Endpoints: []struct {
			URL    string
			Method string
		}{
			{URL: "http://localhost:9100/metrics", Method: "GET"},
		},
		Pushgateway: struct {
			URL          string
			PushInterval string
			Job          string
			Instance     string
		}{
			URL:          getEnvVar("PUSHGATEWAY_URL", "%s"),
			PushInterval: "30s",
			Job:          hostname,
			Instance:     hostname,
		},
		Heartbeat: struct {
			ServerURL string
			Interval  string
			Token     string
		}{
			ServerURL: getEnvVar("HEARTBEAT_SERVER_URL", "%s"),
			Interval:  "30s",
			Token:     getEnvVar("HEARTBEAT_TOKEN", "%s"),
		},
		CollectInterval: "15s",
	}

	if err := StartAgentWithConfig(config); err != nil {
		log.Fatalf("Failed to start agent: %%v", err)
	}
	select {}
}
`, pushgatewayURL, heartbeatURL, heartbeatToken)

	log.Printf("CreateAgentMainFile: 生成的代码长度: %d", len(mainCode))
	log.Printf("CreateAgentMainFile: 代码生成完成，heartbeat URL=%s", heartbeatURL)

	return mainCode
}

// ValidateAgentCode 验证生成的agent代码是否有效
func ValidateAgentCode(code string) error {
	// 检查基本结构
	if !strings.Contains(code, "package main") {
		return fmt.Errorf("缺少package main声明")
	}

	if !strings.Contains(code, "func main()") {
		return fmt.Errorf("缺少main函数")
	}

	// 简化验证：只检查明显的语法错误，不做复杂的解析
	lines := strings.Split(code, "\n")

	// 检查是否有明显的语法问题
	openBraces := 0
	openParens := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		openBraces += strings.Count(trimmed, "{")
		openBraces -= strings.Count(trimmed, "}")
		openParens += strings.Count(trimmed, "(")
		openParens -= strings.Count(trimmed, ")")
	}

	// 检查括号匹配
	if openBraces != 0 {
		return fmt.Errorf("大括号不匹配")
	}
	if openParens != 0 {
		return fmt.Errorf("圆括号不匹配")
	}

	return nil
}
