package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// IK8sServiceService K8s Service服务接口
type IK8sServiceService interface {
	// Service管理
	GetServices(c *gin.Context, clusterId uint, namespaceName string)
	GetServiceDetail(c *gin.Context, clusterId uint, namespaceName string, serviceName string)
	CreateService(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateServiceRequest)
	UpdateService(c *gin.Context, clusterId uint, namespaceName string, serviceName string, req *model.UpdateServiceRequest)
	DeleteService(c *gin.Context, clusterId uint, namespaceName string, serviceName string)
	GetServiceYaml(c *gin.Context, clusterId uint, namespaceName string, serviceName string)
	UpdateServiceYaml(c *gin.Context, clusterId uint, namespaceName string, serviceName string, yamlData map[string]interface{})
	GetServiceEvents(c *gin.Context, clusterId uint, namespaceName string, serviceName string)
}

// K8sServiceServiceImpl K8s Service服务实现
type K8sServiceServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

// NewK8sServiceService 创建K8s Service服务实例
func NewK8sServiceService(db *gorm.DB) IK8sServiceService {
	return &K8sServiceServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// getK8sClient 获取K8s客户端
func (s *K8sServiceServiceImpl) getK8sClient(clusterId uint) (*kubernetes.Clientset, error) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	if cluster.Credential == "" {
		return nil, fmt.Errorf("集群凭证为空")
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("解析集群凭证失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	return clientset, nil
}

// GetServices 获取Service列表
func (s *K8sServiceServiceImpl) GetServices(c *gin.Context, clusterId uint, namespaceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	services, err := clientset.CoreV1().Services(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Service列表失败: "+err.Error())
		return
	}

	var serviceList []model.K8sService
	for _, svc := range services.Items {
		k8sService := s.convertToK8sService(&svc, clientset, namespaceName)
		serviceList = append(serviceList, k8sService)
	}

	response := model.ServiceListResponse{
		Services: serviceList,
		Total:    len(serviceList),
	}

	result.Success(c, response)
}

// GetServiceDetail 获取Service详情
func (s *K8sServiceServiceImpl) GetServiceDetail(c *gin.Context, clusterId uint, namespaceName string, serviceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	svc, err := clientset.CoreV1().Services(namespaceName).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Service不存在: "+err.Error())
		return
	}

	// 获取Service详情
	k8sService := s.convertToK8sService(svc, clientset, namespaceName)

	// 获取相关事件
	events, _ := s.getServiceEvents(clientset, namespaceName, serviceName)

	// 获取关联的Pod列表
	pods := s.getServicePods(clientset, namespaceName, svc.Spec.Selector)

	serviceDetail := model.ServiceDetail{
		K8sService: k8sService,
		Events:     events,
		Pods:       pods,
		Spec:       svc.Spec,
	}

	result.Success(c, serviceDetail)
}

// CreateService 创建Service
func (s *K8sServiceServiceImpl) CreateService(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateServiceRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
			Labels:    req.Labels,
		},
		Spec: corev1.ServiceSpec{
			Type:        corev1.ServiceType(req.Type),
			Selector:    req.Selector,
			ExternalIPs: req.ExternalIPs,
		},
	}

	// 设置端口配置
	var ports []corev1.ServicePort
	for _, portSpec := range req.Ports {
		port := corev1.ServicePort{
			Name:     portSpec.Name,
			Protocol: corev1.Protocol(portSpec.Protocol),
			Port:     portSpec.Port,
		}

		if portSpec.TargetPort != "" {
			port.TargetPort = intstr.FromString(portSpec.TargetPort)
		} else {
			port.TargetPort = intstr.FromInt(int(portSpec.Port))
		}

		if portSpec.NodePort > 0 && req.Type == "NodePort" {
			port.NodePort = portSpec.NodePort
		}

		ports = append(ports, port)
	}
	service.Spec.Ports = ports

	// 创建Service
	createdService, err := clientset.CoreV1().Services(namespaceName).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建Service失败: "+err.Error())
		return
	}

	k8sService := s.convertToK8sService(createdService, clientset, namespaceName)
	result.Success(c, k8sService)
}

