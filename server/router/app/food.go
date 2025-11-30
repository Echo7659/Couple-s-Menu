package app

import "github.com/gin-gonic/gin"

type FoodRouter struct{}

// InitFoodRouter 菜单路由初始化
func (r *FoodRouter) InitFoodRouter(AppRouter *gin.RouterGroup) {
	foodRouter := AppRouter.Group("food")
	{
		foodRouter.GET("list", foodApi.GetFoodList) // 获取菜单列表
		foodRouter.GET(":id", foodApi.GetFoodById)  // 根据ID获取菜单
		foodRouter.POST("", foodApi.Create)         // 新建菜品
		foodRouter.PUT("", foodApi.Update)          // 更新菜品
		foodRouter.DELETE("", foodApi.Delete)       // 删除菜品
	}

	{
		foodRouter.GET("category/list", foodApi.GetCategoryList) // 获取分类列表
		foodRouter.GET("category/:id", foodApi.GetCategoryById)  // 根据ID获取分类
		foodRouter.POST("category", foodApi.CreateCategory)      // 新建分类
		foodRouter.PUT("category", foodApi.UpdateCategory)       // 更新分类
		foodRouter.DELETE("category", foodApi.DeleteCategory)    // 删除分类
	}

}
