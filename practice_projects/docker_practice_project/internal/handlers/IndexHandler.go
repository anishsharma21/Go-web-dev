package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(db *sql.DB, templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "unable to render template", http.StatusInternalServerError)
			log.Printf("Unable to render index.html template: %v\n", err)
			return
		}
	})
}
