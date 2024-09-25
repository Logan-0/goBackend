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

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connectString := "user=postgres password=postgres dbname=review_db sslmode=disable"
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db:db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.createReviewTable()
}

func (s *PostgresStorage) createReviewTable() error {
	query := `CREATE table reviews if not exists (
		id serial primary key,
		title varchar(75),
		director varchar(75),
		rating serial,
		release_date,
		review varchar(20000),
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) dropReviewTable() error {
	query := `DROP table reviews`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateReview(*Review) error {
	return nil
}

func (s *PostgresStorage) UpdateReview(*Review) error {
	return nil
}
func (s *PostgresStorage) DeleteReview(id int) error {
	return nil
}
func (s *PostgresStorage) GetReviewById(id int) (*Review, error){
	return nil, nil
}