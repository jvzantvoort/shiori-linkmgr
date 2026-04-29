package display

import (
	"fmt"
)

// Error prints an error message with formatting
func Error(msg string, args ...interface{}) {
	fmt.Printf("Error: %s\n", fmt.Sprintf(msg, args...))
}

// ErrorWithSuggestion prints an error with a helpful suggestion
func ErrorWithSuggestion(msg, suggestion string) {
	fmt.Printf("Error: %s\n", msg)
	if suggestion != "" {
		fmt.Printf("Suggestion: %s\n", suggestion)
	}
}

// Success prints a success message
func Success(msg string, args ...interface{}) {
	fmt.Printf("✓ %s\n", fmt.Sprintf(msg, args...))
}

// Info prints an informational message
func Info(msg string, args ...interface{}) {
	fmt.Printf("ℹ %s\n", fmt.Sprintf(msg, args...))
}
