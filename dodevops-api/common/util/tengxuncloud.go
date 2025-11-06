package util

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCloudService struct {
	AccessKey    string
	AccessSecret string
}

func NewTencentCloudService(accessKey, accessSecret string) *TencentCloudService {
	return &TencentCloudService{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
	}
}

type TencentInstanceInfo struct {
	InstanceID  string   `json:"instance_id"`
	PrivateIPs  []string `json:"private_ips"`
	PublicIPs   []string `json:"public_ips"`
	CPU         int      `json:"cpu"`
	Memory      int      `json:"memory"`
	SystemDisk  int      `json:"system_disk"`
	DataDisks   []int    `json:"data_disks"`
	OSName      string   `json:"os_name"`
	CreateTime  string   `json:"create_time"`
	ExpireTime  string   `json:"expire_time"`
	Region      string   `json:"region"`
}

func (s *TencentCloudService) GetInstances() ([]TencentInstanceInfo, error) {
	credential := common.NewCredential(s.AccessKey, s.AccessSecret)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

	// 获取所有可用区域
	regionClient, err := cvm.NewClient(credential, "ap-guangzhou", clientProfile)
	if err != nil {
		return nil, fmt.Errorf("创建区域查询客户端失败: %v", err)
	}

	regionRequest := cvm.NewDescribeRegionsRequest()
	regionResponse, err := regionClient.DescribeRegions(regionRequest)
	if err != nil {
		return nil, fmt.Errorf("查询可用区域失败: %v", err)
	}

	var allInstances []TencentInstanceInfo
	
	for _, regionInfo := range regionResponse.Response.RegionSet {
		if regionInfo.Region == nil {
			continue
		}
		regionName := *regionInfo.Region
		client, err := cvm.NewClient(credential, regionName, clientProfile)
		if err != nil {
			continue
		}

		request := cvm.NewDescribeInstancesRequest()
		response, err := client.DescribeInstances(request)
		if err != nil {
			continue
		}

		for _, instance := range response.Response.InstanceSet {
			privateIPs := make([]string, 0)
			publicIPs := make([]string, 0)
			
			if instance.PrivateIpAddresses != nil {
				for _, ip := range instance.PrivateIpAddresses {
					if ip != nil {
						privateIPs = append(privateIPs, *ip)
					}
				}
			}
			if instance.PublicIpAddresses != nil {
				for _, ip := range instance.PublicIpAddresses {
					if ip != nil {
						publicIPs = append(publicIPs, *ip)
					}
				}
			}

			// 获取系统盘大小
			systemDiskSize := 0
			if instance.SystemDisk != nil && instance.SystemDisk.DiskSize != nil {
				systemDiskSize = int(*instance.SystemDisk.DiskSize)
			}

			// 获取数据盘大小列表
			dataDiskSizes := make([]int, 0)
			if instance.DataDisks != nil {
				for _, disk := range instance.DataDisks {
					if disk.DiskSize != nil {
						dataDiskSizes = append(dataDiskSizes, int(*disk.DiskSize))
					}
				}
			}

			instanceInfo := TencentInstanceInfo{
				PrivateIPs:   privateIPs,
				PublicIPs:    publicIPs,
				Region:       regionName,
				SystemDisk:   systemDiskSize,
				DataDisks:    dataDiskSizes,
			}

			// 安全处理可能为nil的字段
			if instance.InstanceId != nil {
				instanceInfo.InstanceID = *instance.InstanceId
			}
			if instance.CPU != nil {
				instanceInfo.CPU = int(*instance.CPU)
			}
			if instance.Memory != nil {
				instanceInfo.Memory = int(*instance.Memory)
			}
			if instance.CreatedTime != nil {
				instanceInfo.CreateTime = *instance.CreatedTime
			}
			if instance.ExpiredTime != nil {
				instanceInfo.ExpireTime = *instance.ExpiredTime
			}
			if instance.OsName != nil {
				instanceInfo.OSName = *instance.OsName
			}

			allInstances = append(allInstances, instanceInfo)
		}
	}

	return allInstances, nil
}

func (s *TencentCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}
