package controller

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type MonitorController struct {
	monitorService service.MonitorServiceInterface
}

func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: service.NewMonitorService(),
	}
}

// GetHostMetrics 获取主机监控数据
// @Summary 获取主机监控数据
// @Description 获取主机的CPU、内存、磁盘使用率
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/host/{id} [get]
// @Security ApiKeyAuth
func (c *MonitorController) GetHostMetrics(ctx *gin.Context) {
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	c.monitorService.GetHostMetrics(ctx, uint(hostID))
}

// BatchGetHostMetrics 批量获取主机监控数据,主机ID列表，逗号分隔，如：1,2,3
// @Summary 批量获取主机监控数据,主机ID列表，逗号分隔，如：1,2,3
// @Description 批量获取主机的CPU、内存、磁盘使用率
// @Tags 监控
// @Accept json
// @Produce json
// @Param ids query string true "主机ID列表，逗号分隔，如：1,2,3"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/hosts [get]
// @Security ApiKeyAuth
func (c *MonitorController) BatchGetHostMetrics(ctx *gin.Context) {
	idsStr, err := url.QueryUnescape(ctx.Query("ids"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数解码失败")
		return
	}

	if idsStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "缺少主机ID参数")
		return
	}

	// 将逗号分隔的字符串转换为uint数组
	var hostIDs []uint
	for _, idStr := range strings.Split(idsStr, ",") {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "无效的主机ID格式")
			return
		}
		hostIDs = append(hostIDs, uint(id))
	}

	c.monitorService.BatchGetHostMetrics(ctx, hostIDs)
}

// GetHostAllMetricsHistory 获取主机所有指标历史数据
// @Summary 获取主机所有指标历史数据
// @Description 获取主机的CPU、内存、磁盘、网络等所有指标的历史数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Param start query string false "开始时间(格式: 2025-08-02 15:00:00)" example("2025-08-02 15:00:00")
// @Param end query string false  "结束时间(格式: 2025-08-02 16:00:00)" example("2025-08-02 16:00:00")
// @Param duration query string false "时间范围(30m/1h/3h/6h/12h/24h)" example("1h")
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/hosts/{id}/all-metrics [get]
// @Security ApiKeyAuth
func (c *MonitorController) GetHostAllMetricsHistory(ctx *gin.Context) {
	// 获取主机ID
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 获取时间参数
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	duration := ctx.Query("duration")

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	userTimeLayout := "2006-01-02 15:04:05"

	// 如果提供了duration参数，则忽略start和end参数
	if duration != "" {
		// 验证duration参数
		validDurations := map[string]time.Duration{
			"30m": 30 * time.Minute,
			"1h":  1 * time.Hour,
			"3h":  3 * time.Hour,
			"6h":  6 * time.Hour,
			"12h": 12 * time.Hour,
			"24h": 24 * time.Hour,
		}

		durationValue, ok := validDurations[duration]
		if !ok {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间范围参数，请使用: 30m, 1h, 3h, 6h, 12h 或 24h")
			return
		}

		// 计算时间范围
		endTime := now
		startTime := endTime.Add(-durationValue)

		startStr = startTime.Format(userTimeLayout)
		endStr = endTime.Format(userTimeLayout)
	} else {
		// 原有逻辑，验证start和end参数
		if startStr == "" || endStr == "" {
			result.Failed(ctx, http.StatusBadRequest, "请提供开始时间和结束时间参数，或时间范围参数(duration)")
			return
		}

		// URL解码时间参数
		startStr, err = url.QueryUnescape(startStr)
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "开始时间参数解码失败")
			return
		}
		endStr, err = url.QueryUnescape(endStr)
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "结束时间参数解码失败")
			return
		}

		// 去除可能的空格和特殊字符
		startStr = strings.TrimSpace(startStr)
		endStr = strings.TrimSpace(endStr)

		// 预处理时间格式：将+号替换为空格
		startStr = strings.ReplaceAll(startStr, "+", " ")
		endStr = strings.ReplaceAll(endStr, "+", " ")

		// 验证时间格式但不保留变量
		if _, err := time.ParseInLocation(userTimeLayout, startStr, loc); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间格式，请使用: 2006-01-02 15:00:00 或 yyyy-mm-dd+HH:MM:SS")
			return
		}
		if _, err := time.ParseInLocation(userTimeLayout, endStr, loc); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间格式，请使用: 2006-01-02 15:00:00 或 yyyy-mm-dd+HH:MM:SS")
			return
		}
	}

	// 查询所有指标历史数据
	c.monitorService.GetHostAllMetricsHistory(
		ctx,
		uint(hostID),
		startStr,
		endStr,
		"") // 传递空字符串，由service层自动计算step
}

