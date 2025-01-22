package order

import (
	"database/sql"
	"ecommerce-microservices/pkg/kafka"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service           *Service
	Kafka             *kafka.Kafka
	ProductServiceURL string
}

func NewHandler(db *sql.DB, kafka *kafka.Kafka, productServiceURL string) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{Service: service, Kafka: kafka, ProductServiceURL: productServiceURL}
}

func (h *Handler) GetAllOrdersByUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userID, err := strconv.Atoi(queryParams.Get("userId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	orders, err := h.Service.GetAllOrdersByUser(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting orders: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing order ID: %v", err), http.StatusBadRequest)
		return
	}

	order, err := h.Service.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting order: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest NewOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing order request: %v", err), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	userID, err := strconv.Atoi(queryParams.Get("userId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}
	orderRequest.UserID = userID

	// calculate total price
	var total float64
	for _, purchase := range orderRequest.Products {
		productID := strings.Split(purchase, ":")[0]
		quantity, err := strconv.Atoi(strings.Split(purchase, ":")[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing quantity: %v", err), http.StatusBadRequest)
			return
		}

		resp, err := http.Get(fmt.Sprintf("%s/get/%s", h.ProductServiceURL, productID))
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting product: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var product Product
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			http.Error(w, fmt.Sprintf("error decoding product: %v", err), http.StatusInternalServerError)
			return
		}

		total += product.Price * float64(quantity)
	}

	order, err := h.Service.CreateOrder(&orderRequest, total)
	if err != nil {
		http.Error(w, fmt.Sprintf("error adding order: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var updateOrderStatusRequest UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&updateOrderStatusRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing update order status request: %v", err), http.StatusBadRequest)
		return
	}

	userID, err := h.Service.UpdateOrderStatus(updateOrderStatusRequest.OrderID, updateOrderStatusRequest.Status)
	if err != nil {
		http.Error(w, fmt.Sprintf("error updating order status: %v", err), http.StatusInternalServerError)
		return
	}

	orderStatusUpdate := OrderStatusUpdate{OrderID: updateOrderStatusRequest.OrderID, Status: updateOrderStatusRequest.Status, UserID: userID}
	message, err := json.Marshal(orderStatusUpdate)
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshaling order status update: %v", err), http.StatusInternalServerError)
		return
	}
	if err := h.Kafka.Producer.SendMessage("order-status-updated", string(message)); err != nil {
		http.Error(w, fmt.Sprintf("error sending order status to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
