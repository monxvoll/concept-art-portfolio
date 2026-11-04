package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const pdfMagic = "%PDF-"

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

	ok, err := isPDF(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read resume file"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume must be a PDF file"})
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

// isPDF reports whether the uploaded file has a .pdf extension and
// starts with the PDF magic bytes.
func isPDF(file *multipart.FileHeader) (bool, error) {
	if !strings.EqualFold(filepath.Ext(file.Filename), ".pdf") {
		return false, nil
	}

	f, err := file.Open()
	if err != nil {
		return false, err
	}
	defer f.Close()

	header := make([]byte, len(pdfMagic))
	if _, err := io.ReadFull(f, header); err != nil {
		// Files shorter than the magic header are not PDFs.
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return false, nil
		}
		return false, err
	}
	return string(header) == pdfMagic, nil
}
