package app

import "github.com/gin-gonic/gin"

type TargetRouter struct{}

// InitTargetRouter 打卡目标路由初始化
func (r *TargetRouter) InitTargetRouter(AppRouter *gin.RouterGroup) {
	targetRouter := AppRouter.Group("target")
	{
		targetRouter.GET("list", targetApi.GetTargetList)           // 获取打卡目标列表
		targetRouter.GET("my", targetApi.GetMyTargets)              // 获取我的目标（包括情侣共享的）
		targetRouter.GET(":id", targetApi.GetTargetById)            // 根据ID获取打卡目标
		targetRouter.GET(":id/records", targetApi.GetTargetRecords) // 获取目标的打卡记录
		targetRouter.POST("", targetApi.Create)                     // 新建打卡目标
		targetRouter.PUT("", targetApi.Update)                      // 更新打卡目标
		targetRouter.DELETE("", targetApi.Delete)                   // 删除打卡目标
		targetRouter.POST("record", targetApi.CreateRecord)         // 创建打卡记录
		targetRouter.POST("complete", targetApi.CompleteTarget)     // 完成打卡目标（标记为已完成）
	}
}
