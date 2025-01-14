package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Handle requests to "/"
	http.HandleFunc("POST /", HelloServer)

	// Log that the server is running
	log.Println("Server is running at http://localhost:8080")

	// Start the server and handle any errors
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	// Extract the name from the URL path (after "/")
	name := r.URL.Path[1:]
	if name == "" {
		name = "World"
	}
	// Send response with "Hello, <name>!"
	fmt.Fprintf(w, "Hello, %s!", name)
}
