package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// IK8sNodesService K8s节点服务接口
type IK8sNodesService interface {
	GetNodes(c *gin.Context, clusterId uint)
	GetNodeDetail(c *gin.Context, clusterId uint, nodeName string)
	GetNodeDetailEnhanced(c *gin.Context, clusterId uint, nodeName string)
	AddTaint(c *gin.Context, clusterId uint, nodeName string, req *model.AddTaintRequest)
	RemoveTaint(c *gin.Context, clusterId uint, nodeName string, req *model.RemoveTaintRequest)
	AddLabel(c *gin.Context, clusterId uint, nodeName string, req *model.AddLabelRequest)
	RemoveLabel(c *gin.Context, clusterId uint, nodeName string, req *model.RemoveLabelRequest)
	CordonNode(c *gin.Context, clusterId uint, nodeName string, req *model.CordonNodeRequest)
	DrainNode(c *gin.Context, clusterId uint, nodeName string, req *model.DrainNodeRequest)
	GetNodeResourceAllocation(c *gin.Context, clusterId uint, nodeName string)
}

// K8sNodesServiceImpl K8s节点服务实现
type K8sNodesServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

func NewK8sNodesService(db *gorm.DB) IK8sNodesService {
	return &K8sNodesServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// GetNodes 获取集群的所有节点信息
func (s *K8sNodesServiceImpl) GetNodes(c *gin.Context, clusterId uint) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法获取节点信息")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点列表
	nodes, err := s.fetchNodes(clientset)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取节点列表失败: "+err.Error())
		return
	}

	result.Success(c, nodes)
}

// GetNodeDetail 获取单个节点的详细信息
func (s *K8sNodesServiceImpl) GetNodeDetail(c *gin.Context, clusterId uint, nodeName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法获取节点信息")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 尝试创建Metrics客户端（可选）
	var metricsClient *metricsclientset.Clientset
	metricsClient, _ = s.createMetricsClient(cluster.Credential) // 忽略错误，允许没有metrics server

	// 获取节点详细信息
	nodeDetail, err := s.fetchNodeDetail(clientset, metricsClient, nodeName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取节点详细信息失败: "+err.Error())
		return
	}

	result.Success(c, nodeDetail)
}

// createK8sClient 创建K8s客户端
func (s *K8sNodesServiceImpl) createK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// fetchNodes 获取所有节点信息
func (s *K8sNodesServiceImpl) fetchNodes(clientset *kubernetes.Clientset) ([]model.K8sNode, error) {
	// 获取节点列表
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 获取所有Pod信息用于计算节点上的Pod数量
	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 构建节点到Pod的映射
	nodePodMap := make(map[string][]corev1.Pod)
	for _, pod := range podList.Items {
		if pod.Spec.NodeName != "" {
			nodePodMap[pod.Spec.NodeName] = append(nodePodMap[pod.Spec.NodeName], pod)
		}
	}

	var nodes []model.K8sNode
	for _, node := range nodeList.Items {
		nodePods := nodePodMap[node.Name]
		k8sNode := s.convertToK8sNode(&node, nodePods)
		nodes = append(nodes, k8sNode)
	}

	return nodes, nil
}

// fetchNodeDetail 获取节点详细信息
func (s *K8sNodesServiceImpl) fetchNodeDetail(clientset *kubernetes.Clientset, metricsClient *metricsclientset.Clientset, nodeName string) (*model.K8sNodeDetail, error) {
	// 获取节点信息
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 获取节点上的Pod列表
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, err
	}

	// 获取节点相关事件
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Node", nodeName),
	})
	// 如果事件获取失败，不影响整体响应，只是事件为空
	if err != nil {
		events = &corev1.EventList{Items: []corev1.Event{}}
	}

	// 转换为详细节点信息
	nodeDetail := &model.K8sNodeDetail{
		K8sNode: s.convertToK8sNode(node, pods.Items),
		Pods:    s.convertToPodInfosWithMetrics(pods.Items, metricsClient),
		Events:  s.convertToEventInfos(events.Items),
		Metrics: s.calculateNodeMetricsWithReal(node, pods.Items, metricsClient),
	}

	return nodeDetail, nil
}

