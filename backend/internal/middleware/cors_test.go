package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
)

// newCORSRouter mirrors main.go: CORS is a global middleware in front of the API routes
func newCORSRouter(allowedOrigin string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CORS(allowedOrigin))
	r.PUT("/api/artworks/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/api/artworks", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestCORS_PreflightFromAllowedOrigin(t *testing.T) {
	r := newCORSRouter("http://localhost:5500")

	req := httptest.NewRequest(http.MethodOptions, "/api/artworks/123", nil)
	req.Header.Set("Origin", "http://localhost:5500")
	req.Header.Set("Access-Control-Request-Method", http.MethodPut)
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5500" {
		t.Fatalf("expected Allow-Origin http://localhost:5500, got %q", got)
	}
	methods := w.Header().Get("Access-Control-Allow-Methods")
	for _, m := range []string{"GET", "POST", "PUT", "DELETE"} {
		if !strings.Contains(methods, m) {
			t.Errorf("expected Allow-Methods to contain %s, got %q", m, methods)
		}
	}
	headers := strings.ToLower(w.Header().Get("Access-Control-Allow-Headers"))
	for _, h := range []string{"authorization", "content-type"} {
		if !strings.Contains(headers, h) {
			t.Errorf("expected Allow-Headers to contain %s, got %q", h, headers)
		}
	}
}

func TestCORS_SimpleRequestFromAllowedOrigin(t *testing.T) {
	r := newCORSRouter("http://localhost:5500")

	req := httptest.NewRequest(http.MethodGet, "/api/artworks", nil)
	req.Header.Set("Origin", "http://localhost:5500")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5500" {
		t.Fatalf("expected Allow-Origin http://localhost:5500, got %q", got)
	}
}

func TestCORS_RejectsOtherOrigins(t *testing.T) {
	r := newCORSRouter("http://localhost:5500")

	req := httptest.NewRequest(http.MethodGet, "/api/artworks", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Allow-Origin for a foreign origin, got %q", got)
	}
}
