package request

import common "git.echol.cn/loser/menu-server/model/common/request"

type CodeLoginReq struct {
	Code     string `json:"code" form:"code" binding:"required"`
	InivCode string `json:"inivCode" form:"inivCode"`
}
type BindLoverReq struct {
	InivCode string `json:"inivCode" form:"inivCode" binding:"required"`
}

type UploadBase64AvatarReq struct {
	Base64 string `json:"base64" form:"base64" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type GetJNRList struct {
	common.PageInfo
	UserId uint `json:"userId" form:"userId"` // 用户ID
}
