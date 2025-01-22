package main

import (
	"ecommerce-microservices/internal/gateway/graph"
	"ecommerce-microservices/internal/gateway/graph/resolvers"
	"ecommerce-microservices/internal/gateway/middlewares"
	"ecommerce-microservices/pkg/kafka"
	"ecommerce-microservices/pkg/redis"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	port := os.Getenv("PORT")
	redisURL := os.Getenv("REDIS_URL")
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	kakfaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	adminSecretKey := os.Getenv("ADMIN_SECRET_KEY")

	if port == "" || redisURL == "" || userServiceURL == "" || notificationServiceURL == "" || len(kakfaBrokers) == 0 || productServiceURL == "" || orderServiceURL == "" || adminSecretKey == "" {
		if err := godotenv.Load("./cmd/gateway-service/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
		redisURL = os.Getenv("REDIS_URL")
		userServiceURL = os.Getenv("USER_SERVICE_URL")
		notificationServiceURL = os.Getenv("NOTIFICATION_SERVICE_URL")
		productServiceURL = os.Getenv("PRODUCT_SERVICE_URL")
		orderServiceURL = os.Getenv("ORDER_SERVICE_URL")
		kakfaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
		jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
		adminSecretKey = os.Getenv("ADMIN_SECRET_KEY")
	}

	kafka, err := kafka.NewKafka(kakfaBrokers, "gateway-service")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer kafka.Close()

	redisClient := redis.NewRedisClient(redisURL)
	defer redisClient.Close()

	resolver := &resolvers.Resolver{
		UserServiceURL:         userServiceURL,
		NotificationServiceURL: notificationServiceURL,
		ProductServiceURL:      productServiceURL,
		OrderServiceURL:        orderServiceURL,
		JWTSecretKey:           jwtSecretKey,
		RedisClient:            redisClient,
	}

	r := mux.NewRouter()
	r.Use(middlewares.JWTAuthMiddleware(jwtSecretKey))
	r.Use(middlewares.AdminAuthMiddleware(adminSecretKey))

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Println("Gateway Service running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
