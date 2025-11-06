package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	cmdbmodel "dodevops-api/api/cmdb/model"
	configcentermodel "dodevops-api/api/configcenter/model"
	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	taskmodel "dodevops-api/api/task/model"
	taskservice "dodevops-api/api/task/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// IKubeClusterService K8s集群服务接口
type IKubeClusterService interface {
	CreateCluster(c *gin.Context, req *model.CreateKubeClusterRequest)
	GetCluster(c *gin.Context, id uint)
	GetClusterList(c *gin.Context, page, size int)
	UpdateCluster(c *gin.Context, id uint, req *model.UpdateKubeClusterRequest)
	DeleteCluster(c *gin.Context, id uint)
	GetClusterStatus(c *gin.Context, id uint)
	SyncCluster(c *gin.Context, id uint)
	GetClusterDetail(c *gin.Context, id uint)
}

// KubeClusterServiceImpl K8s集群服务实现
type KubeClusterServiceImpl struct {
	dao         *dao.KubeClusterDao
	taskService taskservice.ITaskAnsibleService
}

func NewKubeClusterService(db *gorm.DB) IKubeClusterService {
	return &KubeClusterServiceImpl{
		dao:         dao.NewKubeClusterDao(db),
		taskService: taskservice.NewTaskAnsibleService(db),
	}
}

// CreateCluster 创建K8s集群
func (s *KubeClusterServiceImpl) CreateCluster(c *gin.Context, req *model.CreateKubeClusterRequest) {
	// 检查集群名称是否已存在
	exists, err := s.dao.IsClusterNameExists(req.Name)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("检查集群名称失败: %v", err))
		return
	}
	if exists {
		result.Failed(c, http.StatusBadRequest, "集群名称已存在")
		return
	}

	// 根据集群类型处理不同逻辑
	// 默认为自建集群
	if req.ClusterType == 0 {
		req.ClusterType = model.ClusterTypeSelfBuilt
	}
	
	switch req.ClusterType {
	case model.ClusterTypeSelfBuilt:
		s.createSelfBuiltCluster(c, req)
	case model.ClusterTypeImported:
		s.importCluster(c, req)
	default:
		result.Failed(c, http.StatusBadRequest, "无效的集群类型")
	}
}

// GetCluster 获取集群详情
func (s *KubeClusterServiceImpl) GetCluster(c *gin.Context, id uint) {
	cluster, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群失败: %v", err))
		return
	}

	// 获取集群节点信息（通过K8s API）
	nodes, nodeErr := s.getClusterNodes(cluster)
	if nodeErr != nil {
		// 如果获取节点失败，仍然返回集群基本信息，但标记节点获取失败
		nodes = []model.NodeInfo{}
	}

	// 计算集群概要信息
	summary := s.calculateClusterSummary(nodes)

	response := model.ClusterDetailResponse{
		Cluster: *cluster,
		Nodes:   nodes,
		Summary: summary,
	}

	result.Success(c, response)
}

// GetClusterList 获取集群列表
func (s *KubeClusterServiceImpl) GetClusterList(c *gin.Context, page, size int) {
	clusters, total, err := s.dao.List(page, size)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群列表失败: %v", err))
		return
	}

	response := model.KubeClusterListResponse{
		List:  clusters,
		Total: total,
	}

	result.Success(c, response)
}

// UpdateCluster 更新集群信息
func (s *KubeClusterServiceImpl) UpdateCluster(c *gin.Context, id uint, req *model.UpdateKubeClusterRequest) {
	// 检查集群是否存在
	_, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群失败: %v", err))
		return
	}

	// 检查名称是否被其他集群使用
	if req.Name != "" {
		existing, err := s.dao.GetByName(req.Name)
		if err == nil && existing.ID != id {
			result.Failed(c, http.StatusBadRequest, "集群名称已存在")
			return
		}
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Credential != "" {
		updates["credential"] = req.Credential
	}
	if req.Version != "" {
		updates["version"] = req.Version
	}

	// 如果有更新字段，执行更新
	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.dao.Update(id, updates); err != nil {
			result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("更新集群失败: %v", err))
			return
		}
	}

	// 返回更新后的集群信息
	updatedCluster, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取更新后集群信息失败: %v", err))
		return
	}

	result.Success(c, updatedCluster)
}

// DeleteCluster 删除集群
func (s *KubeClusterServiceImpl) DeleteCluster(c *gin.Context, id uint) {
	// 检查集群是否存在
	_, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群失败: %v", err))
		return
	}

	// 直接删除集群，无状态限制
	if err := s.dao.Delete(id); err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("删除集群失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetClusterStatus 获取集群状态
func (s *KubeClusterServiceImpl) GetClusterStatus(c *gin.Context, id uint) {
	cluster, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群失败: %v", err))
		return
	}

	// 获取节点信息
	nodes, _ := s.getClusterNodes(cluster)
	summary := s.calculateClusterSummary(nodes)

	statusInfo := map[string]interface{}{
		"id":         cluster.ID,
		"name":       cluster.Name,
		"status":     cluster.Status,
		"statusText": cluster.GetStatusText(),
		"clusterType": cluster.ClusterType,
		"summary":    summary,
		"nodes":      nodes,
	}

	result.Success(c, statusInfo)
}


