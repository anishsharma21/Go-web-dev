package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/practice_project_2/internal/templates"
	"example/practice_project_2/internal/types"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetUsers(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, created_at FROM users")
		if err != nil {
			http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
			log.Printf("error retrieving user data from db: %v\n", err)
			return
		}
		defer rows.Close()

		users := []types.User{}
		for rows.Next() {
			user := types.User{}

			err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
			if err != nil {
				http.Error(w, "Error processing user data", http.StatusInternalServerError)
				log.Printf("error scanning row from user data: %v\n", err)
				return
			}

			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
			log.Printf("error retrieving user data from db: %v\n", err)
			return
		}

		component := templates.Users(users)
		if r.Header.Get("HX-Request") == "" {
			component = templates.Base(component)
		}
		w.Header().Set("Content-Type", "text/html")
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering view", http.StatusInternalServerError)
			log.Printf("error rendering users view: %v\n", err)
			return
		}
	})
}

func GetUserById(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "User id not provided", http.StatusBadRequest)
			log.Println("user id not provided")
			return
		}

		var user types.User
		row := db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id)
		err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("user not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			log.Printf("error encoding user data: %v\n", err)
		}
	})
}

func AddUserPage(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		component := templates.AddUsers()
		w.Header().Set("Content-Type", "text/html")

		if r.Header.Get("HX-Request") == "" {
			component = templates.Base(component)
		}

		err := component.Render(context.Background(), w)
		if err != nil {
			http.Error(w, "Error rendering view", http.StatusInternalServerError)
			log.Printf("error rendering users view: %v\n", err)
			return
		}
	})
}

func AddUser(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			log.Printf("error parsing form data: %v\n", err)
			return
		}

		var user types.User
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		if user.Name == "" || user.Email == "" {
			http.Error(w, "Name or email not provided", http.StatusBadRequest)
			log.Println("name or email not provided")
			return
		}

		if user.CreatedAt == (time.Time{}) {
			user.CreatedAt = time.Now()
		}

		query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)"
		result, err := db.Exec(query, user.Name, user.Email, user.CreatedAt)
		if err != nil {
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			log.Printf("error adding user %v into users table: %v\n", user, err)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			log.Printf("error retrieving id from result: %v\n", err)
			return
		}

		user.Id = int(id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			log.Printf("error encoding user data: %v\n", err)
		}
	})
}

func DeleteUserById(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Error user id not given", http.StatusBadRequest)
			log.Println("user id not given in url path")
			return
		}

		userId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			log.Printf("invalid user id %q converted to %q: %v\n", id, userId, err)
			return
		}

		var user types.User
		query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
		row := db.QueryRow(query, id)
		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("error finding user to delete: %v\n", err)
			return
		}

		query = "DELETE FROM users WHERE id = ?"
		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, "User not deleted", http.StatusInternalServerError)
			log.Printf("error deleting user %v from users table: %v\n", user, err)
			return
		}

		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("error encoding user data: %v\nUser Data: %v\n", err, user)
			return
		}
	})
}

func UpdateUserById(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "User id not provided", http.StatusBadRequest)
			log.Printf("user id not provided in %q\n", r.URL.RawPath)
			return
		}

		userId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			log.Printf("invalid user id %q converted to %q: %v\n", id, userId, err)
			return
		}

		var user types.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Error decoding data", http.StatusInternalServerError)
			log.Printf("error decoding body data to update user: %v\nbody data: %v\n", err, r.Body)
			return
		}

		query := `UPDATE users
		SET name = ?, email = ?, created_at = ?
		WHERE id = ?;`
		_, err = db.Exec(query, &user.Name, &user.Email, &user.CreatedAt, userId)
		if err != nil {
			http.Error(w, "Error updating user data", http.StatusInternalServerError)
			log.Printf("error updating user data: %v\nuser data: %v\n", err, user)
			return
		}

		query = "SELECT id, name, email, created_at FROM users WHERE id = ?"
		row := db.QueryRow(query, userId)
		if err = row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt); err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("user not found after updating: %v\nuser id: %d\n", err, userId)
			return
		}

		if err = json.NewEncoder(w).Encode(&user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("error encoding user data: %v\nUser Data: %v\n", err, user)
			return
		}
	})
}
