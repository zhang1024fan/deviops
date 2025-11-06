package service

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// IK8sEventsService K8s事件服务接口
type IK8sEventsService interface {
	GetEvents(c *gin.Context, clusterId uint, namespaceName string, kind string, name string, limit int)
	GetClusterEvents(c *gin.Context, clusterId uint, limit int)
}

// K8sEventsServiceImpl K8s事件服务实现
type K8sEventsServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

func NewK8sEventsService(db *gorm.DB) IK8sEventsService {
	return &K8sEventsServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// NewK8sEventsServiceImpl 创建K8s事件服务实现（用于controller直接调用）
func NewK8sEventsServiceImpl(db *gorm.DB) *K8sEventsServiceImpl {
	return &K8sEventsServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// GetEvents 获取指定命名空间的事件列表
func (s *K8sEventsServiceImpl) GetEvents(c *gin.Context, clusterId uint, namespaceName string, kind string, name string, limit int) {
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
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法获取事件信息")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建查询选项
	listOptions := metav1.ListOptions{}
	
	// 如果指定了资源类型和名称，添加字段选择器
	if kind != "" && name != "" {
		listOptions.FieldSelector = "involvedObject.kind=" + kind + ",involvedObject.name=" + name
	} else if kind != "" {
		listOptions.FieldSelector = "involvedObject.kind=" + kind
	} else if name != "" {
		listOptions.FieldSelector = "involvedObject.name=" + name
	}

	// 获取事件列表
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), listOptions)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取事件列表失败: "+err.Error())
		return
	}

	// 转换事件数据
	k8sEvents := s.convertToK8sEvents(events.Items)
	
	// 按时间降序排序（最新的在前面）
	sort.Slice(k8sEvents, func(i, j int) bool {
		return k8sEvents[i].LastTime > k8sEvents[j].LastTime
	})

	// 限制返回数量
	if limit > 0 && len(k8sEvents) > limit {
		k8sEvents = k8sEvents[:limit]
	}

	// 构造响应
	response := model.EventListResponse{
		Events:    k8sEvents,
		Total:     len(k8sEvents),
		Namespace: namespaceName,
		Filter: map[string]string{
			"kind": kind,
			"name": name,
		},
	}

	result.Success(c, response)
}

// GetClusterEvents 获取整个集群的事件列表
func (s *K8sEventsServiceImpl) GetClusterEvents(c *gin.Context, clusterId uint, limit int) {
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
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法获取事件信息")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取所有命名空间的事件
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取集群事件失败: "+err.Error())
		return
	}

	// 转换事件数据
	k8sEvents := s.convertToK8sEvents(events.Items)
	
	// 按时间降序排序
	sort.Slice(k8sEvents, func(i, j int) bool {
		return k8sEvents[i].LastTime > k8sEvents[j].LastTime
	})

	// 限制返回数量
	if limit > 0 && len(k8sEvents) > limit {
		k8sEvents = k8sEvents[:limit]
	}

	// 构造响应
	response := model.EventListResponse{
		Events: k8sEvents,
		Total:  len(k8sEvents),
		Filter: map[string]string{},
	}

	result.Success(c, response)
}

// createK8sClient 创建K8s客户端
func (s *K8sEventsServiceImpl) createK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
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

// convertToK8sEvents 转换K8s事件为模型
func (s *K8sEventsServiceImpl) convertToK8sEvents(events []corev1.Event) []model.K8sEvent {
	var k8sEvents []model.K8sEvent

	for _, event := range events {
		k8sEvent := model.K8sEvent{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			Source:    event.Source.Component,
			Count:     event.Count,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
		}

		// 添加涉及的对象信息
		if event.InvolvedObject.Kind != "" && event.InvolvedObject.Name != "" {
			k8sEvent.Message = strings.TrimSpace(event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name + ": " + k8sEvent.Message)
		}

		k8sEvents = append(k8sEvents, k8sEvent)
	}

	return k8sEvents
}