// main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct to represent the data model
type Book struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Author   string `json:"author,omitempty" bson:"author,omitempty"`
	Year     int    `json:"year,omitempty" bson:"year,omitempty"`
}

// Database configuration
const (
	DatabaseURL    = "mongodb://localhost:27017"
	DatabaseName   = "Golib"
	CollectionName = "book"
)

// MongoDB client
var client *mongo.Client

// Initialize MongoDB connection
func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI(DatabaseURL)

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

// Handlers

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a slice to store the result
	var books []Book

	// Get all documents from the collection
	cur, err := client.Database(DatabaseName).Collection(CollectionName).Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	// Iterate through the cursor and decode each document
	for cur.Next(context.Background()) {
		var book Book
		err := cur.Decode(&book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Convert the result to JSON and send it in the response
	json.NewEncoder(w).Encode(books)
}

// Get a single book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	bookID := params["id"]

	// Find the document with the specified ID
	var book Book
	err := client.Database(DatabaseName).Collection(CollectionName).FindOne(context.Background(), bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert the result to JSON and send it in the response
	json.NewEncoder(w).Encode(book)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a Book object
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set a unique ID for the new book
	book.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// Insert the new book into the collection
	_, err = client.Database(DatabaseName).Collection(CollectionName).InsertOne(context.Background(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the created book in the response
	json.NewEncoder(w).Encode(book)
}

// Update an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	bookID := params["id"]

	// Decode the request body into a Book object
	var updatedBook Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the document with the specified ID
	_, err = client.Database(DatabaseName).Collection(CollectionName).UpdateOne(
		context.Background(),
		bson.M{"_id": bookID},
		bson.D{{Key: "$set", Value: updatedBook}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send the updated book in the response
	json.NewEncoder(w).Encode(updatedBook)
}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	bookID := params["id"]

	// Delete the document with the specified ID
	_, err := client.Database(DatabaseName).Collection(CollectionName).DeleteOne(context.Background(), bson.M{"_id": bookID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send success message in the response
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}

// Main function
func main() {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define API endpoints and their corresponding handlers
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))
}


/*

1. `-json:"id,omitempty" bson:"_id,omitempty"`: This is a struct tag in Go, used to provide metadata about the fields of a struct. In this case, it is used for both JSON and BSON (MongoDB's binary JSON) serialization. The `id` field in the struct will be mapped to the `"_id"` field in MongoDB, and the `omitempty` option ensures that the field won't be included in the serialized output if its value is empty.

2. `-clientOptions := options.Client().ApplyURI(DatabaseURL)`: This line initializes MongoDB client options. It creates a new instance of `options.Client` and sets the MongoDB connection URI using `ApplyURI`. The `clientOptions` variable is then used to configure the MongoDB client.

3. `-w http.ResponseWriter, r *http.Request`: This line is a function signature for an HTTP handler in Go. It represents the `ResponseWriter` and `Request` parameters that are passed to an HTTP handler function. `w` is used to construct an HTTP response, and `r` contains information about the incoming HTTP request.

4. `-var books []Book`: This line declares a variable named `books` that is a slice of the `Book` struct. In Go, a slice is a dynamically-sized, flexible view into an array. Slices are more common than arrays for most applications.

5. `-cur, err := client.Database(DatabaseName).Collection(CollectionName).Find(context.Background(), bson.M{})`: This line retrieves documents from the MongoDB collection specified by `DatabaseName` and `CollectionName`. The `Find` method returns a cursor (`cur`) that can be iterated to access the results.

6. `-defer cur.Close(context.Background())`: This line defers the closing of the cursor until the surrounding function returns. It ensures that the cursor is closed after the loop that iterates over the results is finished.

7. `-var book Book`: Declares a variable named `book` of type `Book`. It is later used to decode a document from the MongoDB cursor.

8. ```
   err := cur.Decode(&book)
   if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }
   ```: Decodes the current document from the MongoDB cursor into the `book` variable. If there is an error during decoding, it returns an internal server error with the error message.

9. `-books = append(books, book)`: Appends the decoded `book` to the `books` slice.

10. `-json.NewEncoder(w).Encode(books)`: Encodes the `books` slice as JSON and writes it to the HTTP response.

11. `-params := mux.Vars(r)`: Retrieves the URL parameters from the request using Gorilla Mux. `params` is a map containing the variables defined in the route pattern.

12. `-book.ID = fmt.Sprintf("%d", time.Now().UnixNano())`: Sets a unique ID for a new book. It uses the current Unix timestamp in nanoseconds to generate a unique identifier.

13. ```
    _, err = client.Database(DatabaseName).Collection(CollectionName).InsertOne(context.Background(), book)
    if err != nil {
    ```: Inserts a new document (book) into the MongoDB collection. If there is an error during the insertion, it returns an error.

These lines cover various aspects of Go programming, HTTP handling, MongoDB operations, and struct serialization. Understanding these concepts will help you grasp how the given CRUD application is structured.

*/