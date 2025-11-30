package repo

import "git.echol.cn/loser/menu-server/model/app"

type LoginVo struct {
	Info  app.User
	Token string `json:"token"`
	// 是否首次登录
	IsFirstLogin bool `json:"isFirstLogin"`
}
