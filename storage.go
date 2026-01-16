// Package main provides the storage layer for the Movie Review API.
// This file implements the database access layer using PostgreSQL with
// connection pooling, prepared statements, and context-based timeouts.
//
// The storage layer follows the Repository pattern, abstracting database
// operations behind the Storage interface for testability and flexibility.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	// PostgreSQL driver - imported for side effects (driver registration)
	_ "github.com/lib/pq"
)

// Storage defines the interface for review persistence operations.
// This interface allows for easy mocking in tests and potential
// swapping of storage backends (e.g., switching from PostgreSQL to MySQL).
//
// All methods accept a context.Context for cancellation and timeout support.
type Storage interface {
	// CreateReview persists a new review to the database.
	// Returns a success message with the creation timestamp, or an error.
	CreateReview(context.Context, *Review) (string, error)

	// UpdateReview modifies an existing review identified by the Review.ID field.
	// Returns an error if the review doesn't exist or the update fails.
	UpdateReview(context.Context, *Review) error

	// DeleteReview removes a review by its ID.
	// Returns an error if the review doesn't exist or the deletion fails.
	DeleteReview(context.Context, int) error

	// GetReviewById retrieves a single review by its unique identifier.
	// Returns the Review and nil error on success, or nil and an error if not found.
	GetReviewById(context.Context, int) (*Review, error)
}

// PgDb implements the Storage interface using PostgreSQL.
// It maintains a connection pool and prepared statements for optimal performance.
//
// Prepared statements are created once during initialization and reused for all
// subsequent queries, reducing parsing overhead and improving throughput.
type PgDb struct {
	// db is the underlying database connection pool managed by database/sql.
	db *sql.DB

	// stmtCreate is the prepared statement for INSERT operations.
	stmtCreate *sql.Stmt

	// stmtUpdate is the prepared statement for UPDATE operations.
	stmtUpdate *sql.Stmt

	// stmtDelete is the prepared statement for DELETE operations.
	stmtDelete *sql.Stmt

	// stmtGetById is the prepared statement for SELECT by ID operations.
	stmtGetById *sql.Stmt
}

const (
	// maxOpenConns is the maximum number of open connections to the database.
	// This limits resource usage and prevents overwhelming the database server.
	maxOpenConns = 25

	// maxIdleConns is the maximum number of idle connections retained in the pool.
	// Setting this equal to maxOpenConns prevents connection churn under load.
	maxIdleConns = 25

	// maxConnLifetime is the maximum duration a connection can be reused.
	// This helps balance connection freshness with connection reuse efficiency.
	maxConnLifetime = 5 * time.Minute

	// defaultTimeout is the maximum duration for database operations.
	// Operations exceeding this timeout will be cancelled and return an error.
	defaultTimeout = 10 * time.Second
)

// getEnv retrieves an environment variable value or returns a fallback default.
// This is used for configuration management, allowing runtime configuration
// without code changes.
//
// Parameters:
//   - key: The environment variable name to look up
//   - fallback: The default value to return if the variable is not set
//
// Returns:
//   - The environment variable value if set, otherwise the fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvInt retrieves an environment variable as an integer or returns a fallback.
// If the environment variable is set but cannot be parsed as an integer,
// the fallback value is returned.
//
// Parameters:
//   - key: The environment variable name to look up
//   - fallback: The default integer value if the variable is not set or invalid
//
// Returns:
//   - The parsed integer value if valid, otherwise the fallback value
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// InitializeClientAndDB creates and configures a new PostgreSQL database connection.
// It reads connection parameters from environment variables with sensible defaults
// for local development.
//
// Environment Variables:
//   - DB_HOST: Database host (default: "localhost")
//   - DB_PORT: Database port (default: 5432)
//   - DB_USER: Database user (default: "postgres")
//   - DB_PASSWORD: Database password (default: "test")
//   - DB_NAME: Database name (default: "postgres")
//
// The function performs the following initialization steps:
//  1. Builds connection string from environment variables
//  2. Opens database connection pool
//  3. Configures connection pool settings (max connections, idle connections, lifetime)
//  4. Verifies connectivity with a ping
//  5. Prepares SQL statements for CRUD operations
//
// Returns:
//   - *PgDb: Configured database client ready for use
//   - error: Non-nil if any initialization step fails
//
// Example:
//
//	client, err := InitializeClientAndDB()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
func InitializeClientAndDB() (*PgDb, error) {
	// Read database configuration from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnvInt("DB_PORT", 5432)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "test")
	dbname := getEnv("DB_NAME", "postgres")

	// Build PostgreSQL connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open database connection pool (does not actually connect yet)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("********************** Failed: Open Connection to PostgreSQL DB: %w", err)
	}

	// Configure connection pool for optimal performance
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxConnLifetime)

	// Verify the connection is actually working
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("********************** Failed: Ping DB: %w", err)
	}

	pgDb := &PgDb{db: db}

	// Prepare statements for better performance (parsed once, executed many times)
	if err := pgDb.prepareStatements(); err != nil {
		db.Close()
		return nil, fmt.Errorf("********************** Failed: Prepare Statements: %w", err)
	}

	return pgDb, nil
}

