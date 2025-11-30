package app

import "git.echol.cn/loser/menu-server/global"

type Menu struct {
	global.GVA_MODEL
	Name   string `json:"name" gorm:"comment:菜单名称"`
	UserId uint   `json:"user_id" gorm:"comment:用户ID"`
	// 封面
	Cover string `json:"cover" gorm:"comment:封面"`
	// 描述
	Description string `json:"description" gorm:"comment:描述"`
}

func (Menu) TableName() string {
	return "app_menus"
}
