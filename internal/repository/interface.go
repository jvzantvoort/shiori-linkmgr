package repository

import (
	"context"

	"github.com/jvzantvoort/linkmgr/internal/models"
)

// BookmarkRepository defines operations for bookmark management
type BookmarkRepository interface {
	Create(ctx context.Context, bookmark *models.Bookmark, tags []string) error
	GetByID(ctx context.Context, id int) (*models.Bookmark, error)
	List(ctx context.Context, filters *ListFilters) ([]models.Bookmark, error)
	Search(ctx context.Context, query string, limit int) ([]models.Bookmark, error)
	Update(ctx context.Context, bookmark *models.Bookmark) error
	UpdateTags(ctx context.Context, bookmarkID int, tags []string) error
	AddTags(ctx context.Context, bookmarkID int, tags []string) error
	RemoveTags(ctx context.Context, bookmarkID int, tags []string) error
	Delete(ctx context.Context, id int) error
}

// TagRepository defines operations for tag management
type TagRepository interface {
	GetAll(ctx context.Context) ([]models.Tag, error)
	GetAllWithCounts(ctx context.Context) ([]TagCount, error)
	GetByName(ctx context.Context, name string) (*models.Tag, error)
	Create(ctx context.Context, name string) (*models.Tag, error)
	GetOrCreate(ctx context.Context, name string) (*models.Tag, error)
	Rename(ctx context.Context, oldName, newName string) error
	DeleteOrphaned(ctx context.Context) (int64, error)
	GetBookmarkTags(ctx context.Context, bookmarkID int) ([]models.Tag, error)
}

// ListFilters represents filters for listing bookmarks
type ListFilters struct {
	Tag        string
	PublicOnly bool
	Limit      int
	Offset     int
}

// TagCount represents a tag with its usage count
type TagCount struct {
	Tag   models.Tag
	Count int
}
