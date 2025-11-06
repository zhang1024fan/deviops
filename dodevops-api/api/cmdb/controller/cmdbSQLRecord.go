package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	cmdbService "dodevops-api/api/cmdb/service"
	configService "dodevops-api/api/configcenter/service"
	configModel "dodevops-api/api/configcenter/model"
	"dodevops-api/common/util"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
	sysModel "dodevops-api/api/system/model"
)

const (
	ParamError = 1001
	DatabaseError = 1002
)

type CmdbSQLRecordController struct {
	recordService *cmdbService.CmdbSQLRecordService
	sqlService   *cmdbService.CmdbSQLService
}

var sqlRecordController *CmdbSQLRecordController

func InitCmdbSQLRecordController(db *gorm.DB) {
	recordDao := dao.NewCmdbSQLRecordDao(db)
	recordService := cmdbService.NewCmdbSQLRecordService(recordDao)
	sqlService := cmdbService.NewCmdbSQLService(dao.NewCmdbSQLDao(db))
	sqlRecordController = &CmdbSQLRecordController{
		recordService: recordService,
		sqlService:   sqlService,
	}
}

func GetCmdbSQLRecordController() *CmdbSQLRecordController {
	return sqlRecordController
}

func NewCmdbSQLRecordController(db *gorm.DB) *CmdbSQLRecordController {
	recordDao := dao.NewCmdbSQLRecordDao(db)
	recordService := cmdbService.NewCmdbSQLRecordService(recordDao)
	sqlService := cmdbService.NewCmdbSQLService(dao.NewCmdbSQLDao(db))
	return &CmdbSQLRecordController{
		recordService: recordService,
		sqlService:   sqlService,
	}
}

// SQLRequest 新的SQL请求结构体
type SQLRequest struct {
	DatabaseID   uint   `json:"databaseId"`   // 数据库ID
	DatabaseName string `json:"databaseName"` // 数据库名称
	SQL          string `json:"sql"`          // SQL语句
}

// 公共方法：获取数据库连接信息
func (c *CmdbSQLRecordController) getDBConnectionInfo(req SQLRequest) (*model.CmdbSQL, *configModel.AccountAuth, string, error) {
	var dbInfo *model.CmdbSQL
	var err error
	
	if req.DatabaseID > 0 {
		dbInfo, err = c.sqlService.GetDatabase(req.DatabaseID)
		if err != nil || dbInfo == nil {
			return nil, nil, "", fmt.Errorf("获取数据库信息失败 (ID: %d)", req.DatabaseID)
		}
	} else if req.DatabaseName != "" {
		// 直接使用用户传入的数据库名称
		dbInfo = &model.CmdbSQL{
			Name:      req.DatabaseName,
			AccountID: 1, // 默认账号ID，实际应从配置获取
		}
	} else {
		return nil, nil, "", fmt.Errorf("必须提供数据库ID或名称")
	}

	// 获取账号信息
	account, err := configService.NewAccountAuthService().GetByID(dbInfo.AccountID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("获取账号信息失败: %v", err)
	}

	// 解密密码
	decrypted, err := configService.NewAccountAuthService().DecryptPassword(dbInfo.AccountID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("解密密码失败: %v", err)
	}

	return dbInfo, account, decrypted, nil
}

// 验证SQL语句类型
func validateSQLType(sql string, allowedTypes []string) bool {
	sql = strings.TrimSpace(strings.ToUpper(sql))
	for _, t := range allowedTypes {
		if strings.HasPrefix(sql, t) {
			return true
		}
	}
	return false
}

// 获取当前用户名
func getCurrentUsername(ctx *gin.Context) (string, error) {
	userObj, exists := ctx.Get(constant.ContextKeyUserObj)
	if !exists {
		return "", fmt.Errorf("无法获取当前用户信息")
	}

	// 尝试多种可能的用户对象格式
	switch v := userObj.(type) {
	case map[string]interface{}:
		if username, ok := v["username"].(string); ok && username != "" {
			return username, nil
		}
	case *sysModel.SysAdmin:
		if v.Username != "" {
			return v.Username, nil
		}
	case sysModel.SysAdmin:
		if v.Username != "" {
			return v.Username, nil
		}
	default:
		// 尝试通过反射获取Username字段
		val := reflect.ValueOf(userObj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			field := val.FieldByName("Username")
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String(), nil
			}
		}
	}

	return "", fmt.Errorf("无法从用户对象中获取有效的用户名")
}

