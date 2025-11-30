package app

import "git.echol.cn/loser/menu-server/global"

type Lovers struct {
	global.GVA_MODEL
	UserID    uint   `json:"userId" gorm:"comment:用户ID"`
	LoverId   uint   `json:"loverId" gorm:"comment:爱人ID"`
	StartDate string `json:"startDate" gorm:"comment:开始日期"` // 开始日期
}

// TableName Lovers表
func (Lovers) TableName() string {
	return "app_lovers"
}
