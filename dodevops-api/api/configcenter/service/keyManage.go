package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	aliyuncloud "dodevops-api/common/util"
	tengxuncloud "dodevops-api/common/util"
	baiducloud "dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type KeyManageService struct {
	dao *dao.KeyManageDao
}

func NewKeyManageService() *KeyManageService {
	return &KeyManageService{
		dao: dao.NewKeyManageDao(),
	}
}

// Create 创建密钥
func (s *KeyManageService) Create(keyManage *model.KeyManage) error {
	return s.dao.Create(keyManage)
}

// Update 更新密钥
func (s *KeyManageService) Update(keyManage *model.KeyManage) error {
	return s.dao.Update(keyManage)
}

// Delete 删除密钥
func (s *KeyManageService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID 根据ID查询密钥
func (s *KeyManageService) GetByID(id uint) (*model.KeyManage, error) {
	return s.dao.GetByID(id)
}

// List 获取密钥列表
func (s *KeyManageService) List() ([]model.KeyManage, error) {
	return s.dao.List()
}

// ListWithPage 获取密钥列表（分页）
func (s *KeyManageService) ListWithPage(page, pageSize int) ([]model.KeyManage, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// DecryptKeys 解密密钥信息
func (s *KeyManageService) DecryptKeys(id uint) (string, string, error) {
	keyManage, err := s.dao.GetByID(id)
	if err != nil {
		return "", "", err
	}
	return keyManage.DecryptKeys()
}

// GetByType 根据云厂商类型查询密钥
func (s *KeyManageService) GetByType(keyType int) ([]model.KeyManage, error) {
	return s.dao.GetByType(keyType)
}

// GetDecryptedKeyForCloudAPI 获取解密后的密钥信息用于调用云API
// 返回解密后的KeyID和KeySecret，专用于云API调用
func (s *KeyManageService) GetDecryptedKeyForCloudAPI(id uint) (keyID string, keySecret string, keyType int, err error) {
	keyManage, err := s.dao.GetByID(id)
	if err != nil {
		return "", "", 0, err
	}
	
	decryptedKeyID, decryptedKeySecret, err := keyManage.DecryptKeys()
	if err != nil {
		return "", "", 0, err
	}
	
	return decryptedKeyID, decryptedKeySecret, keyManage.KeyType, nil
}

// GetDecryptedKeyByType 根据云厂商类型获取第一个可用的解密密钥
// 用于自动选择对应云厂商的密钥调用API
func (s *KeyManageService) GetDecryptedKeyByType(keyType int) (keyID string, keySecret string, err error) {
	keyManages, err := s.dao.GetByType(keyType)
	if err != nil {
		return "", "", err
	}
	
	if len(keyManages) == 0 {
		return "", "", fmt.Errorf("未找到类型为 %d 的密钥配置", keyType)
	}
	
	// 使用第一个找到的密钥
	decryptedKeyID, decryptedKeySecret, err := keyManages[0].DecryptKeys()
	if err != nil {
		return "", "", err
	}
	
	return decryptedKeyID, decryptedKeySecret, nil
}

// SyncAliyunHosts 同步阿里云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncAliyunHosts(c *gin.Context, keyID uint, groupID uint, region string) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "阿里云主机同步请求已接收，正在后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.syncAliyunHosts(keyID, groupID, region)
		if err != nil {
			fmt.Printf("[ERROR] 阿里云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的阿里云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 阿里云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 阿里云主机同步成功: %s\n", host.InstanceID)
		}

		fmt.Printf("[INFO] 阿里云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncAliyunHostsBackground 后台同步阿里云主机（用于定时任务，不依赖gin.Context）
func (s *KeyManageService) SyncAliyunHostsBackground(keyID uint, groupID uint, region string) error {
	log.Printf("开始后台同步阿里云主机: keyID=%d, groupID=%d, region=%s", keyID, groupID, region)

	hostDao := cmdbDao.NewCmdbHostDao()

	hosts, err := s.syncAliyunHosts(keyID, groupID, region)
	if err != nil {
		log.Printf("阿里云主机同步失败: %v", err)
		return err
	}

	// 过滤掉已存在的主机
	newHosts := s.filterExistingHosts(hosts)
	log.Printf("准备同步 %d 台新的阿里云主机（总共获取 %d 台）", len(newHosts), len(hosts))

	// 批量导入新主机
	successCount := 0
	for _, host := range newHosts {
		if err := hostDao.CreateCmdbHost(&host); err != nil {
			log.Printf("阿里云主机保存失败: %v (实例ID: %s)", err, host.InstanceID)
			continue
		}
		successCount++
		log.Printf("阿里云主机同步成功: %s", host.InstanceID)
	}

	log.Printf("阿里云主机同步完成: 成功 %d 台，失败 %d 台", successCount, len(newHosts)-successCount)
	return nil
}

// syncAliyunHosts 获取阿里云主机信息
func (s *KeyManageService) syncAliyunHosts(keyID uint, groupID uint, region string) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}
	
	// 检查密钥类型是否为阿里云
	if keyType != 1 {
		return nil, fmt.Errorf("密钥类型不正确，期望阿里云(1)，实际为: %d", keyType)
	}

	aliyunService := aliyuncloud.NewAliyunCloudService(accessKey, secretKey, region)
	
	instances, err := aliyunService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取阿里云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	var hosts []cmdbModel.CmdbHost
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

		// 过滤已过期的实例
		if !expireTime.IsZero() && expireTime.Before(time.Now()) {
			fmt.Printf("[INFO] 跳过已过期实例: %s (过期时间: %s)\n", instance.InstanceID, expireTime.Format("2006-01-02 15:04:05"))
			continue
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      2, // 阿里云
			Region:      strings.TrimSpace(instance.Region),
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 未认证
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
		return nil, errors.New("没有有效的实例可同步")
	}
	return hosts, nil
}

// SyncTencentHosts 同步腾讯云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncTencentHosts(c *gin.Context, keyID uint, groupID uint) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "腾讯云主机同步请求已接收，正在后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()
		
		hosts, err := s.syncTencentHosts(keyID, groupID)
		if err != nil {
			fmt.Printf("[ERROR] 腾讯云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的腾讯云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 腾讯云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 腾讯云主机同步成功: %s\n", host.InstanceID)
		}

		fmt.Printf("[INFO] 腾讯云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncTencentHostsBackground 后台同步腾讯云主机（用于定时任务，不依赖gin.Context）
func (s *KeyManageService) SyncTencentHostsBackground(keyID uint, groupID uint) error {
	log.Printf("开始后台同步腾讯云主机: keyID=%d, groupID=%d", keyID, groupID)

	hostDao := cmdbDao.NewCmdbHostDao()

	hosts, err := s.syncTencentHosts(keyID, groupID)
	if err != nil {
		log.Printf("腾讯云主机同步失败: %v", err)
		return err
	}

	// 过滤掉已存在的主机
	newHosts := s.filterExistingHosts(hosts)
	log.Printf("准备同步 %d 台新的腾讯云主机（总共获取 %d 台）", len(newHosts), len(hosts))

	// 批量导入新主机
	successCount := 0
	for _, host := range newHosts {
		if err := hostDao.CreateCmdbHost(&host); err != nil {
			log.Printf("腾讯云主机保存失败: %v (实例ID: %s)", err, host.InstanceID)
			continue
		}
		successCount++
		log.Printf("腾讯云主机同步成功: %s", host.InstanceID)
	}

	log.Printf("腾讯云主机同步完成: 成功 %d 台，失败 %d 台", successCount, len(newHosts)-successCount)
	return nil
}

// syncTencentHosts 获取腾讯云主机信息
func (s *KeyManageService) syncTencentHosts(keyID uint, groupID uint) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}
	
	// 检查密钥类型是否为腾讯云
	if keyType != 2 {
		return nil, fmt.Errorf("密钥类型不正确，期望腾讯云(2)，实际为: %d", keyType)
	}

	tencentService := tengxuncloud.NewTencentCloudService(accessKey, secretKey)
	
	instances, err := tencentService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取腾讯云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到指定实例")
	}

	var hosts []cmdbModel.CmdbHost
	for _, instance := range instances {
		sshIP := ""
		if len(instance.PublicIPs) > 0 {
			sshIP = instance.PublicIPs[0]
		} else if len(instance.PrivateIPs) > 0 {
			sshIP = instance.PrivateIPs[0]
		} else {
			continue
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

		// 过滤已过期的实例
		if !expireTime.IsZero() && expireTime.Before(time.Now()) {
			fmt.Printf("[INFO] 跳过已过期实例: %s (过期时间: %s)\n", instance.InstanceID, expireTime.Format("2006-01-02 15:04:05"))
			continue
		}

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      3, // 腾讯云
			Region:      instance.Region,
			InstanceID:  instance.InstanceID,
			OS:          instance.OSName,
			Status:      2, // 未认证
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
		return nil, errors.New("没有有效的实例可同步")
	}
	return hosts, nil
}

