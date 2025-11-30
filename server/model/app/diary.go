package app

import "git.echol.cn/loser/menu-server/global"

type Diary struct {
	global.GVA_MODEL
	Title    string `json:"title" gorm:"comment:标题"`
	Content  string `json:"content" gorm:"comment:内容;type:text"`
	UserId   uint   `json:"userId" gorm:"comment:用户ID"`
	UserName string `json:"userName" gorm:"comment:用户名"` // 用户名
	Avatar   string `json:"avatar" gorm:"comment:用户头像"`  // 用户头像
	// 心情
	Mood int    `json:"mood" gorm:"comment:心情 1-开心 2-兴奋 3-幸福 4-思念 5-心累 6-疲倦 7-失落 8-生气 9-平静;default:1"`
	Img  string `json:"img" gorm:"comment:图片;type:text"` // 图片
}

// TableName 日记表
func (Diary) TableName() string {
	return "app_diary"
}
