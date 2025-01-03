package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// FIXME should be called from main.go in main package based on env
func DbUp() {
	db, err := sql.Open("sqlite3", "posts.db")
	if err != nil {
		log.Fatalf("Error: could not open posts.db: %v\n", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	content TEXT NOT NULL
	likes INTEGER DEFAULT 0
	);`)
	if err != nil {
		log.Fatalf("Error: could not create table: %v\n", err)
	}
}

func DbDown() {
	db, err := sql.Open("sqlite3", "posts.db")
	if err != nil {
		log.Fatalf("Error: could not open posts.db: %v\n", err)
	}
	defer db.Close()

	_, err = db.Exec(`DROP TABLE IF EXISTS posts;`)
	if err != nil {
		log.Fatalf("Error: could not drop table: %v\n", err)
	}
}