package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/models"
	"github.com/jvzantvoort/linkmgr/internal/validator"
	"github.com/spf13/cobra"
)

var (
	importSkipDuplicates   bool
	importUpdateDuplicates bool
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <filename>",
	Short: "Import bookmarks",
	Long: `Import bookmarks from a file.

Supports JSON and CSV formats. The format is automatically detected from
the filename extension.

By default, duplicate URLs will cause an error. Use --skip-duplicates to
ignore them, or --update-duplicates to update existing bookmarks.

Examples:
  linkmgr import bookmarks.json
  linkmgr import bookmarks.csv --skip-duplicates
  linkmgr import backup.json --update-duplicates`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&importSkipDuplicates, "skip-duplicates", false, "skip bookmarks with duplicate URLs")
	importCmd.Flags().BoolVar(&importUpdateDuplicates, "update-duplicates", false, "update existing bookmarks with duplicate URLs")
}

func runImport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Determine format from extension
	var bookmarks []models.Bookmark
	var err error

	if strings.HasSuffix(filename, ".json") {
		bookmarks, err = importJSON(filename)
	} else if strings.HasSuffix(filename, ".csv") {
		bookmarks, err = importCSV(filename)
	} else {
		return fmt.Errorf("unsupported file format (use .json or .csv)")
	}

	if err != nil {
		display.Error("Failed to read file: %v", err)
		return err
	}

	if len(bookmarks) == 0 {
		display.Info("No bookmarks to import")
		return nil
	}

	display.Info("Found %d bookmark(s) to import", len(bookmarks))

	// Import bookmarks
	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	imported := 0
	skipped := 0
	updated := 0
	failed := 0

	for _, bookmark := range bookmarks {
		// Validate URL
		validURL, err := validator.ValidateURL(bookmark.URL)
		if err != nil {
			fmt.Printf("Skipping invalid URL: %s\n", bookmark.URL)
			failed++
			continue
		}
		bookmark.URL = validURL

		// Extract tag names
		tagNames := make([]string, len(bookmark.Tags))
		for i, tag := range bookmark.Tags {
			tagNames[i] = tag.Name
		}

		// Try to create bookmark
		err = repo.Create(ctx, &bookmark, tagNames)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				// Handle duplicate
				if importSkipDuplicates {
					skipped++
					continue
				} else if importUpdateDuplicates {
					// TODO: Find existing bookmark and update it
					// For now, just skip
					skipped++
					continue
				} else {
					fmt.Printf("Duplicate URL: %s\n", bookmark.URL)
					failed++
					continue
				}
			} else {
				fmt.Printf("Failed to import %s: %v\n", bookmark.URL, err)
				failed++
				continue
			}
		}

		imported++
	}

	// Display summary
	fmt.Println()
	if imported > 0 {
		display.Success("Imported %d bookmark(s)", imported)
	}
	if skipped > 0 {
		display.Info("Skipped %d duplicate(s)", skipped)
	}
	if updated > 0 {
		display.Success("Updated %d bookmark(s)", updated)
	}
	if failed > 0 {
		display.Error("Failed to import %d bookmark(s)", failed)
	}

	return nil
}

func importJSON(filename string) ([]models.Bookmark, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bookmarks []models.Bookmark
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&bookmarks)
	return bookmarks, err
}

func importCSV(filename string) ([]models.Bookmark, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Find column indices
	indices := make(map[string]int)
	for i, col := range header {
		indices[col] = i
	}

	// Read records
	bookmarks := []models.Bookmark{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		bookmark := models.Bookmark{}

		// Parse fields
		if idx, ok := indices["URL"]; ok {
			bookmark.URL = record[idx]
		}
		if idx, ok := indices["Title"]; ok {
			bookmark.Title = record[idx]
		}
		if idx, ok := indices["Excerpt"]; ok {
			bookmark.Excerpt = record[idx]
		}
		if idx, ok := indices["Author"]; ok {
			bookmark.Author = record[idx]
		}
		if idx, ok := indices["Public"]; ok {
			bookmark.Public = record[idx] == "true"
		}

		// Parse tags
		if idx, ok := indices["Tags"]; ok && record[idx] != "" {
			tagNames := strings.Split(record[idx], ",")
			for _, tagName := range tagNames {
				tagName = strings.TrimSpace(tagName)
				if tagName != "" {
					bookmark.Tags = append(bookmark.Tags, models.Tag{Name: tagName})
				}
			}
		}

		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}
