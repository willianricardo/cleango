package main

import (
	"api/database"
	"api/handler"
	"api/repository"
	"api/routes"
	"log"

	_ "api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Initialize the database
	db, err := database.InitializeDB("database/database.db", "file://database/migrations")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the repository instances
	productRepository := repository.NewProductRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	orderRepository := repository.NewOrderRepository(db)

	// Create the handler instances
	productHandler := handler.NewProductHandler(productRepository)
	customerHandler := handler.NewCustomerHandler(customerRepository)
	orderHandler := handler.NewOrderHandler(orderRepository)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Define the API routes
	routes.SetupProductRoutes(app, productHandler)
	routes.SetupCustomerRoutes(app, customerHandler)
	routes.SetupOrderRoutes(app, orderHandler)

	// Swagger docs route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start the HTTP server
	log.Fatal(app.Listen(":8080"))
}
