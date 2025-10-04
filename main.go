package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}

	// Start the HTTP server and listen on port 8080
	log.Printf("Server listening on: %s\n", port)
	log.Fatal(server.ListenAndServe())
} 