package main

import (
	"backend/seckill/backend/web/controllers"
	"backend/seckill/common"
	"backend/seckill/repositories"
	"backend/seckill/services"
	"context"

	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"

	"github.com/kataras/iris"
)

func main() {

	//创建iris实例
	app := iris.New()

	//设置错误模式，在MVC模式下提示错误
	app.Logger().SetLevel("debug")

	//注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	//设置模板目标
	app.HandleDir("assets", iris.Dir("./backend/web/assets"))

	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
