package app

import (
	"strconv"
	"strings"
	"time"

	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	orquest "git.echol.cn/loser/menu-server/model/app/request"
	r "git.echol.cn/loser/menu-server/model/common/response"
	"git.echol.cn/loser/menu-server/utils/wechat"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/basicService/subscribeMessage/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderApi struct{}

// Create 创建订单
func (o *OrderApi) Create(ctx *gin.Context) {
	var p app.Order
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("创建订单参数错误", zap.Error(err))
		r.FailWithMessage("创建订单参数错误", ctx)
		return
	}

	order, err := orderService.Create(&p)
	if err != nil {
		global.GVA_LOG.Error("创建订单失败", zap.Error(err))
		r.FailWithMessage("创建订单失败", ctx)
		return
	}

	// 查询当前用户信息
	go func() {
		user, err := userService.GetUserByID(strconv.Itoa(int(p.UserId)))
		if err != nil {
			global.GVA_LOG.Error("获取用户信息失败", zap.Error(err))
			return
		}
		// 判断用户是否为吃货
		if user.Role == 2 {
			if user.LoverId == 0 {
				global.GVA_LOG.Error("用户未绑定饲养员，无法发送订阅消息通知")
				return
			}
			// 用户为吃货，通知饲养员
			lover, err := userService.GetUserByID(strconv.Itoa(int(user.LoverId)))
			if err != nil {
				global.GVA_LOG.Error("获取饲养员信息失败", zap.Error(err))
				return
			}

			// 查询菜品名称
			foodIds := strings.Split(order.FoodIds, ",")
			var foodNameList []string
			global.GVA_DB.Model(&app.Food{}).Select("name").Where("id IN (?)", foodIds).Scan(&foodNameList)
			foodName := strings.Join(foodNameList, ", ")

			data := &power.HashMap{
				"time1": power.StringMap{
					"value": time.Now().Format("2006-01-02 15:04:05"),
				},
				"thing4": power.StringMap{
					"value": user.NickName,
				},
				"thing5": power.StringMap{
					"value": foodName,
				},
			}

			wechat.MiniSDK.SubscribeMessage.Send(ctx, &request.RequestSubscribeMessageSend{
				ToUser:           lover.OpenId,
				TemplateID:       "qwh5Cymj3VYvIfPMjI6zw-h-Pstg1saDBlbDysbODYU",
				Page:             "/pages/home/home",
				MiniProgramState: "developer",
				Lang:             "zh_CN",
				Data:             data,
			})
		}

	}()

	r.OkWithData(order, ctx)
}

// GetOrderList 获取订单列表
func (o *OrderApi) GetOrderList(ctx *gin.Context) {
	var p orquest.GetOrderList
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("获取订单列表参数错误", zap.Error(err))
		r.FailWithMessage("获取订单列表参数错误", ctx)
		return
	}

	orderList, err := orderService.GetOrderList(p)
	if err != nil {
		global.GVA_LOG.Error("获取订单列表失败", zap.Error(err))
		r.FailWithMessage("获取订单列表失败", ctx)
		return
	}

	r.OkWithDetailed(orderList, "获取订单列表成功", ctx)
}

// GetOrderDetail 获取订单详情
func (o *OrderApi) GetOrderDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		global.GVA_LOG.Error("获取订单详情参数错误: ID不能为空")
		r.FailWithMessage("获取订单详情参数错误: ID不能为空", ctx)
		return
	}

	order, err := orderService.GetOrderDetail(id)
	if err != nil {
		global.GVA_LOG.Error("获取订单详情失败", zap.Error(err))
		r.FailWithMessage("获取订单详情失败", ctx)
		return
	}

	r.OkWithData(order, ctx)
}

// Delete 删除订单
func (o *OrderApi) Delete(ctx *gin.Context) {
	var p app.Order
	if err := ctx.ShouldBind(&p); err != nil {
		global.GVA_LOG.Error("删除订单参数错误", zap.Error(err))
		r.FailWithMessage("删除订单参数错误", ctx)
		return
	}

	if err := orderService.DeleteOrder(p); err != nil {
		global.GVA_LOG.Error("删除订单失败", zap.Error(err))
		r.FailWithMessage("删除订单失败", ctx)
		return
	}

	r.OkWithMessage("删除订单成功", ctx)
}
