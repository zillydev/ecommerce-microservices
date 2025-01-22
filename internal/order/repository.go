package order

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Repository struct {
	DB                *sql.DB
	ProductServiceURL string
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateOrder(order *NewOrderRequest, total float64) (*Order, error) {
	timeNow := time.Now()
	var id int
	err := r.DB.QueryRow("INSERT INTO orders (userId, products, total, status, createdAt) VALUES ($1, $2, $3, $4, $5) RETURNING id", order.UserID, pq.Array(order.Products), total, "packaged", timeNow).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &Order{ID: id, UserID: order.UserID, Products: order.Products, Total: total, Status: "packaged", CreatedAt: timeNow}, nil
}

func (r *Repository) GetAllOrdersByUser(userID int) ([]Order, error) {
	rows, err := r.DB.Query("SELECT id, userId, products, total, status, createdAt FROM orders WHERE userId = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, pq.Array(&order.Products), &order.Total, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) GetOrderByID(id int) (*Order, error) {
	var order Order
	err := r.DB.QueryRow("SELECT id, userId, products, total, status, createdAt FROM orders WHERE id = $1", id).Scan(&order.ID, &order.UserID, pq.Array(&order.Products), &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order with id %d not found", id)
		}
		return nil, err
	}
	return &order, nil
}

func (r *Repository) UpdateOrderStatus(id int, status string) (int, error) {
	var userID int
	err := r.DB.QueryRow("UPDATE orders SET status = $1 WHERE id = $2 RETURNING userId", status, id).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
