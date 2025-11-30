package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoversApi struct{}

// GetUserLovers 获取用户的情侣信息
func (l *LoversApi) GetUserLovers(c *gin.Context) {
	userId := user_jwt.GetUserID(c)
	if userId == 0 {
		r.FailWithMessage("用户ID为空", c)
		return
	}

	lover, err := loveService.GetUserLovers(userId)
	if err != nil {
		r.FailWithMessage("获取情侣信息失败: "+err.Error(), c)
		return
	}

	r.OkWithData(lover, c)
}

func (l *LoversApi) BindLover(context *gin.Context) {
	var p request.BindLoverReq
	if err := context.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误: "+err.Error(), zap.Error(err))
		r.FailWithMessage("参数错误", context)
		return
	}

	if err := loveService.BindLover(context, p); err != nil {
		r.FailWithMessage("绑定情侣失败: "+err.Error(), context)
		return
	}

	r.OkWithMessage("绑定情侣成功", context)

}

func (l *LoversApi) UpdateLoverDay(context *gin.Context) {
	var p app.Lovers
	if err := context.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误: "+err.Error(), zap.Error(err))
		r.FailWithMessage("参数错误", context)
		return
	}

	if err := loveService.UpdateLoverDay(p); err != nil {
		r.FailWithMessage("更新纪念日失败: "+err.Error(), context)
		return
	}

	r.OkWithMessage("更新纪念日成功", context)
}
