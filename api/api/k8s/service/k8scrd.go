package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"dodevops-api/api/k8s/dao"

	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

type IK8sCRDService interface {
	GetCRDGroups(clusterId uint) ([]string, error)
	GetCRDList(clusterId uint, params map[string]string) ([]map[string]interface{}, error)
	GetCustomResourceList(clusterId uint, namespaceName, crdName string, params map[string]string) (map[string]interface{}, error)
	GetCustomResourceDetail(clusterId uint, namespaceName, crdName, crName string) (*unstructured.Unstructured, error)
	CreateCustomResource(clusterId uint, namespaceName, crdName string, data map[string]interface{}) (*unstructured.Unstructured, error)
	DeleteCustomResource(clusterId uint, namespaceName, crdName, crName string) error
	GetCustomResourceYaml(clusterId uint, namespaceName, crdName, crName string) (string, error)
	UpdateCustomResourceYaml(clusterId uint, namespaceName, crdName, crName, yamlContent string) (*unstructured.Unstructured, error)
}

type k8sCRDService struct {
	dao dao.IK8sCRDDao
}

func NewK8sCRDService(db *gorm.DB) IK8sCRDService {
	return &k8sCRDService{
		dao: dao.NewK8sCRDDao(db),
	}
}

// 解析 CRD 名字以获取 Group, Resource
// crdName 通常形如 "prometheusrules.monitoring.coreos.com"
func parseCRDName(crdName string) (group, resource string) {
	parts := strings.SplitN(crdName, ".", 2)
	if len(parts) >= 2 {
		return parts[1], parts[0]
	}
	return "", parts[0]
}

// 模拟获取 GroupResourceVersion（假设 v1 可用，不准确但通常需要通过 API 发现，这里简化处理。最好的方式是用 Discovery 获取准确版本）
func (s *k8sCRDService) getGVR(clusterId uint, crdName string) (schema.GroupVersionResource, error) {
	group, resource := parseCRDName(crdName)
	// 动态获取推荐的 version
	clientset, err := s.dao.GetClientSet(clusterId)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	// 这里通过 Discovery 获取推荐版本比较复杂，简化：我们在知道大部分情况下 v1 或 v1beta1 是首选。
	// 这里简单默认使用 API group 发现的第一版，或者硬编码。真实情况可以使用 discoveryClient 获取。
	version := "v1" // 默认退化
	groups, err := clientset.Discovery().ServerGroups()
	if err == nil {
		for _, g := range groups.Groups {
			if g.Name == group {
				version = g.PreferredVersion.Version
				break
			}
		}
	}

	return schema.GroupVersionResource{Group: group, Version: version, Resource: resource}, nil
}

func (s *k8sCRDService) GetCRDGroups(clusterId uint) ([]string, error) {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	list, err := dynClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// 尝试 v1beta1
		gvr.Version = "v1beta1"
		list, err = dynClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("列出 CRD 失败: %v", err)
		}
	}

	groupMap := make(map[string]bool)
	for _, item := range list.Items {
		spec, _, _ := unstructured.NestedMap(item.Object, "spec")
		group, _, _ := unstructured.NestedString(spec, "group")
		if group != "" {
			groupMap[group] = true
		}
	}

	var groups []string
	for group := range groupMap {
		groups = append(groups, group)
	}
	return groups, nil
}

func (s *k8sCRDService) GetCRDList(clusterId uint, params map[string]string) ([]map[string]interface{}, error) {
	// k8s 内置的 CRD 通常属于 apiextensions.k8s.io 组
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	list, err := dynClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// 尝试 v1beta1
		gvr.Version = "v1beta1"
		list, err = dynClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("列出 CRD 失败: %v", err)
		}
	}

	var results []map[string]interface{}
	targetGroup := params["group"]

	for _, item := range list.Items {
		spec, _, _ := unstructured.NestedMap(item.Object, "spec")
		group, _, _ := unstructured.NestedString(spec, "group")

		if targetGroup != "" && group != targetGroup {
			continue
		}

		names, _, _ := unstructured.NestedMap(spec, "names")
		kind, _, _ := unstructured.NestedString(names, "kind")
		plural, _, _ := unstructured.NestedString(names, "plural")
		scope, _, _ := unstructured.NestedString(spec, "scope")

		results = append(results, map[string]interface{}{
			"name":              item.GetName(),
			"group":             group,
			"kind":              kind,
			"plural":            plural,
			"scope":             scope,
			"creationTimestamp": item.GetCreationTimestamp().Time,
		})
	}
	return results, nil
}

