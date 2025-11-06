package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// IK8sIngressService K8s Ingress服务接口
type IK8sIngressService interface {
	// Ingress管理
	GetIngresses(c *gin.Context, clusterId uint, namespaceName string)
	GetIngressDetail(c *gin.Context, clusterId uint, namespaceName string, ingressName string)
	CreateIngress(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateIngressRequest)
	UpdateIngress(c *gin.Context, clusterId uint, namespaceName string, ingressName string, req *model.UpdateIngressRequest)
	DeleteIngress(c *gin.Context, clusterId uint, namespaceName string, ingressName string)
	GetIngressYaml(c *gin.Context, clusterId uint, namespaceName string, ingressName string)
	UpdateIngressYaml(c *gin.Context, clusterId uint, namespaceName string, ingressName string, yamlData map[string]interface{})
	GetIngressEvents(c *gin.Context, clusterId uint, namespaceName string, ingressName string)
	GetIngressMonitoring(c *gin.Context, clusterId uint, namespaceName string, ingressName string)
}

// K8sIngressServiceImpl K8s Ingress服务实现
type K8sIngressServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

// NewK8sIngressService 创建K8s Ingress服务实例
func NewK8sIngressService(db *gorm.DB) IK8sIngressService {
	return &K8sIngressServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// getK8sClient 获取K8s客户端
func (s *K8sIngressServiceImpl) getK8sClient(clusterId uint) (*kubernetes.Clientset, error) {
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

// GetIngresses 获取Ingress列表
func (s *K8sIngressServiceImpl) GetIngresses(c *gin.Context, clusterId uint, namespaceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	ingresses, err := clientset.NetworkingV1().Ingresses(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Ingress列表失败: "+err.Error())
		return
	}

	var ingressList []model.K8sIngress

	// 列表查询时使用轻量级转换，不查询 IngressClass 和 Controller 详细信息
	// 这样可以避免性能问题，详细信息在查看单个 Ingress 详情时再获取
	for _, ing := range ingresses.Items {
		k8sIngress := s.convertToK8sIngress(&ing)
		ingressList = append(ingressList, k8sIngress)
	}

	response := model.IngressListResponse{
		Ingresses: ingressList,
		Total:     len(ingressList),
	}

	result.Success(c, response)
}

// GetIngressDetail 获取Ingress详情
func (s *K8sIngressServiceImpl) GetIngressDetail(c *gin.Context, clusterId uint, namespaceName string, ingressName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	ing, err := clientset.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Ingress不存在: "+err.Error())
		return
	}

	// 获取Ingress详情（使用完整转换，包含 IngressClass 和 Controller 信息）
	k8sIngress := s.convertToK8sIngressWithClient(clientset, ing)

	// 获取相关事件
	events, _ := s.getIngressEvents(clientset, namespaceName, ingressName)

	ingressDetail := model.IngressDetail{
		K8sIngress: k8sIngress,
		Events:     events,
		Spec:       ing.Spec,
	}

	result.Success(c, ingressDetail)
}

// CreateIngress 创建Ingress
func (s *K8sIngressServiceImpl) CreateIngress(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateIngressRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建Ingress对象
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   namespaceName,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: networkingv1.IngressSpec{},
	}

	// 设置IngressClass
	if req.Class != "" {
		ingress.Spec.IngressClassName = &req.Class
	}

	// 构建规则
	var rules []networkingv1.IngressRule
	for _, ruleSpec := range req.Rules {
		rule := networkingv1.IngressRule{
			Host: ruleSpec.Host,
		}

		// 构建HTTP路径
		var paths []networkingv1.HTTPIngressPath
		for _, pathSpec := range ruleSpec.Paths {
			pathType := networkingv1.PathType(pathSpec.PathType)
			path := networkingv1.HTTPIngressPath{
				Path:     pathSpec.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: pathSpec.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: pathSpec.ServicePort,
						},
					},
				},
			}
			paths = append(paths, path)
		}

		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: paths,
		}
		rules = append(rules, rule)
	}
	ingress.Spec.Rules = rules

	// 构建TLS配置
	if len(req.TLS) > 0 {
		var tlsConfigs []networkingv1.IngressTLS
		for _, tlsSpec := range req.TLS {
			tlsConfig := networkingv1.IngressTLS{
				Hosts:      tlsSpec.Hosts,
				SecretName: tlsSpec.SecretName,
			}
			tlsConfigs = append(tlsConfigs, tlsConfig)
		}
		ingress.Spec.TLS = tlsConfigs
	}

	// 创建Ingress
	createdIngress, err := clientset.NetworkingV1().Ingresses(namespaceName).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建Ingress失败: "+err.Error())
		return
	}

	k8sIngress := s.convertToK8sIngress(createdIngress)
	result.Success(c, k8sIngress)
}

