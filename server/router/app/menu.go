package app

import "github.com/gin-gonic/gin"

type MenuRouter struct{}

// InitMenuRouter 菜谱路由初始化
func (r *MenuRouter) InitMenuRouter(AppRouter *gin.RouterGroup) {
	menuRouter := AppRouter.Group("menu")
	{
		menuRouter.GET("list", menuApi.GetMenuList)                     // 获取菜谱列表
		menuRouter.GET(":id", menuApi.GetMenuById)                      // 根据ID获取菜谱
		menuRouter.GET(":id/categories", menuApi.GetMenuWithCategories) // 获取菜谱及其分类
		menuRouter.POST("", menuApi.Create)                             // 新建菜谱
		menuRouter.PUT("", menuApi.Update)                              // 更新菜谱
		menuRouter.DELETE("", menuApi.Delete)                           // 删除菜谱
	}
}
