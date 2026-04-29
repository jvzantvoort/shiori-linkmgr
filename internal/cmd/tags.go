package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	tagCleanup bool
	tagDryRun  bool
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	Long: `List all tags with their usage counts.

Shows how many bookmarks are associated with each tag, sorted by
usage count (most used first).

Examples:
  linkmgr tags
  linkmgr tags --cleanup          # Remove unused tags
  linkmgr tags --cleanup --dry-run`,
	RunE: runTags,
}

func init() {
	rootCmd.AddCommand(tagsCmd)
	tagsCmd.Flags().BoolVar(&tagCleanup, "cleanup", false, "remove orphaned tags (tags with no bookmarks)")
	tagsCmd.Flags().BoolVar(&tagDryRun, "dry-run", false, "show what would be deleted without actually deleting")
}

func runTags(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo := db.NewTagRepository()

	// If cleanup requested
	if tagCleanup {
		return runTagCleanup(ctx, repo)
	}

	// List all tags with counts
	tagCounts, err := repo.GetAllWithCounts(ctx)
	if err != nil {
		display.Error("Failed to retrieve tags: %v", err)
		return err
	}

	if len(tagCounts) == 0 {
		display.Info("No tags found.")
		return nil
	}

	// Display tags in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Tag\tCount\n")
	fmt.Fprintf(w, "---\t-----\n")

	for _, tc := range tagCounts {
		fmt.Fprintf(w, "%s\t%d\n", tc.Tag.Name, tc.Count)
	}

	fmt.Println()
	display.Info("Total tags: %d", len(tagCounts))

	return nil
}

func runTagCleanup(ctx context.Context, repo repository.TagRepository) error {
	// Get all tags with counts
	tagCounts, err := repo.GetAllWithCounts(ctx)
	if err != nil {
		display.Error("Failed to retrieve tags: %v", err)
		return err
	}

	// Find orphaned tags
	orphaned := []string{}
	for _, tc := range tagCounts {
		if tc.Count == 0 {
			orphaned = append(orphaned, tc.Tag.Name)
		}
	}

	if len(orphaned) == 0 {
		display.Success("No orphaned tags found. Database is clean!")
		return nil
	}

	// Show what will be deleted
	fmt.Printf("Found %d orphaned tag(s):\n", len(orphaned))
	for _, tag := range orphaned {
		fmt.Printf("  - %s\n", tag)
	}
	fmt.Println()

	if tagDryRun {
		display.Info("Dry run - no tags were deleted")
		return nil
	}

	// Delete orphaned tags
	deleted, err := repo.DeleteOrphaned(ctx)
	if err != nil {
		display.Error("Failed to delete orphaned tags: %v", err)
		return err
	}

	display.Success("Deleted %d orphaned tag(s)", deleted)
	return nil
}
