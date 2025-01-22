package main

import (
	"ecommerce-microservices/internal/user"
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
		if err := godotenv.Load("./cmd/user-service/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
		kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	}

	db, err := database.Connect("./cmd/user-service/.env")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	kafka, err := kafka.NewKafka(kafkaBrokers, "user-service")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafka.Close()

	handler := user.NewHandler(db, kafka)

	r := mux.NewRouter()
	r.HandleFunc("/register", handler.RegisterUser).Methods("POST")
	r.HandleFunc("/update-preferences/{id}", handler.UpdatePreferences).Methods("PUT")
	r.HandleFunc("/user/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/get-users-by-preference", handler.GetUsersByPreferences).Methods("GET")

	log.Println("User Service running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
