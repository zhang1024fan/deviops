package service

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/cache"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// IK8sNamespaceService K8s命名空间服务接口
type IK8sNamespaceService interface {
	GetNamespaces(c *gin.Context, clusterId uint)
	GetNamespace(c *gin.Context, clusterId uint, namespaceName string)
	CreateNamespace(c *gin.Context, clusterId uint, req *model.CreateNamespaceRequest)
	UpdateNamespace(c *gin.Context, clusterId uint, namespaceName string, req *model.UpdateNamespaceRequest)
	DeleteNamespace(c *gin.Context, clusterId uint, namespaceName string)
	
	// ResourceQuota管理
	GetResourceQuotas(c *gin.Context, clusterId uint, namespaceName string)
	CreateResourceQuota(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateResourceQuotaRequest)
	UpdateResourceQuota(c *gin.Context, clusterId uint, namespaceName string, quotaName string, req *model.UpdateResourceQuotaRequest)
	DeleteResourceQuota(c *gin.Context, clusterId uint, namespaceName string, quotaName string)
	
	// LimitRange管理
	GetLimitRanges(c *gin.Context, clusterId uint, namespaceName string)
	CreateLimitRange(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateLimitRangeRequest)
	UpdateLimitRange(c *gin.Context, clusterId uint, namespaceName string, limitRangeName string, req *model.UpdateLimitRangeRequest)
	DeleteLimitRange(c *gin.Context, clusterId uint, namespaceName string, limitRangeName string)
}

// K8sNamespaceServiceImpl K8s命名空间服务实现
type K8sNamespaceServiceImpl struct {
	clusterDao *dao.KubeClusterDao
	cacheService cache.ICacheService
}

func NewK8sNamespaceService(db *gorm.DB) IK8sNamespaceService {
	return &K8sNamespaceServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
		cacheService: nil, // 将通过依赖注入设置
	}
}

// NewK8sNamespaceServiceWithCache 创建带缓存的命名空间服务
func NewK8sNamespaceServiceWithCache(db *gorm.DB, redisClient *redis.Client) IK8sNamespaceService {
	return &K8sNamespaceServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
		cacheService: cache.NewRedisCache(redisClient),
	}
}

// GetNamespaces 获取命名空间列表
func (s *K8sNamespaceServiceImpl) GetNamespaces(c *gin.Context, clusterId uint) {
	ctx := c.Request.Context()
	
	// 尝试从缓存获取
	if s.cacheService != nil {
		var cachedResponse model.NamespaceListResponse
		err := s.cacheService.GetNamespaceList(ctx, clusterId, &cachedResponse)
		if err == nil {
			result.Success(c, cachedResponse)
			return
		}
		// 缓存未命中或错误，继续从K8s获取
	}

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
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法获取命名空间信息")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取命名空间列表
	namespaces, err := s.fetchNamespaces(clientset)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取命名空间列表失败: "+err.Error())
		return
	}

	response := model.NamespaceListResponse{
		Namespaces: namespaces,
		Total:      len(namespaces),
	}

	// 缓存响应结果
	if s.cacheService != nil {
		_ = s.cacheService.SetNamespaceList(ctx, clusterId, response, cache.NamespaceListExpiration)
	}

	result.Success(c, response)
}

// GetNamespace 获取单个命名空间详情
func (s *K8sNamespaceServiceImpl) GetNamespace(c *gin.Context, clusterId uint, namespaceName string) {
	ctx := c.Request.Context()
	
	// 尝试从缓存获取
	if s.cacheService != nil {
		var cachedNamespace model.K8sNamespace
		err := s.cacheService.GetNamespaceDetail(ctx, clusterId, namespaceName, &cachedNamespace)
		if err == nil {
			result.Success(c, &cachedNamespace)
			return
		}
		// 缓存未命中或错误，继续从K8s获取
	}

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

	// 获取命名空间详情
	namespace, err := s.fetchNamespaceDetail(clientset, namespaceName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取命名空间详情失败: "+err.Error())
		return
	}

	// 缓存命名空间详情
	if s.cacheService != nil {
		_ = s.cacheService.SetNamespaceDetail(ctx, clusterId, namespaceName, namespace, cache.NamespaceDetailExpiration)
	}

	result.Success(c, namespace)
}

