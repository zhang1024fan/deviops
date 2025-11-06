package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type ISysDeptService interface {
	GetSysDeptList(c *gin.Context, DeptName, DeptStatus string) // 查询部门列表
	CreateSysDept(c *gin.Context, sysDept model.SysDept)        // 新增部门
	GetSysDeptById(c *gin.Context, Id int)                      // 根据id查询部门
	UpdateSysDept(c *gin.Context, dept model.SysDept)           // 修改部门
	DeleteSysDeptById(c *gin.Context, dto model.SysDeptIdDto)   // 删除部门
	QuerySysDeptVoList(c *gin.Context)                          // 部门下拉列表
	GetDeptUsers(c *gin.Context, deptId int)                    // 获取某部门下的所有用户
}

type SysDeptServiceImpl struct{}

// 查询部门列表
func (s SysDeptServiceImpl) GetSysDeptList(c *gin.Context, DeptName, DeptStatus string) {
	result.Success(c, dao.GetSysDeptList(DeptName, DeptStatus))
}

// 新增部门
func (s SysDeptServiceImpl) CreateSysDept(c *gin.Context, sysDept model.SysDept) {
	bool := dao.CreateSysDept(sysDept)
	if !bool {
		result.Failed(c, int(result.ApiCode.DEPTISEXIST), result.ApiCode.GetMessage(result.ApiCode.DEPTISEXIST))
		return
	}
	result.Success(c, true)
}

// 根据id查询部门
func (s SysDeptServiceImpl) GetSysDeptById(c *gin.Context, Id int) {
	result.Success(c, dao.GetSysDeptById(Id))
}

// 修改部门
func (s SysDeptServiceImpl) UpdateSysDept(c *gin.Context, dept model.SysDept) {
	sysDept := dao.UpdateSysDept(dept)
	result.Success(c, sysDept)
}

// 删除部门
func (s SysDeptServiceImpl) DeleteSysDeptById(c *gin.Context, dto model.SysDeptIdDto) {
	bool := dao.DeleteSysDeptById(dto)
	if !bool {
		result.Failed(c, int(result.ApiCode.DEPTISDISTRIBUTE), result.ApiCode.GetMessage(result.ApiCode.DEPTISDISTRIBUTE))
		return
	}
	result.Success(c, true)
}

// 部门下拉列表
func (s SysDeptServiceImpl) QuerySysDeptVoList(c *gin.Context) {
	result.Success(c, dao.QuerySysDeptVoList())
}

// 获取某部门下的所有用户
func (s SysDeptServiceImpl) GetDeptUsers(c *gin.Context, deptId int) {
	users := dao.GetUsersByDeptId(deptId)
	result.Success(c, users)
}

var sysDeptService = SysDeptServiceImpl{}

func SysDeptService() ISysDeptService {
	return &sysDeptService
}
