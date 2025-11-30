package app

import (
	"strconv"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"github.com/gin-gonic/gin"
)

type TargetApi struct{}

// Create 创建打卡目标
func (t *TargetApi) Create(ctx *gin.Context) {
	var p app.Target
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := targetService.Create(p); err != nil {
		r.FailWithMessage("创建打卡目标失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("创建打卡目标成功", ctx)
}

// Update 更新打卡目标
func (t *TargetApi) Update(ctx *gin.Context) {
	var p app.Target
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := targetService.Update(p); err != nil {
		r.FailWithMessage("更新打卡目标失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("更新打卡目标成功", ctx)
}

// Delete 删除打卡目标
func (t *TargetApi) Delete(ctx *gin.Context) {
	var p app.Target
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := targetService.Delete(p.ID); err != nil {
		r.FailWithMessage("删除打卡目标失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("删除打卡目标成功", ctx)
}

// GetTargetList 获取打卡目标列表
func (t *TargetApi) GetTargetList(ctx *gin.Context) {
	var p request.GetTargetList
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	targetList, total, err := targetService.GetTargetList(p)
	if err != nil {
		r.FailWithMessage("获取打卡目标列表失败:"+err.Error(), ctx)
		return
	}
	r.OkWithDetailed(
		r.PageResult{
			List:     targetList,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		"获取打卡目标列表成功",
		ctx,
	)
}

// GetMyTargets 获取我的目标（包括我创建的和情侣共享的）
func (t *TargetApi) GetMyTargets(ctx *gin.Context) {
	var p request.GetTargetList
	if err := ctx.ShouldBind(&p); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}

	// 从上下文获取用户ID（需要从JWT中获取）
	userIdStr := ctx.Query("userId")
	if userIdStr == "" {
		r.FailWithMessage("缺少用户ID参数", ctx)
		return
	}
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		r.FailWithMessage("用户ID格式错误", ctx)
		return
	}

	targetList, total, err := targetService.GetMyTargets(uint(userId), p)
	if err != nil {
		r.FailWithMessage("获取我的目标列表失败:"+err.Error(), ctx)
		return
	}
	r.OkWithDetailed(
		r.PageResult{
			List:     targetList,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		"获取我的目标列表成功",
		ctx,
	)
}

// GetTargetById 获取打卡目标详情
func (t *TargetApi) GetTargetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取打卡目标详情参数错误: ID不能为空")
		r.FailWithMessage("获取打卡目标详情参数错误: ID不能为空", ctx)
		return
	}

	target, err := targetService.GetTarget(id)
	if err != nil {
		r.FailWithMessage("获取打卡目标详情失败:"+err.Error(), ctx)
		return
	}
	r.OkWithData(target, ctx)
}

// CreateRecord 创建打卡记录
func (t *TargetApi) CreateRecord(ctx *gin.Context) {
	var req request.CreateTargetRecord
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := targetService.CreateRecord(req); err != nil {
		r.FailWithMessage("创建打卡记录失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("打卡成功", ctx)
}

// GetTargetRecords 获取目标的所有打卡记录
func (t *TargetApi) GetTargetRecords(ctx *gin.Context) {
	targetIdStr := ctx.Param("id")
	if targetIdStr == "" {
		r.FailWithMessage("目标ID不能为空", ctx)
		return
	}

	targetId, err := strconv.ParseUint(targetIdStr, 10, 32)
	if err != nil {
		r.FailWithMessage("目标ID格式错误", ctx)
		return
	}

	records, err := targetService.GetTargetRecords(uint(targetId))
	if err != nil {
		r.FailWithMessage("获取打卡记录失败:"+err.Error(), ctx)
		return
	}
	r.OkWithData(records, ctx)
}

// CompleteTarget 完成打卡目标
func (t *TargetApi) CompleteTarget(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.FailWithMessage("参数错误", ctx)
		return
	}
	if err := targetService.CompleteTarget(req.ID); err != nil {
		r.FailWithMessage("完成打卡目标失败:"+err.Error(), ctx)
		return
	}
	r.OkWithMessage("完成打卡目标成功", ctx)
}
