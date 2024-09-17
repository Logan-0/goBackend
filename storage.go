package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateReview(*Review) error
	UpdateReview(*Review) error
	DeleteReview(int) error
	GetReviewById(int) (*Review, error)
}

type PostgreStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgreStorage, error) {
	connectString := ""
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreStorage{
		db:db,
	}, nil
}


func (s *PostgreStorage) CreateReview(*Review) error {
	return nil
}

func (s *PostgreStorage) UpdateReview(*Review) error {
	return nil
}
func (s *PostgreStorage) DeleteReview(id int) error {
	return nil
}
func (s *PostgreStorage) GetReviewById(id int) (*Review, error){
	return nil, nil
}