// createSelfBuiltCluster 创建自建集群
func (s *KubeClusterServiceImpl) createSelfBuiltCluster(c *gin.Context, req *model.CreateKubeClusterRequest) {
	// 验证必要参数
	if req.NodeConfig == nil {
		result.Failed(c, http.StatusBadRequest, "节点配置不能为空")
		return
	}
	if req.Version == "" {
		result.Failed(c, http.StatusBadRequest, "集群版本不能为空")
		return
	}

	// 验证主机ID并获取主机信息
	_, err := s.getHostInfos(*req.NodeConfig)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, fmt.Sprintf("获取主机信息失败: %v", err))
		return
	}

	// 创建集群记录
	cluster := &model.KubeCluster{
		Name:        req.Name,
		Version:     req.Version,
		Status:      model.ClusterStatusCreating,
		Description: req.Description,
		ClusterType: model.ClusterTypeSelfBuilt,
	}

	if err := s.dao.Create(cluster); err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("创建集群失败: %v", err))
		return
	}

	// 对于自建集群，始终创建Ansible任务（不管是否自动部署）
	if req.ClusterType == 1 { // 1=自建集群
		go func() {
			// 异步任务错误处理
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("异步创建Ansible任务失败: %v\n", r)
					// 更新集群状态为离线
					s.dao.UpdateStatus(cluster.ID, model.ClusterStatusStopped)
				}
			}()
			
			fmt.Printf("开始异步创建Ansible任务: ID=%d, Name=%s\n", cluster.ID, cluster.Name)
			
			// 构建Ansible任务请求参数
			taskReq := &taskservice.CreateK8sTaskRequest{
				Name:              req.Name + "-deployment",
				Description:       req.Description,
				ClusterName:       req.Name,
				ClusterVersion:    req.Version,
				DeploymentMode:    req.DeploymentMode,
				MasterHostIDs:     req.NodeConfig.MasterHostIDs,
				WorkerHostIDs:     req.NodeConfig.WorkerHostIDs,
				EtcdHostIDs:       req.NodeConfig.EtcdHostIDs,
				EnabledComponents: req.EnabledComponents,
				PrivateRegistry:   req.PrivateRegistry,
				RegistryUsername:  req.RegistryUsername,
				RegistryPassword:  req.RegistryPassword,
				RegistryConfig:    convertRegistryConfig(req.RegistryConfig), // 添加新的嵌套配置
			}

			if taskReq.DeploymentMode == 0 {
				taskReq.DeploymentMode = 1
			}

			// 调用专门的异步任务创建方法
			fmt.Printf("创建Ansible任务参数: %+v\n", taskReq)
			
			taskID, err := s.createK8sTaskAsync(taskReq)
			if err != nil {
				fmt.Printf("创建Ansible任务失败: %v\n", err)
				s.dao.UpdateStatus(cluster.ID, model.ClusterStatusStopped)
				return
			}
			
			fmt.Printf("Ansible任务创建成功，任务ID: %d\n", taskID)
			
			// 如果启用自动部署，则立即执行任务
			// 否则等待用户手动点击执行按钮
			fmt.Printf("AutoDeploy=%v, 任务创建完成，等待手动执行\n", req.AutoDeploy)
		}()
	}

	result.Success(c, cluster)
}

// importCluster 导入集群
func (s *KubeClusterServiceImpl) importCluster(c *gin.Context, req *model.CreateKubeClusterRequest) {
	// 验证必要参数
	if req.Kubeconfig == "" {
		result.Failed(c, http.StatusBadRequest, "kubeconfig不能为空")
		return
	}

	// 验证kubeconfig有效性
	if err := s.validateKubeconfig(req.Kubeconfig); err != nil {
		result.Failed(c, http.StatusBadRequest, fmt.Sprintf("kubeconfig验证失败: %v", err))
		return
	}

	// 从kubeconfig获取集群版本信息
	version := req.Version
	if version == "" {
		version = "unknown" // 如果没有提供版本，从kubeconfig中解析或标记为未知
	}

	// 创建集群记录
	cluster := &model.KubeCluster{
		Name:        req.Name,
		Version:     version,
		Status:      model.ClusterStatusRunning, // 导入的集群默认为运行中
		Description: req.Description,
		ClusterType: model.ClusterTypeImported,
		Credential:  req.Kubeconfig,
	}

	if err := s.dao.Create(cluster); err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("导入集群失败: %v", err))
		return
	}

	// 自动同步集群信息
	go func() {
		// 在goroutine中进行同步，避免阻塞响应
		if err := s.syncClusterInternal(cluster.ID); err != nil {
			// 记录同步失败日志，但不影响导入成功的响应
			fmt.Printf("集群[%d]自动同步失败: %v\n", cluster.ID, err)
		} else {
			fmt.Printf("集群[%d]自动同步成功\n", cluster.ID)
		}
	}()

	result.Success(c, cluster)
}

// getHostInfos 获取主机信息
func (s *KubeClusterServiceImpl) getHostInfos(nodeConfig model.NodeConfig) (map[uint]cmdbmodel.CmdbHost, error) {
	// 收集所有主机ID
	allHostIDs := make(map[uint]bool)
	for _, id := range nodeConfig.MasterHostIDs {
		allHostIDs[id] = true
	}
	for _, id := range nodeConfig.WorkerHostIDs {
		allHostIDs[id] = true
	}
	for _, id := range nodeConfig.EtcdHostIDs {
		allHostIDs[id] = true
	}

	// 转换为切片
	var hostIDs []uint
	for id := range allHostIDs {
		hostIDs = append(hostIDs, id)
	}

	// 查询主机信息
	var hosts []cmdbmodel.CmdbHost
	if err := s.dao.DB.Where("id IN ?", hostIDs).Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("查询主机信息失败: %v", err)
	}

	// 构建主机信息映射
	hostMap := make(map[uint]cmdbmodel.CmdbHost)
	for _, host := range hosts {
		hostMap[host.ID] = host
	}

	// 验证所有主机是否都存在
	for _, id := range hostIDs {
		if _, exists := hostMap[id]; !exists {
			return nil, fmt.Errorf("主机ID %d 不存在", id)
		}
	}

	return hostMap, nil
}

