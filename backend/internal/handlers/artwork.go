package handlers

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/models"
)

type ArtworkServiceInterface interface {
	CreateArtwork(artwork models.Artwork, file *multipart.FileHeader) error
}

type ArtworkHandler struct {
	svc ArtworkServiceInterface
}

func NewArtworkHandler(svc ArtworkServiceInterface) *ArtworkHandler {
	return &ArtworkHandler{svc: svc}
}

func (h *ArtworkHandler) CreateArtwork(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	artwork := models.Artwork{
		ID:          uuid.New().String(),
		Title:       title,
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Project:     c.PostForm("project"),
		CreatedAt:   time.Now(),
		ArtworkDate: c.PostForm("artwork_date"),
	}

	err = h.svc.CreateArtwork(artwork, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create artwork"})
		return
	}

	c.JSON(http.StatusCreated, artwork)
}
