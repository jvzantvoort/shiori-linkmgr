package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// Config holds database configuration
type Config struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	MaxConnections int
}

// NewDB creates a new database connection
func NewDB(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	if cfg.MaxConnections > 0 {
		conn.SetMaxOpenConns(cfg.MaxConnections)
	} else {
		conn.SetMaxOpenConns(10)
	}
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping verifies database connectivity
func (db *DB) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.conn.BeginTx(ctx, nil)
}

// NewBookmarkRepository creates a bookmark repository
func (db *DB) NewBookmarkRepository() BookmarkRepository {
	return &bookmarkRepo{db: db.conn}
}

// NewTagRepository creates a tag repository
func (db *DB) NewTagRepository() TagRepository {
	return &tagRepo{db: db.conn}
}
