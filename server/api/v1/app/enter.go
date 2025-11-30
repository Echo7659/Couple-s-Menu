package app

import "git.echol.cn/loser/menu-server/service"

type ApiGroup struct {
	UserApi
	LoversApi
	FoodApi
	AnniversaryApi
	OrderApi
	DiaryApi
	MenuApi
	TargetApi
}

var (
	userService         = service.ServiceGroupApp.AppServiceGroup.UserService
	loveService         = service.ServiceGroupApp.AppServiceGroup.LoversService
	foodService         = service.ServiceGroupApp.AppServiceGroup.FoodService
	annviversaryService = service.ServiceGroupApp.AppServiceGroup.AnniversaryService
	orderService        = service.ServiceGroupApp.AppServiceGroup.OrderService
	diaryService        = service.ServiceGroupApp.AppServiceGroup.DiaryService
	menuService         = service.ServiceGroupApp.AppServiceGroup.MenuService
	targetService       = service.ServiceGroupApp.AppServiceGroup.TargetService
)