// GetHostMetricHistory 获取主机指标历史数据-CPU、内存、磁盘
// @Summary 获取主机指标历史数据-CPU、内存、磁盘
// @Description 获取主机的CPU、内存、磁盘使用率历史数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Param metric query string false "指标类型(cpu/memory/disk)"
// @Param start query string false "开始时间(格式: 2025-08-02 15:00:00)" example("2025-08-02 15:00:00")
// @Param end query string false  "结束时间(格式: 2025-08-02 16:00:00)" example("2025-08-02 16:00:00")
// @Param duration query string false "时间范围(30m/1h/3h/6h/12h/24h)" example("1h")
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/hosts/{id}/history [get]
// @Security ApiKeyAuth
func (c *MonitorController) GetHostMetricHistory(ctx *gin.Context) {
	// 获取主机ID
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 获取指标类型
	metric := ctx.Query("metric")
	// 如果metric为空，则查询所有指标
	if metric == "" {
		metric = "all"
	}

	// 获取时间参数
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	duration := ctx.Query("duration")

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	userTimeLayout := "2006-01-02 15:04:05"

	// 如果提供了duration参数，则忽略start和end参数
	if duration != "" {
		// 验证duration参数
		validDurations := map[string]time.Duration{
			"30m": 30 * time.Minute,
			"1h":  1 * time.Hour,
			"3h":  3 * time.Hour,
			"6h":  6 * time.Hour,
			"12h": 12 * time.Hour,
			"24h": 24 * time.Hour,
		}

		durationValue, ok := validDurations[duration]
		if !ok {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间范围参数，请使用: 30m, 1h, 3h, 6h, 12h 或 24h")
			return
		}

		// 计算时间范围
		endTime := now
		startTime := endTime.Add(-durationValue)

		startStr = startTime.Format(userTimeLayout)
		endStr = endTime.Format(userTimeLayout)
	} else {
		// 原有逻辑，验证start和end参数
		if startStr == "" || endStr == "" {
			result.Failed(ctx, http.StatusBadRequest, "请提供开始时间和结束时间参数，或时间范围参数(duration)")
			return
		}

		// 去除可能的空格和特殊字符
		startStr = strings.TrimSpace(startStr)
		endStr = strings.TrimSpace(endStr)

		// 预处理时间格式：将+号替换为空格
		startStr = strings.ReplaceAll(startStr, "+", " ")
		endStr = strings.ReplaceAll(endStr, "+", " ")

		// 验证时间格式但不保留变量
		if _, err := time.ParseInLocation(userTimeLayout, startStr, loc); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间格式，请使用: 2006-01-02 15:00:00 或 yyyy-mm-dd+HH:MM:SS")
			return
		}
		if _, err := time.ParseInLocation(userTimeLayout, endStr, loc); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "无效的时间格式，请使用: 2006-01-02 15:00:00 或 yyyy-mm-dd+HH:MM:SS")
			return
		}
	}

	// 查询历史数据
	c.monitorService.GetHostMetricHistory(
		ctx,
		uint(hostID),
		metric,
		startStr,
		endStr,
		"") // 传递空字符串，由service层自动计算step
}

// GetTopProcesses 获取主机TOP进程使用率
// @Summary 获取主机TOP进程使用率
// @Description 获取主机CPU和内存使用率前5名的进程信息
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/hosts/{id}/top-processes [get]
// @Security ApiKeyAuth
func (c *MonitorController) GetTopProcesses(ctx *gin.Context) {
	// 获取主机ID
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 调用服务层获取TOP进程
	c.monitorService.GetTopProcesses(ctx, uint(hostID))
}

// GetHostPorts 获取主机端口信息
// @Summary 获取主机端口信息
// @Description 获取主机所有TCP端口的监听状态、服务名称和进程信息
// @Tags 监控
// @Accept json
// @Produce json
// @Param id path uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/hosts/{id}/ports [get]
// @Security ApiKeyAuth
func (c *MonitorController) GetHostPorts(ctx *gin.Context) {
	// 获取主机ID
	hostIDStr := ctx.Param("id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	// 调用服务层获取端口信息
	c.monitorService.GetHostPorts(ctx, uint(hostID))
}
