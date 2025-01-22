// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"time"
)

type AddProductInput struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type CreateOrderInput struct {
	Products []string `json:"products"`
}

type Mutation struct {
}

type Notification struct {
	ID      int       `json:"id"`
	UserID  int       `json:"userId"`
	Type    string    `json:"type"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sentAt"`
	Read    bool      `json:"read"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Products  []string  `json:"products"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostNotificationInput struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type Query struct {
}

type RegisterUserInput struct {
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	PreferredNotifications []string `json:"preferredNotifications"`
}

type RegisterUserResult struct {
	User        *User  `json:"user"`
	AccessToken string `json:"accessToken"`
}

type UpdateOrderStatusInput struct {
	OrderID int    `json:"orderId"`
	Status  string `json:"status"`
}

type UpdatePreferencesInput struct {
	PreferredNotifications []string `json:"preferredNotifications"`
}

type User struct {
	ID                     int      `json:"id"`
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	PreferredNotifications []string `json:"preferredNotifications"`
}