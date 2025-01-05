package handlers

import (
	"database/sql"
	"encoding/json"
	"example/practice_project_2/internal/types"
	"fmt"
	"log"
	"net/http"
)

// TODO fix up logging, expose less info to client, log errors for debugging
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

func AddUser(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user types.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request payload: %v\n", err), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)"
		result, err := db.Exec(query, user.Name, user.Email, user.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Id = int(id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode data: %v\n", err), http.StatusInternalServerError)
		}
	})
}

func DeleteUserById(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Invalid id given: \"\"", http.StatusBadRequest)
			return
		}

		var user types.User
		query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
		row := db.QueryRow(query, id)
		err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		query = "DELETE FROM users WHERE id = ?"
		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to delete user with id %s: %v\n", id, err), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("DeleteUserById: Error encoding response: %v\n", err)
			return
		}
	})
}
