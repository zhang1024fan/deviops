package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"dodevops-api/api/task/model"
	"dodevops-api/api/task/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// ListResponse 任务列表响应
type ListResponse struct {
	List  []model.TaskAnsible `json:"list"`
	Total int64               `json:"total"`
}

type TaskAnsibleController struct {
	service service.ITaskAnsibleService
}

func NewTaskAnsibleController(service service.ITaskAnsibleService) *TaskAnsibleController {
	return &TaskAnsibleController{service: service}
}

// CreateTask 创建Ansible任务
// @Summary 创建Ansible任务
// @Description 创建Ansible任务（1=手动，2=Git导入）。K8s部署任务请使用专门的K8s创建接口
// @Tags 任务作业
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "任务名称"
// @Param type formData int true "任务类型(1=手动，2=Git导入)"
// @Param hostGroups formData string true "主机分组JSON"
// @Param gitRepo formData string false "Git仓库地址(type=2时必填)"
// @Param variables formData string false "全局变量JSON"
// @Param playbooks formData file false "playbook文件(type=1时上传)"
// @Param roles formData file false "roles目录(type=1时上传)"
// @Param extra_vars formData string false "额外变量(JSON/YAML字符串)"
// @Param cli_args formData string false "命令行参数"
// @Param use_config formData int false "是否使用配置中心(0=否，1=是)"
// @Param inventory_config_id formData int false "Inventory配置ID"
// @Param global_vars_config_id formData int false "全局变量配置ID"
// @Param extra_vars_config_id formData int false "额外变量配置ID"
// @Param cli_args_config_id formData int false "命令行参数配置ID"
// @Param cron_expr formData string false "Cron表达式(周期任务必填)"
// @Param is_recurring formData int false "是否为周期任务(0=否, 1=是)"
// @Param playbook_paths formData string false "Playbook文件路径列表(JSON数组字符串, type=2时可选)"
// @Param view_id formData int false "视图ID"
// @Success 200 {object} result.Result{data=model.TaskAnsible}
// @Router /api/v1/task/ansible [post]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) CreateTask(ctx *gin.Context) {
	// 获取表单参数
	taskType, _ := strconv.Atoi(ctx.PostForm("type"))
	name := ctx.PostForm("name")
	hostGroupsJSON := ctx.PostForm("hostGroups")
	gitRepo := ctx.PostForm("gitRepo")
	variablesJSON := ctx.PostForm("variables")

	// 验证参数
	if name == "" {
		result.Failed(ctx, http.StatusBadRequest, "任务名称不能为空")
		return
	}

	// 验证任务类型
	if taskType < 1 || taskType > 3 {
		result.Failed(ctx, http.StatusBadRequest, "任务类型必须是1、2或3")
		return
	}

	// Type=2时，Git仓库地址为必填
	if taskType == 2 && gitRepo == "" {
		result.Failed(ctx, http.StatusBadRequest, "Git任务类型必须提供仓库地址")
		return
	}

	// Type=3 K8s任务不应该通过通用接口创建
	if taskType == 3 {
		result.Failed(ctx, http.StatusBadRequest, "K8s任务请使用专门的创建接口")
		return
	}

	// 解析全局变量
	var variables map[string]string
	if variablesJSON != "" {
		if err := json.Unmarshal([]byte(variablesJSON), &variables); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "variables参数格式错误")
			return
		}
	}

	// 解析主机分组(type=1,2时必需)
	var hostGroups map[string][]uint
	if hostGroupsJSON != "" {
		if err := json.Unmarshal([]byte(hostGroupsJSON), &hostGroups); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "hostGroups参数格式错误")
			return
		}
	}

	// 解析Playbook Paths (type=2)
	playbookPathsJSON := ctx.PostForm("playbook_paths")
	var playbookPaths []string
	if playbookPathsJSON != "" {
		if err := json.Unmarshal([]byte(playbookPathsJSON), &playbookPaths); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "playbook_paths参数格式错误")
			return
		}
	}

	// 获取其他新增参数
	extraVars := ctx.PostForm("extra_vars")
	cliArgs := ctx.PostForm("cli_args")
	useConfig, _ := strconv.Atoi(ctx.PostForm("use_config"))
	cronExpr := ctx.PostForm("cron_expr")
	isRecurring, _ := strconv.Atoi(ctx.PostForm("is_recurring"))

	// 解析ID字段
	var inventoryConfigID, globalVarsConfigID, extraVarsConfigID, cliArgsConfigID, viewID *uint

	if val := ctx.PostForm("inventory_config_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		uid := uint(id)
		inventoryConfigID = &uid
	}
	if val := ctx.PostForm("global_vars_config_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		uid := uint(id)
		globalVarsConfigID = &uid
	}
	if val := ctx.PostForm("extra_vars_config_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		uid := uint(id)
		extraVarsConfigID = &uid
	}
	if val := ctx.PostForm("cli_args_config_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		uid := uint(id)
		cliArgsConfigID = &uid
	}
	if val := ctx.PostForm("view_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		uid := uint(id)
		viewID = &uid
	}

	// 处理文件上传
	var rolesFile *multipart.FileHeader
	var playbookFiles []*multipart.FileHeader
	if taskType == 1 {
		rolesFile, _ = ctx.FormFile("roles")
		form, err := ctx.MultipartForm()
		if err == nil {
			playbookFiles = form.File["playbooks"]
		}
	}

	// 处理文件内容
	var rolesContent []byte
	if rolesFile != nil {
		file, err := rolesFile.Open()
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "打开roles文件失败")
			return
		}
		defer file.Close()
		rolesContent, err = io.ReadAll(file)
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "读取roles文件失败")
			return
		}
	}

	// 处理playbook文件
	var playbookContents [][]byte
	for _, f := range playbookFiles {
		file, err := f.Open()
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "打开playbook文件失败")
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			result.Failed(ctx, http.StatusBadRequest, "读取playbook文件失败")
			return
		}
		playbookContents = append(playbookContents, content)
	}

	// 构建请求参数
	req := &service.CreateTaskRequest{
		TaskType:           taskType,
		Name:               name,
		HostGroups:         hostGroups,
		GitRepo:            gitRepo,
		RolesContent:       rolesContent,
		PlaybookContents:   playbookContents,
		Variables:          variables,
		PlaybookPaths:      playbookPaths,
		ExtraVars:          extraVars,
		CliArgs:            cliArgs,
		UseConfig:          useConfig,
		InventoryConfigID:  inventoryConfigID,
		GlobalVarsConfigID: globalVarsConfigID,
		ExtraVarsConfigID:  extraVarsConfigID,
		CliArgsConfigID:    cliArgsConfigID,
		CronExpr:           cronExpr,
		IsRecurring:        isRecurring,
		ViewID:             viewID,
	}

	// 调用服务层
	c.service.CreateTask(ctx, req)
}

