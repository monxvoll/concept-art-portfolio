package services

import (
	"log/slog"
	"mime/multipart"

	"backend/internal/clients"
	"backend/internal/models"
)

type ArtworkService struct {
	storage clients.StorageClient
	db      clients.DatabaseClient
}

func NewArtworkService(storage clients.StorageClient, db clients.DatabaseClient) *ArtworkService {
	return &ArtworkService{
		storage: storage,
		db:      db,
	}
}

func (s *ArtworkService) CreateArtwork(artwork models.Artwork, file *multipart.FileHeader) error {
	// Upload image to Storage
	url, publicID, err := s.storage.UploadImage(file)
	if err != nil {
		slog.Error("failed to upload image", "error", err)
		return err
	}

	// Update artwork model with storage details
	artwork.ImageURL = url
	artwork.CloudinaryID = publicID

	// Save artwork to Database
	err = s.db.SaveArtwork(artwork)
	if err != nil {
		slog.Error("failed to save artwork to database, initiating rollback", "error", err)
		// Rollback: delete the uploaded image
		rollbackErr := s.storage.DeleteImage(publicID)
		if rollbackErr != nil {
			// If rollback fails, we have an orphaned image. Log critically.
			slog.Error("CRITICAL: failed to rollback image deletion after db save failure",
				"error", rollbackErr,
				"publicID", publicID,
			)
		}
		return err
	}

	return nil
}
