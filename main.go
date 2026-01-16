// Package main is the entry point for the Movie Review REST API server.
//
// This application provides a RESTful API for managing movie reviews, backed by
// a PostgreSQL database. It supports full CRUD operations on review resources.
//
// # Architecture Overview
//
// The application follows a layered architecture:
//   - main.go: Application entry point and initialization
//   - api.go: HTTP routing and request handlers (presentation layer)
//   - storage.go: Database access and persistence (data layer)
//   - types.go: Domain models and DTOs
//
// # Quick Start
//
// 1. Ensure PostgreSQL is running on localhost:5432
// 2. Set environment variables (optional - defaults work for local dev):
//   - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
//
// 3. Run the server:
//
//	go run .
//
// 4. The API will be available at http://localhost:8080
//
// # API Endpoints
//
//	POST   /review      - Create a new review
//	GET    /review/{id} - Get a review by ID
//	PUT    /review/{id} - Update a review
//	DELETE /review/{id} - Delete a review
//
// # Graceful Shutdown
//
// The server handles SIGINT (Ctrl+C) and SIGTERM signals gracefully,
// allowing in-flight requests to complete before shutting down.
package main

import (
	"fmt"
	"log"
)

// main is the application entry point.
// It initializes the database connection and starts the HTTP server.
//
// Initialization sequence:
//  1. Connect to PostgreSQL database
//  2. Configure connection pool
//  3. Prepare SQL statements
//  4. Start HTTP server with graceful shutdown support
//
// The function will log.Fatal and exit if database connection fails.
// Once the server is running, it blocks until a shutdown signal is received.
func main() {
	// Initialize Database connection and prepare statements
	client, err := InitializeClientAndDB()
	if err != nil {
		log.Fatal("********************** Failed: Connection to Database ", err.Error())
	}
	fmt.Println("********************** Success: Database Port 5432")

	// Table creation is commented out for development.
	// In production, tables should be managed via migrations.
	// Uncomment below to auto-create the reviews table on startup:
	//
	// err = client.CreateReviewTable()
	// if err != nil {
	// 	log.Fatal("********************** Failed: Create Review Table" + err.Error())
	// }
	// fmt.Println("********************** Success: Create Review Table")

	// Start the HTTP server (blocks until shutdown signal)
	fmt.Println("********************** Success: Server Running 8080")
	RunNewServer("0.0.0.0:8080", client)
}
