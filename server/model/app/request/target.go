package request

import common "git.echol.cn/loser/menu-server/model/common/request"

type GetTargetList struct {
	common.PageInfo
	CreatorId   uint   `json:"creatorId" form:"creatorId"`     // 创建者ID
	LoverId     uint   `json:"loverId" form:"loverId"`         // 情侣ID
	IsCompleted *bool  `json:"isCompleted" form:"isCompleted"` // 是否已完成
	Title       string `json:"title" form:"title"`             // 目标标题
}

type CreateTargetRecord struct {
	TargetId   uint   `json:"targetId" binding:"required"`  // 目标ID
	UserId     uint   `json:"userId" binding:"required"`    // 打卡用户ID
	RecordImg  string `json:"recordImg" binding:"required"` // 打卡记录图片URL
	RecordNote string `json:"recordNote"`                   // 打卡备注
}
