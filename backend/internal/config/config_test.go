package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"backend/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	// Isolate test from current env vars
	os.Unsetenv("PORT")
	os.Unsetenv("FRONTEND_ORIGIN")
	os.Unsetenv("FIREBASE_PROJECT_ID")
	os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	os.Unsetenv("CLOUDINARY_URL")
	os.Unsetenv("CLOUDINARY_CLOUD_NAME")
	os.Unsetenv("CLOUDINARY_API_KEY")
	os.Unsetenv("CLOUDINARY_API_SECRET")

	// Set required vars so validation passes while testing defaults
	os.Setenv("FIREBASE_PROJECT_ID", "test-project")
	os.Setenv("CLOUDINARY_URL", "cloudinary://test")
	defer func() {
		os.Unsetenv("FIREBASE_PROJECT_ID")
		os.Unsetenv("CLOUDINARY_URL")
	}()

	cfg, err := config.Load("/nonexistent_env_file")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default Port '8080', got %q", cfg.Port)
	}
	if cfg.FrontendOrigin != "http://localhost:5500" {
		t.Errorf("expected default FrontendOrigin 'http://localhost:5500', got %q", cfg.FrontendOrigin)
	}
}

func TestLoad_Validation(t *testing.T) {
	os.Unsetenv("FIREBASE_PROJECT_ID")
	os.Unsetenv("CLOUDINARY_URL")
	os.Unsetenv("CLOUDINARY_CLOUD_NAME")
	os.Unsetenv("CLOUDINARY_API_KEY")
	os.Unsetenv("CLOUDINARY_API_SECRET")

	// Missing CloudinaryURL only (Firebase present)
	os.Setenv("FIREBASE_PROJECT_ID", "test-project")
	_, err := config.Load("/nonexistent_env_file")
	if err == nil {
		t.Fatal("expected error when CLOUDINARY_URL is missing, got nil")
	}

	// Missing FirebaseProjectID only (Cloudinary present)
	os.Unsetenv("FIREBASE_PROJECT_ID")
	os.Setenv("CLOUDINARY_URL", "cloudinary://test")
	defer os.Unsetenv("CLOUDINARY_URL")
	_, err = config.Load("/nonexistent_env_file")
	if err == nil {
		t.Fatal("expected error when FIREBASE_PROJECT_ID is missing, got nil")
	}
}

func TestLoad_EnvFile(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env.test")

	content := `
# Comment line
PORT=9000
FRONTEND_ORIGIN=https://myportfolio.com
FIREBASE_PROJECT_ID="test-project"
CLOUDINARY_URL='cloudinary://key:secret@cloud'
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	os.Unsetenv("PORT")
	os.Unsetenv("FRONTEND_ORIGIN")
	os.Unsetenv("FIREBASE_PROJECT_ID")
	os.Unsetenv("CLOUDINARY_URL")

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "9000" {
		t.Errorf("expected Port '9000', got %q", cfg.Port)
	}
	if cfg.FrontendOrigin != "https://myportfolio.com" {
		t.Errorf("expected FrontendOrigin 'https://myportfolio.com', got %q", cfg.FrontendOrigin)
	}
	if cfg.FirebaseProjectID != "test-project" {
		t.Errorf("expected FirebaseProjectID 'test-project', got %q", cfg.FirebaseProjectID)
	}
	if cfg.CloudinaryURL != "cloudinary://key:secret@cloud" {
		t.Errorf("expected CloudinaryURL 'cloudinary://key:secret@cloud', got %q", cfg.CloudinaryURL)
	}
}

func TestLoad_CloudinaryPartsFallback(t *testing.T) {
	os.Unsetenv("CLOUDINARY_URL")
	os.Setenv("CLOUDINARY_CLOUD_NAME", "mycloud")
	os.Setenv("CLOUDINARY_API_KEY", "mykey")
	os.Setenv("CLOUDINARY_API_SECRET", "mysecret")
	os.Setenv("FIREBASE_PROJECT_ID", "test-project")
	defer func() {
		os.Unsetenv("CLOUDINARY_CLOUD_NAME")
		os.Unsetenv("CLOUDINARY_API_KEY")
		os.Unsetenv("CLOUDINARY_API_SECRET")
		os.Unsetenv("CLOUDINARY_URL")
		os.Unsetenv("FIREBASE_PROJECT_ID")
	}()

	cfg, err := config.Load("/nonexistent_env_file")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "cloudinary://mykey:mysecret@mycloud"
	if cfg.CloudinaryURL != expected {
		t.Errorf("expected %q, got %q", expected, cfg.CloudinaryURL)
	}
	if os.Getenv("CLOUDINARY_URL") != expected {
		t.Errorf("expected env CLOUDINARY_URL to be set to %q, got %q", expected, os.Getenv("CLOUDINARY_URL"))
	}
}
