package service

import (
	"fmt"
	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	configDao "dodevops-api/api/configcenter/dao"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"time"

	"github.com/gin-gonic/gin"
)

type CmdbHostServiceInterface interface {
	GetCmdbHostList(c *gin.Context)                                                                           // 获取主机列表
	GetCmdbHostListWithPage(c *gin.Context, page, pageSize int)                                               // 获取主机列表(分页)
	GetCmdbHostById(c *gin.Context, id uint)                                                                  // 根据ID获取主机
	GetCmdbHostByName(c *gin.Context, name string)                                                            // 根据名称获取主机
	CreateCmdbHost(c *gin.Context, dto *model.CreateCmdbHostDto)                                              // 创建主机
	UpdateCmdbHost(c *gin.Context, id uint, dto *model.UpdateCmdbHostDto)                                     // 更新主机
	DeleteCmdbHost(c *gin.Context, id uint)                                                                   // 删除主机
	GetCmdbHostsByGroupId(c *gin.Context, groupId uint)                                                       // 根据分组ID获取主机列表
	GetCmdbHostsByHostNameLike(c *gin.Context, name string)                                                   // 根据主机名称模糊查询
	GetCmdbHostsByIP(c *gin.Context, ip string)                                                               // 根据IP查询(内网/公网/SSH)
	GetCmdbHostsByStatus(c *gin.Context, status int)                                                          // 根据状态查询
	ImportHostsFromExcel(c *gin.Context, dto *model.ImportHostsFromExcelDto, hosts []model.ExcelHostTemplate) // 从Excel导入主机
	SyncHostInfo(c *gin.Context, id uint)                                                                     // 同步主机基本信息
}

type CmdbHostServiceImpl struct {
	dao      cmdbDao.CmdbHostDao
	groupDao cmdbDao.CmdbGroupDao
}

// 从Excel导入主机
func (s *CmdbHostServiceImpl) ImportHostsFromExcel(c *gin.Context, dto *model.ImportHostsFromExcelDto, hosts []model.ExcelHostTemplate) {
	// 检查分组是否存在
	group, err := s.groupDao.GetCmdbGroupById(dto.GroupID)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "分组不存在")
		return
	}
	fmt.Printf("导入主机到分组: ID=%d, Name=%s\n", group.ID, group.Name)

	// 批量创建主机
	var successCount, failCount int
	var failedHosts []string
	type hostResult struct {
		hostID uint
		err    error
	}
	resultChan := make(chan hostResult, len(hosts))

	for i, excelHost := range hosts {
		fmt.Printf("处理第%d条主机记录: %+v\n", i+1, excelHost)

		// 检查SSH IP是否已存在
		if existHost := s.dao.GetCmdbHostBySSHIP(excelHost.SSHIP); existHost != nil {
			fmt.Printf("警告: SSH地址 %s 已存在(主机别名: %s)，跳过导入\n", excelHost.SSHIP, existHost.HostName)
			failedHosts = append(failedHosts, fmt.Sprintf("%s(IP重复)", excelHost.HostAlias))
			failCount++
			continue
		}

		// 创建主机
		host := model.CmdbHost{
			HostName:   excelHost.HostAlias, // 使用主机别名作为主机名称
			Name:       excelHost.HostAlias, // 默认name也使用主机别名，后续可通过同步更新
			GroupID:    dto.GroupID,
			SSHIP:      excelHost.SSHIP,
			SSHName:    excelHost.SSHName,
			SSHPort:    excelHost.SSHPort,
			Remark:     excelHost.Remark,
			CreateTime: util.HTime{Time: time.Now()},
			Vendor:     1, // 默认创建主机都是为自建主机
			Status:     2, // 初始状态设为未认证
		}

		// 保存主机信息
		fmt.Printf("尝试创建主机: %+v\n", host)
		if err := s.dao.CreateCmdbHost(&host); err != nil {
			fmt.Printf("创建主机失败: %v\n主机信息: %+v\n", err, host)
			failedHosts = append(failedHosts, excelHost.HostAlias)
			failCount++
			continue
		}
		fmt.Printf("主机创建成功, ID=%d\n", host.ID)
		successCount++

		// 异步执行SSH采集（使用Type=3公钥免认证）
		go func(h model.CmdbHost) {
			fmt.Printf("开始SSH采集(公钥免认证): %s@%s:%d\n", h.SSHName, h.SSHIP, h.SSHPort)

			// 准备SSH配置 - 使用Type=3公钥免认证
			sshConfig := util.SSHConfig{
				IP:       h.SSHIP,
				Port:     h.SSHPort,
				Type:     3, // 公钥免认证
				Username: h.SSHName,
			}

			// 获取系统信息
			sshUtil := util.NewSSHUtil()
			systemInfo, err := sshUtil.GetSystemInfo(&sshConfig)
			if err != nil {
				// 更新状态为认证失败
				fmt.Printf("主机ID %d SSH连接失败(公钥免认证): %v\n", h.ID, err)
				s.dao.UpdateCmdbHost(h.ID, &model.CmdbHost{Status: 3})
				resultChan <- hostResult{hostID: h.ID, err: err}
				return
			}

			// 更新主机信息
			updateData := model.CmdbHost{
				PrivateIP:  systemInfo["privateIp"],
				PublicIP:   systemInfo["publicIp"],
				Name:       systemInfo["name"],
				OS:         systemInfo["os"],
				CPU:        systemInfo["cpu"],
				Memory:     systemInfo["memory"],
				Disk:       systemInfo["disk"],
				Status:     1, // 认证成功
				UpdateTime: util.HTime{Time: time.Now()},
			}
			if err := s.dao.UpdateCmdbHost(h.ID, &updateData); err != nil {
				fmt.Printf("主机ID %d 更新信息失败: %v\n", h.ID, err)
				resultChan <- hostResult{hostID: h.ID, err: err}
				return
			}

			fmt.Printf("主机ID %d 信息采集成功\n", h.ID)
			resultChan <- hostResult{hostID: h.ID, err: nil}
		}(host)
	}

	// 等待所有goroutine完成
	go func() {
		for i := 0; i < successCount; i++ {
			res := <-resultChan
			if res.err != nil {
				fmt.Printf("主机ID %d 采集失败: %v\n", res.hostID, res.err)
			}
		}
		close(resultChan)
	}()

	responseData := gin.H{
		"success": successCount,
		"fail":    failCount,
		"total":   len(hosts),
		"message": "主机导入成功，正在使用公钥免认证方式采集主机信息...",
	}

	if len(failedHosts) > 0 {
		responseData["failedHosts"] = failedHosts
	}

	result.Success(c, responseData)
}

