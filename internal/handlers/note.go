package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/maxx-2345/notes-app-backend/internal/service"
)

func GetAllNotes(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := repos.Note.SelectAll(r.Context(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
