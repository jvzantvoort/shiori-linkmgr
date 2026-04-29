package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	validateConcurrency int
	validateTimeout     int
	validateTag         string
	validateUpdateDB    bool
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate bookmark URLs",
	Long: `Check if bookmark URLs are still reachable.

This command sends HTTP HEAD requests to all bookmark URLs to verify
they are still accessible. It reports broken links and optionally
updates bookmark metadata.

Examples:
  linkmgr validate
  linkmgr validate --tag important
  linkmgr validate --concurrency 10
  linkmgr validate --timeout 10`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().IntVarP(&validateConcurrency, "concurrency", "c", 5, "number of concurrent checks")
	validateCmd.Flags().IntVarP(&validateTimeout, "timeout", "t", 5, "timeout in seconds per URL")
	validateCmd.Flags().StringVar(&validateTag, "tag", "", "only validate bookmarks with this tag")
	validateCmd.Flags().BoolVar(&validateUpdateDB, "update", false, "update bookmark metadata with validation results (future)")
}

type ValidationResult struct {
	BookmarkID int
	URL        string
	Title      string
	StatusCode int
	Error      error
	Duration   time.Duration
}

func runValidate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo := db.NewBookmarkRepository()

	// Fetch bookmarks to validate
	filters := &repository.ListFilters{
		Tag:   validateTag,
		Limit: 10000, // Get all bookmarks
	}

	bookmarks, err := repo.List(ctx, filters)
	if err != nil {
		display.Error("Failed to retrieve bookmarks: %v", err)
		return err
	}

	if len(bookmarks) == 0 {
		display.Info("No bookmarks to validate")
		return nil
	}

	display.Info("Validating %d bookmark(s)...", len(bookmarks))
	fmt.Println()

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(validateTimeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow redirects, up to 10
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Validate URLs concurrently
	results := make(chan ValidationResult, len(bookmarks))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, validateConcurrency)

	for _, bookmark := range bookmarks {
		wg.Add(1)
		go func(id int, url, title string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := ValidationResult{
				BookmarkID: id,
				URL:        url,
				Title:      title,
			}

			start := time.Now()
			resp, err := client.Head(url)
			result.Duration = time.Since(start)

			if err != nil {
				result.Error = err
			} else {
				result.StatusCode = resp.StatusCode
				resp.Body.Close()
			}

			results <- result
		}(bookmark.ID, bookmark.URL, bookmark.Title)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display results
	var successful, failed, warnings int
	var failedResults []ValidationResult

	for result := range results {
		if result.Error != nil {
			failed++
			failedResults = append(failedResults, result)
			fmt.Printf("✗ [#%d] %s - ERROR: %v\n", result.BookmarkID, truncateURL(result.URL, 60), result.Error)
		} else if result.StatusCode >= 400 {
			failed++
			failedResults = append(failedResults, result)
			fmt.Printf("✗ [#%d] %s - HTTP %d\n", result.BookmarkID, truncateURL(result.URL, 60), result.StatusCode)
		} else if result.StatusCode >= 300 {
			warnings++
			fmt.Printf("⚠ [#%d] %s - HTTP %d (%.2fs)\n", result.BookmarkID, truncateURL(result.URL, 60), result.StatusCode, result.Duration.Seconds())
		} else {
			successful++
			if verbose {
				fmt.Printf("✓ [#%d] %s - HTTP %d (%.2fs)\n", result.BookmarkID, truncateURL(result.URL, 60), result.StatusCode, result.Duration.Seconds())
			}
		}
	}

	// Display summary
	fmt.Println()
	fmt.Println("Validation Summary:")
	fmt.Println("==================")
	fmt.Printf("Total:      %d\n", len(bookmarks))
	fmt.Printf("Successful: %d (HTTP 2xx)\n", successful)
	if warnings > 0 {
		fmt.Printf("Warnings:   %d (HTTP 3xx)\n", warnings)
	}
	if failed > 0 {
		fmt.Printf("Failed:     %d (HTTP 4xx/5xx or errors)\n", failed)
	}

	// Show detailed failed results
	if len(failedResults) > 0 {
		fmt.Println()
		fmt.Println("Failed Bookmarks:")
		fmt.Println("=================")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()

		fmt.Fprintf(w, "ID\tTitle\tURL\tStatus\n")
		fmt.Fprintf(w, "--\t-----\t---\t------\n")

		for _, result := range failedResults {
			status := fmt.Sprintf("HTTP %d", result.StatusCode)
			if result.Error != nil {
				status = result.Error.Error()
				if len(status) > 40 {
					status = status[:37] + "..."
				}
			}

			title := result.Title
			if len(title) > 30 {
				title = title[:27] + "..."
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				result.BookmarkID,
				title,
				truncateURL(result.URL, 40),
				status,
			)
		}
	}

	if failed > 0 {
		fmt.Println()
		display.Info("Use 'linkmgr show <id>' to view details or 'linkmgr delete <id>' to remove broken bookmarks")
	}

	return nil
}

func truncateURL(url string, maxLen int) string {
	if len(url) <= maxLen {
		return url
	}
	return url[:maxLen-3] + "..."
}