// UpdateTask 修改Ansible任务
// @Summary 修改Ansible任务
// @Description 修改Ansible任务基本信息和配置（运行中任务不可修改）
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param request body service.UpdateTaskRequest true "修改任务请求"
// @Success 200 {object} result.Result{data=model.TaskAnsible}
// @Router /api/v1/task/ansible/{id} [put]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) UpdateTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	var req service.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, fmt.Sprintf("参数错误: %v", err))
		return
	}

	c.service.UpdateTask(ctx, uint(id), &req)
}

// DeleteTask 删除Ansible任务
// @Summary 删除Ansible任务
// @Description 删除指定的Ansible任务（级联删除关联的子任务）
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/task/ansible/{id} [delete]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	c.service.DeleteTask(ctx, uint(id))
}

// StartTask 启动Ansible任务
// @Summary 启动Ansible任务
// @Description 启动指定的Ansible任务
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/task/ansible/{id}/start [post]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) StartTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	c.service.StartJob(ctx, uint(id))
}

// GetJobLog 获取任务日志(SSE实现)
// @Summary 获取Ansible任务日志(SSE)
// @Description 通过SSE协议实时获取Ansible任务执行日志
// @Tags 任务作业
// @Accept json
// @Produce text/event-stream
// @Param id path int true "任务ID"
// @Param work_id path int true "子任务ID"
// @Success 200 {object} string "SSE格式的实时日志"
// @Router /api/v1/task/ansible/{id}/log/{work_id} [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetJobLog(ctx *gin.Context) {
	// 先设置SSE响应头，确保正确的Content-Type
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		// SSE格式的错误响应
		ctx.Writer.WriteString("event: error\n")
		ctx.Writer.WriteString("data: 无效的任务ID\n\n")
		ctx.Writer.Flush()
		return
	}

	workID, err := strconv.ParseUint(ctx.Param("work_id"), 10, 64)
	if err != nil {
		// SSE格式的错误响应
		ctx.Writer.WriteString("event: error\n")
		ctx.Writer.WriteString("data: 无效的子任务ID\n\n")
		ctx.Writer.Flush()
		return
	}

	c.service.GetJobLog(ctx, uint(taskID), uint(workID))
}

// List 获取任务列表
// @Summary 获取Ansible任务列表
// @Description 获取Ansible任务列表
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=ListResponse}
// @Router /api/v1/task/ansiblelist [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	c.service.List(ctx, page, size)
}