// CreateNamespace 创建命名空间
func (s *K8sNamespaceServiceImpl) CreateNamespace(c *gin.Context, clusterId uint, req *model.CreateNamespaceRequest) {
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

	// 验证命名空间名称
	if err := s.validateNamespaceName(req.Name); err != nil {
		result.Failed(c, http.StatusBadRequest, "命名空间名称不符合规范: "+err.Error())
		return
	}
	
	// 验证标签
	if err := s.validateLabels(req.Labels); err != nil {
		result.Failed(c, http.StatusBadRequest, "标签格式不符合规范: "+err.Error())
		return
	}
	
	// 验证注释
	if err := s.validateAnnotations(req.Annotations); err != nil {
		result.Failed(c, http.StatusBadRequest, "注释格式不符合规范: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 检查命名空间是否已存在
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), req.Name, metav1.GetOptions{})
	if err == nil {
		result.Failed(c, http.StatusBadRequest, "命名空间已存在")
		return
	}

	// 创建命名空间对象
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}

	// 创建命名空间
	createdNs, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建命名空间失败: "+err.Error())
		return
	}

	// 转换为响应格式
	k8sNamespace := s.convertToK8sNamespace(createdNs)

	// 使命名空间缓存失效（因为有新命名空间创建）
	if s.cacheService != nil {
		ctx := c.Request.Context()
		_ = s.cacheService.InvalidateNamespaceCache(ctx, clusterId)
	}

	result.Success(c, k8sNamespace)
}

// UpdateNamespace 更新命名空间
func (s *K8sNamespaceServiceImpl) UpdateNamespace(c *gin.Context, clusterId uint, namespaceName string, req *model.UpdateNamespaceRequest) {
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

	// 验证标签
	if err := s.validateLabels(req.Labels); err != nil {
		result.Failed(c, http.StatusBadRequest, "标签格式不符合规范: "+err.Error())
		return
	}
	
	// 验证注释
	if err := s.validateAnnotations(req.Annotations); err != nil {
		result.Failed(c, http.StatusBadRequest, "注释格式不符合规范: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有命名空间
	namespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "命名空间不存在: "+err.Error())
		return
	}

	// 更新标签和注释
	if req.Labels != nil {
		namespace.Labels = req.Labels
	}
	if req.Annotations != nil {
		namespace.Annotations = req.Annotations
	}

	// 更新命名空间
	updatedNs, err := clientset.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新命名空间失败: "+err.Error())
		return
	}

	// 转换为响应格式
	k8sNamespace := s.convertToK8sNamespace(updatedNs)

	// 使该命名空间的缓存失效
	if s.cacheService != nil {
		ctx := c.Request.Context()
		_ = s.cacheService.InvalidateNamespaceCache(ctx, clusterId, namespaceName)
	}

	result.Success(c, k8sNamespace)
}

// DeleteNamespace 删除命名空间
func (s *K8sNamespaceServiceImpl) DeleteNamespace(c *gin.Context, clusterId uint, namespaceName string) {
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

	// 检查命名空间是否存在
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "命名空间不存在: "+err.Error())
		return
	}

	// 删除命名空间
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除命名空间失败: "+err.Error())
		return
	}

	// 使命名空间缓存失效（因为命名空间被删除）
	if s.cacheService != nil {
		ctx := c.Request.Context()
		_ = s.cacheService.InvalidateNamespaceCache(ctx, clusterId, namespaceName)
	}

	result.Success(c, "命名空间删除成功")
}