// getClusterNodes 通过K8s API获取集群节点信息
func (s *KubeClusterServiceImpl) getClusterNodes(cluster *model.KubeCluster) ([]model.NodeInfo, error) {
	// TODO: 实现K8s API调用获取节点信息
	// 这里暂时返回空切片，实际应该调用K8s API
	
	// 示例代码：
	// k8sClient := NewK8sClient(cluster.Credential)
	// nodes, err := k8sClient.GetNodes()
	// if err != nil {
	//     return nil, err
	// }
	// return convertToNodeInfo(nodes), nil
	
	return []model.NodeInfo{}, nil
}

// calculateClusterSummary 计算集群概要信息
func (s *KubeClusterServiceImpl) calculateClusterSummary(nodes []model.NodeInfo) model.ClusterSummary {
	summary := model.ClusterSummary{
		TotalNodes:  len(nodes),
		ReadyNodes:  0,
		MasterNodes: 0,
		WorkerNodes: 0,
	}

	for _, node := range nodes {
		if node.Status == "Ready" {
			summary.ReadyNodes++
		}
		if node.Role == "master" {
			summary.MasterNodes++
		} else if node.Role == "worker" {
			summary.WorkerNodes++
		}
	}

	return summary
}

// createK8sTaskAsync 异步创建K8s任务（不依赖gin.Context）
func (s *KubeClusterServiceImpl) createK8sTaskAsync(req *taskservice.CreateK8sTaskRequest) (uint, error) {
	// 直接调用底层逻辑，避免使用gin.Context
	
	// 1. 验证主机信息
	hostInfos, err := s.getK8sHostInfoForAsync(req)
	if err != nil {
		return 0, fmt.Errorf("获取主机信息失败: %v", err)
	}

	// 2. 创建任务记录
	task := &taskmodel.TaskAnsible{
		Name:        req.Name,
		Description: req.Description,
		Type:        3, // K8s任务类型
		GitRepo:     "git@gitee.com:zhang_fan1024/zf-k8s.git", // K8s部署仓库
		HostGroups:  s.buildK8sHostGroupsForAsync(req),
		AllHostIDs:  s.buildK8sAllHostIDsForAsync(req),
		GlobalVars:  s.buildK8sGlobalVarsForAsync(req),
		Status:      1, // 等待中
		TaskCount:   1, // K8s任务固定为1个
	}

	// 直接使用数据库创建任务
	if err := s.dao.DB.Create(task).Error; err != nil {
		return 0, fmt.Errorf("创建任务失败: %v", err)
	}

	// 3. 创建项目目录和生成配置
	projectDir := fmt.Sprintf("./task/%d/%s", task.ID, task.Name)
	if err := s.setupK8sTaskFiles(projectDir, req, hostInfos, task.ID); err != nil {
		return 0, fmt.Errorf("设置任务文件失败: %v", err)
	}
	
	// 4. 创建子任务（K8s部署脚本）
	if err := s.createK8sSubTaskAsync(task.ID, projectDir); err != nil {
		return 0, fmt.Errorf("创建K8s子任务失败: %v", err)
	}

	return task.ID, nil
}

// startK8sTaskAsync 异步启动K8s任务
func (s *KubeClusterServiceImpl) startK8sTaskAsync(taskID uint) error {
	// 暂时简化，只更新任务状态为运行中
	// 实际应该启动任务执行，但这需要复杂的实现
	
	// 更新任务状态为运行中
	if err := s.dao.DB.Model(&taskmodel.TaskAnsible{}).Where("id = ?", taskID).Update("status", 2).Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}
	
	fmt.Printf("任务 %d 已标记为运行中（简化版实现）\n", taskID)
	
	// TODO: 这里应该实际启动Ansible执行
	// 可以考虑调用shell命令或者其他方式
	
	return nil
}

// createK8sSubTaskAsync 创建K8s子任务（异步版本）
func (s *KubeClusterServiceImpl) createK8sSubTaskAsync(taskID uint, projectDir string) error {
	// 导入TaskAnsibleWork模型
	work := &taskmodel.TaskAnsibleWork{
		TaskID:        taskID,
		EntryFileName: "deploy-simple.sh",
		EntryFilePath: "./scripts/deploy-simple.sh",
		LogPath:       filepath.Join("./logs/ansible", fmt.Sprintf("%d", taskID), "deploy-k8s.log"),
		Status:        1, // 等待中
	}
	
	// 确保日志目录存在
	logDir := filepath.Dir(work.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}
	
	// 保存子任务
	if err := s.dao.DB.Create(work).Error; err != nil {
		return fmt.Errorf("创建子任务失败: %v", err)
	}
	
	return nil
}