// convertToK8sNode 转换为K8sNode结构
func (s *K8sNodesServiceImpl) convertToK8sNode(node *corev1.Node, pods []corev1.Pod) model.K8sNode {
	podCount := len(pods)
	// 获取节点IP地址
	var internalIP, externalIP string
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case corev1.NodeInternalIP:
			internalIP = addr.Address
		case corev1.NodeExternalIP:
			externalIP = addr.Address
		}
	}

	// 获取节点状态
	var status string = "Unknown"
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				status = "Ready"
			} else {
				status = "NotReady"
			}
			break
		}
	}

	// 转换节点状态条件
	var conditions []model.NodeCondition
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, model.NodeCondition{
			Type:   string(condition.Type),
			Status: string(condition.Status),
			Reason: condition.Reason,
		})
	}

	// 获取节点角色 - 支持多种角色组合显示
	var roles []string
	
	// 检查各种角色标签
	if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
		roles = append(roles, "control-plane")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
		roles = append(roles, "master")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/worker"]; exists {
		roles = append(roles, "worker")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/etcd"]; exists {
		roles = append(roles, "etcd")
	}
	
	// 如果没有找到任何角色标签，检查是否为master节点
	if len(roles) == 0 {
		// 通过其他方式判断是否为master节点
		if _, exists := node.Labels["kubernetes.io/role"]; exists {
			if node.Labels["kubernetes.io/role"] == "master" {
				roles = append(roles, "master")
			}
		} else {
			// 默认为worker节点
			roles = append(roles, "worker")
		}
	}
	
	// 将角色数组转换为逗号分隔的字符串，与kubectl显示格式一致
	rolesStr := strings.Join(roles, ",")
	if rolesStr == "" {
		rolesStr = "<none>"
	}
	
	// 保持兼容性，为Configuration.Role设置主要角色
	primaryRole := "worker"
	if len(roles) > 0 {
		if contains(roles, "control-plane") || contains(roles, "master") {
			primaryRole = "master"
		} else {
			primaryRole = roles[0]
		}
	}

	// 计算Pod容量
	maxPods := int64(110) // 默认值
	if podCapacity, ok := node.Status.Capacity[corev1.ResourcePods]; ok {
		maxPods = podCapacity.Value()
	}

	// 转换污点信息
	var taints []model.NodeTaint
	for _, taint := range node.Spec.Taints {
		taints = append(taints, model.NodeTaint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: string(taint.Effect),
		})
	}

	return model.K8sNode{
		Name:       node.Name,
		InternalIP: internalIP,
		ExternalIP: externalIP,
		Status:     status,
		Roles:      rolesStr,
		Conditions: conditions,
		Configuration: model.NodeConfiguration{
			Role:          primaryRole,
			Architecture:  node.Status.NodeInfo.Architecture,
			KernelVersion: node.Status.NodeInfo.KernelVersion,
			OSImage:       node.Status.NodeInfo.OSImage,
			Labels:        node.Labels,
			Annotations:   node.Annotations,
		},
		PodMetrics: model.PodMetrics{
			Allocated: podCount,
			Total:     int(maxPods),
		},
		Resources: s.calculateNodeResources(node, pods),
		Runtime: model.RuntimeInfo{
			KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
			ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
			KubeProxyVersion:        "", // KubeProxyVersion is deprecated, set to empty
			OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
			OSImage:                 node.Status.NodeInfo.OSImage,
		},
		Scheduling: model.NodeSchedulingInfo{
			Unschedulable: node.Spec.Unschedulable,
			Taints:        taints,
		},
		CreatedAt: node.CreationTimestamp.Format(time.RFC3339),
	}
}

// convertToPodInfos 转换为Pod信息列表
func (s *K8sNodesServiceImpl) convertToPodInfos(pods []corev1.Pod) []model.PodInfo {
	var podInfos []model.PodInfo
	for _, pod := range pods {
		// 计算Pod的资源请求作为使用估算
		var totalCPURequests, totalMemoryRequests int64
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				totalCPURequests += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				totalMemoryRequests += memory.Value()
			}
		}

		cpuUsage := "0m"
		if totalCPURequests > 0 {
			cpuUsage = fmt.Sprintf("%dm", totalCPURequests)
		}

		memUsage := "0Mi"
		if totalMemoryRequests > 0 {
			memUsage = fmt.Sprintf("%dMi", totalMemoryRequests/(1024*1024))
		}

		podInfos = append(podInfos, model.PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			CPUUsage:  cpuUsage,
			MemUsage:  memUsage,
		})
	}
	return podInfos
}

// convertToPodInfosWithMetrics 转换为Pod信息列表（包含真实metrics）
func (s *K8sNodesServiceImpl) convertToPodInfosWithMetrics(pods []corev1.Pod, metricsClient *metricsclientset.Clientset) []model.PodInfo {
	var podInfos []model.PodInfo

	for _, pod := range pods {
		cpuUsage := "0m"
		memUsage := "0Mi"

		// 尝试获取真实的metrics数据
		if metricsClient != nil {
			podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if err == nil && podMetrics != nil {
				// 计算总使用量
				var totalCPU, totalMemory resource.Quantity
				for _, containerMetric := range podMetrics.Containers {
					totalCPU.Add(containerMetric.Usage[corev1.ResourceCPU])
					totalMemory.Add(containerMetric.Usage[corev1.ResourceMemory])
				}

				if totalCPU.MilliValue() > 0 {
					cpuUsage = fmt.Sprintf("%dm", totalCPU.MilliValue())
				}
				if totalMemory.Value() > 0 {
					memUsage = fmt.Sprintf("%dMi", totalMemory.Value()/(1024*1024))
				}
			}
		}

		// 如果没有metrics数据，回退到使用请求值
		if cpuUsage == "0m" || memUsage == "0Mi" {
			var totalCPURequests, totalMemoryRequests int64
			for _, container := range pod.Spec.Containers {
				if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
					totalCPURequests += cpu.MilliValue()
				}
				if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
					totalMemoryRequests += memory.Value()
				}
			}

			if cpuUsage == "0m" && totalCPURequests > 0 {
				cpuUsage = fmt.Sprintf("%dm", totalCPURequests)
			}
			if memUsage == "0Mi" && totalMemoryRequests > 0 {
				memUsage = fmt.Sprintf("%dMi", totalMemoryRequests/(1024*1024))
			}
		}

		podInfos = append(podInfos, model.PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			CPUUsage:  cpuUsage,
			MemUsage:  memUsage,
		})
	}
	return podInfos
}

