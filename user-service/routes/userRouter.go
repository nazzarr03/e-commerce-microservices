package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/user-service/controllers"
	"github.com/nazzarr03/user-service/middleware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")
	user := api.Group("/user")

	user.Post("/signup", controllers.SignUp)
	user.Post("/login", controllers.Login)

	user.Use(middleware.Authentication())

	user.Get("/user", controllers.GetUsers)
	user.Get("/user/:id", controllers.GetUserByID)
}
