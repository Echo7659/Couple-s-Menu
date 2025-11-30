package app

import (
	"strconv"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/repo"
	"git.echol.cn/loser/menu-server/model/app/request"
	"git.echol.cn/loser/menu-server/utils"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"git.echol.cn/loser/menu-server/utils/wechat"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserService struct{}

// Login 登录
func (u *UserService) Login(ctx *gin.Context, req request.CodeLoginReq) (vo repo.LoginVo, err error) {
	vo.IsFirstLogin = false
	session, err := wechat.MiniSDK.Auth.Session(ctx, req.Code)
	if err != nil {
		global.GVA_LOG.Error("登录失败", zap.Error(err))
		return
	}
	if session.ErrCode != 0 {
		global.GVA_LOG.Error("登录失败", zap.String("errMsg", session.ErrMsg))
		return
	}

	// 判断用户是否存在, 不存在则创建
	var count int64
	err = global.GVA_DB.Model(&app.User{}).Where("open_id = ?", session.OpenID).Count(&count).Error
	if err != nil {
		global.GVA_LOG.Error("查询用户失败", zap.Error(err))
		return
	}

	var user app.User

	if count == 0 {
		// 准备新用户基础信息，邀请码需要在写入后基于 ID 生成
		user.OpenId = session.OpenID
		user.NickName = "用户" + utils.RandomString(8)
		user.Avatar = "https://menu-echo.oss-cn-chengdu.aliyuncs.com/avatar/def.jpg"

		// 使用事务保障用户创建与邀请关系一致性
		tx := global.GVA_DB.Begin()
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			global.GVA_LOG.Error("创建用户失败", zap.Error(err))
			return vo, err
		}

		// 基于创建后的 ID 生成唯一邀请码并更新
		// 基于创建后的 ID 生成唯一邀请码并更新
		iniv := utils.RandomString(8) + strconv.Itoa(int(user.ID))
		user.InivCode = iniv
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			global.GVA_LOG.Error("更新邀请码失败", zap.Error(err))
			return vo, err
		}

		// 如果通过邀请注册，建立情侣关系和关联表
		if req.InivCode != "" {
			var inviter app.User
			if err := tx.Where("iniv_code = ?", req.InivCode).First(&inviter).Error; err != nil {
				tx.Rollback()
				global.GVA_LOG.Error("查询邀请人失败", zap.Error(err))
				return vo, err
			}

			// 建立情侣关系
			user.LoverId = inviter.ID
			inviter.LoverId = user.ID

			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback()
				global.GVA_LOG.Error("更新用户情侣关系失败", zap.Error(err))
				return vo, err
			}
			if err := tx.Save(&inviter).Error; err != nil {
				tx.Rollback()
				global.GVA_LOG.Error("更新邀请人情侣关系失败", zap.Error(err))
				return vo, err
			}

			// 创建关联表
			lovers := app.Lovers{
				UserID:  inviter.ID,
				LoverId: user.ID,
			}
			if err := tx.Create(&lovers).Error; err != nil {
				tx.Rollback()
				global.GVA_LOG.Error("创建情侣关联失败", zap.Error(err))
				return vo, err
			}
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			global.GVA_LOG.Error("提交事务失败", zap.Error(err))
			return vo, err
		}

		vo.IsFirstLogin = true
	} else {
		err = global.GVA_DB.Where("open_id = ?", session.OpenID).First(&user).Error
		if err != nil {
			global.GVA_LOG.Error("查询用户失败", zap.Error(err))
			return
		}
	}
	vo.Info = user

	// 生成JWT Token
	token, _, err := user_jwt.LoginToken(user)
	if err != nil {
		global.GVA_LOG.Error("生成token失败", zap.Error(err))
		return
	}
	vo.Token = token

	return
}

// GetUserByID 通过ID获取用户
func (u *UserService) GetUserByID(id string) (user app.User, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		global.GVA_LOG.Error("查询用户失败", zap.Error(err))
		return user, err
	}
	return user, nil
}

// UpdateUserInfo 更新用户信息
func (u *UserService) UpdateUserInfo(a *app.User) error {
	// 只更新nick_name和avatar
	err := global.GVA_DB.Model(&app.User{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"nick_name": a.NickName,
			"avatar":    a.Avatar,
			"role":      a.Role,
		}).Error

	if err != nil {
		global.GVA_LOG.Error("更新用户信息失败", zap.Error(err))
		return err
	}
	return nil
}

func (u *UserService) UpdateUserRole(a *app.User) error {

	// 只更新role
	err := global.GVA_DB.Model(&app.User{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"role": a.Role,
		}).Error

	if err != nil {
		global.GVA_LOG.Error("更新用户角色失败", zap.Error(err))
		return err
	}
	return nil

}
