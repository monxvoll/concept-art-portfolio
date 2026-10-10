package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"

	"backend/internal/clients"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"
)

func main() {
	ctx := context.Background()

	// Init Cloudinary
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	storageClient, err := clients.NewCloudinaryClient(cloudinaryURL)
	if err != nil {
		slog.Error("failed to init cloudinary", "err", err)
	}

	// Init Firestore
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	dbClient, err := clients.NewFirestoreClient(ctx, projectID)
	if err != nil {
		slog.Error("failed to init firestore", "err", err)
	}

	// Init Firebase Auth for Middleware
	app, err := firebase.NewApp(ctx, nil)
	var authMiddleware gin.HandlerFunc
	if err == nil {
		authClient, err := app.Auth(ctx)
		if err == nil {
			authMiddleware = middleware.FirebaseAuth(authClient)
		} else {
			slog.Error("failed to init firebase auth", "err", err)
		}
	} else {
		slog.Error("failed to init firebase app", "err", err)
	}

	if authMiddleware == nil {
		// Mock auth for now if firebase fails to load
		authMiddleware = func(c *gin.Context) { c.Next() }
	}

	// Wire up layers
	artworkSvc := services.NewArtworkService(storageClient, dbClient)
	artworkHandler := handlers.NewArtworkHandler(artworkSvc)

	r := gin.Default()

	// Serves the entire frontend folder at the root path "/"
	r.StaticFS("/frontend", http.Dir("../frontend"))

	// Routes
	api := r.Group("/api")

	// Public routes
	api.GET("/artworks", artworkHandler.GetArtworks)

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/artworks", artworkHandler.CreateArtwork)
		protected.DELETE("/artworks/:id", artworkHandler.DeleteArtwork)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
