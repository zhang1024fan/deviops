package service

import (
	"context"
	"fmt"
	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
	"gorm.io/gorm"
)

type K8sConfigService struct {
	DB  *gorm.DB
	dao *dao.KubeClusterDao
}

func NewK8sConfigService(db *gorm.DB) *K8sConfigService {
	return &K8sConfigService{
		DB:  db,
		dao: dao.NewKubeClusterDao(db),
	}
}

// getK8sClient 获取K8s客户端
func (s *K8sConfigService) getK8sClient(clusterId uint) (*kubernetes.Clientset, error) {
	cluster, err := s.dao.GetByID(clusterId)
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

// ===================== ConfigMap 管理服务 =====================

// GetConfigMaps 获取ConfigMap列表
func (s *K8sConfigService) GetConfigMaps(clusterId int, namespaceName string, page, pageSize int) (*model.ConfigMapListResponse, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	configMaps, err := clientset.CoreV1().ConfigMaps(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取ConfigMap列表失败: %v", err)
	}

	var result []model.K8sConfigMap
	for _, cm := range configMaps.Items {
		configMap := s.convertToConfigMapModel(&cm)
		result = append(result, configMap)
	}

	// 简单分页处理
	total := len(result)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		result = []model.K8sConfigMap{}
	} else if end > total {
		result = result[start:]
	} else {
		result = result[start:end]
	}

	return &model.ConfigMapListResponse{
		ConfigMaps: result,
		Total:      total,
	}, nil
}

// GetConfigMapDetail 获取ConfigMap详情
func (s *K8sConfigService) GetConfigMapDetail(clusterId int, namespaceName, configMapName string) (*model.ConfigMapDetail, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取ConfigMap详情失败: %v", err)
	}

	configMap := s.convertToConfigMapModel(cm)

	// 获取相关事件
	events, _ := s.getConfigMapEvents(clientset, namespaceName, configMapName)

	// 获取使用情况
	usage, _ := s.getConfigMapUsage(clientset, namespaceName, configMapName)

	return &model.ConfigMapDetail{
		K8sConfigMap: configMap,
		Events:       events,
		Usage:        usage,
		Spec:         cm,
	}, nil
}

// CreateConfigMap 创建ConfigMap
func (s *K8sConfigService) CreateConfigMap(clusterId int, namespaceName string, req *model.CreateConfigMapRequest) (*model.K8sConfigMap, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
			Labels:    req.Labels,
		},
		Data:       req.Data,
		BinaryData: req.BinaryData,
		Immutable:  &req.Immutable,
	}

	createdCM, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("创建ConfigMap失败: %v", err)
	}

	result := s.convertToConfigMapModel(createdCM)
	return &result, nil
}

// UpdateConfigMap 更新ConfigMap
func (s *K8sConfigService) UpdateConfigMap(clusterId int, namespaceName, configMapName string, req *model.UpdateConfigMapRequest) (*model.K8sConfigMap, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 获取现有ConfigMap
	existingCM, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取ConfigMap失败: %v", err)
	}

	// 更新字段
	if req.Labels != nil {
		existingCM.Labels = req.Labels
	}
	if req.Data != nil {
		existingCM.Data = req.Data
	}
	if req.BinaryData != nil {
		existingCM.BinaryData = req.BinaryData
	}

	updatedCM, err := clientset.CoreV1().ConfigMaps(namespaceName).Update(context.TODO(), existingCM, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("更新ConfigMap失败: %v", err)
	}

	result := s.convertToConfigMapModel(updatedCM)
	return &result, nil
}

// DeleteConfigMap 删除ConfigMap
func (s *K8sConfigService) DeleteConfigMap(clusterId int, namespaceName, configMapName string) error {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	err = clientset.CoreV1().ConfigMaps(namespaceName).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除ConfigMap失败: %v", err)
	}

	return nil
}

// GetConfigMapYaml 获取ConfigMap YAML
func (s *K8sConfigService) GetConfigMapYaml(clusterId int, namespaceName, configMapName string) (string, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return "", fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("获取ConfigMap失败: %v", err)
	}

	// 清理不需要的字段
	s.cleanupConfigMapForYaml(cm)

	yamlBytes, err := yaml.Marshal(cm)
	if err != nil {
		return "", fmt.Errorf("转换YAML失败: %v", err)
	}

	return string(yamlBytes), nil
}

