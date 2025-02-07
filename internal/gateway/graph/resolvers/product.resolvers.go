package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"bytes"
	"context"
	"ecommerce-microservices/internal/gateway/graph"
	"ecommerce-microservices/internal/gateway/middlewares"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AddProduct is the resolver for the addProduct field.
func (r *mutationResolver) AddProduct(ctx context.Context, input *graph.AddProductInput) (*graph.Product, error) {
	adminKey := middlewares.ForAdminContext(ctx)
	if adminKey == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product request: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/add", r.ProductServiceURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to add product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add product: %v", string(body))
	}

	var product graph.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %v", err)
	}

	cacheKey := "product_listings"
	// Invalidate the cache
	if err := r.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		return nil, fmt.Errorf("failed to invalidate cache: %v", err)
	}

	return &product, nil
}

// GetAllProducts is the resolver for the getAllProducts field.
func (r *queryResolver) GetAllProducts(ctx context.Context) ([]*graph.Product, error) {
	cacheKey := "product_listings"
	// Check if the data is already cached
	cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []*graph.Product
		if err := json.Unmarshal([]byte(cachedData), &products); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		fmt.Println("Cache hit")
		return products, nil
	}

	resp, err := http.Get(fmt.Sprintf("%s/get-all", r.ProductServiceURL))
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get products: %s", string(body))
	}

	var products []*graph.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode products response: %v", err)
	}

	// Cache the data
	data, err := json.Marshal(products)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal products: %v", err)
	}
	if err := r.RedisClient.Set(ctx, cacheKey, string(data), 0).Err(); err != nil {
		return nil, fmt.Errorf("failed to cache products: %v", err)
	}

	return products, nil
}

// GetProduct is the resolver for the getProduct field.
func (r *queryResolver) GetProduct(ctx context.Context, productID int) (*graph.Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/get/%d", r.ProductServiceURL, productID))
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get product: %s", string(body))
	}

	var product graph.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %v", err)
	}

	return &product, nil
}
