package app

import (
	"strconv"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	common "git.echol.cn/loser/menu-server/model/common/request"
	"git.echol.cn/loser/menu-server/model/common/response"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DiaryApi struct{}

// Create 新建日记
func (d *DiaryApi) Create(ctx *gin.Context) {
	var p app.Diary
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("新建日记参数错误", zap.Error(err))
		response.FailWithMessage("新建日记参数错误", ctx)
		return
	}

	if err := diaryService.Create(p); err != nil {
		global.GVA_LOG.Error("新建日记失败", zap.Error(err))
		response.FailWithMessage("新建日记失败:"+err.Error(), ctx)
		return
	}

	response.OkWithMessage("新建日记成功", ctx)
}

// GetList 获取日记列表
func (d *DiaryApi) GetList(ctx *gin.Context) {
	var p common.PageInfo
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("获取日记列表参数错误", zap.Error(err))
		response.FailWithMessage("获取日记列表参数错误", ctx)
		return
	}

	if p.UserId == 0 {
		p.UserId = user_jwt.GetUserID(ctx)
	}

	diaryList, total, err := diaryService.GetDiaryList(p)
	if err != nil {
		global.GVA_LOG.Error("获取日记列表失败", zap.Error(err))
		response.FailWithMessage("获取日记列表失败:"+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     diaryList,
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
	}, "获取日记列表成功", ctx)
}

// Update 更新日记
func (d *DiaryApi) Update(ctx *gin.Context) {
	var p app.Diary
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("更新日记参数错误", zap.Error(err))
		response.FailWithMessage("更新日记参数错误", ctx)
		return
	}

	if err := diaryService.Update(p); err != nil {
		global.GVA_LOG.Error("更新日记失败", zap.Error(err))
		response.FailWithMessage("更新日记失败:"+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新日记成功", ctx)
}

// Delete 删除日记
func (d *DiaryApi) Delete(ctx *gin.Context) {
	var p common.GetById
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("删除日记参数错误", zap.Error(err))
		response.FailWithMessage("删除日记参数错误", ctx)
		return
	}

	if err := diaryService.Delete(p.ID); err != nil {
		global.GVA_LOG.Error("删除日记失败", zap.Error(err))
		response.FailWithMessage("删除日记失败:"+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除日记成功", ctx)
}

// GetDetail 获取日记详情
func (d *DiaryApi) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		global.GVA_LOG.Error("获取日记详情参数错误", zap.Error(err))
		response.FailWithMessage("获取日记详情参数错误", ctx)
		return
	}

	diary, err := diaryService.GetDiary(uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取日记详情失败", zap.Error(err))
		response.FailWithMessage("获取日记详情失败:"+err.Error(), ctx)
		return
	}

	response.OkWithDetailed(diary, "获取日记详情成功", ctx)
}
