package handler

import (
	"api/model"
	"api/repository"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerRepository *repository.CustomerRepository
}

func NewCustomerHandler(customerRepository *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{
		customerRepository: customerRepository,
	}
}

// GetCustomers godoc
// @Summary List customers
// @Description Get all customers
// @Tags customers
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Customer
// @Failure 500 {object} map[string]string
// @Router /customers [get]
func (handler *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	customers, err := handler.customerRepository.GetCustomers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve customers",
		})
	}
	return c.Status(fiber.StatusOK).JSON(customers)
}

// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Get a customer by its ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Success 200 {object} model.Customer
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [get]
func (handler *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
	customerID := c.Params("id")
	customer, err := handler.customerRepository.GetCustomerByID(customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve customer",
		})
	}
	return c.Status(fiber.StatusOK).JSON(customer)
}

// CreateCustomer godoc
// @Summary Create customer
// @Description Create a new customer
// @Tags customers
// @Accept  json
// @Produce  json
// @Param customer body model.Customer true "Customer to create"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers [post]
func (handler *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var customer model.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer data",
		})
	}
	if err := handler.customerRepository.CreateCustomer(customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

// UpdateCustomer godoc
// @Summary Update customer
// @Description Update an existing customer
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Param customer body model.Customer true "Customer to update"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [put]
func (handler *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	customerID := c.Params("id")
	var customer model.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer data",
		})
	}
	customer.ID = customerID
	if err := handler.customerRepository.UpdateCustomer(customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteCustomer godoc
// @Summary Delete customer
// @Description Delete a customer by its ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Success 200
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [delete]
func (handler *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	customerID := c.Params("id")
	if err := handler.customerRepository.DeleteCustomer(customerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