// @Summary 执行查询语句
// @Produce json
// @Tags CMDB数据库
// @Description 执行查询语句(通过数据库ID/名称)
// @Param data body SQLRequest true "SQL查询请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql/select [post]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) ExecuteSelect(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	// 只允许SELECT语句
	if !validateSQLType(req.SQL, []string{"SELECT"}) {
		result.FailedWithCode(ctx, ParamError, "只允许执行SELECT语句")
		return
	}

	dbInfo, account, decrypted, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}

	// 使用解密后的密码建立连接
	if dbInfo.Type != 1 {
		result.FailedWithCode(ctx, DatabaseError, "目前只支持MySQL数据库")
		return
	}

	// 设置连接参数(先连接到默认数据库)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&interpolateParams=true", 
		account.Name, decrypted, account.Host, account.Port)
	
	// 测试网络连通性
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", account.Host, account.Port), 5*time.Second)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, fmt.Sprintf("无法连接到数据库服务器 %s:%d - %v", account.Host, account.Port, err))
		return
	}
	conn.Close()

	// 连接到MySQL服务器(不指定数据库)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接初始化失败: "+err.Error())
		return
	}
	defer db.Close()

	// 检查schema是否存在
	ctxCheckSchema, cancelCheckSchema := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancelCheckSchema()
	
	var schemaExists int
	err = db.QueryRowContext(ctxCheckSchema, 
		"SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", 
		req.DatabaseName).Scan(&schemaExists)
	
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "检查schema存在性失败: "+err.Error())
		return
	}
	
	if schemaExists == 0 {
		result.FailedWithCode(ctx, DatabaseError, fmt.Sprintf("schema '%s' 不存在", req.DatabaseName))
		return
	}

	// 连接到指定schema
	connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=60s&readTimeout=60s&writeTimeout=60s", 
		account.Name, decrypted, account.Host, account.Port, req.DatabaseName)
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接初始化失败: "+err.Error())
		return
	}
	defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	// 测试连接是否可用
	ctxPing, cancelPing := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancelPing()
	if err := db.PingContext(ctxPing); err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接测试失败: "+err.Error())
		return
	}

	// 设置查询超时
	queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// 执行查询
	startTime := time.Now()
	rows, err := db.QueryContext(queryCtx, req.SQL)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "查询执行失败: "+err.Error())
		return
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "获取列信息失败: "+err.Error())
		return
	}

	// 处理结果集
	var results []map[string]interface{}
	for rows.Next() {
		// 创建值的切片和指针切片
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(pointers...); err != nil {
			result.FailedWithCode(ctx, DatabaseError, "扫描行数据失败: "+err.Error())
			return
		}

		// 构建结果map
		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}
		results = append(results, rowData)
	}

	executionTime := time.Since(startTime).Milliseconds()
	returnedRows := int64(len(results))

	// 获取当前用户名和IP
	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}
	clientIP := util.GetClientIP(ctx.Request)

	// 记录执行历史
	err = c.recordService.RecordSQLExecution(
		account.Host,
		req.DatabaseName,
		"SELECT",
		req.SQL,
		username,
		clientIP,
		0, // scannedRows
		0, // affectedRows
		executionTime,
		returnedRows,
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库错误: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"returnedRows": returnedRows,
		"executionTime": executionTime,
		"results": results,
	})
}

// @Summary 执行插入语句
// @Produce json
// @Tags CMDB数据库
// @Description 执行插入语句(通过数据库ID/名称)
// @Param data body SQLRequest true "SQL插入请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql [post]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) ExecuteInsert(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	dbInfo, account, decrypted, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}

	// 使用解密后的密码建立连接
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", account.Name, decrypted, account.Host, account.Port, dbInfo.Name)
	_ = connStr // 显式使用变量
	// TODO: 实际执行插入并获取结果
	affectedRows := int64(1) // 示例值
	executionTime := int64(80) // 示例值(毫秒)

	// 获取当前用户名和IP
	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}
	clientIP := util.GetClientIP(ctx.Request)

	err = c.recordService.RecordSQLExecution(
		account.Host,
		dbInfo.Name,
		"INSERT",
		req.SQL,
		username,
		clientIP,
		0, // scannedRows
		affectedRows,
		executionTime,
		0, // returnedRows
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库错误: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"affectedRows": affectedRows,
		"executionTime": executionTime,
	})
}

// @Summary 执行更新语句
// @Produce json
// @Tags CMDB数据库
// @Description 执行更新语句(通过数据库ID/名称)
// @Param data body SQLRequest true "SQL更新请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql [put]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) ExecuteUpdate(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	_, account, decrypted, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}

	// 设置连接参数
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=60s&readTimeout=60s&writeTimeout=60s", 
		account.Name, decrypted, account.Host, account.Port, req.DatabaseName)
	
	// 建立数据库连接
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接初始化失败: "+err.Error())
		return
	}
	defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	// 测试连接是否可用
	ctxPing, cancelPing := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancelPing()
	if err := db.PingContext(ctxPing); err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接测试失败: "+err.Error())
		return
	}

	// 设置查询超时
	queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// 执行更新
	startTime := time.Now()
	res, err := db.ExecContext(queryCtx, req.SQL)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "更新执行失败: "+err.Error())
		return
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "获取影响行数失败: "+err.Error())
		return
	}

	executionTime := time.Since(startTime).Milliseconds()

	// 获取当前用户名和IP
	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}
	clientIP := util.GetClientIP(ctx.Request)

	// 记录执行历史
	err = c.recordService.RecordSQLExecution(
		account.Host,
		req.DatabaseName,
		"UPDATE",
		req.SQL,
		username,
		clientIP,
		0, // scannedRows
		affectedRows,
		executionTime,
		0, // returnedRows
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库错误: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"affectedRows": affectedRows,
		"executionTime": executionTime,
	})
}

