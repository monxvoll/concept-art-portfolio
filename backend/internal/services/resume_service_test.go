package services_test

import (
	"errors"
	"mime/multipart"
	"testing"

	"backend/internal/clients"
	"backend/internal/services"
)

// Mock ResumeDatabaseClient
var _ clients.ResumeDatabaseClient = (*mockResumeDatabaseClient)(nil)

type mockResumeDatabaseClient struct {
	savedURL string
	saveErr  error
	getURL   string
	getErr   error
}

func (m *mockResumeDatabaseClient) SaveResumeURL(url string) error {
	m.savedURL = url
	return m.saveErr
}

func (m *mockResumeDatabaseClient) GetResumeURL() (string, error) {
	return m.getURL, m.getErr
}

func TestResumeService_UploadResume_Success(t *testing.T) {
	mockStorage := &mockStorageClient{
		uploadURL: "https://cloudinary.com/dummy/cv.pdf",
	}
	mockDB := &mockResumeDatabaseClient{}

	svc := services.NewResumeService(mockStorage, mockDB)

	header := &multipart.FileHeader{Filename: "cv.pdf"}
	url, err := svc.UploadResume(header)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "https://cloudinary.com/dummy/cv.pdf" {
		t.Fatalf("expected url 'https://cloudinary.com/dummy/cv.pdf', got '%s'", url)
	}

	if mockDB.savedURL != url {
		t.Fatalf("expected db to save url '%s', got '%s'", url, mockDB.savedURL)
	}
}

func TestResumeService_UploadResume_StorageError(t *testing.T) {
	mockStorage := &mockStorageClient{
		uploadErr: errors.New("upload failed"),
	}
	mockDB := &mockResumeDatabaseClient{}

	svc := services.NewResumeService(mockStorage, mockDB)

	header := &multipart.FileHeader{Filename: "cv.pdf"}
	_, err := svc.UploadResume(header)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if mockDB.savedURL != "" {
		t.Fatalf("expected db not to be called, got '%s'", mockDB.savedURL)
	}
}

func TestResumeService_GetResumeURL_Success(t *testing.T) {
	mockDB := &mockResumeDatabaseClient{
		getURL: "https://cloudinary.com/dummy/cv.pdf",
	}

	svc := services.NewResumeService(nil, mockDB)

	url, err := svc.GetResumeURL()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "https://cloudinary.com/dummy/cv.pdf" {
		t.Fatalf("expected url 'https://cloudinary.com/dummy/cv.pdf', got '%s'", url)
	}
}

func TestResumeService_GetResumeURL_Error(t *testing.T) {
	mockDB := &mockResumeDatabaseClient{
		getErr: errors.New("not found"),
	}

	svc := services.NewResumeService(nil, mockDB)

	_, err := svc.GetResumeURL()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
