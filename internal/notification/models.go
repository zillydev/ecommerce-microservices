package notification

import (
	"time"
)

type NewNotificationRequest struct {
	UserID  int    `json:"userId"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Notification struct {
	ID      int       `json:"id"`
	UserID  int       `json:"userId"`
	Type    string    `json:"type"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sentAt"`
	Read    bool      `json:"read"`
}

type OrderStatusUpdate struct {
	OrderID int    `json:"orderId"`
	Status  string `json:"status"`
	UserID  int    `json:"userId"`
}

type User struct {
	ID                     int      `json:"id"`
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	PreferredNotifications []string `json:"preferredNotifications"`
}

type GetUsersByPreferencesRequest struct {
	PreferredNotifications []string `json:"preferredNotifications"`
}