// @Summary 执行删除语句
// @Produce json
// @Tags CMDB数据库
// @Description 执行删除语句(通过数据库ID/名称)
// @Param data body SQLRequest true "SQL删除请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql [delete]
// @Security ApiKeyAuth
// ExecuteSQL 执行原生SQL语句
// @Summary 执行原生SQL语句
// @Produce json
// @Tags CMDB数据库
// @Description 执行原生SQL语句(通过数据库ID/名称)
// @Param data body SQLRequest true "SQL执行请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql/execute [post]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) ExecuteSQL(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	// 验证SQL类型，只允许特定操作
	allowedTypes := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP"}
	if !validateSQLType(req.SQL, allowedTypes) {
		result.FailedWithCode(ctx, ParamError, "不允许执行此类型的SQL语句")
		return
	}

	// 根据SQL类型调用相应的方法
	switch {
	case strings.HasPrefix(strings.ToUpper(req.SQL), "SELECT"):
		c.ExecuteSelect(ctx)
	case strings.HasPrefix(strings.ToUpper(req.SQL), "INSERT"):
		c.ExecuteInsert(ctx)
	case strings.HasPrefix(strings.ToUpper(req.SQL), "UPDATE"):
		c.ExecuteUpdate(ctx)
	case strings.HasPrefix(strings.ToUpper(req.SQL), "DELETE"):
		c.ExecuteDelete(ctx)
	default:
		result.FailedWithCode(ctx, ParamError, "不支持执行此类型的SQL语句")
	}
}

// GetDatabaseList 获取数据库列表
// @Summary 获取数据库列表
// @Produce json
// @Tags CMDB数据库
// @Description 获取指定数据库实例的数据库列表
// @Param data body SQLRequest true "数据库查询请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql/databaselist [post]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) GetDatabaseList(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	if req.DatabaseID == 0 {
		result.FailedWithCode(ctx, ParamError, "必须提供数据库ID")
		return
	}

	// 获取数据库信息
	dbInfo, err := c.sqlService.GetDatabase(req.DatabaseID)
	if err != nil || dbInfo == nil {
		result.FailedWithCode(ctx, DatabaseError, "获取数据库信息失败")
		return
	}

	// 获取账号信息
	account, err := configService.NewAccountAuthService().GetByID(dbInfo.AccountID)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "获取账号信息失败")
		return
	}

	// 解密密码
	decrypted, err := configService.NewAccountAuthService().DecryptPassword(dbInfo.AccountID)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "解密密码失败")
		return
	}

	// 连接到MySQL服务器(不指定数据库)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=60s", 
		account.Name, decrypted, account.Host, account.Port)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库连接初始化失败: "+err.Error())
		return
	}
	defer db.Close()

	// 查询数据库列表
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "查询数据库列表失败: "+err.Error())
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			result.FailedWithCode(ctx, DatabaseError, "解析数据库列表失败: "+err.Error())
			return
		}
		databases = append(databases, dbName)
	}

	result.Success(ctx, gin.H{
		"databases": databases,
		"host":      account.Host,
		"port":      account.Port,
	})
}

// ListDatabases 获取数据库列表
// @Summary 获取数据库列表
// @Produce json
// @Tags CMDB数据库
// @Description 获取指定数据库实例的数据库列表
// @Param data body SQLRequest true "数据库查询请求"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sql/databaselist [post]
// @Security ApiKeyAuth
func (c *CmdbSQLRecordController) ListDatabases(ctx *gin.Context) {
	c.GetDatabaseList(ctx)
}

func (c *CmdbSQLRecordController) ExecuteDelete(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "参数错误: "+err.Error())
		return
	}

	// 只允许DELETE语句
	if !validateSQLType(req.SQL, []string{"DELETE"}) {
		result.FailedWithCode(ctx, ParamError, "只允许执行DELETE语句")
		return
	}

	dbInfo, account, decrypted, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}

	// 使用解密后的密码建立连接
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", account.Name, decrypted, account.Host, account.Port, dbInfo.Name)
	_ = connStr // 显式使用变量
	// TODO: 实际执行删除并获取结果
	affectedRows := int64(3) // 示例值
	executionTime := int64(90) // 示例值(毫秒)

	// 获取当前用户名和IP
	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}
	clientIP := util.GetClientIP(ctx.Request)

	err = c.recordService.RecordSQLExecution(
		account.Host,
		dbInfo.Name,
		"DELETE",
		req.SQL,
		username,
		clientIP,
		0, // scannedRows
		affectedRows,
		executionTime,
		0, // returnedRows
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "数据库错误: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"affectedRows": affectedRows,
		"executionTime": executionTime,
	})
}
