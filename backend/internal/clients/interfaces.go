package clients

import (
	"mime/multipart"

	"backend/internal/models"
)

type StorageClient interface {
	UploadImage(file *multipart.FileHeader) (url string, publicID string, err error)
	DeleteImage(publicID string) error
}

type DatabaseClient interface {
	SaveArtwork(artwork models.Artwork) error
	GetArtworks() ([]models.Artwork, error)
	DeleteArtwork(id string) error
}
