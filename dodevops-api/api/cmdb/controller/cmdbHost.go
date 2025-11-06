package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type CmdbHostController struct {
	service service.CmdbHostServiceInterface
}

func NewCmdbHostController() *CmdbHostController {
	return &CmdbHostController{
		service: service.GetCmdbHostService(),
	}
}

// 分页参数
type PageParams struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

// 下载主机导入模板
// @Summary 下载主机导入模板
// @Description 下载主机导入Excel模板
// @Tags CMDB资产管理
// @Produce octet-stream
// @Success 200 {file} file
// @Router /api/v1/cmdb/hosttemplate [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) DownloadHostTemplate(ctx *gin.Context) {
	filePath := "upload/xlsl/host.xlsx"

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(404, gin.H{"error": "模板文件不存在"})
		return
	}

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename=host_template.xlsx")
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 返回文件
	ctx.File(filePath)
}

// 获取主机列表(分页)
// @Summary 获取主机列表(分页)
// @Description 获取主机列表(分页)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostlist [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostListWithPage(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "分页参数错误")
		return
	}
	c.service.GetCmdbHostListWithPage(ctx, params.Page, params.PageSize)
}

// 创建主机
// @Summary 创建主机
// @Description 创建主机
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CreateCmdbHostDto true "主机信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostcreate [post]
// @Security ApiKeyAuth
func (c *CmdbHostController) CreateCmdbHost(ctx *gin.Context) {
	var dto model.CreateCmdbHostDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.CreateCmdbHost(ctx, &dto)
}

// 从Excel导入主机
// @Summary 从Excel导入主机
// @Description 通过上传Excel模板批量导入主机（Excel列顺序：主机别名、SSH地址、SSH端口、SSH用户、备注）
// @Tags CMDB资产管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Param groupId formData int true "分组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostimport [post]
// @Security ApiKeyAuth
func (c *CmdbHostController) ImportHostsFromExcel(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "请上传文件")
		return
	}

	// 获取分组ID
	groupId := util.StringToUint(ctx.PostForm("groupId"))
	if groupId == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "分组ID不能为空")
		return
	}

	// 使用项目upload目录作为临时存储，避免系统临时目录的权限问题
	tempDir := "./upload/temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		result.Failed(ctx, constant.FILE_OPERATION_ERROR, "创建临时目录失败: "+err.Error())
		return
	}

	// 生成唯一的临时文件名
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(file.Filename)
	tempFileName := fmt.Sprintf("host_import_%d%s", timestamp, ext)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	// 保存上传的文件
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		result.Failed(ctx, constant.FILE_OPERATION_ERROR, "保存上传文件失败: "+err.Error())
		return
	}
	defer os.Remove(tempFilePath)

	// 打开Excel文件
	f, err := excelize.OpenFile(tempFilePath)
	if err != nil {
		result.Failed(ctx, constant.FILE_OPERATION_ERROR, "打开Excel文件失败")
		return
	}
	defer f.Close()

	// 读取Excel数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		result.Failed(ctx, constant.FILE_OPERATION_ERROR, "读取Excel数据失败")
		return
	}

	// 解析Excel数据
	var hosts []model.ExcelHostTemplate
	for i, row := range rows {
		if i == 0 { // 跳过表头
			continue
		}
		// Excel列顺序: 主机别名, SSH地址, SSH端口, SSH用户, 备注
		if len(row) < 4 { // 至少需要前4列：主机别名、SSH地址、SSH端口、SSH用户
			continue
		}

		// 获取备注（如果存在）
		remark := ""
		if len(row) >= 5 {
			remark = strings.TrimSpace(row[4])
		}

		host := model.ExcelHostTemplate{
			HostAlias: strings.TrimSpace(row[0]),
			SSHIP:     strings.TrimSpace(row[1]),
			SSHPort:   util.StringToInt(row[2]),
			SSHName:   strings.TrimSpace(row[3]),
			Remark:    remark,
		}
		hosts = append(hosts, host)
	}

	// 调用服务层批量导入
	dto := model.ImportHostsFromExcelDto{
		GroupID: groupId,
		File:    file.Filename,
	}
	c.service.ImportHostsFromExcel(ctx, &dto, hosts)
}

// 更新主机
// @Summary 更新主机
// @Description 更新主机
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.UpdateCmdbHostDto true "主机信息(包含ID)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostupdate [put]
// @Security ApiKeyAuth
func (c *CmdbHostController) UpdateCmdbHost(ctx *gin.Context) {
	var dto model.UpdateCmdbHostDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.UpdateCmdbHost(ctx, dto.ID, &dto)
}

// 删除主机
// @Summary 删除主机
// @Description 删除主机
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CmdbHostIdDto true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostdelete [delete]
// @Security ApiKeyAuth
func (c *CmdbHostController) DeleteCmdbHost(ctx *gin.Context) {
	var dto model.CmdbHostIdDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.DeleteCmdbHost(ctx, dto.ID)
}

// 根据ID获取主机
// @Summary 根据ID获取主机
// @Description 根据ID获取主机
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param id query int true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostinfo [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostById(ctx *gin.Context) {
	id := util.StringToUint(ctx.Query("id"))
	c.service.GetCmdbHostById(ctx, id)
}

// 根据分组ID获取主机列表
// @Summary 根据分组ID获取主机列表
// @Description 根据分组ID获取主机列表（包括所有子分组的主机）
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param groupId query int true "分组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostgroup [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostsByGroupId(ctx *gin.Context) {
	groupId := util.StringToUint(ctx.Query("groupId"))
	c.service.GetCmdbHostsByGroupId(ctx, groupId)
}

// 根据主机名称模糊查询
// @Summary 根据主机名称模糊查询
// @Description 根据主机名称模糊查询
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param name query string true "主机名称(模糊匹配)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostbyname [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostsByHostNameLike(ctx *gin.Context) {
	name := ctx.Query("name")
	c.service.GetCmdbHostsByHostNameLike(ctx, name)
}

// 根据IP查询主机
// @Summary 根据IP查询主机
// @Description 根据IP查询主机(匹配内网IP、公网IP或SSH IP)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param ip query string true "IP地址"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostbyip [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostsByIP(ctx *gin.Context) {
	ip := ctx.Query("ip")
	c.service.GetCmdbHostsByIP(ctx, ip)
}

// 根据状态查询主机
// @Summary 根据状态查询主机
// @Description 根据状态查询主机(1->认证成功,2->未认证,3->认证失败)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param status query int true "状态(1/2/3)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostbystatus [get]
// @Security ApiKeyAuth
func (c *CmdbHostController) GetCmdbHostsByStatus(ctx *gin.Context) {
	status := int(util.StringToUint(ctx.Query("status")))
	c.service.GetCmdbHostsByStatus(ctx, status)
}

// 同步主机基本信息
// @Summary 同步主机基本信息
// @Description 根据主机ID自动同步获取目标主机的基本信息(主机名称、操作系统、CPU、内存、磁盘、内网IP、公网IP)
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.CmdbHostIdDto true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostsync [post]
// @Security ApiKeyAuth
func (c *CmdbHostController) SyncHostInfo(ctx *gin.Context) {
	var dto model.CmdbHostIdDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.SyncHostInfo(ctx, dto.ID)
}