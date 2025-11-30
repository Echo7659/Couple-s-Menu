package app

import (
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (e *UserRouter) InitCustomerRouter(AppRouter *gin.RouterGroup, PublicGroup *gin.RouterGroup) {
	userRouter := AppRouter.Group("user")
	publicUserRouter := PublicGroup.Group("user")
	{
		publicUserRouter.POST("login", userApi.Login)
		userRouter.GET("info", userApi.GetInfo)
		userRouter.GET("info/:id", userApi.GetInfoById)
		userRouter.PUT("info", userApi.UpdateUserInfo)
		userRouter.PUT("role", userApi.UpdateUserRole)
		userRouter.POST("upload", exaFileUploadAndDownloadApi.UploadFile)
		userRouter.POST("avatar", userApi.UploadAvatarOSS)          // 上传头像
		publicUserRouter.POST("refreshToken", userApi.RefreshToken) // 刷新token
	}
	{
		// 恋人相关接口
		userRouter.GET("lover", loverApi.GetUserLovers)       // 查询恋人
		userRouter.POST("/lover/bind", loverApi.BindLover)    // 绑定恋人
		userRouter.PUT("/lover/day", loverApi.UpdateLoverDay) // 更新恋爱纪念日
	}

}
