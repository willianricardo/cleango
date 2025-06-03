package handler

import (
	"api/model"
	"api/repository"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderRepository *repository.OrderRepository
}

func NewOrderHandler(orderRepository *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderRepository: orderRepository,
	}
}

// GetOrders godoc
// @Summary List orders
// @Description Get all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Order
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (handler *OrderHandler) GetOrders(c *fiber.Ctx) error {
	orders, err := handler.orderRepository.GetOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get an order by its ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (handler *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("id")
	order, err := handler.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(order)
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body model.Order true "Order to create"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (handler *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order model.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := handler.orderRepository.CreateOrder(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created",
	})
}

// UpdateOrder godoc
// @Summary Update order
// @Description Update an existing order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param order body model.Order true "Order to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [put]
func (handler *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var order model.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	order.ID = orderID
	err := handler.orderRepository.UpdateOrder(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order updated",
	})
}

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete an order by its ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [delete]
func (handler *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	err := handler.orderRepository.DeleteOrder(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order deleted",
	})
}
