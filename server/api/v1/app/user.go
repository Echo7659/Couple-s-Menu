package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"git.echol.cn/loser/menu-server/utils/user_jwt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// Login 小程序登录
func (u *UserApi) Login(ctx *gin.Context) {
	var p request.CodeLoginReq
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("登录请求参数错误："+err.Error(), ctx)
		return
	}

	token, err := userService.Login(ctx, p)
	if err != nil {
		global.GVA_LOG.Error("登录失败", zap.Error(err))
		r.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}
	r.OkWithData(token, ctx)
}

// GetInfo 获取用户信息
func (u *UserApi) GetInfo(ctx *gin.Context) {
	id := user_jwt.GetUserID(ctx)
	if id == 0 {
		r.FailWithMessage("用户ID获取失败", ctx)
		return
	}

	user, err := userService.GetUserByID(strconv.Itoa(int(id)))
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败：", zap.Error(err))
		r.FailWithMessage("获取用户信息失败："+err.Error(), ctx)
		return
	}
	r.OkWithData(user, ctx)
}

// GetInfoById 根据用户ID获取用户信息
func (u *UserApi) GetInfoById(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		r.FailWithMessage("用户ID不能为空", ctx)
		return
	}

	user, err := userService.GetUserByID(userId)
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败：", zap.Error(err))
		r.FailWithMessage("获取用户信息失败："+err.Error(), ctx)
		return
	}
	r.OkWithData(user, ctx)
}

// UpdateUserInfo 更新用户信息
func (u *UserApi) UpdateUserInfo(ctx *gin.Context) {
	var user app.User
	if err := ctx.ShouldBind(&user); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("更新用户信息请求参数错误："+err.Error(), ctx)
		return
	}
	user.ID = user_jwt.GetUserID(ctx)

	if user.ID == 0 {
		global.GVA_LOG.Error("用户ID不能为空")
		r.FailWithMessage("用户ID不能为空", ctx)
		return
	}

	err := userService.UpdateUserInfo(&user)
	if err != nil {
		global.GVA_LOG.Error("更新用户信息失败：", zap.Error(err))
		r.FailWithMessage("更新用户信息失败："+err.Error(), ctx)
		return
	}
	r.OkWithMessage("更新用户信息成功", ctx)
}

// 云存储配置
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	Domain          string
}

var ossConfig = OSSConfig{
	Endpoint:        "oss-cn-chengdu.aliyuncs.com",
	AccessKeyID:     "LTAI5tNSrjY28s3ax9F5RegJ",
	AccessKeySecret: "dnrgeAlMpMpePDkYODrfklMSNzvwDD",
	BucketName:      "menu-echo",
	Domain:          "https://menu-echo.oss-cn-chengdu.aliyuncs.com",
}

// UploadAvatarOSS 头像上传处理（云存储版本）
func (u *UserApi) UploadAvatarOSS(c *gin.Context) {
	// 1. 获取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "获取文件失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "文件大小不能超过5MB",
			"data": nil,
		})
		return
	}

	// 3. 创建临时文件
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "创建临时目录失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 生成临时文件路径
	tempFile := filepath.Join(tempDir, fmt.Sprintf("temp_%d_%s", time.Now().UnixNano(), file.Filename))

	// 4. 保存文件到临时位置
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "打开文件失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	defer src.Close()

	dst, err := os.Create(tempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "创建临时文件失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	defer dst.Close()
	defer os.Remove(tempFile) // 确保临时文件被删除

	// 5. 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "保存临时文件失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 6. 初始化OSS客户端
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyID, ossConfig.AccessKeySecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "初始化OSS客户端失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	bucket, err := client.Bucket(ossConfig.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取OSS存储桶失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 7. 生成OSS文件名
	ext := filepath.Ext(file.Filename)
	ossFileName := fmt.Sprintf("avatar/%d%s", time.Now().UnixNano(), ext)

	// 8. 上传到OSS
	if err := bucket.PutObjectFromFile(ossFileName, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "上传到OSS失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 9. 返回成功结果
	fileURL := fmt.Sprintf("%s/%s", ossConfig.Domain, ossFileName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "上传成功",
		"data": gin.H{
			"url":      fileURL,
			"filename": ossFileName,
			"size":     file.Size,
		},
	})
}

// RefreshToken 刷新Token
func (u *UserApi) RefreshToken(ctx *gin.Context) {
	token := user_jwt.GetToken(ctx)
	if token == "" {
		r.FailWithMessage("Token获取失败", ctx)
		return
	}

	j := user_jwt.NewUserJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("解析Token失败", zap.Error(err))
		r.FailWithMessage("Token解析失败："+err.Error(), ctx)
		return
	}

	var user app.User
	// 根据用户ID获取用户信息
	err = global.GVA_DB.Model(&app.User{}).Where("id = ?", claims.BaseClaims.ID).First(&user).Error
	if err != nil {
		r.FailWithMessage("获取用户信息失败："+err.Error(), ctx)
		return
	}
	// 生成新的Token
	newToken, _, err := user_jwt.LoginToken(user)
	if err != nil {
		global.GVA_LOG.Error("生成新Token失败", zap.Error(err))
		r.FailWithMessage("生成新Token失败："+err.Error(), ctx)
		return
	}

	r.OkWithData(newToken, ctx)
}

func (u *UserApi) UpdateUserRole(context *gin.Context) {
	var req app.User
	if err := context.ShouldBind(&req); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		r.FailWithMessage("更新用户角色请求参数错误："+err.Error(), context)
		return
	}

	// 用户ID从请求中获取
	req.ID = user_jwt.GetUserID(context)

	if req.ID == 0 {
		global.GVA_LOG.Error("获取用户ID失败")
		r.FailWithMessage("获取用户ID失败", context)
		return
	}

	err := userService.UpdateUserRole(&req)
	if err != nil {
		global.GVA_LOG.Error("更新用户角色失败：", zap.Error(err))
		r.FailWithMessage("更新用户角色失败："+err.Error(), context)
		return
	}
	r.OkWithMessage("更新用户角色成功", context)
}
