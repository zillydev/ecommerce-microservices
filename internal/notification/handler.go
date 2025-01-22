package notification

import (
	"bytes"
	"database/sql"
	"ecommerce-microservices/pkg/kafka"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	Service        *Service
	Kafka          *kafka.Kafka
	UserServiceURL string
}

func NewHandler(db *sql.DB, kafka *kafka.Kafka, userServiceURL string) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{Service: service, Kafka: kafka, UserServiceURL: userServiceURL}
}

func (h *Handler) PostNotification(w http.ResponseWriter, r *http.Request) {
	var notificationRequest NewNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&notificationRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing notification request: %v", err), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	userID, err := strconv.Atoi(queryParams.Get("userId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	notificationRequest.UserID = userID

	notification, err := h.Service.CreateNotification(&notificationRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating notification: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	var notificationID int
	if err := json.NewDecoder(r.Body).Decode(&notificationID); err != nil {
		http.Error(w, fmt.Sprintf("error parsing notification ID: %v", err), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	userID, err := strconv.Atoi(queryParams.Get("userId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.Service.MarkNotificationRead(notificationID, userID); err != nil {
		http.Error(w, fmt.Sprintf("error marking notification as read: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userID, err := strconv.Atoi(queryParams.Get("userId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	notifications, err := h.Service.GetUnreadNotifications(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting unread notifications: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}

func (h *Handler) HandleOrderStatusUpdate(orderStatusUpdate OrderStatusUpdate) error {
	resp, err := http.Get(fmt.Sprintf("%s/user/%d", h.UserServiceURL, orderStatusUpdate.UserID))
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	defer resp.Body.Close()

	var user *User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Errorf("failed to decode user: %v", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	for _, notification := range user.PreferredNotifications {
		if notification == "order_updates" {
			notification := &NewNotificationRequest{
				UserID:  orderStatusUpdate.UserID,
				Type:    "order_updates",
				Content: "Your order with ID " + strconv.Itoa(orderStatusUpdate.OrderID) + " has been " + orderStatusUpdate.Status,
			}
			_, err = h.Service.CreateNotification(notification)
			return err
		}
	}

	return nil
}

func (h *Handler) SendDailyPromotionNotifications() error {
	preferencesRequest := GetUsersByPreferencesRequest{
		PreferredNotifications: []string{"promotions"},
	}
	body, err := json.Marshal(preferencesRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences request: %v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", h.UserServiceURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	for _, user := range users {
		notification := &NewNotificationRequest{
			UserID:  user.ID,
			Type:    "promotions",
			Content: "Get 20% off your next purchase!",
		}
		_, err = h.Service.CreateNotification(notification)
		if err != nil {
			return fmt.Errorf("failed to create notification: %v", err)
		}
	}

	return nil
}