// GetTask 获取任务详情
// @Summary 获取Ansible任务详情
// @Description 获取Ansible任务详情
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} result.Result{data=model.TaskAnsible}
// @Router /api/v1/task/ansible/{id} [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的任务ID")
		return
	}

	c.service.GetTaskDetail(ctx, uint(id))
}

// GetTasksByName 根据名称模糊查询任务
// @Summary 根据名称模糊查询Ansible任务
// @Description 根据任务名称进行模糊查询
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param name query string true "任务名称（支持模糊查询）"
// @Success 200 {object} result.Result{data=[]model.TaskAnsible}
// @Router /api/v1/task/ansible/query/name [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTasksByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务名称")
		return
	}

	c.service.GetTasksByName(ctx, name)
}

// GetTasksByType 根据类型查询任务
// @Summary 根据类型查询Ansible任务
// @Description 根据任务类型查询（1=手动，2=Git导入，3=K8s部署）
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param type query int true "任务类型（1=手动，2=Git导入，3=K8s部署）"
// @Success 200 {object} result.Result{data=[]model.TaskAnsible}
// @Router /api/v1/task/ansible/query/type [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTasksByType(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	if typeStr == "" {
		result.Failed(ctx, http.StatusBadRequest, "必须提供任务类型")
		return
	}

	taskType, err := strconv.Atoi(typeStr)
	if err != nil || (taskType < 1 || taskType > 3) {
		result.Failed(ctx, http.StatusBadRequest, "任务类型必须是1、2或3（1=手动，2=Git导入，3=K8s部署）")
		return
	}

	c.service.GetTasksByType(ctx, taskType)
}

// GetTasks 查询任务列表 (多条件)
// @Summary 查询任务列表
// @Description 支持任务名称、类型、视图名称的多条件查询和分页
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param name query string false "任务名称(模糊)"
// @Param type query int false "任务类型"
// @Param viewName query string false "视图名称"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=ListResponse}
// @Router /api/v1/task/ansible/query [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTasks(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	name := ctx.Query("name")
	taskType, _ := strconv.Atoi(ctx.Query("type"))
	viewName := ctx.Query("viewName")

	c.service.GetTasks(ctx, name, taskType, viewName, page, size)
}

// CreateK8sTask 创建K8s部署任务
// @Summary 创建K8s部署任务
// @Description 创建K8s集群部署任务
// @Tags 任务作业
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "任务名称"
// @Param description formData string false "任务描述"
// @Param cluster_name formData string true "K8s集群名称"
// @Param cluster_version formData string true "K8s集群版本"
// @Param deployment_mode formData int true "部署模式(1=单节点,2=集群)"
// @Param master_host_ids formData string true "Master节点主机ID数组JSON"
// @Param worker_host_ids formData string false "Worker节点主机ID数组JSON"
// @Param etcd_host_ids formData string true "ETCD节点主机ID数组JSON"
// @Param enabled_components formData string false "启用组件数组JSON"
// @Param private_registry formData string false "私有仓库地址"
// @Param registry_username formData string false "仓库用户名"
// @Param registry_password formData string false "仓库密码"
// @Success 200 {object} result.Result{data=model.TaskAnsible}
// @Router /api/v1/task/k8s [post]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) CreateK8sTask(ctx *gin.Context) {
	// 获取K8s特有参数
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	clusterName := ctx.PostForm("cluster_name")
	clusterVersion := ctx.PostForm("cluster_version")
	deploymentMode, _ := strconv.Atoi(ctx.PostForm("deployment_mode"))
	masterHostIDsJSON := ctx.PostForm("master_host_ids")
	workerHostIDsJSON := ctx.PostForm("worker_host_ids")
	etcdHostIDsJSON := ctx.PostForm("etcd_host_ids")
	enabledComponentsJSON := ctx.PostForm("enabled_components")
	privateRegistry := ctx.PostForm("private_registry")
	registryUsername := ctx.PostForm("registry_username")
	registryPassword := ctx.PostForm("registry_password")

	// 验证K8s必填参数
	if name == "" {
		result.Failed(ctx, http.StatusBadRequest, "任务名称不能为空")
		return
	}
	if clusterName == "" {
		result.Failed(ctx, http.StatusBadRequest, "K8s集群名称不能为空")
		return
	}
	if clusterVersion == "" {
		result.Failed(ctx, http.StatusBadRequest, "K8s集群版本不能为空")
		return
	}
	if deploymentMode < 1 || deploymentMode > 2 {
		result.Failed(ctx, http.StatusBadRequest, "部署模式必须是1(单节点)或2(集群)")
		return
	}
	if masterHostIDsJSON == "" {
		result.Failed(ctx, http.StatusBadRequest, "Master节点主机ID不能为空")
		return
	}
	if etcdHostIDsJSON == "" {
		result.Failed(ctx, http.StatusBadRequest, "ETCD节点主机ID不能为空")
		return
	}

	// 解析主机ID数组
	var masterHostIDs, workerHostIDs, etcdHostIDs []uint

	if err := json.Unmarshal([]byte(masterHostIDsJSON), &masterHostIDs); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "Master节点主机ID格式错误")
		return
	}

	if workerHostIDsJSON != "" {
		if err := json.Unmarshal([]byte(workerHostIDsJSON), &workerHostIDs); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "Worker节点主机ID格式错误")
			return
		}
	}

	if err := json.Unmarshal([]byte(etcdHostIDsJSON), &etcdHostIDs); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "ETCD节点主机ID格式错误")
		return
	}

	// 解析启用组件
	var enabledComponents []string
	if enabledComponentsJSON != "" {
		if err := json.Unmarshal([]byte(enabledComponentsJSON), &enabledComponents); err != nil {
			result.Failed(ctx, http.StatusBadRequest, "启用组件格式错误")
			return
		}
	} else {
		// 默认组件
		enabledComponents = []string{"calico", "coredns"}
	}

	// 构建K8s任务请求
	req := &service.CreateK8sTaskRequest{
		Name:              name,
		Description:       description,
		ClusterName:       clusterName,
		ClusterVersion:    clusterVersion,
		DeploymentMode:    deploymentMode,
		MasterHostIDs:     masterHostIDs,
		WorkerHostIDs:     workerHostIDs,
		EtcdHostIDs:       etcdHostIDs,
		EnabledComponents: enabledComponents,
		PrivateRegistry:   privateRegistry,
		RegistryUsername:  registryUsername,
		RegistryPassword:  registryPassword,
	}

	// 调用服务层
	c.service.CreateK8sTask(ctx, req)
}

