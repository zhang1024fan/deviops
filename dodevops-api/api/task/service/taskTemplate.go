package service

import (
	"fmt"
	"net/http"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	sysmodel "dodevops-api/api/system/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type TaskServiceInterface interface {
	CreateTemplate(ctx *gin.Context, template *model.TaskTemplate) *result.Result
	GetAllTemplates(ctx *gin.Context) *result.Result
	UpdateTemplate(ctx *gin.Context, template *model.TaskTemplate) *result.Result
	DeleteTemplate(ctx *gin.Context, id uint) *result.Result
	GetTemplateByID(ctx *gin.Context, id uint) *result.Result
	GetTemplateContent(ctx *gin.Context, id uint)
	GetTemplatesByName(ctx *gin.Context, name string) *result.Result
	GetTemplatesByType(ctx *gin.Context, templateType int) *result.Result
}

type TaskService struct {
	hostSSHService    service.CmdbHostSSHServiceInterface
	taskTemplateDao   dao.TaskTemplateDaoInterface
}



func NewTaskServiceImpl(
	hostSSHService service.CmdbHostSSHServiceInterface,
	taskTemplateDao dao.TaskTemplateDaoInterface,

) *TaskService {
	return &TaskService{
		hostSSHService:    hostSSHService,
		taskTemplateDao:   taskTemplateDao,

	}
}

func (s *TaskService) CreateTemplate(ctx *gin.Context, template *model.TaskTemplate) *result.Result {
	// 检查模板名称是否已存在
	existingTemplates, err := s.taskTemplateDao.GetTemplatesByName(template.Name)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "检查模板名称失败")
		return nil
	}
	if len(existingTemplates) > 0 {
		result.Failed(ctx, http.StatusConflict, "模板名称已存在")
		return nil
	}

	// 获取当前用户
	userObj, exists := ctx.Get(constant.ContextKeyUserObj)
	if !exists {
		result.Failed(ctx, http.StatusUnauthorized, "未获取到用户信息")
		return nil
	}
	admin, ok := userObj.(*sysmodel.JwtAdmin)
	if !ok {
		result.Failed(ctx, http.StatusUnauthorized, "无效的用户信息格式")
		return nil
	}
	template.CreatedBy = admin.Username

	if err := s.taskTemplateDao.CreateTemplate(template); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, fmt.Sprintf("创建任务模板失败: %v", err))
		return nil
	}
	
	result.Success(ctx, gin.H{
		"id":      template.ID,
		"name":    template.Name,
		"type":    template.Type,
		"message": "任务模板创建成功",
	})
	return nil
}

func (s *TaskService) GetAllTemplates(ctx *gin.Context) *result.Result {
	templates, err := s.taskTemplateDao.GetAllTemplates()
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取模板列表失败")
		return nil
	}

	// 格式化时间
	var formattedTemplates []gin.H
	for _, t := range templates {
		formattedTemplates = append(formattedTemplates, gin.H{
			"id":        t.ID,
			"name":      t.Name,
			"type":      t.Type,
			"content":   t.Content,
			"remark":    t.Remark,
			"createdBy": t.CreatedBy,
			"updatedBy": t.UpdatedBy,
			"createdAt": t.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": t.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	result.Success(ctx, formattedTemplates)
	return nil
}

func (s *TaskService) UpdateTemplate(ctx *gin.Context, template *model.TaskTemplate) *result.Result {
	// 检查模板是否存在
	_, err := s.taskTemplateDao.GetTemplateByID(template.ID)
	if err != nil {
		result.Failed(ctx, http.StatusNotFound, "模板不存在")
		return nil
	}

	// 获取当前用户
	userObj, exists := ctx.Get(constant.ContextKeyUserObj)
	if !exists {
		result.Failed(ctx, http.StatusUnauthorized, "未获取到用户信息")
		return nil
	}
	admin, ok := userObj.(*sysmodel.JwtAdmin)
	if !ok {
		result.Failed(ctx, http.StatusUnauthorized, "无效的用户信息格式")
		return nil
	}
	template.UpdatedBy = admin.Username

	if err := s.taskTemplateDao.UpdateTemplate(template); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, fmt.Sprintf("更新模板失败: %v", err))
		return nil
	}
	result.Success(ctx, gin.H{
		"id":      template.ID,
		"name":    template.Name,
		"type":    template.Type,
		"remark":  template.Remark,
		"message": "模板更新成功",
	})
	return nil
}

func (s *TaskService) DeleteTemplate(ctx *gin.Context, id uint) *result.Result {
	if err := s.taskTemplateDao.DeleteTemplate(id); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "删除模板失败")
		return nil
	}
	result.Success(ctx, gin.H{
		"id":      id,
		"message": "模板删除成功",
	})
	return nil
}

func (s *TaskService) GetTemplateByID(ctx *gin.Context, id uint) *result.Result {
	template, err := s.taskTemplateDao.GetTemplateByID(id)
	if err != nil {
		result.Failed(ctx, http.StatusNotFound, "模板不存在")
		return nil
	}

	// 格式化时间
	formattedTemplate := gin.H{
		"id":        template.ID,
		"name":      template.Name,
		"type":      template.Type,
		"content":   template.Content,
		"remark":    template.Remark,
		"createdBy": template.CreatedBy,
		"updatedBy": template.UpdatedBy,
		"createdAt": template.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt": template.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	result.Success(ctx, formattedTemplate)
	return nil
}

func (s *TaskService) GetTemplateContent(ctx *gin.Context, id uint) {
	template, err := s.taskTemplateDao.GetTemplateByID(id)
	if err != nil {
		ctx.String(http.StatusNotFound, "模板不存在")
		return
	}
	ctx.String(http.StatusOK, template.Content)
}

func (s *TaskService) GetTemplatesByName(ctx *gin.Context, name string) *result.Result {
	templates, err := s.taskTemplateDao.GetTemplatesByName(name)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询模板失败")
		return nil
	}
	result.Success(ctx, templates)
	return nil
}

func (s *TaskService) GetTemplatesByType(ctx *gin.Context, templateType int) *result.Result {
	templates, err := s.taskTemplateDao.GetTemplatesByType(templateType)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询模板失败")
		return nil
	}
	result.Success(ctx, templates)
	return nil
}