// UpdateIngress 更新Ingress
func (s *K8sIngressServiceImpl) UpdateIngress(c *gin.Context, clusterId uint, namespaceName string, ingressName string, req *model.UpdateIngressRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有Ingress
	existingIngress, err := clientset.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Ingress不存在: "+err.Error())
		return
	}

	// 更新Ingress属性
	if req.Labels != nil {
		existingIngress.Labels = req.Labels
	}
	if req.Annotations != nil {
		existingIngress.Annotations = req.Annotations
	}
	if req.Class != "" {
		existingIngress.Spec.IngressClassName = &req.Class
	}

	// 更新规则
	if req.Rules != nil {
		var rules []networkingv1.IngressRule
		for _, ruleSpec := range req.Rules {
			rule := networkingv1.IngressRule{
				Host: ruleSpec.Host,
			}

			// 构建HTTP路径
			var paths []networkingv1.HTTPIngressPath
			for _, pathSpec := range ruleSpec.Paths {
				pathType := networkingv1.PathType(pathSpec.PathType)
				path := networkingv1.HTTPIngressPath{
					Path:     pathSpec.Path,
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: pathSpec.ServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: pathSpec.ServicePort,
							},
						},
					},
				}
				paths = append(paths, path)
			}

			rule.HTTP = &networkingv1.HTTPIngressRuleValue{
				Paths: paths,
			}
			rules = append(rules, rule)
		}
		existingIngress.Spec.Rules = rules
	}

	// 更新TLS配置
	if req.TLS != nil {
		var tlsConfigs []networkingv1.IngressTLS
		for _, tlsSpec := range req.TLS {
			tlsConfig := networkingv1.IngressTLS{
				Hosts:      tlsSpec.Hosts,
				SecretName: tlsSpec.SecretName,
			}
			tlsConfigs = append(tlsConfigs, tlsConfig)
		}
		existingIngress.Spec.TLS = tlsConfigs
	}

	// 更新Ingress
	updatedIngress, err := clientset.NetworkingV1().Ingresses(namespaceName).Update(context.TODO(), existingIngress, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新Ingress失败: "+err.Error())
		return
	}

	k8sIngress := s.convertToK8sIngress(updatedIngress)
	result.Success(c, k8sIngress)
}

// DeleteIngress 删除Ingress
func (s *K8sIngressServiceImpl) DeleteIngress(c *gin.Context, clusterId uint, namespaceName string, ingressName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.NetworkingV1().Ingresses(namespaceName).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除Ingress失败: "+err.Error())
		return
	}

	result.Success(c, "Ingress删除成功")
}

// GetIngressYaml 获取Ingress的YAML配置
func (s *K8sIngressServiceImpl) GetIngressYaml(c *gin.Context, clusterId uint, namespaceName string, ingressName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	ing, err := clientset.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Ingress不存在: "+err.Error())
		return
	}

	// 清理不需要的字段，但保留 status
	ing.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if ing.APIVersion == "" {
		ing.APIVersion = "networking.k8s.io/v1"
	}
	if ing.Kind == "" {
		ing.Kind = "Ingress"
	}

	// 清理元数据字段
	ing.ResourceVersion = ""
	ing.UID = ""
	ing.Generation = 0
	ing.CreationTimestamp = metav1.Time{}
	ing.SelfLink = ""

	yamlData, err := yaml.Marshal(ing)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "转换YAML失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": string(yamlData),
	})
}

