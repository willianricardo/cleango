package repository

import (
	"api/model"
	"database/sql"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repository *OrderRepository) GetOrders() ([]model.Order, error) {
	var orders []model.Order = []model.Order{}

	orderRows, err := repository.db.Query(`
		SELECT o.id, o.customer_id, o.order_date,
		       c.id, c.name
		FROM orders o
		INNER JOIN customers c ON o.customer_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer orderRows.Close()

	var customer model.Customer

	for orderRows.Next() {
		order := model.Order{}
		err := orderRows.Scan(
			&order.ID, &order.CustomerID, &order.OrderDate,
			&customer.ID, &customer.Name,
		)
		if err != nil {
			return nil, err
		}

		order.Customer = customer
		order.OrderItems = []model.OrderItem{}
		orders = append(orders, order)
	}

	orderItemRows, err := repository.db.Query(`
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
			   p.id, p.name, p.price
		FROM order_items oi
		INNER JOIN products p ON oi.product_id = p.id
	`)
	if err != nil {
		return nil, err
	}
	defer orderItemRows.Close()

	var orderItem model.OrderItem
	var product model.Product

	for orderItemRows.Next() {
		err := orderItemRows.Scan(
			&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price,
			&product.ID, &product.Name, &product.Price,
		)
		if err != nil {
			return nil, err
		}

		for i := range orders {
			if orders[i].ID == orderItem.OrderID {
				orderItem.Product = product
				orders[i].OrderItems = append(orders[i].OrderItems, orderItem)
				break
			}
		}
	}

	return orders, nil
}

func (repository *OrderRepository) GetOrderByID(orderID string) (model.Order, error) {
	var order model.Order

	orderRow := repository.db.QueryRow(`
		SELECT o.id, o.customer_id, o.order_date,
			   c.id, c.name
		FROM orders o
		INNER JOIN customers c ON o.customer_id = c.id
		WHERE o.id = ?
	`, orderID)

	customer := model.Customer{}
	err := orderRow.Scan(
		&order.ID, &order.CustomerID, &order.OrderDate,
		&customer.ID, &customer.Name,
	)
	if err != nil {
		return order, err
	}

	order.Customer = customer

	orderItemRows, err := repository.db.Query(`
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
			   p.id, p.name, p.price
		FROM order_items oi
		INNER JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		return order, err
	}
	defer orderItemRows.Close()

	var orderItem model.OrderItem
	var product model.Product

	for orderItemRows.Next() {
		err := orderItemRows.Scan(
			&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price,
			&product.ID, &product.Name, &product.Price,
		)
		if err != nil {
			return order, err
		}

		orderItem.Product = product
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	return order, nil
}

func (repository *OrderRepository) CreateOrder(order model.Order) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	// Insert order
	_, err = tx.Exec("INSERT INTO orders (id, customer_id, order_date) VALUES (?, ?, ?)", order.ID, order.CustomerID, order.OrderDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert order items
	for _, orderItem := range order.OrderItems {
		_, err = tx.Exec("INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?)", orderItem.ID, order.ID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *OrderRepository) UpdateOrder(order model.Order) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	// Update order
	_, err = tx.Exec("UPDATE orders SET customer_id = ?, order_date = ? WHERE id = ?", order.CustomerID, order.OrderDate, order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing order items
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert updated order items
	for _, orderItem := range order.OrderItems {
		_, err = tx.Exec("INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?)", orderItem.ID, order.ID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *OrderRepository) DeleteOrder(orderID string) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	// Delete order items
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete order
	_, err = tx.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
