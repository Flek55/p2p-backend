package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type CustomClaims struct {
	UserID uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

func CreateService(jwtSecret string, accessExpirySeconds int) *Service {
	return &Service{
		jwtSecret:    []byte(jwtSecret),
		accessExpiry: time.Duration(accessExpirySeconds) * time.Second,
	}
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
		ExpiresAt:    time.Now().Add(time.Second * s.accessExpiry),
		CreatedAt:    time.Now(),
	}

	if err := CreateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("session creation failed: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) GenerateAccessToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(s.accessExpiry)
	log.Printf("Generating token for user %s, expires at: %s", userID, expirationTime.Format(time.RFC3339))

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Printf("Valid token for user: %s", claims.UserID)
	} else {
		return nil, errors.New("invalid token claims")
	}

	return token, nil
}

func (s *Service) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	return userID, nil
}