func (s *k8sCRDService) GetCustomResourceList(clusterId uint, namespaceName, crdName string, params map[string]string) (map[string]interface{}, error) {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr, err := s.getGVR(clusterId, crdName)
	if err != nil {
		return nil, err
	}

	listOptions := metav1.ListOptions{}
	if labelsStr, ok := params["labels"]; ok && labelsStr != "" {
		listOptions.LabelSelector = labelsStr
	}

	var list *unstructured.UnstructuredList
	if namespaceName != "" && namespaceName != "all" {
		list, err = dynClient.Resource(gvr).Namespace(namespaceName).List(context.Background(), listOptions)
	} else {
		list, err = dynClient.Resource(gvr).List(context.Background(), listOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("列出自定义资源失败: %v", err)
	}

	var results []map[string]interface{}
	for _, item := range list.Items {
		name := item.GetName()
		kind := item.GetKind()

		// Filter
		match := true
		if exactName, ok := params["name"]; ok && exactName != "" {
			if name != exactName {
				match = false
			}
		}
		if keyword, ok := params["keyword"]; ok && keyword != "" {
			if !strings.Contains(name, keyword) {
				match = false
			}
		}
		if exactKind, ok := params["kind"]; ok && exactKind != "" {
			// CRD kind is uniform usually, but just in case
			if kind != exactKind {
				match = false
			}
		}

		if match {
			results = append(results, item.Object)
		}
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	total := len(results)

	pageStr := params["page"]
	pageSizeStr := params["pageSize"]
	if pageStr != "" && pageSizeStr != "" {
		page, _ := strconv.Atoi(pageStr)
		pageSize, _ := strconv.Atoi(pageSizeStr)
		if page > 0 && pageSize > 0 {
			start := (page - 1) * pageSize
			if start > total {
				results = []map[string]interface{}{}
			} else {
				end := start + pageSize
				if end > total {
					end = total
				}
				results = results[start:end]
			}
		}
	}

	return map[string]interface{}{
		"items": results,
		"total": total,
	}, nil
}

func (s *k8sCRDService) GetCustomResourceDetail(clusterId uint, namespaceName, crdName, crName string) (*unstructured.Unstructured, error) {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr, err := s.getGVR(clusterId, crdName)
	if err != nil {
		return nil, err
	}

	var item *unstructured.Unstructured
	if namespaceName != "" && namespaceName != "all" {
		item, err = dynClient.Resource(gvr).Namespace(namespaceName).Get(context.Background(), crName, metav1.GetOptions{})
	} else {
		item, err = dynClient.Resource(gvr).Get(context.Background(), crName, metav1.GetOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("获取自定义资源详情失败: %v", err)
	}
	// 去除 managedFields 减小体积
	item.SetManagedFields(nil)
	return item, nil
}

func (s *k8sCRDService) CreateCustomResource(clusterId uint, namespaceName, crdName string, data map[string]interface{}) (*unstructured.Unstructured, error) {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr, err := s.getGVR(clusterId, crdName)
	if err != nil {
		return nil, err
	}

	var finalData map[string]interface{}
	if yamlContent, ok := data["yamlContent"].(string); ok && yamlContent != "" {
		err = yaml.Unmarshal([]byte(yamlContent), &finalData)
		if err != nil {
			return nil, fmt.Errorf("解析 YAML 失败: %v", err)
		}
	} else {
		finalData = data
	}

	obj := &unstructured.Unstructured{Object: finalData}

	var created *unstructured.Unstructured
	if namespaceName != "" && namespaceName != "all" {
		created, err = dynClient.Resource(gvr).Namespace(namespaceName).Create(context.Background(), obj, metav1.CreateOptions{})
	} else {
		created, err = dynClient.Resource(gvr).Create(context.Background(), obj, metav1.CreateOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("创建自定义资源失败: %v", err)
	}
	return created, nil
}

func (s *k8sCRDService) DeleteCustomResource(clusterId uint, namespaceName, crdName, crName string) error {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return err
	}

	gvr, err := s.getGVR(clusterId, crdName)
	if err != nil {
		return err
	}

	if namespaceName != "" && namespaceName != "all" {
		err = dynClient.Resource(gvr).Namespace(namespaceName).Delete(context.Background(), crName, metav1.DeleteOptions{})
	} else {
		err = dynClient.Resource(gvr).Delete(context.Background(), crName, metav1.DeleteOptions{})
	}

	if err != nil {
		return fmt.Errorf("删除自定义资源失败: %v", err)
	}
	return nil
}

func (s *k8sCRDService) GetCustomResourceYaml(clusterId uint, namespaceName, crdName, crName string) (string, error) {
	item, err := s.GetCustomResourceDetail(clusterId, namespaceName, crdName, crName)
	if err != nil {
		return "", err
	}

	// 移除一些通常不必要展示的状态和元数据
	unstructured.RemoveNestedField(item.Object, "metadata", "managedFields")
	// 可以选择是否移除 status
	// unstructured.RemoveNestedField(item.Object, "status")

	yamlBytes, err := yaml.Marshal(item.Object)
	if err != nil {
		return "", fmt.Errorf("格式化 YAML 失败: %v", err)
	}
	return string(yamlBytes), nil
}

func (s *k8sCRDService) UpdateCustomResourceYaml(clusterId uint, namespaceName, crdName, crName, yamlContent string) (*unstructured.Unstructured, error) {
	dynClient, err := s.dao.GetDynamicClient(clusterId)
	if err != nil {
		return nil, err
	}

	gvr, err := s.getGVR(clusterId, crdName)
	if err != nil {
		return nil, err
	}

	// 先将其转换为 JSON map
	var obj map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlContent), &obj)
	if err != nil {
		return nil, fmt.Errorf("YAML 解析为 JSON 错误: %v", err)
	}

	unstructObj := &unstructured.Unstructured{Object: obj}

	// 获取现有版本，为了获取 resourceVersion 这是 Update 必须的
	var existing *unstructured.Unstructured
	if namespaceName != "" && namespaceName != "all" {
		existing, err = dynClient.Resource(gvr).Namespace(namespaceName).Get(context.Background(), crName, metav1.GetOptions{})
	} else {
		existing, err = dynClient.Resource(gvr).Get(context.Background(), crName, metav1.GetOptions{})
	}
	if err != nil {
		return nil, fmt.Errorf("获取现有自定义资源失败: %v", err)
	}

	unstructObj.SetResourceVersion(existing.GetResourceVersion())

	var updated *unstructured.Unstructured
	if namespaceName != "" && namespaceName != "all" {
		updated, err = dynClient.Resource(gvr).Namespace(namespaceName).Update(context.Background(), unstructObj, metav1.UpdateOptions{})
	} else {
		updated, err = dynClient.Resource(gvr).Update(context.Background(), unstructObj, metav1.UpdateOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("更新自定义资源失败: %v", err)
	}
	return updated, nil
}