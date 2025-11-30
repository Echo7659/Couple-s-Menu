package app

import (
	"errors"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoversService struct{}

func (s LoversService) GetUserLovers(id uint) (lover app.Lovers, err error) {
	err = global.GVA_DB.Where("user_id = ? OR lover_id = ?", id, id).First(&lover).Error
	// 判断是否为空 无数据
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return app.Lovers{}, errors.New("用户暂未绑定恋人")
	}

	if err != nil {
		global.GVA_LOG.Error("查询用户恋人信息失败", zap.Error(err))
		return
	}
	return
}

func (s LoversService) BindLover(context *gin.Context, p request.BindLoverReq) error {
	// 获取当前用户ID
	userId := user_jwt.GetUserID(context)
	if userId == 0 {
		return errors.New("用户ID获取失败")
	}

	// 查询当前用户信息
	var user app.User
	if err := global.GVA_DB.Where("id = ?", userId).First(&user).Error; err != nil {
		global.GVA_LOG.Error("查询用户信息失败", zap.Error(err))
		return errors.New("查询用户信息失败")
	}

	// 查询要绑定的恋人信息
	var lover app.User
	if err := global.GVA_DB.Where("iniv_code = ?", p.InivCode).First(&lover).Error; err != nil {
		global.GVA_LOG.Error("查询恋人信息失败", zap.Error(err))
		return errors.New("查询恋人信息失败")
	}

	// 检查是否已经绑定
	if user.LoverId != 0 || lover.LoverId != 0 {
		return errors.New("您或您的恋人已绑定其他用户，无法重复绑定")
	}

	// 创建情侣关系
	lovers := app.Lovers{
		UserID:  user.ID,
		LoverId: lover.ID,
	}
	if err := global.GVA_DB.Create(&lovers).Error; err != nil {
		global.GVA_LOG.Error("创建情侣关系失败", zap.Error(err))
		return errors.New("创建情侣关系失败")
	}
	// 更新用户的LoverId字段
	if err := global.GVA_DB.Model(&user).Update("lover_id", lover.ID).Error; err != nil {
		global.GVA_LOG.Error("更新用户恋人ID失败", zap.Error(err))
		return errors.New("更新用户恋人ID失败")
	}
	if err := global.GVA_DB.Model(&lover).Update("lover_id", user.ID).Error; err != nil {
		global.GVA_LOG.Error("更新恋人用户ID失败", zap.Error(err))
		return errors.New("更新恋人用户ID失败")
	}

	return nil
}

func (s LoversService) UpdateLoverDay(p app.Lovers) error {

	err := global.GVA_DB.Model(&app.Lovers{}).Where("id = ?", p.ID).Update("start_date", p.StartDate).Error
	if err != nil {
		global.GVA_LOG.Error("更新恋人纪念日失败", zap.Error(err))
		return err
	}
	return nil

}
