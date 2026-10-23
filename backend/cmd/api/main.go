package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"

	"backend/internal/clients"
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Init Cloudinary
	storageClient, err := clients.NewCloudinaryClient(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("failed to init cloudinary: %v", err)
	}

	// Init Firestore
	dbClient, err := clients.NewFirestoreClient(ctx, cfg.FirebaseProjectID)
	if err != nil {
		log.Fatalf("failed to init firestore: %v", err)
	}

	// Init Firebase Auth for Middleware.
	// Fail fast: never start the server with protected routes left unprotected.
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("failed to init firebase app: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to init firebase auth: %v", err)
	}
	authMiddleware := middleware.FirebaseAuth(authClient)

	// Wire up layers
	artworkSvc := services.NewArtworkService(storageClient, dbClient)
	artworkHandler := handlers.NewArtworkHandler(artworkSvc)

	resumeSvc := services.NewResumeService(storageClient, dbClient)
	resumeHandler := handlers.NewResumeHandler(resumeSvc)

	r := gin.Default()

	// CORS: only the configured frontend origin may call the API
	r.Use(middleware.CORS(cfg.FrontendOrigin))

	// Serves the entire frontend folder at the root path "/"
	r.StaticFS("/frontend", http.Dir("../frontend"))

	// Routes
	api := r.Group("/api")

	// Public routes
	api.GET("/artworks", artworkHandler.GetArtworks)
	api.GET("/resume", resumeHandler.GetResume)

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/artworks", artworkHandler.CreateArtwork)
		protected.PUT("/artworks/:id", artworkHandler.UpdateArtwork)
		protected.DELETE("/artworks/:id", artworkHandler.DeleteArtwork)
		protected.POST("/resume", resumeHandler.UploadResume)
	}

	// Start the server
	log.Printf("Server starting on port %s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
