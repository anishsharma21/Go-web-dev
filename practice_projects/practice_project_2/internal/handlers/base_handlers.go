package handlers

import (
	"example/practice_project_2/internal/templates"
	"log"
	"net/http"
)

func Base() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		component := templates.Base(nil)
		w.Header().Set("Content-Type", "text/html")
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering view", http.StatusInternalServerError)
			log.Printf("error rendering nav/base view: %v\n", err)
			return
		}
	})
}
