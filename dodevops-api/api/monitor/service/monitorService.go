package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	monitormodel "dodevops-api/api/monitor/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

type hostResult struct {
	id      uint
	metrics *monitormodel.HostMetrics
}

// MonitorServiceInterface 定义监控服务接口
type MonitorServiceInterface interface {
	GetHostMetrics(c *gin.Context, hostID uint)                      // 获取主机监控指标
	BatchGetHostMetrics(c *gin.Context, hostIDs []uint)              // 批量获取主机监控指标
	GetHostMetricHistory(c *gin.Context, hostID uint, metric string, // 获取主机指标历史数据
		start, end string, step string)
	GetHostAllMetricsHistory(c *gin.Context, hostID uint, // 获取主机所有指标历史数据
		start, end string, step string)
	GetTopProcesses(c *gin.Context, hostID uint) // 获取主机TOP进程使用率
	GetHostPorts(c *gin.Context, hostID uint)    // 获取主机端口信息
}

// MonitorServiceImpl 监控服务实现
type MonitorServiceImpl struct {
	hostDao       dao.CmdbHostDao
	prometheusURL string
	cacheExpire   time.Duration
}

// 缓存键名常量
const (
	hostCacheKey = "monitor:host:%d" // 主机信息缓存键
	cacheTimeout = 5 * time.Minute   // 缓存过期时间
)

// NewMonitorService 创建监控服务实例
func NewMonitorService() MonitorServiceInterface {
	return &MonitorServiceImpl{
		hostDao:       dao.NewCmdbHostDao(),
		prometheusURL: config.Config.Monitor.Prometheus.URL,
		cacheExpire:   cacheTimeout,
	}
}

// checkHostOnline 检查主机是否在线(通过TCP端口检测)
func (s *MonitorServiceImpl) checkHostOnline(host model.CmdbHost) bool {
	if host.SSHIP == "" || host.SSHPort == 0 {
		return false
	}

	// 使用更快的TCP探测，超时500ms
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort),
		500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// getHostFromCache 从缓存获取主机信息
func (s *MonitorServiceImpl) getHostFromCache(hostID uint) (*model.CmdbHost, error) {
	cacheKey := fmt.Sprintf(hostCacheKey, hostID)

	// 从 Redis 获取缓存数据
	cachedData, err := redis.RedisDb.Get(context.Background(), cacheKey).Result()
	if err != nil {
		// 缓存未命中或出错，返回 nil
		return nil, err
	}

	// 反序列化主机信息
	var host model.CmdbHost
	if err := json.Unmarshal([]byte(cachedData), &host); err != nil {
		return nil, err
	}

	return &host, nil
}

// setHostToCache 将主机信息存入缓存
func (s *MonitorServiceImpl) setHostToCache(hostID uint, host *model.CmdbHost) error {
	cacheKey := fmt.Sprintf(hostCacheKey, hostID)

	// 序列化主机信息
	hostData, err := json.Marshal(host)
	if err != nil {
		return err
	}

	// 存入 Redis，设置过期时间
	return redis.RedisDb.Set(context.Background(), cacheKey, hostData, s.cacheExpire).Err()
}

// getHostWithCache 优先从缓存获取主机信息，缓存未命中则查询数据库
func (s *MonitorServiceImpl) getHostWithCache(hostID uint) (*model.CmdbHost, error) {
	// 1. 先尝试从缓存获取
	if host, err := s.getHostFromCache(hostID); err == nil && host != nil {
		return host, nil
	}

	// 2. 缓存未命中，查询数据库
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		return nil, err
	}

	// 3. 将查询结果存入缓存(异步执行，不影响主流程)
	go func() {
		if err := s.setHostToCache(hostID, &host); err != nil {
			fmt.Printf("缓存主机信息失败: %v\n", err)
		}
	}()

	return &host, nil
}