// UpdateConfigMapYaml 更新ConfigMap YAML
func (s *K8sConfigService) UpdateConfigMapYaml(clusterId int, namespaceName, configMapName string, yamlData map[string]interface{}) (*model.K8sConfigMap, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 将map转换为YAML字节
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("YAML数据格式错误: %v", err)
	}

	// 解析为ConfigMap对象
	var configMap corev1.ConfigMap
	err = yaml.Unmarshal(yamlBytes, &configMap)
	if err != nil {
		return nil, fmt.Errorf("YAML解析失败: %v", err)
	}

	// 确保名称和命名空间正确
	configMap.Name = configMapName
	configMap.Namespace = namespaceName

	updatedCM, err := clientset.CoreV1().ConfigMaps(namespaceName).Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("更新ConfigMap失败: %v", err)
	}

	result := s.convertToConfigMapModel(updatedCM)
	return &result, nil
}

// ===================== Secret 管理服务 =====================

// GetSecrets 获取Secret列表
func (s *K8sConfigService) GetSecrets(clusterId int, namespaceName string, page, pageSize int) (*model.SecretListResponse, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	secrets, err := clientset.CoreV1().Secrets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Secret列表失败: %v", err)
	}

	var result []model.K8sSecret
	for _, secret := range secrets.Items {
		secretModel := s.convertToSecretModel(&secret)
		result = append(result, secretModel)
	}

	// 简单分页处理
	total := len(result)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		result = []model.K8sSecret{}
	} else if end > total {
		result = result[start:]
	} else {
		result = result[start:end]
	}

	return &model.SecretListResponse{
		Secrets: result,
		Total:   total,
	}, nil
}

// GetSecretDetail 获取Secret详情
func (s *K8sConfigService) GetSecretDetail(clusterId int, namespaceName, secretName string) (*model.SecretDetail, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	secret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Secret详情失败: %v", err)
	}

	secretModel := s.convertToSecretModel(secret)

	// 获取相关事件
	events, _ := s.getSecretEvents(clientset, namespaceName, secretName)

	// 获取使用情况
	usage, _ := s.getSecretUsage(clientset, namespaceName, secretName)

	return &model.SecretDetail{
		K8sSecret: secretModel,
		Events:    events,
		Usage:     usage,
		Spec:      secret,
	}, nil
}

// CreateSecret 创建Secret
func (s *K8sConfigService) CreateSecret(clusterId int, namespaceName string, req *model.CreateSecretRequest) (*model.K8sSecret, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespaceName,
			Labels:    req.Labels,
		},
		Type:       corev1.SecretType(req.Type),
		Data:       req.Data,
		StringData: req.StringData,
		Immutable:  &req.Immutable,
	}

	createdSecret, err := clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("创建Secret失败: %v", err)
	}

	result := s.convertToSecretModel(createdSecret)
	return &result, nil
}

// UpdateSecret 更新Secret
func (s *K8sConfigService) UpdateSecret(clusterId int, namespaceName, secretName string, req *model.UpdateSecretRequest) (*model.K8sSecret, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 获取现有Secret
	existingSecret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取Secret失败: %v", err)
	}

	// 更新字段
	if req.Labels != nil {
		existingSecret.Labels = req.Labels
	}
	if req.Data != nil {
		existingSecret.Data = req.Data
	}
	if req.StringData != nil {
		existingSecret.StringData = req.StringData
	}

	updatedSecret, err := clientset.CoreV1().Secrets(namespaceName).Update(context.TODO(), existingSecret, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("更新Secret失败: %v", err)
	}

	result := s.convertToSecretModel(updatedSecret)
	return &result, nil
}

// DeleteSecret 删除Secret
func (s *K8sConfigService) DeleteSecret(clusterId int, namespaceName, secretName string) error {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	err = clientset.CoreV1().Secrets(namespaceName).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除Secret失败: %v", err)
	}

	return nil
}

// GetSecretYaml 获取Secret YAML
func (s *K8sConfigService) GetSecretYaml(clusterId int, namespaceName, secretName string) (string, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return "", fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	secret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("获取Secret失败: %v", err)
	}

	// 清理不需要的字段
	s.cleanupSecretForYaml(secret)

	yamlBytes, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("转换YAML失败: %v", err)
	}

	return string(yamlBytes), nil
}

