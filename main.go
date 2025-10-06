package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"


	mux := http.NewServeMux()

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}
	
	filesystem := http.Dir(filepathRoot)
	handler := http.FileServer(filesystem)
	mux.Handle("/", handler)
	mux.Handle("/assets/", handler)

	// Start the HTTP server and listen on port 8080
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
} 