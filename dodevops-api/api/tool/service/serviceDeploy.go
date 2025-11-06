// 服务部署 Service层
// author xiaoRui
package service

import (
	"bytes"
	"dodevops-api/api/cmdb/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	ccDao "dodevops-api/api/configcenter/dao"
	"dodevops-api/api/tool/dao1"
	toolModel "dodevops-api/api/tool/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// IServiceDeployService 服务部署接口
type IServiceDeployService interface {
	GetServicesList(c *gin.Context)                                // 获取可部署服务列表
	GetServiceDetail(c *gin.Context, serviceID string)             // 获取服务详情
	CreateDeploy(c *gin.Context, dto toolModel.CreateDeployDto)    // 创建部署任务
	GetDeployList(c *gin.Context, dto toolModel.DeployQueryDto)    // 获取部署历史列表
	GetDeployStatus(c *gin.Context, id uint)                       // 获取部署状态
	DeleteDeploy(c *gin.Context, id uint)                          // 删除部署记录（卸载服务）
}

type ServiceDeployServiceImpl struct{}

// GetServicesList 获取可部署服务列表
func (s ServiceDeployServiceImpl) GetServicesList(c *gin.Context) {
	// 读取 services.json 文件
	servicesFile := "common/templates/services.json"
	data, err := os.ReadFile(servicesFile)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "读取服务配置失败: "+err.Error())
		return
	}

	var servicesConfig toolModel.ServicesConfig
	if err := json.Unmarshal(data, &servicesConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "解析服务配置失败: "+err.Error())
		return
	}

	result.Success(c, servicesConfig)
}

// GetServiceDetail 获取服务详情
func (s ServiceDeployServiceImpl) GetServiceDetail(c *gin.Context, serviceID string) {
	if serviceID == "" {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "服务ID不能为空")
		return
	}

	// 读取 services.json
	servicesFile := "common/templates/services.json"
	data, err := os.ReadFile(servicesFile)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "读取服务配置失败: "+err.Error())
		return
	}

	var servicesConfig toolModel.ServicesConfig
	if err := json.Unmarshal(data, &servicesConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "解析服务配置失败: "+err.Error())
		return
	}

	// 查找指定服务
	for _, service := range servicesConfig.Services {
		if service.ID == serviceID {
			result.Success(c, service)
			return
		}
	}

	result.Failed(c, int(result.ApiCode.FAILED), "服务不存在")
}

// CreateDeploy 创建部署任务
func (s ServiceDeployServiceImpl) CreateDeploy(c *gin.Context, dto toolModel.CreateDeployDto) {
	// 参数校验
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	// 获取主机信息
	cmdbHostDao := cmdbDao.NewCmdbHostDao()
	host, err := cmdbHostDao.GetCmdbHostById(dto.HostID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "主机不存在: "+err.Error())
		return
	}

	// 检查主机连接状态
	if host.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "主机未认证或认证失败，无法部署")
		return
	}

	// 读取服务配置
	servicesFile := "common/templates/services.json"
	data, err := os.ReadFile(servicesFile)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "读取服务配置失败: "+err.Error())
		return
	}

	var servicesConfig toolModel.ServicesConfig
	if err := json.Unmarshal(data, &servicesConfig); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "解析服务配置失败: "+err.Error())
		return
	}

	// 查找服务和版本
	var serviceInfo *toolModel.ServiceInfo
	var versionInfo *toolModel.ServiceVersion
	for _, service := range servicesConfig.Services {
		if service.ID == dto.ServiceID {
			serviceInfo = &service
			for _, version := range service.Versions {
				if version.ID == dto.Version {
					versionInfo = &version
					break
				}
			}
			break
		}
	}

	if serviceInfo == nil {
		result.Failed(c, int(result.ApiCode.FAILED), "服务不存在")
		return
	}
	if versionInfo == nil {
		result.Failed(c, int(result.ApiCode.FAILED), "服务版本不存在")
		return
	}

	// 检查是否已部署
	exists, err := dao.CheckServiceDeployExists(dto.HostID, dto.ServiceID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "检查部署状态失败: "+err.Error())
		return
	}
	if exists {
		result.Failed(c, int(result.ApiCode.FAILED), "该服务已在此主机上部署")
		return
	}

	// 创建部署记录
	envVarsJSON, _ := json.Marshal(dto.EnvVars)
	deployRecord := &toolModel.ServiceDeploy{
		ServiceName: serviceInfo.Name,
		ServiceID:   dto.ServiceID,
		Version:     dto.Version,
		HostID:      dto.HostID,
		HostIP:      host.SSHIP,
		InstallDir:  dto.InstallDir,
		EnvVars:     string(envVarsJSON),
		Status:      0, // 部署中
		CreateTime:  util.HTime{Time: time.Now()},
		UpdateTime:  time.Now(),
	}

	if err := dao.CreateServiceDeploy(deployRecord); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建部署记录失败: "+err.Error())
		return
	}

	// 异步执行部署任务
	go s.executeDeploy(deployRecord.ID, &host, serviceInfo, versionInfo, dto)

	result.Success(c, map[string]interface{}{
		"deployId": deployRecord.ID,
		"message":  "部署任务已创建，正在后台执行",
	})
}

