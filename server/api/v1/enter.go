package v1

import (
	"git.echol.cn/loser/menu-server/api/v1/app"
	"git.echol.cn/loser/menu-server/api/v1/example"
	"git.echol.cn/loser/menu-server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	AppApiGroup     app.ApiGroup
}
