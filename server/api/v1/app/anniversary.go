package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AnniversaryApi 纪念日接口
type AnniversaryApi struct {
}

// GetAnniversaryList 获取纪念日列表
func (a *AnniversaryApi) GetAnniversaryList(ctx *gin.Context) {
	var p request.GetJNRList
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("参数错误", ctx)
		return
	}
	list, total, err := annviversaryService.GetAnniversaryList(p)
	if err != nil {
		r.FailWithMessage("获取纪念日列表失败: "+err.Error(), ctx)
		return
	}

	r.OkWithDetailed(
		r.PageResult{
			List:     list,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		}, "获取纪念日列表成功", ctx)
}

// Create 创建纪念日
func (a *AnniversaryApi) Create(ctx *gin.Context) {
	var p app.Anniversary
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("参数错误", ctx)
		return
	}

	if err := annviversaryService.Create(p); err != nil {
		r.FailWithMessage("创建纪念日失败: "+err.Error(), ctx)
		return
	}

	r.OkWithMessage("创建纪念日成功", ctx)
}

// Delete 删除纪念日
func (a *AnniversaryApi) Delete(ctx *gin.Context) {
	var p app.Anniversary
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("参数错误", ctx)
		return
	}

	if err := annviversaryService.Delete(p.ID); err != nil {
		r.FailWithMessage("删除纪念日失败: "+err.Error(), ctx)
		return
	}

	r.OkWithMessage("删除纪念日成功", ctx)
}

// Update 更新纪念日
func (a *AnniversaryApi) Update(ctx *gin.Context) {
	var p app.Anniversary
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("参数错误", ctx)
		return
	}

	if err := annviversaryService.Update(p); err != nil {
		r.FailWithMessage("更新纪念日失败: "+err.Error(), ctx)
		return
	}

	r.OkWithMessage("更新纪念日成功", ctx)
}

// Get 获取纪念日详情
func (a *AnniversaryApi) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取纪念日详情参数错误: ID不能为空")
		r.FailWithMessage("获取纪念日详情参数错误: ID不能为空", ctx)
		return
	}

	anniversary, err := annviversaryService.GetAnniversary(id)
	if err != nil {
		r.FailWithMessage("获取纪念日详情失败: "+err.Error(), ctx)
		return
	}

	r.OkWithData(anniversary, ctx)
}
