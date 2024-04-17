package main

import (
	"encoding/json"
	"math/rand"
	"strconv"

	// "fmt"
	"log"

	"github.com/gorilla/mux"

	// "math/rand"
	"net/http"
	// "strconv"
)

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:lastname`
}

// Init books var as a slice Book struct
var books []Book

// get Allbooks
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

}

// get single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) //get Params
	for _, item := range books {
		if item.ID == params["id"] {
			// Create a new JSON encoder
			encoder := json.NewEncoder(w)
			// Encode the item and write it to the response writer
			if err := encoder.Encode(item); err != nil {
				// Handle error if encoding fails
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	// If book with given ID is not found, respond with an empty object
	json.NewEncoder(w).Encode(&Book{})
}

// create newbooks
func createBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)

}

// Delete Books
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Remove the book from the slice
			books = append(books[:index], books[index+1:]...)
			// Encode the updated list of books and send it as response
			json.NewEncoder(w).Encode(books)
			return
		}
	}

	// If the book with the given ID is not found, respond with an empty object
	json.NewEncoder(w).Encode(&Book{})
}

// update Books
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Decode the incoming JSON data into a Book struct
			var updatedBook Book
			err := json.NewDecoder(r.Body).Decode(&updatedBook)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Update the book with new values
			updatedBook.ID = params["id"]
			books[index] = updatedBook
			// Encode the updated book and send it as response
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	json.NewEncoder(w).Encode(books)

}
func main() {
	r := mux.NewRouter()

	// test Data
	books = append(books, Book{ID: "1", Isbn: "4423", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})

	books = append(books, Book{ID: "2", Isbn: "4424", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	r.HandleFunc("/api/books", getBooks).Methods("GET")

	r.HandleFunc("/api/books", createBooks).Methods("POST")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")

	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")

	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
