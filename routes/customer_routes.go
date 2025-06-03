package routes

import (
	"api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupCustomerRoutes(app *fiber.App, customerHandler *handler.CustomerHandler) {
	router := app.Group("/customers")
	router.Get("", customerHandler.GetCustomers)
	router.Get("/:id", customerHandler.GetCustomerByID)
	router.Post("", customerHandler.CreateCustomer)
	router.Put("/:id", customerHandler.UpdateCustomer)
	router.Delete("/:id", customerHandler.DeleteCustomer)
}
