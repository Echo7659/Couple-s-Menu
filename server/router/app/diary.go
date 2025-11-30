package app

import "github.com/gin-gonic/gin"

type DiaryRouter struct{}

// InitDiaryRouter 初始化日记路由
func (r *DiaryRouter) InitDiaryRouter(AppRouter *gin.RouterGroup) {
	diaryRouter := AppRouter.Group("diary")
	{
		diaryRouter.POST("", diaryApi.Create)      // 新建日记
		diaryRouter.DELETE("", diaryApi.Delete)    // 删除日记
		diaryRouter.PUT("", diaryApi.Update)       // 更新日记
		diaryRouter.GET("list", diaryApi.GetList)  // 获取日记列表
		diaryRouter.GET(":id", diaryApi.GetDetail) // 根据ID获取日记
	}
}