// convertToEventInfos 转换为事件信息列表
func (s *K8sNodesServiceImpl) convertToEventInfos(events []corev1.Event) []model.EventInfo {
	var eventInfos []model.EventInfo
	for _, event := range events {
		eventInfos = append(eventInfos, model.EventInfo{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
			Count:     int(event.Count),
		})
	}
	return eventInfos
}

// calculateNodeMetrics 计算节点指标
func (s *K8sNodesServiceImpl) calculateNodeMetrics(node *corev1.Node, pods []corev1.Pod) model.NodeMetrics {
	metrics := model.NodeMetrics{
		CPUUsagePercentage:    0.0,
		MemoryUsagePercentage: 0.0,
		DiskUsagePercentage:   0.0,
		NetworkInBytes:        0,
		NetworkOutBytes:       0,
	}

	// 计算基于资源请求的使用率估算
	var totalCPURequests, totalMemoryRequests int64
	var allocatableCPU, allocatableMemory int64

	// 获取可分配资源
	if cpu, exists := node.Status.Allocatable[corev1.ResourceCPU]; exists {
		allocatableCPU = cpu.MilliValue()
	}
	if memory, exists := node.Status.Allocatable[corev1.ResourceMemory]; exists {
		allocatableMemory = memory.Value()
	}

	// 计算所有Pod的资源请求总和
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				totalCPURequests += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				totalMemoryRequests += memory.Value()
			}
		}
	}

	// 基于请求计算使用率估算（不是真实使用率，但比0更有意义）
	if allocatableCPU > 0 {
		metrics.CPUUsagePercentage = float64(totalCPURequests) / float64(allocatableCPU) * 100
		if metrics.CPUUsagePercentage > 100 {
			metrics.CPUUsagePercentage = 100
		}
	}

	if allocatableMemory > 0 {
		metrics.MemoryUsagePercentage = float64(totalMemoryRequests) / float64(allocatableMemory) * 100
		if metrics.MemoryUsagePercentage > 100 {
			metrics.MemoryUsagePercentage = 100
		}
	}

	// 简单的磁盘使用率估算（基于Pod数量）
	podCount := len(pods)
	if podCount > 0 {
		// 假设每个Pod平均占用1%的磁盘空间，最大不超过80%
		metrics.DiskUsagePercentage = float64(podCount) * 1.0
		if metrics.DiskUsagePercentage > 80 {
			metrics.DiskUsagePercentage = 80
		}
	}

	return metrics
}

// calculateNodeMetricsWithReal 计算节点指标（包含真实数据）
func (s *K8sNodesServiceImpl) calculateNodeMetricsWithReal(node *corev1.Node, pods []corev1.Pod, metricsClient *metricsclientset.Clientset) model.NodeMetrics {
	metrics := model.NodeMetrics{
		CPUUsagePercentage:    0.0,
		MemoryUsagePercentage: 0.0,
		DiskUsagePercentage:   0.0,
		NetworkInBytes:        0,
		NetworkOutBytes:       0,
	}

	var totalCPURequests, totalMemoryRequests int64
	var allocatableCPU, allocatableMemory int64

	// 获取可分配资源
	if cpu, exists := node.Status.Allocatable[corev1.ResourceCPU]; exists {
		allocatableCPU = cpu.MilliValue()
	}
	if memory, exists := node.Status.Allocatable[corev1.ResourceMemory]; exists {
		allocatableMemory = memory.Value()
	}

	// 如果有metrics客户端，尝试获取真实使用数据
	var totalCPUUsage, totalMemoryUsage int64

	if metricsClient != nil {
		for _, pod := range pods {
			// 累加资源请求
			for _, container := range pod.Spec.Containers {
				if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
					totalCPURequests += cpu.MilliValue()
				}
				if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
					totalMemoryRequests += memory.Value()
				}
			}

			// 尝试获取真实使用数据
			podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if err == nil && podMetrics != nil {
				for _, containerMetric := range podMetrics.Containers {
					if cpu, exists := containerMetric.Usage[corev1.ResourceCPU]; exists {
						totalCPUUsage += cpu.MilliValue()
					}
					if memory, exists := containerMetric.Usage[corev1.ResourceMemory]; exists {
						totalMemoryUsage += memory.Value()
					}
				}
			}
		}

		// 如果获取到真实使用数据，使用真实数据计算使用率
		if totalCPUUsage > 0 && allocatableCPU > 0 {
			metrics.CPUUsagePercentage = float64(totalCPUUsage) / float64(allocatableCPU) * 100
			if metrics.CPUUsagePercentage > 100 {
				metrics.CPUUsagePercentage = 100
			}
		} else if totalCPURequests > 0 && allocatableCPU > 0 {
			// 回退到使用请求值
			metrics.CPUUsagePercentage = float64(totalCPURequests) / float64(allocatableCPU) * 100
			if metrics.CPUUsagePercentage > 100 {
				metrics.CPUUsagePercentage = 100
			}
		}

		if totalMemoryUsage > 0 && allocatableMemory > 0 {
			metrics.MemoryUsagePercentage = float64(totalMemoryUsage) / float64(allocatableMemory) * 100
			if metrics.MemoryUsagePercentage > 100 {
				metrics.MemoryUsagePercentage = 100
			}
		} else if totalMemoryRequests > 0 && allocatableMemory > 0 {
			// 回退到使用请求值
			metrics.MemoryUsagePercentage = float64(totalMemoryRequests) / float64(allocatableMemory) * 100
			if metrics.MemoryUsagePercentage > 100 {
				metrics.MemoryUsagePercentage = 100
			}
		}
	} else {
		// 没有metrics客户端，使用原来的逻辑
		return s.calculateNodeMetrics(node, pods)
	}

	// 简单的磁盘使用率估算（基于Pod数量）
	podCount := len(pods)
	if podCount > 0 {
		metrics.DiskUsagePercentage = float64(podCount) * 1.0
		if metrics.DiskUsagePercentage > 80 {
			metrics.DiskUsagePercentage = 80
		}
	}

	return metrics
}