// createK8sClient 创建K8s客户端
func (s *K8sNamespaceServiceImpl) createK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	// 设置API版本兼容配置,确保兼容旧版本K8s集群(如v1.22.5)
	// 忽略API版本弃用警告,提高与旧版本K8s的兼容性
	config.WarningHandler = rest.NoWarnings{}

	// 设置超时,避免长时间等待
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// fetchNamespaces 获取所有命名空间
func (s *K8sNamespaceServiceImpl) fetchNamespaces(clientset *kubernetes.Clientset) ([]model.K8sNamespace, error) {
	// 设置超时上下文,避免长时间等待
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取命名空间列表
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 只返回基本命名空间信息,不获取子资源详情
	// 子资源详情在获取单个命名空间时再获取,避免列表接口响应过慢
	var namespaces []model.K8sNamespace
	for _, ns := range nsList.Items {
		k8sNamespace := s.convertToK8sNamespace(&ns)
		namespaces = append(namespaces, k8sNamespace)
	}

	return namespaces, nil
}

// fetchNamespaceDetail 获取命名空间详情
func (s *K8sNamespaceServiceImpl) fetchNamespaceDetail(clientset *kubernetes.Clientset, namespaceName string) (*model.K8sNamespace, error) {
	// 设置超时上下文,避免长时间等待
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取命名空间
	ns, err := clientset.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	k8sNamespace := s.convertToK8sNamespace(ns)

	// 获取资源配额详情(忽略错误,不影响主要功能)
	quotas, _ := s.getResourceQuotaDetails(clientset, namespaceName)
	k8sNamespace.ResourceQuotas = quotas

	// 获取LimitRange详情(忽略错误,不影响主要功能)
	limitRanges, _ := s.getLimitRangeDetails(clientset, namespaceName)
	k8sNamespace.LimitRanges = limitRanges

	// 获取资源统计(忽略错误,不影响主要功能)
	resourceCount, _ := s.getNamespaceResourceCount(clientset, namespaceName)
	k8sNamespace.ResourceCount = resourceCount

	return &k8sNamespace, nil
}

// convertToK8sNamespace 转换为K8sNamespace结构
func (s *K8sNamespaceServiceImpl) convertToK8sNamespace(ns *corev1.Namespace) model.K8sNamespace {
	return model.K8sNamespace{
		Name:        ns.Name,
		Status:      string(ns.Status.Phase),
		Labels:      ns.Labels,
		Annotations: ns.Annotations,
		CreatedAt:   ns.CreationTimestamp.Format(time.RFC3339),
	}
}


// validateNamespaceName 验证命名空间名称
func (s *K8sNamespaceServiceImpl) validateNamespaceName(name string) error {
	if name == "" {
		return fmt.Errorf("命名空间名称不能为空")
	}
	
	if len(name) > 63 {
		return fmt.Errorf("命名空间名称长度不能超过63个字符")
	}
	
	// K8s命名空间名称规则：只能包含小写字母、数字和连字符，且必须以字母数字开头和结尾
	if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(string(name[0]))) {
		return fmt.Errorf("命名空间名称必须以字母或数字开头")
	}
	
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("命名空间名称只能包含小写字母、数字和连字符")
		}
	}
	
	if strings.HasSuffix(name, "-") {
		return fmt.Errorf("命名空间名称不能以连字符结尾")
	}
	
	return nil
}

// GetResourceQuotas 获取ResourceQuota列表
func (s *K8sNamespaceServiceImpl) GetResourceQuotas(c *gin.Context, clusterId uint, namespaceName string) {
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

	// 获取ResourceQuota列表
	quotaList, err := clientset.CoreV1().ResourceQuotas(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取ResourceQuota列表失败: "+err.Error())
		return
	}

	var quotas []model.ResourceQuotaDetail
	for _, quota := range quotaList.Items {
		quotaDetail := s.convertToResourceQuotaDetail(&quota)
		quotas = append(quotas, quotaDetail)
	}

	response := model.ResourceQuotaListResponse{
		ResourceQuotas: quotas,
		Total:          len(quotas),
	}

	result.Success(c, response)
}

