package controller

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CmdbSQLController struct {
	service *service.CmdbSQLService
}

func NewCmdbSQLController() *CmdbSQLController {
	sqlService := service.NewCmdbSQLService(dao.NewCmdbSQLDao(common.GetDB()))
	return &CmdbSQLController{service: sqlService}
}

// CreateDatabase godoc
// @Summary 创建数据库
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param data body model.CmdbSQL true "数据库信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database [post]
// @Security ApiKeyAuth
func (c *CmdbSQLController) CreateDatabase(ctx *gin.Context) {
	var db model.CmdbSQL
	if err := ctx.ShouldBindJSON(&db); err != nil {
		result.Failed(ctx, 400, "参数错误")
		return
	}

	if err := c.service.CreateDatabase(&db); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	result.Success(ctx, db)
}

// UpdateDatabase godoc
// @Summary 修改数据库记录
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param data body model.CmdbSQL true "数据库信息(必须包含ID)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database [put]
// @Security ApiKeyAuth
func (c *CmdbSQLController) UpdateDatabase(ctx *gin.Context) {
	var db model.CmdbSQL
	if err := ctx.ShouldBindJSON(&db); err != nil {
		result.Failed(ctx, 400, "参数错误")
		return
	}

	if db.ID == 0 {
		result.Failed(ctx, 400, "ID不能为空")
		return
	}

	// 获取现有记录以保留created_at
	existing, err := c.service.GetDatabase(db.ID)
	if err != nil {
		result.Failed(ctx, 500, "获取记录失败: "+err.Error())
		return
	}

	// 保留原有时间字段
	db.CreatedAt = existing.CreatedAt
	db.UpdatedAt = util.HTime{Time: time.Now()}

	if err := c.service.UpdateDatabase(&db); err != nil {
		result.Failed(ctx, 500, "更新失败: "+err.Error())
		return
	}

	result.Success(ctx, db)
}

// DeleteDatabase godoc
// @Summary 删除数据库记录
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param id query int false "数据库ID（query参数）"
// @Param body body object false "请求体（包含id字段）"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database [delete]
// @Security ApiKeyAuth
func (c *CmdbSQLController) DeleteDatabase(ctx *gin.Context) {
	var id int
	var err error

	// 优先从 query 参数获取 id
	if queryID := ctx.Query("id"); queryID != "" {
		id, err = strconv.Atoi(queryID)
		if err != nil {
			result.Failed(ctx, 400, "ID格式错误")
			return
		}
	} else {
		// 如果 query 中没有，尝试从 body 中获取
		var req struct {
			ID int `json:"id" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			result.Failed(ctx, 400, "ID不能为空或格式错误")
			return
		}
		id = req.ID
	}

	if id == 0 {
		result.Failed(ctx, 400, "ID不能为0")
		return
	}

	if err := c.service.DeleteDatabase(uint(id)); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// GetDatabase godoc
// @Summary 根据ID获取数据库详情
// @Description 根据ID获取数据库详情
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param id query int true "数据库ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database/info [get]
// @Security ApiKeyAuth
func (c *CmdbSQLController) GetDatabase(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil || id == 0 {
		result.Failed(ctx, 400, "ID不能为空或格式错误")
		return
	}

	database, err := c.service.GetDatabase(uint(id))
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, database)
}

// ListDatabases godoc
// @Summary 获取数据库列表[分页]
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/databaselist [get]
// @Security ApiKeyAuth
func (c *CmdbSQLController) ListDatabases(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	dbs, count, err := c.service.ListDatabases(page, pageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"list":  dbs,
		"total": count,
	})
}

// GetDatabasesByName godoc
// @Summary 根据名称查询数据库
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param name query string true "数据库名称"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database/byname [get]
// @Security ApiKeyAuth
func (c *CmdbSQLController) GetDatabasesByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		result.Failed(ctx, 400, "名称不能为空")
		return
	}

	dbs, err := c.service.GetDatabasesByName(name)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, dbs)
}

// GetDatabasesByType godoc
// @Summary 根据类型查询数据库
// @Tags CMDB数据库
// @Accept json
// @Produce json
// @Param type query int true "数据库类型(1=MySQL 2=PostgreSQL 3=Redis 4=MongoDB 5=Elasticsearch)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/database/bytype [get]
// @Security ApiKeyAuth
func (c *CmdbSQLController) GetDatabasesByType(ctx *gin.Context) {
	dbType, err := strconv.Atoi(ctx.Query("type"))
	if err != nil || dbType < 1 || dbType > 5 {
		result.Failed(ctx, 400, "无效的数据库类型")
		return
	}

	dbs, err := c.service.GetDatabasesByType(dbType)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, dbs)
}
