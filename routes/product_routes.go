package routes

import (
	"api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	router := app.Group("/products")
	router.Get("", productHandler.GetProducts)
	router.Get("/:id", productHandler.GetProductByID)
	router.Post("", productHandler.CreateProduct)
	router.Put("/:id", productHandler.UpdateProduct)
	router.Delete("/:id", productHandler.DeleteProduct)
}
