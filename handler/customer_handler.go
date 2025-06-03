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

func (handler *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	customers, err := handler.customerRepository.GetCustomers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve customers",
		})
	}
	return c.Status(fiber.StatusOK).JSON(customers)
}

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

func (handler *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	customerID := c.Params("id")
	if err := handler.customerRepository.DeleteCustomer(customerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