// CreateResourceQuota 创建ResourceQuota
func (s *K8sNamespaceServiceImpl) CreateResourceQuota(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateResourceQuotaRequest) {
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

	// 检查ResourceQuota是否已存在
	_, err = clientset.CoreV1().ResourceQuotas(namespaceName).Get(context.TODO(), req.Name, metav1.GetOptions{})
	if err == nil {
		result.Failed(c, http.StatusBadRequest, "ResourceQuota已存在")
		return
	}

	// 创建ResourceQuota对象
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: make(corev1.ResourceList),
		},
	}

	// 设置资源限制
	for resourceName, quantityStr := range req.Hard {
		quantity, err := resource.ParseQuantity(quantityStr)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, "无效的资源量格式: "+err.Error())
			return
		}
		quota.Spec.Hard[corev1.ResourceName(resourceName)] = quantity
	}

	// 创建ResourceQuota
	createdQuota, err := clientset.CoreV1().ResourceQuotas(namespaceName).Create(context.TODO(), quota, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建ResourceQuota失败: "+err.Error())
		return
	}

	quotaDetail := s.convertToResourceQuotaDetail(createdQuota)
	result.Success(c, quotaDetail)
}

// UpdateResourceQuota 更新ResourceQuota
func (s *K8sNamespaceServiceImpl) UpdateResourceQuota(c *gin.Context, clusterId uint, namespaceName string, quotaName string, req *model.UpdateResourceQuotaRequest) {
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

	// 获取现有ResourceQuota
	quota, err := clientset.CoreV1().ResourceQuotas(namespaceName).Get(context.TODO(), quotaName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "ResourceQuota不存在: "+err.Error())
		return
	}

	// 更新资源限制
	if req.Hard != nil {
		quota.Spec.Hard = make(corev1.ResourceList)
		for resourceName, quantityStr := range req.Hard {
			quantity, err := resource.ParseQuantity(quantityStr)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, "无效的资源量格式: "+err.Error())
				return
			}
			quota.Spec.Hard[corev1.ResourceName(resourceName)] = quantity
		}
	}

	// 更新ResourceQuota
	updatedQuota, err := clientset.CoreV1().ResourceQuotas(namespaceName).Update(context.TODO(), quota, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新ResourceQuota失败: "+err.Error())
		return
	}

	quotaDetail := s.convertToResourceQuotaDetail(updatedQuota)
	result.Success(c, quotaDetail)
}

// DeleteResourceQuota 删除ResourceQuota
func (s *K8sNamespaceServiceImpl) DeleteResourceQuota(c *gin.Context, clusterId uint, namespaceName string, quotaName string) {
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

	// 删除ResourceQuota
	err = clientset.CoreV1().ResourceQuotas(namespaceName).Delete(context.TODO(), quotaName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除ResourceQuota失败: "+err.Error())
		return
	}

	result.Success(c, "ResourceQuota删除成功")
}

// convertToResourceQuotaDetail 转换为ResourceQuotaDetail结构
func (s *K8sNamespaceServiceImpl) convertToResourceQuotaDetail(quota *corev1.ResourceQuota) model.ResourceQuotaDetail {
	hard := make(map[string]string)
	used := make(map[string]string)

	for resource, quantity := range quota.Status.Hard {
		hard[string(resource)] = quantity.String()
	}

	for resource, quantity := range quota.Status.Used {
		used[string(resource)] = quantity.String()
	}

	return model.ResourceQuotaDetail{
		Name:      quota.Name,
		Hard:      hard,
		Used:      used,
		CreatedAt: quota.CreationTimestamp.Format(time.RFC3339),
	}
}

// GetLimitRanges 获取LimitRange列表
func (s *K8sNamespaceServiceImpl) GetLimitRanges(c *gin.Context, clusterId uint, namespaceName string) {
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

	// 获取LimitRange列表
	limitRangeList, err := clientset.CoreV1().LimitRanges(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取LimitRange列表失败: "+err.Error())
		return
	}

	var limitRanges []model.LimitRangeDetail
	for _, lr := range limitRangeList.Items {
		limitRangeDetail := s.convertToLimitRangeDetail(&lr)
		limitRanges = append(limitRanges, limitRangeDetail)
	}

	response := model.LimitRangeListResponse{
		LimitRanges: limitRanges,
		Total:       len(limitRanges),
	}

	result.Success(c, response)
}

