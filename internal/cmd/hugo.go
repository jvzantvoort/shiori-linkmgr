package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	hugoOutputDir   string
	hugoSectionName string
	hugoContentType string
	hugoAuthorName  string
	hugoDraft       bool
	hugoPerTagFile  bool
	hugoSingleFile  bool
)

// hugoCmd represents the hugo command
var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Export bookmarks to Hugo-compatible markdown",
	Long: `Export bookmarks to Hugo CMS-compatible markdown files.

Bookmarks can be exported in different modes:
- Per tag: One markdown file per tag (default)
- Single file: All bookmarks in one file
- Individual: One markdown file per bookmark

The generated markdown includes Hugo front matter with metadata.

Examples:
  # Export by tag to content/bookmarks/
  linkmgr hugo --output content/bookmarks

  # Export all in single file
  linkmgr hugo --output content/bookmarks --single-file

  # Specify content type
  linkmgr hugo --output content/links --type links

  # Export as drafts
  linkmgr hugo --output content/bookmarks --draft`,
	RunE: runHugo,
}

func init() {
	rootCmd.AddCommand(hugoCmd)

	hugoCmd.Flags().StringVarP(&hugoOutputDir, "output", "o", "content/bookmarks", "output directory for markdown files")
	hugoCmd.Flags().StringVar(&hugoSectionName, "section", "bookmarks", "Hugo section name")
	hugoCmd.Flags().StringVar(&hugoContentType, "type", "bookmark", "Hugo content type")
	hugoCmd.Flags().StringVar(&hugoAuthorName, "author", "", "author name for front matter")
	hugoCmd.Flags().BoolVar(&hugoDraft, "draft", false, "mark content as draft")
	hugoCmd.Flags().BoolVar(&hugoPerTagFile, "per-tag", true, "create one file per tag (default)")
	hugoCmd.Flags().BoolVar(&hugoSingleFile, "single-file", false, "create single file with all bookmarks")
}

func runHugo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo := db.NewBookmarkRepository()
	tagRepo := db.NewTagRepository()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(hugoOutputDir, 0755); err != nil {
		display.Error("Failed to create output directory: %v", err)
		return err
	}

	if hugoSingleFile {
		return exportHugoSingleFile(ctx, repo, hugoOutputDir)
	}

	// Default: export by tag
	return exportHugoByTag(ctx, repo, tagRepo, hugoOutputDir)
}

func linkTag(tag string) string {
	rootpath := "/links"
	return fmt.Sprintf("[%s](%s/%s)", tag, rootpath, tag)
}

