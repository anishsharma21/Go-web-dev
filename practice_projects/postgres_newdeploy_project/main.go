package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	_, err := SetupDb()
	if err != nil {
		log.Fatalf("Failed to initialise connection to the database: %v\n", err)
	}
	log.Printf("Initialised the database connection succesfully!\n")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			log.Printf("Server started on port 8080...\n")
			return context.Background()
		},
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server closed early: %v", err)
	}
	log.Println("Server shutdown.")
}

func SetupDb() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping the database: %v\n", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS questions (
	id SERIAL PRIMARY KEY,
	question TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the 'questions' table: %v\n", err)
	}

	return db, nil
}
