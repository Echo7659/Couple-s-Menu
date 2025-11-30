package app

import "git.echol.cn/loser/menu-server/global"

type Order struct {
	global.GVA_MODEL
	FoodIds string `json:"food_ids" gorm:"comment:菜品ID列表;type:varchar(255);NOT NULL"` // 菜品ID列表，逗号分隔
	UserId  uint   `json:"user_id" gorm:"comment:用户ID"`
	Remark  string `json:"remark" gorm:"comment:留言备注"`
	Status  int    `json:"status" gorm:"comment:订单状态;default:1"`               // 订单状态 1-待接单 2-制作中 3-已完成 4-已取消
	Date    string `json:"date" gorm:"comment:下单日期;type:varchar(20);NOT NULL"` // 下单日期
}

// TableName Order表
func (Order) TableName() string {
	return "app_order"
}