func exportHugoSingleFile(ctx context.Context, repo repository.BookmarkRepository, outputDir string) error {
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

	// Create filename
	filename := filepath.Join(outputDir, "_index.md")

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		display.Error("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Write front matter
	fmt.Fprintf(file, "---\n")
	fmt.Fprintf(file, "title: \"Bookmarks\"\n")
	fmt.Fprintf(file, "date: %s\n", time.Now().Format("2006-01-02T15:04:05-07:00"))
	if hugoAuthorName != "" {
		fmt.Fprintf(file, "author: %s\n", hugoAuthorName)
	}
	fmt.Fprintf(file, "type: %s\n", hugoContentType)
	if hugoDraft {
		fmt.Fprintf(file, "draft: true\n")
	}
	fmt.Fprintf(file, "description: \"Collection of %d bookmarks\"\n", len(bookmarks))
	fmt.Fprintf(file, "---\n\n")

	// Write content
	fmt.Fprintf(file, "# Bookmarks\n\n")
	fmt.Fprintf(file, "Total bookmarks: %d\n\n", len(bookmarks))

	// Group by first tag
	tagGroups := make(map[string][]int)
	for i, bookmark := range bookmarks {
		if len(bookmark.Tags) > 0 {
			tagName := bookmark.Tags[0].Name
			tagGroups[tagName] = append(tagGroups[tagName], i)
		} else {
			tagGroups["untagged"] = append(tagGroups["untagged"], i)
		}
	}

	// Write bookmarks by tag
	for tagName, indices := range tagGroups {
		fmt.Fprintf(file, "## %s\n\n", strings.Title(tagName))

		for _, i := range indices {
			bookmark := bookmarks[i]
			fmt.Fprintf(file, "### [%s](%s)\n\n", bookmark.Title, bookmark.URL)

			if bookmark.Excerpt != "" {
				fmt.Fprintf(file, "%s\n\n", bookmark.Excerpt)
			}

			if bookmark.Author != "" {
				fmt.Fprintf(file, "**Author:** %s\n\n", bookmark.Author)
			}

			// Write tags
			if len(bookmark.Tags) > 0 {
				fmt.Fprintf(file, "**Tags:** ")
				for i, tag := range bookmark.Tags {
					if i > 0 {
						fmt.Fprintf(file, ", ")
					}
					fmt.Fprintf(file, "%s", linkTag(tag.Name))
				}
				fmt.Fprintf(file, "\n\n")
			}

			fmt.Fprintf(file, "---\n\n")
		}
	}

	display.Success("Exported %d bookmarks to %s", len(bookmarks), filename)
	return nil
}

func exportHugoByTag(ctx context.Context, repo repository.BookmarkRepository, tagRepo repository.TagRepository, outputDir string) error {
	// Get all tags with counts
	tagCounts, err := tagRepo.GetAllWithCounts(ctx)
	if err != nil {
		display.Error("Failed to retrieve tags: %v", err)
		return err
	}

	if len(tagCounts) == 0 {
		display.Info("No tags found")
		return nil
	}

	totalFiles := 0
	totalBookmarks := 0

	// Export bookmarks for each tag
	for _, tc := range tagCounts {
		if tc.Count == 0 {
			continue
		}

		tagName := tc.Tag.Name

		// Fetch bookmarks for this tag
		filters := &repository.ListFilters{
			Tag:   tagName,
			Limit: 10000,
		}

		bookmarks, err := repo.List(ctx, filters)
		if err != nil {
			display.Error("Failed to retrieve bookmarks for tag '%s': %v", tagName, err)
			continue
		}

		if len(bookmarks) == 0 {
			continue
		}

		// Create safe filename from tag name
		safeTagName := strings.ToLower(tagName)
		safeTagName = strings.ReplaceAll(safeTagName, " ", "-")
		safeTagName = strings.ReplaceAll(safeTagName, "_", "-")

		filename := filepath.Join(outputDir, safeTagName+".md")

		// Create file
		file, err := os.Create(filename)
		if err != nil {
			display.Error("Failed to create file for tag '%s': %v", tagName, err)
			continue
		}

		// Write front matter
		fmt.Fprintf(file, "---\n")
		fmt.Fprintf(file, "title: \"%s Bookmarks\"\n", strings.Title(tagName))
		fmt.Fprintf(file, "date: %s\n", time.Now().Format("2006-01-02T15:04:05-07:00"))
		if hugoAuthorName != "" {
			fmt.Fprintf(file, "author: %s\n", hugoAuthorName)
		}
		fmt.Fprintf(file, "type: %s\n", hugoContentType)
		fmt.Fprintf(file, "tags:\n")
		fmt.Fprintf(file, "  - %s\n", tagName)
		if hugoDraft {
			fmt.Fprintf(file, "draft: true\n")
		}
		fmt.Fprintf(file, "description: \"Collection of %d bookmarks tagged with %s\"\n", len(bookmarks), tagName)
		fmt.Fprintf(file, "---\n\n")

		// Write content
		fmt.Fprintf(file, "# %s Bookmarks\n\n", strings.Title(tagName))
		fmt.Fprintf(file, "A curated collection of %d links related to %s.\n\n", len(bookmarks), tagName)

		// Write bookmarks
		for _, bookmark := range bookmarks {
			fmt.Fprintf(file, "## [%s](%s)\n\n", bookmark.Title, bookmark.URL)

			if bookmark.Excerpt != "" {
				fmt.Fprintf(file, "%s\n\n", bookmark.Excerpt)
			}

			if bookmark.Author != "" {
				fmt.Fprintf(file, "**Author:** %s\n\n", bookmark.Author)
			}

			// Write additional tags
			if len(bookmark.Tags) > 1 {
				fmt.Fprintf(file, "**Also tagged:** ")
				first := true
				for _, tag := range bookmark.Tags {
					if tag.Name != tagName {
						if !first {
							fmt.Fprintf(file, ", ")
						}
						// fmt.Fprintf(file, "`%s`", tag.Name)
						fmt.Fprintf(file, "%s", linkTag(tag.Name))
						first = false
					}
				}
				fmt.Fprintf(file, "\n\n")
			}

			if bookmark.Content != "" && len(bookmark.Content) > 0 {
				fmt.Fprintf(file, "### Summary\n\n")
				// Truncate content if too long
				content := bookmark.Content
				if len(content) > 500 {
					content = content[:500] + "..."
				}
				fmt.Fprintf(file, "%s\n\n", content)
			}

			fmt.Fprintf(file, "---\n\n")
		}

		file.Close()

		totalFiles++
		totalBookmarks += len(bookmarks)

		if verbose {
			fmt.Printf("Created %s (%d bookmarks)\n", filename, len(bookmarks))
		}
	}

	display.Success("Exported %d bookmarks to %d file(s) in %s", totalBookmarks, totalFiles, outputDir)
	return nil
}
