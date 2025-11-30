package router

import (
	"git.echol.cn/loser/menu-server/router/app"
	"git.echol.cn/loser/menu-server/router/example"
	"git.echol.cn/loser/menu-server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	App     app.RouterGroup
}