// calculateNodeResources 计算节点资源信息
func (s *K8sNodesServiceImpl) calculateNodeResources(node *corev1.Node, pods []corev1.Pod) model.NodeResources {
	// 获取节点容量和可分配资源
	var capacityCPU, allocatableCPU, totalCPURequests int64
	var totalMemoryRequests int64

	if cpu, exists := node.Status.Capacity[corev1.ResourceCPU]; exists {
		capacityCPU = cpu.MilliValue()
	}
	if cpu, exists := node.Status.Allocatable[corev1.ResourceCPU]; exists {
		allocatableCPU = cpu.MilliValue()
	}

	// 计算所有Pod的资源请求总和
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				totalCPURequests += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				totalMemoryRequests += memory.Value()
			}
		}
	}

	return model.NodeResources{
		CPU: model.ResourceInfo{
			Capacity:    fmt.Sprintf("%.1f", float64(capacityCPU)/1000),
			Allocatable: fmt.Sprintf("%dm", allocatableCPU),
			Requests:    fmt.Sprintf("%dm", totalCPURequests),
			Usage:       fmt.Sprintf("%dm", totalCPURequests), // 使用请求值作为使用估算
		},
		Memory: model.ResourceInfo{
			Capacity:    node.Status.Capacity.Memory().String(),
			Allocatable: node.Status.Allocatable.Memory().String(),
			Requests:    fmt.Sprintf("%dKi", totalMemoryRequests/1024),
			Usage:       fmt.Sprintf("%dKi", totalMemoryRequests/1024), // 使用请求值作为使用估算
		},
	}
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// AddTaint 为节点添加污点
func (s *K8sNodesServiceImpl) AddTaint(c *gin.Context, clusterId uint, nodeName string, req *model.AddTaintRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 添加污点
	newTaint := corev1.Taint{
		Key:    req.Key,
		Value:  req.Value,
		Effect: corev1.TaintEffect(req.Effect),
	}

	// 检查污点是否已存在
	for _, taint := range node.Spec.Taints {
		if taint.Key == req.Key && taint.Effect == corev1.TaintEffect(req.Effect) {
			result.Failed(c, http.StatusBadRequest, "污点已存在")
			return
		}
	}

	node.Spec.Taints = append(node.Spec.Taints, newTaint)

	// 更新节点
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "添加污点失败: "+err.Error())
		return
	}

	result.Success(c, "污点添加成功")
}

// RemoveTaint 移除节点污点
func (s *K8sNodesServiceImpl) RemoveTaint(c *gin.Context, clusterId uint, nodeName string, req *model.RemoveTaintRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 移除污点
	var newTaints []corev1.Taint
	found := false
	for _, taint := range node.Spec.Taints {
		if taint.Key == req.Key && (req.Effect == "" || taint.Effect == corev1.TaintEffect(req.Effect)) {
			found = true
			continue // 跳过要删除的污点
		}
		newTaints = append(newTaints, taint)
	}

	if !found {
		result.Failed(c, http.StatusNotFound, "污点不存在")
		return
	}

	node.Spec.Taints = newTaints

	// 更新节点
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "移除污点失败: "+err.Error())
		return
	}

	result.Success(c, "污点移除成功")
}

// AddLabel 为节点添加标签
func (s *K8sNodesServiceImpl) AddLabel(c *gin.Context, clusterId uint, nodeName string, req *model.AddLabelRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 添加标签
	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}
	node.Labels[req.Key] = req.Value

	// 更新节点
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "添加标签失败: "+err.Error())
		return
	}

	result.Success(c, "标签添加成功")
}

// RemoveLabel 移除节点标签
func (s *K8sNodesServiceImpl) RemoveLabel(c *gin.Context, clusterId uint, nodeName string, req *model.RemoveLabelRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 检查标签是否存在
	if node.Labels == nil || node.Labels[req.Key] == "" {
		result.Failed(c, http.StatusNotFound, "标签不存在")
		return
	}

	// 移除标签
	delete(node.Labels, req.Key)

	// 更新节点
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "移除标签失败: "+err.Error())
		return
	}

	result.Success(c, "标签移除成功")
}