// UpdateIngressYaml 通过YAML更新Ingress
func (s *K8sIngressServiceImpl) UpdateIngressYaml(c *gin.Context, clusterId uint, namespaceName string, ingressName string, yamlData map[string]interface{}) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 首先获取当前的 Ingress 以获取 resourceVersion
	existingIngress, err := clientset.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Ingress不存在: "+err.Error())
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

	// 处理API版本兼容性
	yamlBytes, err = s.convertIngressAPIVersion(yamlBytes)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "解码Kubernetes对象失败: "+err.Error())
		return
	}

	// 解析 YAML 到 Ingress 对象
	var ingress networkingv1.Ingress
	err = yaml.Unmarshal(yamlBytes, &ingress)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}

	// 验证必填字段
	if len(ingress.Spec.Rules) == 0 && ingress.Spec.DefaultBackend == nil {
		result.Failed(c, http.StatusBadRequest, "Ingress必须指定rules或defaultBackend")
		return
	}

	// 确保名称和命名空间正确
	ingress.Name = ingressName
	ingress.Namespace = namespaceName

	// 保留 resourceVersion 用于更新
	ingress.ResourceVersion = existingIngress.ResourceVersion

	// 更新Ingress
	_, err = clientset.NetworkingV1().Ingresses(namespaceName).Update(context.TODO(), &ingress, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新Ingress失败: "+err.Error())
		return
	}

	result.Success(c, "Ingress YAML更新成功")
}

// GetIngressEvents 获取Ingress事件
func (s *K8sIngressServiceImpl) GetIngressEvents(c *gin.Context, clusterId uint, namespaceName string, ingressName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	events, err := s.getIngressEvents(clientset, namespaceName, ingressName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Ingress事件失败: "+err.Error())
		return
	}

	result.Success(c, events)
}

// GetIngressMonitoring 获取Ingress监控信息
func (s *K8sIngressServiceImpl) GetIngressMonitoring(c *gin.Context, clusterId uint, namespaceName string, ingressName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	ing, err := clientset.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Ingress不存在: "+err.Error())
		return
	}

	monitoring := map[string]interface{}{
		"name":      ing.Name,
		"namespace": ing.Namespace,
		"status":    s.getIngressStatus(ing),
		"endpoints": s.getIngressEndpoints(ing),
		"rules":     len(ing.Spec.Rules),
		"tls":       len(ing.Spec.TLS),
		"createdAt": ing.CreationTimestamp.Format(time.RFC3339),
	}

	result.Success(c, monitoring)
}

