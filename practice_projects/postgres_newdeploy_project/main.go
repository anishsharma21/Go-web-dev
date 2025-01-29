package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
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
