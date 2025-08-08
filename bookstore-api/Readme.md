# Bookstore CRUD API

![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)
![Docker](https://img.shields.io/badge/Docker-Compatible-blue.svg)

A lightweight RESTful API for managing bookstore inventory, built with Go's standard library and PostgreSQL.

## Features

- **CRUD Operations**: Full Create, Read, Update, Delete functionality
- **Lightweight**: Uses only Go's standard `net/http` package
- **PostgreSQL**: Robust relational database backend
- **Docker Support**: Easy containerized development and deployment
- **Comprehensive Testing**: 90%+ test coverage
- **Clean Architecture**: Separation of concerns with handlers, models, and database layers

## API Documentation

### Endpoints

| Method | Endpoint       | Description                     | Status Codes               |
|--------|----------------|---------------------------------|----------------------------|
| GET    | `/books`       | List all books                  | 200 OK                     |
| GET    | `/books/{id}`  | Get book by ID                  | 200 OK, 404 Not Found      |
| POST   | `/books`       | Create new book                 | 201 Created, 400 Bad Request |
| PUT    | `/books/{id}`  | Update existing book            | 200 OK, 400 Bad Request, 404 Not Found |
| DELETE | `/books/{id}`  | Delete book                     | 204 No Content, 404 Not Found |

### Request/Response Examples

**Create Book:**
```http
POST /books HTTP/1.1
Content-Type: application/json

{
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "published_date": "2008-08-01",
    "isbn": "978-0132350884",
    "price": 35.99
}

Response:

json
{
    "id": 1,
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "published_date": "2008-08-01",
    "isbn": "978-0132350884",
    "price": 35.99
}
Getting Started
Prerequisites
Docker and Docker Compose

Go 1.21+ (optional if using Docker)

Installation
Clone the repository:

bash
git clone https://github.com/Obasegun123/bookstore-api.git
cd bookstore-api
Copy the example environment file:

bash
cp .env.example .env
Running with Docker (Recommended)
bash
# Start the application and database
docker-compose up -d --build

# Run tests
docker-compose run test
The API will be available at http://localhost:8080

Running Locally
Set up PostgreSQL and update .env file

Install dependencies:

bash
go mod download
Run the application:

bash
go run main.go
Run tests:

bash
go test -v ./...
Project Structure
text
bookstore-api/
├── db/               # Database connection and setup
│   └── database.go
├── handlers/         # HTTP request handlers
│   └── book_handlers.go
├── models/           # Data models
│   └── book.go
├── tests/            # Test files
│   └── main_test.go
├── .env.example      # Environment configuration example
├── docker-compose.yml # Docker orchestration
├── Dockerfile        # Container configuration
├── go.mod           # Go dependencies
├── go.sum           # Go dependencies checksum
└── main.go          # Application entry point
Testing
The project includes comprehensive tests with 90%+ coverage:

bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
Environment Variables
Variable	Description	Default
DB_HOST	Database host	localhost
DB_PORT	Database port	5432
DB_USER	Database username	postgres
DB_PASSWORD	Database password	postgres
DB_NAME	Database name	bookstore


