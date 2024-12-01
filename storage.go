package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateReview(*Review) error
	UpdateReview(*Review) error
	DeleteReview(int) error
	GetReviewById(int) (*Review, error)
}

type PgDb struct {
	db *sql.DB
}

const (
	port = 5432
	host = "localhost"
	user = "postgres"
	password = "test"
	dbname = "reviewdb"
	)

func InitializeClientAndDB() (*PgDb, error) {
	// Connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("********************** Failed: Open Connection to PostgreSQL DB: %w", err)
	}
	
	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("********************** Failed: Ping DB: %w",err)
	}
    return &PgDb{db: db,}, nil
}

func (pg *PgDb) CreateReviewTable() error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS public.reviewtable (
	id SERIAL PRIMARY KEY,
	title VARCHAR NOT NULL,
	director VARCHAR NOT NULL,
	rating VARCHAR NOT NULL,
	releaseDate VARCHAR NOT NULL,
	reviewNotes VARCHAR NOT NULL,
	createdAt VARCHAR NOT NULL
	);`

	_, err := pg.db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("********************** Failed: Create Review Table: %w",err)
	}

	return nil
}

func (pg *PgDb) DropReviewTable() error {
	dropTableQuery := `DROP TABLE reviewdb`
	_, err := pg.db.Query(dropTableQuery)
	if err != nil {
		fmt.Println("********************** Failed: Drop Review Table")
		return err
	}
	fmt.Println("********************** Success: Drop Review Table")
	return nil
}

func (pg *PgDb) CreateReview(review *Review) error {
	createReviewQuery := `INSERT INTO public.reviewtable (
	title,director,releaseDate,rating,reviewNotes,createdAt
	) VALUES ($1, $2, $3, $4, $5, $6);`

	response, err := pg.db.Query(createReviewQuery,
		review.Title,
		review.Director,
		review.ReleaseDate,
		review.Rating,
		review.ReviewNotes,
		review.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", response)
	return nil
}

func (pg *PgDb) UpdateReview(*Review) error {
	return nil
}
func (pg *PgDb) DeleteReview(id int) error {
	return nil
}
func (pg *PgDb) GetReviewById(id int) (*Review, error){
	return nil, nil
}