// CordonNode 封锁/解封节点
func (s *K8sNodesServiceImpl) CordonNode(c *gin.Context, clusterId uint, nodeName string, req *model.CordonNodeRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 设置调度状态
	node.Spec.Unschedulable = req.Unschedulable

	// 更新节点
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新节点调度状态失败: "+err.Error())
		return
	}

	message := "节点解封成功"
	if req.Unschedulable {
		message = "节点封锁成功"
	}

	result.Success(c, message)
}

// DrainNode 驱逐节点
func (s *K8sNodesServiceImpl) DrainNode(c *gin.Context, clusterId uint, nodeName string, req *model.DrainNodeRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 首先封锁节点
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	node.Spec.Unschedulable = true
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "封锁节点失败: "+err.Error())
		return
	}

	// 获取节点上的Pod列表
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	// 驱逐Pod
	var drainedPods []string
	for _, pod := range pods.Items {
		// 跳过DaemonSet管理的Pod（如果设置了忽略）
		if req.IgnoreDaemonSets {
			for _, ownerRef := range pod.OwnerReferences {
				if ownerRef.Kind == "DaemonSet" {
					continue
				}
			}
		}

		// 跳过系统Pod（除非强制删除）
		if !req.Force && pod.Namespace == "kube-system" {
			continue
		}

		// 设置优雅终止时间
		gracePeriodSeconds := int64(req.GracePeriodSeconds)
		if gracePeriodSeconds <= 0 {
			gracePeriodSeconds = 30 // 默认30秒
		}

		// 删除Pod
		err = clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriodSeconds,
		})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("驱逐Pod %s/%s失败: %s", pod.Namespace, pod.Name, err.Error()))
			return
		}

		drainedPods = append(drainedPods, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	}

	result.Success(c, map[string]interface{}{
		"message":     "节点驱逐成功",
		"drainedPods": drainedPods,
	})
}

// GetNodeResourceAllocation 获取节点资源分配详情
func (s *K8sNodesServiceImpl) GetNodeResourceAllocation(c *gin.Context, clusterId uint, nodeName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取节点信息
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 获取节点上的Pod列表
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	// 计算资源分配
	allocation := s.calculateResourceAllocation(node, pods.Items)

	result.Success(c, allocation)
}

// calculateResourceAllocation 计算节点资源分配
func (s *K8sNodesServiceImpl) calculateResourceAllocation(node *corev1.Node, pods []corev1.Pod) *model.NodeResourceAllocation {
	// 获取节点容量和可分配资源
	capacity := make(map[string]string)
	allocatable := make(map[string]string)
	
	for resource, quantity := range node.Status.Capacity {
		capacity[string(resource)] = quantity.String()
	}
	
	for resource, quantity := range node.Status.Allocatable {
		allocatable[string(resource)] = quantity.String()
	}

	// 计算已分配资源
	allocated := make(map[string]string)
	cpuRequests := int64(0)
	memoryRequests := int64(0)
	
	var podList []model.PodResourceInfo
	
	for _, pod := range pods {
		podCPU := int64(0)
		podMemory := int64(0)
		requests := make(map[string]string)
		limits := make(map[string]string)
		
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				if cpu := container.Resources.Requests.Cpu(); cpu != nil {
					podCPU += cpu.MilliValue()
				}
				if memory := container.Resources.Requests.Memory(); memory != nil {
					podMemory += memory.Value()
				}
			}
			
			// 记录requests和limits
			for resource, quantity := range container.Resources.Requests {
				requests[string(resource)] = quantity.String()
			}
			for resource, quantity := range container.Resources.Limits {
				limits[string(resource)] = quantity.String()
			}
		}
		
		cpuRequests += podCPU
		memoryRequests += podMemory
		
		podList = append(podList, model.PodResourceInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Requests:  requests,
			Limits:    limits,
		})
	}
	
	allocated["cpu"] = fmt.Sprintf("%dm", cpuRequests)
	allocated["memory"] = fmt.Sprintf("%d", memoryRequests)

	return &model.NodeResourceAllocation{
		NodeName:    node.Name,
		Capacity:    capacity,
		Allocatable: allocatable,
		Allocated:   allocated,
		PodList:     podList,
	}
}

