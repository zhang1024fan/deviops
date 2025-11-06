package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/app/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IJenkinsService Jenkins服务接口
type IJenkinsService interface {
	// 服务器管理
	GetJenkinsServers(c *gin.Context)
	GetJenkinsServerDetail(c *gin.Context, serverID uint)
	TestJenkinsConnection(c *gin.Context, req *model.TestJenkinsConnectionRequest)

	// 任务管理
	GetJobs(c *gin.Context, serverID uint)
	SearchJobs(c *gin.Context, serverID uint, keyword string)
	GetJobDetail(c *gin.Context, serverID uint, jobName string)
	StartJob(c *gin.Context, serverID uint, jobName string, req *model.StartJobRequest)
	StopBuild(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.StopBuildRequest)

	// 构建管理
	GetBuildDetail(c *gin.Context, serverID uint, jobName string, buildNumber int)
	GetBuildLog(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.GetBuildLogRequest)

	// 系统信息
	GetSystemInfo(c *gin.Context, serverID uint)
	GetQueueInfo(c *gin.Context, serverID uint)
}

// JenkinsServiceImpl Jenkins服务实现
type JenkinsServiceImpl struct {
	db         *gorm.DB
	accountDao *dao.AccountAuthDao
}

// NewJenkinsService 创建Jenkins服务实例
func NewJenkinsService(db *gorm.DB) IJenkinsService {
	return &JenkinsServiceImpl{
		db:         db,
		accountDao: dao.NewAccountAuthDao(),
	}
}

// JenkinsClient Jenkins客户端
type JenkinsClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewJenkinsClient 创建Jenkins客户端
func NewJenkinsClient(host string, port int, username, password string) *JenkinsClient {
	baseURL := fmt.Sprintf("http://%s:%d", host, port)
	if port == 443 {
		baseURL = fmt.Sprintf("https://%s", host)
	}

	return &JenkinsClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 发送HTTP请求
func (jc *JenkinsClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := jc.BaseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jc.Username, jc.Password)
	req.Header.Set("Content-Type", "application/json")

	return jc.HTTPClient.Do(req)
}

// getJenkinsClient 获取Jenkins客户端
func (s *JenkinsServiceImpl) getJenkinsClient(serverID uint) (*JenkinsClient, error) {
	account, err := s.accountDao.GetByID(serverID)
	if err != nil {
		return nil, err
	}

	if account.Type != model.JenkinsAccountType {
		return nil, fmt.Errorf("账号类型不是Jenkins")
	}

	password, err := account.DecryptPassword()
	if err != nil {
		return nil, err
	}

	return NewJenkinsClient(account.Host, account.Port, account.Name, password), nil
}

// GetJenkinsServers 获取Jenkins服务器列表
func (s *JenkinsServiceImpl) GetJenkinsServers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	accounts, total, err := s.accountDao.GetAccountsByType(model.JenkinsAccountType, page, pageSize)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器列表失败: "+err.Error())
		return
	}

	var serverList []model.JenkinsServerInfo
	for _, account := range accounts {
		serverInfo := model.JenkinsServerInfo{
			ID:          account.ID,
			Alias:       account.Alias,
			Host:        account.Host,
			Port:        account.Port,
			Username:    account.Name,
			Description: account.Remark,
			CreatedAt:   account.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   account.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		serverList = append(serverList, serverInfo)
	}

	response := model.JenkinsServerListResponse{
		List:  serverList,
		Total: total,
	}

	result.Success(c, response)
}

