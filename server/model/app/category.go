package app

import "git.echol.cn/loser/menu-server/global"

type Category struct {
	global.GVA_MODEL
	Name   string `json:"name" gorm:"comment:分类名称"`   // 分类名称
	MenuId uint   `json:"menuId" gorm:"comment:菜谱ID"` // 菜谱ID
}

// TableName Category表
func (Category) TableName() string {
	return "app_category"
}