// GetNodeDetailEnhanced 获取增强的节点详细信息
func (s *K8sNodesServiceImpl) GetNodeDetailEnhanced(c *gin.Context, clusterId uint, nodeName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群凭证
	if cluster.Credential == "" {
		result.Failed(c, http.StatusBadRequest, "集群缺少kubeconfig配置")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClientWithTimeout(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 尝试创建Metrics客户端（可选）
	var metricsClient *metricsclientset.Clientset
	metricsClient, _ = s.createMetricsClient(cluster.Credential) // 忽略错误，允许没有metrics server

	// 获取增强的节点详细信息
	nodeDetail, err := s.fetchNodeDetailEnhanced(clientset, metricsClient, nodeName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取节点详细信息失败: "+err.Error())
		return
	}

	result.Success(c, nodeDetail)
}

// createK8sClientWithTimeout 创建带超时的K8s客户端
func (s *K8sNodesServiceImpl) createK8sClientWithTimeout(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	// 设置超时
	config.Timeout = 10 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// createMetricsClient 创建Metrics客户端
func (s *K8sNodesServiceImpl) createMetricsClient(kubeconfig string) (*metricsclientset.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return metricsClient, nil
}

// fetchNodeDetailEnhanced 获取增强的节点详细信息
func (s *K8sNodesServiceImpl) fetchNodeDetailEnhanced(clientset *kubernetes.Clientset, metricsClient *metricsclientset.Clientset, nodeName string) (*model.NodeDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取节点信息
	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取节点信息失败: %v", err)
	}

	// 获取节点上的Pod列表
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, fmt.Errorf("获取Pod列表失败: %v", err)
	}

	// 构建增强的节点详情响应
	nodeDetail := &model.NodeDetailResponse{}

	// 填充基本信息
	nodeDetail.Name = node.Name
	nodeDetail.CreatedAt = node.CreationTimestamp.Format(time.RFC3339)
	nodeDetail.UID = string(node.UID)
	nodeDetail.ProviderID = node.Spec.ProviderID

	// 填充IP地址信息
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case corev1.NodeInternalIP:
			nodeDetail.InternalIP = addr.Address
		case corev1.NodeExternalIP:
			nodeDetail.ExternalIP = addr.Address
		case corev1.NodeHostName:
			nodeDetail.Hostname = addr.Address
		}
	}

	// 填充系统信息
	nodeDetail.OSImage = node.Status.NodeInfo.OSImage
	nodeDetail.KernelVersion = node.Status.NodeInfo.KernelVersion
	nodeDetail.Architecture = node.Status.NodeInfo.Architecture
	nodeDetail.OperatingSystem = node.Status.NodeInfo.OperatingSystem
	nodeDetail.MachineID = node.Status.NodeInfo.MachineID
	nodeDetail.SystemUUID = node.Status.NodeInfo.SystemUUID
	nodeDetail.BootID = node.Status.NodeInfo.BootID

	// 填充K8s组件版本
	nodeDetail.KubeletVersion = node.Status.NodeInfo.KubeletVersion
	nodeDetail.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
	nodeDetail.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion

	// 填充资源信息
	nodeDetail.Capacity = make(map[string]string)
	nodeDetail.Allocatable = make(map[string]string)
	for resourceName, quantity := range node.Status.Capacity {
		nodeDetail.Capacity[string(resourceName)] = quantity.String()
	}
	for resourceName, quantity := range node.Status.Allocatable {
		nodeDetail.Allocatable[string(resourceName)] = quantity.String()
	}

	// 填充调度信息
	nodeDetail.Unschedulable = node.Spec.Unschedulable
	nodeDetail.Labels = node.Labels
	nodeDetail.Annotations = node.Annotations

	// 填充污点信息
	var taints []model.NodeTaint
	for _, taint := range node.Spec.Taints {
		taints = append(taints, model.NodeTaint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: string(taint.Effect),
		})
	}
	nodeDetail.Taints = taints

	// 填充状态信息
	nodeDetail.Status = "NotReady"
	var conditions []model.NodeCondition
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			nodeDetail.Status = "Ready"
		}
		conditions = append(conditions, model.NodeCondition{
			Type:   string(condition.Type),
			Status: string(condition.Status),
			Reason: condition.Reason,
		})
	}
	nodeDetail.Conditions = conditions

	// 填充Pod CIDR信息
	nodeDetail.PodCIDR = node.Spec.PodCIDR
	nodeDetail.PodCIDRs = node.Spec.PodCIDRs

	// 填充监控信息
	monitoring, err := s.calculateNodeMonitoring(node, pods.Items)
	if err == nil {
		nodeDetail.Monitoring = monitoring
	}

	// 填充Pod信息
	podInfo, podList := s.calculateNodePodInfoWithMetrics(pods.Items, metricsClient)
	nodeDetail.PodInfo = podInfo
	nodeDetail.PodList = podList

	return nodeDetail, nil
}

