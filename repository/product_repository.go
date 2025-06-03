package repository

import (
	"api/model"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repository *ProductRepository) GetProducts() ([]model.Product, error) {
	var products []model.Product = []model.Product{}
	rows, err := repository.db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *ProductRepository) GetProductByID(id string) (model.Product, error) {
	var product model.Product
	row := repository.db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id)
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (repository *ProductRepository) CreateProduct(product model.Product) error {
	_, err := repository.db.Exec("INSERT INTO products (id, name, price) VALUES (?, ?, ?)", product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ProductRepository) UpdateProduct(product model.Product) error {
	_, err := repository.db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ProductRepository) DeleteProduct(id string) error {
	_, err := repository.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