// filterExistingHosts 过滤掉已存在的主机，避免重复导入（基于instance_id去重）
func (s *KeyManageService) filterExistingHosts(hosts []cmdbModel.CmdbHost) []cmdbModel.CmdbHost {
	var newHosts []cmdbModel.CmdbHost
	hostDao := cmdbDao.NewCmdbHostDao()

	for _, host := range hosts {
		// 通过instance_id来判断是否存在（更可靠的去重方式）
		existing, _ := hostDao.GetCmdbHostByInstanceID(host.InstanceID)
		if existing == nil {
			// 主机不存在，可以添加
			newHosts = append(newHosts, host)
		} else {
			fmt.Printf("[INFO] 主机已存在，跳过: %s (区域: %s)\n", host.InstanceID, host.Region)
		}
	}

	return newHosts
}

// isHostExists 检查主机是否已存在
func (s *KeyManageService) isHostExists(hostName string) bool {
	hostDao := cmdbDao.NewCmdbHostDao()
	return hostDao.CheckNameExists(hostName)
}

// SyncBaiduHosts 同步百度云主机（使用密钥管理中的凭据）
func (s *KeyManageService) SyncBaiduHosts(c *gin.Context, keyID uint, groupID uint) {
	// 立即返回成功响应
	result.Success(c, gin.H{
		"status": "processing",
		"msg":    "百度云主机同步请求已接收，正在扫描所有区域并后台处理",
		"data": gin.H{
			"task_id":    time.Now().UnixNano(),
			"start_time": time.Now().Format(time.RFC3339),
		},
	})

	// 异步处理同步
	go func() {
		hostDao := cmdbDao.NewCmdbHostDao()

		hosts, err := s.syncBaiduHosts(keyID, groupID)
		if err != nil {
			fmt.Printf("[ERROR] 百度云主机同步失败: %v\n", err)
			return
		}

		// 过滤掉已存在的主机
		newHosts := s.filterExistingHosts(hosts)
		fmt.Printf("[DEBUG] 准备同步 %d 台新的百度云主机（总共获取 %d 台）\n", len(newHosts), len(hosts))

		// 批量导入新主机
		successCount := 0
		for _, host := range newHosts {
			if err := hostDao.CreateCmdbHost(&host); err != nil {
				fmt.Printf("[ERROR] 百度云主机保存失败: %v (实例ID: %s)\n", err, host.InstanceID)
				continue
			}
			successCount++
			fmt.Printf("[INFO] 百度云主机同步成功: %s (区域: %s)\n", host.InstanceID, host.Region)
		}

		fmt.Printf("[INFO] 百度云主机同步完成: 成功 %d 台, 失败 %d 台, 跳过已存在 %d 台\n",
			successCount, len(newHosts)-successCount, len(hosts)-len(newHosts))
	}()
}

