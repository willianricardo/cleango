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

func (handler *OrderHandler) GetOrders(c *fiber.Ctx) error {
	orders, err := handler.orderRepository.GetOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

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
