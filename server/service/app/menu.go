package app

import (
	"errors"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	common "git.echol.cn/loser/menu-server/model/common/request"
	"go.uber.org/zap"
)

type MenuService struct{}

// Create 创建菜谱
func (m *MenuService) Create(req app.Menu) (err error) {
	err = global.GVA_DB.Create(&req).Error
	if err != nil {
		global.GVA_LOG.Error("创建菜谱失败", zap.Error(err))
		return err
	}
	return
}

// Delete 删除菜谱
func (m *MenuService) Delete(id uint) (err error) {
	err = global.GVA_DB.Where("id = ?", id).Delete(&app.Menu{}).Error
	if err != nil {
		global.GVA_LOG.Error("删除菜谱失败", zap.Error(err))
		return err
	}
	return
}

// Update 更新菜谱
func (m *MenuService) Update(req app.Menu) (err error) {
	err = global.GVA_DB.Model(&app.Menu{}).Where("id = ?", req.ID).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新菜谱失败", zap.Error(err))
		return err
	}
	return
}

// GetMenu 获取菜谱详情（需要权限验证）
func (m *MenuService) GetMenu(id string, userId uint) (menu app.Menu, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜谱详情失败", zap.Error(err))
		return menu, err
	}

	// 检查权限
	hasPermission, err := m.checkMenuPermission(menu.ID, userId)
	if err != nil {
		return menu, err
	}
	if !hasPermission {
		global.GVA_LOG.Error("无权限查看该菜谱")
		return menu, errors.New("无权限查看该菜谱")
	}

	return menu, nil
}

// GetMenuList 获取菜谱列表（包括自己和恋人的菜谱）
func (m *MenuService) GetMenuList(userId uint, pageInfo common.PageInfo) (list []app.Menu, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)

	// 获取恋人的ID
	loverId := m.getLoverId(userId)

	// 创建db，查询自己或恋人的菜谱
	db := global.GVA_DB.Model(&app.Menu{})
	if loverId > 0 {
		// 如果有恋人，查询自己和恋人的菜谱
		db = db.Where("user_id = ? OR user_id = ?", userId, loverId)
	} else {
		// 如果没有恋人，只查询自己的菜谱
		db = db.Where("user_id = ?", userId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageInfo.Keyword != "" {
		db = db.Where("name LIKE ? ", "%"+pageInfo.Keyword+"%")
	}

	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜谱列表失败", zap.Error(err))
		return list, total, err
	}
	return list, total, nil
}

// getLoverId 获取用户的恋人ID
func (m *MenuService) getLoverId(userId uint) uint {
	var lover app.Lovers
	err := global.GVA_DB.Where("user_id = ? OR lover_id = ?", userId, userId).First(&lover).Error
	if err != nil {
		// 没有恋人关系，返回0
		return 0
	}
	// 判断当前用户是user_id还是lover_id，返回对方的ID
	if lover.UserID == userId {
		return lover.LoverId
	}
	return lover.UserID
}

// checkMenuPermission 检查用户是否有权限查看/操作菜谱
func (m *MenuService) checkMenuPermission(menuId uint, userId uint) (bool, error) {
	var menu app.Menu
	err := global.GVA_DB.Where("id = ?", menuId).First(&menu).Error
	if err != nil {
		return false, err
	}

	// 如果是自己的菜谱，有权限
	if menu.UserId == userId {
		return true, nil
	}

	// 检查是否是恋人的菜谱
	loverId := m.getLoverId(userId)
	if loverId > 0 && menu.UserId == loverId {
		return true, nil
	}

	// 既不是自己的也不是恋人的，无权限
	return false, nil
}

// GetMenuWithCategories 获取菜谱及其分类和菜品（用于选择菜谱后获取完整数据，需要权限验证）
func (m *MenuService) GetMenuWithCategories(menuId string, userId uint) (menu app.Menu, categories []app.Category, foods []app.Food, err error) {
	// 获取菜谱信息
	err = global.GVA_DB.Where("id = ?", menuId).First(&menu).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜谱详情失败", zap.Error(err))
		return menu, categories, foods, err
	}

	// 检查权限
	hasPermission, err := m.checkMenuPermission(menu.ID, userId)
	if err != nil {
		return menu, categories, foods, err
	}
	if !hasPermission {
		global.GVA_LOG.Error("无权限查看该菜谱")
		return menu, categories, foods, errors.New("无权限查看该菜谱")
	}

	// 获取该菜谱下的所有分类
	err = global.GVA_DB.Where("menu_id = ?", menuId).Order("created_at ASC").Find(&categories).Error
	if err != nil {
		global.GVA_LOG.Error("获取分类列表失败", zap.Error(err))
		return menu, categories, foods, err
	}

	// 获取该菜谱下的所有菜品（通过分类关联）
	err = global.GVA_DB.Table("app_food").
		Joins("JOIN app_category ON app_food.category_id = app_category.id").
		Where("app_category.menu_id = ?", menuId).
		Order("app_food.created_at ASC").
		Find(&foods).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜品列表失败", zap.Error(err))
		return menu, categories, foods, err
	}

	return menu, categories, foods, nil
}
