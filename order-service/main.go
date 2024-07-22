package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/order-service/routes"
)

func main() {
	app := fiber.New()

	routes.Setup(app)
	routes.Router(app)

	app.Listen(":8083")

}
