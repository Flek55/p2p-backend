package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
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

func (s *Service) Login(ctx context.Context, email, password, userAgent, ip string) (accessToken, refreshToken string, err error) {
	user, err := FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err = s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("access token generation failed: %w", err)
	}

	refreshToken, err = GenerateRandomString(64)
	if err != nil {
		return "", "", fmt.Errorf("refresh token generation failed: %w", err)
	}

	session := &Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ip,
		ExpiresAt:    time.Now().Add(time.Hour * 2),
		CreatedAt:    time.Now(),
	}

	if err := CreateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("session creation failed: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(s.accessExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
}
