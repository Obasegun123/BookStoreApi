// db/database.go (updated)
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "bookstore"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnvWithDefault("DB_HOST", "localhost"),
		getEnvWithDefault("DB_PORT", "5432"),
		getEnvWithDefault("DB_USER", "postgres"),
		getEnvWithDefault("DB_PASSWORD", "postgres"),
		dbName,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL database")
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func CreateBooksTable() {
	query := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			published_date VARCHAR(20),
			isbn VARCHAR(20),
			price DECIMAL(10, 2)
		)
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Books table created or already exists")
}

// TestDBInit initializes database for testing
func TestDBInit(t *testing.T) *sql.DB {
	testDB, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=bookstore_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	_, err = testDB.Exec("DROP TABLE IF EXISTS books")
	if err != nil {
		t.Fatal(err)
	}

	query := `
		CREATE TABLE books (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			published_date VARCHAR(20),
			isbn VARCHAR(20),
			price DECIMAL(10, 2)
		)
	`

	_, err = testDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}

	return testDB
}