// GetHostMetrics 获取单个主机监控指标
func (s *MonitorServiceImpl) GetHostMetrics(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	if host.Name == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机名称为空")
		return
	}

	metrics := &monitormodel.HostMetrics{}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		metrics.CPUUsage = 0
		metrics.MemoryUsage = 0
		metrics.DiskUsage = 0
		result.Success(c, metrics)
		return
	}

	// ... existing code ...
	now := time.Now()
	startTime := now.Add(-5 * time.Minute) // 查询最近5分钟的数据
	endTime := now
	step := "30s" // 30秒间隔

	// 查询CPU使用率
	cpuQuery := fmt.Sprintf("system_cpu_usage_percent{instance=\"%s\"}", host.Name)
	cpu, err := s.queryPrometheus(cpuQuery, startTime, endTime, step)
	if err == nil && len(cpu.Data.Result) > 0 && len(cpu.Data.Result[0].Values) > 0 {
		lastPoint := cpu.Data.Result[0].Values[len(cpu.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.CPUUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.CPUUsage = 0
	}

	// 查询内存使用率
	memQuery := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", host.Name)
	mem, err := s.queryPrometheus(memQuery, startTime, endTime, step)
	if err == nil && len(mem.Data.Result) > 0 && len(mem.Data.Result[0].Values) > 0 {
		lastPoint := mem.Data.Result[0].Values[len(mem.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.MemoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.MemoryUsage = 0
	}

	// 查询磁盘使用率
	diskQuery := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", host.Name)
	disk, err := s.queryPrometheus(diskQuery, startTime, endTime, step)
	if err == nil && len(disk.Data.Result) > 0 && len(disk.Data.Result[0].Values) > 0 {
		lastPoint := disk.Data.Result[0].Values[len(disk.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.DiskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.DiskUsage = 0
	}

	result.Success(c, metrics)
}

// BatchGetHostMetrics 批量获取主机监控指标
func (s *MonitorServiceImpl) BatchGetHostMetrics(c *gin.Context, hostIDs []uint) {
	// 参数校验
	if len(hostIDs) == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 限制并发数量，防止过多并发请求导致数据库压力
	maxConcurrent := 10
	if len(hostIDs) < maxConcurrent {
		maxConcurrent = len(hostIDs)
	}
	semaphore := make(chan struct{}, maxConcurrent)

	// 使用sync.WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup
	resultMap := make(map[uint]*monitormodel.HostMetrics, len(hostIDs))
	var mu sync.Mutex // 用于保护resultMap的并发访问

	// 为每个主机ID启动一个goroutine
	for _, id := range hostIDs {
		wg.Add(1)
		go func(hostID uint) {
			defer wg.Done()

			// 获取信号量，限制并发数
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 优先从缓存获取主机信息
			host, err := s.getHostWithCache(hostID)
			if err != nil {
				// 单个主机失败不影响其他主机
				return
			}

			metrics := &monitormodel.HostMetrics{}

			// 检测主机在线状态
			metrics.OnlineStatus = 1 // 默认离线
			if host.SSHIP != "" && host.SSHPort != 0 {
				conn, err := net.DialTimeout("tcp",
					fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort),
					500*time.Millisecond)
				if err == nil {
					metrics.OnlineStatus = 0 // 在线
					conn.Close()
				}
			}

			// 如果主机离线，设置监控数据为0
			if metrics.OnlineStatus == 1 {
				metrics.CPUUsage = 0
				metrics.MemoryUsage = 0
				metrics.DiskUsage = 0
				mu.Lock()
				resultMap[hostID] = metrics
				mu.Unlock()
				return
			}

			// 主机在线，获取监控数据（优化查询频率）
			// 获取当前时间
			now := time.Now()
			startTime := now.Add(-3 * time.Minute) // 缩短查询时间范围从5分钟到3分钟
			endTime := now
			step := "60s" // 增大步长从30秒到60秒

			// 优化查询：只获取最后一个数据点
			if host.Name != "" {
				// 并发查询三个指标
				var metricWg sync.WaitGroup
				var metricMu sync.Mutex

				// 查询CPU使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					cpuQuery := fmt.Sprintf("system_cpu_usage_percent{instance=\"%s\"}", host.Name)
					cpu, err := s.queryPrometheus(cpuQuery, startTime, endTime, step)
					if err == nil && len(cpu.Data.Result) > 0 && len(cpu.Data.Result[0].Values) > 0 {
						lastPoint := cpu.Data.Result[0].Values[len(cpu.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.CPUUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 查询内存使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					memQuery := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", host.Name)
					mem, err := s.queryPrometheus(memQuery, startTime, endTime, step)
					if err == nil && len(mem.Data.Result) > 0 && len(mem.Data.Result[0].Values) > 0 {
						lastPoint := mem.Data.Result[0].Values[len(mem.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.MemoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 查询磁盘使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					diskQuery := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", host.Name)
					disk, err := s.queryPrometheus(diskQuery, startTime, endTime, step)
					if err == nil && len(disk.Data.Result) > 0 && len(disk.Data.Result[0].Values) > 0 {
						lastPoint := disk.Data.Result[0].Values[len(disk.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.DiskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 等待所有指标查询完成
				metricWg.Wait()
			}

			// 安全地将结果写入map
			mu.Lock()
			resultMap[hostID] = metrics
			mu.Unlock()
		}(id)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 返回结果，即使部分主机查询失败也返回成功查询的结果
	result.Success(c, resultMap)
}

// queryPrometheus 查询Prometheus指标数据
func (s *MonitorServiceImpl) queryPrometheus(query string, start, end time.Time, step string) (*monitormodel.PrometheusQueryResult, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	// 检查Prometheus URL配置
	if s.prometheusURL == "" {
		return nil, fmt.Errorf("prometheus URL not configured")
	}

	var req *http.Request
	var q url.Values
	var err error

	// 对于时间范围查询使用query_range端点
	if !start.IsZero() && !end.IsZero() {
		req, err = http.NewRequest("GET", strings.TrimSuffix(s.prometheusURL, "/")+"/api/v1/query_range", nil)
		if err != nil {
			return nil, err
		}

		q = req.URL.Query()
		q.Add("query", query)
		q.Add("start", strconv.FormatInt(start.Unix(), 10))
		q.Add("end", strconv.FormatInt(end.Unix(), 10))
		q.Add("step", step)
	} else {
		// 即时查询使用query端点
		req, err = http.NewRequest("GET", s.prometheusURL+"/api/v1/query", nil)
		if err != nil {
			return nil, err
		}
		q = req.URL.Query()
		q.Set("query", query)
	}

	// 强制不使用缓存并获取最新数据
	q.Add("nocache", "1")
	q.Add("time", fmt.Sprintf("%d", time.Now().Unix()))  // 添加时间戳防止缓存
	q.Add("_", fmt.Sprintf("%d", time.Now().UnixNano())) // 添加纳秒级时间戳确保唯一性
	req.URL.RawQuery = q.Encode()

	// 设置请求头确保不使用缓存
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")
	req.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano())) // 添加唯一请求ID

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result monitormodel.PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Prometheus response: %v\nRaw response: %s", err, string(body))
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", string(body))
	}

	// 兼容vector和matrix类型结果
	if result.Data.ResultType != "matrix" && result.Data.ResultType != "vector" {
		return nil, fmt.Errorf("unexpected result type: %s", result.Data.ResultType)
	}

	return &result, nil
}

// GetTopProcesses 获取主机TOP进程使用率(前5名)
func (s *MonitorServiceImpl) GetTopProcesses(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	if host.Name == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机名称为空")
		return
	}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机当前不在线")
		return
	}

	// ... existing code ...
	var wg sync.WaitGroup
	var topCPU, topMemory []monitormodel.ProcessInfo
	var mu sync.Mutex

	// 查询CPU使用率TOP5进程
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuQuery := fmt.Sprintf("topk(5, process_cpu_percent{host=\"%s\"})", host.Name)
		cpuResult, err := s.queryPrometheus(cpuQuery, time.Time{}, time.Time{}, "")
		if err == nil && len(cpuResult.Data.Result) > 0 {
			mu.Lock()
			for _, item := range cpuResult.Data.Result {
				if len(item.Value) >= 2 {
					// 解析指标数据
					var processInfo monitormodel.ProcessInfo
					processInfo.Host = item.Metric["host"]
					processInfo.Name = item.Metric["name"]

					// 解析PID
					if pidStr, ok := item.Metric["pid"]; ok {
						if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
							processInfo.PID = uint(pid)
						}
					}

					// 解析CPU使用率
					if valStr, ok := item.Value[1].(string); ok {
						if cpuVal, err := strconv.ParseFloat(valStr, 64); err == nil {
							processInfo.CPUPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuVal), 64)
						}
					}

					topCPU = append(topCPU, processInfo)
				}
			}
			mu.Unlock()
		}
	}()

	// 查询内存使用率TOP5进程
	wg.Add(1)
	go func() {
		defer wg.Done()
		memQuery := fmt.Sprintf("topk(5, process_memory_percent{host=\"%s\"})", host.Name)
		memResult, err := s.queryPrometheus(memQuery, time.Time{}, time.Time{}, "")
		if err == nil && len(memResult.Data.Result) > 0 {
			mu.Lock()
			for _, item := range memResult.Data.Result {
				if len(item.Value) >= 2 {
					// 解析指标数据
					var processInfo monitormodel.ProcessInfo
					processInfo.Host = item.Metric["host"]
					processInfo.Name = item.Metric["name"]

					// 解析PID
					if pidStr, ok := item.Metric["pid"]; ok {
						if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
							processInfo.PID = uint(pid)
						}
					}

					// 解析内存使用率
					if valStr, ok := item.Value[1].(string); ok {
						if memVal, err := strconv.ParseFloat(valStr, 64); err == nil {
							processInfo.MemPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", memVal), 64)
						}
					}

					topMemory = append(topMemory, processInfo)
				}
			}
			mu.Unlock()
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	// 构建返回结果
	topResult := &monitormodel.TopProcessesResult{
		HostID:     hostID,
		HostName:   host.Name,
		TopCPU:     topCPU,
		TopMemory:  topMemory,
		UpdateTime: time.Now().Unix(),
	}

	result.Success(c, topResult)
}

// GetHostAllMetricsHistory 获取主机所有指标历史数据
func (s *MonitorServiceImpl) GetHostAllMetricsHistory(c *gin.Context, hostID uint, startTime, endTime string, step string) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 解析时间参数
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "开始时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}
	end, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "结束时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}

	// 自动计算step参数
	if step == "" {
		step = calculateStep(end.Sub(start))
	} else {
		step = strings.Trim(step, `"`) // 移除可能存在的额外引号
	}

	// 验证step格式
	if _, err := time.ParseDuration(step); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError),
			"无效的步长格式，请使用如: 15s, 1m, 5m")
		return
	}

	// 使用WaitGroup并行查询所有指标
	var wg sync.WaitGroup
	var mu sync.Mutex
	allMetrics := &monitormodel.AllMetricsHistory{
		HostID: hostID,
	}

	// 查询CPU指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("rate(system_cpu_usage_percent{instance=\"%s\"}[1m])", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.CPU = append(allMetrics.CPU, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询内存指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Memory = append(allMetrics.Memory, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Disk = append(allMetrics.Disk, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询网络接收指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_network_receive_kb_per_second{instance=\"%s\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.NetworkReceive = append(allMetrics.NetworkReceive, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询网络发送指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_network_send_kb_per_second{instance=\"%s\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.NetworkSend = append(allMetrics.NetworkSend, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询1分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average{instance=\"%s\",period=\"1min\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load1min = append(allMetrics.Load1min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询5分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average{instance=\"%s\",period=\"5min\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load5min = append(allMetrics.Load5min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询15分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average{instance=\"%s\",period=\"15min\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load15min = append(allMetrics.Load15min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘读取速率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_read_kb_per_second{instance=\"%s\",device=\"vda\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.DiskReadKB = append(allMetrics.DiskReadKB, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘写入速率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_write_kb_per_second{instance=\"%s\",device=\"vda\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.DiskWriteKB = append(allMetrics.DiskWriteKB, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询系统进程总数
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_total_processes{instance=\"%s\"}", host.Name)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.TotalProcesses = append(allMetrics.TotalProcesses, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	result.Success(c, allMetrics)
}

// calculateStep 根据时间范围计算步长，固定返回60个数据点
func calculateStep(duration time.Duration) string {
	step := duration / 60
	return fmt.Sprintf("%.0fs", step.Seconds())
}

// GetHostMetricHistory 获取主机指标历史数据
func (s *MonitorServiceImpl) GetHostMetricHistory(c *gin.Context, hostID uint, metric string, startTime, endTime string, step string) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 验证指标名称
	switch metric {
	case "cpu", "memory", "disk":
		// 有效指标
	default:
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "不支持的指标类型")
		return
	}

	// 解析时间参数
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "开始时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}
	end, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "结束时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}

	// 自动计算step参数
	if step == "" {
		step = calculateStep(end.Sub(start))
	} else {
		step = strings.Trim(step, `"`) // 移除可能存在的额外引号
	}

	// 验证step格式
	if _, err := time.ParseDuration(step); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError),
			"无效的步长格式，请使用如: 15s, 1m, 5m")
		return
	}

	// 构建Prometheus查询
	var query string
	switch metric {
	case "cpu":
		query = fmt.Sprintf("rate(system_cpu_usage_percent{instance=\"%s\"}[1m])", host.Name)
	case "memory":
		query = fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", host.Name)
	case "disk":
		query = fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", host.Name)
	}

	// 执行Prometheus查询
	promResult, err := s.queryPrometheus(query, start, end, step)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED),
			fmt.Sprintf("Prometheus查询失败: %v", err))
		return
	}

	history := &monitormodel.HostMetricHistory{
		HostID: hostID,
		Metric: metric,
	}

	// 提取时间序列数据
	if len(promResult.Data.Result) > 0 {
		for _, resultItem := range promResult.Data.Result {
			switch promResult.Data.ResultType {
			case "matrix":
				// 处理matrix类型结果(时间范围查询)
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								history.TimeData = append(history.TimeData, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			case "vector":
				// 处理vector类型结果(即时查询)
				if len(resultItem.Value) >= 2 {
					if timestamp, ok := resultItem.Value[0].(float64); ok {
						if valStr, ok := resultItem.Value[1].(string); ok {
							val, _ := strconv.ParseFloat(valStr, 64)
							history.TimeData = append(history.TimeData, monitormodel.MetricDataPoint{
								Timestamp: int64(timestamp),
								Value:     val,
							})
						}
					}
				}
			default:
				// 未知结果类型，不做处理
			}
		}
	}

	result.Success(c, history)
}

