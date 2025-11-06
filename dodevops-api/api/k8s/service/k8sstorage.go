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
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// IK8sStorageService K8s存储管理服务接口
type IK8sStorageService interface {
	// PVC管理
	GetPVCs(c *gin.Context, clusterId uint, namespaceName string)
	GetPVCDetail(c *gin.Context, clusterId uint, namespaceName string, pvcName string)
	CreatePVC(c *gin.Context, clusterId uint, namespaceName string, req *model.CreatePVCRequest)
	UpdatePVC(c *gin.Context, clusterId uint, namespaceName string, pvcName string, req *model.UpdatePVCRequest)
	DeletePVC(c *gin.Context, clusterId uint, namespaceName string, pvcName string)
	GetPVCYaml(c *gin.Context, clusterId uint, namespaceName string, pvcName string)
	UpdatePVCYaml(c *gin.Context, clusterId uint, namespaceName string, pvcName string, yamlData map[string]interface{})

	// PV管理
	GetPVs(c *gin.Context, clusterId uint)
	GetPVDetail(c *gin.Context, clusterId uint, pvName string)
	CreatePV(c *gin.Context, clusterId uint, req *model.CreatePVRequest)
	UpdatePV(c *gin.Context, clusterId uint, pvName string, req *model.UpdatePVRequest)
	DeletePV(c *gin.Context, clusterId uint, pvName string)
	GetPVYaml(c *gin.Context, clusterId uint, pvName string)
	UpdatePVYaml(c *gin.Context, clusterId uint, pvName string, yamlData map[string]interface{})

	// StorageClass管理
	GetStorageClasses(c *gin.Context, clusterId uint)
	GetStorageClassDetail(c *gin.Context, clusterId uint, storageClassName string)
	CreateStorageClass(c *gin.Context, clusterId uint, req *model.CreateStorageClassRequest)
	UpdateStorageClass(c *gin.Context, clusterId uint, storageClassName string, req *model.UpdateStorageClassRequest)
	DeleteStorageClass(c *gin.Context, clusterId uint, storageClassName string)
	GetStorageClassYaml(c *gin.Context, clusterId uint, storageClassName string)
	UpdateStorageClassYaml(c *gin.Context, clusterId uint, storageClassName string, yamlData map[string]interface{})
}

// K8sStorageServiceImpl K8s存储管理服务实现
type K8sStorageServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

// NewK8sStorageService 创建K8s存储管理服务实例
func NewK8sStorageService(db *gorm.DB) IK8sStorageService {
	return &K8sStorageServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// getK8sClient 获取K8s客户端
func (s *K8sStorageServiceImpl) getK8sClient(clusterId uint) (*kubernetes.Clientset, error) {
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

// ===================== PVC 管理实现 =====================

// GetPVCs 获取PVC列表
func (s *K8sStorageServiceImpl) GetPVCs(c *gin.Context, clusterId uint, namespaceName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取PVC列表失败: "+err.Error())
		return
	}

	var pvcList []model.K8sPersistentVolumeClaim
	for _, pvc := range pvcs.Items {
		k8sPVC := s.convertToK8sPVC(&pvc)
		pvcList = append(pvcList, k8sPVC)
	}

	response := model.PVCListResponse{
		PVCs:  pvcList,
		Total: len(pvcList),
	}

	result.Success(c, response)
}

// GetPVCDetail 获取PVC详情
func (s *K8sStorageServiceImpl) GetPVCDetail(c *gin.Context, clusterId uint, namespaceName string, pvcName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PVC不存在: "+err.Error())
		return
	}

	k8sPVC := s.convertToK8sPVC(pvc)

	// 获取相关事件
	events, _ := s.getPVCEvents(clientset, namespaceName, pvcName)

	pvcDetail := model.PVCDetail{
		K8sPersistentVolumeClaim: k8sPVC,
		Events:                   events,
		Spec:                     pvc.Spec,
	}

	result.Success(c, pvcDetail)
}

