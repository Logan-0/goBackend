package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateReview(context.Context, *Review) (string, error)
	UpdateReview(context.Context, *Review) error
	DeleteReview(context.Context, int) error
	GetReviewById(context.Context, int) (*Review, error)
}

type PgDb struct {
	db *sql.DB
}

const (
	port     = 5432
	host     = "localhost"
	user     = "postgres"
	password = "test"
	dbname   = "postgres"

	// Connection pool settings
	maxOpenConns    = 25
	maxIdleConns    = 5
	maxConnLifetime = 5 * time.Minute

	// Context timeout settings
	defaultTimeout = 10 * time.Second
)

func InitializeClientAndDB() (*PgDb, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("********************** Failed: Open Connection to PostgreSQL DB: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxConnLifetime)

	// Verify the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("********************** Failed: Ping DB: %w", err)
	}
	return &PgDb{db: db}, nil
}

// To be used after testing when Database when restarting dbs is no longer as simple.

// func (pg *PgDb) CreateReviewTable() error {
// 	createTableQuery := `CREATE TABLE IF NOT EXISTS public.reviews (
// 	id SERIAL PRIMARY KEY,
// 	title VARCHAR NOT NULL,
// 	director VARCHAR NOT NULL,
// 	rating VARCHAR NOT NULL,
// 	releaseDate VARCHAR NOT NULL,
// 	reviewNotes VARCHAR NOT NULL,
// 	createdAt VARCHAR NOT NULL
// 	);`

// 	_, err := pg.db.Exec(createTableQuery)
// 	if err != nil {
// 		return fmt.Errorf("********************** Failed: Create Review Table: %w",err)
// 	}

// 	return nil
// }

func (pg *PgDb) DropReviewTable() error {
	dropTableQuery := `DROP TABLE reviews`
	_, err := pg.db.Query(dropTableQuery)
	if err != nil {
		fmt.Println("********************** Failed: Drop Review Table")
		return err
	}
	fmt.Println("********************** Success: Drop Review Table")
	return nil
}

func (pg *PgDb) CreateReview(ctx context.Context, review *Review) (string, error) {
	createReviewQuery := `INSERT INTO public.reviews (
	title,director,releaseDate,rating,reviewNotes,dateCreated
	) VALUES ($1, $2, $3, $4, $5, $6);`

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := pg.db.ExecContext(ctx, createReviewQuery,
		review.Title,
		review.Director,
		review.ReleaseDate,
		review.Rating,
		review.ReviewNotes,
		review.DateCreated)

	if err != nil {
		return "", fmt.Errorf("failed to create review: %w", err)
	}

	success := "Review Created :: Recorded In DB:: " + review.DateCreated
	return success, nil
}

func (pg *PgDb) UpdateReview(ctx context.Context, review *Review) error {
	updateReviewQuery := `UPDATE public.reviews 
        SET title=$1, director=$2, releaseDate=$3, rating=$4, reviewNotes=$5 
        WHERE id=$6`

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := pg.db.ExecContext(ctx, updateReviewQuery,
		review.Title,
		review.Director,
		review.ReleaseDate,
		review.Rating,
		review.ReviewNotes,
		review.ID)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}
	return nil
}

func (pg *PgDb) DeleteReview(ctx context.Context, id int) error {
	deleteReviewQuery := `DELETE FROM public.reviews WHERE id=$1`

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := pg.db.ExecContext(ctx, deleteReviewQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	fmt.Println("********************** Success: Deleted Review ", id)
	return nil
}

func (pg *PgDb) GetReviewById(ctx context.Context, id int) (*Review, error) {
	getReviewQuery := `SELECT id, title, director, releaseDate, rating, reviewNotes, dateCreated 
		FROM public.reviews WHERE id=$1`

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	review := &Review{}
	err := pg.db.QueryRowContext(ctx, getReviewQuery, id).Scan(
		&review.ID,
		&review.Title,
		&review.Director,
		&review.ReleaseDate,
		&review.Rating,
		&review.ReviewNotes,
		&review.DateCreated)

	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return review, nil
}