// executeDeploy 执行部署任务（异步）
func (s ServiceDeployServiceImpl) executeDeploy(
	deployID uint,
	host *model.CmdbHost,
	serviceInfo *toolModel.ServiceInfo,
	versionInfo *toolModel.ServiceVersion,
	dto toolModel.CreateDeployDto,
) {
	var deployLog bytes.Buffer
	status := 1 // 默认成功

	defer func() {
		// 更新部署状态
		dao.UpdateServiceDeployStatus(deployID, status, deployLog.String())
	}()

	deployLog.WriteString(fmt.Sprintf("[%s] 开始部署 %s %s\n", time.Now().Format("2006-01-02 15:04:05"), serviceInfo.Name, versionInfo.Name))

	// 1. 获取SSH凭据
	ecsAuthDao := ccDao.NewEcsAuthDao()
	key, err := ecsAuthDao.GetById(host.SSHKeyID)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 获取SSH凭据失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		status = 3
		return
	}

	// 2. 建立SSH连接
	deployLog.WriteString(fmt.Sprintf("[%s] 连接主机 %s...\n", time.Now().Format("2006-01-02 15:04:05"), host.SSHIP))
	sshConfig := &util.SSHConfig{
		IP:        host.SSHIP,
		Port:      host.SSHPort,
		Username:  host.SSHName,
		Type:      key.Type,
		Password:  key.Password,
		PublicKey: key.PublicKey,
	}
	sshClient, err := util.GetSSHClientByConfig(sshConfig)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] SSH连接失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		status = 3
		return
	}
	defer sshClient.Close()
	deployLog.WriteString(fmt.Sprintf("[%s] SSH连接成功\n", time.Now().Format("2006-01-02 15:04:05")))

	// 根据部署类型选择不同的部署方式
	deployType := versionInfo.DeployType
	if deployType == "" {
		deployType = "container" // 默认容器部署
	}

	if deployType == "binary" {
		// 二进制部署方式
		status = s.executeBinaryDeploy(sshClient, &deployLog, serviceInfo, versionInfo, dto)
	} else {
		// 容器部署方式（原有逻辑）
		status = s.executeContainerDeploy(sshClient, &deployLog, serviceInfo, versionInfo, dto)
	}
}

