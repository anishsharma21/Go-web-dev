package handlers_test

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anishsharma21/Go-web-dev/practice_projects/postgres_newdeploy_project/internal/handlers"
	_ "github.com/lib/pq"
)

var db *sql.DB
var server *httptest.Server

func init() {
	dbConnStr := os.Getenv("DATABASE_URL")

	if dbConnStr == "" {
		dbConnStr = "host=localhost port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v\n", err)
	}

	// Setup test handlers
	handler := handlers.AddUser(db)
	server = httptest.NewServer(handler)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v\n", err)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()

	server.Close()
	db.Close()

	os.Exit(code)
}
