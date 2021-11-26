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
	index.Post("/getOrders", getOrders)
	index.Get("/getOrdersForViewer", getOrdersForViewer)
	index.Post("/getOrderSubOrders", getOrderSubOrders)
	index.Get("/getDefaultList", getDefaultList)
	index.Post("/getShippedSubOrders", getShippedSubOrders)
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
	index.Post("/getStatistics", getStatistics)
	index.Get("/getVendors", getVendors)
	index.Post("/updateVendors", updateVendors)

	// TLS:
	// app.Run(iris.TLS("0.0.0.0:2083", "/root/cert/cf.pem", "/root/cert/cf.key"))
	app.Run(iris.Addr("0.0.0.0:5050"))
}
