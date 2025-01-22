package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"bytes"
	"context"
	"ecommerce-microservices/internal/gateway/graph"
	"ecommerce-microservices/internal/gateway/middlewares"
	"ecommerce-microservices/pkg/jwt"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input *graph.RegisterUserInput) (*graph.RegisterUserResult, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/register", r.UserServiceURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s", string(body))
	}

	var user graph.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %v", err)
	}

	token, err := jwt.GenerateToken(fmt.Sprintf("%d", user.ID), r.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	result := &graph.RegisterUserResult{
		User:        &user,
		AccessToken: token,
	}
	return result, nil
}

// UpdatePreferences is the resolver for the updatePreferences field.
func (r *mutationResolver) UpdatePreferences(ctx context.Context, input *graph.UpdatePreferencesInput) (bool, error) {
	userid := middlewares.ForJWTContext(ctx)
	if userid == "" {
		return false, fmt.Errorf("unauthorized")
	}

	body, _ := json.Marshal(input)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/update-preferences/%s", r.UserServiceURL, userid), bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("%s", string(body))
	}

	return true, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*graph.User, error) {
	userid := middlewares.ForJWTContext(ctx)
	if userid == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	resp, err := http.Get(fmt.Sprintf("%s/user/%s", r.UserServiceURL, userid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s", string(body))
	}

	var user graph.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %v", err)
	}
	return &user, nil
}
