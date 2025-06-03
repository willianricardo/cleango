package handler

import (
	"api/model"
	"api/repository"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productRepository *repository.ProductRepository
}

func NewProductHandler(productRepository *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepository: productRepository,
	}
}

func (handler *ProductHandler) GetProducts(c *fiber.Ctx) error {
	products, err := handler.productRepository.GetProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

func (handler *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	productID := c.Params("id")
	product, err := handler.productRepository.GetProductByID(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve product",
		})
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

func (handler *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product model.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product data",
		})
	}
	if err := handler.productRepository.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (handler *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	var product model.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product data",
		})
	}
	product.ID = productID
	if err := handler.productRepository.UpdateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (handler *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	if err := handler.productRepository.DeleteProduct(productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
