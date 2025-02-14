package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anishsharma21/Go-web-dev/practice_projects/postgres_newdeploy_project/internal/handlers"
	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app, err := initializeApp()
	if err != nil {
		slog.Error("Failed to initialise the application", "error", err)
		return
	}
	defer func() {
		if app.DB != nil {
			if err := app.DB.Close(); err != nil {
				slog.Error("Failed to close the database connection", "error", err)
			}
		}
	}()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.setupRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			slog.Info("Server started on port 8080...")
			return context.Background()
		},
	}

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server closed early", "error", err)
		}
		slog.Info("Stopped serving new connections.")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Warn("Received signal", "signal", sig.String())

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	slog.Info("Shutting down server gracefully...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}

func initializeApp() (*App, error) {
	db, err := SetupDb()
	if err != nil {
		slog.Error("Failed to initialise connection to the database", "error", err)
		return nil, err
	}
	slog.Info("Initialised the database connection successfully!")

	return &App{
		DB: db,
	}, nil
}

func (app *App) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /users", app.HandlerWrapper(handlers.GetUsers))
	mux.Handle("POST /users", app.HandlerWrapper(handlers.AddUser))
	mux.Handle("GET /", app.HandlerWrapper(handlers.BaseHandler))

	return mux
}

func SetupDb() (*sql.DB, error) {
	env := os.Getenv("ENV")

	var dsn string
	if env == "production" || env == "cicd" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		dsn = "host=localhost port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping the database: %v\n", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the 'users' table: %v\n", err)
	}

	return db, nil
}

func (app *App) HandlerWrapper(handlerFunc func(*sql.DB) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handlerFunc(app.DB).ServeHTTP(w, r)
		duration := time.Since(start).Milliseconds()
		slog.Info("Request processed", "method", r.Method, "url", r.URL.Path, "processing_duration", fmt.Sprintf("%vms", duration))
	})
}
