// 通用返回结构
// author xiaoRui

package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 结构体
type Result struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回的数据
}

// 分页返回结构
type PageResult struct {
	Total    int64       `json:"total"`    // 总记录数
	List     interface{} `json:"list"`     // 数据列表
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页数量
}

// 返回成功
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Code = int(ApiCode.SUCCESS)
	res.Message = ApiCode.GetMessage(ApiCode.SUCCESS)
	res.Data = data
	c.JSON(http.StatusOK, res)
}

// 返回分页成功
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	res := Result{
		Code:    int(ApiCode.SUCCESS),
		Message: ApiCode.GetMessage(ApiCode.SUCCESS),
		Data: PageResult{
			Total:    total,
			List:     list,
			Page:     page,
			PageSize: pageSize,
		},
	}
	c.JSON(http.StatusOK, res)
}

// 返回失败
func Failed(c *gin.Context, code int, message string) {
	res := Result{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}
	c.JSON(http.StatusOK, res)
}

// 新增方法：支持直接传入 code 和 message
func FailedWithCode(c *gin.Context, code int, message string) {
	Failed(c, code, message)
}