// convertToK8sIngress 转换Ingress对象
func (s *K8sIngressServiceImpl) convertToK8sIngress(ing *networkingv1.Ingress) model.K8sIngress {
	// 转换规则
	var rules []model.IngressRule
	for _, rule := range ing.Spec.Rules {
		ingressRule := model.IngressRule{
			Host: rule.Host,
		}

		if rule.HTTP != nil {
			var paths []model.IngressPath
			for _, path := range rule.HTTP.Paths {
				ingressPath := model.IngressPath{
					Path:     path.Path,
					PathType: string(*path.PathType),
				}

				if path.Backend.Service != nil {
					ingressPath.Backend = model.IngressBackend{
						Service: model.IngressServiceBackend{
							Name: path.Backend.Service.Name,
							Port: model.IngressServicePort{
								Number: path.Backend.Service.Port.Number,
								Name:   path.Backend.Service.Port.Name,
							},
						},
					}
				}

				paths = append(paths, ingressPath)
			}

			ingressRule.HTTP = model.IngressRuleValue{
				Paths: paths,
			}
		}

		rules = append(rules, ingressRule)
	}

	// 转换TLS配置
	var tlsConfigs []model.IngressTLS
	for _, tls := range ing.Spec.TLS {
		tlsConfig := model.IngressTLS{
			Hosts:      tls.Hosts,
			SecretName: tls.SecretName,
		}
		tlsConfigs = append(tlsConfigs, tlsConfig)
	}

	// 获取IngressClass
	var ingressClass string
	if ing.Spec.IngressClassName != nil {
		ingressClass = *ing.Spec.IngressClassName
	} else if ing.Annotations != nil {
		// 兼容旧版本从 annotation 获取
		if class, ok := ing.Annotations["kubernetes.io/ingress.class"]; ok {
			ingressClass = class
		}
	}

	// 构建负载均衡器信息
	var loadBalancer model.IngressLoadBalancer
	var ingressIngresses []model.IngressLoadBalancerIngress
	for _, lbIngress := range ing.Status.LoadBalancer.Ingress {
		lbIngressInfo := model.IngressLoadBalancerIngress{
			IP:       lbIngress.IP,
			Hostname: lbIngress.Hostname,
		}

		var lbPorts []model.IngressLoadBalancerPort
		for _, port := range lbIngress.Ports {
			lbPort := model.IngressLoadBalancerPort{
				Port:     port.Port,
				Protocol: string(port.Protocol),
			}
			lbPorts = append(lbPorts, lbPort)
		}
		lbIngressInfo.Ports = lbPorts

		ingressIngresses = append(ingressIngresses, lbIngressInfo)
	}
	loadBalancer.Ingress = ingressIngresses

	// 生成访问端点
	endpoints := s.getIngressEndpoints(ing)

	// 获取 IngressClass 详细信息、Controller 版本和类型
	controllerName := ""
	controllerVersion := ""
	ingressType := ""

	if ingressClass != "" {
		// 获取 IngressClass 信息
		controllerName, controllerVersion = s.getIngressClassInfo(ing)
		// 判断类型
		ingressType = s.getIngressType(ing, ingressClass)
	}

	return model.K8sIngress{
		Name:              ing.Name,
		Namespace:         ing.Namespace,
		Labels:            ing.Labels,
		Class:             ingressClass,
		ControllerName:    controllerName,
		ControllerVersion: controllerVersion,
		Type:              ingressType,
		Rules:             rules,
		TLS:               tlsConfigs,
		LoadBalancer:      loadBalancer,
		Endpoints:         endpoints,
		CreatedAt:         ing.CreationTimestamp.Format(time.RFC3339),
		Status:            s.getIngressStatus(ing),
	}
}

