package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/models"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export <filename>",
	Short: "Export bookmarks",
	Long: `Export all bookmarks to a file.

Supports JSON and CSV formats. The format is determined by the --format flag
or automatically detected from the filename extension.

Examples:
  linkmgr export bookmarks.json
  linkmgr export bookmarks.csv
  linkmgr export backup.json --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "", "export format (json or csv, auto-detected from extension)")
}

func runExport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Determine format
	format := exportFormat
	if format == "" {
		// Auto-detect from extension
		if len(filename) > 5 && filename[len(filename)-5:] == ".json" {
			format = "json"
		} else if len(filename) > 4 && filename[len(filename)-4:] == ".csv" {
			format = "csv"
		} else {
			format = "json" // Default to JSON
		}
	}

	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Fetch all bookmarks
	filters := &repository.ListFilters{Limit: 100000}
	bookmarks, err := repo.List(ctx, filters)
	if err != nil {
		display.Error("Failed to retrieve bookmarks: %v", err)
		return err
	}

	if len(bookmarks) == 0 {
		display.Info("No bookmarks to export")
		return nil
	}

	// Export based on format
	switch format {
	case "json":
		err = exportJSON(filename, bookmarks)
	case "csv":
		err = exportCSV(filename, bookmarks)
	default:
		return fmt.Errorf("unsupported format: %s (use json or csv)", format)
	}

	if err != nil {
		display.Error("Failed to export: %v", err)
		return err
	}

	// Get file size
	stat, _ := os.Stat(filename)
	size := stat.Size()

	display.Success("Exported %d bookmark(s) to %s (%d bytes)", len(bookmarks), filename, size)
	return nil
}

func exportJSON(filename string, bookmarks []models.Bookmark) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(bookmarks)
}

func exportCSV(filename string, bookmarks []models.Bookmark) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "URL", "Title", "Excerpt", "Author", "Public", "Tags", "Created", "Modified"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write bookmarks
	for _, b := range bookmarks {
		tags := ""
		for i, tag := range b.Tags {
			if i > 0 {
				tags += ","
			}
			tags += tag.Name
		}

		public := "false"
		if b.Public {
			public = "true"
		}

		row := []string{
			fmt.Sprintf("%d", b.ID),
			b.URL,
			b.Title,
			b.Excerpt,
			b.Author,
			public,
			tags,
			b.CreatedAt.Format("2006-01-02 15:04:05"),
			b.ModifiedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