// CreatePVC 创建PVC
func (s *K8sStorageServiceImpl) CreatePVC(c *gin.Context, clusterId uint, namespaceName string, req *model.CreatePVCRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建PVC对象
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
			Labels:    req.Labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{},
			},
		},
	}

	// 设置访问模式
	for _, mode := range req.AccessModes {
		pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, corev1.PersistentVolumeAccessMode(mode))
	}

	// 设置资源请求
	for key, value := range req.Resources.Requests {
		quantity, err := resource.ParseQuantity(value)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("无效的资源请求 %s: %v", key, err))
			return
		}
		pvc.Spec.Resources.Requests[corev1.ResourceName(key)] = quantity
	}

	// 设置资源限制
	if len(req.Resources.Limits) > 0 {
		pvc.Spec.Resources.Limits = corev1.ResourceList{}
		for key, value := range req.Resources.Limits {
			quantity, err := resource.ParseQuantity(value)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, fmt.Sprintf("无效的资源限制 %s: %v", key, err))
				return
			}
			pvc.Spec.Resources.Limits[corev1.ResourceName(key)] = quantity
		}
	}

	// 设置存储类
	if req.StorageClass != "" {
		pvc.Spec.StorageClassName = &req.StorageClass
	}

	// 设置卷模式
	if req.VolumeMode != "" {
		volumeMode := corev1.PersistentVolumeMode(req.VolumeMode)
		pvc.Spec.VolumeMode = &volumeMode
	}

	// 设置选择器
	if req.Selector != nil {
		pvc.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: req.Selector.MatchLabels,
		}
		if len(req.Selector.MatchExpressions) > 0 {
			for _, expr := range req.Selector.MatchExpressions {
				pvc.Spec.Selector.MatchExpressions = append(pvc.Spec.Selector.MatchExpressions, metav1.LabelSelectorRequirement{
					Key:      expr.Key,
					Operator: metav1.LabelSelectorOperator(expr.Operator),
					Values:   expr.Values,
				})
			}
		}
	}

	// 创建PVC
	createdPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建PVC失败: "+err.Error())
		return
	}

	k8sPVC := s.convertToK8sPVC(createdPVC)
	result.Success(c, k8sPVC)
}

// UpdatePVC 更新PVC
func (s *K8sStorageServiceImpl) UpdatePVC(c *gin.Context, clusterId uint, namespaceName string, pvcName string, req *model.UpdatePVCRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有PVC
	existingPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PVC不存在: "+err.Error())
		return
	}

	// 更新标签
	if req.Labels != nil {
		existingPVC.Labels = req.Labels
	}

	// 更新资源请求（仅支持存储扩容）
	if len(req.Resources.Requests) > 0 {
		for key, value := range req.Resources.Requests {
			quantity, err := resource.ParseQuantity(value)
			if err != nil {
				result.Failed(c, http.StatusBadRequest, fmt.Sprintf("无效的资源请求 %s: %v", key, err))
				return
			}
			existingPVC.Spec.Resources.Requests[corev1.ResourceName(key)] = quantity
		}
	}

	// 更新PVC
	updatedPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Update(context.TODO(), existingPVC, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新PVC失败: "+err.Error())
		return
	}

	k8sPVC := s.convertToK8sPVC(updatedPVC)
	result.Success(c, k8sPVC)
}

// DeletePVC 删除PVC
func (s *K8sStorageServiceImpl) DeletePVC(c *gin.Context, clusterId uint, namespaceName string, pvcName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.CoreV1().PersistentVolumeClaims(namespaceName).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除PVC失败: "+err.Error())
		return
	}

	result.Success(c, "PVC删除成功")
}

// GetPVCYaml 获取PVC的YAML配置
func (s *K8sStorageServiceImpl) GetPVCYaml(c *gin.Context, clusterId uint, namespaceName string, pvcName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PVC不存在: "+err.Error())
		return
	}

	// 清理不需要的字段，但保留 status
	pvc.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if pvc.APIVersion == "" {
		pvc.APIVersion = "v1"
	}
	if pvc.Kind == "" {
		pvc.Kind = "PersistentVolumeClaim"
	}

	// 清理元数据字段
	pvc.ResourceVersion = ""
	pvc.UID = ""
	pvc.Generation = 0
	pvc.CreationTimestamp = metav1.Time{}
	pvc.SelfLink = ""

	yamlData, err := yaml.Marshal(pvc)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "转换YAML失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": string(yamlData),
	})
}

