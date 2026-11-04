package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// mockResumeService is a stub for the ResumeServiceInterface
type mockResumeService struct {
	uploadErr error
	fileURL   string
	called    bool
}

func (m *mockResumeService) UploadResume(file *multipart.FileHeader) (string, error) {
	m.called = true
	if m.uploadErr != nil {
		return "", m.uploadErr
	}
	return m.fileURL, nil
}

func (m *mockResumeService) GetResumeURL() (string, error) {
	return m.fileURL, nil
}

// performResumeUpload sends a multipart request with the given file in the
// "resume" field and returns the recorded response.
func performResumeUpload(t *testing.T, svc *mockResumeService, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("resume", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write(content)
	writer.Close()

	c.Request = httptest.NewRequest("POST", "/api/resume", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	handlers.NewResumeHandler(svc).UploadResume(c)
	return w
}

func TestResumeHandler_UploadResume_Success(t *testing.T) {
	mockSvc := &mockResumeService{
		fileURL: "https://cloudinary.com/dummy/cv.pdf",
	}

	w := performResumeUpload(t, mockSvc, "cv.pdf", []byte("%PDF-1.4 dummy"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestResumeHandler_UploadResume_WrongExtension(t *testing.T) {
	mockSvc := &mockResumeService{}

	w := performResumeUpload(t, mockSvc, "photo.png", []byte("%PDF-1.4 dummy"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
	if mockSvc.called {
		t.Fatal("service must not be called for a non-PDF file")
	}
}

func TestResumeHandler_UploadResume_NotPDFContent(t *testing.T) {
	mockSvc := &mockResumeService{}

	w := performResumeUpload(t, mockSvc, "cv.pdf", []byte("\x89PNG\r\n\x1a\nnot a pdf"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
	if mockSvc.called {
		t.Fatal("service must not be called for a file without PDF magic bytes")
	}
}

func TestResumeHandler_UploadResume_InvalidFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close() // Missing 'resume' file part

	c.Request = httptest.NewRequest("POST", "/api/resume", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	mockSvc := &mockResumeService{}
	h := handlers.NewResumeHandler(mockSvc)

	h.UploadResume(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}
