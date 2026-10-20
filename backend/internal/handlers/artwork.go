package handlers

import (
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"backend/internal/models"
)

type ArtworkServiceInterface interface {
	CreateArtwork(artwork models.Artwork, file *multipart.FileHeader) error
	GetArtworks() ([]models.Artwork, error)
	DeleteArtwork(id string) error
	UpdateArtwork(id string, update models.ArtworkUpdate) (models.Artwork, error)
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

func (h *ArtworkHandler) GetArtworks(c *gin.Context) {
	artworks, err := h.svc.GetArtworks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artworks"})
		return
	}
	// Return an empty array instead of null if there are no artworks
	if artworks == nil {
		artworks = []models.Artwork{}
	}
	c.JSON(http.StatusOK, artworks)
}

// isNotFound reports whether err means the artwork does not exist
// (Firestore gRPC NotFound, or a plain "not found" error from lower layers).
func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound || strings.Contains(strings.ToLower(err.Error()), "not found")
}

func (h *ArtworkHandler) DeleteArtwork(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artwork id is required"})
		return
	}

	err := h.svc.DeleteArtwork(id)
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "artwork not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete artwork"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "artwork deleted"})
}

func (h *ArtworkHandler) UpdateArtwork(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artwork id is required"})
		return
	}

	var update models.ArtworkUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	update.Title = strings.TrimSpace(update.Title)
	if update.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	artwork, err := h.svc.UpdateArtwork(id, update)
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "artwork not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update artwork"})
		return
	}

	c.JSON(http.StatusOK, artwork)
}