// UpdateService 更新Service
func (s *K8sServiceServiceImpl) UpdateService(c *gin.Context, clusterId uint, namespaceName string, serviceName string, req *model.UpdateServiceRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有Service
	existingService, err := clientset.CoreV1().Services(namespaceName).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Service不存在: "+err.Error())
		return
	}

	// 更新Service属性
	if req.Labels != nil {
		existingService.Labels = req.Labels
	}
	if req.Type != "" {
		existingService.Spec.Type = corev1.ServiceType(req.Type)
	}
	if req.Selector != nil {
		existingService.Spec.Selector = req.Selector
	}
	if req.ExternalIPs != nil {
		existingService.Spec.ExternalIPs = req.ExternalIPs
	}

	// 更新端口配置
	if req.Ports != nil {
		var ports []corev1.ServicePort
		for _, portSpec := range req.Ports {
			port := corev1.ServicePort{
				Name:     portSpec.Name,
				Protocol: corev1.Protocol(portSpec.Protocol),
				Port:     portSpec.Port,
			}

			if portSpec.TargetPort != "" {
				port.TargetPort = intstr.FromString(portSpec.TargetPort)
			} else {
				port.TargetPort = intstr.FromInt(int(portSpec.Port))
			}

			if portSpec.NodePort > 0 && req.Type == "NodePort" {
				port.NodePort = portSpec.NodePort
			}

			ports = append(ports, port)
		}
		existingService.Spec.Ports = ports
	}

	// 更新Service
	updatedService, err := clientset.CoreV1().Services(namespaceName).Update(context.TODO(), existingService, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新Service失败: "+err.Error())
		return
	}

	k8sService := s.convertToK8sService(updatedService, clientset, namespaceName)
	result.Success(c, k8sService)
}

// DeleteService 删除Service
func (s *K8sServiceServiceImpl) DeleteService(c *gin.Context, clusterId uint, namespaceName string, serviceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.CoreV1().Services(namespaceName).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除Service失败: "+err.Error())
		return
	}

	result.Success(c, "Service删除成功")
}

// GetServiceYaml 获取Service的YAML配置
func (s *K8sServiceServiceImpl) GetServiceYaml(c *gin.Context, clusterId uint, namespaceName string, serviceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	svc, err := clientset.CoreV1().Services(namespaceName).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Service不存在: "+err.Error())
		return
	}

	// 清理不需要的字段，但保留 status
	svc.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if svc.APIVersion == "" {
		svc.APIVersion = "v1"
	}
	if svc.Kind == "" {
		svc.Kind = "Service"
	}

	// 清理元数据字段
	svc.ResourceVersion = ""
	svc.UID = ""
	svc.Generation = 0
	svc.CreationTimestamp = metav1.Time{}
	svc.SelfLink = ""

	yamlData, err := yaml.Marshal(svc)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "转换YAML失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": string(yamlData),
	})
}

// UpdateServiceYaml 通过YAML更新Service
func (s *K8sServiceServiceImpl) UpdateServiceYaml(c *gin.Context, clusterId uint, namespaceName string, serviceName string, yamlData map[string]interface{}) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 首先获取当前的 Service 以获取 resourceVersion 和 clusterIP
	existingService, err := clientset.CoreV1().Services(namespaceName).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Service不存在: "+err.Error())
		return
	}

	var yamlBytes []byte

	// 检查是否是 {yaml: "..."} 格式
	if yamlStr, ok := yamlData["yaml"].(string); ok {
		// 前端发送的是 YAML 字符串
		yamlBytes = []byte(yamlStr)
	} else {
		// 前端发送的是 JSON 对象
		yamlBytes, err = json.Marshal(yamlData)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, "YAML数据格式错误: "+err.Error())
			return
		}
	}

	// 解析 YAML 到 Service 对象
	var service corev1.Service
	err = yaml.Unmarshal(yamlBytes, &service)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}

	// 确保名称和命名空间正确
	service.Name = serviceName
	service.Namespace = namespaceName

	// 保留 resourceVersion 和 clusterIP（Service 的 clusterIP 是不可变的）
	service.ResourceVersion = existingService.ResourceVersion
	service.Spec.ClusterIP = existingService.Spec.ClusterIP
	service.Spec.ClusterIPs = existingService.Spec.ClusterIPs

	// 更新Service
	_, err = clientset.CoreV1().Services(namespaceName).Update(context.TODO(), &service, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新Service失败: "+err.Error())
		return
	}

	result.Success(c, "Service YAML更新成功")
}

// GetServiceEvents 获取Service事件
func (s *K8sServiceServiceImpl) GetServiceEvents(c *gin.Context, clusterId uint, namespaceName string, serviceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	events, err := s.getServiceEvents(clientset, namespaceName, serviceName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Service事件失败: "+err.Error())
		return
	}

	result.Success(c, events)
}

