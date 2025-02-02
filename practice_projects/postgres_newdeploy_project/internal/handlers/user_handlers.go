package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/anishsharma21/Go-web-dev/practice_projects/postgres_newdeploy_project/internal/types"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func AddUser(db *sql.DB, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		insertUserQuery := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;"
		randomName := randString(rand.Intn(5) + 5)
		randomEmail := fmt.Sprintf("%s@gmail.com", randomName)

		var id int
		err := db.QueryRow(insertUserQuery, randomName, randomEmail).Scan(&id)
		if err != nil {
			logger.Error("Error inserting user into table: %v\n", slog.String("error", err.Error()))
			http.Error(w, "error adding user", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User added to database:\nID: %d\nName: %s\nEmail: %s\n", id, randomName, randomEmail)
	})
}

func GetUsers(db *sql.DB, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var users []types.User
		getUsersError := "error getting users data"

		getUsersQuery := "SELECT * FROM users"
		rows, err := db.Query(getUsersQuery)
		if err != nil {
			logger.Error("Error getting users from 'users' table", slog.String("error", err.Error()))
			http.Error(w, getUsersError, http.StatusInternalServerError)
			return
		}

		var user types.User
		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
			if err != nil {
				logger.Error("Error scanning user row", slog.String("error", err.Error()))
				http.Error(w, getUsersError, http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			logger.Error("Error with rows: %v\n", slog.String("error", err.Error()))
			http.Error(w, getUsersError, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			logger.Error("Error encoding users to JSON: %v\n", slog.String("error", err.Error()))
			http.Error(w, getUsersError, http.StatusInternalServerError)
			return
		}
	})
}
