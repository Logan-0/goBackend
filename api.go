package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/review", makeHttpHandleFunc(s.handleReview))

	log.Println("Server Running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleReview(w http.ResponseWriter,r *http.Request) error {
	method := r.Method
	switch method {
		case "GET":
			return s.handleGetReview(w, r)
		case "POST": 
			return s.handleCreateReview(w, r)
		case "DELETE": 
			return s.handleDeleteReview(w, r)
		case "PUT":
			return s.handleTransportReview(w, r)
		default:
			return fmt.Errorf("method_Denied %s", r.Method)
		}
}

func (s *APIServer) handleGetReview(w http.ResponseWriter, r *http.Request) error {
	review := NewReview("Requiem for A Dream", "NewDirector", 1999, 4.5, "Great")
	return WriteJSON(w, http.StatusOK, review)
}

func (s *APIServer) handleCreateReview(w http.ResponseWriter,r *http.Request) error {
	review := NewReview("Requiem for A Dream", "NewDirector", 1999, 4.5, "Great")
	return WriteJSON(w, http.StatusOK, review)
}

func (s *APIServer) handleDeleteReview(w http.ResponseWriter,r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransportReview(w http.ResponseWriter,r *http.Request) error {
	return nil
}