package handlers

import (
        "mime/multipart"
        "net/http"

        "github.com/gin-gonic/gin"
)

type ResumeServiceInterface interface {
        UploadResume(file *multipart.FileHeader) (string, error)
        GetResumeURL() (string, error)
}

type ResumeHandler struct {
        svc ResumeServiceInterface
}

func NewResumeHandler(svc ResumeServiceInterface) *ResumeHandler {
        return &ResumeHandler{svc: svc}
}

func (h *ResumeHandler) UploadResume(c *gin.Context) {
        file, err := c.FormFile("resume")
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "resume file is required"})
                return
        }

        url, err := h.svc.UploadResume(file)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload resume"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *ResumeHandler) GetResume(c *gin.Context) {
        url, err := h.svc.GetResumeURL()
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
                return
        }
        c.JSON(http.StatusOK, gin.H{"url": url})
}
