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

var testCustomerDB *sql.DB

func setupCustomerTestApp() *fiber.App {
	dbFile := "test_customer_integration.db"
	os.Remove(dbFile)
	migrationDir := "file://../database/migrations"
	db, err := database.InitializeDB(dbFile, migrationDir)
	if err != nil {
		panic(err)
	}
	testCustomerDB = db
	customerRepo := repository.NewCustomerRepository(db)
	customerHandler := handler.NewCustomerHandler(customerRepo)
	app := fiber.New()
	routes.SetupCustomerRoutes(app, customerHandler)
	return app
}

func teardownCustomerTestDB() {
	if testCustomerDB != nil {
		testCustomerDB.Close()
		os.Remove("test_customer_integration.db")
	}
}

func TestCustomerIntegration(t *testing.T) {
	app := setupCustomerTestApp()
	defer teardownCustomerTestDB()

	// Test POST /customers
	customer := model.Customer{ID: "1", Name: "Test Customer"}
	body, _ := json.Marshal(customer)
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Test GET /customers
	req = httptest.NewRequest(http.MethodGet, "/customers", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var customers []model.Customer
	json.NewDecoder(resp.Body).Decode(&customers)
	assert.Len(t, customers, 1)
	assert.Equal(t, customer.ID, customers[0].ID)

	// Test GET /customers/:id
	req = httptest.NewRequest(http.MethodGet, "/customers/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var gotCustomer model.Customer
	json.NewDecoder(resp.Body).Decode(&gotCustomer)
	assert.Equal(t, customer.ID, gotCustomer.ID)

	// Test PUT /customers/:id
	updated := model.Customer{ID: "1", Name: "Updated Customer"}
	body, _ = json.Marshal(updated)
	req = httptest.NewRequest(http.MethodPut, "/customers/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Test DELETE /customers/:id
	req = httptest.NewRequest(http.MethodDelete, "/customers/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Confirm deletion
	req = httptest.NewRequest(http.MethodGet, "/customers/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
