package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
}

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