// prepareStatements creates prepared statements for all CRUD operations.
// Prepared statements are parsed and planned once by PostgreSQL, then reused
// for subsequent executions, providing significant performance benefits.
//
// This is called automatically by InitializeClientAndDB and should not be
// called directly.
//
// Returns:
//   - error: Non-nil if any statement preparation fails
func (pg *PgDb) prepareStatements() error {
	var err error

	// Prepare INSERT statement for creating new reviews
	pg.stmtCreate, err = pg.db.Prepare(`INSERT INTO public.reviews (
		title,director,releaseDate,rating,reviewNotes,dateCreated
	) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return fmt.Errorf("prepare create: %w", err)
	}

	// Prepare UPDATE statement for modifying existing reviews
	pg.stmtUpdate, err = pg.db.Prepare(`UPDATE public.reviews 
		SET title=$1, director=$2, releaseDate=$3, rating=$4, reviewNotes=$5 
		WHERE id=$6`)
	if err != nil {
		return fmt.Errorf("prepare update: %w", err)
	}

	// Prepare DELETE statement for removing reviews
	pg.stmtDelete, err = pg.db.Prepare(`DELETE FROM public.reviews WHERE id=$1`)
	if err != nil {
		return fmt.Errorf("prepare delete: %w", err)
	}

	// Prepare SELECT statement for fetching reviews by ID
	pg.stmtGetById, err = pg.db.Prepare(`SELECT id, title, director, releaseDate, rating, reviewNotes, dateCreated 
		FROM public.reviews WHERE id=$1`)
	if err != nil {
		return fmt.Errorf("prepare getById: %w", err)
	}

	return nil
}

// Close releases all database resources including prepared statements and
// the connection pool. This should be called when the application shuts down
// to ensure clean resource cleanup.
//
// It is safe to call Close multiple times; subsequent calls are no-ops for
// already-closed statements.
//
// Returns:
//   - error: Non-nil if closing the database connection fails
//
// Example:
//
//	client, _ := InitializeClientAndDB()
//	defer client.Close() // Ensure cleanup on exit
func (pg *PgDb) Close() error {
	// Close all prepared statements first
	if pg.stmtCreate != nil {
		pg.stmtCreate.Close()
	}
	if pg.stmtUpdate != nil {
		pg.stmtUpdate.Close()
	}
	if pg.stmtDelete != nil {
		pg.stmtDelete.Close()
	}
	if pg.stmtGetById != nil {
		pg.stmtGetById.Close()
	}
	// Close the underlying database connection pool
	return pg.db.Close()
}

// CreateReviewTable creates the reviews table if it doesn't exist.
// This is commented out for development but can be enabled for production
// deployments where manual table creation is not feasible.
//
// Table Schema:
//   - id: SERIAL PRIMARY KEY (auto-incrementing integer)
//   - title: VARCHAR NOT NULL
//   - director: VARCHAR NOT NULL
//   - rating: VARCHAR NOT NULL
//   - releaseDate: VARCHAR NOT NULL
//   - reviewNotes: VARCHAR NOT NULL
//   - createdAt: VARCHAR NOT NULL

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

// DropReviewTable removes the reviews table from the database.
// WARNING: This is a destructive operation that permanently deletes all review data.
// Use with caution, primarily intended for development and testing purposes.
//
// Returns:
//   - error: Non-nil if the DROP TABLE operation fails
func (pg *PgDb) DropReviewTable() error {
	dropTableQuery := `DROP TABLE reviews`
	_, err := pg.db.Exec(dropTableQuery)
	if err != nil {
		fmt.Println("********************** Failed: Drop Review Table")
		return err
	}
	fmt.Println("********************** Success: Drop Review Table")
	return nil
}

// CreateReview inserts a new review into the database.
// The review's ID field is ignored as the database auto-generates it.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - review: The review data to persist (ID field is ignored)
//
// Returns:
//   - string: Success message including the creation timestamp
//   - error: Non-nil if the insert operation fails
//
// The operation is subject to the defaultTimeout (10 seconds).
func (pg *PgDb) CreateReview(ctx context.Context, review *Review) (string, error) {
	// Apply timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Execute the prepared INSERT statement
	_, err := pg.stmtCreate.ExecContext(ctx,
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

// UpdateReview modifies an existing review in the database.
// The review is identified by its ID field; all other fields are updated.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - review: The review data with updated values (ID must be set)
//
// Returns:
//   - error: Non-nil if the update fails or no review exists with the given ID
//
// The operation verifies that exactly one row was affected. If no rows are
// affected, an error is returned indicating the review was not found.
func (pg *PgDb) UpdateReview(ctx context.Context, review *Review) error {
	// Apply timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Execute the prepared UPDATE statement
	result, err := pg.stmtUpdate.ExecContext(ctx,
		review.Title,
		review.Director,
		review.ReleaseDate,
		review.Rating,
		review.ReviewNotes,
		review.ID)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	// Verify the update actually modified a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review with id %d not found", review.ID)
	}
	return nil
}

// DeleteReview removes a review from the database by its ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The unique identifier of the review to delete
//
// Returns:
//   - error: Non-nil if the deletion fails or no review exists with the given ID
//
// The operation verifies that exactly one row was affected. If no rows are
// affected, an error is returned indicating the review was not found.
func (pg *PgDb) DeleteReview(ctx context.Context, id int) error {
	// Apply timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Execute the prepared DELETE statement
	result, err := pg.stmtDelete.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// Verify the delete actually removed a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review with id %d not found", id)
	}
	fmt.Println("********************** Success: Deleted Review ", id)
	return nil
}

// GetReviewById retrieves a single review from the database by its unique ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The unique identifier of the review to retrieve
//
// Returns:
//   - *Review: The retrieved review data
//   - error: Non-nil if the query fails or no review exists with the given ID
//
// If no review is found with the given ID, an error wrapping sql.ErrNoRows
// is returned.
func (pg *PgDb) GetReviewById(ctx context.Context, id int) (*Review, error) {
	// Apply timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Execute the prepared SELECT statement and scan results into Review struct
	review := &Review{}
	err := pg.stmtGetById.QueryRowContext(ctx, id).Scan(
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
