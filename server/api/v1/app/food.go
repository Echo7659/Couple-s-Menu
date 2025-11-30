package app

import (
	"strconv"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	common "git.echol.cn/loser/menu-server/model/common/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"github.com/gin-gonic/gin"
)

type FoodApi struct{}

// Create 创建菜品
func (f *FoodApi) Create(ctx *gin.Context) {
	var p app.Food
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.Create(p); err != nil {
		r.FailWithMessage("创建菜品失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("创建菜品成功", ctx)
}

// Update 更新菜品
func (f *FoodApi) Update(ctx *gin.Context) {
	var p app.Food
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.Update(p); err != nil {
		r.FailWithMessage("更新菜品失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("更新菜品成功", ctx)
}

// Delete 删除菜品
func (f *FoodApi) Delete(ctx *gin.Context) {
	var p app.Food
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.Delete(p.ID); err != nil {
		r.FailWithMessage("删除菜品失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("删除菜品成功", ctx)
}

// GetFoodList 获取菜品列表
func (f *FoodApi) GetFoodList(ctx *gin.Context) {
	var p request.GetFoodList
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	foodList, total, err := foodService.GetFoodList(p)
	if err != nil {
		r.FailWithMessage("获取菜品列表失败:"+err.Error(), ctx)
		return
	}
	r.OkWithDetailed(
		r.PageResult{
			List:     foodList,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		"获取菜品列表成功",
		ctx,
	)
}

// GetFoodById 获取菜品详情
func (f *FoodApi) GetFoodById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取纪念日详情参数错误: ID不能为空")
		r.FailWithMessage("获取纪念日详情参数错误: ID不能为空", ctx)
		return
	}

	food, err := foodService.GetFood(id)
	if err != nil {
		r.FailWithMessage("获取菜品详情失败:"+err.Error(), ctx)
		return
	}
	r.OkWithData(food, ctx)
}

// CreateCategory 创建菜品分类
func (f *FoodApi) CreateCategory(ctx *gin.Context) {
	var p app.Category
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.CreateCategory(p); err != nil {
		r.FailWithMessage("创建菜品分类失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("创建菜品分类成功", ctx)
}

// UpdateCategory 更新菜品分类
func (f *FoodApi) UpdateCategory(ctx *gin.Context) {
	var p app.Category
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.UpdateCategory(p); err != nil {
		r.FailWithMessage("更新菜品分类失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("更新菜品分类成功", ctx)
}

// DeleteCategory 删除菜品分类
func (f *FoodApi) DeleteCategory(ctx *gin.Context) {
	var p app.Category
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := foodService.DeleteCategory(p.ID); err != nil {
		r.FailWithMessage("删除菜品分类失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("删除菜品分类成功", ctx)
}

// GetCategoryList 获取菜品分类列表
func (f *FoodApi) GetCategoryList(ctx *gin.Context) {
	var p common.PageInfo
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	// 从查询参数中获取MenuId
	menuIdStr := ctx.Query("menuId")
	var menuId uint = 0
	if menuIdStr != "" {
		if parsedId, err := strconv.ParseUint(menuIdStr, 10, 32); err == nil {
			menuId = uint(parsedId)
		}
	}
	categoryList, total, err := foodService.GetCategoryList(menuId, p)
	if err != nil {
		r.FailWithMessage("获取菜品分类列表失败:"+err.Error(), ctx)
		return
	}
	r.OkWithDetailed(
		r.PageResult{
			List:     categoryList,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		"获取菜品分类列表成功",
		ctx,
	)
}

// GetCategoryById 获取菜品分类详情
func (f *FoodApi) GetCategoryById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取菜品分类失败: ID不能为空")
		r.FailWithMessage("获取菜品分类失败: ID不能为空", ctx)
		return
	}

	category, err := foodService.GetCategory(id)
	if err != nil {
		r.FailWithMessage("获取菜品分类详情失败:"+err.Error(), ctx)
		return
	}
	r.OkWithData(category, ctx)
}
