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
	defer app.DB.Close()

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

		select {
		case shutdownChan <- true:
		default:
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("Received signal", "signal", sig)

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

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
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	return mux
}

func SetupDb() (*sql.DB, error) {
	env := os.Getenv("ENV")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var dsn string
	if env == "production" || env == "cicd" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
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
	return handlerFunc(app.DB)
}
