package models

import "time"

// Bookmark represents a saved bookmark with metadata
type Bookmark struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	Excerpt    string    `json:"excerpt"`
	Author     string    `json:"author"`
	Public     bool      `json:"public"`
	Content    string    `json:"content"`
	HTML       string    `json:"html"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	HasContent bool      `json:"has_content"`
	Tags       []Tag     `json:"tags,omitempty"`
}
