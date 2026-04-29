package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/models"
	"github.com/jvzantvoort/linkmgr/internal/validator"
	"github.com/spf13/cobra"
)

var (
	addTitle   string
	addTags    string
	addExcerpt string
	addAuthor  string
	addPublic  bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a new bookmark",
	Long: `Add a new bookmark to the database.
	
The URL is required. Other metadata like title, tags, excerpt, and author
can be specified using flags. If no title is provided, the URL will be used.

Examples:
  linkmgr add https://example.com
  linkmgr add https://example.com --title "Example Site"
  linkmgr add https://golang.org --title "Go Lang" --tags "programming,golang"
  linkmgr add https://example.com --excerpt "A great example" --public`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "bookmark title")
	addCmd.Flags().StringVar(&addTags, "tags", "", "comma-separated tags")
	addCmd.Flags().StringVarP(&addExcerpt, "excerpt", "e", "", "brief description")
	addCmd.Flags().StringVarP(&addAuthor, "author", "a", "", "author name")
	addCmd.Flags().BoolVarP(&addPublic, "public", "p", false, "mark as public")
}

func runAdd(cmd *cobra.Command, args []string) error {
	rawURL := args[0]

	// Validate URL
	validURL, err := validator.ValidateURL(rawURL)
	if err != nil {
		display.ErrorWithSuggestion(err.Error(), "Make sure the URL includes http:// or https://")
		return err
	}

	// Prepare bookmark
	bookmark := &models.Bookmark{
		URL:     validURL,
		Title:   addTitle,
		Excerpt: addExcerpt,
		Author:  addAuthor,
		Public:  addPublic,
	}

	// Use URL as title if not provided
	if bookmark.Title == "" {
		bookmark.Title = validURL
	}

	// Parse and normalize tags
	tags := validator.NormalizeTags(addTags)

	// Create bookmark
	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	err = repo.Create(ctx, bookmark, tags)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			display.ErrorWithSuggestion(
				"A bookmark with this URL already exists",
				fmt.Sprintf("Use 'linkmgr search \"%s\"' to find it or 'linkmgr update' to modify it", validURL),
			)
			return err
		}
		display.Error("Failed to create bookmark: %v", err)
		return err
	}

	// Success message
	display.Success("Bookmark #%d created", bookmark.ID)
	fmt.Printf("Title: %s\n", bookmark.Title)
	fmt.Printf("URL:   %s\n", bookmark.URL)
	if len(tags) > 0 {
		fmt.Printf("Tags:  %s\n", strings.Join(tags, ", "))
	}

	return nil
}