// UpdateSecretYaml 更新Secret YAML
func (s *K8sConfigService) UpdateSecretYaml(clusterId int, namespaceName, secretName string, yamlData map[string]interface{}) (*model.K8sSecret, error) {
	clientset, err := s.getK8sClient(uint(clusterId))
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %v", err)
	}

	// 获取现有Secret以保留不可变字段
	existingSecret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取现有Secret失败: %v", err)
	}

	// 将map转换为YAML字节
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("YAML数据格式错误: %v", err)
	}

	// 解析为Secret对象
	var secret corev1.Secret
	err = yaml.Unmarshal(yamlBytes, &secret)
	if err != nil {
		return nil, fmt.Errorf("YAML解析失败: %v", err)
	}

	// 确保名称和命名空间正确
	secret.Name = secretName
	secret.Namespace = namespaceName

	// 保留不可变字段：type
	secret.Type = existingSecret.Type

	// 初始化 Data 字段（如果为 nil）
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	// 合并现有的 data 字段：保留现有数据，只更新/添加新数据
	for key, value := range existingSecret.Data {
		if _, exists := secret.Data[key]; !exists {
			// 如果新 YAML 中没有这个 key，保留旧值
			secret.Data[key] = value
		}
	}

	// 验证特定类型的Secret必需字段
	if err := s.validateSecretData(&secret); err != nil {
		return nil, err
	}

	// 保留resourceVersion以支持乐观锁
	secret.ResourceVersion = existingSecret.ResourceVersion

	updatedSecret, err := clientset.CoreV1().Secrets(namespaceName).Update(context.TODO(), &secret, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("更新Secret失败: %v", err)
	}

	result := s.convertToSecretModel(updatedSecret)
	return &result, nil
}

// ===================== 工具方法 =====================

// convertToConfigMapModel 转换ConfigMap为模型
func (s *K8sConfigService) convertToConfigMapModel(cm *corev1.ConfigMap) model.K8sConfigMap {
	return model.K8sConfigMap{
		Name:        cm.Name,
		Namespace:   cm.Namespace,
		Labels:      cm.Labels,
		Data:        cm.Data,
		BinaryData:  cm.BinaryData,
		Immutable:   cm.Immutable != nil && *cm.Immutable,
		CreatedTime: cm.CreationTimestamp.Format(time.RFC3339),
	}
}

// convertToSecretModel 转换Secret为模型
func (s *K8sConfigService) convertToSecretModel(secret *corev1.Secret) model.K8sSecret {
	// 隐藏敏感数据，只显示键名
	stringData := make(map[string]string)
	for key := range secret.Data {
		stringData[key] = "***"
	}

	return model.K8sSecret{
		Name:        secret.Name,
		Namespace:   secret.Namespace,
		Labels:      secret.Labels,
		Type:        string(secret.Type),
		Data:        secret.Data,
		StringData:  stringData,
		Immutable:   secret.Immutable != nil && *secret.Immutable,
		CreatedTime: secret.CreationTimestamp.Format(time.RFC3339),
	}
}

// cleanupConfigMapForYaml 清理ConfigMap的YAML输出
func (s *K8sConfigService) cleanupConfigMapForYaml(cm *corev1.ConfigMap) {
	// 清理不需要的字段
	cm.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if cm.APIVersion == "" {
		cm.APIVersion = "v1"
	}
	if cm.Kind == "" {
		cm.Kind = "ConfigMap"
	}

	// 清理元数据字段
	cm.ResourceVersion = ""
	cm.UID = ""
	cm.Generation = 0
	cm.CreationTimestamp = metav1.Time{}
	cm.SelfLink = ""
}

// cleanupSecretForYaml 清理Secret的YAML输出
func (s *K8sConfigService) cleanupSecretForYaml(secret *corev1.Secret) {
	// 清理不需要的字段
	secret.ManagedFields = nil

	// 确保 apiVersion 和 kind 字段存在
	if secret.APIVersion == "" {
		secret.APIVersion = "v1"
	}
	if secret.Kind == "" {
		secret.Kind = "Secret"
	}

	// 清理元数据字段
	secret.ResourceVersion = ""
	secret.UID = ""
	secret.Generation = 0
	secret.CreationTimestamp = metav1.Time{}
	secret.SelfLink = ""
}

