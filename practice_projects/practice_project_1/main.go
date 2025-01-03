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

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(":"+port, mux)
}