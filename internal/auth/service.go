package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type Service struct {
	jwtSecret    []byte
	accessExpiry time.Duration
}

func CreateService(jwtSecret string, accessExpiry time.Duration) *Service {
	return &Service{jwtSecret: []byte(jwtSecret), accessExpiry: accessExpiry}
}

func (s *Service) Register(ctx context.Context, email, password string) error {
	if _, err := FindUserByEmail(ctx, email); err == nil {
		return ErrUserExists
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := CreateUser(ctx, &user); err != nil {
		return fmt.Errorf("user creation failed: %w", err)
	}

	return nil
}

func Login(c *gin.Context) {

}

func Logout(c *gin.Context) {

}
