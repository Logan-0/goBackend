// Package main provides the HTTP API layer for the Movie Review API.
// This file implements the REST API server using the chi router, including
// route definitions, request handlers, and graceful shutdown support.
//
// API Endpoints:
//   - POST   /review      - Create a new review
//   - GET    /review/{id} - Retrieve a review by ID
//   - PUT    /review/{id} - Update an existing review
//   - DELETE /review/{id} - Delete a review by ID
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

// WriteJSON is a helper function that writes a JSON response to the client.
// It sets the Content-Type header to application/json and encodes the provided
// value as JSON in the response body.
//
// Parameters:
//   - writer: The HTTP response writer
//   - status: The HTTP status code to send (e.g., http.StatusOK)
//   - anyVar: Any value that can be marshaled to JSON
//
// Returns:
//   - error: Non-nil if JSON encoding fails
//
// Note: Headers must be set before WriteHeader is called, which this function
// handles correctly.
func WriteJSON(writer http.ResponseWriter, status int, anyVar any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(anyVar)
}

// apiFunc is a function signature for HTTP handlers that can return errors.
// This allows handlers to return errors instead of handling them inline,
// which are then processed by makeHttpHandleFunc.
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents a JSON error response returned to clients.
// It provides a consistent error format across all API endpoints.
//
// Example JSON response:
//
//	{"Error": "review with id 123 not found"}
type ApiError struct {
	// Error contains the human-readable error message
	Error string
}

// makeHttpHandleFunc wraps an apiFunc to create a standard http.HandlerFunc.
// It provides centralized error handling - if the wrapped function returns
// an error, it's automatically converted to a JSON error response with
// HTTP 400 Bad Request status.
//
// This pattern allows handlers to focus on business logic and simply return
// errors, rather than handling HTTP response writing for error cases.
//
// Parameters:
//   - function: The apiFunc to wrap
//
// Returns:
//   - http.HandlerFunc: A standard handler that can be registered with the router
func makeHttpHandleFunc(function apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := function(writer, request); err != nil {
			WriteJSON(writer, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// APIServer represents the HTTP API server instance.
// It holds the server configuration, database connection, and HTTP server reference.
type APIServer struct {
	// listenAddr is the address and port the server listens on (e.g., "0.0.0.0:8080")
	listenAddr string

	// dbInstance is the storage backend implementing the Storage interface
	dbInstance Storage

	// httpServer is the underlying HTTP server for graceful shutdown support
	httpServer *http.Server
}

// RunNewServer creates, configures, and starts the HTTP API server.
// It sets up routing, configures timeouts, and implements graceful shutdown
// on SIGINT or SIGTERM signals.
//
// Parameters:
//   - listenAddr: The address to listen on (e.g., "0.0.0.0:8080")
//   - dbInstance: The storage backend for persisting reviews
//
// Returns:
//   - error: Non-nil if the server fails to start or shutdown fails
//
// The server runs until it receives an interrupt signal (Ctrl+C) or SIGTERM,
// at which point it gracefully shuts down with a 30-second timeout to allow
// in-flight requests to complete.
//
// Configured Timeouts:
//   - ReadTimeout: 15 seconds - max time to read request
//   - WriteTimeout: 15 seconds - max time to write response
//   - IdleTimeout: 60 seconds - max time for keep-alive connections
//   - ShutdownTimeout: 30 seconds - max time for graceful shutdown
func RunNewServer(listenAddr string, dbInstance Storage) error {
	// Create chi router - lightweight and fast HTTP router
	router := chi.NewRouter()
	server := &APIServer{
		listenAddr: listenAddr,
		dbInstance: dbInstance,
	}

	// Register API routes
	// All routes use makeHttpHandleFunc for consistent error handling
	router.Post("/review", makeHttpHandleFunc(server.handleCreateReview))
	router.Get("/review/{id}", makeHttpHandleFunc(server.handleGetReview))
	router.Delete("/review/{id}", makeHttpHandleFunc(server.handleDeleteReview))
	router.Put("/review/{id}", makeHttpHandleFunc(server.handleUpdateReview))

	// Configure HTTP server with security-conscious timeouts
	server.httpServer = &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,  // Prevents slow client attacks
		WriteTimeout: 15 * time.Second,  // Prevents slow response attacks
		IdleTimeout:  60 * time.Second,  // Keep-alive connection timeout
	}

	// Set up graceful shutdown signal handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in background goroutine
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("********************** Server Error: %v\n", err)
		}
	}()

	// Block until shutdown signal received
	<-shutdownChan
	fmt.Println("\n********************** Shutting down gracefully...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown - allows in-flight requests to complete
	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("********************** Server stopped")
	return nil
}

// handleGetReview handles GET /review/{id} requests.
// It retrieves a single review by its unique identifier.
//
// URL Parameters:
//   - id: The numeric ID of the review to retrieve
//
// Response:
//   - 200 OK: Returns the review as JSON
//   - 400 Bad Request: If the ID is invalid or the review is not found
//
// Example Request:
//
//	GET /review/42
//
// Example Response:
//
//	{
//	    "id": 42,
//	    "title": "Inception",
//	    "director": "Christopher Nolan",
//	    "releaseDate": "16 Jul 10 00:00",
//	    "rating": "9/10",
//	    "reviewNotes": "Mind-bending masterpiece",
//	    "dateCreated": "15 Jan 26 10:30"
//	}
func (server *APIServer) handleGetReview(writer http.ResponseWriter, request *http.Request) error {
	// Extract and validate the ID from URL path
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	// Fetch the review from the database
	review, err := server.dbInstance.GetReviewById(context.Background(), id)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}
	return WriteJSON(writer, http.StatusOK, review)
}

