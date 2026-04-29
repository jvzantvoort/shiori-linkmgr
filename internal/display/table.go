package display

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jvzantvoort/linkmgr/internal/models"
)

// FormatBookmarkTable formats bookmarks as a table
func FormatBookmarkTable(bookmarks []models.Bookmark) {
	if len(bookmarks) == 0 {
		fmt.Println("No bookmarks found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Header
	fmt.Fprintf(w, "ID\tTitle\tURL\tTags\tCreated\n")
	fmt.Fprintf(w, "--\t-----\t---\t----\t-------\n")

	// Rows
	for _, b := range bookmarks {
		title := truncate(b.Title, 40)
		url := truncate(b.URL, 50)
		tags := formatTags(b.Tags)
		created := b.CreatedAt.Format("2006-01-02")

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", b.ID, title, url, tags, created)
	}
}

// formatTags formats tags as comma-separated string
func formatTags(tags []models.Tag) string {
	if len(tags) == 0 {
		return "-"
	}

	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return strings.Join(names, ",")
}

// truncate truncates a string to maxLen with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
