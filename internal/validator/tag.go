package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateTagName validates a tag name
func ValidateTagName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if len(name) > 250 {
		return fmt.Errorf("tag name cannot exceed 250 characters")
	}

	// Check for invalid characters (allow alphanumeric, dash, underscore)
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return fmt.Errorf("tag name can only contain letters, numbers, dashes, and underscores")
		}
	}

	return nil
}

// NormalizeTagName normalizes a tag name (trim and lowercase)
func NormalizeTagName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// NormalizeTags parses and normalizes a comma-separated list of tags
func NormalizeTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}

	parts := strings.Split(tagsStr, ",")
	tags := make([]string, 0, len(parts))

	for _, tag := range parts {
		normalized := NormalizeTagName(tag)
		if normalized != "" {
			tags = append(tags, normalized)
		}
	}

	return tags
}
