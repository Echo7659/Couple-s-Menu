package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	common "git.echol.cn/loser/menu-server/model/common/request"
	"go.uber.org/zap"
)

type DiaryService struct{}

// Create 新建日记
func (s *DiaryService) Create(diary app.Diary) error {
	if diary.UserName == "" {
		global.GVA_DB.Model(&app.User{}).Select("nick_name").Where("id = ?", diary.UserId).Scan(&diary.UserName)
	}
	if diary.Avatar == "" {
		global.GVA_DB.Model(&app.User{}).Select("avatar").Where("id = ?", diary.UserId).Scan(&diary.Avatar)
	}

	if err := global.GVA_DB.Create(&diary).Error; err != nil {
		global.GVA_LOG.Error("创建日记失败", zap.Error(err))
		return err
	}
	return nil
}

// Update 更新日记
func (s *DiaryService) Update(diary app.Diary) error {
	if err := global.GVA_DB.Model(&app.Diary{}).Where("id = ?", diary.ID).Updates(&diary).Error; err != nil {
		global.GVA_LOG.Error("更新日记失败", zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除日记
func (s *DiaryService) Delete(id int) error {
	if err := global.GVA_DB.Where("id = ?", id).Delete(&app.Diary{}).Error; err != nil {
		global.GVA_LOG.Error("删除日记失败", zap.Error(err))
		return err
	}
	return nil
}

// GetDiary 获取日记详情
func (s *DiaryService) GetDiary(id uint) (app.Diary, error) {
	var diary app.Diary
	if err := global.GVA_DB.Where("id = ?", id).First(&diary).Error; err != nil {
		global.GVA_LOG.Error("获取日记详情失败", zap.Error(err))
		return diary, err
	}
	return diary, nil
}

// GetDiaryList 获取日记列表
func (s *DiaryService) GetDiaryList(req common.PageInfo) ([]app.Diary, int64, error) {
	var diaries []app.Diary
	var total int64
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.GVA_DB.Model(&app.Diary{})
	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		global.GVA_LOG.Error("获取日记总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 查询用户Lover信息
	if req.UserId != 0 {
		loverId := 0
		err := global.GVA_DB.Table("app_user").Select("lover_id").Where("id =  ?", req.UserId).Scan(&loverId).Error
		if err != nil {
			global.GVA_LOG.Error("获取LoverId失败")
		}

		if err := db.Limit(limit).Offset(offset).Where("user_id = ? OR user_id = ?", req.UserId, loverId).Find(&diaries).Error; err != nil {
			global.GVA_LOG.Error("获取日记列表失败", zap.Error(err))
			return nil, 0, err
		}

		return diaries, total, nil
	}
	return diaries, total, nil
}
