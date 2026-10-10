package services_test

import (
	"errors"
	"mime/multipart"
	"testing"

	"backend/internal/clients"
	"backend/internal/models"
	"backend/internal/services"
)

// Mock StorageClient
var _ clients.StorageClient = (*mockStorageClient)(nil)
type mockStorageClient struct {
	uploadErr  error
	deleteErr  error
	uploadURL  string
	uploadID   string
	deletedIDs []string
}

func (m *mockStorageClient) UploadImage(file *multipart.FileHeader) (string, string, error) {
	return m.uploadURL, m.uploadID, m.uploadErr
}

func (m *mockStorageClient) DeleteImage(publicID string) error {
	m.deletedIDs = append(m.deletedIDs, publicID)
	return m.deleteErr
}

// Mock DatabaseClient
var _ clients.DatabaseClient = (*mockDatabaseClient)(nil)
type mockDatabaseClient struct {
	saveErr   error
	saved     []models.Artwork
	getErr    error
	artworks  []models.Artwork
	deleteErr error
}

func (m *mockDatabaseClient) SaveArtwork(artwork models.Artwork) error {
	m.saved = append(m.saved, artwork)
	return m.saveErr
}

func (m *mockDatabaseClient) GetArtworks() ([]models.Artwork, error) {
	return m.artworks, m.getErr
}

func (m *mockDatabaseClient) GetArtworkByID(id string) (models.Artwork, error) {
	if m.getErr != nil {
		return models.Artwork{}, m.getErr
	}
	for _, art := range m.saved {
		if art.ID == id {
			return art, nil
		}
	}
	for _, art := range m.artworks {
		if art.ID == id {
			return art, nil
		}
	}
	return models.Artwork{}, errors.New("not found")
}

func (m *mockDatabaseClient) DeleteArtwork(id string) error {
	return m.deleteErr
}

func TestArtworkService_CreateArtwork_Success(t *testing.T) {
	storage := &mockStorageClient{uploadURL: "url", uploadID: "pid"}
	db := &mockDatabaseClient{}
	svc := services.NewArtworkService(storage, db)

	err := svc.CreateArtwork(models.Artwork{ID: "123"}, &multipart.FileHeader{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(db.saved) == 0 {
		t.Fatal("expected artwork to be saved to db")
	}
	if db.saved[0].CloudinaryID != "pid" || db.saved[0].ImageURL != "url" {
		t.Errorf("expected storage details to be set on artwork")
	}
}

func TestArtworkService_CreateArtwork_StorageFailure(t *testing.T) {
	storage := &mockStorageClient{uploadErr: errors.New("upload failed")}
	db := &mockDatabaseClient{}
	svc := services.NewArtworkService(storage, db)

	err := svc.CreateArtwork(models.Artwork{ID: "123"}, &multipart.FileHeader{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(db.saved) > 0 {
		t.Fatal("expected db save to not be called on storage failure")
	}
}

func TestArtworkService_CreateArtwork_DBFailureRollback(t *testing.T) {
	storage := &mockStorageClient{uploadURL: "url", uploadID: "pid"}
	db := &mockDatabaseClient{saveErr: errors.New("db save failed")}
	svc := services.NewArtworkService(storage, db)

	err := svc.CreateArtwork(models.Artwork{ID: "123"}, &multipart.FileHeader{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(storage.deletedIDs) == 0 || storage.deletedIDs[0] != "pid" {
		t.Fatalf("expected rollback deletion of image, deleted: %v", storage.deletedIDs)
	}
}
