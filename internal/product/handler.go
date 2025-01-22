package product

import (
	"database/sql"
	"ecommerce-microservices/pkg/kafka"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *Service
	Kafka   *kafka.Kafka
}

func NewHandler(db *sql.DB, kafka *kafka.Kafka) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{Service: service, Kafka: kafka}
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.GetAllProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting products: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing product ID: %v", err), http.StatusBadRequest)
		return
	}

	product, err := h.Service.GetProductByID(productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting product: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var productRequest NewProductRequest
	if err := json.NewDecoder(r.Body).Decode(&productRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing product request: %v", err), http.StatusBadRequest)
		return
	}

	product, err := h.Service.AddProduct(&productRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error adding product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