// handleCreateReview handles POST /review requests.
// It creates a new review from the JSON request body.
//
// Request Body:
//   - JSON object matching CreateReviewRequest structure
//
// Response:
//   - 200 OK: Returns the created review (with auto-generated ID and dateCreated)
//   - 400 Bad Request: If the request body is invalid or database insert fails
//
// Example Request:
//
//	POST /review
//	Content-Type: application/json
//
//	{
//	    "title": "Inception",
//	    "director": "Christopher Nolan",
//	    "releaseDate": "16 Jul 10 00:00 UTC",
//	    "rating": "9/10",
//	    "reviewNotes": "Mind-bending masterpiece"
//	}
func (server *APIServer) handleCreateReview(writer http.ResponseWriter, request *http.Request) error {
	// Parse the JSON request body into CreateReviewRequest
	createReviewRequest := new(CreateReviewRequest)
	if err := json.NewDecoder(request.Body).Decode(createReviewRequest); err != nil {
		return err
	}

	// Create a new Review with the provided data and auto-generated timestamps
	review := NewReview(
		createReviewRequest.Title,
		createReviewRequest.Director,
		createReviewRequest.ReleaseDate,
		createReviewRequest.Rating,
		createReviewRequest.ReviewNotes,
	)

	// Persist the review to the database
	if _, err := server.dbInstance.CreateReview(context.Background(), review); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, review)
}

// handleDeleteReview handles DELETE /review/{id} requests.
// It permanently removes a review from the database.
//
// URL Parameters:
//   - id: The numeric ID of the review to delete
//
// Response:
//   - 200 OK: Returns {"deleted": "success"}
//   - 400 Bad Request: If the ID is invalid or the review is not found
//
// Example Request:
//
//	DELETE /review/42
//
// Example Response:
//
//	{"deleted": "success"}
func (server *APIServer) handleDeleteReview(writer http.ResponseWriter, request *http.Request) error {
	// Extract and validate the ID from URL path
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	// Delete the review from the database
	if err := server.dbInstance.DeleteReview(context.Background(), id); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, map[string]string{"deleted": "success"})
}

// handleUpdateReview handles PUT /review/{id} requests.
// It updates an existing review with the provided JSON data.
// The ID in the URL takes precedence over any ID in the request body.
//
// URL Parameters:
//   - id: The numeric ID of the review to update
//
// Request Body:
//   - JSON object with review fields to update
//
// Response:
//   - 200 OK: Returns the updated review
//   - 400 Bad Request: If the ID is invalid, body is malformed, or review not found
//
// Example Request:
//
//	PUT /review/42
//	Content-Type: application/json
//
//	{
//	    "title": "Inception (Director's Cut)",
//	    "director": "Christopher Nolan",
//	    "releaseDate": "16 Jul 10 00:00",
//	    "rating": "10/10",
//	    "reviewNotes": "Even better on rewatch"
//	}
func (server *APIServer) handleUpdateReview(writer http.ResponseWriter, request *http.Request) error {
	// Extract and validate the ID from URL path
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	// Parse the JSON request body
	updateReview := new(Review)
	if err := json.NewDecoder(request.Body).Decode(updateReview); err != nil {
		return err
	}

	// Use the URL ID (overrides any ID in the body)
	updateReview.ID = id

	// Update the review in the database
	if err := server.dbInstance.UpdateReview(context.Background(), updateReview); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, updateReview)
}
