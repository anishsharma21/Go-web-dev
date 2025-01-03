package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.Handle("GET /posts", GetPosts())
	mux.Handle("GET /posts/{id}", GetPostById())

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	mux.Handle("GET /", http.FileServer(http.Dir(".")))

	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(":"+port, mux)
}

func GetPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get posts\n")
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