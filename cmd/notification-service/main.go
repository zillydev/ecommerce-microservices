package main

import (
	"ecommerce-microservices/internal/notification"
	"ecommerce-microservices/pkg/database"
	"ecommerce-microservices/pkg/kafka"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	port := os.Getenv("PORT")
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if port == "" || len(kafkaBrokers) == 0 || userServiceURL == "" {
		if err := godotenv.Load("./cmd/notification-service/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
		kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
		userServiceURL = os.Getenv("USER_SERVICE_URL")
	}

	db, err := database.Connect("./cmd/notification-service/.env")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	kafka, err := kafka.NewKafka(kafkaBrokers, "notification-service")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafka.Close()

	handler := notification.NewHandler(db, kafka, userServiceURL)

	r := mux.NewRouter()
	r.HandleFunc("/post", handler.PostNotification).Methods("POST")
	r.HandleFunc("/mark-read", handler.MarkRead).Methods("PUT")
	r.HandleFunc("/unread", handler.GetUnreadNotifications).Methods("GET")

	err = kafka.Consumer.Subscribe("order-status-updated", func(message *sarama.ConsumerMessage) error {
		log.Printf("Received message: %s", message.Value)

		var orderStatusUpdate notification.OrderStatusUpdate
		if err := json.Unmarshal(message.Value, &orderStatusUpdate); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}

		err = handler.HandleOrderStatusUpdate(orderStatusUpdate)
		if err != nil {
			return fmt.Errorf("failed to handle order status update: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	c := cron.New()

	// Daily promotion notifications, at 9:00 AM every day
	_, err = c.AddFunc("0 9 * * *", func() {
		fmt.Printf("Running cron job at: %v\n", time.Now().Format("15:04:05"))
		handler.SendDailyPromotionNotifications()
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	log.Println("Notification Service running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