// UpdatePVCYaml 通过YAML更新PVC
func (s *K8sStorageServiceImpl) UpdatePVCYaml(c *gin.Context, clusterId uint, namespaceName string, pvcName string, yamlData map[string]interface{}) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 将map转换为YAML字符串
	yamlBytes, err := json.Marshal(yamlData)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML数据格式错误: "+err.Error())
		return
	}

	var pvc corev1.PersistentVolumeClaim
	err = yaml.Unmarshal(yamlBytes, &pvc)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}

	// 确保名称和命名空间正确
	pvc.Name = pvcName
	pvc.Namespace = namespaceName

	// 更新PVC
	_, err = clientset.CoreV1().PersistentVolumeClaims(namespaceName).Update(context.TODO(), &pvc, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新PVC失败: "+err.Error())
		return
	}

	result.Success(c, "PVC YAML更新成功")
}

// ===================== PV 管理实现 =====================

// GetPVs 获取PV列表
func (s *K8sStorageServiceImpl) GetPVs(c *gin.Context, clusterId uint) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取PV列表失败: "+err.Error())
		return
	}

	var pvList []model.K8sPersistentVolume
	for _, pv := range pvs.Items {
		k8sPV := s.convertToK8sPV(&pv)
		pvList = append(pvList, k8sPV)
	}

	response := model.PVListResponse{
		PVs:   pvList,
		Total: len(pvList),
	}

	result.Success(c, response)
}

// GetPVDetail 获取PV详情
func (s *K8sStorageServiceImpl) GetPVDetail(c *gin.Context, clusterId uint, pvName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pv, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PV不存在: "+err.Error())
		return
	}

	k8sPV := s.convertToK8sPV(pv)

	// 获取相关事件
	events, _ := s.getPVEvents(clientset, pvName)

	pvDetail := model.PVDetail{
		K8sPersistentVolume: k8sPV,
		Events:              events,
		Spec:                pv.Spec,
	}

	result.Success(c, pvDetail)
}

// CreatePV 创建PV
func (s *K8sStorageServiceImpl) CreatePV(c *gin.Context, clusterId uint, req *model.CreatePVRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建PV对象
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: req.Labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity:    corev1.ResourceList{},
			AccessModes: []corev1.PersistentVolumeAccessMode{},
		},
	}

	// 设置容量
	for key, value := range req.Capacity {
		quantity, err := resource.ParseQuantity(value)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("无效的容量 %s: %v", key, err))
			return
		}
		pv.Spec.Capacity[corev1.ResourceName(key)] = quantity
	}

	// 设置访问模式
	for _, mode := range req.AccessModes {
		pv.Spec.AccessModes = append(pv.Spec.AccessModes, corev1.PersistentVolumeAccessMode(mode))
	}

	// 设置回收策略
	if req.ReclaimPolicy != "" {
		reclaimPolicy := corev1.PersistentVolumeReclaimPolicy(req.ReclaimPolicy)
		pv.Spec.PersistentVolumeReclaimPolicy = reclaimPolicy
	}

	// 设置存储类
	if req.StorageClassName != "" {
		pv.Spec.StorageClassName = req.StorageClassName
	}

	// 设置卷模式
	if req.VolumeMode != "" {
		volumeMode := corev1.PersistentVolumeMode(req.VolumeMode)
		pv.Spec.VolumeMode = &volumeMode
	}

	// 设置挂载选项
	pv.Spec.MountOptions = req.MountOptions

	// 设置存储源
	if req.PersistentVolumeSource.HostPath != nil {
		pv.Spec.PersistentVolumeSource.HostPath = &corev1.HostPathVolumeSource{
			Path: req.PersistentVolumeSource.HostPath.Path,
			Type: (*corev1.HostPathType)(&req.PersistentVolumeSource.HostPath.Type),
		}
	}

	if req.PersistentVolumeSource.NFS != nil {
		pv.Spec.PersistentVolumeSource.NFS = &corev1.NFSVolumeSource{
			Server:   req.PersistentVolumeSource.NFS.Server,
			Path:     req.PersistentVolumeSource.NFS.Path,
			ReadOnly: req.PersistentVolumeSource.NFS.ReadOnly,
		}
	}

	// 创建PV
	createdPV, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建PV失败: "+err.Error())
		return
	}

	k8sPV := s.convertToK8sPV(createdPV)
	result.Success(c, k8sPV)
}

