package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
)

// mockAuthClient implements the verify method we need
type mockAuthClient struct {
	verifyErr error
}

func (m *mockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	return &auth.Token{UID: "test-uid"}, nil
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	mockAuth := &mockAuthClient{}
	mw := middleware.FirebaseAuth(mockAuth)

	mw(c)

	if c.IsAborted() {
		t.Fatal("expected request to proceed, but it was aborted")
	}
	if uid, exists := c.Get("UID"); !exists || uid != "test-uid" {
		t.Fatalf("expected UID to be test-uid, got %v", uid)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	mockAuth := &mockAuthClient{verifyErr: errors.New("invalid token")}
	mw := middleware.FirebaseAuth(mockAuth)

	mw(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
