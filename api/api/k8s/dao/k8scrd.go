package dao

import (
	"fmt"

	"gorm.io/gorm"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type IK8sCRDDao interface {
	GetDynamicClient(clusterId uint) (dynamic.Interface, error)
	GetClientSet(clusterId uint) (*kubernetes.Clientset, error)
}

type k8sCRDDao struct {
	db *gorm.DB
}

func NewK8sCRDDao(db *gorm.DB) IK8sCRDDao {
	return &k8sCRDDao{db: db}
}

// 获取集群配置并返回 dynamic client
func (d *k8sCRDDao) GetDynamicClient(clusterId uint) (dynamic.Interface, error) {
	clusterDao := NewKubeClusterDao(d.db)
	cluster, err := clusterDao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("解析 kubeconfig 失败: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 dynamic client 失败: %v", err)
	}

	return dynClient, nil
}

// 获取集群配置并返回 clientset (用于获取系统 CRD 列表)
func (d *k8sCRDDao) GetClientSet(clusterId uint) (*kubernetes.Clientset, error) {
	clusterDao := NewKubeClusterDao(d.db)
	cluster, err := clusterDao.GetByID(clusterId)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("解析 kubeconfig 失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 clientset 失败: %v", err)
	}

	return clientset, nil
}