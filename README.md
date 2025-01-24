# E-Commerce Notification System

This project is a personalized notification system for an e-commerce platform, designed as a microservices-based architecture. It incorporates various technologies to ensure scalability, reliability, and flexibility.

## Tech Stack

- **Programming Language**: Golang
- **Databases**: PostgreSQL, Redis
- **Database Migrations**: golang-migrate
- **Message Broker**: Kafka
- **Containerization**: Docker & Docker Compose
- **API Gateway**: GraphQL (gqlgen)

## Microservices Overview

### User Service

- **Features**:
    - Register a user (name, email, preferences)
    - Get user details
    - Update notification preferences
- **Database**: PostgreSQL
- **Authentication**: JWT

### Notification Service

- **Features**:
    - Post, mark read, and fetch unread notifications
    - Scheduled promotion notifications (cron job) for users with relevant preferences
    - Subscribes to the `order-status-updated` Kafka topic to create notifications for users with the `order_updates` preference
- **Database**: PostgreSQL

### Product Service

- **Features**:
    - Add product
    - Get product details
    - Fetch all products
- **Database**: PostgreSQL

### Order Service

- **Features**:
    - Create order
    - Fetch order details and orders by user
    - Update order status
    - Emits `order-status-updated` events via Kafka when order status is updated
- **Database**: PostgreSQL

### Gateway Service

- **Features**:
    - Unified GraphQL API for client interaction
    - JWT middleware for protected routes
    - Admin middleware for restricted actions
    - Redis caching for frequently queried data (e.g., product listings)
    - Data aggregation from all microservices
- **GraphQL Framework**: gqlgen
- **Cache**: Redis

## Postman Collection

https://www.postman.com/blue-escape-5551/workspace/public-workspace/collection/678ff7e709d730b05b8d4e35?action=share&creator=35043396

Queries marked with "JWT" require a Bearer token in Authorization header, which can be received after registering a user.

Queries marked with "Admin" require `x-api-key` in headers, which should match `ADMIN_SECRET_KEY` in Gateway service .env file.

## Running with Docker

The entire architecture can be run with a single command:

```sh
docker-compose up
```

## Local development

Use the Makefile to build and run all services at once:

```sh
make build
make run
```

Or run an individual service:

```sh
go run cmd/user-service/main.go
```
