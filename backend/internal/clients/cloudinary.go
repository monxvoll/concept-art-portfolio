package clients

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryClient initializes a new Cloudinary client using the Cloudinary URL.
func NewCloudinaryClient(url string) (*CloudinaryClient, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	return &CloudinaryClient{cld: cld}, nil
}

// UploadImage uploads a file to Cloudinary and returns its secure URL and public ID.
func (c *CloudinaryClient) UploadImage(file *multipart.FileHeader) (string, string, error) {
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	ctx := context.Background()
	resp, err := c.cld.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: "concept-art-portfolio",
	})
	if err != nil {
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

// DeleteImage deletes an image from Cloudinary using its public ID.
func (c *CloudinaryClient) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
