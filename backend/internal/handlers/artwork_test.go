package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"backend/internal/handlers"
	"backend/internal/models"
)

// mockArtworkService simulates the service layer
type mockArtworkService struct {
	createErr error
	deleteErr error
}

func (m *mockArtworkService) CreateArtwork(artwork models.Artwork, file *multipart.FileHeader) error {
	return m.createErr
}

func (m *mockArtworkService) GetArtworks() ([]models.Artwork, error) {
	return nil, nil
}

func (m *mockArtworkService) DeleteArtwork(id string) error {
	return m.deleteErr
}

func TestArtworkHandler_CreateArtwork_MissingTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a multipart form with missing title
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("category", "concept") // missing title
	// add a dummy file
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("dummy image content"))
	writer.Close()

	c.Request = httptest.NewRequest("POST", "/api/artworks", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.CreateArtwork(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestArtworkHandler_CreateArtwork_MissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Test Title")
	// missing file
	writer.Close()

	c.Request = httptest.NewRequest("POST", "/api/artworks", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.CreateArtwork(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestArtworkHandler_CreateArtwork_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Test Title")
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("dummy image content"))
	writer.Close()

	c.Request = httptest.NewRequest("POST", "/api/artworks", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.CreateArtwork(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}
}

func TestArtworkHandler_GetArtworks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/api/artworks", nil)

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.GetArtworks(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestArtworkHandler_DeleteArtwork_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	c.Request = httptest.NewRequest("DELETE", "/api/artworks/123", nil)

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.DeleteArtwork(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestArtworkHandler_DeleteArtwork_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	c.Request = httptest.NewRequest("DELETE", "/api/artworks/123", nil)

	// Simulate real Firestore gRPC NotFound error
	grpcErr := status.Error(codes.NotFound, "projects/my-project/databases/(default)/documents/artworks/123 not found")
	mockSvc := &mockArtworkService{deleteErr: grpcErr}
	h := handlers.NewArtworkHandler(mockSvc)

	h.DeleteArtwork(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", w.Code)
	}
}
