package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/maxx-2345/notes-app-backend/internal/handlers"
	"github.com/maxx-2345/notes-app-backend/internal/service"
)

func SetupNoteRoutes(ctx context.Context, r *chi.Mux, repos *service.Repositories) {
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", handlers.GetAllNotes(ctx,repos))
	})	
}