// executeBinaryDeploy 执行二进制部署（从镜像提取）
func (s ServiceDeployServiceImpl) executeBinaryDeploy(
	sshClient *ssh.Client,
	deployLog *bytes.Buffer,
	serviceInfo *toolModel.ServiceInfo,
	versionInfo *toolModel.ServiceVersion,
	dto toolModel.CreateDeployDto,
) int {
	// 1. 检查Docker环境
	deployLog.WriteString(fmt.Sprintf("[%s] 检查Docker环境...\n", time.Now().Format("2006-01-02 15:04:05")))
	if err := util.ExecuteSSHCommand(sshClient, "docker --version"); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] Docker未安装，请先安装Docker\n", time.Now().Format("2006-01-02 15:04:05")))
		return 3
	}

	// 2. 获取镜像地址
	imageRegistry := config.GetImageRegistry()
	versionNumber := strings.TrimPrefix(versionInfo.ID, serviceInfo.ID+"-")
	dockerImage := imageRegistry.GetImage(serviceInfo.ID, versionNumber, true)
	deployLog.WriteString(fmt.Sprintf("[%s] 使用镜像: %s\n", time.Now().Format("2006-01-02 15:04:05"), dockerImage))

	// 3. 拉取镜像
	deployLog.WriteString(fmt.Sprintf("[%s] 拉取镜像...\n", time.Now().Format("2006-01-02 15:04:05")))
	pullCmd := fmt.Sprintf("docker pull %s", dockerImage)
	output, err := util.ExecuteSSHCommandWithOutput(sshClient, pullCmd)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 拉取镜像失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		deployLog.WriteString(fmt.Sprintf("输出: %s\n", output))
		return 3
	}
	deployLog.WriteString(fmt.Sprintf("[%s] 镜像拉取成功\n", time.Now().Format("2006-01-02 15:04:05")))

	// 4. 从镜像提取二进制文件
	if len(versionInfo.ExtractPaths) == 0 {
		deployLog.WriteString(fmt.Sprintf("[%s] 错误: 未配置extract_paths\n", time.Now().Format("2006-01-02 15:04:05")))
		return 3
	}

	containerName := fmt.Sprintf("temp_%s_%d", serviceInfo.ID, time.Now().Unix())

	// 5. 创建临时容器
	deployLog.WriteString(fmt.Sprintf("[%s] 创建临时容器...\n", time.Now().Format("2006-01-02 15:04:05")))
	createCmd := fmt.Sprintf("docker create --name %s %s", containerName, dockerImage)
	if err := util.ExecuteSSHCommand(sshClient, createCmd); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 创建临时容器失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 6. 从容器复制文件到宿主机
	for srcPath, destPath := range versionInfo.ExtractPaths {
		deployLog.WriteString(fmt.Sprintf("[%s] 提取文件 %s -> %s...\n", time.Now().Format("2006-01-02 15:04:05"), srcPath, destPath))

		// 删除目标目录（如果存在）
		util.ExecuteSSHCommand(sshClient, fmt.Sprintf("rm -rf %s", destPath))

		// 复制文件
		copyCmd := fmt.Sprintf("docker cp %s:%s %s", containerName, srcPath, destPath)
		output, err := util.ExecuteSSHCommandWithOutput(sshClient, copyCmd)
		if err != nil {
			deployLog.WriteString(fmt.Sprintf("[%s] 提取文件失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			deployLog.WriteString(fmt.Sprintf("输出: %s\n", output))
			// 清理临时容器
			util.ExecuteSSHCommand(sshClient, fmt.Sprintf("docker rm %s", containerName))
			return 3
		}
	}

	// 7. 删除临时容器
	deployLog.WriteString(fmt.Sprintf("[%s] 清理临时容器...\n", time.Now().Format("2006-01-02 15:04:05")))
	util.ExecuteSSHCommand(sshClient, fmt.Sprintf("docker rm %s", containerName))

	// 8. 读取并执行安装脚本
	templatePath := filepath.Join("common/templates", versionInfo.File)
	deployLog.WriteString(fmt.Sprintf("[%s] 读取安装脚本 %s...\n", time.Now().Format("2006-01-02 15:04:05"), templatePath))
	scriptData, err := os.ReadFile(templatePath)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 读取安装脚本失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 9. 创建SFTP客户端上传脚本
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] SFTP连接失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}
	defer sftpClient.Close()

	// 10. 上传安装脚本
	scriptPath := fmt.Sprintf("/tmp/install_%s.sh", serviceInfo.ID)
	deployLog.WriteString(fmt.Sprintf("[%s] 上传安装脚本...\n", time.Now().Format("2006-01-02 15:04:05")))
	if err := s.uploadFile(sftpClient, scriptPath, scriptData); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 上传脚本失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 11. 执行安装脚本
	deployLog.WriteString(fmt.Sprintf("[%s] 执行安装脚本...\n", time.Now().Format("2006-01-02 15:04:05")))

	// 获取安装路径和版本号（用于环境变量文件命名）
	var installPath string
	var profileName string
	for _, destPath := range versionInfo.ExtractPaths {
		installPath = destPath
		break
	}
	// 从 version.ID 提取版本号，如 java-11 -> 11
	versionNum := strings.TrimPrefix(versionInfo.ID, serviceInfo.ID+"-")
	profileName = fmt.Sprintf("%s%s", serviceInfo.ID, versionNum)

	execCmd := fmt.Sprintf("chmod +x %s && bash %s %s %s", scriptPath, scriptPath, installPath, profileName)
	output, err = util.ExecuteSSHCommandWithOutput(sshClient, execCmd)
	deployLog.WriteString(fmt.Sprintf("[%s] 安装脚本输出:\n%s\n", time.Now().Format("2006-01-02 15:04:05"), output))
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 安装失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 12. 清理安装脚本
	util.ExecuteSSHCommand(sshClient, fmt.Sprintf("rm -f %s", scriptPath))

	// 13. 验证安装
	deployLog.WriteString(fmt.Sprintf("[%s] 验证安装...\n", time.Now().Format("2006-01-02 15:04:05")))
	var verifyCmd string

	// 先加载环境变量，再执行验证命令
	profilePath := fmt.Sprintf("/etc/profile.d/%s.sh", profileName)
	sourceCmd := fmt.Sprintf(". %s", profilePath)

	switch serviceInfo.ID {
	case "golang":
		verifyCmd = fmt.Sprintf("%s && go version", sourceCmd)
	case "java":
		verifyCmd = fmt.Sprintf("%s && java -version 2>&1", sourceCmd)
	case "nodejs":
		verifyCmd = fmt.Sprintf("%s && node -v && npm -v", sourceCmd)
	default:
		verifyCmd = fmt.Sprintf("%s && %s --version", sourceCmd, serviceInfo.ID)
	}

	verifyOutput, err := util.ExecuteSSHCommandWithOutput(sshClient, verifyCmd)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 验证失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}
	deployLog.WriteString(fmt.Sprintf("[%s] 验证结果:\n%s\n", time.Now().Format("2006-01-02 15:04:05"), verifyOutput))

	deployLog.WriteString(fmt.Sprintf("[%s] 部署完成！\n", time.Now().Format("2006-01-02 15:04:05")))
	return 1
}

