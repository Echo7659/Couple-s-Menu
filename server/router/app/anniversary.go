package app

import "github.com/gin-gonic/gin"

type AnniversaryRouter struct{}

// InitAnniversaryRouter 纪念日路由
func (a *AnniversaryRouter) InitAnniversaryRouter(AppRouter *gin.RouterGroup) {
	jnrRouter := AppRouter.Group("anniversary")

	{
		jnrRouter.POST("", anniversaryApi.Create)                // 创建纪念日
		jnrRouter.DELETE("", anniversaryApi.Delete)              // 删除纪念日
		jnrRouter.PUT("", anniversaryApi.Update)                 // 更新纪念日
		jnrRouter.GET("list", anniversaryApi.GetAnniversaryList) // 获取纪念日
		jnrRouter.GET(":id", anniversaryApi.Get)                 // 根据ID获取纪

	}
}
