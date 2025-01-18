package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	_, err := DBInit()
	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
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
