package app

import (
	"errors"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	common "git.echol.cn/loser/menu-server/model/common/request"
	"go.uber.org/zap"
)

type FoodService struct{}

// Create 创建菜品
func (f *FoodService) Create(req app.Food) (err error) {
	err = global.GVA_DB.Create(&req).Error
	if err != nil {
		global.GVA_LOG.Error("创建菜品失败", zap.Error(err))
		return err
	}
	return
}

// Delete 删除菜品
func (f *FoodService) Delete(id uint) (err error) {
	err = global.GVA_DB.Where("id = ?", id).Delete(&app.Food{}).Error
	if err != nil {
		global.GVA_LOG.Error("删除菜品失败", zap.Error(err))
		return err
	}
	return
}

// Update 更新菜品
func (f *FoodService) Update(req app.Food) (err error) {
	err = global.GVA_DB.Model(&app.Food{}).Where("id = ?", req.ID).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新菜品失败", zap.Error(err))
		return err
	}
	return
}

// GetFood 获取菜品详情
func (f *FoodService) GetFood(id string) (food app.Food, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&food).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜品详情失败", zap.Error(err))
		return food, err
	}
	return food, nil
}

// GetFoodList 获取菜品列表
func (f *FoodService) GetFoodList(pageInfo request.GetFoodList) (list []app.Food, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&app.Food{})

	// 如果指定了MenuId，需要通过Category关联查询
	if pageInfo.MenuId != 0 {
		db = db.Joins("JOIN app_category ON app_food.category_id = app_category.id").
			Where("app_category.menu_id = ?", pageInfo.MenuId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageInfo.Keyword != "" {
		db = db.Where("app_food.name LIKE ? ", "%"+pageInfo.Keyword+"%")
	}
	if pageInfo.CategoryId != 0 {
		db = db.Where("app_food.category_id = ?", pageInfo.CategoryId)
	}
	if pageInfo.Name != "" {
		db = db.Where("app_food.name LIKE ?", "%"+pageInfo.Name+"%")
	}

	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜品列表失败", zap.Error(err))
		return list, total, err
	}
	return list, total, nil
}

// CreateCategory 创建菜品分类
func (f *FoodService) CreateCategory(req app.Category) (err error) {
	// 验证MenuId是否存在
	if req.MenuId == 0 {
		global.GVA_LOG.Error("创建菜品分类失败: MenuId不能为空")
		return errors.New("MenuId不能为空")
	}
	err = global.GVA_DB.Create(&req).Error
	if err != nil {
		global.GVA_LOG.Error("创建菜品分类失败", zap.Error(err))
		return err
	}
	return
}

// GetCategoryList 获取菜品分类列表
func (f *FoodService) GetCategoryList(menuId uint, pageInfo common.PageInfo) (list []app.Category, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&app.Category{})

	// 如果指定了MenuId，则按MenuId过滤
	if menuId != 0 {
		db = db.Where("menu_id = ?", menuId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageInfo.Keyword != "" {
		db = db.Where("name LIKE ? ", "%"+pageInfo.Keyword+"%")
	}

	err = db.Limit(limit).Offset(offset).Order("created_at ASC").Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜品分类列表失败", zap.Error(err))
		return list, total, err
	}
	return list, total, nil
}

// UpdateCategory 更新菜品分类
func (f *FoodService) UpdateCategory(p app.Category) (err error) {
	err = global.GVA_DB.Model(&app.Category{}).Where("id = ?", p.ID).Updates(&p).Error
	if err != nil {
		global.GVA_LOG.Error("更新菜品分类失败", zap.Error(err))
		return
	}
	return
}

func (f *FoodService) DeleteCategory(id uint) (err error) {
	err = global.GVA_DB.Where("id = ?", id).Delete(&app.Category{}).Error
	if err != nil {
		global.GVA_LOG.Error("删除菜品分类失败", zap.Error(err))
		return err
	}
	return
}

func (f *FoodService) GetCategory(id string) (category app.Category, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		global.GVA_LOG.Error("获取菜品分类详情失败", zap.Error(err))
		return category, err
	}
	return category, nil
}
