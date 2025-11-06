package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type CmdbHostCloudServiceInterface interface {
	CreateAliyunHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
	CreateTencentHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
	CreateBaiduHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto)
}

type CmdbHostCloudServiceImpl struct {
	hostDao cmdbDao.CmdbHostDao
}

// 创建阿里云主机
func (s *CmdbHostCloudServiceImpl) CreateAliyunHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "阿里云主机导入请求已接收，正在后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.createAliyunHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 阿里云主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台阿里云主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入
		successCount := 0
		for _, host := range hosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 阿里云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 阿里云主机导入成功: %s\n", host.InstanceID)
		}

		// 记录最终结果
		fmt.Printf("[INFO] 阿里云主机导入完成: 成功 %d 台, 失败 %d 台\n", 
			successCount, len(hosts)-successCount)
		
		// 这里可以添加将结果写入Redis或数据库的逻辑，供前端查询
	}()
}

func (s *CmdbHostCloudServiceImpl) createAliyunHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	aliyunService := util.NewAliyunCloudService(dto.AccessKey, dto.SecretKey, dto.Region)
	
	instances, err := aliyunService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取阿里云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	var hosts []model.CmdbHost
	for _, instance := range instances {
		if strings.TrimSpace(instance.Region) == "" {
			instance.Region = "all"
		}

		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			continue
		}

		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		host := model.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     1,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      2,
			Region:      strings.TrimSpace(instance.Region),
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.MemoryGB),
			Disk:        strings.TrimSuffix(instance.DiskGB, "GB"),
			BillingType: instance.BillingType,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的实例可导入")
	}
	return hosts, nil
}

// 创建腾讯云主机
func (s *CmdbHostCloudServiceImpl) CreateTencentHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "success",
		"msg":    "腾讯云主机导入请求已接收",
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.createTencentHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 腾讯云主机导入失败: %v\n", err)
			return
		}

		for _, host := range hosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 腾讯云主机保存失败: %v\n", err)
				continue
			}
			fmt.Printf("[INFO] 腾讯云主机导入成功: %s\n", host.InstanceID)
		}
	}()
}

func (s *CmdbHostCloudServiceImpl) createTencentHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	tencentService := util.NewTencentCloudService(dto.AccessKey, dto.SecretKey)
	
	instances, err := tencentService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取腾讯云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	// 处理所有实例
	var hosts []model.CmdbHost
	for _, instance := range instances {
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			continue // 跳过无效IP的实例
		}

		totalDisk := instance.SystemDisk
		for _, disk := range instance.DataDisks {
			totalDisk += disk
		}

		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		host := model.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     1,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      3,
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: "包年包月",
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的实例可导入")
	}
	return hosts, nil
}

// 创建百度云主机
func (s *CmdbHostCloudServiceImpl) CreateBaiduHost(c *gin.Context, dto *model.CreateCmdbHostCloudDto) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "百度云主机导入请求已接收，正在扫描所有区域",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理导入
	go func() {
		// 初始化DAO
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.createBaiduHosts(dto)
		if err != nil {
			fmt.Printf("[ERROR] 百度云主机导入失败: %v\n", err)
			return
		}

		// 打印所有要导入的主机信息
		fmt.Printf("[DEBUG] 准备导入 %d 台百度云主机\n", len(hosts))
		for i, host := range hosts {
			fmt.Printf("[DEBUG] 主机 %d: %+v\n", i+1, host)
		}

		// 批量导入（自动去重）
		successCount := 0
		duplicateCount := 0
		for _, host := range hosts {
			// 检查是否已存在（根据InstanceID去重）
			existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
			if existing != nil {
				fmt.Printf("[INFO] 百度云主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
				duplicateCount++
				continue
			}

			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 百度云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 百度云主机导入成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		// 记录最终结果
		fmt.Printf("[INFO] 百度云主机导入完成: 成功 %d 台, 重复 %d 台, 失败 %d 台\n",
			successCount, duplicateCount, len(hosts)-successCount-duplicateCount)
	}()
}

func (s *CmdbHostCloudServiceImpl) createBaiduHosts(dto *model.CreateCmdbHostCloudDto) ([]model.CmdbHost, error) {
	baiduService := util.NewBaiduCloudService(dto.AccessKey, dto.SecretKey)

	// 获取所有区域的实例（自动扫描所有区域）
	instances, err := baiduService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取百度云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何百度云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台百度云实例\n", len(instances))

	var hosts []model.CmdbHost
	for _, instance := range instances {
		// 确定SSH连接IP
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			fmt.Printf("[WARN] 百度云实例无有效IP，跳过: %s\n", instance.InstanceID)
			continue
		}

		// 计算总磁盘容量
		totalDisk := instance.SystemDisk
		for _, disk := range instance.DataDisks {
			totalDisk += disk
		}

		// 解析时间
		createTime := time.Now()
		expireTime := time.Now().AddDate(1, 0, 0)

		if instance.CreateTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.CreateTime); err == nil {
				createTime = t
			}
		}
		if instance.ExpireTime != "" {
			if t, err := time.Parse("2006-01-02T15:04Z", instance.ExpireTime); err == nil {
				expireTime = t
			}
		}

		// 确定主机状态
		status := 2 // 默认在线
		if instance.Status == "Stopped" {
			status = 3 // 离线
		}

		host := model.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     1,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      4, // 百度云厂商代码为4
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      status,
			CPU:         strconv.Itoa(instance.CPU),
			Memory:      strconv.Itoa(instance.Memory),
			Disk:        strconv.Itoa(totalDisk),
			BillingType: instance.BillingType,
			CreateTime:  util.HTime{Time: createTime.UTC()},
			ExpireTime:  util.HTime{Time: expireTime.UTC()},
		}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return nil, errors.New("没有有效的百度云实例可导入")
	}

	fmt.Printf("[INFO] 成功解析 %d 台百度云主机\n", len(hosts))
	return hosts, nil
}

func GetCmdbHostCloudService() CmdbHostCloudServiceInterface {
	return &CmdbHostCloudServiceImpl{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}
