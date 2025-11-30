package app

import "git.echol.cn/loser/menu-server/global"

// Anniversary 纪念日
type Anniversary struct {
	global.GVA_MODEL
	UserId  uint   `json:"user_id" gorm:"comment:用户ID"`         // 用户ID
	LoverId uint   `json:"lover_id" gorm:"comment:爱人ID"`        // 爱人ID
	Title   string `json:"title" gorm:"comment:纪念日标题"`          // 纪念日标题
	Date    string `json:"date" gorm:"comment:纪念日日期"`           // 纪念日日期
	Type    int    `json:"type" gorm:"comment:纪念日类型;default:1"` // 纪念日类型 1-恋爱纪念日 2-生日 3-其他

}

// TableName Anniversary表
func (Anniversary) TableName() string {
	return "app_anniversary"
}