// SyncBaiduHostsBackground 后台同步百度云主机（用于定时任务，不依赖gin.Context）
func (s *KeyManageService) SyncBaiduHostsBackground(keyID uint, groupID uint) error {
	log.Printf("开始后台同步百度云主机: keyID=%d, groupID=%d", keyID, groupID)

	hostDao := cmdbDao.NewCmdbHostDao()

	hosts, err := s.syncBaiduHosts(keyID, groupID)
	if err != nil {
		log.Printf("百度云主机同步失败: %v", err)
		return err
	}

	// 过滤掉已存在的主机
	newHosts := s.filterExistingHosts(hosts)
	log.Printf("准备同步 %d 台新的百度云主机（总共获取 %d 台）", len(newHosts), len(hosts))

	// 批量导入新主机
	successCount := 0
	for _, host := range newHosts {
		if err := hostDao.CreateCmdbHost(&host); err != nil {
			log.Printf("百度云主机保存失败: %v (实例ID: %s)", err, host.InstanceID)
			continue
		}
		successCount++
		log.Printf("百度云主机同步成功: %s (区域: %s)", host.InstanceID, host.Region)
	}

	log.Printf("百度云主机同步完成: 成功 %d 台，失败 %d 台", successCount, len(newHosts)-successCount)
	return nil
}

// syncBaiduHosts 获取百度云主机信息（自动扫描所有区域）
func (s *KeyManageService) syncBaiduHosts(keyID uint, groupID uint) ([]cmdbModel.CmdbHost, error) {
	// 获取解密后的密钥
	accessKey, secretKey, keyType, err := s.GetDecryptedKeyForCloudAPI(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取密钥失败: %v", err)
	}

	// 检查密钥类型是否为百度云
	if keyType != 3 {
		return nil, fmt.Errorf("密钥类型不正确，期望百度云(3)，实际为: %d", keyType)
	}

	baiduService := baiducloud.NewBaiduCloudService(accessKey, secretKey)

	// 获取所有区域的实例（自动扫描所有区域）
	instances, err := baiduService.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("获取百度云实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, errors.New("未找到任何百度云实例")
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台百度云实例\n", len(instances))

	var hosts []cmdbModel.CmdbHost
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

		host := cmdbModel.CmdbHost{
			HostName:    instance.InstanceID,
			GroupID:     groupID,
			PrivateIP:   strings.Join(instance.PrivateIPs, ","),
			PublicIP:    strings.Join(instance.PublicIPs, ","),
			SSHIP:       sshIP,
			SSHPort:     22,
			SSHName:     "root",
			Vendor:      3, // 百度云厂商代码为3
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
		return nil, errors.New("没有有效的百度云实例可同步")
	}

	fmt.Printf("[INFO] 成功解析 %d 台百度云主机\n", len(hosts))
	return hosts, nil
}