// GetJenkinsServerDetail 获取Jenkins服务器详情
func (s *JenkinsServiceImpl) GetJenkinsServerDetail(c *gin.Context, serverID uint) {
	account, err := s.accountDao.GetByID(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器详情失败: "+err.Error())
		return
	}

	if account.Type != model.JenkinsAccountType {
		result.Failed(c, 400, "账号类型不是Jenkins")
		return
	}

	serverInfo := model.JenkinsServerInfo{
		ID:          account.ID,
		Alias:       account.Alias,
		Host:        account.Host,
		Port:        account.Port,
		Username:    account.Name,
		Description: account.Remark,
		CreatedAt:   account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	result.Success(c, serverInfo)
}

// TestJenkinsConnection 测试Jenkins连接
func (s *JenkinsServiceImpl) TestJenkinsConnection(c *gin.Context, req *model.TestJenkinsConnectionRequest) {
	// 解析URL获取host和port
	u, err := url.Parse(req.URL)
	if err != nil {
		result.Failed(c, 500, "无效的Jenkins URL: "+err.Error())
		return
	}

	host := u.Hostname()
	port := 80
	if u.Port() != "" {
		port, _ = strconv.Atoi(u.Port())
	} else if u.Scheme == "https" {
		port = 443
	}

	client := NewJenkinsClient(host, port, req.Username, req.Password)

	// 获取系统信息来测试连接
	resp, err := client.makeRequest("GET", "/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "连接Jenkins失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var systemInfo model.JenkinsSystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&systemInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	response := model.TestJenkinsConnectionResponse{
		Success:    true,
		Message:    "连接成功",
		SystemInfo: &systemInfo,
	}

	result.Success(c, response)
}

// GetJobs 获取任务列表
func (s *JenkinsServiceImpl) GetJobs(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/api/json?tree=jobs[name,displayName,description,url,buildable,color,_class,lastBuild[number,url,displayName,result,building,duration,timestamp],lastStableBuild[number,url,displayName,result],lastSuccessfulBuild[number,url,displayName,result],lastFailedBuild[number,url,displayName,result]]", nil)
	if err != nil {
		result.Failed(c, 500, "获取任务列表失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var jenkinsData struct {
		Jobs []model.JenkinsJob `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jenkinsData); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobListResponse{
		Jobs:   jenkinsData.Jobs,
		Total:  len(jenkinsData.Jobs),
		Server: serverName,
	}

	result.Success(c, response)
}

// SearchJobs 搜索Jenkins任务
func (s *JenkinsServiceImpl) SearchJobs(c *gin.Context, serverID uint, keyword string) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	// 获取所有任务
	resp, err := client.makeRequest("GET", "/api/json?tree=jobs[name,displayName,description,url,buildable,color,_class,lastBuild[number,url,displayName,result,building,duration,timestamp],lastStableBuild[number,url,displayName,result],lastSuccessfulBuild[number,url,displayName,result],lastFailedBuild[number,url,displayName,result]]", nil)
	if err != nil {
		result.Failed(c, 500, "获取任务列表失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var jenkinsData struct {
		Jobs []model.JenkinsJob `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jenkinsData); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 过滤任务 - 支持名称、显示名称、描述的模糊匹配
	var filteredJobs []model.JenkinsJob
	keyword = strings.ToLower(keyword)

	for _, job := range jenkinsData.Jobs {
		// 检查任务名称
		if strings.Contains(strings.ToLower(job.Name), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
		// 检查显示名称
		if strings.Contains(strings.ToLower(job.DisplayName), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
		// 检查描述
		if strings.Contains(strings.ToLower(job.Description), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobListResponse{
		Jobs:   filteredJobs,
		Total:  len(filteredJobs),
		Server: serverName,
	}

	result.Success(c, response)
}

// GetJobDetail 获取任务详情
func (s *JenkinsServiceImpl) GetJobDetail(c *gin.Context, serverID uint, jobName string) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	// 获取任务详情
	jobURL := fmt.Sprintf("/job/%s/api/json", url.PathEscape(jobName))
	resp, err := client.makeRequest("GET", jobURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取任务详情失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var job model.JenkinsJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取构建历史
	buildsURL := fmt.Sprintf("/job/%s/api/json?tree=builds[number,url,displayName,result,building,duration,timestamp,keepLog,queueId]", url.PathEscape(jobName))
	resp, err = client.makeRequest("GET", buildsURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建历史失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var buildsData struct {
		Builds []model.JenkinsBuild `json:"builds"`
	}

	if resp.StatusCode == http.StatusOK {
		json.NewDecoder(resp.Body).Decode(&buildsData)
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobDetailResponse{
		Job:    job,
		Builds: buildsData.Builds,
		Server: serverName,
	}

	result.Success(c, response)
}

// StartJob 启动任务
func (s *JenkinsServiceImpl) StartJob(c *gin.Context, serverID uint, jobName string, req *model.StartJobRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	var buildURL string
	if req.Parameters != nil && len(req.Parameters) > 0 {
		// 带参数构建
		buildURL = fmt.Sprintf("/job/%s/buildWithParameters", url.PathEscape(jobName))

		// 构建表单数据
		formData := url.Values{}
		for key, value := range req.Parameters {
			formData.Set(key, value)
		}

		resp, err := client.makeRequest("POST", buildURL+"?"+formData.Encode(), nil)
		if err != nil {
			result.Failed(c, 500, "启动任务失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
			return
		}
	} else {
		// 无参数构建
		buildURL = fmt.Sprintf("/job/%s/build", url.PathEscape(jobName))

		resp, err := client.makeRequest("POST", buildURL, nil)
		if err != nil {
			result.Failed(c, 500, "启动任务失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
			return
		}
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.StartJobResponse{
		Success: true,
		Message: "任务启动成功",
		JobName: jobName,
		Server:  serverName,
	}

	result.Success(c, response)
}

// StopBuild 停止构建
func (s *JenkinsServiceImpl) StopBuild(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.StopBuildRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	stopURL := fmt.Sprintf("/job/%s/%d/stop", url.PathEscape(jobName), buildNumber)
	resp, err := client.makeRequest("POST", stopURL, nil)
	if err != nil {
		result.Failed(c, 500, "停止构建失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.StopBuildResponse{
		Success:     true,
		Message:     "构建停止成功",
		JobName:     jobName,
		BuildNumber: buildNumber,
		Server:      serverName,
	}

	result.Success(c, response)
}

// GetBuildDetail 获取构建详情
func (s *JenkinsServiceImpl) GetBuildDetail(c *gin.Context, serverID uint, jobName string, buildNumber int) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	buildURL := fmt.Sprintf("/job/%s/%d/api/json", url.PathEscape(jobName), buildNumber)
	resp, err := client.makeRequest("GET", buildURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建详情失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var build model.JenkinsBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsBuildDetailResponse{
		Build:  build,
		Server: serverName,
	}

	result.Success(c, response)
}

// GetBuildLog 获取构建日志
func (s *JenkinsServiceImpl) GetBuildLog(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.GetBuildLogRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	logURL := fmt.Sprintf("/job/%s/%d/logText/progressiveText?start=%d", url.PathEscape(jobName), buildNumber, req.Start)
	resp, err := client.makeRequest("GET", logURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建日志失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	logBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, 500, "读取日志内容失败: "+err.Error())
		return
	}

	// 检查是否有更多日志
	moreData := resp.Header.Get("X-More-Data") == "true"
	textSize, _ := strconv.Atoi(resp.Header.Get("X-Text-Size"))

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.GetBuildLogResponse{
		Log:         string(logBytes),
		HasMore:     moreData,
		TextSize:    textSize,
		MoreData:    moreData,
		JobName:     jobName,
		BuildNumber: buildNumber,
		Server:      serverName,
	}

	result.Success(c, response)
}

// GetSystemInfo 获取系统信息
func (s *JenkinsServiceImpl) GetSystemInfo(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "获取系统信息失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var systemInfo model.JenkinsSystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&systemInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	result.Success(c, systemInfo)
}

// GetQueueInfo 获取队列信息
func (s *JenkinsServiceImpl) GetQueueInfo(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/queue/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "获取队列信息失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var queueInfo model.JenkinsQueue
	if err := json.NewDecoder(resp.Body).Decode(&queueInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	result.Success(c, queueInfo)
}