// 辅助方法：获取主机信息（异步版本）
func (s *KubeClusterServiceImpl) getK8sHostInfoForAsync(req *taskservice.CreateK8sTaskRequest) (map[string][]taskservice.K8sNodeInfo, error) {
	// 收集所有主机ID
	allIDs := make(map[uint]bool)
	for _, id := range req.MasterHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.WorkerHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.EtcdHostIDs {
		allIDs[id] = true
	}

	var hostIDs []uint
	for id := range allIDs {
		hostIDs = append(hostIDs, id)
	}

	// 查询主机信息
	var hosts []cmdbmodel.CmdbHost
	if err := s.dao.DB.Where("id IN ?", hostIDs).Find(&hosts).Error; err != nil {
		return nil, err
	}

	// 构建映射
	hostMap := make(map[uint]cmdbmodel.CmdbHost)
	sshKeyIDs := make([]uint, 0)
	for _, host := range hosts {
		hostMap[host.ID] = host
		if host.SSHKeyID > 0 {
			sshKeyIDs = append(sshKeyIDs, host.SSHKeyID)
		}
	}
	
	// 查询SSH凭据信息
	var sshKeys []configcentermodel.EcsAuth
	if len(sshKeyIDs) > 0 {
		s.dao.DB.Where("id IN ?", sshKeyIDs).Find(&sshKeys)
	}
	
	// 构建SSH凭据映射
	sshKeyMap := make(map[uint]string)
	for _, key := range sshKeys {
		sshKeyMap[key.ID] = key.Password
	}

	// 按角色分组
	result := map[string][]taskservice.K8sNodeInfo{
		"masters": {},
		"workers": {},
		"etcd":    {},
	}

	for _, id := range req.MasterHostIDs {
		if host, exists := hostMap[id]; exists {
			password := ""
			if host.SSHKeyID > 0 {
				password = sshKeyMap[host.SSHKeyID]
			}
			result["masters"] = append(result["masters"], taskservice.K8sNodeInfo{
				Name:     host.HostName,
				IP:       host.SSHIP,
				User:     host.SSHName,
				Password: password,
			})
		}
	}

	for _, id := range req.WorkerHostIDs {
		if host, exists := hostMap[id]; exists {
			password := ""
			if host.SSHKeyID > 0 {
				password = sshKeyMap[host.SSHKeyID]
			}
			result["workers"] = append(result["workers"], taskservice.K8sNodeInfo{
				Name:     host.HostName,
				IP:       host.SSHIP,
				User:     host.SSHName,
				Password: password,
			})
		}
	}

	for _, id := range req.EtcdHostIDs {
		if host, exists := hostMap[id]; exists {
			password := ""
			if host.SSHKeyID > 0 {
				password = sshKeyMap[host.SSHKeyID]
			}
			result["etcd"] = append(result["etcd"], taskservice.K8sNodeInfo{
				Name:     host.HostName,
				IP:       host.SSHIP,
				User:     host.SSHName,
				Password: password,
			})
		}
	}

	return result, nil
}

// 辅助方法：构建主机分组
func (s *KubeClusterServiceImpl) buildK8sHostGroupsForAsync(req *taskservice.CreateK8sTaskRequest) string {
	hostGroups := map[string][]uint{
		"masters": req.MasterHostIDs,
		"workers": req.WorkerHostIDs,
		"etcd":    req.EtcdHostIDs,
	}
	
	jsonData, _ := json.Marshal(hostGroups)
	return string(jsonData)
}

// 辅助方法：构建所有主机ID
func (s *KubeClusterServiceImpl) buildK8sAllHostIDsForAsync(req *taskservice.CreateK8sTaskRequest) string {
	allIDs := make(map[uint]bool)
	for _, id := range req.MasterHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.WorkerHostIDs {
		allIDs[id] = true
	}
	for _, id := range req.EtcdHostIDs {
		allIDs[id] = true
	}

	var ids []uint
	for id := range allIDs {
		ids = append(ids, id)
	}
	
	jsonData, _ := json.Marshal(ids)
	return string(jsonData)
}

// 辅助方法：构建全局变量
func (s *KubeClusterServiceImpl) buildK8sGlobalVarsForAsync(req *taskservice.CreateK8sTaskRequest) string {
	vars := map[string]string{
		"cluster_name":       req.ClusterName,
		"cluster_version":    req.ClusterVersion,
		"deployment_mode":    fmt.Sprintf("%d", req.DeploymentMode),
		"enabled_components": strings.Join(req.EnabledComponents, ","),
	}
	
	jsonData, _ := json.Marshal(vars)
	return string(jsonData)
}

// 辅助方法：设置任务文件
func (s *KubeClusterServiceImpl) setupK8sTaskFiles(projectDir string, req *taskservice.CreateK8sTaskRequest, hostInfos map[string][]taskservice.K8sNodeInfo, taskID uint) error {
	// 创建目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}

	// 2. 克隆K8s Git仓库
	gitRepo := "git@gitee.com:zhang_fan1024/zf-k8s.git"
	if err := s.cloneGitRepositoryAsync(gitRepo, projectDir); err != nil {
		return fmt.Errorf("克隆K8s仓库失败: %v", err)
	}
	
	// 3. 生成config.json配置文件
	if err := s.generateK8sConfigAsync(projectDir, req, hostInfos, taskID); err != nil {
		return fmt.Errorf("生成K8s配置失败: %v", err)
	}
	
	return nil
}



// cloneGitRepositoryAsync 异步版本的Git仓库克隆
func (s *KubeClusterServiceImpl) cloneGitRepositoryAsync(gitRepo, projectDir string) error {
	fmt.Printf("[DEBUG] 开始克隆Git仓库: %s 到 %s\n", gitRepo, projectDir)
	// 使用git命令克隆仓库
	cmd := exec.Command("git", "clone", gitRepo, projectDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone失败: %v", err)
	}
	fmt.Printf("[DEBUG] Git仓库克隆成功\n")
	return nil
}

