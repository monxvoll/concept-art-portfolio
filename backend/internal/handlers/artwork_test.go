package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/handlers"
	"backend/internal/models"
)

// mockArtworkService simulates the service layer
type mockArtworkService struct {
	createErr error
}

func (m *mockArtworkService) CreateArtwork(artwork models.Artwork, file *multipart.FileHeader) error {
	return m.createErr
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
