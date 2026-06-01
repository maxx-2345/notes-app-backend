package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/maxx-2345/notes-app-backend/internal/config"
	"github.com/maxx-2345/notes-app-backend/internal/database"
	"github.com/maxx-2345/notes-app-backend/internal/models/repository"
	"github.com/maxx-2345/notes-app-backend/internal/router"
	"github.com/maxx-2345/notes-app-backend/internal/service"
)

func main() {
	// Initialize context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Fatal("Failed to load config", err)
		return
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&repository.Note{}); err != nil {
		log.Fatal("Failed to run database migrations: ", err)
	}

	repos := service.NewRepositories(db)

	r := chi.NewRouter()

	// Setup health check routes (public)
	// router.SetupHealthRoutes(ctx, r, db)

	// setup Note routes
	router.SetupNoteRoutes(ctx, r, repos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