// generateK8sConfigAsync 异步版本的K8s配置生成
func (s *KubeClusterServiceImpl) generateK8sConfigAsync(projectDir string, req *taskservice.CreateK8sTaskRequest, hostInfos map[string][]taskservice.K8sNodeInfo, taskID uint) error {
	// 构建配置数据 - 按照正确的格式结构
	// 处理镜像仓库配置 - 支持新旧两种格式
	var registryConfig map[string]interface{}
	fmt.Printf("[DEBUG] Registry配置检查: RegistryConfig=%+v, UsePrivateRegistry=%v\n", req.RegistryConfig, req.RegistryConfig != nil && req.RegistryConfig.UsePrivateRegistry)
	fmt.Printf("[DEBUG] 旧格式Registry: PrivateRegistry=%s, Username=%s\n", req.PrivateRegistry, req.RegistryUsername)
	
	if req.RegistryConfig != nil && req.RegistryConfig.PrivateRegistry != "" {
		// 使用新的嵌套格式
		fmt.Printf("[DEBUG] 使用新的嵌套Registry格式\n")
		registryConfig = map[string]interface{}{
			"private_registry": req.RegistryConfig.PrivateRegistry,
			"username":        req.RegistryConfig.RegistryUsername,
			"password":        req.RegistryConfig.RegistryPassword,
		}
	} else {
		// 兼容旧格式
		fmt.Printf("[DEBUG] 使用旧的平铺Registry格式\n")
		registryConfig = map[string]interface{}{
			"private_registry": req.PrivateRegistry,
			"username":        req.RegistryUsername,
			"password":        req.RegistryPassword,
		}
	}

	config := map[string]interface{}{
		"cluster": map[string]interface{}{
			"name":            req.ClusterName,
			"version":         req.ClusterVersion,
			"deployment_mode": req.DeploymentMode,
		},
		"nodes": map[string]interface{}{},
		"components": map[string]interface{}{
			"enabled": req.EnabledComponents,
		},
		"registry": registryConfig,
	}
	
	// 添加节点信息到nodes结构中
	nodes := config["nodes"].(map[string]interface{})
	if masters, ok := hostInfos["masters"]; ok && len(masters) > 0 {
		nodes["masters"] = masters
	}
	if workers, ok := hostInfos["workers"]; ok && len(workers) > 0 {
		nodes["workers"] = workers
	}
	if etcd, ok := hostInfos["etcd"]; ok && len(etcd) > 0 {
		nodes["etcd"] = etcd
	}
	
	// 生成JSON配置文件
	configPath := filepath.Join(projectDir, "config.json")
	configData, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	fmt.Printf("[DEBUG] K8s配置文件生成成功: %s\n", configPath)
	return nil
}

// SyncCluster 同步集群信息
func (s *KubeClusterServiceImpl) SyncCluster(c *gin.Context, id uint) {
	// 获取集群信息
	cluster, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
		} else {
			result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("查询集群失败: %v", err))
		}
		return
	}

	// 检查集群是否有kubeconfig凭证，有凭证的集群都可以同步
	if cluster.Credential == "" {
		result.Failed(c, http.StatusBadRequest, "集群缺少kubeconfig配置，无法同步")
		return
	}

	// 通过K8s API同步集群信息
	syncData, err := s.syncFromK8sAPI(cluster.Credential)
	if err != nil {
		// 同步失败时，判断是否为超时错误，将集群状态设为离线
		now := time.Now()
		updates := map[string]interface{}{
			"status":       model.ClusterStatusStopped, // 3-离线
			"last_sync_at": &now,
			"updated_at":   now,
		}

		// 更新集群状态为离线
		if updateErr := s.dao.Update(id, updates); updateErr != nil {
			result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("同步失败且更新状态失败: %v, %v", err, updateErr))
			return
		}

		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("同步集群信息失败，已将集群状态设为离线: %v", err))
		return
	}

	// 更新集群信息到数据库
	now := time.Now()
	updates := map[string]interface{}{
		"version":      syncData.Version,
		"status":       syncData.Status,
		"node_count":   syncData.NodeCount,
		"ready_nodes":  syncData.ReadyNodes,
		"master_nodes": syncData.MasterNodes,
		"worker_nodes": syncData.WorkerNodes,
		"last_sync_at": &now,
		"updated_at":   now,
	}

	if err := s.dao.Update(id, updates); err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("更新集群信息失败: %v", err))
		return
	}

	// 重新查询更新后的集群信息
	updatedCluster, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("查询更新后的集群失败: %v", err))
		return
	}

	result.Success(c, updatedCluster)
}

// K8sSyncData K8s同步数据结构
type K8sSyncData struct {
	Version     string
	Status      int
	NodeCount   int
	ReadyNodes  int
	MasterNodes int
	WorkerNodes int
}

// syncFromK8sAPI 从K8s API同步集群信息
func (s *KubeClusterServiceImpl) syncFromK8sAPI(kubeconfig string) (*K8sSyncData, error) {
	// 解析kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("解析kubeconfig失败: %v", err)
	}

	// 设置超时配置
	config.Timeout = 5 * time.Second

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取集群版本信息
	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("获取集群版本失败: %v", err)
	}

	// 获取节点列表
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 统计节点信息
	var readyNodes, masterNodes, workerNodes int
	for _, node := range nodes.Items {
		// 统计就绪节点
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				readyNodes++
				break
			}
		}

		// 统计节点类型
		if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
			masterNodes++
		} else if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
			masterNodes++
		} else {
			workerNodes++
		}
	}

	// 确定集群状态
	status := model.ClusterStatusRunning
	if readyNodes == 0 {
		status = model.ClusterStatusStopped // 离线
	} else if readyNodes < len(nodes.Items) {
		status = model.ClusterStatusRunning // 部分节点不可用但仍保持运行状态
	}

	return &K8sSyncData{
		Version:     versionInfo.GitVersion,
		Status:      status,
		NodeCount:   len(nodes.Items),
		ReadyNodes:  readyNodes,
		MasterNodes: masterNodes,
		WorkerNodes: workerNodes,
	}, nil
}

// convertRegistryConfig 转换RegistryConfig类型
func convertRegistryConfig(src *model.RegistryConfig) *taskservice.RegistryConfig {
	if src == nil {
		return nil
	}
	return &taskservice.RegistryConfig{
		PrivateRegistry:    src.PrivateRegistry,
		RegistryUsername:   src.RegistryUsername,
		RegistryPassword:   src.RegistryPassword,
		UsePrivateRegistry: src.UsePrivateRegistry,
	}
}

