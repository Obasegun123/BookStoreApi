// main.go
package main

import (
	"bookstore-api/db"
	"bookstore-api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db.InitDB()
	defer db.DB.Close()
	db.CreateBooksTable()

	// Create router
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
