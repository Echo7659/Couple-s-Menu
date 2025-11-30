package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	common "git.echol.cn/loser/menu-server/model/common/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MenuApi struct{}

// Create 创建菜谱
func (m *MenuApi) Create(ctx *gin.Context) {
	var p app.Menu
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	// 从JWT中获取用户ID
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}
	p.UserId = userId

	if err := menuService.Create(p); err != nil {
		r.FailWithMessage("创建菜谱失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("创建菜谱成功", ctx)
}

// Update 更新菜谱
func (m *MenuApi) Update(ctx *gin.Context) {
	var p app.Menu
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	// 验证用户是否有权限修改该菜谱
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}
	// 检查菜谱是否属于当前用户
	var existingMenu app.Menu
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", p.ID, userId).First(&existingMenu).Error; err != nil {
		global.GVA_LOG.Error("菜谱不存在或无权限", zap.Error(err))
		r.FailWithMessage("菜谱不存在或无权限", ctx)
		return
	}

	if err := menuService.Update(p); err != nil {
		r.FailWithMessage("更新菜谱失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("更新菜谱成功", ctx)
}

// Delete 删除菜谱
func (m *MenuApi) Delete(ctx *gin.Context) {
	var p app.Menu
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	// 验证用户是否有权限删除该菜谱
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}
	// 检查菜谱是否属于当前用户
	var existingMenu app.Menu
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", p.ID, userId).First(&existingMenu).Error; err != nil {
		global.GVA_LOG.Error("菜谱不存在或无权限", zap.Error(err))
		r.FailWithMessage("菜谱不存在或无权限", ctx)
		return
	}

	if err := menuService.Delete(p.ID); err != nil {
		r.FailWithMessage("删除菜谱失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("删除菜谱成功", ctx)
}

// GetMenuList 获取菜谱列表（当前用户的）
func (m *MenuApi) GetMenuList(ctx *gin.Context) {
	var p common.PageInfo
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	// 从JWT中获取用户ID
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}

	menuList, total, err := menuService.GetMenuList(userId, p)
	if err != nil {
		r.FailWithMessage("获取菜谱列表失败:"+err.Error(), ctx)
		return
	}
	r.OkWithDetailed(
		r.PageResult{
			List:     menuList,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		"获取菜谱列表成功",
		ctx,
	)
}

// GetMenuById 获取菜谱详情
func (m *MenuApi) GetMenuById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取菜谱详情参数错误: ID不能为空")
		r.FailWithMessage("获取菜谱详情参数错误: ID不能为空", ctx)
		return
	}

	// 从JWT中获取用户ID
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}

	menu, err := menuService.GetMenu(id, userId)
	if err != nil {
		r.FailWithMessage("获取菜谱详情失败:"+err.Error(), ctx)
		return
	}
	r.OkWithData(menu, ctx)
}

// GetMenuWithCategories 获取菜谱及其分类和菜品（选择菜谱后调用）
func (m *MenuApi) GetMenuWithCategories(ctx *gin.Context) {
	menuId := ctx.Param("id")
	if menuId == "" {
		global.GVA_LOG.Error("获取菜谱详情参数错误: ID不能为空")
		r.FailWithMessage("获取菜谱详情参数错误: ID不能为空", ctx)
		return
	}

	// 从JWT中获取用户ID
	userId := user_jwt.GetUserID(ctx)
	if userId == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}

	menu, categories, foods, err := menuService.GetMenuWithCategories(menuId, userId)
	if err != nil {
		r.FailWithMessage("获取菜谱及分类失败:"+err.Error(), ctx)
		return
	}

	type MenuWithCategoriesResponse struct {
		Menu       app.Menu       `json:"menu"`
		Categories []app.Category `json:"categories"`
		Foods      []app.Food     `json:"foods"`
	}

	r.OkWithData(MenuWithCategoriesResponse{
		Menu:       menu,
		Categories: categories,
		Foods:      foods,
	}, ctx)
}
