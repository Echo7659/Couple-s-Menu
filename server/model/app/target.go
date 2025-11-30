package app

import "git.echol.cn/loser/menu-server/global"

type Target struct {
	global.GVA_MODEL
	Title       string         `json:"title" gorm:"comment:目标标题"`                     // 目标标题
	Description string         `json:"description" gorm:"comment:目标描述"`               // 目标描述
	CreatorId   uint           `json:"creatorId" gorm:"comment:创建者用户ID"`              // 创建者用户ID
	LoverId     uint           `json:"loverId" gorm:"comment:情侣ID;index"`             // 情侣ID（自动设置为创建者的情侣）
	IsCompleted bool           `json:"isCompleted" gorm:"default:false;comment:是否完成"` // 是否已完成
	Records     []TargetRecord `json:"records" gorm:"foreignKey:TargetId"`            // 打卡记录列表
}

// TargetRecord 打卡记录
type TargetRecord struct {
	global.GVA_MODEL
	TargetId   uint   `json:"targetId" gorm:"comment:目标ID;index"`      // 目标ID
	UserId     uint   `json:"userId" gorm:"comment:打卡用户ID;index"`      // 打卡用户ID
	RecordImg  string `json:"recordImg" gorm:"comment:打卡记录图片URL"`      // 打卡记录图片URL
	RecordNote string `json:"recordNote" gorm:"comment:打卡备注"`          // 打卡备注
	User       *User  `json:"user,omitempty" gorm:"foreignKey:UserId"` // 关联用户信息
}

// TableName Target表
func (Target) TableName() string {
	return "app_target"
}

// TableName TargetRecord表
func (TargetRecord) TableName() string {
	return "app_target_record"
}
