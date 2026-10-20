package services

import (
	"mime/multipart"

	"backend/internal/clients"
)

type ResumeService struct {
	storage clients.StorageClient
	db      clients.ResumeDatabaseClient
}

func NewResumeService(storage clients.StorageClient, db clients.ResumeDatabaseClient) *ResumeService {
	return &ResumeService{
		storage: storage,
		db:      db,
	}
}

func (s *ResumeService) UploadResume(file *multipart.FileHeader) (string, error) {
	url, _, err := s.storage.UploadImage(file)
	if err != nil {
		return "", err
	}

	if s.db != nil {
		_ = s.db.SaveResumeURL(url)
	}
	return url, nil
}

func (s *ResumeService) GetResumeURL() (string, error) {
	if s.db == nil {
		return "", nil
	}
	return s.db.GetResumeURL()
}