// UpdatePV 更新PV
func (s *K8sStorageServiceImpl) UpdatePV(c *gin.Context, clusterId uint, pvName string, req *model.UpdatePVRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有PV
	existingPV, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PV不存在: "+err.Error())
		return
	}

	// 更新标签
	if req.Labels != nil {
		existingPV.Labels = req.Labels
	}

	// 更新回收策略
	if req.ReclaimPolicy != "" {
		reclaimPolicy := corev1.PersistentVolumeReclaimPolicy(req.ReclaimPolicy)
		existingPV.Spec.PersistentVolumeReclaimPolicy = reclaimPolicy
	}

	// 更新挂载选项
	if req.MountOptions != nil {
		existingPV.Spec.MountOptions = req.MountOptions
	}

	// 更新PV
	updatedPV, err := clientset.CoreV1().PersistentVolumes().Update(context.TODO(), existingPV, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新PV失败: "+err.Error())
		return
	}

	k8sPV := s.convertToK8sPV(updatedPV)
	result.Success(c, k8sPV)
}

// DeletePV 删除PV
func (s *K8sStorageServiceImpl) DeletePV(c *gin.Context, clusterId uint, pvName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除PV失败: "+err.Error())
		return
	}

	result.Success(c, "PV删除成功")
}

// GetPVYaml 获取PV的YAML配置
func (s *K8sStorageServiceImpl) GetPVYaml(c *gin.Context, clusterId uint, pvName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pv, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "PV不存在: "+err.Error())
		return
	}

	// 清理不需要的字段，但保留 status
	pv.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if pv.APIVersion == "" {
		pv.APIVersion = "v1"
	}
	if pv.Kind == "" {
		pv.Kind = "PersistentVolume"
	}

	// 清理元数据字段
	pv.ResourceVersion = ""
	pv.UID = ""
	pv.Generation = 0
	pv.CreationTimestamp = metav1.Time{}
	pv.SelfLink = ""

	yamlData, err := yaml.Marshal(pv)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "转换YAML失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": string(yamlData),
	})
}

// UpdatePVYaml 通过YAML更新PV
func (s *K8sStorageServiceImpl) UpdatePVYaml(c *gin.Context, clusterId uint, pvName string, yamlData map[string]interface{}) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 将map转换为YAML字符串
	yamlBytes, err := json.Marshal(yamlData)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML数据格式错误: "+err.Error())
		return
	}

	var pv corev1.PersistentVolume
	err = yaml.Unmarshal(yamlBytes, &pv)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}

	// 确保名称正确
	pv.Name = pvName

	// 更新PV
	_, err = clientset.CoreV1().PersistentVolumes().Update(context.TODO(), &pv, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新PV失败: "+err.Error())
		return
	}

	result.Success(c, "PV YAML更新成功")
}

// ===================== StorageClass 管理实现 =====================

// GetStorageClasses 获取存储类列表
func (s *K8sStorageServiceImpl) GetStorageClasses(c *gin.Context, clusterId uint) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	storageClasses, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取存储类列表失败: "+err.Error())
		return
	}

	var scList []model.K8sStorageClass
	for _, sc := range storageClasses.Items {
		k8sSC := s.convertToK8sStorageClass(&sc)
		scList = append(scList, k8sSC)
	}

	response := model.StorageClassListResponse{
		StorageClasses: scList,
		Total:          len(scList),
	}

	result.Success(c, response)
}

