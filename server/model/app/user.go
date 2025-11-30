package app

import "git.echol.cn/loser/menu-server/global"

type User struct {
	global.GVA_MODEL
	NickName string `json:"nick_name" form:"nick_name" gorm:"comment:用户昵称"`
	Phone    string `json:"phone" form:"phone" gorm:"comment:用户手机号"`
	UnionId  string `json:"union_id" form:"union_id" gorm:"comment:微信UnionId"`
	OpenId   string `gorm:"column:open_id;default:'';comment:openId" json:"-"`
	Password string `json:"password" form:"password" gorm:"comment:用户密码"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"comment:用户头像"`
	Status   int8   `gorm:"column:status;default:1;NOT NULL;comment:用户状态  0 封禁 1 正常" json:"status" form:"status"`
	Role     int8   `gorm:"column:role;default:1;NOT NULL;comment:用户角色  1 饲养员 2 吃货" json:"role" form:"role"`
	LoverId  uint   `json:"lover_id" form:"lover_id" gorm:"comment:爱人ID"`
	InivCode string `json:"iniv_code" form:"iniv_code" gorm:"comment:邀请码"`
}

func (User) TableName() string {
	return "app_user"
}