// CreateLimitRange 创建LimitRange
func (s *K8sNamespaceServiceImpl) CreateLimitRange(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateLimitRangeRequest) {
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

	// 检查LimitRange是否已存在
	_, err = clientset.CoreV1().LimitRanges(namespaceName).Get(context.TODO(), req.Name, metav1.GetOptions{})
	if err == nil {
		result.Failed(c, http.StatusBadRequest, "LimitRange已存在")
		return
	}

	// 创建LimitRange对象
	limitRange := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{},
		},
	}

	// 设置限制规则
	for _, limit := range req.Spec.Limits {
		limitRangeItem := corev1.LimitRangeItem{
			Type:                 corev1.LimitType(limit.Type),
			Max:                  make(corev1.ResourceList),
			Min:                  make(corev1.ResourceList),
			Default:              make(corev1.ResourceList),
			DefaultRequest:       make(corev1.ResourceList),
			MaxLimitRequestRatio: make(corev1.ResourceList),
		}

		// 设置各种资源限制
		for resourceName, quantityStr := range limit.Max {
			quantity, err := resource.ParseQuantity(quantityStr)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, "无效的Max资源量格式: "+err.Error())
				return
			}
			limitRangeItem.Max[corev1.ResourceName(resourceName)] = quantity
		}

		for resourceName, quantityStr := range limit.Min {
			quantity, err := resource.ParseQuantity(quantityStr)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, "无效的Min资源量格式: "+err.Error())
				return
			}
			limitRangeItem.Min[corev1.ResourceName(resourceName)] = quantity
		}

		for resourceName, quantityStr := range limit.Default {
			quantity, err := resource.ParseQuantity(quantityStr)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, "无效的Default资源量格式: "+err.Error())
				return
			}
			limitRangeItem.Default[corev1.ResourceName(resourceName)] = quantity
		}

		for resourceName, quantityStr := range limit.DefaultRequest {
			quantity, err := resource.ParseQuantity(quantityStr)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, "无效的DefaultRequest资源量格式: "+err.Error())
				return
			}
			limitRangeItem.DefaultRequest[corev1.ResourceName(resourceName)] = quantity
		}

		limitRange.Spec.Limits = append(limitRange.Spec.Limits, limitRangeItem)
	}

	// 创建LimitRange
	createdLimitRange, err := clientset.CoreV1().LimitRanges(namespaceName).Create(context.TODO(), limitRange, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建LimitRange失败: "+err.Error())
		return
	}

	limitRangeDetail := s.convertToLimitRangeDetail(createdLimitRange)
	result.Success(c, limitRangeDetail)
}

// UpdateLimitRange 更新LimitRange
func (s *K8sNamespaceServiceImpl) UpdateLimitRange(c *gin.Context, clusterId uint, namespaceName string, limitRangeName string, req *model.UpdateLimitRangeRequest) {
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

	// 获取现有LimitRange
	limitRange, err := clientset.CoreV1().LimitRanges(namespaceName).Get(context.TODO(), limitRangeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "LimitRange不存在: "+err.Error())
		return
	}

	// 更新限制规则
	if req.Spec.Limits != nil {
		limitRange.Spec.Limits = []corev1.LimitRangeItem{}
		for _, limit := range req.Spec.Limits {
			limitRangeItem := corev1.LimitRangeItem{
				Type:                 corev1.LimitType(limit.Type),
				Max:                  make(corev1.ResourceList),
				Min:                  make(corev1.ResourceList),
				Default:              make(corev1.ResourceList),
				DefaultRequest:       make(corev1.ResourceList),
				MaxLimitRequestRatio: make(corev1.ResourceList),
			}

			// 设置各种资源限制
			for resourceName, quantityStr := range limit.Max {
				quantity, err := resource.ParseQuantity(quantityStr)
				if err != nil {
					result.Failed(c, http.StatusBadRequest, "无效的Max资源量格式: "+err.Error())
					return
				}
				limitRangeItem.Max[corev1.ResourceName(resourceName)] = quantity
			}

			for resourceName, quantityStr := range limit.Min {
				quantity, err := resource.ParseQuantity(quantityStr)
				if err != nil {
					result.Failed(c, http.StatusBadRequest, "无效的Min资源量格式: "+err.Error())
					return
				}
				limitRangeItem.Min[corev1.ResourceName(resourceName)] = quantity
			}

			for resourceName, quantityStr := range limit.Default {
				quantity, err := resource.ParseQuantity(quantityStr)
				if err != nil {
					result.Failed(c, http.StatusBadRequest, "无效的Default资源量格式: "+err.Error())
					return
				}
				limitRangeItem.Default[corev1.ResourceName(resourceName)] = quantity
			}

			for resourceName, quantityStr := range limit.DefaultRequest {
				quantity, err := resource.ParseQuantity(quantityStr)
				if err != nil {
					result.Failed(c, http.StatusBadRequest, "无效的DefaultRequest资源量格式: "+err.Error())
					return
				}
				limitRangeItem.DefaultRequest[corev1.ResourceName(resourceName)] = quantity
			}

			limitRange.Spec.Limits = append(limitRange.Spec.Limits, limitRangeItem)
		}
	}

	// 更新LimitRange
	updatedLimitRange, err := clientset.CoreV1().LimitRanges(namespaceName).Update(context.TODO(), limitRange, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新LimitRange失败: "+err.Error())
		return
	}

	limitRangeDetail := s.convertToLimitRangeDetail(updatedLimitRange)
	result.Success(c, limitRangeDetail)
}