// GetStorageClassDetail 获取存储类详情
func (s *K8sStorageServiceImpl) GetStorageClassDetail(c *gin.Context, clusterId uint, storageClassName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	sc, err := clientset.StorageV1().StorageClasses().Get(context.TODO(), storageClassName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "存储类不存在: "+err.Error())
		return
	}

	k8sSC := s.convertToK8sStorageClass(sc)

	// 获取相关事件
	events, _ := s.getStorageClassEvents(clientset, storageClassName)

	scDetail := model.StorageClassDetail{
		K8sStorageClass: k8sSC,
		Events:          events,
		Spec:            sc,
	}

	result.Success(c, scDetail)
}

// CreateStorageClass 创建存储类
func (s *K8sStorageServiceImpl) CreateStorageClass(c *gin.Context, clusterId uint, req *model.CreateStorageClassRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建StorageClass对象
	sc := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: req.Labels,
		},
		Provisioner:  req.Provisioner,
		Parameters:   req.Parameters,
		MountOptions: req.MountOptions,
	}

	// 设置回收策略
	if req.ReclaimPolicy != "" {
		reclaimPolicy := corev1.PersistentVolumeReclaimPolicy(req.ReclaimPolicy)
		sc.ReclaimPolicy = &reclaimPolicy
	}

	// 设置卷绑定模式
	if req.VolumeBindingMode != "" {
		volumeBindingMode := storagev1.VolumeBindingMode(req.VolumeBindingMode)
		sc.VolumeBindingMode = &volumeBindingMode
	}

	// 设置允许卷扩展
	sc.AllowVolumeExpansion = &req.AllowVolumeExpansion

	// 设置允许的拓扑
	if len(req.AllowedTopologies) > 0 {
		var topologies []corev1.TopologySelectorTerm
		for _, topology := range req.AllowedTopologies {
			term := corev1.TopologySelectorTerm{}
			for _, expr := range topology.MatchLabelExpressions {
				term.MatchLabelExpressions = append(term.MatchLabelExpressions, corev1.TopologySelectorLabelRequirement{
					Key:    expr.Key,
					Values: expr.Values,
				})
			}
			topologies = append(topologies, term)
		}
		sc.AllowedTopologies = topologies
	}

	// 创建StorageClass
	createdSC, err := clientset.StorageV1().StorageClasses().Create(context.TODO(), sc, metav1.CreateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建存储类失败: "+err.Error())
		return
	}

	k8sSC := s.convertToK8sStorageClass(createdSC)
	result.Success(c, k8sSC)
}

// UpdateStorageClass 更新存储类
func (s *K8sStorageServiceImpl) UpdateStorageClass(c *gin.Context, clusterId uint, storageClassName string, req *model.UpdateStorageClassRequest) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取现有StorageClass
	existingSC, err := clientset.StorageV1().StorageClasses().Get(context.TODO(), storageClassName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "存储类不存在: "+err.Error())
		return
	}

	// 更新标签
	if req.Labels != nil {
		existingSC.Labels = req.Labels
	}

	// 更新参数
	if req.Parameters != nil {
		existingSC.Parameters = req.Parameters
	}

	// 更新挂载选项
	if req.MountOptions != nil {
		existingSC.MountOptions = req.MountOptions
	}

	// 更新允许卷扩展
	existingSC.AllowVolumeExpansion = &req.AllowVolumeExpansion

	// 更新StorageClass
	updatedSC, err := clientset.StorageV1().StorageClasses().Update(context.TODO(), existingSC, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新存储类失败: "+err.Error())
		return
	}

	k8sSC := s.convertToK8sStorageClass(updatedSC)
	result.Success(c, k8sSC)
}

// DeleteStorageClass 删除存储类
func (s *K8sStorageServiceImpl) DeleteStorageClass(c *gin.Context, clusterId uint, storageClassName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.StorageV1().StorageClasses().Delete(context.TODO(), storageClassName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除存储类失败: "+err.Error())
		return
	}

	result.Success(c, "存储类删除成功")
}

