package util

import (
	"encoding/json"
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

// BaiduCloudService 提供百度云API调用服务
type BaiduCloudService struct {
	AccessKey    string
	AccessSecret string
}

// NewBaiduCloudService 创建新的百度云服务实例
func NewBaiduCloudService(accessKey, accessSecret string) *BaiduCloudService {
	return &BaiduCloudService{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
	}
}

// BaiduInstanceInfo 百度云实例信息结构
type BaiduInstanceInfo struct {
	InstanceID  string   `json:"instance_id"`
	PrivateIPs  []string `json:"private_ips"`
	PublicIPs   []string `json:"public_ips"`
	CPU         int      `json:"cpu"`
	Memory      int      `json:"memory"` // 单位: GB
	SystemDisk  int      `json:"system_disk"` // 单位: GB
	DataDisks   []int    `json:"data_disks"` // 单位: GB
	BillingType string   `json:"billing_type"`
	OSName      string   `json:"os_name"`
	CreateTime  string   `json:"create_time"`
	ExpireTime  string   `json:"expire_time"`
	Region      string   `json:"region"` // 实例所在区域
	Zone        string   `json:"zone"` // 可用区
	Status      string   `json:"status"` // 实例状态
}

// GetInstances 获取百度云实例列表（扫描所有区域）
func (s *BaiduCloudService) GetInstances() ([]BaiduInstanceInfo, error) {
	// 百度云支持的所有区域
	regions := []string{
		"bj",      // 北京
		"gz",      // 广州
		"su",      // 苏州
		"hkg",     // 香港
		"fwh",     // 武汉
		"bd",      // 保定
		"sin",     // 新加坡
		"fsh",     // 上海
	}

	var allInstances []BaiduInstanceInfo

	// 遍历所有区域获取实例
	for _, region := range regions {
		instances, err := s.GetInstancesByRegion(region)
		if err != nil {
			// 如果某个区域查询失败，记录日志但继续查询其他区域
			fmt.Printf("[WARN] 查询百度云区域 %s 失败: %v\n", region, err)
			continue
		}
		allInstances = append(allInstances, instances...)
	}

	return allInstances, nil
}

// GetInstancesByRegion 获取指定区域的百度云实例列表
func (s *BaiduCloudService) GetInstancesByRegion(region string) ([]BaiduInstanceInfo, error) {
	// 使用百度云官方SDK
	endpoint := fmt.Sprintf("bcc.%s.baidubce.com", region)

	// 创建BCC客户端
	client, err := bcc.NewClient(s.AccessKey, s.AccessSecret, endpoint)
	if err != nil {
		return nil, fmt.Errorf("创建BCC客户端失败: %v", err)
	}

	// 查询实例列表
	args := &api.ListInstanceArgs{}
	result, err := client.ListInstances(args)
	if err != nil {
		return nil, fmt.Errorf("查询实例列表失败: %v", err)
	}

	// 转换为统一格式
	var instances []BaiduInstanceInfo
	for _, inst := range result.Instances {
		privateIPs := []string{}
		if inst.InternalIP != "" {
			privateIPs = append(privateIPs, inst.InternalIP)
		}

		publicIPs := []string{}
		if inst.PublicIP != "" {
			publicIPs = append(publicIPs, inst.PublicIP)
		}

		billingType := "后付费"
		if inst.PaymentTiming == "Prepaid" {
			billingType = "预付费"
		}

		// 计算磁盘容量 - 使用CDS磁盘信息
		systemDisk := 40 // BCC默认系统盘大小
		var dataDisks []int

		instances = append(instances, BaiduInstanceInfo{
			InstanceID:  inst.InstanceId,
			PrivateIPs:  privateIPs,
			PublicIPs:   publicIPs,
			CPU:         inst.CpuCount,
			Memory:      inst.MemoryCapacityInGB,
			SystemDisk:  systemDisk,
			DataDisks:   dataDisks,
			BillingType: billingType,
			OSName:      inst.ImageId,
			CreateTime:  inst.CreationTime,
			ExpireTime:  inst.ExpireTime,
			Region:      region,
			Zone:        inst.ZoneName,
			Status:      string(inst.Status),
		})
	}

	fmt.Printf("[INFO] 从区域 %s 获取到 %d 台百度云实例\n", region, len(instances))
	return instances, nil
}

// GetInstancesJSON 获取百度云实例列表的JSON格式
func (s *BaiduCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}

// GetAllRegions 获取百度云支持的所有区域列表
func (s *BaiduCloudService) GetAllRegions() []string {
	return []string{
		"bj",      // 北京
		"gz",      // 广州
		"su",      // 苏州
		"hkg",     // 香港
		"fwh",     // 武汉
		"bd",      // 保定
		"sin",     // 新加坡
		"fsh",     // 上海
	}
}
