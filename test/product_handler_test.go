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

var testDB *sql.DB

func setupTestApp() *fiber.App {
	dbFile := "test_product_integration.db"
	os.Remove(dbFile)
	migrationDir := "file://../database/migrations"
	db, err := database.InitializeDB(dbFile, migrationDir)
	if err != nil {
		panic(err)
	}
	testDB = db
	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)
	app := fiber.New()
	routes.SetupProductRoutes(app, productHandler)
	return app
}

func teardownTestDB() {
	if testDB != nil {
		testDB.Close()
		os.Remove("test_product_integration.db")
	}
}

func TestProductIntegration(t *testing.T) {
	app := setupTestApp()
	defer teardownTestDB()

	// Test POST /products
	product := model.Product{ID: "1", Name: "Test Product", Price: 9.99}
	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Test GET /products
	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var products []model.Product
	json.NewDecoder(resp.Body).Decode(&products)
	assert.Len(t, products, 1)
	assert.Equal(t, product.ID, products[0].ID)

	// Test GET /products/:id
	req = httptest.NewRequest(http.MethodGet, "/products/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var gotProduct model.Product
	json.NewDecoder(resp.Body).Decode(&gotProduct)
	assert.Equal(t, product.ID, gotProduct.ID)

	// Test PUT /products/:id
	updated := model.Product{ID: "1", Name: "Updated Product", Price: 19.99}
	body, _ = json.Marshal(updated)
	req = httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Test DELETE /products/:id
	req = httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Confirm deletion
	req = httptest.NewRequest(http.MethodGet, "/products/1", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
