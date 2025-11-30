package app

import (
	"errors"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	"go.uber.org/zap"
)

type TargetService struct{}

// Create 创建打卡目标（情侣共享）
func (t *TargetService) Create(req app.Target) (err error) {
	// 查询创建者的情侣ID
	var user app.User
	if err := global.GVA_DB.Where("id = ?", req.CreatorId).First(&user).Error; err != nil {
		global.GVA_LOG.Error("查询用户信息失败", zap.Error(err))
		return errors.New("查询用户信息失败")
	}
	if user.LoverId == 0 {
		return errors.New("您还未绑定情侣，无法创建打卡目标")
	}
	// 自动设置情侣ID
	req.LoverId = user.LoverId

	err = global.GVA_DB.Create(&req).Error
	if err != nil {
		global.GVA_LOG.Error("创建打卡目标失败", zap.Error(err))
		return err
	}
	return
}

// Delete 删除打卡目标
func (t *TargetService) Delete(id uint) (err error) {
	err = global.GVA_DB.Where("id = ?", id).Delete(&app.Target{}).Error
	if err != nil {
		global.GVA_LOG.Error("删除打卡目标失败", zap.Error(err))
		return err
	}
	return
}

// Update 更新打卡目标
func (t *TargetService) Update(req app.Target) (err error) {
	err = global.GVA_DB.Model(&app.Target{}).Where("id = ?", req.ID).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新打卡目标失败", zap.Error(err))
		return err
	}
	return
}

// GetTarget 获取打卡目标详情
func (t *TargetService) GetTarget(id string) (target app.Target, err error) {
	err = global.GVA_DB.Preload("Records").Preload("Records.User").Where("id = ?", id).First(&target).Error
	if err != nil {
		global.GVA_LOG.Error("获取打卡目标详情失败", zap.Error(err))
		return target, err
	}
	return target, nil
}

// GetTargetList 获取打卡目标列表（情侣共享）
func (t *TargetService) GetTargetList(info request.GetTargetList) (list []app.Target, total int64, err error) {
	db := global.GVA_DB.Model(&app.Target{})

	// 如果指定了loverId，查询该用户作为创建者的目标，以及情侣创建的目标
	if info.LoverId != 0 {
		db = db.Where("creator_id = ? OR lover_id = ?", info.LoverId, info.LoverId)
	}

	// 根据创建者ID筛选
	if info.CreatorId != 0 {
		db = db.Where("creator_id = ?", info.CreatorId)
	}

	// 根据完成状态筛选
	if info.IsCompleted != nil {
		db = db.Where("is_completed = ?", *info.IsCompleted)
	}

	// 根据目标标题模糊搜索
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取打卡目标列表失败", zap.Error(err))
		return list, total, err
	}

	// 分页查询（列表不需要加载打卡记录，提高性能）
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取打卡目标列表失败", zap.Error(err))
		return list, total, err
	}
	return list, total, nil
}

// GetMyTargets 获取我的目标（包括我创建的和情侣创建的）
func (t *TargetService) GetMyTargets(userId uint, info request.GetTargetList) (list []app.Target, total int64, err error) {
	// 查询当前用户的情侣ID
	var user app.User
	if err := global.GVA_DB.Where("id = ?", userId).First(&user).Error; err != nil {
		global.GVA_LOG.Error("查询用户信息失败", zap.Error(err))
		return list, total, err
	}

	db := global.GVA_DB.Model(&app.Target{})

	// 查询我创建的目标 或者 情侣创建的目标
	if user.LoverId != 0 {
		db = db.Where("creator_id = ? OR creator_id = ?", userId, user.LoverId)
	} else {
		db = db.Where("creator_id = ?", userId)
	}

	// 根据完成状态筛选
	if info.IsCompleted != nil {
		db = db.Where("is_completed = ?", *info.IsCompleted)
	}

	// 根据目标标题模糊搜索
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取我的目标列表失败", zap.Error(err))
		return list, total, err
	}

	// 分页查询（列表不需要加载打卡记录）
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取我的目标列表失败", zap.Error(err))
		return list, total, err
	}
	return list, total, nil
}

// CreateRecord 创建打卡记录
func (t *TargetService) CreateRecord(req request.CreateTargetRecord) (err error) {
	// 先查询目标是否存在
	var target app.Target
	if err := global.GVA_DB.Where("id = ?", req.TargetId).First(&target).Error; err != nil {
		global.GVA_LOG.Error("目标不存在", zap.Error(err))
		return errors.New("目标不存在")
	}

	// 创建打卡记录
	record := app.TargetRecord{
		TargetId:   req.TargetId,
		UserId:     req.UserId,
		RecordImg:  req.RecordImg,
		RecordNote: req.RecordNote,
	}

	err = global.GVA_DB.Create(&record).Error
	if err != nil {
		global.GVA_LOG.Error("创建打卡记录失败", zap.Error(err))
		return err
	}
	return
}

// GetTargetRecords 获取目标的所有打卡记录
func (t *TargetService) GetTargetRecords(targetId uint) (records []app.TargetRecord, err error) {
	err = global.GVA_DB.Preload("User").Where("target_id = ?", targetId).Order("created_at desc").Find(&records).Error
	if err != nil {
		global.GVA_LOG.Error("获取打卡记录失败", zap.Error(err))
		return records, err
	}
	return records, nil
}

// CompleteTarget 完成打卡目标（标记为已完成）
func (t *TargetService) CompleteTarget(id uint) (err error) {
	err = global.GVA_DB.Model(&app.Target{}).Where("id = ?", id).Update("is_completed", true).Error
	if err != nil {
		global.GVA_LOG.Error("完成打卡目标失败", zap.Error(err))
		return err
	}
	return
}
