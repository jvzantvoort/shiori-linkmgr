package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/spf13/cobra"
)

var openURL bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show bookmark details",
	Long: `Display detailed information about a specific bookmark.

Shows all metadata including title, URL, author, excerpt, content,
tags, and timestamps.

Examples:
  linkmgr show 5
  linkmgr show 10 --open`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&openURL, "open", "o", false, "open URL in default browser")
}

func runShow(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		display.Error("Invalid bookmark ID: %s", args[0])
		return err
	}

	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	bookmark, err := repo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "bookmark not found" {
			display.ErrorWithSuggestion(
				fmt.Sprintf("Bookmark #%d not found", id),
				"Use 'linkmgr list' to see available bookmarks",
			)
			return err
		}
		display.Error("Failed to retrieve bookmark: %v", err)
		return err
	}

	// Display bookmark details
	display.FormatBookmarkDetail(bookmark)

	// Open URL if requested
	if openURL {
		fmt.Println()
		if err := openBrowser(bookmark.URL); err != nil {
			display.Error("Failed to open URL: %v", err)
			return err
		}
		display.Success("Opened %s in browser", bookmark.URL)
	}

	return nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}
