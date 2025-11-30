package wechat

import (
	"git.echol.cn/loser/menu-server/global"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"go.uber.org/zap"
)

var MiniSDK *miniProgram.MiniProgram

// InitMiniSDK 初始化小程序SDK
func InitMiniSDK() {
	MiniProgramApp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     "", // 小程序appid
		Secret:    "", // 小程序app secret
		HttpDebug: true,
		Log: miniProgram.Log{
			Level: "debug",
			// 可以重定向到你的目录下，如果设置File和Error，默认会在当前目录下的wechat文件夹下生成日志
			File:   "/logs/info.log",
			Error:  "/logs/error.log",
			Stdout: false, //  是否打印在终端
		},
		// 可选，不传默认走程序内存
		//Cache: kernel.NewRedisClient(&kernel.UniversalOptions{
		//	Addrs:    []string{"127.0.0.1:6379"},
		//	Password: "",
		//	DB:       0,
		//}),
	})
	MiniSDK = MiniProgramApp

	if err != nil {
		global.GVA_LOG.Error("初始化小程序SDK失败", zap.Any("err", err))
	}
	return
}
