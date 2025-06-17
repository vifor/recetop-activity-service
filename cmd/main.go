package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// main is the entrypoint for our application.
func main() {
	// http.NewServeMux() creates a new HTTP request router. It's good practice
	// to create our own mux instead of using the default one.
	mux := http.NewServeMux()

	// Register a handler function for the "/health" route.
	mux.HandleFunc("/health", healthCheckHandler)

	// Define the port we want the server to listen on.
	// We'll use 8081 to avoid conflicts with your other Spring Boot app on 8080.
	port := ":8081"

	log.Printf("Starting server on port %s", port)

	// http.ListenAndServe starts the HTTP server.
	// It will block forever unless it encounters an error.
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// healthCheckHandler is the function that will handle all requests to "/health".
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a simple map for our JSON response.
	response := map[string]string{"status": "UP"}

	// Set the Content-Type header so the client knows we are sending JSON.
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code to 200 OK.
	w.WriteHeader(http.StatusOK)

	// Encode our map into JSON and write it to the response.
	json.NewEncoder(w).Encode(response)
}