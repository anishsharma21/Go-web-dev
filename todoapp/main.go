package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := NewDb()

	mux := http.NewServeMux()

	mux.Handle("GET /", ShowTodos(db))
	mux.Handle("POST /", AddTodo(db))

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", mux)
}

func ShowTodos(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var content string
		var done bool

		err := db.QueryRow("select content, done from todos").Scan(&content, &done)
		if err != nil {
			http.Error(w, "data not found", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "The first todo is %s which is %v\n", content, done)
	})
}

func AddTodo(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		todo := params.Get("todo")
		if todo == "" {
			http.Error(w, "empty todo", http.StatusBadRequest)
		}

		_, err := db.Exec("insert into todos(content, done) values(?, ?)", todo, false)
		if err != nil {
			http.Error(w, "failed to insert todo", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "todo %s added\n", todo)
	})
}

func NewDb() *sql.DB {
	db, err := sql.Open("sqlite3", "todoapp.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`create table if not exists todos(
		id integer primary key autoincrement,
		content text not null,
		done boolean not null
	);`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}