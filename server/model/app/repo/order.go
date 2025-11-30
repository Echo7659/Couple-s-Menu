package repo

import "git.echol.cn/loser/menu-server/model/app"

type OrderVo struct {
	ID     uint `json:"id"`
	Foods  []app.Food
	Remark string `json:"remark" `
	Date   string `json:"date" ` // 下单日期
}
