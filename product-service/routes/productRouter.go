package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/product-service/controllers"
	"github.com/nazzarr03/product-service/middleware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")
	product := api.Group("/product")

	product.Use(middleware.Authentication())

	product.Get("/product", controllers.GetProducts)
	product.Get("/product/:id", controllers.GetProductByID)
	product.Post("/product", controllers.CreateProduct)
	product.Put("/product/:id", controllers.UpdateProduct)
}
