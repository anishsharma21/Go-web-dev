package handlers

import (
	"database/sql"
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

		userStr := ""
		for _, user := range users {
			userStr += fmt.Sprintf("%d\t%s\t%s\t%v\n", user.Id, user.Name, user.Email, user.CreatedAt)
		}

		fmt.Fprintf(w, userStr)
	})
}