package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("Incoming request: %s %s \n", c.Method(), c.Path())
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "*",
		ExposeHeaders:    "*",
		AllowCredentials: false,
	}))

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading .env file")
		panic(err)
	}

	app.Post("/api/v1/user/signup", forwardToService("USER_SERVICE_URL"))
	app.Post("/api/v1/user/login", forwardToService("USER_SERVICE_URL"))
	app.Get("/api/v1/user/user", forwardToService("USER_SERVICE_URL"))
	app.Get("/api/v1/user/user/:id", forwardToService("USER_SERVICE_URL"))

	app.Get("/api/v1/product/product", forwardToService("PRODUCT_SERVICE_URL"))
	app.Get("/api/v1/product/product/:id", forwardToService("PRODUCT_SERVICE_URL"))
	app.Post("/api/v1/product/product", forwardToService("PRODUCT_SERVICE_URL"))
	app.Put("/api/v1/product/product/:id", forwardToService("PRODUCT_SERVICE_URL"))

	app.Listen(":8080")

}

func forwardToService(serviceEnvVar string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceURL := os.Getenv(serviceEnvVar)
		if serviceURL == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is not set", serviceEnvVar),
			})
		}

		targetURL := serviceURL + c.Path()
		body := bytes.NewReader(c.Body())
		req, err := http.NewRequest(c.Method(), targetURL, body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		req.Header = c.GetReqHeaders()
		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Status(resp.StatusCode)

		for key, value := range resp.Header {
			c.Set(key, value[0])
		}

		return c.Send(bodyBytes)
	}
}