// convertToK8sService 转换Service对象
func (s *K8sServiceServiceImpl) convertToK8sService(svc *corev1.Service, clientset *kubernetes.Clientset, namespaceName string) model.K8sService {
	// 转换端口配置
	var ports []model.ServicePort
	for _, port := range svc.Spec.Ports {
		servicePort := model.ServicePort{
			Name:       port.Name,
			Protocol:   string(port.Protocol),
			Port:       port.Port,
			TargetPort: port.TargetPort.String(),
			NodePort:   port.NodePort,
		}
		ports = append(ports, servicePort)
	}

	// 获取端点信息
	endpoints := s.getServiceEndpoints(clientset, namespaceName, svc.Name)

	// 确定状态
	status := "Active"
	if len(endpoints) == 0 {
		status = "No Endpoints"
	}

	return model.K8sService{
		Name:        svc.Name,
		Namespace:   svc.Namespace,
		Labels:      svc.Labels,
		Type:        string(svc.Spec.Type),
		Selector:    svc.Spec.Selector,
		ClusterIP:   svc.Spec.ClusterIP,
		ExternalIPs: svc.Spec.ExternalIPs,
		Ports:       ports,
		Endpoints:   endpoints,
		CreatedAt:   svc.CreationTimestamp.Format(time.RFC3339),
		Status:      status,
	}
}

// getServiceEndpoints 获取Service端点信息
func (s *K8sServiceServiceImpl) getServiceEndpoints(clientset *kubernetes.Clientset, namespaceName string, serviceName string) []model.ServiceEndpoint {
	var endpoints []model.ServiceEndpoint

	eps, err := clientset.CoreV1().Endpoints(namespaceName).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return endpoints
	}

	for _, subset := range eps.Subsets {
		// 处理就绪的地址
		for _, addr := range subset.Addresses {
			endpoint := model.ServiceEndpoint{
				IP:       addr.IP,
				Hostname: addr.Hostname,
				Ready:    true,
			}
			if addr.NodeName != nil {
				endpoint.NodeName = *addr.NodeName
			}

			// 添加端口信息
			var endpointPorts []model.EndpointPort
			for _, port := range subset.Ports {
				endpointPort := model.EndpointPort{
					Name:     port.Name,
					Port:     port.Port,
					Protocol: string(port.Protocol),
				}
				endpointPorts = append(endpointPorts, endpointPort)
			}
			endpoint.Ports = endpointPorts

			endpoints = append(endpoints, endpoint)
		}

		// 处理未就绪的地址
		for _, addr := range subset.NotReadyAddresses {
			endpoint := model.ServiceEndpoint{
				IP:       addr.IP,
				Hostname: addr.Hostname,
				Ready:    false,
			}
			if addr.NodeName != nil {
				endpoint.NodeName = *addr.NodeName
			}

			// 添加端口信息
			var endpointPorts []model.EndpointPort
			for _, port := range subset.Ports {
				endpointPort := model.EndpointPort{
					Name:     port.Name,
					Port:     port.Port,
					Protocol: string(port.Protocol),
				}
				endpointPorts = append(endpointPorts, endpointPort)
			}
			endpoint.Ports = endpointPorts

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// getServiceEvents 获取Service相关事件
func (s *K8sServiceServiceImpl) getServiceEvents(clientset *kubernetes.Clientset, namespaceName string, serviceName string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Service", serviceName),
	})
	if err != nil {
		return nil, err
	}

	var k8sEvents []model.K8sEvent
	for _, event := range events.Items {
		k8sEvent := model.K8sEvent{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			Source:    event.Source.Component,
			Count:     event.Count,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
		}
		k8sEvents = append(k8sEvents, k8sEvent)
	}

	return k8sEvents, nil
}

// getServicePods 获取Service关联的Pod列表
func (s *K8sServiceServiceImpl) getServicePods(clientset *kubernetes.Clientset, namespaceName string, selector map[string]string) []model.K8sPodInfo {
	var pods []model.K8sPodInfo

	if len(selector) == 0 {
		return pods
	}

	// 构建标签选择器
	var selectorParts []string
	for key, value := range selector {
		selectorParts = append(selectorParts, fmt.Sprintf("%s=%s", key, value))
	}
	labelSelector := fmt.Sprintf("%s", selectorParts[0])
	for i := 1; i < len(selectorParts); i++ {
		labelSelector += "," + selectorParts[i]
	}

	podList, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return pods
	}

	for _, pod := range podList.Items {
		podInfo := model.K8sPodInfo{
			Name:         pod.Name,
			Status:       string(pod.Status.Phase),
			Phase:        string(pod.Status.Phase),
			RestartCount: 0,
			NodeName:     pod.Spec.NodeName,
			PodIP:        pod.Status.PodIP,
			HostIP:       pod.Status.HostIP,
			CreatedAt:    pod.CreationTimestamp.Format(time.RFC3339),
			Labels:       pod.Labels,
		}

		// 计算重启次数
		for _, containerStatus := range pod.Status.ContainerStatuses {
			podInfo.RestartCount += containerStatus.RestartCount
		}

		// 计算运行时间
		if pod.Status.StartTime != nil {
			duration := time.Since(pod.Status.StartTime.Time)
			podInfo.RunningTime = duration.String()
		}

		pods = append(pods, podInfo)
	}

	return pods
}