// executeContainerDeploy 执行容器部署（原有逻辑）
func (s ServiceDeployServiceImpl) executeContainerDeploy(
	sshClient *ssh.Client,
	deployLog *bytes.Buffer,
	serviceInfo *toolModel.ServiceInfo,
	versionInfo *toolModel.ServiceVersion,
	dto toolModel.CreateDeployDto,
) int {
	// 1. 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] SFTP连接失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}
	defer sftpClient.Close()

	// 2. 创建安装目录
	deployLog.WriteString(fmt.Sprintf("[%s] 创建安装目录 %s...\n", time.Now().Format("2006-01-02 15:04:05"), dto.InstallDir))
	if err := util.ExecuteSSHCommand(sshClient, fmt.Sprintf("mkdir -p %s", dto.InstallDir)); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 创建目录失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 3. 读取docker-compose模板
	templatePath := filepath.Join("common/templates", versionInfo.File)
	deployLog.WriteString(fmt.Sprintf("[%s] 读取模板文件 %s...\n", time.Now().Format("2006-01-02 15:04:05"), templatePath))
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 读取模板失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 4. 生成.env文件内容（包含镜像地址）
	deployLog.WriteString(fmt.Sprintf("[%s] 生成环境变量配置...\n", time.Now().Format("2006-01-02 15:04:05")))
	envContent := s.generateEnvFile(serviceInfo, versionInfo, dto)

	// 5. 上传docker-compose.yml
	composePath := filepath.Join(dto.InstallDir, "docker-compose.yml")
	deployLog.WriteString(fmt.Sprintf("[%s] 上传 docker-compose.yml...\n", time.Now().Format("2006-01-02 15:04:05")))
	if err := s.uploadFile(sftpClient, composePath, templateData); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 上传失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 6. 上传.env文件
	envPath := filepath.Join(dto.InstallDir, ".env")
	deployLog.WriteString(fmt.Sprintf("[%s] 上传 .env...\n", time.Now().Format("2006-01-02 15:04:05")))
	if err := s.uploadFile(sftpClient, envPath, []byte(envContent)); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 上传失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}

	// 7. 检查Docker是否安装
	deployLog.WriteString(fmt.Sprintf("[%s] 检查Docker环境...\n", time.Now().Format("2006-01-02 15:04:05")))
	if err := util.ExecuteSSHCommand(sshClient, "docker --version"); err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] Docker未安装，请先安装Docker\n", time.Now().Format("2006-01-02 15:04:05")))
		return 3
	}

	// 8. 执行docker-compose up
	deployLog.WriteString(fmt.Sprintf("[%s] 启动服务容器...\n", time.Now().Format("2006-01-02 15:04:05")))
	composeCmd := fmt.Sprintf("cd %s && docker-compose up -d", dto.InstallDir)
	output, err := util.ExecuteSSHCommandWithOutput(sshClient, composeCmd)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 启动失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		deployLog.WriteString(fmt.Sprintf("输出: %s\n", output))
		return 3
	}
	deployLog.WriteString(fmt.Sprintf("[%s] 容器启动输出:\n%s\n", time.Now().Format("2006-01-02 15:04:05"), output))

	// 9. 验证容器状态
	deployLog.WriteString(fmt.Sprintf("[%s] 验证容器状态...\n", time.Now().Format("2006-01-02 15:04:05")))
	time.Sleep(3 * time.Second)
	checkCmd := fmt.Sprintf("cd %s && docker-compose ps", dto.InstallDir)
	psOutput, err := util.ExecuteSSHCommandWithOutput(sshClient, checkCmd)
	if err != nil {
		deployLog.WriteString(fmt.Sprintf("[%s] 获取容器状态失败: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return 3
	}
	deployLog.WriteString(fmt.Sprintf("[%s] 容器状态:\n%s\n", time.Now().Format("2006-01-02 15:04:05"), psOutput))

	deployLog.WriteString(fmt.Sprintf("[%s] 部署完成！\n", time.Now().Format("2006-01-02 15:04:05")))
	return 1
}

// generateEnvFile 生成.env文件内容
func (s ServiceDeployServiceImpl) generateEnvFile(serviceInfo *toolModel.ServiceInfo, versionInfo *toolModel.ServiceVersion, dto toolModel.CreateDeployDto) string {
	var envLines []string

	// 添加基础配置
	envLines = append(envLines, fmt.Sprintf("# %s 环境配置", serviceInfo.Name))
	envLines = append(envLines, fmt.Sprintf("# 生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	envLines = append(envLines, "")

	// 添加镜像地址（从配置文件获取）
	imageRegistry := config.GetImageRegistry()
	// 从 version.ID 提取真实版本号（如 mysql-5.7 -> 5.7）
	versionNumber := strings.TrimPrefix(versionInfo.ID, serviceInfo.ID+"-")
	dockerImage := imageRegistry.GetImage(serviceInfo.ID, versionNumber, true)
	envLines = append(envLines, fmt.Sprintf("# Docker镜像"))
	envLines = append(envLines, fmt.Sprintf("DOCKER_IMAGE=%s", dockerImage))
	envLines = append(envLines, "")

	// 添加安装目录
	envLines = append(envLines, fmt.Sprintf("DATA_DIR=%s", dto.InstallDir))

	// 添加用户自定义环境变量
	for key, value := range dto.EnvVars {
		envLines = append(envLines, fmt.Sprintf("%s=%v", key, value))
	}

	// 添加默认环境变量（如果用户未提供）
	for _, envVar := range serviceInfo.EnvVars {
		if _, exists := dto.EnvVars[envVar.Name]; !exists && envVar.Default != "" {
			envLines = append(envLines, fmt.Sprintf("%s=%s", envVar.Name, envVar.Default))
		}
	}

	return strings.Join(envLines, "\n")
}

// uploadFile 上传文件到远程主机
func (s ServiceDeployServiceImpl) uploadFile(sftpClient *sftp.Client, remotePath string, data []byte) error {
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	_, err = remoteFile.Write(data)
	return err
}

// GetDeployList 获取部署历史列表
func (s ServiceDeployServiceImpl) GetDeployList(c *gin.Context, dto toolModel.DeployQueryDto) {
	// 设置默认分页参数
	if dto.PageNum <= 0 {
		dto.PageNum = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	deploys, total, err := dao.GetServiceDeployList(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	// 转换为VO
	cmdbHostDao := cmdbDao.NewCmdbHostDao()
	var deployVos []toolModel.ServiceDeployVo
	for _, deploy := range deploys {
		// 获取主机信息
		host, _ := cmdbHostDao.GetCmdbHostById(deploy.HostID)
		hostName := ""
		if err == nil {
			hostName = host.HostName
		}

		// 解析环境变量
		var envVars map[string]interface{}
		if deploy.EnvVars != "" {
			json.Unmarshal([]byte(deploy.EnvVars), &envVars)
		}

		// 状态文本
		statusTexts := map[int]string{
			0: "部署中",
			1: "运行中",
			2: "已停止",
			3: "部署失败",
		}

		deployVos = append(deployVos, toolModel.ServiceDeployVo{
			ID:            deploy.ID,
			ServiceName:   deploy.ServiceName,
			ServiceID:     deploy.ServiceID,
			Version:       deploy.Version,
			HostID:        deploy.HostID,
			HostIP:        deploy.HostIP,
			HostName:      hostName,
			InstallDir:    deploy.InstallDir,
			ContainerName: deploy.ContainerName,
			Ports:         deploy.Ports,
			EnvVars:       envVars,
			Status:        deploy.Status,
			StatusText:    statusTexts[deploy.Status],
			DeployLog:     deploy.DeployLog,
			CreateTime:    deploy.CreateTime,
			UpdateTime:    deploy.UpdateTime,
		})
	}

	data := map[string]interface{}{
		"list":     deployVos,
		"total":    total,
		"pageNum":  dto.PageNum,
		"pageSize": dto.PageSize,
	}

	result.Success(c, data)
}

// GetDeployStatus 获取部署状态
func (s ServiceDeployServiceImpl) GetDeployStatus(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	deploy, err := dao.GetServiceDeployByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	// 获取主机信息
	cmdbHostDao := cmdbDao.NewCmdbHostDao()
	host, _ := cmdbHostDao.GetCmdbHostById(deploy.HostID)
	hostName := ""
	if err == nil {
		hostName = host.HostName
	}

	// 解析环境变量
	var envVars map[string]interface{}
	if deploy.EnvVars != "" {
		json.Unmarshal([]byte(deploy.EnvVars), &envVars)
	}

	// 状态文本
	statusTexts := map[int]string{
		0: "部署中",
		1: "运行中",
		2: "已停止",
		3: "部署失败",
	}

	deployVo := toolModel.ServiceDeployVo{
		ID:            deploy.ID,
		ServiceName:   deploy.ServiceName,
		ServiceID:     deploy.ServiceID,
		Version:       deploy.Version,
		HostID:        deploy.HostID,
		HostIP:        deploy.HostIP,
		HostName:      hostName,
		InstallDir:    deploy.InstallDir,
		ContainerName: deploy.ContainerName,
		Ports:         deploy.Ports,
		EnvVars:       envVars,
		Status:        deploy.Status,
		StatusText:    statusTexts[deploy.Status],
		DeployLog:     deploy.DeployLog,
		CreateTime:    deploy.CreateTime,
		UpdateTime:    deploy.UpdateTime,
	}

	result.Success(c, deployVo)
}

// DeleteDeploy 删除部署记录（卸载服务）
func (s ServiceDeployServiceImpl) DeleteDeploy(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	// 获取部署记录
	deploy, err := dao.GetServiceDeployByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "部署记录不存在")
		return
	}

	// 获取主机信息
	cmdbHostDao := cmdbDao.NewCmdbHostDao()
	host, err := cmdbHostDao.GetCmdbHostById(deploy.HostID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "主机不存在")
		return
	}

	// 尝试停止并删除容器
	if deploy.Status == 1 { // 如果正在运行
		// 获取SSH凭据
		ecsAuthDao := ccDao.NewEcsAuthDao()
		key, err := ecsAuthDao.GetById(host.SSHKeyID)
		if err == nil {
			sshConfig := &util.SSHConfig{
				IP:        host.SSHIP,
				Port:      host.SSHPort,
				Username:  host.SSHName,
				Type:      key.Type,
				Password:  key.Password,
				PublicKey: key.PublicKey,
			}
			sshClient, err := util.GetSSHClientByConfig(sshConfig)
			if err == nil {
				defer sshClient.Close()
				// 执行 docker-compose down
				downCmd := fmt.Sprintf("cd %s && docker-compose down -v", deploy.InstallDir)
				util.ExecuteSSHCommand(sshClient, downCmd)
			}
		}
	}

	// 删除部署记录
	if err := dao.DeleteServiceDeploy(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

var serviceDeployService = ServiceDeployServiceImpl{}

// ServiceDeployService 获取服务实例
func ServiceDeployService() IServiceDeployService {
	return &serviceDeployService
}
