package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a course participant, authenticated via GitHub OAuth.
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GitHubID    int64     `gorm:"uniqueIndex;not null"`
	Username    string    `gorm:"uniqueIndex;not null"`
	DisplayName string
	AvatarURL   string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Session stores a user's authenticated session (JWT reference, for revocation).
type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	TokenHash string    `gorm:"uniqueIndex;not null"` // SHA-256 hash of the JWT
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// LessonProgress tracks a user's progress through a lesson.
type LessonProgress struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_lesson"`
	User        User       `gorm:"foreignKey:UserID"`
	LessonSlug  string     `gorm:"not null;uniqueIndex:idx_user_lesson"`
	ModuleSlug  string     `gorm:"not null;index"`
	Status      string     `gorm:"not null;default:'not_started'"` // not_started | in_progress | completed
	CompletedAt *time.Time // nullable
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ExerciseSubmission records a user's code submission for an exercise.
type ExerciseSubmission struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	User         User      `gorm:"foreignKey:UserID"`
	ExerciseSlug string    `gorm:"not null;index"`
	Code         string    `gorm:"type:text;not null"`
	Passed       bool      `gorm:"not null;default:false"`
	Output       string    `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
