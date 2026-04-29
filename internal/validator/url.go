package validator

import (
	"fmt"
	"net/url"
)

// ValidateURL validates and normalizes a URL
func ValidateURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Require scheme
	if u.Scheme == "" {
		return "", fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	// Validate scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("URL scheme must be http or https")
	}

	// Require host
	if u.Host == "" {
		return "", fmt.Errorf("URL must include host")
	}

	return u.String(), nil
}

// IsValidURL checks if a string is a valid URL without normalizing
func IsValidURL(rawURL string) bool {
	_, err := ValidateURL(rawURL)
	return err == nil
}
