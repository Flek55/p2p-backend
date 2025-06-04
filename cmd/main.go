package main

import (
	"os"
	"time"

	"github.com/Flek55/p2p-backend/handlers"
	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/Flek55/p2p-backend/internal/signaling"
	"github.com/gin-gonic/gin"
)

func main() {
	jwtSecret := os.Getenv("jwt_secret")
	auth.InitDb()
	// Create auth service

	accessExpiry := 300 * time.Second
	authService := auth.CreateService(jwtSecret, int(accessExpiry))

	// Create signaling server
	signalingServer := signaling.NewServer()

	r := gin.Default()

	// Existing routes
	r.POST("/login", handlers.LoginUser(authService))
	r.POST("/register", handlers.RegisterUser(authService))

	// New WebSocket route
	r.GET("/ws", handlers.WebSocketHandler(signalingServer, authService))

	r.Run(":8080")
}
