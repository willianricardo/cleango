package handler_test

import (
	"api/handler"
	"api/model"
	"api/repository"
	"api/routes"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"api/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var testOrderDB *sql.DB

func setupOrderTestApp() *fiber.App {
	dbFile := "test_order_integration.db"
	os.Remove(dbFile)
	migrationDir := "file://../database/migrations"
	db, err := database.InitializeDB(dbFile, migrationDir)
	if err != nil {
		panic(err)
	}
	testOrderDB = db

	// Create repositories and handlers for products, customers, and orders
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	productHandler := handler.NewProductHandler(productRepo)
	customerHandler := handler.NewCustomerHandler(customerRepo)
	orderHandler := handler.NewOrderHandler(orderRepo)

	app := fiber.New()
	// Register all routes needed for the integration test
	routes.SetupProductRoutes(app, productHandler)
	routes.SetupCustomerRoutes(app, customerHandler)
	routes.SetupOrderRoutes(app, orderHandler)
	return app
}

func teardownOrderTestDB() {
	if testOrderDB != nil {
		testOrderDB.Close()
		os.Remove("test_order_integration.db")
	}
}

func TestOrderIntegration(t *testing.T) {
	app := setupOrderTestApp()
	defer teardownOrderTestDB()

	// Test POST /products
	product := model.Product{ID: "p1", Name: "Test Product", Price: 9.99}
	productBody, _ := json.Marshal(product)
	productReq := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(productBody))
	productReq.Header.Set("Content-Type", "application/json")
	productResp, err := app.Test(productReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, productResp.StatusCode)

	// Test POST /customers
	customer := model.Customer{ID: "c1", Name: "Test Customer"}
	customerBody, _ := json.Marshal(customer)
	customerReq := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(customerBody))
	customerReq.Header.Set("Content-Type", "application/json")
	customerResp, err := app.Test(customerReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, customerResp.StatusCode)

	// Test POST /orders
	order := model.Order{
		ID:         "1",
		OrderDate:  "2024-01-01",
		CustomerID: "c1",
		OrderItems: []model.OrderItem{
			{
				ID:        "oi1",
				OrderID:   "1",
				ProductID: "p1",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}
	body, _ := json.Marshal(order)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Test GET /orders
	req = httptest.NewRequest(http.MethodGet, "/orders", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var orders []model.Order
	json.NewDecoder(resp.Body).Decode(&orders)
	assert.Len(t, orders, 1)
	assert.Equal(t, order.ID, orders[0].ID)
	assert.Len(t, orders[0].OrderItems, 1)
	assert.Equal(t, order.OrderItems[0].ID, orders[0].OrderItems[0].ID)

	// Test GET /orders/:id
	req = httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var gotOrder model.Order
	json.NewDecoder(resp.Body).Decode(&gotOrder)
	assert.Equal(t, order.ID, gotOrder.ID)
	assert.Len(t, gotOrder.OrderItems, 1)
	assert.Equal(t, order.OrderItems[0].ID, gotOrder.OrderItems[0].ID)

	// Test PUT /orders/:id
	updated := model.Order{
		ID:         "1",
		OrderDate:  "2024-01-02",
		CustomerID: "c1",
		OrderItems: []model.OrderItem{
			{
				ID:        "oi1",
				OrderID:   "1",
				ProductID: "p1",
				Quantity:  3,
				Price:     12.0,
			},
		},
	}
	body, _ = json.Marshal(updated)
	req = httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Test DELETE /orders/:id
	req = httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Confirm deletion
	req = httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
