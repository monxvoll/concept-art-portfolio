package clients_test

import (
	"testing"

	"backend/internal/clients"
)

// We verify that the FirestoreClient struct complies with the DatabaseClient interface
func TestFirestoreClient_ImplementsDatabaseClient(t *testing.T) {
	// This will fail to compile if FirestoreClient does not exist or doesn't implement the interface.
	var _ clients.DatabaseClient = (*clients.FirestoreClient)(nil)
}
