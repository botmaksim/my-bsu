package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Fatal error connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Database not reachable: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "connected", "database": "%s"}`, os.Getenv("DB_NAME"))
	})

	log.Printf("Starting server on port %s...", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}
