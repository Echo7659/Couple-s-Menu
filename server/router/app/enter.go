package app

import api "git.echol.cn/loser/menu-server/api/v1"

type RouterGroup struct {
	UserRouter
	FoodRouter
	AnniversaryRouter
	OrderRouter
	DiaryRouter
	MenuRouter
	TargetRouter
}

var (
	userApi                     = api.ApiGroupApp.AppApiGroup.UserApi
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
	loverApi                    = api.ApiGroupApp.AppApiGroup.LoversApi
	foodApi                     = api.ApiGroupApp.AppApiGroup.FoodApi
	anniversaryApi              = api.ApiGroupApp.AppApiGroup.AnniversaryApi
	orderApi                    = api.ApiGroupApp.AppApiGroup.OrderApi
	diaryApi                    = api.ApiGroupApp.AppApiGroup.DiaryApi
	menuApi                     = api.ApiGroupApp.AppApiGroup.MenuApi
	targetApi                   = api.ApiGroupApp.AppApiGroup.TargetApi
)
