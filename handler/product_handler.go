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

// GetProducts godoc
// @Summary List products
// @Description Get all products
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (handler *ProductHandler) GetProducts(c *fiber.Ctx) error {
	products, err := handler.productRepository.GetProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
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

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body model.Product true "Product to create"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
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

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body model.Product true "Product to update"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product by its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (handler *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	if err := handler.productRepository.DeleteProduct(productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
