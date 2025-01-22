package resolvers

import "github.com/redis/go-redis/v9"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserServiceURL         string
	NotificationServiceURL string
	ProductServiceURL      string
	OrderServiceURL        string
	JWTSecretKey           string
	RedisClient            *redis.Client
}
