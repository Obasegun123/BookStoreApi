// main_test.go
package main

import (
	"bookstore-api/db"
	"bookstore-api/handlers"
	"bookstore-api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router
var testBook models.Book

func TestMain(m *testing.M) {
	// Setup
	setup()

	// Run tests
	code := m.Run()

	// Cleanup
	teardown()

	os.Exit(code)
}

func setup() {
	// Initialize database connection
	db.InitDB()
	db.CreateBooksTable()

	// Create router
	router = mux.NewRouter()

	// Register routes
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Create a test book for testing GET, UPDATE, DELETE
	testBook = models.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedDate: "2023-01-01",
		ISBN:          "1234567890",
		Price:         19.99,
	}

	// Insert test book into database
	query := `
		INSERT INTO books (title, author, published_date, isbn, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := db.DB.QueryRow(
		query,
		testBook.Title,
		testBook.Author,
		testBook.PublishedDate,
		testBook.ISBN,
		testBook.Price,
	).Scan(&testBook.ID)
	if err != nil {
		panic("Failed to create test book: " + err.Error())
	}
}

func teardown() {
	// Clean up database
	_, err := db.DB.Exec("DELETE FROM books")
	if err != nil {
		panic("Failed to clean up test database: " + err.Error())
	}
	db.DB.Close()
}

func TestGetBooks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/books", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var books []models.Book
	err := json.Unmarshal(response.Body.Bytes(), &books)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(books), 1)
}

func TestGetBook(t *testing.T) {
	req, _ := http.NewRequest("GET", "/books/"+strconv.Itoa(testBook.ID), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var book models.Book
	err := json.Unmarshal(response.Body.Bytes(), &book)
	assert.Nil(t, err)
	assert.Equal(t, testBook.ID, book.ID)
	assert.Equal(t, testBook.Title, book.Title)
}

func TestGetBookNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/books/9999", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestCreateBook(t *testing.T) {
	newBook := models.Book{
		Title:         "New Book",
		Author:        "New Author",
		PublishedDate: "2023-02-01",
		ISBN:          "0987654321",
		Price:         29.99,
	}

	jsonValue, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusCreated, response.Code)

	var createdBook models.Book
	err := json.Unmarshal(response.Body.Bytes(), &createdBook)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, createdBook.ID)
	assert.Equal(t, newBook.Title, createdBook.Title)
}

func TestCreateBookInvalidData(t *testing.T) {
	invalidBook := map[string]interface{}{
		"title":  123, // Invalid type for title
		"author": "Invalid Author",
	}

	jsonValue, _ := json.Marshal(invalidBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUpdateBook(t *testing.T) {
	updatedBook := models.Book{
		Title:         "Updated Book Title",
		Author:        testBook.Author,
		PublishedDate: testBook.PublishedDate,
		ISBN:          testBook.ISBN,
		Price:         25.99,
	}

	jsonValue, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest("PUT", "/books/"+strconv.Itoa(testBook.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var book models.Book
	err := json.Unmarshal(response.Body.Bytes(), &book)
	assert.Nil(t, err)
	assert.Equal(t, testBook.ID, book.ID)
	assert.Equal(t, updatedBook.Title, book.Title)
	assert.Equal(t, updatedBook.Price, book.Price)
}

func TestUpdateBookNotFound(t *testing.T) {
	updatedBook := models.Book{
		Title:         "Updated Book Title",
		Author:        testBook.Author,
		PublishedDate: testBook.PublishedDate,
		ISBN:          testBook.ISBN,
		Price:         25.99,
	}

	jsonValue, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest("PUT", "/books/9999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestDeleteBook(t *testing.T) {
	// First create a book to delete
	newBook := models.Book{
		Title:         "Book to Delete",
		Author:        "Delete Author",
		PublishedDate: "2023-03-01",
		ISBN:          "1122334455",
		Price:         15.99,
	}

	jsonValue, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	var createdBook models.Book
	json.Unmarshal(response.Body.Bytes(), &createdBook)

	// Now delete the book
	req, _ = http.NewRequest("DELETE", "/books/"+strconv.Itoa(createdBook.ID), nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNoContent, response.Code)

	// Verify the book is deleted
	req, _ = http.NewRequest("GET", "/books/"+strconv.Itoa(createdBook.ID), nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestDeleteBookNotFound(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/books/9999", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}
