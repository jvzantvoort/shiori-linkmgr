package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jvzantvoort/linkmgr/internal/models"
)

type bookmarkRepo struct {
	db *sql.DB
}

// Create inserts a new bookmark with tags
func (r *bookmarkRepo) Create(ctx context.Context, bookmark *models.Bookmark, tags []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert bookmark
	result, err := tx.ExecContext(ctx, `
		INSERT INTO bookmark (url, title, excerpt, author, public, content, html, has_content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, bookmark.URL, bookmark.Title, bookmark.Excerpt, bookmark.Author, bookmark.Public,
		bookmark.Content, bookmark.HTML, bookmark.HasContent)

	if err != nil {
		return fmt.Errorf("failed to insert bookmark: %w", err)
	}

	bookmarkID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get bookmark ID: %w", err)
	}
	bookmark.ID = int(bookmarkID)

	// Insert tags
	if len(tags) > 0 {
		tagRepo := &tagRepo{db: r.db}
		for _, tagName := range tags {
			tagName = strings.TrimSpace(strings.ToLower(tagName))
			if tagName == "" {
				continue
			}

			tag, err := tagRepo.GetOrCreate(ctx, tagName)
			if err != nil {
				return fmt.Errorf("failed to get/create tag %q: %w", tagName, err)
			}

			_, err = tx.ExecContext(ctx, `
				INSERT INTO bookmark_tag (bookmark_id, tag_id) VALUES (?, ?)
			`, bookmarkID, tag.ID)
			if err != nil {
				return fmt.Errorf("failed to associate tag %q: %w", tagName, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a bookmark by ID with its tags
func (r *bookmarkRepo) GetByID(ctx context.Context, id int) (*models.Bookmark, error) {
	bookmark := &models.Bookmark{}

	err := r.db.QueryRowContext(ctx, `
		SELECT id, url, title, excerpt, author, public, content, html, 
		       created_at, modified_at, has_content
		FROM bookmark
		WHERE id = ?
	`, id).Scan(
		&bookmark.ID, &bookmark.URL, &bookmark.Title, &bookmark.Excerpt,
		&bookmark.Author, &bookmark.Public, &bookmark.Content, &bookmark.HTML,
		&bookmark.CreatedAt, &bookmark.ModifiedAt, &bookmark.HasContent,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bookmark not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query bookmark: %w", err)
	}

	// Load tags
	tagRepo := &tagRepo{db: r.db}
	tags, err := tagRepo.GetBookmarkTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	bookmark.Tags = tags

	return bookmark, nil
}

// List retrieves bookmarks with optional filters
func (r *bookmarkRepo) List(ctx context.Context, filters *ListFilters) ([]models.Bookmark, error) {
	if filters == nil {
		filters = &ListFilters{Limit: 50}
	}
	if filters.Limit <= 0 {
		filters.Limit = 50
	}

	query := `
		SELECT DISTINCT b.id, b.url, b.title, b.excerpt, b.author, b.public,
		       b.created_at, b.modified_at
		FROM bookmark b
	`
	args := []interface{}{}

	if filters.Tag != "" {
		query += `
			INNER JOIN bookmark_tag bt ON b.id = bt.bookmark_id
			INNER JOIN tag t ON bt.tag_id = t.id
		`
	}

	where := []string{}
	if filters.Tag != "" {
		where = append(where, "t.name = ?")
		args = append(args, strings.ToLower(filters.Tag))
	}
	if filters.PublicOnly {
		where = append(where, "b.public = 1")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query += " ORDER BY b.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookmarks: %w", err)
	}
	defer rows.Close()

	bookmarks := []models.Bookmark{}
	tagRepo := &tagRepo{db: r.db}

	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(
			&b.ID, &b.URL, &b.Title, &b.Excerpt, &b.Author, &b.Public,
			&b.CreatedAt, &b.ModifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		// Load tags for each bookmark
		tags, err := tagRepo.GetBookmarkTags(ctx, b.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags for bookmark %d: %w", b.ID, err)
		}
		b.Tags = tags

		bookmarks = append(bookmarks, b)
	}

	return bookmarks, nil
}

// Search performs fulltext search on bookmarks
func (r *bookmarkRepo) Search(ctx context.Context, query string, limit int) ([]models.Bookmark, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url, title, excerpt, author, public, created_at, modified_at
		FROM bookmark
		WHERE MATCH(title, excerpt, content) AGAINST(? IN NATURAL LANGUAGE MODE)
		ORDER BY created_at DESC
		LIMIT ?
	`, query, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to search bookmarks: %w", err)
	}
	defer rows.Close()

	bookmarks := []models.Bookmark{}
	tagRepo := &tagRepo{db: r.db}

	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(
			&b.ID, &b.URL, &b.Title, &b.Excerpt, &b.Author, &b.Public,
			&b.CreatedAt, &b.ModifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		tags, err := tagRepo.GetBookmarkTags(ctx, b.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags: %w", err)
		}
		b.Tags = tags

		bookmarks = append(bookmarks, b)
	}

	return bookmarks, nil
}