// GetStorageClassYaml 获取存储类的YAML配置
func (s *K8sStorageServiceImpl) GetStorageClassYaml(c *gin.Context, clusterId uint, storageClassName string) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	sc, err := clientset.StorageV1().StorageClasses().Get(context.TODO(), storageClassName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "存储类不存在: "+err.Error())
		return
	}

	// 清理不需要的字段
	sc.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if sc.APIVersion == "" {
		sc.APIVersion = "storage.k8s.io/v1"
	}
	if sc.Kind == "" {
		sc.Kind = "StorageClass"
	}

	// 清理元数据字段
	sc.ResourceVersion = ""
	sc.UID = ""
	sc.Generation = 0
	sc.CreationTimestamp = metav1.Time{}
	sc.SelfLink = ""

	yamlData, err := yaml.Marshal(sc)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "转换YAML失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": string(yamlData),
	})
}

// UpdateStorageClassYaml 通过YAML更新存储类
func (s *K8sStorageServiceImpl) UpdateStorageClassYaml(c *gin.Context, clusterId uint, storageClassName string, yamlData map[string]interface{}) {
	clientset, err := s.getK8sClient(clusterId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 将map转换为YAML字符串
	yamlBytes, err := json.Marshal(yamlData)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML数据格式错误: "+err.Error())
		return
	}

	var sc storagev1.StorageClass
	err = yaml.Unmarshal(yamlBytes, &sc)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}

	// 确保名称正确
	sc.Name = storageClassName

	// 更新StorageClass
	_, err = clientset.StorageV1().StorageClasses().Update(context.TODO(), &sc, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新存储类失败: "+err.Error())
		return
	}

	result.Success(c, "存储类YAML更新成功")
}

// ===================== 辅助方法 =====================

// convertToK8sPVC 转换PVC对象
func (s *K8sStorageServiceImpl) convertToK8sPVC(pvc *corev1.PersistentVolumeClaim) model.K8sPersistentVolumeClaim {
	// 获取容量
	var capacity string
	if storage, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		capacity = storage.String()
	}

	// 获取访问模式
	var accessModes []string
	for _, mode := range pvc.Spec.AccessModes {
		accessModes = append(accessModes, string(mode))
	}

	// 获取存储类
	var storageClass string
	if pvc.Spec.StorageClassName != nil {
		storageClass = *pvc.Spec.StorageClassName
	}

	// 获取卷模式
	var volumeMode string
	if pvc.Spec.VolumeMode != nil {
		volumeMode = string(*pvc.Spec.VolumeMode)
	}

	return model.K8sPersistentVolumeClaim{
		Name:         pvc.Name,
		Namespace:    pvc.Namespace,
		Labels:       pvc.Labels,
		Capacity:     capacity,
		AccessModes:  accessModes,
		Status:       string(pvc.Status.Phase),
		StorageClass: storageClass,
		VolumeMode:   volumeMode,
		VolumeName:   pvc.Spec.VolumeName,
		CreatedAt:    pvc.CreationTimestamp.Format(time.RFC3339),
	}
}

