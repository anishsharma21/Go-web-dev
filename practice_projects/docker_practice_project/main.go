package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	db, err := DBInit()
	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
	}

	mux := setupRoutes(db)

	log.Printf("HTTP server started on port %v...", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func DBInit() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, fmt.Errorf("could not open the database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %v\n", err)
	}

	return db, nil
}

func setupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