// Update updates a bookmark
func (r *bookmarkRepo) Update(ctx context.Context, bookmark *models.Bookmark) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE bookmark
		SET url = ?, title = ?, excerpt = ?, author = ?, public = ?,
		    content = ?, html = ?, has_content = ?, modified_at = ?
		WHERE id = ?
	`, bookmark.URL, bookmark.Title, bookmark.Excerpt, bookmark.Author, bookmark.Public,
		bookmark.Content, bookmark.HTML, bookmark.HasContent, time.Now(), bookmark.ID)

	if err != nil {
		return fmt.Errorf("failed to update bookmark: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("bookmark not found")
	}

	return nil
}

// UpdateTags replaces all tags for a bookmark
func (r *bookmarkRepo) UpdateTags(ctx context.Context, bookmarkID int, tags []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing tags
	_, err = tx.ExecContext(ctx, "DELETE FROM bookmark_tag WHERE bookmark_id = ?", bookmarkID)
	if err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	// Add new tags
	tagRepo := &tagRepo{db: r.db}
	for _, tagName := range tags {
		tagName = strings.TrimSpace(strings.ToLower(tagName))
		if tagName == "" {
			continue
		}

		tag, err := tagRepo.GetOrCreate(ctx, tagName)
		if err != nil {
			return fmt.Errorf("failed to get/create tag: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO bookmark_tag (bookmark_id, tag_id) VALUES (?, ?)
		`, bookmarkID, tag.ID)
		if err != nil {
			return fmt.Errorf("failed to associate tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddTags adds tags to a bookmark
func (r *bookmarkRepo) AddTags(ctx context.Context, bookmarkID int, tags []string) error {
	tagRepo := &tagRepo{db: r.db}

	for _, tagName := range tags {
		tagName = strings.TrimSpace(strings.ToLower(tagName))
		if tagName == "" {
			continue
		}

		tag, err := tagRepo.GetOrCreate(ctx, tagName)
		if err != nil {
			return fmt.Errorf("failed to get/create tag: %w", err)
		}

		_, err = r.db.ExecContext(ctx, `
			INSERT IGNORE INTO bookmark_tag (bookmark_id, tag_id) VALUES (?, ?)
		`, bookmarkID, tag.ID)
		if err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}
	}

	return nil
}

// RemoveTags removes specific tags from a bookmark
func (r *bookmarkRepo) RemoveTags(ctx context.Context, bookmarkID int, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	placeholders := make([]string, len(tags))
	args := []interface{}{bookmarkID}

	for i, tagName := range tags {
		placeholders[i] = "?"
		args = append(args, strings.TrimSpace(strings.ToLower(tagName)))
	}

	query := fmt.Sprintf(`
		DELETE bt FROM bookmark_tag bt
		INNER JOIN tag t ON bt.tag_id = t.id
		WHERE bt.bookmark_id = ? AND t.name IN (%s)
	`, strings.Join(placeholders, ","))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove tags: %w", err)
	}

	return nil
}

// Delete deletes a bookmark and its tag associations
func (r *bookmarkRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete tag associations
	_, err = tx.ExecContext(ctx, "DELETE FROM bookmark_tag WHERE bookmark_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete tag associations: %w", err)
	}

	// Delete bookmark
	result, err := tx.ExecContext(ctx, "DELETE FROM bookmark WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("bookmark not found")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