// convertToK8sIngressWithClient 使用 clientset 转换Ingress对象，获取更详细的信息
func (s *K8sIngressServiceImpl) convertToK8sIngressWithClient(clientset *kubernetes.Clientset, ing *networkingv1.Ingress) model.K8sIngress {
	// 转换规则
	var rules []model.IngressRule
	for _, rule := range ing.Spec.Rules {
		ingressRule := model.IngressRule{
			Host: rule.Host,
		}

		if rule.HTTP != nil {
			var paths []model.IngressPath
			for _, path := range rule.HTTP.Paths {
				ingressPath := model.IngressPath{
					Path:     path.Path,
					PathType: string(*path.PathType),
				}

				if path.Backend.Service != nil {
					ingressPath.Backend = model.IngressBackend{
						Service: model.IngressServiceBackend{
							Name: path.Backend.Service.Name,
							Port: model.IngressServicePort{
								Number: path.Backend.Service.Port.Number,
								Name:   path.Backend.Service.Port.Name,
							},
						},
					}
				}

				paths = append(paths, ingressPath)
			}

			ingressRule.HTTP = model.IngressRuleValue{
				Paths: paths,
			}
		}

		rules = append(rules, ingressRule)
	}

	// 转换TLS配置
	var tlsConfigs []model.IngressTLS
	for _, tls := range ing.Spec.TLS {
		tlsConfig := model.IngressTLS{
			Hosts:      tls.Hosts,
			SecretName: tls.SecretName,
		}
		tlsConfigs = append(tlsConfigs, tlsConfig)
	}

	// 获取IngressClass
	var ingressClass string
	if ing.Spec.IngressClassName != nil {
		ingressClass = *ing.Spec.IngressClassName
	} else if ing.Annotations != nil {
		// 兼容旧版本从 annotation 获取
		if class, ok := ing.Annotations["kubernetes.io/ingress.class"]; ok {
			ingressClass = class
		}
	}

	// 构建负载均衡器信息
	var loadBalancer model.IngressLoadBalancer
	var ingressIngresses []model.IngressLoadBalancerIngress
	for _, lbIngress := range ing.Status.LoadBalancer.Ingress {
		lbIngressInfo := model.IngressLoadBalancerIngress{
			IP:       lbIngress.IP,
			Hostname: lbIngress.Hostname,
		}

		var lbPorts []model.IngressLoadBalancerPort
		for _, port := range lbIngress.Ports {
			lbPort := model.IngressLoadBalancerPort{
				Port:     port.Port,
				Protocol: string(port.Protocol),
			}
			lbPorts = append(lbPorts, lbPort)
		}
		lbIngressInfo.Ports = lbPorts

		ingressIngresses = append(ingressIngresses, lbIngressInfo)
	}
	loadBalancer.Ingress = ingressIngresses

	// 生成访问端点
	endpoints := s.getIngressEndpoints(ing)

	// 获取 IngressClass 详细信息、Controller 版本和类型（使用 clientset）
	controllerName := ""
	controllerVersion := ""
	ingressType := ""

	if ingressClass != "" {
		// 使用 clientset 获取更详细的信息
		controllerName, controllerVersion = s.getIngressClassInfoWithClient(clientset, ing)
		// 判断类型
		ingressType = s.getIngressType(ing, ingressClass)
	}

	return model.K8sIngress{
		Name:              ing.Name,
		Namespace:         ing.Namespace,
		Labels:            ing.Labels,
		Class:             ingressClass,
		ControllerName:    controllerName,
		ControllerVersion: controllerVersion,
		Type:              ingressType,
		Rules:             rules,
		TLS:               tlsConfigs,
		LoadBalancer:      loadBalancer,
		Endpoints:         endpoints,
		CreatedAt:         ing.CreationTimestamp.Format(time.RFC3339),
		Status:            s.getIngressStatus(ing),
	}
}

// getIngressStatus 获取Ingress状态
func (s *K8sIngressServiceImpl) getIngressStatus(ing *networkingv1.Ingress) string {
	if len(ing.Status.LoadBalancer.Ingress) > 0 {
		return "Active"
	}
	return "Pending"
}

// getIngressEndpoints 获取Ingress访问端点
func (s *K8sIngressServiceImpl) getIngressEndpoints(ing *networkingv1.Ingress) []string {
	var endpoints []string

	// 从LoadBalancer状态获取端点
	for _, lbIngress := range ing.Status.LoadBalancer.Ingress {
		if lbIngress.IP != "" {
			for _, rule := range ing.Spec.Rules {
				if rule.Host != "" {
					protocol := "http"
					if len(ing.Spec.TLS) > 0 {
						for _, tls := range ing.Spec.TLS {
							for _, host := range tls.Hosts {
								if host == rule.Host {
									protocol = "https"
									break
								}
							}
							if protocol == "https" {
								break
							}
						}
					}
					endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, rule.Host))
				}
			}
		}
		if lbIngress.Hostname != "" {
			for _, rule := range ing.Spec.Rules {
				if rule.Host != "" {
					protocol := "http"
					if len(ing.Spec.TLS) > 0 {
						for _, tls := range ing.Spec.TLS {
							for _, host := range tls.Hosts {
								if host == rule.Host {
									protocol = "https"
									break
								}
							}
							if protocol == "https" {
								break
							}
						}
					}
					endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, rule.Host))
				}
			}
		}
	}

	// 如果没有LoadBalancer状态，使用规则中的主机名
	if len(endpoints) == 0 {
		for _, rule := range ing.Spec.Rules {
			if rule.Host != "" {
				protocol := "http"
				if len(ing.Spec.TLS) > 0 {
					for _, tls := range ing.Spec.TLS {
						for _, host := range tls.Hosts {
							if host == rule.Host {
								protocol = "https"
								break
							}
						}
						if protocol == "https" {
							break
						}
					}
				}
				endpoints = append(endpoints, fmt.Sprintf("%s://%s", protocol, rule.Host))
			}
		}
	}

	return endpoints
}

