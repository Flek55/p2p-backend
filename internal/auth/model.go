package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
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
