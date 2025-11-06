package util

import (
	"encoding/json"
)

// AliyunCloudService 提供阿里云API调用服务
type AliyunCloudService struct {
	AccessKey    string
	AccessSecret string
	Region      string
}

// NewAliyunCloudService 创建新的阿里云服务实例
func NewAliyunCloudService(accessKey, accessSecret, region string) *AliyunCloudService {
	return &AliyunCloudService{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
		Region:      region,
	}
}

// AliyunInstanceInfo 阿里云实例信息结构
type AliyunInstanceInfo struct {
	InstanceID  string   `json:"instance_id"`
	PrivateIPs  []string `json:"private_ips"`
	PublicIPs   []string `json:"public_ips"`
	CPU         int      `json:"cpu"`
	MemoryGB    int      `json:"memory_gb"`
	DiskGB      string   `json:"disk_gb"` // 注意: 这里应该返回纯数字，不带"GB"单位
	BillingType string   `json:"billing_type"`
	OSName      string   `json:"os_name"`
	CreateTime  string   `json:"create_time"`
	ExpireTime  string   `json:"expire_time"`
	Region      string   `json:"region"` // 实例所在区域
}

// GetInstances 获取阿里云实例列表
func (s *AliyunCloudService) GetInstances() ([]AliyunInstanceInfo, error) {
	// 1. 调用阿里云API获取所有区域的ECS实例
	// 这里应该实现实际的API调用逻辑，查询所有区域的主机
	// 示例数据 - 实际应该调用阿里云API
	instances := []AliyunInstanceInfo{
		{
			InstanceID:  "i-0jl9blax8btv92iklhbj",
			PrivateIPs:  []string{"172.20.236.121"},
			PublicIPs:   []string{"8.130.14.34"},
			CPU:         2,
			MemoryGB:    2,
			DiskGB:      "40",
			BillingType: "包年包月",
			OSName:      "Ubuntu 24.04 64位",
			CreateTime:  "2025-05-23T03:16Z",
			ExpireTime:  "2026-05-23T16:00Z",
			Region:      "cn-wulanchabu",
		},
		{
			InstanceID:  "i-2ndinstanceid123456",
			PrivateIPs:  []string{"172.20.236.122"},
			PublicIPs:   []string{"8.130.14.35"},
			CPU:         4,
			MemoryGB:    8,
			DiskGB:      "100",
			BillingType: "按量付费",
			OSName:      "CentOS 7.9 64位",
			CreateTime:  "2025-06-01T10:00Z",
			ExpireTime:  "2025-08-01T10:00Z", // 已过期
			Region:      "cn-hangzhou",
		},
		{
			InstanceID:  "i-3rdinstanceid789012",
			PrivateIPs:  []string{"172.20.236.123"},
			PublicIPs:   []string{"8.130.14.36"},
			CPU:         8,
			MemoryGB:    16,
			DiskGB:      "200",
			BillingType: "包年包月",
			OSName:      "Windows Server 2019",
			CreateTime:  "2025-04-15T08:30Z",
			ExpireTime:  "2025-07-15T08:30Z", // 已过期
			Region:      "cn-beijing",
		},
	}

	// 2. 确保每个实例都有region字段
	for i := range instances {
		if instances[i].Region == "" {
			instances[i].Region = "unknown"
		}
	}
	
	return instances, nil
}

// GetInstancesJSON 获取阿里云实例列表的JSON格式
func (s *AliyunCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}
