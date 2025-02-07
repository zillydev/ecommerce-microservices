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

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *graph.CreateOrderInput) (*graph.Order, error) {
	userID := middlewares.ForJWTContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/create?userId=%s", r.OrderServiceURL, userID), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create order: %v", string(body))
	}

	var order graph.Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &order, nil
}

// UpdateOrderStatus is the resolver for the updateOrderStatus field.
func (r *mutationResolver) UpdateOrderStatus(ctx context.Context, input *graph.UpdateOrderStatusInput) (bool, error) {
	adminKey := middlewares.ForAdminContext(ctx)
	if adminKey == "" {
		return false, fmt.Errorf("unauthorized")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return false, fmt.Errorf("failed to marshal input: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/update", r.OrderServiceURL), bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to update order status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to update order status: %v", string(body))
	}

	return true, nil
}

// GetAllOrdersByUser is the resolver for the getAllOrdersByUser field.
func (r *queryResolver) GetAllOrdersByUser(ctx context.Context) ([]*graph.Order, error) {
	userID := middlewares.ForJWTContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	resp, err := http.Get(fmt.Sprintf("%s/get-by-user?userId=%s", r.OrderServiceURL, userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get orders: %v", string(body))
	}

	var orders []*graph.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return orders, nil
}

// GetOrder is the resolver for the getOrder field.
func (r *queryResolver) GetOrder(ctx context.Context, orderID int) (*graph.Order, error) {
	userID := middlewares.ForJWTContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	resp, err := http.Get(fmt.Sprintf("%s/get/%d", r.OrderServiceURL, orderID))
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get order: %v", string(body))
	}

	var order graph.Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &order, nil
}
