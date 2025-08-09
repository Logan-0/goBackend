package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteJSON(writer http.ResponseWriter, status int, anyVar any) error {
	writer.WriteHeader(status)
	writer.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(writer).Encode(anyVar)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(function apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := function(writer, request); err != nil {
			WriteJSON(writer, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	dbInstance Storage
}

func RunNewServer(listenAddr string, dbInstance Storage) error {

	// Create Router and Server listenAddr -> Port
	router := mux.NewRouter()
	server := &APIServer{
		listenAddr: listenAddr,
		dbInstance: dbInstance,
	}

	// Single endpoint with different methods
	router.HandleFunc("/review", makeHttpHandleFunc(server.handleCreateReview)).Methods("POST")
	router.HandleFunc("/review/{id}", makeHttpHandleFunc(server.handleGetReview)).Methods("GET")
	router.HandleFunc("/review/{id}", makeHttpHandleFunc(server.handleDeleteReview)).Methods("DELETE")
	router.HandleFunc("/review/{id}", makeHttpHandleFunc(server.handleUpdateReview)).Methods("PUT")

	// Always returns non-nil error
	return http.ListenAndServe(server.listenAddr, router)
}

func (server *APIServer) handleGetReview(writer http.ResponseWriter, request *http.Request) error {

	id, err := strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	review, err := server.dbInstance.GetReviewById(context.Background(), id)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}
	return WriteJSON(writer, http.StatusOK, review)
}

func (server *APIServer) handleCreateReview(writer http.ResponseWriter, request *http.Request) error {

	createReviewRequest := new(CreateReviewRequest)

	if err := json.NewDecoder(request.Body).Decode(createReviewRequest); err != nil {
		return err
	}

	review := NewReview(
		createReviewRequest.Title,
		createReviewRequest.Director,
		createReviewRequest.ReleaseDate,
		createReviewRequest.Rating,
		createReviewRequest.ReviewNotes,
	)

	if _, err := server.dbInstance.CreateReview(context.Background(), review); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, review)
}

func (server *APIServer) handleDeleteReview(writer http.ResponseWriter, request *http.Request) error {

	id, err := strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	if err := server.dbInstance.DeleteReview(context.Background(), id); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, map[string]string{"deleted": "success"})
}

// Rename handleTransportReview to handleUpdateReview
func (server *APIServer) handleUpdateReview(writer http.ResponseWriter, request *http.Request) error {

	id, err := strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	updateReview := new(Review)

	if err := json.NewDecoder(request.Body).Decode(updateReview); err != nil {
		return err
	}

	updateReview.ID = id

	if err := server.dbInstance.UpdateReview(context.Background(), updateReview); err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, updateReview)
}
