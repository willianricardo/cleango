package model

type Order struct {
	ID         string      `json:"id"`
	OrderDate  string      `json:"order_date"`
	CustomerID string      `json:"customer_id"`
	Customer   Customer    `json:"customer"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