// DeleteLimitRange 删除LimitRange
func (s *K8sNamespaceServiceImpl) DeleteLimitRange(c *gin.Context, clusterId uint, namespaceName string, limitRangeName string) {
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

	// 删除LimitRange
	err = clientset.CoreV1().LimitRanges(namespaceName).Delete(context.TODO(), limitRangeName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除LimitRange失败: "+err.Error())
		return
	}

	result.Success(c, "LimitRange删除成功")
}

// convertToLimitRangeDetail 转换为LimitRangeDetail结构
func (s *K8sNamespaceServiceImpl) convertToLimitRangeDetail(lr *corev1.LimitRange) model.LimitRangeDetail {
	var limits []model.LimitRangeItem
	
	for _, item := range lr.Spec.Limits {
		limitItem := model.LimitRangeItem{
			Type:           string(item.Type),
			Max:            make(map[string]string),
			Min:            make(map[string]string),
			Default:        make(map[string]string),
			DefaultRequest: make(map[string]string),
		}

		for resource, quantity := range item.Max {
			limitItem.Max[string(resource)] = quantity.String()
		}

		for resource, quantity := range item.Min {
			limitItem.Min[string(resource)] = quantity.String()
		}

		for resource, quantity := range item.Default {
			limitItem.Default[string(resource)] = quantity.String()
		}

		for resource, quantity := range item.DefaultRequest {
			limitItem.DefaultRequest[string(resource)] = quantity.String()
		}

		limits = append(limits, limitItem)
	}

	return model.LimitRangeDetail{
		Name:      lr.Name,
		Limits:    limits,
		CreatedAt: lr.CreationTimestamp.Format(time.RFC3339),
	}
}

