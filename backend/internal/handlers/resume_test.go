package handlers_test

import (
        "bytes"
        "mime/multipart"
        "net/http"
        "net/http/httptest"
        "testing"

        "github.com/gin-gonic/gin"
        "backend/internal/handlers"
)

// mockResumeService is a stub for the ResumeServiceInterface
type mockResumeService struct {
        uploadErr error
        fileURL   string
}

func (m *mockResumeService) UploadResume(file *multipart.FileHeader) (string, error) {
        if m.uploadErr != nil {
                return "", m.uploadErr
        }
        return m.fileURL, nil
}

func (m *mockResumeService) GetResumeURL() (string, error) {
        return m.fileURL, nil
}

func TestResumeHandler_UploadResume_Success(t *testing.T) {
        gin.SetMode(gin.TestMode)
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        part, _ := writer.CreateFormFile("resume", "cv.pdf")
        part.Write([]byte("dummy pdf content"))
        writer.Close()

        c.Request = httptest.NewRequest("POST", "/api/resume", body)
        c.Request.Header.Set("Content-Type", writer.FormDataContentType())

        mockSvc := &mockResumeService{
                fileURL: "https://cloudinary.com/dummy/cv.pdf",
        }
        h := handlers.NewResumeHandler(mockSvc)

        h.UploadResume(c)

        if w.Code != http.StatusOK {
                t.Fatalf("expected 200 OK, got %d", w.Code)
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