// getIngressEvents 获取Ingress相关事件
func (s *K8sIngressServiceImpl) getIngressEvents(clientset *kubernetes.Clientset, namespaceName string, ingressName string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Ingress", ingressName),
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

// convertIngressAPIVersion 转换Ingress API版本从v1beta1到v1
func (s *K8sIngressServiceImpl) convertIngressAPIVersion(yamlBytes []byte) ([]byte, error) {
	// 解析YAML为通用对象
	var obj map[string]interface{}
	if err := yaml.Unmarshal(yamlBytes, &obj); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %v", err)
	}

	// 检查API版本
	apiVersion, ok := obj["apiVersion"].(string)
	if !ok {
		return yamlBytes, nil // 如果没有apiVersion字段，直接返回
	}

	// 如果是v1beta1版本，需要转换到v1
	if apiVersion == "networking.k8s.io/v1beta1" {
		obj["apiVersion"] = "networking.k8s.io/v1"

		// 转换spec.rules[].http.paths[].backend格式
		if spec, ok := obj["spec"].(map[string]interface{}); ok {
			if rules, ok := spec["rules"].([]interface{}); ok {
				for _, rule := range rules {
					if ruleMap, ok := rule.(map[string]interface{}); ok {
						if http, ok := ruleMap["http"].(map[string]interface{}); ok {
							if paths, ok := http["paths"].([]interface{}); ok {
								for _, path := range paths {
									if pathMap, ok := path.(map[string]interface{}); ok {
										if backend, ok := pathMap["backend"].(map[string]interface{}); ok {
											// v1beta1格式: backend.serviceName, backend.servicePort
											// v1格式: backend.service.name, backend.service.port.number
											if serviceName, hasServiceName := backend["serviceName"]; hasServiceName {
												if servicePort, hasServicePort := backend["servicePort"]; hasServicePort {
													// 删除旧格式字段
													delete(backend, "serviceName")
													delete(backend, "servicePort")

													// 添加新格式字段
													backend["service"] = map[string]interface{}{
														"name": serviceName,
														"port": map[string]interface{}{
															"number": servicePort,
														},
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// 重新序列化为YAML
	convertedBytes, err := yaml.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("序列化转换后的YAML失败: %v", err)
	}

	return convertedBytes, nil
}

// getIngressClassInfo 获取 IngressClass 信息和 Controller 版本
func (s *K8sIngressServiceImpl) getIngressClassInfo(ing *networkingv1.Ingress) (controllerName string, version string) {
	// 需要从 Ingress 获取 namespace，才能获取对应的 clientset
	// 由于这个方法在 convertToK8sIngress 中调用，此时没有 clientset
	// 我们需要重新获取 clientset，但这会影响性能
	// 暂时返回空字符串，后续可以优化

	// 从 Ingress 的 annotations 中尝试获取 controller 信息
	if ing.Annotations != nil {
		// 某些 controller 会在 annotation 中标记自己的信息
		if ctrl, ok := ing.Annotations["meta.helm.sh/release-name"]; ok {
			controllerName = ctrl
		}
	}

	// 尝试从 IngressClass 名称推断 controller 类型
	if ing.Spec.IngressClassName != nil {
		className := *ing.Spec.IngressClassName
		if strings.Contains(strings.ToLower(className), "nginx") {
			controllerName = "nginx-ingress"
		} else if strings.Contains(strings.ToLower(className), "traefik") {
			controllerName = "traefik"
		}
	}

	return controllerName, ""
}

// getIngressClassInfoWithClient 使用 clientset 获取详细的 IngressClass 信息和 Controller 版本
func (s *K8sIngressServiceImpl) getIngressClassInfoWithClient(clientset *kubernetes.Clientset, ing *networkingv1.Ingress) (controllerName string, version string) {
	controllerName = ""
	version = ""

	// 获取 IngressClass 名称
	var ingressClassName string
	if ing.Spec.IngressClassName != nil {
		ingressClassName = *ing.Spec.IngressClassName
	} else if ing.Annotations != nil {
		if class, ok := ing.Annotations["kubernetes.io/ingress.class"]; ok {
			ingressClassName = class
		}
	}

	// 如果没有 IngressClass，使用基础推断
	if ingressClassName == "" {
		controllerName, _ = s.getIngressClassInfo(ing)
		return
	}

	// 1. 尝试获取 IngressClass 资源
	ingressClass, err := clientset.NetworkingV1().IngressClasses().Get(context.TODO(), ingressClassName, metav1.GetOptions{})
	if err == nil {
		// 成功获取 IngressClass，提取 controller 名称
		controllerName = ingressClass.Spec.Controller
	} else {
		// 无法获取 IngressClass，从名称推断
		if strings.Contains(strings.ToLower(ingressClassName), "nginx") {
			controllerName = "k8s.io/ingress-nginx"
		} else if strings.Contains(strings.ToLower(ingressClassName), "traefik") {
			controllerName = "traefik.io/ingress-controller"
		} else {
			// 使用 IngressClass 名称作为 controller 名称
			controllerName = ingressClassName
		}
	}

	// 2. 尝试获取 controller 版本
	version = s.findControllerVersion(clientset, ingressClassName, controllerName)

	// 3. 如果还是找不到版本，尝试从 Ingress 的 annotations 获取
	if version == "" && ing.Annotations != nil {
		// 某些 controller 会在 annotation 中记录版本信息
		if ver, ok := ing.Annotations["ingress.kubernetes.io/controller-version"]; ok {
			version = ver
		}
	}

	return
}

// findControllerVersion 查找 Controller 版本
func (s *K8sIngressServiceImpl) findControllerVersion(clientset *kubernetes.Clientset, ingressClassName string, controllerName string) string {
	// 根据 controller 类型确定查找策略
	var labelSelectors []string
	var namespaces []string

	lowerControllerName := strings.ToLower(controllerName)
	lowerClassName := strings.ToLower(ingressClassName)

	// 根据 controller 类型配置查找参数
	if strings.Contains(lowerControllerName, "nginx") || strings.Contains(lowerClassName, "nginx") {
		labelSelectors = []string{
			"app.kubernetes.io/name=ingress-nginx",
			"app=ingress-nginx",
			"app.kubernetes.io/component=controller",
		}
		namespaces = []string{"ingress-nginx", "kube-system", "default", "nginx-ingress"}
	} else if strings.Contains(lowerControllerName, "traefik") || strings.Contains(lowerClassName, "traefik") {
		labelSelectors = []string{
			"app.kubernetes.io/name=traefik",
			"app=traefik",
		}
		namespaces = []string{"traefik", "kube-system", "default"}
	} else {
		// 通用查找
		labelSelectors = []string{
			"app.kubernetes.io/component=controller",
			fmt.Sprintf("app.kubernetes.io/name=%s", ingressClassName),
		}
		namespaces = []string{"kube-system", "default", ingressClassName}
	}

	// 尝试各种组合查找 Controller Pod
	for _, ns := range namespaces {
		for _, labelSelector := range labelSelectors {
			pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
				LabelSelector: labelSelector,
				Limit:         1, // 只需要一个 Pod
			})

			if err == nil && len(pods.Items) > 0 {
				version := s.extractVersionFromPod(&pods.Items[0])
				if version != "" {
					return version
				}
			}
		}
	}

	return ""
}

// extractVersionFromPod 从 Pod 中提取版本信息
func (s *K8sIngressServiceImpl) extractVersionFromPod(pod *corev1.Pod) string {
	// 1. 尝试从 labels 获取版本
	if pod.Labels != nil {
		// 优先级顺序：app.kubernetes.io/version > version > app.version
		if ver, ok := pod.Labels["app.kubernetes.io/version"]; ok && ver != "" {
			return ver
		}
		if ver, ok := pod.Labels["version"]; ok && ver != "" {
			return ver
		}
		if ver, ok := pod.Labels["app.version"]; ok && ver != "" {
			return ver
		}
		if ver, ok := pod.Labels["helm.sh/chart"]; ok && ver != "" {
			// 从 Helm chart 标签提取版本，例如: ingress-nginx-4.0.1
			parts := strings.Split(ver, "-")
			if len(parts) > 1 {
				lastPart := parts[len(parts)-1]
				// 检查最后一部分是否是版本号格式
				if matched, _ := regexp.MatchString(`^\d+\.\d+`, lastPart); matched {
					return lastPart
				}
			}
		}
	}

	// 2. 从容器镜像中提取版本
	if len(pod.Spec.Containers) > 0 {
		for _, container := range pod.Spec.Containers {
			// 检查容器名称是否包含 controller 或 ingress
			if strings.Contains(strings.ToLower(container.Name), "controller") ||
				strings.Contains(strings.ToLower(container.Name), "ingress") {
				version := s.extractVersionFromImage(container.Image)
				if version != "" {
					return version
				}
			}
		}

		// 如果没有匹配的容器名，使用第一个容器
		version := s.extractVersionFromImage(pod.Spec.Containers[0].Image)
		if version != "" {
			return version
		}
	}

	// 3. 尝试从 Pod annotations 获取
	if pod.Annotations != nil {
		if ver, ok := pod.Annotations["version"]; ok && ver != "" {
			return ver
		}
	}

	return ""
}

// extractVersionFromImage 从镜像名称中提取版本
func (s *K8sIngressServiceImpl) extractVersionFromImage(image string) string {
	// 镜像格式示例:
	// registry.k8s.io/ingress-nginx/controller:v1.1.3
	// nginx/nginx-ingress:1.2.0
	// k8s.gcr.io/ingress-nginx/controller:v0.49.0@sha256:xxx

	// 移除 digest 部分
	if idx := strings.Index(image, "@"); idx > 0 {
		image = image[:idx]
	}

	// 提取 tag
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return ""
	}

	tag := parts[len(parts)-1]

	// 过滤掉非版本标签
	nonVersionTags := []string{"latest", "stable", "master", "main", "dev", "develop"}
	lowerTag := strings.ToLower(tag)
	for _, nvt := range nonVersionTags {
		if lowerTag == nvt {
			return ""
		}
	}

	// 返回版本标签（保留 v 前缀）
	return tag
}

// getIngressType 判断 Ingress 类型（公网/内网等）
func (s *K8sIngressServiceImpl) getIngressType(ing *networkingv1.Ingress, ingressClass string) string {
	// 1. 从 annotation 获取自定义类型标记
	if ing.Annotations != nil {
		if ingressType, ok := ing.Annotations["lockin.com/ingress-type"]; ok {
			return ingressType
		}
		if ingressType, ok := ing.Annotations["kubernetes.io/ingress-type"]; ok {
			return ingressType
		}
	}

	// 2. 从 IngressClass 名称推断
	lowerClass := strings.ToLower(ingressClass)

	if strings.Contains(lowerClass, "public") || strings.Contains(lowerClass, "external") {
		return "公网Nginx"
	}

	if strings.Contains(lowerClass, "internal") || strings.Contains(lowerClass, "private") {
		return "内网Nginx"
	}

	// 3. 从 LoadBalancer IP 判断（可以根据 IP 段判断内外网）
	if len(ing.Status.LoadBalancer.Ingress) > 0 {
		ip := ing.Status.LoadBalancer.Ingress[0].IP
		if ip != "" {
			// 判断是否为内网 IP
			if s.isPrivateIP(ip) {
				return "内网Nginx"
			}
			return "公网Nginx"
		}
	}

	// 4. 默认根据 controller 类型返回
	if strings.Contains(lowerClass, "nginx") {
		return "公网Nginx"
	}

	return "Nginx"
}

// isPrivateIP 判断是否为内网IP
func (s *K8sIngressServiceImpl) isPrivateIP(ip string) bool {
	// 内网 IP 段:
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16

	privateIPPatterns := []string{
		`^10\.`,
		`^172\.(1[6-9]|2[0-9]|3[0-1])\.`,
		`^192\.168\.`,
		`^127\.`,
	}

	for _, pattern := range privateIPPatterns {
		matched, _ := regexp.MatchString(pattern, ip)
		if matched {
			return true
		}
	}

	return false
}