package request

import (
	"git.echol.cn/loser/menu-server/model/common/request"
	"git.echol.cn/loser/menu-server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
