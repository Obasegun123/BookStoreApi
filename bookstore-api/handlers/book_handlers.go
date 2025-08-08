// handlers/book_handlers.go
package handlers

import (
	"bookstore-api/db"
	"bookstore-api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// validateBook validates the book data
func validateBook(book models.Book) error {
	if strings.TrimSpace(book.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(book.Author) == "" {
		return fmt.Errorf("author is required")
	}
	if book.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if len(book.Title) > 255 {
		return fmt.Errorf("title cannot exceed 255 characters")
	}
	if len(book.Author) > 255 {
		return fmt.Errorf("author cannot exceed 255 characters")
	}
	return nil
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, author, published_date, isbn, price FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.ISBN, &book.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	err = db.DB.QueryRow("SELECT id, title, author, published_date, isbn, price FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.ISBN, &book.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate book data
	if err := validateBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO books (title, author, published_date, isbn, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = db.DB.QueryRow(
		query,
		book.Title,
		book.Author,
		book.PublishedDate,
		book.ISBN,
		book.Price,
	).Scan(&book.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate book data
	if err := validateBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		UPDATE books
		SET title = $1, author = $2, published_date = $3, isbn = $4, price = $5
		WHERE id = $6
		RETURNING id, title, author, published_date, isbn, price
	`

	err = db.DB.QueryRow(
		query,
		book.Title,
		book.Author,
		book.PublishedDate,
		book.ISBN,
		book.Price,
		id,
	).Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.ISBN, &book.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