// calculateNodeMonitoring 计算节点监控信息
func (s *K8sNodesServiceImpl) calculateNodeMonitoring(node *corev1.Node, pods []corev1.Pod) (model.NodeMonitoringInfo, error) {
	monitoring := model.NodeMonitoringInfo{}

	// 获取节点总资源
	var totalCPU, totalMemory int64
	var allocatableCPU, allocatableMemory int64

	if cpu, exists := node.Status.Capacity[corev1.ResourceCPU]; exists {
		totalCPU = cpu.MilliValue()
	}
	if memory, exists := node.Status.Capacity[corev1.ResourceMemory]; exists {
		totalMemory = memory.Value()
	}
	if cpu, exists := node.Status.Allocatable[corev1.ResourceCPU]; exists {
		allocatableCPU = cpu.MilliValue()
	}
	if memory, exists := node.Status.Allocatable[corev1.ResourceMemory]; exists {
		allocatableMemory = memory.Value()
	}

	// 计算Pod资源请求和限制
	var cpuRequests, cpuLimits, memoryRequests, memoryLimits int64
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				cpuRequests += cpu.MilliValue()
			}
			if cpu, exists := container.Resources.Limits[corev1.ResourceCPU]; exists {
				cpuLimits += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				memoryRequests += memory.Value()
			}
			if memory, exists := container.Resources.Limits[corev1.ResourceMemory]; exists {
				memoryLimits += memory.Value()
			}
		}
	}

	// CPU监控信息
	monitoring.CPU = model.NodeResourceUsage{
		Total:     fmt.Sprintf("%.1f cores", float64(totalCPU)/1000),
		Used:      "0 cores", // 需要metrics-server支持
		Available: fmt.Sprintf("%.1f cores", float64(allocatableCPU)/1000),
		Requests:  fmt.Sprintf("%.1f cores", float64(cpuRequests)/1000),
		Limits:    fmt.Sprintf("%.1f cores", float64(cpuLimits)/1000),
		UsageRate: 0, // 需要metrics-server支持
		RequestRate: float64(cpuRequests) / float64(allocatableCPU) * 100,
	}

	// 内存监控信息
	monitoring.Memory = model.NodeResourceUsage{
		Total:     fmt.Sprintf("%d Mi", totalMemory/(1024*1024)),
		Used:      "0 Mi", // 需要metrics-server支持
		Available: fmt.Sprintf("%d Mi", allocatableMemory/(1024*1024)),
		Requests:  fmt.Sprintf("%d Mi", memoryRequests/(1024*1024)),
		Limits:    fmt.Sprintf("%d Mi", memoryLimits/(1024*1024)),
		UsageRate: 0, // 需要metrics-server支持
		RequestRate: float64(memoryRequests) / float64(allocatableMemory) * 100,
	}

	// 存储监控信息（示例数据）
	monitoring.Storage = model.NodeResourceUsage{
		Total:     "100 Gi",
		Used:      "0 Gi",
		Available: "100 Gi",
		Requests:  "0 Gi",
		Limits:    "0 Gi",
		UsageRate: 0,
		RequestRate: 0,
	}

	// 网络监控信息（示例数据）
	monitoring.Network = model.NodeNetworkUsage{
		InboundBytes:    0,
		OutboundBytes:   0,
		InboundPackets:  0,
		OutboundPackets: 0,
	}

	return monitoring, nil
}

// calculateNodePodInfo 计算节点Pod信息
func (s *K8sNodesServiceImpl) calculateNodePodInfo(pods []corev1.Pod) (model.NodePodInfo, []model.NodePodDetail) {
	podInfo := model.NodePodInfo{
		TotalPods: len(pods),
	}

	var podList []model.NodePodDetail

	for _, pod := range pods {
		// 统计Pod状态
		switch pod.Status.Phase {
		case corev1.PodRunning:
			podInfo.RunningPods++
		case corev1.PodPending:
			podInfo.PendingPods++
		case corev1.PodFailed:
			podInfo.FailedPods++
		case corev1.PodSucceeded:
			podInfo.SucceededPods++
		}

		// 计算重启次数
		var restartCount int32
		for _, containerStatus := range pod.Status.ContainerStatuses {
			restartCount += containerStatus.RestartCount
		}

		// 计算资源请求和限制
		var cpuRequests, cpuLimits, memoryRequests, memoryLimits int64
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				cpuRequests += cpu.MilliValue()
			}
			if cpu, exists := container.Resources.Limits[corev1.ResourceCPU]; exists {
				cpuLimits += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				memoryRequests += memory.Value()
			}
			if memory, exists := container.Resources.Limits[corev1.ResourceMemory]; exists {
				memoryLimits += memory.Value()
			}
		}

		// 构建容器状态列表
		var containers []model.ContainerStatus
		for _, containerStatus := range pod.Status.ContainerStatuses {
			state := "Unknown"
			if containerStatus.State.Running != nil {
				state = "Running"
			} else if containerStatus.State.Waiting != nil {
				state = "Waiting"
			} else if containerStatus.State.Terminated != nil {
				state = "Terminated"
			}

			containers = append(containers, model.ContainerStatus{
				Name:         containerStatus.Name,
				Image:        containerStatus.Image,
				State:        state,
				Ready:        containerStatus.Ready,
				RestartCount: containerStatus.RestartCount,
			})
		}

		podDetail := model.NodePodDetail{
			Name:           pod.Name,
			Namespace:      pod.Namespace,
			Status:         string(pod.Status.Phase),
			Phase:          string(pod.Status.Phase),
			RestartCount:   restartCount,
			CreatedAt:      pod.CreationTimestamp.Format(time.RFC3339),
			Labels:         pod.Labels,
			CPURequests:    fmt.Sprintf("%.1f cores", float64(cpuRequests)/1000),
			CPULimits:      fmt.Sprintf("%.1f cores", float64(cpuLimits)/1000),
			MemoryRequests: fmt.Sprintf("%d Mi", memoryRequests/(1024*1024)),
			MemoryLimits:   fmt.Sprintf("%d Mi", memoryLimits/(1024*1024)),
			Containers:     containers,
		}

		podList = append(podList, podDetail)
	}

	return podInfo, podList
}

