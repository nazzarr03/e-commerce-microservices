package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/product-service/config"
	"github.com/nazzarr03/product-service/models"
)

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	config.Db.Find(&products)

	if len(products) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Products not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":    products,
		"message": "Successfully retrieved products",
	})
}

func GetProductByID(c *fiber.Ctx) error {
	var product models.Product
	idStr := c.Params("product_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	config.Db.Where("id = ?", id).First(&product)
	if product.ProductID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":    product,
		"message": "Successfully retrieved product",
	})
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if product.Name == "" || product.Price == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and price are required",
		})
	}

	if err := config.Db.Create(&product).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create product",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data":    product,
		"message": "Successfully created product",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product
	var newProduct models.Product

	if err := c.BodyParser(&newProduct); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	idStr := c.Params("product_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	config.Db.Where("id = ?", id).First(&product)
	if product.ProductID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if newProduct.Name != "" {
		product.Name = newProduct.Name
	}

	if newProduct.Price != 0 {
		product.Price = newProduct.Price
	}

	if newProduct.Stock != 0 {
		product.Stock = newProduct.Stock
	}

	config.Db.Save(&product)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":    product,
		"message": "Successfully updated product",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	var product models.Product
	config.Db.First(&product, c.Params("product_id"))
	if product.ProductID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	config.Db.Delete(&product)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted product",
	})
}
