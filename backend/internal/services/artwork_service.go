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

func (s *ArtworkService) GetArtworks() ([]models.Artwork, error) {
	return s.db.GetArtworks()
}

func (s *ArtworkService) DeleteArtwork(id string) error {
	artwork, err := s.db.GetArtworkByID(id)
	if err != nil {
		slog.Error("failed to get artwork from database", "error", err)
		return err
	}

	err = s.db.DeleteArtwork(id)
	if err != nil {
		slog.Error("failed to delete artwork from database", "error", err)
		return err
	}

	if artwork.CloudinaryID != "" {
		storageErr := s.storage.DeleteImage(artwork.CloudinaryID)
		if storageErr != nil {
			slog.Error("CRITICAL: failed to delete image from storage", "error", storageErr, "publicID", artwork.CloudinaryID)
		}
	}
	return nil
}

// UpdateArtwork overwrites the editable fields of an existing artwork and
// returns the resulting artwork. Image, slug and identifiers are preserved.
func (s *ArtworkService) UpdateArtwork(id string, update models.ArtworkUpdate) (models.Artwork, error) {
	artwork, err := s.db.GetArtworkByID(id)
	if err != nil {
		slog.Error("failed to get artwork from database", "error", err)
		return models.Artwork{}, err
	}

	if err := s.db.UpdateArtwork(id, update); err != nil {
		slog.Error("failed to update artwork in database", "error", err)
		return models.Artwork{}, err
	}

	artwork.Title = update.Title
	artwork.Description = update.Description
	artwork.Category = update.Category
	artwork.ArtworkDate = update.ArtworkDate
	return artwork, nil
}
