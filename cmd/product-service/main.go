package main

import (
	"ecommerce-microservices/internal/product"
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
	if port == "" || len(kafkaBrokers) == 0 {
		if err := godotenv.Load("./cmd/product-service/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
		kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	}

	db, err := database.Connect("./cmd/product-service/.env")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	kafka, err := kafka.NewKafka(kafkaBrokers, "product-service")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafka.Close()

	handler := product.NewHandler(db, kafka)

	r := mux.NewRouter()
	r.HandleFunc("/get-all", handler.GetAllProducts).Methods("GET")
	r.HandleFunc("/get/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/add", handler.AddProduct).Methods("POST")

	log.Println("Product Service running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
