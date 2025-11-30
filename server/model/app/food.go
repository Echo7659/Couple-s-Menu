package app

import "git.echol.cn/loser/menu-server/global"

type Food struct {
	global.GVA_MODEL
	Name       string `json:"name" gorm:"comment:菜品名称"`         // 菜品名称
	Desc       string `json:"desc" gorm:"comment:菜品描述"`         //
	CategoryId uint   `json:"categoryId" gorm:"comment:菜品分类ID"` // 菜品分类ID
	Img        string `json:"img" gorm:"comment:菜品图片"`          // 菜品图片
}

// TableName Food表
func (Food) TableName() string {
	return "app_food"
}
