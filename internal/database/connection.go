package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB wraps the sqlx.DB connection
type DB struct {
	*sqlx.DB
	driver string
}

// Connect creates a new database connection
func Connect(driver, dsn string) (*DB, error) {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{
		DB:     db,
		driver: driver,
	}, nil
}

// Driver returns the database driver name
func (db *DB) Driver() string {
	return db.driver
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
