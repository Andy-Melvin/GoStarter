package main

import (
	"fmt"
	"log"
	"net/http"
)

// We write the Hello Handler func
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return // Move the return statement here
	}

	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return // Move the return statement here
	}

	fmt.Fprintf(w, "HELLO GO LANG APIS")
}

// Now we write the Form Handler also
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { // Corrected the method name to ParseForm
		http.Error(w, fmt.Sprintf("ParseForm() err: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "POST REQUEST SUCCESSFUL\n")
	name := r.FormValue("name")
	fmt.Fprintf(w, "Name = %s", name)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static")) // Added variable declaration
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting the server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { // Corrected the port format
		log.Fatal(err)
	}
}
