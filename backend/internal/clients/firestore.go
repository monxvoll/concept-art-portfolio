package clients

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/internal/models"
)

type FirestoreClient struct {
	client *firestore.Client
	ctx    context.Context
}

// NewFirestoreClient initializes a new Firestore client.
func NewFirestoreClient(ctx context.Context, projectID string) (*FirestoreClient, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// SaveArtwork saves an artwork document to Firestore.
func (f *FirestoreClient) SaveArtwork(artwork models.Artwork) error {
	_, err := f.client.Collection("artworks").Doc(artwork.ID).Set(f.ctx, artwork)
	return err
}

// GetArtworks retrieves all artwork documents from Firestore.
func (f *FirestoreClient) GetArtworks() ([]models.Artwork, error) {
	var artworks []models.Artwork
	iter := f.client.Collection("artworks").Documents(f.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var artwork models.Artwork
		if err := doc.DataTo(&artwork); err != nil {
			return nil, err
		}
		artworks = append(artworks, artwork)
	}

	return artworks, nil
}

// DeleteArtwork deletes an artwork document from Firestore by ID.
func (f *FirestoreClient) DeleteArtwork(id string) error {
	_, err := f.client.Collection("artworks").Doc(id).Delete(f.ctx)
	return err
}

// UpdateArtwork overwrites only the editable fields of an artwork document.
// Firestore returns a NotFound error if the document does not exist.
func (f *FirestoreClient) UpdateArtwork(id string, update models.ArtworkUpdate) error {
	_, err := f.client.Collection("artworks").Doc(id).Update(f.ctx, []firestore.Update{
		{Path: "title", Value: update.Title},
		{Path: "description", Value: update.Description},
		{Path: "category", Value: update.Category},
		{Path: "artwork_date", Value: update.ArtworkDate},
	})
	return err
}

// GetArtworkByID retrieves a single artwork document from Firestore by ID.
func (f *FirestoreClient) GetArtworkByID(id string) (models.Artwork, error) {
	doc, err := f.client.Collection("artworks").Doc(id).Get(f.ctx)
	if err != nil {
		return models.Artwork{}, err
	}
	var artwork models.Artwork
	err = doc.DataTo(&artwork)
	return artwork, err
}
