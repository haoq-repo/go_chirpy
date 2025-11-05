package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/haoq-repo/go_chirpy/internal/database"
	"github.com/joho/godotenv"
	// importing for side effects, not usage
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db 				*database.Queries
	platform		string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	
	// load the .env file
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("sql.Open error: %v", err)
		return
	}

	if err:= db.Ping(); err != nil {
		log.Printf("DB ping failed: %v", err)
		return
	}
	defer db.Close()
	dbQueries := database.New(db)


	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform:		platform,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}
	
	filesystem := http.Dir(filepathRoot)
	handler := http.StripPrefix("/app", http.FileServer(filesystem))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGet)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	
	// Start the HTTP server and listen on port 8080
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}