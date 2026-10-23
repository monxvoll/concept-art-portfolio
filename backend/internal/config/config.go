// Package config loads environment variables from a .env file
// and exposes them through a typed Config struct.
// This is the single source of truth for all external service credentials.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration values loaded from environment variables.
type Config struct {
	Port                         string
	FrontendOrigin               string
	FirebaseProjectID            string
	GoogleApplicationCredentials string
	CloudinaryURL                string
}

// LoadEnv loads key-value pairs from the first matching .env file found in filenames
// (defaulting to .env and ../.env) into environment variables without overriding
// any variables that are already set.
func LoadEnv(filenames ...string) {
	if len(filenames) == 0 {
		filenames = []string{".env", "../.env"}
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Strip matching quotes if present
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
				(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				if len(val) >= 2 {
					val = val[1 : len(val)-1]
				}
			}

			// Do not overwrite already set environment variables
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
		if err := scanner.Err(); err != nil {
			return
		}
		return
	}
}

// Load loads the environment and resolves the configuration into a Config struct.
// It bridges legacy/alternative variable names and sets sensible defaults:
// - PORT defaults to "8080"
// - FRONTEND_ORIGIN defaults to "http://localhost:5500"
// - CLOUDINARY_URL falls back to CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET
func Load(envFiles ...string) (*Config, error) {
	LoadEnv(envFiles...)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5500"
	}

	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
	googleCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// Resolve Cloudinary URL
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
		if cloudName != "" && apiKey != "" && apiSecret != "" {
			cloudinaryURL = fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
			os.Setenv("CLOUDINARY_URL", cloudinaryURL)
		}
	}

	cfg := &Config{
		Port:                         port,
		FrontendOrigin:               frontendOrigin,
		FirebaseProjectID:            firebaseProjectID,
		GoogleApplicationCredentials: googleCreds,
		CloudinaryURL:                cloudinaryURL,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration values are present.
func (c *Config) Validate() error {
	if c.CloudinaryURL == "" {
		return fmt.Errorf("missing required configuration: CLOUDINARY_URL (or CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET)")
	}
	if c.FirebaseProjectID == "" {
		return fmt.Errorf("missing required configuration: FIREBASE_PROJECT_ID")
	}
	return nil
}
