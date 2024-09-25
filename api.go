package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
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
	storage Storage
}

func NewAPIServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage: storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	
	router.HandleFunc("/review", makeHttpHandleFunc(s.handleReview))
	router.HandleFunc("/review/:id", makeHttpHandleFunc(s.handleGetReview))

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
	id := mux.Vars(r)["id"]
	fmt.Println("Review Id:", id)
	// db.get (id)
	return WriteJSON(w, http.StatusOK, &Review{})
}

func (s *APIServer) handleCreateReview(w http.ResponseWriter,r *http.Request) error {
	dateLayout := "2006-01-02 15:04:05"
	title := mux.Vars(r)["title"]
	director := mux.Vars(r)["director"]
	releaseDateAsString := mux.Vars(r)["releaseDate"]
	rating := mux.Vars(r)["rating"]
	releaseDateAsDate, err := time.Parse(dateLayout, releaseDateAsString)
	if err != nil {
		fmt.Println("releaseDateInvalid: ",err)
    }
	reviewNotes := mux.Vars(r)["reviewNotes"]
	review := NewReview(title, director, releaseDateAsDate.String(), rating, reviewNotes)
	return WriteJSON(w, http.StatusOK, review)
}

func (s *APIServer) handleDeleteReview(w http.ResponseWriter,r *http.Request) error {
	return WriteJSON(w, http.StatusOK, r)
}

func (s *APIServer) handleTransportReview(w http.ResponseWriter,r *http.Request) error {
	return WriteJSON(w, http.StatusOK, r)
}