// 获取主机列表(分页)
func (s *CmdbHostServiceImpl) GetCmdbHostListWithPage(c *gin.Context, page, pageSize int) {
	list, total := s.dao.GetCmdbHostListWithPage(page, pageSize)
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// 获取主机列表
func (s *CmdbHostServiceImpl) GetCmdbHostList(c *gin.Context) {
	list := s.dao.GetCmdbHostList()
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 创建主机
func (s *CmdbHostServiceImpl) CreateCmdbHost(c *gin.Context, dto *model.CreateCmdbHostDto) {
	// 检查名称是否已存在
	if s.dao.CheckNameExists(dto.HostName) {
		result.FailedWithCode(c, constant.CMDB_HOST_NAME_EXISTS, "主机名称已存在")
		return
	}

	// 获取SSH凭据 (前端已确保SSHKeyID有效)
	authDao := configDao.NewEcsAuthDao()
	auth, _ := authDao.GetEcsAuthById(dto.SSHKeyID)

	// 初始保存连接信息
	host := model.CmdbHost{
		HostName:   dto.HostName,
		GroupID:    dto.GroupID,
		SSHIP:      dto.SSHIP,
		SSHName:    dto.SSHName,
		SSHKeyID:   dto.SSHKeyID,
		SSHPort:    dto.SSHPort,
		Remark:     dto.Remark,
		CreateTime: util.HTime{Time: time.Now()},
		Vendor:     1, // 默认创建主机都是为自建主机
		Status:     2, // 初始状态设为未认证
	}

	// 先保存基本信息
	if err := s.dao.CreateCmdbHost(&host); err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_CREATE_FAILED, err.Error())
		return
	}

	// 立即返回成功响应，后台异步执行SSH操作
	go func() {
		// 准备SSH配置
		sshConfig := util.SSHConfig{
			IP:        dto.SSHIP,
			Port:      dto.SSHPort,
			Type:      auth.Type,
			Username:  dto.SSHName,
			Password:  auth.Password,
			PublicKey: auth.PublicKey,
		}

		// 获取系统信息
		fmt.Println("开始尝试SSH连接获取系统信息...")
		fmt.Printf("SSH配置: %+v\n", sshConfig)

		sshUtil := util.NewSSHUtil()
		systemInfo, err := sshUtil.GetSystemInfo(&sshConfig)
		if err != nil {
			fmt.Printf("SSH获取系统信息失败: %v\n", err)
			// 更新状态为认证失败
			s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
			return
		}

		fmt.Printf("成功获取系统信息: %+v\n", systemInfo)

		// 验证必要字段是否存在
		if systemInfo["privateIp"] == "" || systemInfo["os"] == "" {
			fmt.Println("警告: 获取的系统信息不完整")
		}

		// 更新主机信息
		updateData := model.CmdbHost{
			PrivateIP:  systemInfo["privateIp"],
			PublicIP:   systemInfo["publicIp"],
			Name:       systemInfo["name"], // 添加name字段
			OS:         systemInfo["os"],
			CPU:        systemInfo["cpu"],
			Memory:     systemInfo["memory"],
			Disk:       systemInfo["disk"],
			Status:     1, // 认证成功
			UpdateTime: util.HTime{Time: time.Now()},
		}
		s.dao.UpdateCmdbHost(host.ID, &updateData)
	}()

	// 返回成功响应，前端可以通过轮询获取最新状态
	result.Success(c, gin.H{
		"id":     host.ID,
		"status": host.Status,
		"msg":    "主机创建成功，系统信息收集中...",
	})
}

