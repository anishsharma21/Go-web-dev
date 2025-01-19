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
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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
			log.Printf("Server started on port %v...\n", port)
			return context.Background()
		},
	}

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server closed early: %v\n", err)
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

func setupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
