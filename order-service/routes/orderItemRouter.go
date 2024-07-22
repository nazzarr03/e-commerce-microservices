package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/order-service/controllers"
	"github.com/nazzarr03/order-service/middleware"
)

func Router(app *fiber.App) {
	api := app.Group("/api/v1")
	orderItem := api.Group("/order-item")

	orderItem.Use(middleware.Authentication())

	orderItem.Get("/order-item", controllers.GetOrderItem)
	orderItem.Get("/order-item/:id", controllers.GetOrderItem)
	orderItem.Put("/order-item/:id", controllers.UpdateOrderItem)
	orderItem.Delete("/order-item/:id", controllers.DeleteOrderItem)
}
