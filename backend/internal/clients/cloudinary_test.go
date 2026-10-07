package clients_test

import (
	"testing"

	"backend/internal/clients"
)

// We verify that the CloudinaryClient struct complies with the StorageClient interface
func TestCloudinaryClient_ImplementsStorageClient(t *testing.T) {
	// This will fail to compile if CloudinaryClient does not exist or doesn't implement the interface.
	var _ clients.StorageClient = (*clients.CloudinaryClient)(nil)
}
