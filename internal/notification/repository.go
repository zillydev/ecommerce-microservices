package notification

import (
	"database/sql"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateNotification(notification *NewNotificationRequest) (*Notification, error) {
	query := "INSERT INTO notifications (userId, type, content, sentAt, read) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var notificationId int
	timeNow := time.Now()
	err := r.DB.QueryRow(query, notification.UserID, notification.Type, notification.Content, timeNow, false).Scan(&notificationId)
	if err != nil {
		return nil, err
	}
	return &Notification{ID: notificationId, UserID: notification.UserID, Type: notification.Type, Content: notification.Content, SentAt: timeNow, Read: false}, nil
}

func (r *Repository) MarkNotificationRead(notificationID int, userID int) error {
	query := "UPDATE notifications SET read = true WHERE id = $1 AND userId = $2"
	_, err := r.DB.Exec(query, notificationID, userID)
	return err
}

func (r *Repository) GetUnreadNotifications(userID int) ([]Notification, error) {
	query := "SELECT id, userId, type, content, sentAt, read FROM notifications WHERE userId = $1 AND read = false"
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.Content, &notification.SentAt, &notification.Read); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
