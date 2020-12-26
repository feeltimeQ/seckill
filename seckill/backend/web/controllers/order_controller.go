package controllers

import (
	"backend/seckill/services"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

//OrderController 订单控制
type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

//Get 查询订单相关信息
func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}

}
