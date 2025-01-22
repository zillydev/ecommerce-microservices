package order

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Products  []string  `json:"products"` // array of product IDs and quantities, e.g. ["1:2", "2:3"]
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewOrderRequest struct {
	UserID   int      `json:"userId"`
	Products []string `json:"products"` // array of product IDs and quantities, e.g. ["1:2", "2:3"]
}

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type UpdateOrderStatusRequest struct {
	OrderID int    `json:"orderId"`
	Status  string `json:"status"`
}

type OrderStatusUpdate struct {
	OrderID int    `json:"orderId"`
	Status  string `json:"status"`
	UserID  int    `json:"userId"`
}
