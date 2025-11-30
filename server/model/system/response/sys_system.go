package response

import "git.echol.cn/loser/menu-server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
