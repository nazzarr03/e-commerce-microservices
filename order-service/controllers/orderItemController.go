package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/order-service/config"
	"github.com/nazzarr03/order-service/models"
	"github.com/nazzarr03/order-service/utils"
	"github.com/streadway/amqp"
)

func CreateOrderItem(c *fiber.Ctx) error {
	var orderItem models.OrderItem
	var order models.Order
	if err := c.BodyParser(&orderItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if orderItem.Quantity == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity is required",
		})
	}

	config.Db.First(&order, c.Params("order_id"))
	if order.OrderID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	orderItem.OrderID = order.OrderID

	productID := c.Params("product_id")
	productIDUint, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot parse product ID",
		})
	}
	orderItem.ProductID = uint(productIDUint)

	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "RabbitMQ channel error",
		})
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "RabbitMQ queue error",
		})
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "RabbitMQ consume error",
		})
	}

	corrId := utils.RandomString(32)

	err = ch.Publish(
		"",
		"product_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(productID),
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "RabbitMQ publish error",
		})
	}

	for d := range msgs {
		if d.CorrelationId == corrId {
			var response map[string]interface{}
			err := json.Unmarshal(d.Body, &response)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot parse response",
				})
			}

			productPrice, ok := response["price"].(float64)
			if !ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Invalid price format",
				})
			}

			productStock, ok := response["stock"].(float64)
			if !ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Invalid stock format",
				})
			}

			if productStock < float64(orderItem.Quantity) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Not enough stock",
				})
			}

			orderItem.Price = productPrice
			orderItem.TotalPrice = productPrice * float64(orderItem.Quantity)

			break
		}
	}

	if err := config.Db.Create(&orderItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create order item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    orderItem,
		"message": "Successfully created order item",
	})
}

func GetOrderItems(c *fiber.Ctx) error {
	var orderItems []models.OrderItem
	config.Db.Find(&orderItems)

	if len(orderItems) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order items not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    orderItems,
		"message": "Successfully retrieved order items",
	})
}

func GetOrderItem(c *fiber.Ctx) error {
	var orderItem models.OrderItem
	config.Db.First(&orderItem, c.Params("id"))

	if orderItem.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    orderItem,
		"message": "Successfully retrieved order item",
	})
}

func UpdateOrderItem(c *fiber.Ctx) error {
	var orderItem models.OrderItem
	var newOrderItem models.OrderItem

	if err := c.BodyParser(&newOrderItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	config.Db.First(&orderItem, c.Params("id"))
	if orderItem.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	orderItem.Quantity = newOrderItem.Quantity
	orderItem.TotalPrice = orderItem.Price * float64(orderItem.Quantity)

	config.Db.Save(&orderItem)

	return c.JSON(fiber.Map{
		"data":    orderItem,
		"message": "Successfully updated order item",
	})
}

func DeleteOrderItem(c *fiber.Ctx) error {
	var orderItem models.OrderItem
	config.Db.First(&orderItem, c.Params("id"))
	if orderItem.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	config.Db.Delete(&orderItem)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted order item",
	})
}
