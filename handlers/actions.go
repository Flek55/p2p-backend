package handlers

import (
	"net/http"
	"os"
	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

var service *auth.Service

func init() {
	service = auth.CreateService(os.Getenv("jwt_secret"), 300)
}

type LoginAndRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(c *gin.Context) {
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

func RigisterUser(c *gin.Context) {
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
