package main

import (
	"context"
	"database/sql"
	"errors"
	"example/practice_project_2/internal/handlers"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := dbInit()
	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
	}

	mux := setupRoutes(db)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			log.Println(fmt.Sprintf("Server started on port %s...", port))
			return context.Background()
		},
	}

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server closed: %v\n", err)
		}
		log.Println("Stopped serving new connections.")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v\n", err)
	}

	<-shutdownChan
	log.Println("Graceful shutdown complete.")
}

func setupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /users", handlers.GetUsers(db))
	mux.Handle("GET /users/{id}", handlers.GetUserById(db))
	mux.Handle("POST /users", handlers.AddUser(db))
	mux.Handle("DELETE /users/{id}", handlers.DeleteUserById(db))
	mux.Handle("PUT /users/{id}", handlers.UpdateUserById(db))

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	mux.Handle("GET /", handlers.GetUsers(db))

	return mux
}

func dbInit() (*sql.DB, error) {
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
