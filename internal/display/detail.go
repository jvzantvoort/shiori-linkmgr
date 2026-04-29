package display

import (
	"fmt"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/models"
)

// FormatBookmarkDetail formats a bookmark with full details
func FormatBookmarkDetail(bookmark *models.Bookmark) {
	fmt.Printf("Bookmark #%d\n", bookmark.ID)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Title:      %s\n", bookmark.Title)
	fmt.Printf("URL:        %s\n", bookmark.URL)

	if bookmark.Author != "" {
		fmt.Printf("Author:     %s\n", bookmark.Author)
	}

	if bookmark.Excerpt != "" {
		fmt.Printf("Excerpt:    %s\n", bookmark.Excerpt)
	}

	tags := formatTags(bookmark.Tags)
	fmt.Printf("Tags:       %s\n", tags)

	fmt.Printf("Public:     %t\n", bookmark.Public)
	fmt.Printf("Created:    %s\n", bookmark.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Modified:   %s\n", bookmark.ModifiedAt.Format("2006-01-02 15:04:05"))

	if bookmark.HasContent && bookmark.Content != "" {
		fmt.Println("\nContent:")
		fmt.Println(strings.Repeat("-", 60))
		// Show first 500 characters of content
		content := bookmark.Content
		if len(content) > 500 {
			content = content[:500] + "..."
		}
		fmt.Println(content)
	}
}
