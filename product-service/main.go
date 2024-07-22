package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/product-service/config"
	"github.com/nazzarr03/product-service/routes"
)

func main() {
	config.ConnectRabbitMQ()
	defer config.RabbitMQConn.Close()

	go config.ProductConsumer()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8082")

	select {}
}
