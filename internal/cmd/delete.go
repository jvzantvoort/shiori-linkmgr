package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id> [id...]",
	Short: "Delete bookmarks",
	Long: `Delete one or more bookmarks from the database.

By default, you will be prompted to confirm deletion. Use --force
to skip the confirmation prompt.

When deleting a bookmark, its tag associations are also removed.
Tags themselves are not deleted (use 'linkmgr tags --cleanup' for that).

Examples:
  # Delete single bookmark with confirmation
  linkmgr delete 5
  
  # Delete multiple bookmarks
  linkmgr delete 5 10 15
  
  # Delete without confirmation
  linkmgr delete 5 --force`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
	// Parse and validate bookmark IDs
	ids := make([]int, 0, len(args))
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			display.Error("Invalid bookmark ID: %s", arg)
			return err
		}
		ids = append(ids, id)
	}

	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Fetch bookmarks to show what will be deleted
	bookmarksToDelete := make(map[int]string) // id -> title
	notFound := []int{}

	for _, id := range ids {
		bookmark, err := repo.GetByID(ctx, id)
		if err != nil {
			if err.Error() == "bookmark not found" {
				notFound = append(notFound, id)
			} else {
				display.Error("Failed to retrieve bookmark #%d: %v", id, err)
				return err
			}
		} else {
			bookmarksToDelete[id] = bookmark.Title
		}
	}

	// Report not found bookmarks
	if len(notFound) > 0 {
		fmt.Printf("Warning: The following bookmark IDs were not found: %v\n", notFound)
	}

	// Check if any bookmarks to delete
	if len(bookmarksToDelete) == 0 {
		display.Error("No valid bookmarks to delete")
		return fmt.Errorf("no bookmarks found")
	}

	// Show what will be deleted
	fmt.Println("The following bookmarks will be deleted:")
	fmt.Println()
	for id, title := range bookmarksToDelete {
		fmt.Printf("  #%d: %s\n", id, title)
	}
	fmt.Println()

	// Confirm deletion unless --force is set
	if !deleteForce {
		confirmed, err := confirmDeletion(len(bookmarksToDelete))
		if err != nil {
			return err
		}
		if !confirmed {
			display.Info("Deletion cancelled")
			return nil
		}
	}

	// Delete bookmarks
	deleted := 0
	failed := 0

	for id := range bookmarksToDelete {
		if err := repo.Delete(ctx, id); err != nil {
			display.Error("Failed to delete bookmark #%d: %v", id, err)
			failed++
		} else {
			deleted++
		}
	}

	// Report results
	if deleted > 0 {
		display.Success("Deleted %d bookmark(s)", deleted)
	}
	if failed > 0 {
		display.Error("Failed to delete %d bookmark(s)", failed)
	}

	return nil
}

// confirmDeletion prompts the user to confirm deletion
func confirmDeletion(count int) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	if count == 1 {
		fmt.Print("Delete this bookmark? (yes/no): ")
	} else {
		fmt.Printf("Delete these %d bookmarks? (yes/no): ", count)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y", nil
}
