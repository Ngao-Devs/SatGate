package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

const (
	connectionTimeout     = 10 * time.Second
	maxOpenConnections    = 25
	maxIdleConnections    = 5
	connectionMaxLifetime = 5 * time.Minute
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/satgo?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL database")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS forms (
		id VARCHAR(255) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created forms table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id VARCHAR(255) PRIMARY KEY,
		form_id VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created messages table")

	// Health endpoint used as a smoke test by deployments and load balancers.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("go-api: listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("start http server: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
		log.Printf("go-api: write health response: %v", err)
	}
}
