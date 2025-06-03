package routes

import (
	"api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
	router := app.Group("/orders")
	router.Get("", orderHandler.GetOrders)
	router.Get("/:id", orderHandler.GetOrderByID)
	router.Post("", orderHandler.CreateOrder)
	router.Put("/:id", orderHandler.UpdateOrder)
	router.Delete("/:id", orderHandler.DeleteOrder)
}
