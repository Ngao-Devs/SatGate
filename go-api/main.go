package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pressly/goose/v3"
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
	if err := run(); err != nil {
		log.Printf("go-api: %v", err)
		os.Exit(1)
	}
}

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL must be set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connectionMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("configure migrations: %w", err)
	}
	goose.SetBaseFS(migrationFiles)
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	log.Print("PostgreSQL connection ready; migrations applied")

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
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
		log.Printf("go-api: write health response: %v", err)
	}
}