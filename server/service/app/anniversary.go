package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	"go.uber.org/zap"
)

type AnniversaryService struct{}

// Create 创建纪念日
func (s *AnniversaryService) Create(req app.Anniversary) (err error) {
	err = global.GVA_DB.Create(&req).Error
	if err != nil {
		global.GVA_LOG.Error("创建纪念日失败", zap.Error(err))
		return err
	}
	return
}

// Delete 删除纪念日
func (s *AnniversaryService) Delete(id uint) (err error) {
	err = global.GVA_DB.Delete(&app.Anniversary{}, "id = ?", id).Error
	if err != nil {
		global.GVA_LOG.Error("删除纪念日失败", zap.Error(err))
		return
	}
	return
}

// Update 更新纪念日
func (s *AnniversaryService) Update(req app.Anniversary) (err error) {
	err = global.GVA_DB.Model(&app.Anniversary{}).Where("id = ?", req.ID).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新纪念日失败", zap.Error(err))
		return err
	}
	return
}

// GetAnniversaryList 获取纪念日列表
func (s *AnniversaryService) GetAnniversaryList(pageInfo request.GetJNRList) (list []app.Anniversary, total int64, err error) {

	if pageInfo.Page <= 0 && pageInfo.PageSize <= 0 {
		pageInfo.PageSize = 100 // 默认每页10条
	}

	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&app.Anniversary{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	db.Where("user_id = ? OR lover_id = ?", pageInfo.UserId, pageInfo.UserId)

	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取纪念日列表失败", zap.Error(err))
		return
	}
	return
}

// GetAnniversary 获取纪念日
func (s *AnniversaryService) GetAnniversary(id string) (anniversary app.Anniversary, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&anniversary).Error
	if err != nil {
		global.GVA_LOG.Error("获取纪念日失败", zap.Error(err))
		return
	}
	return
}