// syncClusterInternal 内部同步集群信息（不依赖gin.Context）
func (s *KubeClusterServiceImpl) syncClusterInternal(clusterID uint) error {
	// 获取集群信息
	cluster, err := s.dao.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("查询集群失败: %v", err)
	}

	// 检查集群是否有kubeconfig凭证
	if cluster.Credential == "" {
		return fmt.Errorf("集群缺少kubeconfig配置，无法同步")
	}

	// 通过K8s API同步集群信息
	syncData, err := s.syncFromK8sAPI(cluster.Credential)
	if err != nil {
		// 同步失败时，将集群状态设为离线
		now := time.Now()
		updates := map[string]interface{}{
			"status":       model.ClusterStatusStopped, // 3-离线
			"last_sync_at": &now,
		}

		// 更新集群状态为离线
		if updateErr := s.dao.Update(clusterID, updates); updateErr != nil {
			return fmt.Errorf("同步失败且更新状态失败: %v, %v", err, updateErr)
		}

		return fmt.Errorf("同步集群信息失败，已将集群状态设为离线: %v", err)
	}

	// 更新集群信息到数据库
	now := time.Now()
	updates := map[string]interface{}{
		"version":      syncData.Version,
		"node_count":   syncData.NodeCount,
		"ready_nodes":  syncData.ReadyNodes,
		"master_nodes": syncData.MasterNodes,
		"worker_nodes": syncData.WorkerNodes,
		"status":       syncData.Status,
		"last_sync_at": &now,
	}

	return s.dao.Update(clusterID, updates)
}

// validateKubeconfig 验证kubeconfig有效性
func (s *KubeClusterServiceImpl) validateKubeconfig(kubeconfig string) error {
	// 解析kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return fmt.Errorf("kubeconfig格式错误: %v", err)
	}

	// 创建K8s客户端测试连接
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 测试连接 - 尝试获取集群版本信息
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("连接K8s集群失败: %v", err)
	}
	
	// 测试节点访问权限
	_, err = clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return fmt.Errorf("缺少节点访问权限: %v", err)
	}

	return nil
}

// GetClusterDetail 获取集群详细信息
func (s *KubeClusterServiceImpl) GetClusterDetail(c *gin.Context, id uint) {
	// 获取基本集群信息
	cluster, err := s.dao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群失败: %v", err))
		return
	}

	// 检查集群是否有kubeconfig凭证
	if cluster.Credential == "" {
		result.Failed(c, http.StatusBadRequest, "集群缺少kubeconfig配置，无法获取详细信息")
		return
	}

	// 获取详细信息
	detail, err := s.getClusterDetailInfo(cluster)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("获取集群详细信息失败: %v", err))
		return
	}

	result.Success(c, detail)
}

// getClusterDetailInfo 获取集群详细信息
func (s *KubeClusterServiceImpl) getClusterDetailInfo(cluster *model.KubeCluster) (*model.ClusterDetailResponse, error) {
	// 解析kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("解析kubeconfig失败: %v", err)
	}

	// 设置超时
	config.Timeout = 10 * time.Second

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 并行获取各种信息
	var nodes []model.NodeInfo
	var components []model.ComponentInfo
	var workloads model.WorkloadSummary
	var network model.NetworkInfo
	var events []model.ClusterEvent
	var monitoring model.MonitoringInfo
	var runtime model.RuntimeSummary

	// 获取节点信息
	nodes, _ = s.getDetailedNodes(ctx, clientset)

	// 获取组件信息
	components, _ = s.getClusterComponents(ctx, clientset)

	// 获取工作负载统计
	workloads, _ = s.getWorkloadSummary(ctx, clientset)

	// 获取网络配置
	network, _ = s.getNetworkInfo(ctx, clientset)

	// 获取集群事件
	events, _ = s.getClusterEvents(ctx, clientset)

	// 获取监控信息
	monitoring, _ = s.getMonitoringInfo(ctx, clientset)

	// 获取运行时信息
	runtime, _ = s.getRuntimeInfo(ctx, clientset)

	// 计算概要信息
	summary := s.calculateClusterSummary(nodes)

	return &model.ClusterDetailResponse{
		Cluster:     *cluster,
		Nodes:       nodes,
		Summary:     summary,
		Components:  components,
		Workloads:   workloads,
		Network:     network,
		Events:      events,
		Monitoring:  monitoring,
		Runtime:     runtime,
	}, nil
}

// getDetailedNodes 获取详细的节点信息
func (s *KubeClusterServiceImpl) getDetailedNodes(ctx context.Context, clientset *kubernetes.Clientset) ([]model.NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeInfos []model.NodeInfo
	for _, node := range nodes.Items {
		nodeInfo := model.NodeInfo{
			Name:        node.Name,
			InternalIP:  "",
			ExternalIP:  "",
			Status:      "NotReady",
			Version:     node.Status.NodeInfo.KubeletVersion,
			OS:          node.Status.NodeInfo.OSImage,
			Capacity:    make(map[string]string),
			Allocatable: make(map[string]string),
			Conditions:  make([]model.NodeCondition, 0),
		}

		// 获取节点IP地址
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				nodeInfo.InternalIP = addr.Address
			} else if addr.Type == "ExternalIP" {
				nodeInfo.ExternalIP = addr.Address
			}
		}

		// 获取节点状态
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				nodeInfo.Status = "Ready"
			}
			nodeInfo.Conditions = append(nodeInfo.Conditions, model.NodeCondition{
				Type:   string(condition.Type),
				Status: string(condition.Status),
				Reason: condition.Reason,
			})
		}

		// 获取节点角色
		if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
			nodeInfo.Role = "master"
		} else if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
			nodeInfo.Role = "master"
		} else {
			nodeInfo.Role = "worker"
		}

		// 获取资源信息
		for resourceName, quantity := range node.Status.Capacity {
			nodeInfo.Capacity[string(resourceName)] = quantity.String()
		}
		for resourceName, quantity := range node.Status.Allocatable {
			nodeInfo.Allocatable[string(resourceName)] = quantity.String()
		}

		nodeInfos = append(nodeInfos, nodeInfo)
	}

	return nodeInfos, nil
}

