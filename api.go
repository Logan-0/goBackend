package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (server *APIServer) handleReview(writer http.ResponseWriter,request *http.Request) error {
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
	id := mux.Vars(request)["id"]
	fmt.Println("Review Id:", id)
	// TODO: db.get (id)

	// if err != null {
	// 	log.Println("Error returning Review Id:", id)
	// } else {
	// log.Println("Review Id Returned:", id)
	// }
	return WriteJSON(writer, http.StatusOK, &Review{})
}

func (server *APIServer) handleCreateReview(writer http.ResponseWriter,request *http.Request) error {
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
	if err := server.dbInstance.CreateReview(review); err != nil {
		return err;
	}
	return WriteJSON(writer, http.StatusOK, review)
}

func (server *APIServer) handleDeleteReview(writer http.ResponseWriter,request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, request)
}

func (server *APIServer) handleTransportReview(writer http.ResponseWriter,request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, request)
}