// calculateNodePodInfoWithMetrics 计算节点Pod信息（包含真实metrics）
func (s *K8sNodesServiceImpl) calculateNodePodInfoWithMetrics(pods []corev1.Pod, metricsClient *metricsclientset.Clientset) (model.NodePodInfo, []model.NodePodDetail) {
	podInfo := model.NodePodInfo{
		TotalPods: len(pods),
	}

	var podList []model.NodePodDetail

	for _, pod := range pods {
		// 统计Pod状态
		switch pod.Status.Phase {
		case corev1.PodRunning:
			podInfo.RunningPods++
		case corev1.PodPending:
			podInfo.PendingPods++
		case corev1.PodFailed:
			podInfo.FailedPods++
		case corev1.PodSucceeded:
			podInfo.SucceededPods++
		}

		// 计算重启次数
		var restartCount int32
		for _, containerStatus := range pod.Status.ContainerStatuses {
			restartCount += containerStatus.RestartCount
		}

		// 计算资源请求和限制
		var cpuRequests, cpuLimits, memoryRequests, memoryLimits int64
		for _, container := range pod.Spec.Containers {
			if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
				cpuRequests += cpu.MilliValue()
			}
			if cpu, exists := container.Resources.Limits[corev1.ResourceCPU]; exists {
				cpuLimits += cpu.MilliValue()
			}
			if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
				memoryRequests += memory.Value()
			}
			if memory, exists := container.Resources.Limits[corev1.ResourceMemory]; exists {
				memoryLimits += memory.Value()
			}
		}

		// 构建容器状态列表
		var containers []model.ContainerStatus
		for _, containerStatus := range pod.Status.ContainerStatuses {
			state := "Unknown"
			if containerStatus.State.Running != nil {
				state = "Running"
			} else if containerStatus.State.Waiting != nil {
				state = "Waiting"
			} else if containerStatus.State.Terminated != nil {
				state = "Terminated"
			}

			containers = append(containers, model.ContainerStatus{
				Name:         containerStatus.Name,
				Image:        containerStatus.Image,
				State:        state,
				Ready:        containerStatus.Ready,
				RestartCount: containerStatus.RestartCount,
			})
		}

		// 计算运行时间
		runningTime := "未知"
		if !pod.CreationTimestamp.Time.IsZero() {
			duration := time.Since(pod.CreationTimestamp.Time)
			if duration < time.Hour {
				runningTime = fmt.Sprintf("%.0fm", duration.Minutes())
			} else if duration < 24*time.Hour {
				runningTime = fmt.Sprintf("%.0fh%.0fm", duration.Hours(), duration.Minutes()-60*duration.Hours())
			} else {
				days := int(duration.Hours() / 24)
				hours := int(duration.Hours()) % 24
				runningTime = fmt.Sprintf("%dd%dh", days, hours)
			}
		}

		// 尝试获取真实使用数据
		cpuUsageStr := fmt.Sprintf("%.1f cores", float64(cpuRequests)/1000)
		memoryUsageStr := fmt.Sprintf("%d Mi", memoryRequests/(1024*1024))

		if metricsClient != nil {
			podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if err == nil && podMetrics != nil {
				var totalCPU, totalMemory resource.Quantity
				for _, containerMetric := range podMetrics.Containers {
					totalCPU.Add(containerMetric.Usage[corev1.ResourceCPU])
					totalMemory.Add(containerMetric.Usage[corev1.ResourceMemory])
				}

				if totalCPU.MilliValue() > 0 {
					cpuUsageStr = fmt.Sprintf("%.1f cores", float64(totalCPU.MilliValue())/1000)
				}
				if totalMemory.Value() > 0 {
					memoryUsageStr = fmt.Sprintf("%d Mi", totalMemory.Value()/(1024*1024))
				}
			}
		}

		podDetail := model.NodePodDetail{
			Name:           pod.Name,
			Namespace:      pod.Namespace,
			Status:         string(pod.Status.Phase),
			Phase:          string(pod.Status.Phase),
			RestartCount:   restartCount,
			CreatedAt:      pod.CreationTimestamp.Format(time.RFC3339),
			Labels:         pod.Labels,
			CPURequests:    fmt.Sprintf("%.1f cores", float64(cpuRequests)/1000),
			CPULimits:      fmt.Sprintf("%.1f cores", float64(cpuLimits)/1000),
			MemoryRequests: fmt.Sprintf("%d Mi", memoryRequests/(1024*1024)),
			MemoryLimits:   fmt.Sprintf("%d Mi", memoryLimits/(1024*1024)),
			Containers:     containers,
		}

		// 添加运行时间和实际使用量字段（这些字段可能需要在model中定义）
		if podDetail.Labels == nil {
			podDetail.Labels = make(map[string]string)
		}
		podDetail.Labels["__runtime"] = runningTime
		podDetail.Labels["__cpu_usage"] = cpuUsageStr
		podDetail.Labels["__memory_usage"] = memoryUsageStr

		podList = append(podList, podDetail)
	}

	return podInfo, podList
}