// getConfigMapEvents 获取ConfigMap相关事件
func (s *K8sConfigService) getConfigMapEvents(clientset *kubernetes.Clientset, namespace, name string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=ConfigMap", name),
	})
	if err != nil {
		fmt.Printf("获取ConfigMap事件失败: %v", err)
		return nil, err
	}

	var result []model.K8sEvent
	for _, event := range events.Items {
		result = append(result, model.K8sEvent{
			Type:           event.Type,
			Reason:         event.Reason,
			Message:        event.Message,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
			Count:          event.Count,
		})
	}

	return result, nil
}

// getSecretEvents 获取Secret相关事件
func (s *K8sConfigService) getSecretEvents(clientset *kubernetes.Clientset, namespace, name string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Secret", name),
	})
	if err != nil {
		fmt.Printf("获取Secret事件失败: %v", err)
		return nil, err
	}

	var result []model.K8sEvent
	for _, event := range events.Items {
		result = append(result, model.K8sEvent{
			Type:           event.Type,
			Reason:         event.Reason,
			Message:        event.Message,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
			Count:          event.Count,
		})
	}

	return result, nil
}

// getConfigMapUsage 获取ConfigMap使用情况
func (s *K8sConfigService) getConfigMapUsage(clientset *kubernetes.Clientset, namespace, name string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var usage []string
	for _, pod := range pods.Items {
		// 检查Volume引用
		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil && volume.ConfigMap.Name == name {
				usage = append(usage, fmt.Sprintf("Pod: %s (Volume: %s)", pod.Name, volume.Name))
			}
		}

		// 检查环境变量引用
		for _, container := range pod.Spec.Containers {
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == name {
					usage = append(usage, fmt.Sprintf("Pod: %s, Container: %s (Env: %s)", pod.Name, container.Name, env.Name))
				}
			}
		}
	}

	return usage, nil
}

// getSecretUsage 获取Secret使用情况
func (s *K8sConfigService) getSecretUsage(clientset *kubernetes.Clientset, namespace, name string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var usage []string
	for _, pod := range pods.Items {
		// 检查Volume引用
		for _, volume := range pod.Spec.Volumes {
			if volume.Secret != nil && volume.Secret.SecretName == name {
				usage = append(usage, fmt.Sprintf("Pod: %s (Volume: %s)", pod.Name, volume.Name))
			}
		}

		// 检查环境变量引用
		for _, container := range pod.Spec.Containers {
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil && env.ValueFrom.SecretKeyRef.Name == name {
					usage = append(usage, fmt.Sprintf("Pod: %s, Container: %s (Env: %s)", pod.Name, container.Name, env.Name))
				}
			}
		}

		// 检查ImagePullSecrets引用
		for _, imagePullSecret := range pod.Spec.ImagePullSecrets {
			if imagePullSecret.Name == name {
				usage = append(usage, fmt.Sprintf("Pod: %s (ImagePullSecret)", pod.Name))
			}
		}
	}

	return usage, nil
}

// validateSecretData 验证Secret数据的完整性
func (s *K8sConfigService) validateSecretData(secret *corev1.Secret) error {
	// 辅助函数：检查字段是否存在于 Data 或 StringData 中
	hasField := func(key string) bool {
		if _, ok := secret.Data[key]; ok {
			return true
		}
		if _, ok := secret.StringData[key]; ok {
			return true
		}
		return false
	}

	switch secret.Type {
	case corev1.SecretTypeTLS:
		// TLS类型Secret必须包含tls.crt和tls.key
		if !hasField("tls.crt") {
			return fmt.Errorf("TLS类型Secret必须包含 tls.crt 字段")
		}
		if !hasField("tls.key") {
			return fmt.Errorf("TLS类型Secret必须包含 tls.key 字段")
		}
	case corev1.SecretTypeBasicAuth:
		// BasicAuth类型Secret必须包含username和password
		if !hasField("username") {
			return fmt.Errorf("BasicAuth类型Secret必须包含 username 字段")
		}
		if !hasField("password") {
			return fmt.Errorf("BasicAuth类型Secret必须包含 password 字段")
		}
	case corev1.SecretTypeDockerConfigJson:
		// DockerConfigJson类型Secret必须包含.dockerconfigjson
		if !hasField(".dockerconfigjson") {
			return fmt.Errorf("DockerConfigJson类型Secret必须包含 .dockerconfigjson 字段")
		}
	case corev1.SecretTypeSSHAuth:
		// SSHAuth类型Secret必须包含ssh-privatekey
		if !hasField("ssh-privatekey") {
			return fmt.Errorf("SSHAuth类型Secret必须包含 ssh-privatekey 字段")
		}
	}
	return nil
}