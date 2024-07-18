package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/user-service/config"
	"github.com/nazzarr03/user-service/middleware"
	"github.com/nazzarr03/user-service/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	config.Db.Find(&users)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": users,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	var user models.User
	config.Db.First(&user, c.Params("id"))
	if user.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": user,
	})
}

func SignUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email and password are required",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot hash password",
		})
	}

	user.Password = string(hashedPassword)

	if err := config.Db.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create user",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data":    user,
		"message": "Successfully created user",
	})
}

func Login(c *fiber.Ctx) error {
	var user models.User
	var existingUser models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	config.Db.Where("username = ?", user.Username).First(&existingUser)
	if existingUser.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	token, err := middleware.GenerateToken(existingUser.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot generate JWT",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "Successfully logged in",
	})
}
