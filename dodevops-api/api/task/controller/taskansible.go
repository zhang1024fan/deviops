package controller

import (
	"encoding/json"
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
		TaskType:         taskType,
		Name:             name,
		HostGroups:       hostGroups,
		GitRepo:          gitRepo,
		RolesContent:     rolesContent,
		PlaybookContents: playbookContents,
		Variables:        variables,
	}

	// 调用服务层
	c.service.CreateTask(ctx, req)
}

// AnsibleTaskRequest 创建Ansible任务请求
// AnsibleTaskRequest 创建任务请求参数
type AnsibleTaskRequest struct {
	Type       int               `json:"type" form:"type" binding:"required,oneof=1 2"`   // 任务类型(1=手动，2=Git导入)
	Name       string            `json:"name" form:"name" binding:"required"`             // 任务名称
	HostGroups map[string][]uint `json:"hostGroups" form:"hostGroups" binding:"required"` // 主机分组
	GitRepo    string            `json:"gitRepo,omitempty" form:"gitRepo"`                // Git仓库地址(类型为2时必填)
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
		Name:             name,
		Description:      description,
		ClusterName:      clusterName,
		ClusterVersion:   clusterVersion,
		DeploymentMode:   deploymentMode,
		MasterHostIDs:    masterHostIDs,
		WorkerHostIDs:    workerHostIDs,
		EtcdHostIDs:      etcdHostIDs,
		EnabledComponents: enabledComponents,
		PrivateRegistry:  privateRegistry,
		RegistryUsername: registryUsername,
		RegistryPassword: registryPassword,
	}

	// 调用服务层
	c.service.CreateK8sTask(ctx, req)
}
