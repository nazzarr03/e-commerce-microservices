package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/logger/config"
)

func main() {
	go config.ConsumeFromRabbitMQ()

	app := fiber.New()

	log.Fatal(app.Listen(":8084"))
}
