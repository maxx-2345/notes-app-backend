package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/maxx-2345/notes-app-backend/internal/handlers"
	"github.com/maxx-2345/notes-app-backend/internal/service"
)

func SetupNoteRoutes(ctx context.Context, r *chi.Mux, repos *service.Repositories) {
	r.Route("/notes", func(r chi.Router) {

		// get all
		r.Get("/", handlers.GetAllNotes(ctx, repos))

		// get by id
		r.Get("/{id}", handlers.GetNoteByID(ctx, repos))

		// create note
		r.Post("/", handlers.CreateNote(ctx, repos))

		// delete note
		r.Delete("/{id}", handlers.DeleteNote(ctx, repos))
	})
}
