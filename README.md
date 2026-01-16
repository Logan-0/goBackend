# Movie Review REST API

A high-performance RESTful API server for managing movie reviews, built with Go and PostgreSQL.

## Features

- **Full CRUD Operations** - Create, Read, Update, and Delete movie reviews
- **PostgreSQL Backend** - Reliable data persistence with connection pooling
- **Prepared Statements** - Optimized database queries for better performance
- **Graceful Shutdown** - Clean server shutdown with in-flight request completion
- **Environment Configuration** - Flexible configuration via environment variables
- **Chi Router** - Fast, lightweight HTTP routing (~3x faster than gorilla/mux)

## Quick Start

### Prerequisites

- Go 1.23+ 
- PostgreSQL 12+

### Database Setup

1. Start PostgreSQL on `localhost:5432`
2. Create the reviews table:

```sql
CREATE TABLE IF NOT EXISTS public.reviews (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    director VARCHAR NOT NULL,
    rating VARCHAR NOT NULL,
    releaseDate VARCHAR NOT NULL,
    reviewNotes VARCHAR NOT NULL,
    dateCreated VARCHAR NOT NULL
);
```

### Running the Server

```bash
# Clone and navigate to the project
cd goBackend

# Install dependencies
go mod download

# Run the server
go run .

# Or build and run
make build
./bin/goBackend
```

The server will start on `http://localhost:8080`

## Configuration

Configure the application using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `test` |
| `DB_NAME` | Database name | `postgres` |

Example:
```bash
export DB_HOST=mydb.example.com
export DB_PASSWORD=securepassword
go run .
```

## API Reference

### Create a Review

```http
POST /review
Content-Type: application/json

{
    "title": "Inception",
    "director": "Christopher Nolan",
    "releaseDate": "16 Jul 10 00:00 UTC",
    "rating": "9/10",
    "reviewNotes": "A mind-bending masterpiece about dreams within dreams."
}
```

**Response:** `200 OK`
```json
{
    "id": 1,
    "title": "Inception",
    "director": "Christopher Nolan",
    "releaseDate": "16 Jul 10 00:00",
    "rating": "9/10",
    "reviewNotes": "A mind-bending masterpiece about dreams within dreams.",
    "dateCreated": "16 Jan 26 17:30"
}
```

### Get a Review

```http
GET /review/{id}
```

**Response:** `200 OK`
```json
{
    "id": 1,
    "title": "Inception",
    "director": "Christopher Nolan",
    "releaseDate": "16 Jul 10 00:00",
    "rating": "9/10",
    "reviewNotes": "A mind-bending masterpiece about dreams within dreams.",
    "dateCreated": "16 Jan 26 17:30"
}
```

### Update a Review

```http
PUT /review/{id}
Content-Type: application/json

{
    "title": "Inception (Director's Cut)",
    "director": "Christopher Nolan",
    "releaseDate": "16 Jul 10 00:00",
    "rating": "10/10",
    "reviewNotes": "Even better on the second viewing!"
}
```

**Response:** `200 OK` - Returns the updated review

### Delete a Review

```http
DELETE /review/{id}
```

**Response:** `200 OK`
```json
{
    "deleted": "success"
}
```

### Error Responses

All endpoints return errors in a consistent format:

```json
{
    "Error": "review with id 999 not found"
}
```

## Project Structure

```
goBackend/
├── main.go      # Application entry point
├── api.go       # HTTP routing and handlers
├── storage.go   # Database access layer
├── types.go     # Domain models and DTOs
├── go.mod       # Go module definition
├── go.sum       # Dependency checksums
├── Makefile     # Build automation
└── README.md    # This file
```

## Architecture

The application follows a clean layered architecture:

- **Presentation Layer** (`api.go`) - HTTP handlers, routing, JSON serialization
- **Data Layer** (`storage.go`) - PostgreSQL access, connection pooling, prepared statements
- **Domain Layer** (`types.go`) - Business entities and data transfer objects

## Performance Optimizations

This server includes several performance optimizations:

1. **Prepared Statements** - SQL queries are parsed once at startup
2. **Connection Pooling** - 25 max connections with efficient reuse
3. **Chi Router** - Lightweight router with radix tree matching
4. **Context Timeouts** - 10-second database operation timeouts
5. **HTTP Timeouts** - Read/Write/Idle timeouts prevent resource exhaustion

## Development

```bash
# Build the binary
make build

# Run tests (when available)
go test ./...

# Generate documentation
go doc -all
```

## License

MIT License

