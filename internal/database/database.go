package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM database connection
type DB struct {
	*gorm.DB
}

// New creates a new database connection and runs auto-migrations
func New(dbPath string) (*DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	db := &DB{DB: gormDB}

	// Disable foreign keys temporarily for migration
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return nil, fmt.Errorf("failed to disable foreign keys: %w", err)
	}

	// Run auto-migrations
	if err := db.autoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Re-enable foreign keys
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// autoMigrate runs GORM auto-migrations for all models
func (db *DB) autoMigrate() error {
	return db.AutoMigrate(
		&Site{},
		&Hand{},
		&Player{},
		&Action{},
	)
}
