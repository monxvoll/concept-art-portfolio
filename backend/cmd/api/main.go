package main

import (
        "context"
        "log"
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
                log.Fatalf("failed to init cloudinary: %v", err)
        }

        // Init Firestore
        projectID := os.Getenv("FIREBASE_PROJECT_ID")
        dbClient, err := clients.NewFirestoreClient(ctx, projectID)
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
        frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
        if frontendOrigin == "" {
                frontendOrigin = "http://localhost:5500"
        }
        r.Use(middleware.CORS(frontendOrigin))

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
        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }
        log.Printf("Server starting on port %s...\n", port)
        if err := r.Run(":" + port); err != nil {
                log.Fatalf("Server failed to start: %v", err)
        }
}
