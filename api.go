package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
}

func RunNewServer(listenAddr string) {

	// Create Router and Server listenAddr -> Port
	router := mux.NewRouter()
	server := &APIServer{
		listenAddr: listenAddr,
	}
	
	router.HandleFunc("/review", makeHttpHandleFunc(server.handleReview))
	router.HandleFunc("/review/:id", makeHttpHandleFunc(server.handleGetReview))

	err := http.ListenAndServe(server.listenAddr, router)
	if err != nil {
		log.Fatal("Failed To Listen and Serve On Port: " + listenAddr)
	}
	fmt.Println("Server Running on Port: ", server.listenAddr)
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
			return fmt.Errorf("method_Denied %server", request.Method)
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
	dateLayout := "2006-01-02 15:04:05"
	vars := mux.Vars(request)
	title := vars["title"]
	director := vars["director"]
	releaseDateAsString := vars["releaseDate"]
	rating := vars["rating"]
	releaseDateAsDate, err := time.Parse(dateLayout, releaseDateAsString)
	if err != nil {
		fmt.Println("releaseDateInvalid: ",err)
    }
	reviewNotes := vars["reviewNotes"]
	review := NewReview(title, director, releaseDateAsDate.String(), rating, reviewNotes)
	return WriteJSON(writer, http.StatusOK, review)
}

func (server *APIServer) handleDeleteReview(writer http.ResponseWriter,request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, request)
}

func (server *APIServer) handleTransportReview(writer http.ResponseWriter,request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, request)
}