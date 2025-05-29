package main

import (
	"time"

	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	auth.InitDb()

	exUser := auth.User{
		ID:           uuid.New(),
		Email:        "pipka",
		PasswordHash: "hhhhh",
		CreatedAt:    time.Now(),
	}

	if err := auth.CreateUser(&exUser); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Run()
}
