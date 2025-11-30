package request

import common "git.echol.cn/loser/menu-server/model/common/request"

type GetFoodList struct {
	common.PageInfo
	CategoryId int    `json:"categoryId" form:"categoryId"` // 分类ID
	MenuId     int    `json:"menuId" form:"menuId"`         // 菜谱ID
	Name       string `json:"name" form:"name"`             // 菜名
}
