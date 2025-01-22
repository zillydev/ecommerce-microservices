package main

import (
	"ecommerce-microservices/internal/order"
	"ecommerce-microservices/pkg/database"
	"ecommerce-microservices/pkg/kafka"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	port := os.Getenv("PORT")
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if port == "" || len(kafkaBrokers) == 0 || productServiceURL == "" {
		if err := godotenv.Load("./cmd/order-service/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
		kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
		productServiceURL = os.Getenv("PRODUCT_SERVICE_URL")
	}

	db, err := database.Connect("./cmd/order-service/.env")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	kafka, err := kafka.NewKafka(kafkaBrokers, "order-service")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafka.Close()

	handler := order.NewHandler(db, kafka, productServiceURL)

	r := mux.NewRouter()
	r.HandleFunc("/get-by-user", handler.GetAllOrdersByUser).Methods("GET")
	r.HandleFunc("/get/{id}", handler.GetOrder).Methods("GET")
	r.HandleFunc("/create", handler.CreateOrder).Methods("POST")
	r.HandleFunc("/update", handler.UpdateOrderStatus).Methods("PUT")

	log.Println("Order Service running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
