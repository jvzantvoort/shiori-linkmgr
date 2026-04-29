package cmd

import (
	"context"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	listTag        string
	listLimit      int
	listOffset     int
	listPublicOnly bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bookmarks",
	Long: `List bookmarks with optional filtering.

Examples:
  linkmgr list
  linkmgr list --limit 20
  linkmgr list --tag golang
  linkmgr list --public-only
  linkmgr list --tag dev --limit 10`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listTag, "tag", "", "filter by tag")
	listCmd.Flags().IntVarP(&listLimit, "limit", "n", 0, "maximum number of results (default from config)")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "result offset for pagination")
	listCmd.Flags().BoolVar(&listPublicOnly, "public-only", false, "show only public bookmarks")
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Use config default if limit not specified
	limit := listLimit
	if limit == 0 {
		limit = cfg.Display.DefaultLimit
	}

	filters := &repository.ListFilters{
		Tag:        listTag,
		PublicOnly: listPublicOnly,
		Limit:      limit,
		Offset:     listOffset,
	}

	bookmarks, err := repo.List(ctx, filters)
	if err != nil {
		display.Error("Failed to list bookmarks: %v", err)
		return err
	}

	display.FormatBookmarkTable(bookmarks)

	// Show pagination info if results might be truncated
	if len(bookmarks) == limit {
		display.Info("Showing %d results. Use --offset to see more.", limit)
	}

	return nil
}
