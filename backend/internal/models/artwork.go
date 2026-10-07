// Package models defines the data structures (structs) that represent
// the core entities of the application: Artwork and Resume.
// These structs mirror the Firestore document schema defined in the data model diagram.
package models

import "time"

type Artwork struct {
	ID           string    `json:"id" firestore:"id"`
	Title        string    `json:"title" firestore:"title"`
	Slug         string    `json:"slug" firestore:"slug"`
	Description  string    `json:"description" firestore:"description"`
	Category     string    `json:"category" firestore:"category"`
	Project      string    `json:"project,omitempty" firestore:"project,omitempty"`
	ImageURL     string    `json:"image_url" firestore:"image_url"`
	Width        int       `json:"width" firestore:"width"`
	Height       int       `json:"height" firestore:"height"`
	CloudinaryID string    `json:"cloudinary_id" firestore:"cloudinary_id"`
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	ArtworkDate  string    `json:"artwork_date" firestore:"artwork_date"`
}
