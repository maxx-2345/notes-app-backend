package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/maxx-2345/notes-app-backend/internal/models/repository"
	"github.com/maxx-2345/notes-app-backend/internal/models/request"
	"github.com/maxx-2345/notes-app-backend/internal/service"
	"gorm.io/gorm"
)

// Get All Notes
func GetAllNotes(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := repos.Note.SelectAll(r.Context(), nil)
		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithCtx(ctx, w, notes, http.StatusOK)
	}
}

// Get Note By ID
func GetNoteByID(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		note, err := repos.Note.SelectByID(r.Context(), uint(id), nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				respondWithCtx(ctx, w, "record not found", http.StatusNotFound)
				return
			}
			respondWithCtx(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithCtx(ctx, w, note, http.StatusOK)
	}
}

// Create Note
func CreateNote(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request request.CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err := validate.Struct(request)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			respondWithCtx(ctx, w, fmt.Sprintf("Validation error: %s", validationErrors), http.StatusBadRequest)
			return
		}

		note := repository.Note{
			Title:   request.Title,
			Content: request.Content,
		}

		if err := repos.Note.Create(r.Context(), &note); err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithCtx(ctx, w, note, http.StatusCreated)
	}
}

func DeleteNote(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondWithCtx(ctx, w, "invalid note ID", http.StatusBadRequest)
			return
		}

		if err := repos.Note.Delete(ctx, uint(id)); err != nil {
			respondWithCtx(ctx, w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateNote(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondWithCtx(ctx, w, "Invalid Note ID", http.StatusBadRequest)
			return
		}

		var request request.UpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err = validate.Struct(request)
		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		note, err := repos.Note.SelectByID(ctx, uint(id), []string{})
		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusNotFound)
			return
		}

		note.Title = request.Title
		note.Content = request.Content

		if err := repos.Note.Update(ctx, note); err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithCtx(ctx, w, note, http.StatusOK)
	}
}

func PatchNote(ctx context.Context, repos *service.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondWithCtx(ctx, w, "Invalid Note ID", http.StatusBadRequest)
			return
		}

		var request request.PatchNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err = validate.Struct(request)
		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repos.Note.PatchUpdate(ctx, uint(id), &repository.Note{
			Title:   request.Title,
			Content: request.Content,
		})

		if err != nil {
			respondWithCtx(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithCtx(ctx, w, request, http.StatusOK)
	}
}