// 更新主机
func (s *CmdbHostServiceImpl) UpdateCmdbHost(c *gin.Context, id uint, dto *model.UpdateCmdbHostDto) {
	// 不再需要查询认证凭证信息，直接从dto获取SSHName和SSHPort

	host := model.CmdbHost{
		HostName: dto.HostName,
		GroupID:  dto.GroupID,
		SSHIP:    dto.SSHIP,
		SSHName:  dto.SSHName,
		SSHKeyID: dto.SSHKeyID,
		SSHPort:  dto.SSHPort,
		Vendor:   dto.Vendor,
		Remark:   dto.Remark,
	}
	err := s.dao.UpdateCmdbHost(id, &host)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_UPDATE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

// 根据ID获取主机
func (s *CmdbHostServiceImpl) GetCmdbHostById(c *gin.Context, id uint) {
	host, err := s.dao.GetCmdbHostById(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
	vo := model.CmdbHostVo{
		ID:          host.ID,
		HostName:    host.HostName,
		Name:        host.Name,
		GroupID:     host.GroupID,
		GroupName:   group.Name,
		PrivateIP:   host.PrivateIP,
		PublicIP:    host.PublicIP,
		SSHIP:       host.SSHIP,
		SSHName:     host.SSHName,
		SSHKeyID:    host.SSHKeyID,
		SSHPort:     host.SSHPort,
		Remark:      host.Remark,
		Vendor:      fmt.Sprintf("%d", host.Vendor),
		Region:      host.Region,
		InstanceID:  host.InstanceID,
		OS:          host.OS,
		Status:      host.Status,
		CPU:         host.CPU,
		Memory:      host.Memory,
		Disk:        host.Disk,
		BillingType: host.BillingType,
		CreateTime:  host.CreateTime,
		ExpireTime:  host.ExpireTime,
		UpdateTime:  host.UpdateTime,
	}
	result.Success(c, vo)
}

// 根据名称获取主机
func (s *CmdbHostServiceImpl) GetCmdbHostByName(c *gin.Context, name string) {
	host, err := s.dao.GetCmdbHostByName(name)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
	vo := model.CmdbHostVo{
		ID:          host.ID,
		HostName:    host.HostName,
		Name:        host.Name,
		GroupID:     host.GroupID,
		GroupName:   group.Name,
		PrivateIP:   host.PrivateIP,
		PublicIP:    host.PublicIP,
		SSHIP:       host.SSHIP,
		SSHName:     host.SSHName,
		SSHKeyID:    host.SSHKeyID,
		SSHPort:     host.SSHPort,
		Remark:      host.Remark,
		Vendor:      fmt.Sprintf("%d", host.Vendor),
		Region:      host.Region,
		InstanceID:  host.InstanceID,
		OS:          host.OS,
		Status:      host.Status,
		CPU:         host.CPU,
		Memory:      host.Memory,
		Disk:        host.Disk,
		BillingType: host.BillingType,
		CreateTime:  host.CreateTime,
		ExpireTime:  host.ExpireTime,
		UpdateTime:  host.UpdateTime,
	}
	result.Success(c, vo)
}

// 删除主机
func (s *CmdbHostServiceImpl) DeleteCmdbHost(c *gin.Context, id uint) {
	err := s.dao.DeleteCmdbHost(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_DELETE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

// 根据分组ID获取主机列表
func (s *CmdbHostServiceImpl) GetCmdbHostsByGroupId(c *gin.Context, groupId uint) {
	list := s.dao.GetCmdbHostsByGroupId(groupId)
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 根据主机名称模糊查询
func (s *CmdbHostServiceImpl) GetCmdbHostsByHostNameLike(c *gin.Context, name string) {
	list := s.dao.GetCmdbHostsByHostNameLike(name)
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 根据IP查询(内网/公网/SSH)
func (s *CmdbHostServiceImpl) GetCmdbHostsByIP(c *gin.Context, ip string) {
	list := s.dao.GetCmdbHostsByIP(ip)
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

// 根据状态查询
func (s *CmdbHostServiceImpl) GetCmdbHostsByStatus(c *gin.Context, status int) {
	list := s.dao.GetCmdbHostsByStatus(status)
	var vos []model.CmdbHostVo
	for _, host := range list {
		group, _ := s.groupDao.GetCmdbGroupById(host.GroupID)
		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Vendor:      fmt.Sprintf("%d", host.Vendor),
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}
	result.Success(c, vos)
}

func GetCmdbHostService() CmdbHostServiceInterface {
	return &CmdbHostServiceImpl{
		dao:      cmdbDao.NewCmdbHostDao(),
		groupDao: cmdbDao.NewCmdbGroupDao(),
	}
}

// 同步主机基本信息
func (s *CmdbHostServiceImpl) SyncHostInfo(c *gin.Context, id uint) {
	// 1. 验证主机是否存在
	host, err := s.dao.GetCmdbHostById(id)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "主机不存在")
		return
	}

	// 2. 检查SSH连接信息是否完整
	if host.SSHIP == "" || host.SSHName == "" || host.SSHKeyID == 0 {
		result.FailedWithCode(c, constant.INVALID_PARAMS, "主机SSH连接信息不完整，无法同步")
		return
	}

	// 3. 获取SSH认证信息
	authDao := configDao.NewEcsAuthDao()
	auth, err := authDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		result.FailedWithCode(c, constant.CMDB_HOST_NOT_FOUND, "SSH认证凭据不存在")
		return
	}

	// 4. 立即返回成功响应，后台异步执行同步操作
	result.Success(c, gin.H{
		"id":      host.ID,
		"status":  host.Status,
		"message": "开始同步主机信息，请稍后查看结果",
	})

	// 5. 异步执行同步操作
	go func() {
		fmt.Printf("开始同步主机信息: ID=%d, Name=%s\n", host.ID, host.HostName)

		// 设置状态为同步中（使用状态2表示同步中）
		s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 2})

		// 准备SSH配置
		sshConfig := util.SSHConfig{
			IP:        host.SSHIP,
			Port:      host.SSHPort,
			Type:      auth.Type,
			Username:  host.SSHName,
			Password:  auth.Password,
			PublicKey: auth.PublicKey,
		}

		// 获取系统信息
		fmt.Printf("开始SSH连接获取系统信息: %s@%s:%d\n", host.SSHName, host.SSHIP, host.SSHPort)
		sshUtil := util.NewSSHUtil()
		systemInfo, err := sshUtil.GetSystemInfo(&sshConfig)
		if err != nil {
			fmt.Printf("SSH获取系统信息失败: %v\n", err)
			// 更新状态为同步失败（使用状态3表示同步失败）
			s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
			return
		}

		fmt.Printf("成功获取系统信息: %+v\n", systemInfo)

		// 更新主机信息
		updateData := model.CmdbHost{
			Name:       systemInfo["name"],      // 主机名称
			OS:         systemInfo["os"],        // 操作系统
			CPU:        systemInfo["cpu"],       // CPU
			Memory:     systemInfo["memory"],    // 内存
			Disk:       systemInfo["disk"],      // 磁盘
			PrivateIP:  systemInfo["privateIp"], // 内网IP
			PublicIP:   systemInfo["publicIp"],  // 公网IP
			Status:     1,                       // 同步成功
			UpdateTime: util.HTime{Time: time.Now()},
		}

		if err := s.dao.UpdateCmdbHost(host.ID, &updateData); err != nil {
			fmt.Printf("更新主机信息失败: %v\n", err)
			// 更新状态为同步失败
			s.dao.UpdateCmdbHost(host.ID, &model.CmdbHost{Status: 3})
			return
		}

		fmt.Printf("主机信息同步完成: ID=%d\n", host.ID)
	}()
}
