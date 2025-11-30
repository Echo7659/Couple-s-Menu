package response

import (
	"git.echol.cn/loser/menu-server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
