package app

import "github.com/gin-gonic/gin"

type OrderRouter struct{}

// InitOrderRouter initializes the order routes
func (r *OrderRouter) InitOrderRouter(AppRouter *gin.RouterGroup) {
	orderRouter := AppRouter.Group("order")
	{
		orderRouter.POST("", orderApi.Create)           // Create a new order
		orderRouter.GET("list", orderApi.GetOrderList)  // List orders
		orderRouter.GET(":id", orderApi.GetOrderDetail) // Get order by ID
		orderRouter.DELETE("", orderApi.Delete)         // Delete an order
	}
}
