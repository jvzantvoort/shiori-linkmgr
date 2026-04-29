package cmd

import (
	"context"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/spf13/cobra"
)

var searchLimit int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search bookmarks",
	Long: `Search bookmarks using fulltext search across title, excerpt, and content.

The search uses MySQL's FULLTEXT search capabilities, so you can use
natural language queries or boolean mode operators.

Examples:
  linkmgr search "golang tutorial"
  linkmgr search kubernetes
  linkmgr search "web development" --limit 20`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 0, "maximum number of results (default from config)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Use config default if limit not specified
	limit := searchLimit
	if limit == 0 {
		limit = cfg.Display.DefaultLimit
	}

	bookmarks, err := repo.Search(ctx, query, limit)
	if err != nil {
		display.Error("Failed to search bookmarks: %v", err)
		return err
	}

	if len(bookmarks) == 0 {
		display.Info("No bookmarks found matching '%s'", query)
		return nil
	}

	display.FormatBookmarkTable(bookmarks)

	return nil
}