// getClusterComponents 获取集群组件信息
func (s *KubeClusterServiceImpl) getClusterComponents(ctx context.Context, clientset *kubernetes.Clientset) ([]model.ComponentInfo, error) {
	var components []model.ComponentInfo

	// 获取系统命名空间的Pod作为组件
	systemNamespaces := []string{"kube-system", "kube-public", "kube-node-lease"}

	for _, namespace := range systemNamespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		// 按组件名分组
		componentMap := make(map[string]*model.ComponentInfo)

		for _, pod := range pods.Items {
			componentName := ""
			// 从Pod名称推断组件名
			if strings.Contains(pod.Name, "kube-apiserver") {
				componentName = "kube-apiserver"
			} else if strings.Contains(pod.Name, "etcd") {
				componentName = "etcd"
			} else if strings.Contains(pod.Name, "kube-controller-manager") {
				componentName = "kube-controller-manager"
			} else if strings.Contains(pod.Name, "kube-scheduler") {
				componentName = "kube-scheduler"
			} else if strings.Contains(pod.Name, "kube-proxy") {
				componentName = "kube-proxy"
			} else if strings.Contains(pod.Name, "coredns") {
				componentName = "coredns"
			} else {
				continue
			}

			if comp, exists := componentMap[componentName]; exists {
				// 更新状态 - 如果有任何Pod运行，组件状态为Running
				if pod.Status.Phase == "Running" && comp.Status != "Running" {
					comp.Status = "Running"
				}
			} else {
				status := "NotReady"
				if pod.Status.Phase == "Running" {
					status = "Running"
				}

				componentMap[componentName] = &model.ComponentInfo{
					Name:      componentName,
					Namespace: namespace,
					Status:    status,
					Version:   "unknown",
					Type:      "system",
				}
			}
		}

		// 将组件添加到结果中
		for _, comp := range componentMap {
			components = append(components, *comp)
		}
	}

	return components, nil
}

// getWorkloadSummary 获取工作负载统计
func (s *KubeClusterServiceImpl) getWorkloadSummary(ctx context.Context, clientset *kubernetes.Clientset) (model.WorkloadSummary, error) {
	summary := model.WorkloadSummary{}

	// 获取Deployments
	deployments, _ := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	summary.TotalDeployments = len(deployments.Items)

	// 获取StatefulSets
	statefulSets, _ := clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	summary.TotalStatefulSets = len(statefulSets.Items)

	// 获取DaemonSets
	daemonSets, _ := clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	summary.TotalDaemonSets = len(daemonSets.Items)

	// 获取Jobs
	jobs, _ := clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	summary.TotalJobs = len(jobs.Items)

	// 获取CronJobs
	cronJobs, _ := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	summary.TotalCronJobs = len(cronJobs.Items)

	// 获取Pods
	pods, _ := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	summary.TotalPods = len(pods.Items)

	runningPods := 0
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			runningPods++
		}
	}
	summary.RunningPods = runningPods

	return summary, nil
}

// getNetworkInfo 获取网络配置信息
func (s *KubeClusterServiceImpl) getNetworkInfo(ctx context.Context, clientset *kubernetes.Clientset) (model.NetworkInfo, error) {
	network := model.NetworkInfo{}

	// 获取kube-system命名空间中的ConfigMap来获取网络配置
	configMaps, err := clientset.CoreV1().ConfigMaps("kube-system").List(ctx, metav1.ListOptions{})
	if err == nil {
		for _, cm := range configMaps.Items {
			if cm.Name == "kubeadm-config" {
				// 从kubeadm-config中解析网络配置
				if clusterConfig, exists := cm.Data["ClusterConfiguration"]; exists {
					if strings.Contains(clusterConfig, "serviceSubnet") {
						// 简单解析ServiceCIDR（实际应该用yaml解析）
						lines := strings.Split(clusterConfig, "\n")
						for _, line := range lines {
							if strings.Contains(line, "serviceSubnet") {
								parts := strings.Split(line, ":")
								if len(parts) > 1 {
									network.ServiceCIDR = strings.TrimSpace(parts[1])
								}
							}
						}
					}
				}
			}
		}
	}

	// 获取API Server端点
	endpoints, err := clientset.CoreV1().Endpoints("default").Get(ctx, "kubernetes", metav1.GetOptions{})
	if err == nil && len(endpoints.Subsets) > 0 && len(endpoints.Subsets[0].Addresses) > 0 {
		network.APIServerEndpoint = endpoints.Subsets[0].Addresses[0].IP
		if len(endpoints.Subsets[0].Ports) > 0 {
			network.APIServerEndpoint += fmt.Sprintf(":%d", endpoints.Subsets[0].Ports[0].Port)
		}
	}

	// 默认值设置
	if network.ServiceCIDR == "" {
		network.ServiceCIDR = "10.96.0.0/12" // Kubernetes默认
	}
	network.PodCIDR = "10.244.0.0/16"      // Flannel默认
	network.DNSService = "CoreDNS"
	network.NetworkPlugin = "Flannel"      // 需要实际检测
	network.ProxyMode = "iptables"

	return network, nil
}

