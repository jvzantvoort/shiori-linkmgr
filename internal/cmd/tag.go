package cmd

import (
	"context"
	"fmt"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Long: `Manage individual tags.

Subcommands allow you to rename tags globally across all bookmarks.

Examples:
  linkmgr tag rename "old-name" "new-name"`,
}

// tagRenameCmd represents the tag rename command
var tagRenameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a tag",
	Long: `Rename a tag globally across all bookmarks.

This will update the tag name for all bookmarks that use it.
The old tag name will be replaced with the new tag name.

Examples:
  linkmgr tag rename "golang" "go"
  linkmgr tag rename "k8s" "kubernetes"`,
	Args: cobra.ExactArgs(2),
	RunE: runTagRename,
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagRenameCmd)
}

func runTagRename(cmd *cobra.Command, args []string) error {
	oldName := args[0]
	newName := args[1]

	ctx := context.Background()
	repo := db.NewTagRepository()

	// Check if old tag exists
	oldTag, err := repo.GetByName(ctx, oldName)
	if err != nil {
		display.Error("Failed to check tag: %v", err)
		return err
	}
	if oldTag == nil {
		display.ErrorWithSuggestion(
			fmt.Sprintf("Tag '%s' not found", oldName),
			"Use 'linkmgr tags' to see available tags",
		)
		return fmt.Errorf("tag not found")
	}

	// Rename the tag
	err = repo.Rename(ctx, oldName, newName)
	if err != nil {
		display.Error("Failed to rename tag: %v", err)
		return err
	}

	display.Success("Renamed tag '%s' to '%s'", oldName, newName)
	return nil
}
