package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Post struct {
	ID      int
	Content string
	Likes   int
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := NewSQLiteDB("posts.db")

	// NOTE for now its fine to set routes here, but in the future, set them in a SetRoutes func
	mux := http.NewServeMux()

	mux.Handle("GET /posts", GetPosts(db))
	mux.Handle("GET /posts/{id}", GetPostById())
	mux.Handle("POST /posts", AddPost(db))

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	mux.Handle("GET /", http.FileServer(http.Dir(".")))

	fmt.Printf("Server starting on port %s...\n", port)
	// TODO signals to gracefully exit and delete db at the end of session based on cli arg
	http.ListenAndServe(":"+port, mux)
}

func NewSQLiteDB(dbSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v\n", err)
	}

	return db
}

func AddPost(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		content := r.PostForm.Get("content")
		if content == "" {
			http.Error(w, "Content not given", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO posts (content) VALUES (?);", content)
		if err != nil {
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	})
}

func GetPosts(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var posts []Post
		rows, err := db.Query("SELECT * FROM posts;")
		if err != nil {
			http.Error(w, "Failed to query posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Content, &post.Likes); err != nil {
				http.Error(w, "Failed to scan post", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}
	})
}

func GetPostById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id not given", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Post ID: %s\n", id)
	})
}