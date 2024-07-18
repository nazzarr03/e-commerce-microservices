package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/product-service/routes"
)

func main() {
	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8082")
}
