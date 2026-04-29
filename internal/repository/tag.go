package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/models"
)

type tagRepo struct {
	db *sql.DB
}

// GetAll retrieves all tags
func (r *tagRepo) GetAll(ctx context.Context) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM tag ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetAllWithCounts retrieves all tags with usage counts
func (r *tagRepo) GetAllWithCounts(ctx context.Context) ([]TagCount, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name, COUNT(bt.bookmark_id) as count
		FROM tag t
		LEFT JOIN bookmark_tag bt ON t.id = bt.tag_id
		GROUP BY t.id, t.name
		ORDER BY count DESC, t.name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags with counts: %w", err)
	}
	defer rows.Close()

	tagCounts := []TagCount{}
	for rows.Next() {
		var tc TagCount
		if err := rows.Scan(&tc.Tag.ID, &tc.Tag.Name, &tc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan tag count: %w", err)
		}
		tagCounts = append(tagCounts, tc)
	}

	return tagCounts, nil
}

// GetByName retrieves a tag by name
func (r *tagRepo) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	tag := &models.Tag{}
	name = strings.TrimSpace(strings.ToLower(name))

	err := r.db.QueryRowContext(ctx, "SELECT id, name FROM tag WHERE name = ?", name).
		Scan(&tag.ID, &tag.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tag: %w", err)
	}

	return tag, nil
}

// Create creates a new tag
func (r *tagRepo) Create(ctx context.Context, name string) (*models.Tag, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("tag name cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, "INSERT INTO tag (name) VALUES (?)", name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tag: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get tag ID: %w", err)
	}

	return &models.Tag{ID: int(id), Name: name}, nil
}

// GetOrCreate gets an existing tag or creates a new one
func (r *tagRepo) GetOrCreate(ctx context.Context, name string) (*models.Tag, error) {
	tag, err := r.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if tag != nil {
		return tag, nil
	}

	return r.Create(ctx, name)
}

// Rename renames a tag
func (r *tagRepo) Rename(ctx context.Context, oldName, newName string) error {
	oldName = strings.TrimSpace(strings.ToLower(oldName))
	newName = strings.TrimSpace(strings.ToLower(newName))

	if oldName == "" || newName == "" {
		return fmt.Errorf("tag names cannot be empty")
	}

	// Check if new name already exists
	existing, err := r.GetByName(ctx, newName)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("tag %q already exists", newName)
	}

	result, err := r.db.ExecContext(ctx, "UPDATE tag SET name = ? WHERE name = ?", newName, oldName)
	if err != nil {
		return fmt.Errorf("failed to rename tag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tag %q not found", oldName)
	}

	return nil
}

// DeleteOrphaned deletes tags with no bookmark associations
func (r *tagRepo) DeleteOrphaned(ctx context.Context) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM tag
		WHERE id NOT IN (SELECT DISTINCT tag_id FROM bookmark_tag)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to delete orphaned tags: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return count, nil
}

// GetBookmarkTags retrieves all tags for a bookmark
func (r *tagRepo) GetBookmarkTags(ctx context.Context, bookmarkID int) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name
		FROM tag t
		INNER JOIN bookmark_tag bt ON t.id = bt.tag_id
		WHERE bt.bookmark_id = ?
		ORDER BY t.name
	`, bookmarkID)

	if err != nil {
		return nil, fmt.Errorf("failed to query bookmark tags: %w", err)
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
