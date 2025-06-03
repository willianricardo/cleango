package repository

import (
	"api/model"
	"database/sql"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (repository *CustomerRepository) GetCustomers() ([]model.Customer, error) {
	var customers []model.Customer = []model.Customer{}
	rows, err := repository.db.Query("SELECT id, name FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(&customer.ID, &customer.Name)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil

}

func (repository *CustomerRepository) GetCustomerByID(id string) (model.Customer, error) {
	var customer model.Customer
	row := repository.db.QueryRow("SELECT id, name FROM customers WHERE id = ?", id)
	err := row.Scan(&customer.ID, &customer.Name)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (repository *CustomerRepository) CreateCustomer(customer model.Customer) error {
	_, err := repository.db.Exec("INSERT INTO customers (id, name) VALUES (?, ?)", customer.ID, customer.Name)
	if err != nil {
		return err
	}
	return nil
}

func (repository *CustomerRepository) UpdateCustomer(customer model.Customer) error {
	_, err := repository.db.Exec("UPDATE customers SET name = ? WHERE id = ?", customer.Name, customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repository *CustomerRepository) DeleteCustomer(id string) error {
	_, err := repository.db.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
