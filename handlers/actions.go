package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/Flek55/p2p-backend/internal/signaling"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func WebSocketHandler(s *signaling.Server, authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get room ID from query
		roomID := c.Query("room")
		if roomID == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("room ID required"))
			return
		}

		// Authenticate user
		token := c.Query("token")
		if token == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token required"))
			return
		}

		userID, err := validateToken(authService, token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Upgrade to WebSocket
		conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		signaling.NewClient(conn, roomID, userID.String(), s)
	}
}

func validateToken(service *auth.Service, token string) (uuid.UUID, error) {
	jwtToken, err := service.ValidateAccessToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	return userID, nil
}

var service *auth.Service

func init() {
	service = auth.CreateService(os.Getenv("jwt_secret"), 300)
}

type LoginAndRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginAndRegisterRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		userAgent := c.Request.UserAgent()
		ip := c.ClientIP()
		ctx := c.Request.Context()

		accessToken, refreshToken, err := service.Login(ctx, req.Email, req.Password, userAgent, ip)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Login successful",
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func RegisterUser(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginAndRegisterRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		err := service.Register(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
		})
	}
}
