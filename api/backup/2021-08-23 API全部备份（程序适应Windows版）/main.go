package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {

	go WebSocket()

	crs := cors.New(cors.Options{ //crs相当于一个中间件，允许所有主机通过
		AllowedOrigins:   []string{"*"}, //
		AllowCredentials: true,
	})
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	index := app.Party("/", crs) //所有请求先过crs中间件
	index.Get("/getOrders", getOrders)
	index.Get("/getOrdersForViewer", getOrdersForViewer)
	index.Get("/getOrderSubOrders", getOrderSubOrders)
	index.Get("/getDefaultList", getDefaultList)
	index.Get("/getShippedSubOrders", getShippedSubOrders)
	index.Get("/getProducts", getProducts)
	index.Get("/getProductsClassList", getProductsClassList)
	index.Post("/deleteProduct", deleteProduct)
	index.Post("/updateProductPosition", updateProductPosition)
	index.Post("/updateProduct", updateProduct)
	index.Post("/getLogistics", getLogistics)
	index.Post("/updateLogistics", updateLogistics)
	index.Post("/updateOrder", updateOrder)
	index.Post("/deleteOrder", deleteOrder)
	index.Post("/updatePackets", updatePackets)
	_ = app.Run(iris.Addr("192.168.0.2:8081"), iris.WithoutServerError(iris.ErrServerClosed))
}
