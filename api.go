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

	router.HandleFunc("/review", makeHttpHandleFunc(server.handleReview))
	router.HandleFunc("/review/:id", makeHttpHandleFunc(server.handleGetReview))
	router.HandleFunc("/createReview", makeHttpHandleFunc(server.handleCreateReview))

	// Always returns non-nil error
	return http.ListenAndServe(server.listenAddr, router)
}

func (server *APIServer) handleReview(writer http.ResponseWriter, request *http.Request) error {
	method := request.Method
	switch method {
	case "GET":
		return server.handleGetReview(writer, request)
	case "POST":
		return server.handleCreateReview(writer, request)
	case "DELETE":
		return server.handleDeleteReview(writer, request)
	case "PUT":
		return server.handleTransportReview(writer, request)
	default:
		return fmt.Errorf(`method denied %s failed`, request.Method)
	}
}

func (server *APIServer) handleGetReview(writer http.ResponseWriter, request *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	review, err := server.dbInstance.GetReviewById(context.Background(), id)
	if err != nil {
		return err
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

	if err := server.dbInstance.CreateReview(context.Background(), review); err != nil {
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

func (server *APIServer) handleTransportReview(writer http.ResponseWriter, request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, request)
}