// convertToK8sPV 转换PV对象
func (s *K8sStorageServiceImpl) convertToK8sPV(pv *corev1.PersistentVolume) model.K8sPersistentVolume {
	// 获取容量
	var capacity string
	if storage, ok := pv.Spec.Capacity[corev1.ResourceStorage]; ok {
		capacity = storage.String()
	}

	// 获取访问模式
	var accessModes []string
	for _, mode := range pv.Spec.AccessModes {
		accessModes = append(accessModes, string(mode))
	}

	// 获取卷模式
	var volumeMode string
	if pv.Spec.VolumeMode != nil {
		volumeMode = string(*pv.Spec.VolumeMode)
	}

	// 构建声明引用
	var claimRef *model.PVClaimRef
	if pv.Spec.ClaimRef != nil {
		claimRef = &model.PVClaimRef{
			Kind:       pv.Spec.ClaimRef.Kind,
			Namespace:  pv.Spec.ClaimRef.Namespace,
			Name:       pv.Spec.ClaimRef.Name,
			UID:        string(pv.Spec.ClaimRef.UID),
			APIVersion: pv.Spec.ClaimRef.APIVersion,
		}
	}

	// 构建存储源
	pvSource := model.PVSource{}
	if pv.Spec.PersistentVolumeSource.HostPath != nil {
		pvSource.HostPath = &model.PVHostPathVolumeSource{
			Path: pv.Spec.PersistentVolumeSource.HostPath.Path,
			Type: string(*pv.Spec.PersistentVolumeSource.HostPath.Type),
		}
	}
	if pv.Spec.PersistentVolumeSource.NFS != nil {
		pvSource.NFS = &model.PVNFSVolumeSource{
			Server:   pv.Spec.PersistentVolumeSource.NFS.Server,
			Path:     pv.Spec.PersistentVolumeSource.NFS.Path,
			ReadOnly: pv.Spec.PersistentVolumeSource.NFS.ReadOnly,
		}
	}

	return model.K8sPersistentVolume{
		Name:                   pv.Name,
		Labels:                 pv.Labels,
		Capacity:               capacity,
		AccessModes:            accessModes,
		ReclaimPolicy:          string(pv.Spec.PersistentVolumeReclaimPolicy),
		Status:                 string(pv.Status.Phase),
		StorageClass:           pv.Spec.StorageClassName,
		VolumeMode:             volumeMode,
		ClaimRef:               claimRef,
		PersistentVolumeSource: pvSource,
		MountOptions:           pv.Spec.MountOptions,
		CreatedAt:              pv.CreationTimestamp.Format(time.RFC3339),
	}
}

// convertToK8sStorageClass 转换StorageClass对象
func (s *K8sStorageServiceImpl) convertToK8sStorageClass(sc *storagev1.StorageClass) model.K8sStorageClass {
	// 获取回收策略
	var reclaimPolicy string
	if sc.ReclaimPolicy != nil {
		reclaimPolicy = string(*sc.ReclaimPolicy)
	}

	// 获取卷绑定模式
	var volumeBindingMode string
	if sc.VolumeBindingMode != nil {
		volumeBindingMode = string(*sc.VolumeBindingMode)
	}

	// 获取允许卷扩展
	var allowVolumeExpansion bool
	if sc.AllowVolumeExpansion != nil {
		allowVolumeExpansion = *sc.AllowVolumeExpansion
	}

	// 构建允许的拓扑
	var allowedTopologies []model.StorageClassTopology
	for _, topology := range sc.AllowedTopologies {
		scTopology := model.StorageClassTopology{}
		for _, expr := range topology.MatchLabelExpressions {
			scTopology.MatchLabelExpressions = append(scTopology.MatchLabelExpressions, model.StorageClassTopologyExp{
				Key:    expr.Key,
				Values: expr.Values,
			})
		}
		allowedTopologies = append(allowedTopologies, scTopology)
	}

	return model.K8sStorageClass{
		Name:                 sc.Name,
		Labels:               sc.Labels,
		Provisioner:          sc.Provisioner,
		Parameters:           sc.Parameters,
		ReclaimPolicy:        reclaimPolicy,
		VolumeBindingMode:    volumeBindingMode,
		AllowVolumeExpansion: allowVolumeExpansion,
		MountOptions:         sc.MountOptions,
		AllowedTopologies:    allowedTopologies,
		CreatedAt:            sc.CreationTimestamp.Format(time.RFC3339),
	}
}

// getPVCEvents 获取PVC相关事件
func (s *K8sStorageServiceImpl) getPVCEvents(clientset *kubernetes.Clientset, namespaceName string, pvcName string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=PersistentVolumeClaim", pvcName),
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

// getPVEvents 获取PV相关事件
func (s *K8sStorageServiceImpl) getPVEvents(clientset *kubernetes.Clientset, pvName string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=PersistentVolume", pvName),
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

// getStorageClassEvents 获取StorageClass相关事件
func (s *K8sStorageServiceImpl) getStorageClassEvents(clientset *kubernetes.Clientset, storageClassName string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=StorageClass", storageClassName),
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