// getClusterEvents 获取集群事件
func (s *KubeClusterServiceImpl) getClusterEvents(ctx context.Context, clientset *kubernetes.Clientset) ([]model.ClusterEvent, error) {
	events, err := clientset.CoreV1().Events("").List(ctx, metav1.ListOptions{
		Limit: 50, // 限制最近50个事件
	})
	if err != nil {
		return nil, err
	}

	var clusterEvents []model.ClusterEvent
	for _, event := range events.Items {
		clusterEvent := model.ClusterEvent{
			Type:           event.Type,
			Reason:         event.Reason,
			Message:        event.Message,
			Source:         event.Source.Component,
			FirstTime:      event.FirstTimestamp.Format(time.RFC3339),
			LastTime:       event.LastTimestamp.Format(time.RFC3339),
			Count:          event.Count,
			InvolvedObject: fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name),
		}
		clusterEvents = append(clusterEvents, clusterEvent)
	}

	return clusterEvents, nil
}

// getMonitoringInfo 获取监控信息
func (s *KubeClusterServiceImpl) getMonitoringInfo(ctx context.Context, clientset *kubernetes.Clientset) (model.MonitoringInfo, error) {
	monitoring := model.MonitoringInfo{}

	// 获取节点资源统计
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return monitoring, err
	}

	var totalCPU, totalMemory int64
	var allocatableCPU, allocatableMemory int64

	for _, node := range nodes.Items {
		if cpu, exists := node.Status.Capacity["cpu"]; exists {
			// 使用MilliValue获取毫核数
			totalCPU += cpu.MilliValue()
		}
		if memory, exists := node.Status.Capacity["memory"]; exists {
			totalMemory += memory.Value()
		}
		if cpu, exists := node.Status.Allocatable["cpu"]; exists {
			allocatableCPU += cpu.MilliValue()
		}
		if memory, exists := node.Status.Allocatable["memory"]; exists {
			allocatableMemory += memory.Value()
		}
	}

	// CPU监控信息
	monitoring.CPU = model.ClusterResourceMetrics{
		Total:     fmt.Sprintf("%.1f cores", float64(totalCPU)/1000),        // 毫核转换为核
		Used:      "0 cores", // 需要metrics-server支持
		Available: fmt.Sprintf("%.1f cores", float64(allocatableCPU)/1000),  // 毫核转换为核
		UsageRate: 0, // 需要实际监控数据
		RequestRate: 0,
	}

	// 内存监控信息
	monitoring.Memory = model.ClusterResourceMetrics{
		Total:     fmt.Sprintf("%d Mi", totalMemory/(1024*1024)),
		Used:      "0 Mi", // 需要metrics-server支持
		Available: fmt.Sprintf("%d Mi", allocatableMemory/(1024*1024)),
		UsageRate: 0, // 需要实际监控数据
		RequestRate: 0,
	}

	// 网络监控（示例数据）
	monitoring.Network = model.NetworkMetrics{
		InboundTraffic:  "0 MB/s",
		OutboundTraffic: "0 MB/s",
		PacketsIn:       0,
		PacketsOut:      0,
	}

	// 存储监控
	pvs, _ := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	pvcs, _ := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	storageClasses, _ := clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})

	boundPVs := 0
	for _, pv := range pvs.Items {
		if pv.Status.Phase == "Bound" {
			boundPVs++
		}
	}

	var scNames []string
	for _, sc := range storageClasses.Items {
		scNames = append(scNames, sc.Name)
	}

	monitoring.Storage = model.StorageMetrics{
		TotalPVs:       len(pvs.Items),
		BoundPVs:       boundPVs,
		TotalPVCs:      len(pvcs.Items),
		StorageClasses: scNames,
	}

	return monitoring, nil
}

// getRuntimeInfo 获取运行时信息
func (s *KubeClusterServiceImpl) getRuntimeInfo(ctx context.Context, clientset *kubernetes.Clientset) (model.RuntimeSummary, error) {
	runtime := model.RuntimeSummary{}

	// 获取集群版本
	version, err := clientset.Discovery().ServerVersion()
	if err == nil {
		runtime.KubernetesVersion = version.GitVersion
		runtime.APIServerVersion = version.GitVersion
	}

	// 获取第一个节点的运行时信息作为示例
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err == nil && len(nodes.Items) > 0 {
		node := nodes.Items[0]
		runtime.ContainerRuntime = node.Status.NodeInfo.ContainerRuntimeVersion
		runtime.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
	}

	// 从kube-system命名空间获取etcd和CoreDNS版本
	pods, _ := clientset.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "etcd") && len(pod.Spec.Containers) > 0 {
			image := pod.Spec.Containers[0].Image
			if strings.Contains(image, ":") {
				parts := strings.Split(image, ":")
				if len(parts) > 1 {
					runtime.EtcdVersion = parts[len(parts)-1]
				}
			}
		}
		if strings.Contains(pod.Name, "coredns") && len(pod.Spec.Containers) > 0 {
			image := pod.Spec.Containers[0].Image
			if strings.Contains(image, ":") {
				parts := strings.Split(image, ":")
				if len(parts) > 1 {
					runtime.CoreDNSVersion = parts[len(parts)-1]
				}
			}
		}
	}

	// 计算集群运行时间（从最老的节点创建时间算起）
	if len(nodes.Items) > 0 {
		oldestTime := nodes.Items[0].CreationTimestamp.Time
		for _, node := range nodes.Items {
			if node.CreationTimestamp.Time.Before(oldestTime) {
				oldestTime = node.CreationTimestamp.Time
			}
		}
		upTime := time.Since(oldestTime)
		runtime.UpTime = fmt.Sprintf("%.0f days", upTime.Hours()/24)
	}

	return runtime, nil
}
