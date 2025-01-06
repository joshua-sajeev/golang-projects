package main

import (
	"fmt"
	"log"
	"net/http"
)

// Parse the form and print name and address
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// Get form values
	name := r.FormValue("name")
	address := r.FormValue("address")

	// Print the form values
	fmt.Fprintf(w, "POST Request Successful\n")
	fmt.Fprintf(w, "Name: %s\n", name)
	fmt.Fprintf(w, "Address: %s\n", address)
}

// Prints Hello World
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	// Serve static files from the "static" directory
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("GET /", fileServer) // Serve static files at the root

	// Handle /hello route
	http.HandleFunc("GET /hello", helloHandler)

	// Handle /form POST request
	http.HandleFunc("POST /form", formHandler)

	log.Println("Listening on http://localhost:8080")

	// Start the server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
