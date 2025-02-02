package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	DB     *sql.DB
	Logger *slog.Logger
}

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize the application: %v\n", err)
	}
	defer app.DB.Close()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.setupRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			app.Logger.Info("Server started on port 8080...\n")
			return context.Background()
		},
	}

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("HTTP server closed early: %v", slog.String("error", err.Error()))
		}
		app.Logger.Info("Stopped serving new connections.")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("HTTP shutdown error: %v\n", slog.String("error", err.Error()))
	}
	<-shutdownChan
	close(shutdownChan)

	app.Logger.Info("Graceful server shutdown complete.")
}

func initializeApp() (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := SetupDb()
	if err != nil {
		logger.Error("Failed to initialise connection to the database: %v\n", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Info("Initialised the database connection successfully!\n")

	return &App{
		DB:     db,
		Logger: logger,
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

func (app *App) HandlerWrapper(handlerFunc func(*sql.DB, *slog.Logger) http.Handler) http.Handler {
	return handlerFunc(app.DB, app.Logger)
}
