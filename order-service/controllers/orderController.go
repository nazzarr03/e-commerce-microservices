package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/order-service/config"
	"github.com/nazzarr03/order-service/models"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	order.UserID = c.Locals("userID").(uint)
	order.Date = time.Now()

	for i := range order.OrderItems {
		order.Total += order.OrderItems[i].TotalPrice
	}

	if err := config.Db.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    order,
		"message": "Successfully created order",
	})
}

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order
	config.Db.Preload("OrderItems").Find(&orders)

	if len(orders) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Orders not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    orders,
		"message": "Successfully retrieved orders",
	})
}

func GetOrder(c *fiber.Ctx) error {
	var order models.Order
	config.Db.Preload("OrderItems").First(&order, c.Params("id"))

	if order.OrderID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    order,
		"message": "Successfully retrieved order",
	})
}

func DeleteOrder(c *fiber.Ctx) error {
	var order models.Order
	config.Db.First(&order, c.Params("id"))
	if order.OrderID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	config.Db.Delete(&order)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted order",
	})
}
