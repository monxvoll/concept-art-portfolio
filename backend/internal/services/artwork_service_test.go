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
	saveErr    error
	saved      []models.Artwork
	getErr     error
	artworks   []models.Artwork
	deleteErr  error
	deletedIDs []string
	updateErr  error
	updatedIDs []string
	updates    []models.ArtworkUpdate
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
	m.deletedIDs = append(m.deletedIDs, id)
	return m.deleteErr
}

func (m *mockDatabaseClient) UpdateArtwork(id string, update models.ArtworkUpdate) error {
	m.updatedIDs = append(m.updatedIDs, id)
	m.updates = append(m.updates, update)
	return m.updateErr
}

func TestArtworkService_UpdateArtwork_Success(t *testing.T) {
	db := &mockDatabaseClient{
		artworks: []models.Artwork{
			{ID: "123", Title: "Old", Slug: "old-slug", ImageURL: "url", CloudinaryID: "pid", Project: "P"},
		},
	}
	svc := services.NewArtworkService(&mockStorageClient{}, db)

	update := models.ArtworkUpdate{Title: "New", Description: "Desc", Category: "Environment", ArtworkDate: "2026-01-02"}
	got, err := svc.UpdateArtwork("123", update)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(db.updatedIDs) != 1 || db.updatedIDs[0] != "123" || db.updates[0] != update {
		t.Fatalf("expected db update of '123' with %+v, got ids=%v updates=%+v", update, db.updatedIDs, db.updates)
	}

	want := models.Artwork{
		ID: "123", Title: "New", Description: "Desc", Category: "Environment", ArtworkDate: "2026-01-02",
		Slug: "old-slug", ImageURL: "url", CloudinaryID: "pid", Project: "P",
	}
	if got != want {
		t.Fatalf("expected merged artwork %+v, got %+v", want, got)
	}
}

func TestArtworkService_UpdateArtwork_NotFound(t *testing.T) {
	db := &mockDatabaseClient{getErr: errors.New("not found")}
	svc := services.NewArtworkService(&mockStorageClient{}, db)

	_, err := svc.UpdateArtwork("missing", models.ArtworkUpdate{Title: "New"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(db.updatedIDs) != 0 {
		t.Fatalf("db update must not run for a missing artwork, got %v", db.updatedIDs)
	}
}

func TestArtworkService_UpdateArtwork_DBFailure(t *testing.T) {
	db := &mockDatabaseClient{
		artworks:  []models.Artwork{{ID: "123", Title: "Old"}},
		updateErr: errors.New("firestore down"),
	}
	svc := services.NewArtworkService(&mockStorageClient{}, db)

	_, err := svc.UpdateArtwork("123", models.ArtworkUpdate{Title: "New"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
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

func TestArtworkService_GetArtworks_Success(t *testing.T) {
	db := &mockDatabaseClient{
		artworks: []models.Artwork{
			{ID: "1", Title: "Art 1"},
			{ID: "2", Title: "Art 2"},
		},
	}
	svc := services.NewArtworkService(nil, db)

	artworks, err := svc.GetArtworks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(artworks) != 2 {
		t.Fatalf("expected 2 artworks, got %d", len(artworks))
	}
}

func TestArtworkService_DeleteArtwork_Success(t *testing.T) {
	storage := &mockStorageClient{}
	db := &mockDatabaseClient{
		artworks: []models.Artwork{
			{ID: "123", CloudinaryID: "pid123"},
		},
	}
	svc := services.NewArtworkService(storage, db)

	err := svc.DeleteArtwork("123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(storage.deletedIDs) == 0 || storage.deletedIDs[0] != "pid123" {
		t.Errorf("expected storage.DeleteImage to be called with 'pid123', got %v", storage.deletedIDs)
	}
}

func TestArtworkService_DeleteArtwork_NotFound(t *testing.T) {
	storage := &mockStorageClient{}
	db := &mockDatabaseClient{getErr: errors.New("not found")}
	svc := services.NewArtworkService(storage, db)

	err := svc.DeleteArtwork("invalid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestArtworkService_DeleteArtwork_StorageFails_DBStillDeleted(t *testing.T) {
	storage := &mockStorageClient{deleteErr: errors.New("cloudinary down")}
	db := &mockDatabaseClient{
		artworks: []models.Artwork{
			{ID: "123", CloudinaryID: "pid123"},
		},
	}

	svc := services.NewArtworkService(storage, db)

	err := svc.DeleteArtwork("123")

	if err != nil {
		t.Fatalf("no error expected, got: %v", err)
	}

	if len(db.deletedIDs) == 0 || db.deletedIDs[0] != "123" {
		t.Fatalf("se esperaba que la DB borrara '123', pero la libreta tiene: %v", db.deletedIDs)
	}
}
