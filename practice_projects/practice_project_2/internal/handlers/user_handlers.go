package handlers

import (
	"database/sql"
	"encoding/json"
	"example/practice_project_2/internal/types"
	"fmt"
	"net/http"
)

func GetUsers(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, created_at FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []types.User{}
		for rows.Next() {
			user := types.User{}

			err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func GetUserById(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id not provided", http.StatusBadRequest)
			return
		}

		var user types.User
		row := db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id)
		err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("User not found: %v\n", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
