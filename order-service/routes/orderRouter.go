package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/order-service/controllers"
	"github.com/nazzarr03/order-service/middleware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")
	order := api.Group("/order")

	order.Use(middleware.Authentication())

	order.Post("/order", controllers.CreateOrder)
	order.Get("/order", controllers.GetOrders)
	order.Get("/order/:id", controllers.GetOrder)
	order.Delete("/order/:id", controllers.DeleteOrder)

	order.Post("/order/:order_id/product/:product_id", controllers.CreateOrderItem)
}
