package handlers

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
)

func BaseHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("templates/index.html")
		if err != nil {
			slog.Error("Error parsing template file", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := template.Execute(w, nil); err != nil {
			slog.Error("Error executing template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
