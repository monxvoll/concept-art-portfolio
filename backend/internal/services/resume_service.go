package services

import (
        "context"
        "mime/multipart"

        "backend/internal/clients"
)

type ResumeService struct {
        storage clients.StorageClient
        db      clients.DBClient
}

func NewResumeService(storage clients.StorageClient, db clients.DBClient) *ResumeService {
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
        
        ctx := context.Background()
        // Our DBClient interface doesn't have an Upsert, but Create maps to Set in firestore.go
        data := map[string]interface{}{"url": url}
        _ = s.db.Create(ctx, "config", "resume", data) // Ignore error, simple overwrite
        return url, nil
}

func (s *ResumeService) GetResumeURL() (string, error) {
        ctx := context.Background()
        doc, err := s.db.Get(ctx, "config", "resume")
        if err != nil {
                return "", err
        }
        if url, ok := doc["url"].(string); ok {
                return url, nil
        }
        return "", nil
}