// GetHostPorts 获取主机端口信息
func (s *MonitorServiceImpl) GetHostPorts(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 首先检查Redis缓存
	cacheKey := fmt.Sprintf("monitor:ports:%d", hostID)
	if cachedData, err := redis.RedisDb.Get(context.Background(), cacheKey).Result(); err == nil {
		var cachedResult monitormodel.HostPortsResult
		if err := json.Unmarshal([]byte(cachedData), &cachedResult); err == nil {
			// 检查缓存是否过期（5分钟内的数据认为有效）
			if time.Now().Unix()-cachedResult.UpdateTime < 300 {
				result.Success(c, &cachedResult)
				return
			}
		}
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	if host.Name == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机名称为空")
		return
	}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机当前不在线")
		return
	}

	// 查询TCP端口监听状态
	query := fmt.Sprintf("tcp_port_listening{instance=\"%s\"}", host.Name)
	now := time.Now()
	startTime := now.Add(-5 * time.Minute) // 查询最近5分钟的数据
	endTime := now
	step := "30s"

	promResult, err := s.queryPrometheus(query, startTime, endTime, step)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED),
			fmt.Sprintf("Prometheus查询失败: %v", err))
		return
	}

	// 批量查询所有进程的CPU和内存使用率（一次查询获取所有数据）
	cpuQuery := fmt.Sprintf("process_cpu_usage_percent{instance=\"%s\"}", host.Name)
	memQuery := fmt.Sprintf("process_memory_usage_percent{instance=\"%s\"}", host.Name)

	// 并行查询CPU和内存数据
	var cpuResult, memResult *monitormodel.PrometheusQueryResult
	var cpuErr, memErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		cpuResult, cpuErr = s.queryPrometheus(cpuQuery, startTime, endTime, step)
	}()
	go func() {
		defer wg.Done()
		memResult, memErr = s.queryPrometheus(memQuery, startTime, endTime, step)
	}()
	wg.Wait()

	// 构建PID到CPU/内存使用率的映射
	cpuMap := make(map[string]float64)
	memMap := make(map[string]float64)

	if cpuErr == nil && len(cpuResult.Data.Result) > 0 {
		for _, resultItem := range cpuResult.Data.Result {
			pid := resultItem.Metric["pid"]
			var cpuUsage float64
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						cpuUsage, _ = strconv.ParseFloat(valStr, 64)
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					cpuUsage, _ = strconv.ParseFloat(valStr, 64)
				}
			}
			cpuMap[pid] = cpuUsage
		}
	}

	if memErr == nil && len(memResult.Data.Result) > 0 {
		for _, resultItem := range memResult.Data.Result {
			pid := resultItem.Metric["pid"]
			var memUsage float64
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						memUsage, _ = strconv.ParseFloat(valStr, 64)
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					memUsage, _ = strconv.ParseFloat(valStr, 64)
				}
			}
			memMap[pid] = memUsage
		}
	}

	var ports []monitormodel.PortInfo

	// 解析端口数据，并从映射中获取CPU/内存使用率
	if len(promResult.Data.Result) > 0 {
		for _, resultItem := range promResult.Data.Result {
			port := resultItem.Metric["port"]
			pid := resultItem.Metric["pid"]
			service := resultItem.Metric["service"]

			// 获取最新的状态值
			var status int
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						if val, err := strconv.ParseFloat(valStr, 64); err == nil {
							status = int(val)
						}
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					if val, err := strconv.ParseFloat(valStr, 64); err == nil {
						status = int(val)
					}
				}
			}

			// 从映射中获取CPU和内存使用率
			cpuUsage := cpuMap[pid]
			memUsage := memMap[pid]

			ports = append(ports, monitormodel.PortInfo{
				Port:     port,
				PID:      pid,
				Service:  service,
				Status:   status,
				CPUUsage: cpuUsage,
				MemUsage: memUsage,
			})
		}
	}

	// 构造结果
	portResult := &monitormodel.HostPortsResult{
		HostID:     hostID,
		HostName:   host.Name,
		Ports:      ports,
		Total:      len(ports),
		UpdateTime: time.Now().Unix(),
	}

	// 将结果存储到Redis缓存（5分钟过期）
	if cacheData, err := json.Marshal(portResult); err == nil {
		redis.RedisDb.Set(context.Background(), cacheKey, cacheData, 5*time.Minute)
	}

	result.Success(c, portResult)
}
