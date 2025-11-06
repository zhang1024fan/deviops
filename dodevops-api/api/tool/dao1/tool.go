// 快捷导航工具 DAO层
// author xiaoRui
package dao

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
	"time"
)

// CreateTool 创建导航工具
func CreateTool(dto model.AddToolDto) error {
	tool := model.Tool{
		Title:      dto.Title,
		Icon:       dto.Icon,
		Link:       dto.Link,
		Sort:       dto.Sort,
		CreateTime: util.HTime{Time: time.Now()},
		UpdateTime: time.Now(),
	}
	return Db.Create(&tool).Error
}

// GetToolByID 根据ID获取导航工具
func GetToolByID(id uint) (tool model.Tool, err error) {
	err = Db.Where("id = ?", id).First(&tool).Error
	return
}

// UpdateTool 更新导航工具
func UpdateTool(dto model.UpdateToolDto) error {
	updates := map[string]interface{}{
		"title":       dto.Title,
		"icon":        dto.Icon,
		"link":        dto.Link,
		"sort":        dto.Sort,
		"update_time": time.Now(),
	}
	return Db.Model(&model.Tool{}).Where("id = ?", dto.ID).Updates(updates).Error
}

// DeleteTool 删除导航工具
func DeleteTool(id uint) error {
	return Db.Delete(&model.Tool{}, id).Error
}

// GetToolList 获取导航工具列表（分页）
func GetToolList(dto model.ToolQueryDto) (tools []model.Tool, total int64, err error) {
	db := Db.Model(&model.Tool{})

	// 条件查询
	if dto.Title != "" {
		db = db.Where("title LIKE ?", "%"+dto.Title+"%")
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询
	if dto.PageNum > 0 && dto.PageSize > 0 {
		offset := (dto.PageNum - 1) * dto.PageSize
		db = db.Offset(offset).Limit(dto.PageSize)
	}

	// 按排序和创建时间排序
	err = db.Order("sort ASC, create_time DESC").Find(&tools).Error
	return
}

// GetAllTools 获取所有导航工具（不分页）
func GetAllTools() (tools []model.Tool, err error) {
	err = Db.Order("sort ASC, create_time DESC").Find(&tools).Error
	return
}
