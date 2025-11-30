package app

import (
	"git.echol.cn/loser/menu-server/global"
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/repo"
	"git.echol.cn/loser/menu-server/model/app/request"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type OrderService struct {
}

func (s OrderService) Create(a *app.Order) (*app.Order, error) {
	a.Date = time.Now().Format("2006-01-02")
	err := global.GVA_DB.Create(&a).Error
	if err != nil {
		global.GVA_LOG.Error("创建订单失败", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (s OrderService) GetOrderList(p request.GetOrderList) (list []repo.OrderVo, err error) {
	db := global.GVA_DB.Model(&app.Order{})

	db.Where("user_id = ?", p.UserId)
	if p.Date != "" {
		db.Where("date = ?", p.Date)
	}

	var orders []app.Order
	err = db.Find(&orders).Error
	if err != nil {
		global.GVA_LOG.Error("查询订单列表失败", zap.Error(err))
		return nil, err
	}

	for _, order := range orders {
		var foods []app.Food

		if order.FoodIds != "" {
			foodIds := strings.Split(order.FoodIds, ",")

			err = global.GVA_DB.Model(&app.Food{}).Where("id IN (?)", foodIds).Find(&foods).Error
			if err != nil {
				global.GVA_LOG.Error("查询菜品详情失败", zap.Error(err))
				return nil, err
			}
		} else {
			global.GVA_LOG.Warn("订单没有菜品信息", zap.String("order_id", strconv.Itoa(int(order.ID))))
		}
		// 查询菜品详情
		temp := repo.OrderVo{
			ID:     order.ID,
			Remark: order.Remark,
			Date:   order.Date,
			Foods:  foods,
		}
		list = append(list, temp)
	}

	return
}

func (s OrderService) GetOrderDetail(id string) (order repo.OrderVo, err error) {
	var tempOrder app.Order
	err = global.GVA_DB.Model(&app.Order{}).Where("id = ?", id).First(&tempOrder).Error
	if err != nil {
		global.GVA_LOG.Error("查询订单详情失败", zap.Error(err))
		return repo.OrderVo{}, err
	}

	err = global.GVA_DB.Model(&app.Food{}).Where("id IN (?)", tempOrder.FoodIds).Find(&order.Foods).Error
	if err != nil {
		global.GVA_LOG.Error("查询菜品详情失败", zap.Error(err))
		return repo.OrderVo{}, err
	}
	order.ID = tempOrder.ID
	order.Remark = tempOrder.Remark
	order.Date = tempOrder.Date
	return
}

func (s OrderService) DeleteOrder(p app.Order) error {
	err := global.GVA_DB.Where("id = ?", p.ID).Delete(&app.Order{}).Error
	if err != nil {
		global.GVA_LOG.Error("删除订单失败", zap.Error(err))
		return err
	}
	return nil
}