// getResourceQuotaDetails 获取ResourceQuota详情列表
func (s *K8sNamespaceServiceImpl) getResourceQuotaDetails(clientset *kubernetes.Clientset, namespaceName string) ([]model.ResourceQuotaDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quotas, err := clientset.CoreV1().ResourceQuotas(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil || len(quotas.Items) == 0 {
		return nil, err
	}

	var quotaDetails []model.ResourceQuotaDetail
	for _, quota := range quotas.Items {
		quotaDetail := s.convertToResourceQuotaDetail(&quota)
		quotaDetails = append(quotaDetails, quotaDetail)
	}

	return quotaDetails, nil
}

// getLimitRangeDetails 获取LimitRange详情列表
func (s *K8sNamespaceServiceImpl) getLimitRangeDetails(clientset *kubernetes.Clientset, namespaceName string) ([]model.LimitRangeDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limitRanges, err := clientset.CoreV1().LimitRanges(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil || len(limitRanges.Items) == 0 {
		return nil, err
	}

	var limitRangeDetails []model.LimitRangeDetail
	for _, lr := range limitRanges.Items {
		limitRangeDetail := s.convertToLimitRangeDetail(&lr)
		limitRangeDetails = append(limitRangeDetails, limitRangeDetail)
	}

	return limitRangeDetails, nil
}

// getNamespaceResourceCount 获取命名空间资源统计
func (s *K8sNamespaceServiceImpl) getNamespaceResourceCount(clientset *kubernetes.Clientset, namespaceName string) (model.NamespaceResourceCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resourceCount := model.NamespaceResourceCount{}

	// 获取Pod数量
	pods, err := clientset.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{})
	if err == nil {
		resourceCount.PodCount = len(pods.Items)
	}

	// 获取Service数量
	services, err := clientset.CoreV1().Services(namespaceName).List(ctx, metav1.ListOptions{})
	if err == nil {
		resourceCount.ServiceCount = len(services.Items)
	}

	// 获取Secret数量
	secrets, err := clientset.CoreV1().Secrets(namespaceName).List(ctx, metav1.ListOptions{})
	if err == nil {
		resourceCount.SecretCount = len(secrets.Items)
	}

	// 获取ConfigMap数量
	configMaps, err := clientset.CoreV1().ConfigMaps(namespaceName).List(ctx, metav1.ListOptions{})
	if err == nil {
		resourceCount.ConfigMapCount = len(configMaps.Items)
	}

	return resourceCount, nil
}

// validateLabels 验证标签格式
func (s *K8sNamespaceServiceImpl) validateLabels(labels map[string]string) error {
	if labels == nil {
		return nil
	}
	
	// K8s标签键的正则表达式：可选前缀/键名
	keyRegex := regexp.MustCompile(`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/)?[a-z0-9A-Z]([-._a-z0-9A-Z]*[a-z0-9A-Z])?$`)
	// K8s标签值的正则表达式：空字符串或字母数字字符、连字符、下划线和点号
	valueRegex := regexp.MustCompile(`^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`)
	
	for key, value := range labels {
		// 验证键名长度
		if len(key) == 0 {
			return fmt.Errorf("标签键不能为空")
		}
		if len(key) > 253 {
			return fmt.Errorf("标签键'%s'长度超过253个字符", key)
		}
		
		// 验证键名格式
		if !keyRegex.MatchString(key) {
			return fmt.Errorf("标签键'%s'格式不符合规范，只能包含小写字母、数字、连字符、点号，可选的DNS前缀", key)
		}
		
		// 验证值长度
		if len(value) > 63 {
			return fmt.Errorf("标签值'%s'长度超过63个字符", value)
		}
		
		// 验证值格式
		if !valueRegex.MatchString(value) {
			return fmt.Errorf("标签值'%s'格式不符合规范，只能包含字母数字字符、连字符、下划线和点号，且必须以字母数字字符开头和结尾", value)
		}
	}
	
	return nil
}

// validateAnnotations 验证注释格式
func (s *K8sNamespaceServiceImpl) validateAnnotations(annotations map[string]string) error {
	if annotations == nil {
		return nil
	}
	
	// K8s注释键的正则表达式（与标签键相同）
	keyRegex := regexp.MustCompile(`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/)?[a-z0-9A-Z]([-._a-z0-9A-Z]*[a-z0-9A-Z])?$`)
	
	for key, value := range annotations {
		// 验证键名长度
		if len(key) == 0 {
			return fmt.Errorf("注释键不能为空")
		}
		if len(key) > 253 {
			return fmt.Errorf("注释键'%s'长度超过253个字符", key)
		}
		
		// 验证键名格式
		if !keyRegex.MatchString(key) {
			return fmt.Errorf("注释键'%s'格式不符合规范，只能包含小写字母、数字、连字符、点号，可选的DNS前缀", key)
		}
		
		// 注释值的长度限制更宽松，但仍有限制
		if len(value) > 262144 { // 256KB
			return fmt.Errorf("注释值长度超过256KB限制")
		}
	}
	
	return nil
}