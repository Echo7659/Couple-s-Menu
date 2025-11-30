package service

import (
	"git.echol.cn/loser/menu-server/service/app"
	"git.echol.cn/loser/menu-server/service/example"
	"git.echol.cn/loser/menu-server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	AppServiceGroup     app.ServiceGroup
}
