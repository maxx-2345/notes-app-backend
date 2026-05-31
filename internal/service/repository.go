package service

import (
	"github.com/maxx-2345/notes-app-backend/internal/database"
	"github.com/maxx-2345/notes-app-backend/internal/models/repository"
)

// Repositories holds all repository instances for different models.
// This provides a centralized way to access all repositories.
// Add new repositories here as you create new models.
type Repositories struct {
	Note *database.Repository[repository.Note]
}

// NewRepositories initializes all repositories for all models.
// This is a convenient way to set up all repositories at once.
// When you add a new model:
// 1. Add it to the Repositories struct above
// 2. Initialize it in this function
func NewRepositories(db *database.Database) *Repositories {
	return &Repositories{
		Note: database.NewRepository[repository.Note](db),
	}
}