// GetTaskHistoryList 获取任务历史记录列表
// @Summary 获取任务历史记录列表
// @Description 获取任务的历史执行记录列表，支持分页
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Router /api/v1/task/ansible/{id}/history [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTaskHistoryList(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的任务ID")
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	c.service.GetTaskHistoryList(ctx, uint(taskID), page, limit)
}

// GetTaskHistoryDetail 获取任务历史记录详情
// @Summary 获取任务历史记录详情
// @Description 获取任务的历史执行详情，包含每个主机的执行日志
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param history_id path int true "历史记录ID"
// @Success 200 {object} result.Result{data=model.TaskAnsibleHistory}
// @Router /api/v1/task/ansible/history/{history_id} [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTaskHistoryDetail(ctx *gin.Context) {
	historyID, err := strconv.ParseUint(ctx.Param("history_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的历史ID")
		return
	}

	c.service.GetTaskHistoryDetail(ctx, uint(historyID))
}

// GetTaskHistoryLog 获取历史记录日志内容
// @Summary 获取历史记录日志内容
// @Description 获取指定子任务历史记录的日志内容
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param work_history_id path int true "子任务历史记录ID"
// @Success 200 {object} result.Result{data=string}
// @Router /api/v1/task/ansible/history/work/{work_history_id}/log [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTaskHistoryLog(ctx *gin.Context) {
	workHistoryID, err := strconv.ParseUint(ctx.Param("work_history_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的ID")
		return
	}

	c.service.GetTaskHistoryLog(ctx, uint(workHistoryID))
}

// GetTaskHistoryLogByDetails 获取历史记录日志内容(通过详细信息)
// @Summary 获取历史记录日志内容(通过详细信息)
// @Description 根据任务ID、WORKID和HistoryID获取历史任务日志
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Param work_id path int true "子任务ID"
// @Param history_id path int true "历史记录ID"
// @Success 200 {object} result.Result{data=string}
// @Router /api/v1/task/ansible/history/detail/task/{task_id}/work/{work_id}/history/{history_id}/log [get]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) GetTaskHistoryLogByDetails(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("task_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的任务ID")
		return
	}
	workID, err := strconv.ParseUint(ctx.Param("work_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的子任务ID")
		return
	}
	historyID, err := strconv.ParseUint(ctx.Param("history_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的历史记录ID")
		return
	}

	c.service.GetTaskHistoryLogByDetails(ctx, uint(taskID), uint(workID), uint(historyID))
}

// DeleteTaskHistory 删除任务历史记录
// @Summary 删除任务历史记录
// @Description 删除指定的任务历史记录及关联的日志文件
// @Tags 任务作业
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param history_id path int true "历史记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/task/ansible/{id}/history/{history_id} [delete]
// @Security ApiKeyAuth
func (c *TaskAnsibleController) DeleteTaskHistory(ctx *gin.Context) {
	historyID, err := strconv.ParseUint(ctx.Param("history_id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的历史记录ID")
		return
	}

	c.service.DeleteTaskHistory(ctx, uint(historyID))
}