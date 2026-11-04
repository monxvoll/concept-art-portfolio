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

	if err := s.db.SaveResumeURL(url); err != nil {
		return "", err
	}
	return url, nil
}

func (s *ResumeService) GetResumeURL() (string, error) {
	return s.db.GetResumeURL()
}
