package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/validator"
	"github.com/spf13/cobra"
)

var (
	updateURL        string
	updateTitle      string
	updateExcerpt    string
	updateAuthor     string
	updateContent    string
	updateTags       string
	updateAddTags    string
	updateRemoveTags string
	updatePublic     bool
	updateNoPublic   bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update bookmark",
	Long: `Update an existing bookmark's metadata.

You can update the URL, title, excerpt, author, content, and public status.
For tags, you can replace all tags, add new ones, or remove specific ones.

Examples:
  # Update URL
  linkmgr update 5 --url "https://new-url.com"
  
  # Update title
  linkmgr update 5 --title "New Title"
  
  # Update multiple fields
  linkmgr update 5 --title "New Title" --excerpt "Updated description"
  
  # Replace all tags
  linkmgr update 5 --tags "newtag1,newtag2"
  
  # Add tags without removing existing ones
  linkmgr update 5 --add-tags "extra,another"
  
  # Remove specific tags
  linkmgr update 5 --remove-tags "oldtag"
  
  # Toggle public status
  linkmgr update 5 --public
  linkmgr update 5 --no-public`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateURL, "url", "u", "", "new URL")
	updateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "new title")
	updateCmd.Flags().StringVarP(&updateExcerpt, "excerpt", "e", "", "new excerpt")
	updateCmd.Flags().StringVarP(&updateAuthor, "author", "a", "", "new author")
	updateCmd.Flags().StringVar(&updateContent, "content", "", "new content")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "replace all tags (comma-separated)")
	updateCmd.Flags().StringVar(&updateAddTags, "add-tags", "", "add tags (comma-separated)")
	updateCmd.Flags().StringVar(&updateRemoveTags, "remove-tags", "", "remove tags (comma-separated)")
	updateCmd.Flags().BoolVarP(&updatePublic, "public", "p", false, "mark as public")
	updateCmd.Flags().BoolVar(&updateNoPublic, "no-public", false, "mark as private")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		display.Error("Invalid bookmark ID: %s", args[0])
		return err
	}

	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Fetch existing bookmark
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

	// Track what was updated
	updated := []string{}

	// Update fields if provided
	if cmd.Flags().Changed("url") {
		// Validate new URL
		validURL, err := validator.ValidateURL(updateURL)
		if err != nil {
			display.ErrorWithSuggestion(err.Error(), "Make sure the URL includes http:// or https://")
			return err
		}
		bookmark.URL = validURL
		updated = append(updated, "url")
	}
	if cmd.Flags().Changed("title") {
		bookmark.Title = updateTitle
		updated = append(updated, "title")
	}
	if cmd.Flags().Changed("excerpt") {
		bookmark.Excerpt = updateExcerpt
		updated = append(updated, "excerpt")
	}
	if cmd.Flags().Changed("author") {
		bookmark.Author = updateAuthor
		updated = append(updated, "author")
	}
	if cmd.Flags().Changed("content") {
		bookmark.Content = updateContent
		bookmark.HasContent = updateContent != ""
		updated = append(updated, "content")
	}
	if cmd.Flags().Changed("public") {
		bookmark.Public = true
		updated = append(updated, "public")
	}
	if cmd.Flags().Changed("no-public") {
		bookmark.Public = false
		updated = append(updated, "public")
	}

	// Update bookmark if any field changed
	if len(updated) > 0 {
		if err := repo.Update(ctx, bookmark); err != nil {
			display.Error("Failed to update bookmark: %v", err)
			return err
		}
	}

	// Handle tag updates
	tagsUpdated := false

	if cmd.Flags().Changed("tags") {
		// Replace all tags
		tags := validator.NormalizeTags(updateTags)
		if err := repo.UpdateTags(ctx, id, tags); err != nil {
			display.Error("Failed to update tags: %v", err)
			return err
		}
		updated = append(updated, "tags")
		tagsUpdated = true
	}

	if cmd.Flags().Changed("add-tags") {
		// Add tags
		tags := validator.NormalizeTags(updateAddTags)
		if err := repo.AddTags(ctx, id, tags); err != nil {
			display.Error("Failed to add tags: %v", err)
			return err
		}
		updated = append(updated, "tags")
		tagsUpdated = true
	}

	if cmd.Flags().Changed("remove-tags") {
		// Remove tags
		tags := validator.NormalizeTags(updateRemoveTags)
		if err := repo.RemoveTags(ctx, id, tags); err != nil {
			display.Error("Failed to remove tags: %v", err)
			return err
		}
		updated = append(updated, "tags")
		tagsUpdated = true
	}

	// Check if anything was actually updated
	if len(updated) == 0 {
		display.Info("No changes specified. Use --help to see available options.")
		return nil
	}

	// Fetch updated bookmark to show changes
	bookmark, err = repo.GetByID(ctx, id)
	if err != nil {
		display.Error("Failed to retrieve updated bookmark: %v", err)
		return err
	}

	// Display success message
	display.Success("Bookmark #%d updated", id)

	// Show what was updated
	fmt.Printf("Title: %s\n", bookmark.Title)
	fmt.Printf("URL:   %s\n", bookmark.URL)

	if tagsUpdated {
		tagNames := make([]string, len(bookmark.Tags))
		for i, tag := range bookmark.Tags {
			tagNames[i] = tag.Name
		}
		if len(tagNames) > 0 {
			fmt.Printf("Tags:  %s\n", formatTagList(tagNames))
		} else {
			fmt.Printf("Tags:  (none)\n")
		}
	}

	return nil
}

func formatTagList(tags []string) string {
	result := ""
	for i, tag := range tags {
		if i > 0 {
			result += ", "
		}
		result += tag
	}
	return result
}
