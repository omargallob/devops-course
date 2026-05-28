// Package database provides PostgreSQL connectivity, GORM model definitions,
// and automatic schema migrations for the devops-course backend.
package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to the database at dsn and runs AutoMigrate for all models.
// It retries the connection up to maxRetries times with exponential backoff,
// which is useful when the database container is still starting.
func Open(dsn string, log *slog.Logger) (*gorm.DB, error) {
	const maxRetries = 5

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	var db *gorm.DB
	var err error

	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn), gormCfg)
		if err == nil {
			break
		}
		wait := time.Duration(1<<uint(i)) * time.Second
		log.Warn("database connection failed, retrying", "attempt", i+1, "wait", wait, "error", err)
		time.Sleep(wait)
	}
	if err != nil {
		return nil, fmt.Errorf("database connection failed after %d attempts: %w", maxRetries, err)
	}

	log.Info("database connected")

	// AutoMigrate creates/updates tables to match the model structs.
	if err := db.AutoMigrate(
		&User{},
		&Session{},
		&LessonProgress{},
		&ExerciseSubmission{},
	); err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	log.Info("database migrations applied")
	return db, nil
}
