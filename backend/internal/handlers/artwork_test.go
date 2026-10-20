package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
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
	updateErr error

	// Recorded arguments of the last UpdateArtwork call
	updateCalled bool
	updatedID    string
	updateInput  models.ArtworkUpdate
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

func (m *mockArtworkService) UpdateArtwork(id string, update models.ArtworkUpdate) (models.Artwork, error) {
	m.updateCalled = true
	m.updatedID = id
	m.updateInput = update
	if m.updateErr != nil {
		return models.Artwork{}, m.updateErr
	}
	return models.Artwork{
		ID:          id,
		Title:       update.Title,
		Description: update.Description,
		Category:    update.Category,
		ArtworkDate: update.ArtworkDate,
	}, nil
}

// newPutContext builds a gin test context for PUT /api/artworks/:id with a raw JSON body
func newPutContext(id, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: id}}
	c.Request = httptest.NewRequest("PUT", "/api/artworks/"+id, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestArtworkHandler_UpdateArtwork_Success(t *testing.T) {
	c, w := newPutContext("123", `{"title":"New Title","description":"New desc","category":"Environment","artwork_date":"2026-01-02"}`)

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.UpdateArtwork(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", w.Code, w.Body.String())
	}

	want := models.ArtworkUpdate{Title: "New Title", Description: "New desc", Category: "Environment", ArtworkDate: "2026-01-02"}
	if mockSvc.updatedID != "123" || mockSvc.updateInput != want {
		t.Fatalf("service got id=%q input=%+v, want id=%q input=%+v", mockSvc.updatedID, mockSvc.updateInput, "123", want)
	}

	var got models.Artwork
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not an artwork JSON: %v", err)
	}
	if got.ID != "123" || got.Title != "New Title" {
		t.Fatalf("expected updated artwork in response, got %+v", got)
	}
}

func TestArtworkHandler_UpdateArtwork_NotFound(t *testing.T) {
	c, w := newPutContext("missing", `{"title":"New Title"}`)

	grpcErr := status.Error(codes.NotFound, "projects/my-project/databases/(default)/documents/artworks/missing not found")
	mockSvc := &mockArtworkService{updateErr: grpcErr}
	h := handlers.NewArtworkHandler(mockSvc)

	h.UpdateArtwork(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", w.Code)
	}
}

func TestArtworkHandler_UpdateArtwork_MissingTitle(t *testing.T) {
	c, w := newPutContext("123", `{"title":"   ","description":"desc"}`)

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.UpdateArtwork(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
	if mockSvc.updateCalled {
		t.Fatal("service must not be called when title is missing")
	}
}

func TestArtworkHandler_UpdateArtwork_InvalidJSON(t *testing.T) {
	c, w := newPutContext("123", `{"title": `)

	mockSvc := &mockArtworkService{}
	h := handlers.NewArtworkHandler(mockSvc)

	h.UpdateArtwork(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
	if mockSvc.updateCalled {
		t.Fatal("service must not be called when body is malformed")
	}
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
